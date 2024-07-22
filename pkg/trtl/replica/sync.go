package replica

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/rotationalio/honu/replica"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/jitter"
	prom "github.com/trisacrypto/directory/pkg/trtl/metrics"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

//===========================================================================
// AntiEntropy (initiator/client-side) Methods
//===========================================================================

// AntiEntropy is a service that periodically selects a remote peer to synchronize with
// via bilateral anti-entropy using the Gossip service. Jitter is applied to the
// interval between anti-entropy synchronizations to ensure that message traffic isn't
// bursty to disrupt normal messages to the GDS service.
//
// The AntiEntropy background routine accepts a stop channel that can be used to stop
// the routine before the process shuts down. This is primarily used in tests, but is
// also used for graceful shutdown of the anti-entropy service.
func (r *Service) AntiEntropy(stop chan struct{}) {
	// If Replica is not enabled, do not start the AntiEntropy routine
	if !r.conf.Enabled {
		log.Info().Msg("anti-entropy not enabled")
		return
	}

	// Create the anti-entropy ticker and store the channel for shutdown
	ticker := jitter.New(r.conf.GossipInterval, r.conf.GossipSigma)
	r.aestop = stop

	// Log the start of the anti-entropy routine
	log.Info().
		Dur("interval", r.conf.GossipInterval).
		Dur("sigma", r.conf.GossipSigma).
		Msg("anti-entropy routine started")

		// Run anti-entropy at a stochastic interval
		// NOTE: bayou is the name of the original academic system that described anti-entropy
bayou:
	for {
		// Block until next tick or until stop signal is received
		select {
		case <-stop:
			log.Info().Msg("stopping anti-entropy service")
			break bayou
		case <-ticker.C:
		}

		// Randomly select a remote peer to synchronize with, continuing if we cannot
		// select a peer or no remote peers exist yet.
		var peer *peers.Peer
		if peer = r.SelectPeer(context.Background()); peer == nil {
			log.Debug().Msg("no remote peer available, skipping synchronization")
			continue bayou
		}

		// Create a logctx with the peer information for future logging
		logctx := sentry.With(nil).
			Dict("peer", sentry.Dict().Uint64("id", peer.Id).Str("addr", peer.Addr).Str("name", peer.Name)).
			Str("service", "anti-entropy").
			Bool("initiator", true)

		// Perform the anti-entropy synchronization session with the remote peer.
		if err := r.AntiEntropySync(peer, logctx); err != nil {
			logctx.Warn().Err(err).Msg("anti-entropy synchronization was unsuccessful")
		}

		// Update prometheus metrics
		prom.PmAESyncs.WithLabelValues(peer.Name, peer.Region, "initiator").Inc()
	}
}

// SelectPeer randomly to perform anti-entropy with, ensuring that the current replica
// is not selected if it is stored in the database. If a peer cannot be selected, then
// nil is returned. This method handles logging.
func (r *Service) SelectPeer(ctx context.Context) (peer *peers.Peer) {
	// Create a list of keys to select from so we don't unmarshal all peers in the
	// database in the case where we have a very large number of peers.
	keys := make([][]byte, 0)
	iter, err := r.db.Iter(nil, options.WithNamespace(wire.NamespaceReplicas))
	if err != nil {
		sentry.Error(ctx).Err(err).Msg("could not fetch peers from database")
		return nil
	}
	defer iter.Release()

	for iter.Next() {
		// Fetch the key and copy it into a new byte slice; otherwise on the next
		// iteration the pointer to the key byte slice will reference the next key.
		src := iter.Key()
		dst := make([]byte, len(src))
		copy(dst, src)
		keys = append(keys, dst)
	}

	if err = iter.Error(); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not iterate over peers in the database")
		return nil
	}

	// If we have at least 1 peer in the database, randomly select one. However, the
	// database probably contains a peer record that describes the current replica,
	// which we don't want to select - given a database of at least 2 items, trying 10
	// times to select a peer that is not self balances the amount of unmarshalling work
	// we do with the likelihood that we will probably get a peer that is not current,
	// given a large enough number of peers in the database.
	if len(keys) > 0 {
		// 10 attempts to select a random peer that is not the current peer.
		for i := 0; i < 10; i++ {
			var key, data []byte
			key = keys[rand.Intn(len(keys))]
			if data, err = r.db.Get(key, options.WithNamespace(wire.NamespaceReplicas)); err != nil {
				sentry.Warn(ctx).Str("key", string(key)).Err(err).Msg("could not fetch peer from the database")
				continue
			}

			peer = new(peers.Peer)
			if err = proto.Unmarshal(data, peer); err != nil {
				sentry.Warn(ctx).Str("key", string(key)).Err(err).Msg("could not unmarshal peer from database")
				continue
			}

			if peer.Id != r.conf.PID {
				return peer
			}
		}

		sentry.Warn(ctx).Int("nPeers", len(keys)).Msg("could not select peer after 10 attempts")
		return nil
	}

	sentry.Warn(ctx).Msg("database does not contain any peers")
	return nil
}

// AntiEntropySync performs bilateral anti-entropy with the specified remote peer using
// the streaming Gossip RPC. This method initiates the Gossip stream with the remote
// peer, exiting if it cannot connect to the replica (e.g. this method acts as the
// client in an anti-entropy session).
//
// The sync method for the initiator has two phases. In the first phase, the initiator
// loops over all objects in its local database and sends check requests to the remote,
// collecting all repair messages sent back from the remote (sometimes this is referred
// to as the pull phase of bilateral anti-entropy). In the second phase, the initiator
// waits for check messages from the remote and returns any objects that the remote
// requests (the push phase of bilateral anti-entropy).
//
// Both phases and the sending of messages are run in their own go routines, so 4 go
// routines are operating on the initiator side to handle the sync. The first phase go
// routine ends when it finishes looping over its database, the second phase go routine
// is also the recv go routine so it starts shortly after the first phase go routine and
// runs concurrently with it. The second phase ends when it receives complete from the
// remote. The send go routine ends when there are no more messages on its channel. Once
// all go routines are completed the initiator closes the channel, ending the
// synchronization between the initiator and the remote.
func (r *Service) AntiEntropySync(peer *peers.Peer, logctx *sentry.Logger) (err error) {
	// Start a timer to track latency
	start := time.Now()

	// Create a context with a timeout that is sooner than 95% of the timeouts selected
	// by the normally distributed jittered interval, to ensure anti-entropy gossip
	// sessions do not span multiple anti-entropy intervals.
	timeout := r.conf.GossipInterval - (2 * r.conf.GossipSigma)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Dial the remote peer and establish a connection
	var cc *grpc.ClientConn
	if cc, err = r.connect(peer); err != nil {
		return err
	}
	defer cc.Close()

	// Initiate the Gossip Stream
	client := replica.NewReplicationClient(cc)
	var stream replica.Replication_GossipClient
	if stream, err = client.Gossip(ctx); err != nil {
		return fmt.Errorf("could not connect gossip stream: %s", err)
	}

	// Report successful connection
	logctx.Debug().Msg("dialed remote peer and connected to the gossip stream")

	// Create a stream sender to ensure that both phase1 and phase2 go routines can send
	// messages concurrently without violating the grpc semantic that only one go
	// routine can send messages on the stream.
	// NOTE: newStreamSender calls wg.Add(1) and runs the sender go routine.
	wg := new(sync.WaitGroup)
	sender := newStreamSender(wg, logctx, stream)

	// Start phase 1: loop over all objects in the local database and send check
	// requests to the remote replica. This is also called the "pull" phase, since we're
	// asking the remote for its objects that are later than our own, e.g. pulling the
	// objects to this replica from the remote. Phase 1 ends when we've completed
	// looping over the local database.
	wg.Add(1)
	go r.initiatorPhase1(ctx, wg, logctx, sender)

	// Start phase 2: this phase is concurrent with phase 1 since it listens for and
	// responds to all messages from the remote replica. This is also called the "push"
	// phase, since we're pushing objects back to the remote replica. This phase ends
	// when the remote replica sends a COMPLETE message, meaning it is done sending
	// messages. At that point, we will no longer send any messages so this phase will
	// close the sender go routine, which will stop when all messages have been sent.
	wg.Add(1)
	go r.initiatorPhase2(ctx, wg, logctx, sender, stream)

	// Wait for the initiatorPhase1, initiatorPhase2, and sender anti-entropy routines
	wg.Wait()

	// Close the stream gracefully and cleanup the stream connections
	if err = stream.CloseSend(); err != nil {
		return fmt.Errorf("could not close gossip stream gracefully: %s", err)
	}

	if err = cc.Close(); err != nil {
		return fmt.Errorf("could not close the client connection correctly: %s", err)
	}

	// Compute latency in milliseconds
	// NOTE: we're only tracking latency for successful AE sessions
	latency := float64(time.Since(start)/1000) / 1000.0
	prom.PmAESyncLatency.WithLabelValues(peer.Name, peer.Region).Observe(latency)

	// Anti-entropy session complete
	return nil
}

// Connect to a remote peer using mTLS credentials or in insecure mode as necessary.
// This method blocks until the connection has been established to prevent any
// anti-entropy work from happening until we know the remote peer is live.
func (r *Service) connect(peer *peers.Peer) (cc *grpc.ClientConn, err error) {
	// Create the base dial options - ensure blocking
	opts := make([]grpc.DialOption, 0, 2)

	// Add mTLS credentials if required
	if r.mtls.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		var certPool *x509.CertPool
		if certPool, err = r.mtls.GetCertPool(); err != nil {
			return nil, fmt.Errorf("could not get cert pool: %v", err)
		}

		var cert tls.Certificate
		if cert, err = r.mtls.GetCert(); err != nil {
			return nil, fmt.Errorf("could not get cert: %v", err)
		}

		var u *url.URL
		if u, err = url.Parse(peer.Addr); err != nil {
			return nil, fmt.Errorf("could not parse %q: %v", peer.Addr, err)
		}

		conf := &tls.Config{
			ServerName:   u.Host,
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(conf)))
	}

	// Dial the remote peer and establish a connection
	return grpc.NewClient(peer.Addr, opts...)
}

// initiatorPhase1 is the go routine that starts the anti-entropy synchronization
// between the initiator replica (run by AntiEntropySync) and the remote replica
// (handled by Gossip). In this phase, the initiator loops over all objects in its local
// database (possibly modified by a sampling methodology) and sends CHECK requests to
// the remote peer. After looping through all objects in the database it sends a
// COMPLETE message to the remote, allowing it to begin its phase 2. This phase is the
// initiators anti-entropy "pull" component of bilateral anti-entropy, since it is
// asking the remote replica to send its later version.
//
// Note that this go routine does not handle any of the replies from the remote replica,
// all replies are handled in initiatorPhase2 whether they are replies to phase1 or
// messages sent in the remote's phase2.
func (r *Service) initiatorPhase1(ctx context.Context, wg *sync.WaitGroup, logctx *sentry.Logger, sender *streamSender) {
	// Start a timer to track latency
	start := time.Now()

	// Ensure that this routine signals when it exits
	defer wg.Done()
	logctx.Trace().Msg("starting initiator phase 1")

	// Track how many namespaces and versions we attempt to synchronize for logging.
	var nNamespaces, nVersions uint64

	// Access the objects in the object-store by namespace
namespaces:
	for _, namespace := range r.replicatedNamespaces {
		iter, err := r.db.Iter(nil, options.WithNamespace(namespace), options.WithTombstones())
		if err != nil {
			logctx.Error().Err(err).Str("namespace", namespace).Msg("could not iterate over namespace")
			continue namespaces
		}
		logctx.Trace().Str("namespace", namespace).Msg("sending namespace")

	objects:
		for iter.Next() {
			// Check if the context is done, and if so, break
			select {
			case <-ctx.Done():
				break namespaces
			default:
			}

			// Load the object metadata without the data itself, otherwise anti-
			// entropy would exchange way more data than required, putting pressure
			// on pod memory and increasing our cloud bill.
			obj, err := iter.Object()
			if err != nil {
				logctx.Error().Err(err).
					Str("namespace", namespace).
					Str("key", b64e(iter.Key())).
					Msg("could not unmarshal honu metadata")
				continue objects
			}

			// Remove the data from the object and create the check message
			obj.Data = nil
			msg := &replica.Sync{
				Status: replica.Sync_CHECK,
				Object: obj,
			}

			// Ensure that we can send, otherwise stop.
			// NOTE: this is not fully safe, since there is no logical barrier
			// between this check and the send on the chanel, however the sender
			// keeps consuming messages to prevent race conditions.
			if ok := sender.Send(msg); !ok {
				// NOTE: breaking objects loop not the namespaces loop or returning to
				// ensure that all iters get released. This may cause some additional
				// Send calls to happen, but the sender will simply ignore them.
				break objects
			}
			nVersions++
		}

		if err = iter.Error(); err != nil {
			logctx.Error().Err(err).Str("namespace", namespace).Msg("could not iterate over namespace")
		}

		// Ensure the iterator is released, note even if break objects occurs, the
		// iterator should be released at this line of code since there is no return.
		iter.Release()
		nNamespaces++
	}

	// Send a sync complete message to let the remote know that the pull phase is
	// complete and they can start the push phase.
	sender.Send(&replica.Sync{Status: replica.Sync_COMPLETE})
	logctx.Trace().Msg("initiator phase 1 complete")

	// Compute latency in milliseconds
	latency := float64(time.Since(start)/1000) / 1000.0
	prom.PmAEPhase1Latency.WithLabelValues(r.conf.Name).Observe(latency)

	logctx.Debug().
		Uint64("versions", nVersions).
		Uint64("namespaces", nNamespaces).
		Msg("version vectors sent to remote peer")
}

// initiatorPhase2 starts right after initiatorPhase1 and runs as long as messages will
// be coming from the remote. This is the only initiator go routine that will receive
// messages, so an intermediate read routine is not necessary. This phase is also
// responsible for closing the sender go routine because the phase is finished when it
// gets a COMPLETE message from the remote. This phase handles incoming messages from
// the remote by responding to CHECK requests sending later versions to the remote (but
// ignoring if local versions are equal or earlier), and handling REPAIR and error.
func (r *Service) initiatorPhase2(ctx context.Context, wg *sync.WaitGroup, logctx *sentry.Logger, sender *streamSender, stream gossipStream) {
	// Ensure that this routine signals when it exits
	defer wg.Done()
	logctx.Trace().Msg("starting initiator phase 2")

	// Ensure that we close the sending channel when this routine exits to prevent
	// deadlocks if this phase ends prematurely (e.g. the timeout expires).
	defer sender.Close()

	// Track how many check versions, updates, and repairs we get in phase 2
	var versions, updates, repairs uint64
	var err error

gossip:
	for {
		// Check to make sure the deadline isn't over
		select {
		case <-ctx.Done():
			logctx.Debug().Msg("context canceled while trying to recv messages from remote")
			return
		default:
		}

		// Read the next message from the remote replica
		var sync *replica.Sync
		if sync, err = stream.Recv(); err != nil {
			if err != io.EOF {
				// If the error is not EOF then something has gone wrong.
				logctx.Error().Err(err).Msg("anti-entropy aborted early with recv error")
			}
			return
		}

		switch sync.Status {
		case replica.Sync_CHECK:
			// Check to see if this replica's version is later than the remote's, if
			// so, send our later version back to that replica.
			versions++
			var local *object.Object
			if local, err = r.db.Object(sync.Object.Key, options.WithNamespace(sync.Object.Namespace)); err != nil {
				logctx.Warn().Err(err).
					Str("namespace", sync.Object.Namespace).
					Str("key", b64e(sync.Object.Key)).
					Msg("failed check sync on initiator: could not fetch object meta")

				// Even if the error is not found, since this routine is the
				// initiator, it will return an error for a check request because
				// we do not want to bounce check requests back and forth. E.g.
				// check requests only initiate from the "pull" routine (phase 1
				// locally, and phase 2 on the remote).
				sync.Status = replica.Sync_ERROR
				sync.Error = err.Error()
				sender.Send(sync)
				continue gossip
			}

			// Check if the local version is later; because this is the initiating
			// routine, do nothing if the remote is a later version.
			if local.Version.IsLater(sync.Object.Version) {
				// Send the local object back to the remote
				sync.Status = replica.Sync_REPAIR
				sync.Object = local
				sync.Error = ""
				if ok := sender.Send(sync); ok {
					updates++
				}
			}

		case replica.Sync_REPAIR:
			// Receiving a repaired object, check if it's still later than our local
			// replica's and if so, save it to disk with an update instead of a put.
			//
			// NOTE: honu.Update performs the version checking in a transaction.

			var updateType honu.UpdateType
			if updateType, err = r.db.Update(sync.Object, options.WithNamespace(sync.Object.Namespace)); err != nil {
				logctx.Warn().Err(err).
					Str("namespace", sync.Object.Namespace).
					Str("key", b64e(sync.Object.Key)).
					Msg("could not update object from remote peer")
				continue gossip
			}
			// Log update type in prometheus metrics.
			switch updateType {
			case honu.UpdateStomp:
				prom.PmAEStomps.WithLabelValues(r.conf.Name, r.conf.Region).Inc()
			case honu.UpdateSkip:
				prom.PmAESkips.WithLabelValues(r.conf.Name, r.conf.Region).Inc()
			}
			repairs++

		case replica.Sync_ERROR:
			// Something went wrong on the remote, log and continue
			if sync.Object != nil {
				logctx.Warn().Str("error", sync.Error).
					Str("key", b64e(sync.Object.Key)).
					Str("namespace", sync.Object.Namespace).
					Msg("a replication error occurred")
			} else {
				logctx.Warn().Str("error", sync.Error).Msg("a replication error occurred")
			}

		case replica.Sync_COMPLETE:
			// The remote replica is done synchronizing and since we are the
			// initiating replica, we can safely quit receiving.
			logctx.Debug().Uint64("versions", versions).Msg("received version vectors from remote peer")
			if updates > 0 || repairs > 0 {
				r.synchronizedNow()
				logctx.Info().
					Uint64("local_repairs", repairs).
					Uint64("remote_updates", updates).
					Uint64("versions", versions).
					Msg("anti-entropy synchronization complete")
			} else {
				logctx.Debug().Msg("anti-entropy complete with no synchronization")
			}

			// When we receive the COMPLETE message from the remote replica, we're done:
			// exit the for loop and close the sender (via the defer above). Once all
			// messages are sent, we can close the stream and finish.
			logctx.Trace().Msg("initiator phase 2 complete")

			// Update Prometheus metrics
			prom.PmAEVersions.WithLabelValues(r.conf.Name, r.conf.Region, "initiator").Observe(float64(versions))
			prom.PmAEUpdates.WithLabelValues(r.conf.Name, r.conf.Region, "initiator").Observe(float64(updates))
			prom.PmAERepairs.WithLabelValues(r.conf.Name, r.conf.Region, "initiator").Observe(float64(repairs))

			return

		default:
			logctx.Error().Str("status", sync.Status.String()).Msg("unhandled sync status")
		}
	}
}
