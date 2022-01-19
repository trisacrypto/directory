package trtl

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rotationalio/honu"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/rotationalio/honu/replica"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/jitter"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// A ReplicaService manages anti-entropy replication between peers.
type ReplicaService struct {
	sync.RWMutex
	replica.UnimplementedReplicationServer
	parent       *Server
	conf         config.ReplicaConfig
	db           *honu.DB
	aestop       chan struct{}
	synchronized time.Time
}

func NewReplicaService(s *Server) (*ReplicaService, error) {
	return &ReplicaService{
		parent: s,
		conf:   s.conf.Replica,
		db:     s.db,
	}, nil
}

// AntiEntropy is a service that periodically selects a remote peer to synchronize with
// via bilateral anti-entropy using the Gossip service. Jitter is applied to the
// interval between anti-entropy synchronizations to ensure that message traffic isn't
// bursty to disrupt normal messages to the GDS service.
//
// The AntiEntropy background routine accepts a stop channel that can be used to stop
// the routine before the process shuts down. This is primarily used in tests, but is
// also used for graceful shutdown of the anti-entropy service.
// TODO: this background routine is currently untested.
func (r *ReplicaService) AntiEntropy(stop chan struct{}) {
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
		if peer = r.SelectPeer(); peer == nil {
			log.Debug().Msg("no remote peer available, skipping synchronization")
			continue bayou
		}

		// Create a logctx with the peer information for future logging
		logctx := log.With().
			Dict("peer", zerolog.Dict().Uint64("id", peer.Id).Str("addr", peer.Addr).Str("name", peer.Name)).
			Str("service", "anti-entropy").
			Bool("initiator", true).
			Logger()

		// Perform the anti-entropy synchronization session with the remote peer.
		if err := r.AntiEntropySync(peer, logctx); err != nil {
			logctx.Warn().Err(err).Msg("anti-entropy synchronization was unsuccessful")
		}

		// Update prometheus metrics
		pmAESyncs.WithLabelValues(peer.Name, peer.Region).Inc()
	}
}

// AntiEntropySync performs bilateral anti-entropy with the specified remote peer using
// the streaming Gossip RPC. This method initiates the Gossip stream with the remote
// peer, exiting if it cannot connect to the replica. In the pull phase, this method
// sends check sync messages for all objects stored locally; the remote replica responds
// with repairs. Then in the push phase, the method waits until all requested remote
// repairs are complete before exiting.
func (r *ReplicaService) AntiEntropySync(peer *peers.Peer, log zerolog.Logger) (err error) {
	// Start a timer to track latency
	start := time.Now()

	// Create a context with a timeout that is sooner than 95% of the timeouts selected
	// by the normally distributed jittered interval, to ensure anti-entropy gossip
	// sessions do not span multiple anti-entropy intervals.
	timeout := r.conf.GossipInterval - (2 * r.conf.GossipSigma)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Dial the remote peer and establish a connection
	// TODO: add client-side mTLS here.
	var cc *grpc.ClientConn
	if cc, err = grpc.DialContext(ctx, peer.Addr, grpc.WithInsecure(), grpc.WithBlock()); err != nil {
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
	log.Debug().Msg("dialed remote peer and connected to the gossip stream")

	// Kick off two go routines - one to send messages on and one to recv them on, this
	// is necessary because gRPC send and recv are blocking operations that can only be
	// used by one go routine, e.g. they are not thread safe. These go routines place
	// channels between the send and recv operations so that multiple go routines can
	// send and recv concurrently.
	// TODO: standardize the functionality in a utility package.
	var stop uint32
	var wg sync.WaitGroup
	send := make(chan *replica.Sync, 8)
	recv := make(chan *replica.Sync, 8)

	// Kick off the stream sender: will continue to send messages until the send channel is closed
	go func(msgs <-chan *replica.Sync) {
		for msg := range msgs {
			log.Trace().Str("status", msg.Status.String()).Msg("sending sync")
			if err := stream.Send(msg); err != nil {
				if err != io.EOF {
					log.Error().Err(err).Msg("could not send gossip message to remote peer")
				}

				// Let external go routines know that they can no longer send on this channel.
				atomic.AddUint32(&stop, 1)
				break
			}
		}

		// Continue to consume messages until the msgs channel is closed to ensure that
		// go routines that do not check if they can send messages don't get blocked if
		// the message channel is filled up. This won't do anything with the messages.
		// NOTE: ranging over a closed channel immediately returns so even if the msgs
		// channel was closed in the first loop, this second loop will not execute nor
		// will it panic.
		for msg := range msgs {
			log.Trace().Str("status", msg.Status.String()).Msg("dropped sync message (not sent)")
		}
	}(send)

	// Kick off the stream receiver: will continue to recv messages until the timeout or
	// the grpc channel is closed. Note that this go routine closes the recv channel.
	go func(msgs chan<- *replica.Sync) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if msg, err := stream.Recv(); err != nil {
				if err != io.EOF {
					log.Error().Err(err).Msg("could not recv gossip message from remote peer")
				}

				// Stop any external go routines from recv on this channel
				close(msgs)
				return
			} else {
				msgs <- msg
			}
		}
	}(recv)

	// Pull sends check object sync requests to the remote server
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Track how many namespaces and versions we attempt to synchronize for logging.
		var namespaces, versions uint64

		// Access the objects in the object-store by namespace
	namespaces:
		for _, namespace := range replicatedNamespaces {
			iter, err := r.db.Iter(nil, options.WithNamespace(namespace))
			if err != nil {
				log.Error().Err(err).Str("namespace", namespace).Msg("could not iterate over namespace")
				continue namespaces
			}

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
					log.Error().Err(err).
						Str("namespace", namespace).
						Str("key", b64e(iter.Key())).
						Msg("could not unmarshal honu metadata")
					continue objects
				}

				// Remove the data from the object
				obj.Data = nil

				// Ensure that we can send, otherwise stop.
				// TODO: this is not fully safe, since there is no logical barrier
				// between this check and the send on the chanel. This needs to be
				// improved but is ok for now.
				if atomic.LoadUint32(&stop) > 0 {
					break objects
				}

				// Send the synchronization request
				send <- &replica.Sync{
					Status: replica.Sync_CHECK,
					Object: obj,
				}
				versions++
			}

			if err = iter.Error(); err != nil {
				log.Error().Err(err).Str("namespace", namespace).Msg("could not iterate over namespace")
			}

			iter.Release()
			namespaces++
		}

		// Send a sync complete message to let the remote know that the pull phase is
		// complete and they can start the push phase.
		if atomic.LoadUint32(&stop) == 0 {
			send <- &replica.Sync{Status: replica.Sync_COMPLETE}
		}

		log.Debug().
			Uint64("versions", versions).
			Uint64("namespaces", namespaces).
			Msg("sending version vectors to remote peer")
	}()

	// Push listens for sync requests from the remote server and makes any necessary repairs
	wg.Add(1)
	go func() {
		defer wg.Done()

		var versions, updates, repairs uint64

	gossip:
		for sync := range recv {
			switch sync.Status {
			case replica.Sync_CHECK:
				// Check to see if this replica's version is later than the remote's, if
				// so, send our later version back to that replica.
				versions++
				var local *object.Object
				if local, err = r.db.Object(sync.Object.Key, options.WithNamespace(sync.Object.Namespace)); err != nil {
					log.Warn().Err(err).
						Str("namespace", sync.Object.Namespace).
						Str("key", b64e(sync.Object.Key)).
						Msg("failed check sync on initiator: could not fetch object meta")

					// Even if the error is not found, since this routine is the
					// initiator, it will return an error for a check request because
					// we do not want to bounce check requests back and forth. E.g.
					// check requests only initiate from the Pull routine.
					sync.Status = replica.Sync_ERROR
					sync.Error = err.Error()
					if atomic.LoadUint32(&stop) == 0 {
						send <- sync
					}

					continue gossip
				}

				// Check if the local version is later; because this is the initiating
				// routine, do nothing if the remote is a later version.
				if local.Version.IsLater(sync.Object.Version) {
					// Send the local object back to the remote
					sync.Status = replica.Sync_REPAIR
					sync.Object = local
					sync.Error = ""
					if atomic.LoadUint32(&stop) == 0 {
						send <- sync
						updates++
					}
					// Update Prometheus metrics
					pmAEStomps.WithLabelValues(r.conf.Name, r.conf.Region).Inc()
				}

			case replica.Sync_REPAIR:
				// Receiving a repaired object, check if it's still later than our local
				// replica's and if so, save it to disk with an update instead of a put.
				// TODO: should we confirm that the incoming version is still later than
				// our local version? E.g. is it possible that the local version has
				// been updated since we sent the check request?
				if err = r.db.Update(sync.Object, options.WithNamespace(sync.Object.Namespace)); err != nil {
					log.Error().Err(err).
						Str("namespace", sync.Object.Namespace).
						Str("key", b64e(sync.Object.Key)).
						Msg("could not update object from remote peer")
					continue gossip
				}
				repairs++

			case replica.Sync_ERROR:
				// Something went wrong on the remote, log and continue
				if sync.Object != nil {
					log.Warn().Str("error", sync.Error).
						Str("key", b64e(sync.Object.Key)).
						Str("namespace", sync.Object.Namespace).
						Msg("a replication error occurred")
				} else {
					log.Warn().Str("error", sync.Error).Msg("a replication error occurred")
				}

			case replica.Sync_COMPLETE:
				// The remote replica is done synchronizing and since we are the
				// initiating replica, we can safely quit receiving.
				log.Debug().Uint64("versions", versions).Msg("received version vectors from remote peer")
				if updates > 0 || repairs > 0 {
					r.synchronizedNow()
					log.Info().
						Uint64("local_repairs", repairs).
						Uint64("remote_updates", updates).
						Msg("anti-entropy synchronization complete")
				} else {
					log.Debug().Msg("anti-entropy complete with no synchronization")
				}
				return
			default:
				log.Error().Str("status", sync.Status.String()).Msg("unhandled sync status")
			}
		}
	}()

	// Wait for anti-entropy routines to complete.
	wg.Wait()

	// Close the send stream since we will no longer be sending messages
	close(send)

	// TODO: wait for all messages to finish sending

	// Cleanup the stream and stream connections
	if err = stream.CloseSend(); err != nil {
		return fmt.Errorf("could not close gossip stream gracefully: %s", err)
	}

	if err = cc.Close(); err != nil {
		return fmt.Errorf("could not close the client connection correctly: %s", err)
	}

	// Compute latency in milliseconds
	// NOTE: we're only tracking latency for successful AE sessions
	latency := float64(time.Since(start)/1000) / 1000.0
	pmAESyncLatency.WithLabelValues(peer.Name).Observe(latency)

	// Anti-entropy session complete
	return nil
}

// SelectPeer randomly that is not self to perform anti-entropy with. If a peer
// cannot be selected, then nil is returned.
func (r *ReplicaService) SelectPeer() (peer *peers.Peer) {
	// Select a random peer that is not self to perform anti entropy with.
	keys := make([][]byte, 0)
	iter, err := r.db.Iter(nil, options.WithNamespace(NamespacePeers))
	if err != nil {
		log.Error().Err(err).Msg("could not fetch peers from database")
		return nil
	}
	defer iter.Release()

	for iter.Next() {
		keys = append(keys, iter.Key())
	}

	if err = iter.Error(); err != nil {
		log.Error().Err(err).Msg("could not iterate over peers in the database")
		return nil
	}

	// TODO: it appears that the keys are not correctly being constructed and that
	// there are duplicates in the slice.
	if len(keys) > 1 {
		// 10 attempts to select a random peer that is not self.
		for i := 0; i < 10; i++ {
			var key, data []byte
			key = keys[rand.Intn(len(keys))]
			if data, err = r.db.Get(key, options.WithNamespace(NamespacePeers)); err != nil {
				log.Warn().Str("key", string(key)).Err(err).Msg("could not fetch peer from the database")
				continue
			}

			peer = new(peers.Peer)
			if err = proto.Unmarshal(data, peer); err != nil {
				log.Warn().Str("key", string(key)).Err(err).Msg("could not unmarshal peer from database")
			}

			if peer.Id != r.conf.PID {
				return peer
			}
		}
		log.Warn().Int("nPeers", len(keys)).Msg("could not select peer after 10 attempts")
		return nil
	}

	log.Warn().Msg("could not select peer from the database")
	return nil
}

// Gossip implements biltateral anti-entropy: during a Gossip session the initiating
// replica pushes updates to the remote peer and pulls requested changes. Using
// bidirectional streaming, the initiating peer sends data-less sync messages with
// the versions of objects it stores locally. The remote replica then responds with
// data if its local version is later or sends a sync message back requesting the
// data from the initiating replica if its local version is earlier. No exchange
// occurs if both replicas have the same version. At the end of a gossip session,
// both replicas should have synchronized and have identical underlying data stores.
func (r *ReplicaService) Gossip(stream replica.Replication_GossipServer) (err error) {
	var once sync.Once
	var wg sync.WaitGroup
	ctx := stream.Context()

	var versions, updates, repairs uint64

	// Set up the log context for consistent logging
	// TODO: get remote peer information from mTLS context and add to logging.
	logctx := log.With().
		Str("service", "gossip").
		Bool("remote", true).
		Logger()

	// Create data structures for determining what objects are sync'd in phase 1
	// TODO: replace this map with a bloom filter
	var seen nsmap
	phase1 := true

	// Describe the phase 2 go routine
	phase2 := func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			logctx.Trace().Msg("starting phase 2")

			// Loop over all objects in all namespaces and determine what to push back
		namespaces:
			for _, namespace := range replicatedNamespaces {
				iter, err := r.db.Iter(nil, options.WithNamespace(namespace))
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

					// Check if we've already handled this object
					if seen.In(namespace, iter.Key()) {
						// This was already sent to us, so keep going
						continue objects
					}

					// If this key hasn't been seen then it is a new key local to this
					// replica that needs to be pushed back to the client replica.
					obj, err := iter.Object()
					if err != nil {
						logctx.Error().Err(err).
							Str("namespace", namespace).
							Str("key", b64e(iter.Key())).
							Msg("could not unmarshal honu metadata")
						continue objects
					}

					if err := stream.Send(&replica.Sync{Status: replica.Sync_REPAIR, Object: obj}); err != nil {
						logctx.Error().Err(err).Msg("could not send gossip message, prematurely quitting phase 2")
						iter.Release()
						return
					}

					atomic.AddUint64(&updates, 1)
					atomic.AddUint64(&versions, 1)
				}

				if err = iter.Error(); err != nil {
					logctx.Error().Err(err).Str("namespace", namespace).Msg("could not iterate over namespace")
				}
				iter.Release()
			}

			// Phase 2 is complete! Send a complete message and return
			if err := stream.Send(&replica.Sync{Status: replica.Sync_COMPLETE}); err != nil {
				logctx.Error().Err(err).Msg("could not send gossip complete message")
				return
			}
			logctx.Trace().Msg("phase 2 complete")
		}()
	}

	// Kick off a go routine to handle all incoming messages from the client.
	wg.Add(1)
	go func() {
		defer wg.Done()
		logctx.Debug().Msg("starting phase 1")

	gossip:
		for {
			// Check to make sure the deadline isn't over
			select {
			case <-ctx.Done():
				logctx.Debug().Bool("phase1", phase1).Msg("context canceled during gossip recv")
				return
			default:
			}

			// Read the next message from the initiating replica
			var sync *replica.Sync
			if sync, err = stream.Recv(); err != nil {
				if err != io.EOF {
					// If the error is not EOF then something has gone wrong. Otherwise
					// it means that the client closed the stream gracefully and is done
					// sending messages to the server.
					logctx.Error().Err(err).Msg("gossip aborted early with error")
				}
				return
			}

			// Handle the messages coming from the initiating replica
			switch sync.Status {
			case replica.Sync_CHECK:
				// If we're out of phase 1 - meaning that the client sent a complete
				// message, then we ignore check messages to prevent checks bouncing
				// back and forth from poorly behaved clients. This also ensures that
				// both phase1 and phase2 go routines are not sending messages on the
				// stream, which causes a grpc error.
				// TODO: should we create a send go routine anyway for even more safety?
				if !phase1 {
					logctx.Warn().Msg("received check message after phase 1 complete from client")
					continue gossip
				}

				// Mark the object as seen if we're in phase 1 to prevent duplication in phase 2
				seen.Add(sync.Object)

				// Increment the number of versions seen
				atomic.AddUint64(&versions, 1)

				// If we're in phase 1 - that means we're receiving version vectors from
				// the initiating replica, we should compare the incoming version to the
				// local version on this replica.
				var local *object.Object
				if local, err = r.db.Object(sync.Object.Key, options.WithNamespace(sync.Object.Namespace)); err != nil {
					if err == engine.ErrNotFound {
						// If this is a not found error, then this object exists on the
						// initiating replica, but not locally, so request a repair.
						// Note that we have to set the object version to zero to
						// indicate that we don't have a version of it, so the client
						// replica's version is later and it sends a repair message.
						sync.Object.Version = &object.VersionZero
						if err = stream.Send(sync); err != nil {
							logctx.Error().Err(err).Msg("could not send gossip message, prematurely quitting phase 1")
							return
						}
					} else {
						// This is an unhandled error, log and return error information.
						logctx.Warn().Err(err).
							Str("namespace", sync.Object.Namespace).
							Str("key", b64e(sync.Object.Key)).
							Msg("failed check sync on initiator: could not fetch object meta")

						// Even if the error is not found, since this routine is the
						// initiator, it will return an error for a check request because
						// we do not want to bounce check requests back and forth. E.g.
						// check requests only initiate from the Pull routine.
						sync.Status = replica.Sync_ERROR
						sync.Error = err.Error()
						if err = stream.Send(sync); err != nil {
							logctx.Error().Err(err).Msg("could not send gossip message, prematurely quitting phase 1")
							return
						}
					}
					continue gossip
				}

				// Check which version is later, local or remote
				switch {
				case local.Version.IsLater(sync.Object.Version):
					// Send the local object back to the initiating replica as a repair
					sync.Status = replica.Sync_REPAIR
					sync.Object = local
					sync.Error = ""
					if err = stream.Send(sync); err != nil {
						logctx.Error().Err(err).Msg("could not send gossip message, prematurely quitting phase 1")
						return
					}
					atomic.AddUint64(&updates, 1)
				case sync.Object.Version.IsLater(local.Version):
					// Send a check request back to the initiating replica to fetch its version
					if err = stream.Send(sync); err != nil {
						logctx.Error().Err(err).Msg("could not send gossip message, prematurely quitting phase 1")
						return
					}
				default:
					// The versions are equal, do nothing
				}
			case replica.Sync_REPAIR:
				// Receiving a repaired object, check if it's still later than our local
				// replica's and if so, save it to disk with an update instead of a put.
				// TODO: should we confirm that the incoming version is still later than
				// our local version? E.g. is it possible that the local version has
				// been updated since we sent the check request?
				if err = r.db.Update(sync.Object, options.WithNamespace(sync.Object.Namespace)); err != nil {
					logctx.Error().Err(err).
						Str("namespace", sync.Object.Namespace).
						Str("key", b64e(sync.Object.Key)).
						Msg("could not update object from client replica")
					continue gossip
				}
				atomic.AddUint64(&repairs, 1)
			case replica.Sync_ERROR:
				// Something went wrong on the client-side, log and continue
				if sync.Object != nil {
					log.Warn().Str("error", sync.Error).
						Str("key", b64e(sync.Object.Key)).
						Str("namespace", sync.Object.Namespace).
						Msg("a replication error occurred")
				} else {
					log.Warn().Str("error", sync.Error).Msg("a replication error occurred")
				}

			case replica.Sync_COMPLETE:
				// Phase 1 is complete! Begin Phase 2!
				// The phase1 bool causes this routine to stop saving seen messages and
				// the once sync ensures that multiple COMPLETE messages from the client
				// don't kick off more than one go routine.
				phase1 = false
				once.Do(phase2)
			default:
				logctx.Error().Str("status", sync.Status.String()).Msg("unhandled sync status")
			}
		}
	}()

	// Wait for all go routines to finish
	wg.Wait()

	// Log and complete gossip session
	log.Debug().Uint64("versions", versions).Msg("received version vectors from remote peer")
	if updates > 0 || repairs > 0 {
		r.synchronizedNow()
		log.Info().
			Uint64("local_repairs", repairs).
			Uint64("remote_updates", updates).
			Msg("anti-entropy synchronization complete")
	} else {
		log.Debug().Msg("anti-entropy complete with no synchronization")
	}

	// Update Prometheus metrics
	pmAEPushes.WithLabelValues(r.conf.Name, r.conf.Region).Observe(float64(updates))
	pmAEPulls.WithLabelValues(r.conf.Name, r.conf.Region).Observe(float64(repairs))
	pmAEPushVSPull.Add(float64(updates))
	pmAEPushVSPull.Sub(float64(repairs))

	return nil
}

func (r *ReplicaService) Shutdown() error {
	if r.aestop != nil {
		r.aestop <- struct{}{}
	}
	return nil
}

// Helper function to get the timestamp of last synchronization in a thread-safe manner
func (r *ReplicaService) lastSynchronization() string {
	r.RLock()
	defer r.RUnlock()
	if r.synchronized.IsZero() {
		return "never"
	}
	return r.synchronized.Format(time.RFC3339)
}

// Helper function to set the synchronized timestamp to now in a thread-safe manner
func (r *ReplicaService) synchronizedNow() {
	r.Lock()
	r.synchronized = time.Now()
	r.Unlock()
}
