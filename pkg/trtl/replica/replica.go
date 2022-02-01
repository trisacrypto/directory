package replica

import (
	"context"
	"io"
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
	"github.com/trisacrypto/directory/pkg/trtl/metrics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service manages anti-entropy replication between peers. It has two primary functions:
// an Anti-Entropy routine that periodically selects another peer to initiate
// anti-entropy with and a Gossip server method that responds to anti-entropy requests.
// The Service implements bi-lateral anti-entropy to increase consistency.
// NOTE: the name of this struct follows the convention of the other trtl services.
type Service struct {
	sync.RWMutex
	replica.UnimplementedReplicationServer
	conf                 config.ReplicaConfig
	mtls                 config.MTLSConfig
	db                   *honu.DB
	aestop               chan struct{}
	synchronized         time.Time
	replicatedNamespaces []string
}

// New creates a new replica.Service that is completely decoupled from the trtl.Server.
// This breaks the pattern of the PeersService, MetricsService, and TrtlService but
// allows replication to be completely encapsulated in a single package.
func New(conf config.Config, db *honu.DB, replicatedNamespaces []string) (*Service, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		conf:                 conf.Replica,
		mtls:                 conf.MTLS,
		db:                   db,
		aestop:               make(chan struct{}),
		replicatedNamespaces: replicatedNamespaces,
	}, nil
}

// Shutdown the replica server (stops the anti-entropy go-routine)
func (r *Service) Shutdown() error {
	// If anti-entropy is enabled, send a stop signal to it. Do not send the signal if
	// not enabled, otherwise Shutdown() will block forever and cause a deadlock.
	if r.conf.Enabled && r.aestop != nil {
		r.aestop <- struct{}{}
	}
	return nil
}

//===========================================================================
// Gossip (server-side) Methods
//===========================================================================

// Gossip is a server method that responds to anti-entropy requests.
// The initiating replica will engage `Gossip` to enable the remote/receiving
// replica to receive incoming version vectors for all objects in the initiating
// replica's Trtl store in phase one. The final step of phase one triggers phase
// two, when the remote replica responds with data if its local version is later.
// Concurrently with these phases, the remote sends a sync message back
// requesting data from the initiating replica if its local version is earlier.
func (r *Service) Gossip(stream replica.Replication_GossipServer) (err error) {
	// Set up the log context for consistent logging
	// TODO: get remote peer information from mTLS context and add to logging.
	logctx := log.With().
		Str("service", "gossip").
		Bool("remote", true).
		Logger()

	// If anti-entropy is not enabled on this replica, reject the request so that the
	// underlying database does not get modified. This method should not be served in
	// maintenance mode, so no need to check that state.
	if !r.conf.Enabled {
		logctx.Debug().Msg("rejecting gossip request: anti-entropy not enabled")
		return status.Error(codes.FailedPrecondition, "anti-entropy not enabled on remote")
	}

	// Create a stream sender to ensure that both phase1 and phase2 go routines can send
	// messages concurrently without violating the grpc semantic that only one go
	// routine can send messages on the stream.
	// NOTE: newStreamSender calls wg.Add(1) and runs the sender go routine.
	wg := new(sync.WaitGroup)
	sender := newStreamSender(wg, logctx, stream)

	// Start phase 1: receive object version vectors from the initiator. This go routine
	// will kick off phase 2: sending unchecked versions back to the initiator.
	wg.Add(1)
	go r.remotePhase1(stream.Context(), wg, logctx, stream, sender)

	// Wait for all go routines to finish
	wg.Wait()

	return nil
}

// remotePhase1 is the counterpart to initiatorPhase2 and runs as long as messages will
// be coming from the initiator. This is the only remote go routine that will receive
// messages, so an intermediate read routine is not necessary. This phase has two parts
// (so could be considered both phase 1 and phase 3). The difference is that the phase
// is allowed to send messages in the first part (phase 1) but cannot send messages in
// the second part (phase 3). In phase 1, CHECK messages from the initiator are
// responded to with either a mirror CHECK if the initiator version is later, a REPAIR
// if the remote is later, or nothing if they are equal. When a COMPLETE message is
// received from the initiator, phase1 ends and phase3 begins to handle REPAIR and ERROR
// messages that are still coming from the initiator but this routine will no longer
// send any messages. Receipt of the COMPLETE message from the initiator also kicks off
// the remotePhase2 go routine, ensuring it only runs once. This phase ends when the
// initiator closes the stream with CloseSend, ending gossip.
func (r *Service) remotePhase1(ctx context.Context, wg *sync.WaitGroup, log zerolog.Logger, stream replica.Replication_GossipServer, sender *streamSender) {
	defer wg.Done()
	log.Trace().Msg("starting phase 1")

	// TODO: replace the seen nsmap with a bloom filter or a more efficient data structure
	var (
		err      error     // error handling
		seen     nsmap     // track objects seen in phase 1
		once     sync.Once // ensure that only one phase2 routine executes
		phase3   bool      // check if we're in phase 1 or phase 3
		versions uint64    // number of versions received from initiator
		updates  uint64    // number of updates sent back to the initiator
		repairs  uint64    // number of repairs received from initiator
	)

	// When phase 3 is complete (or if phase 1 ends early) log anti-entropy
	defer func() {
		nUpdates := atomic.LoadUint64(&updates)
		rRepairs := atomic.LoadUint64(&repairs)
		if nUpdates > 0 || rRepairs > 0 {
			r.synchronizedNow()
			log.Info().
				Uint64("local_repairs", rRepairs).
				Uint64("remote_updates", nUpdates).
				Msg("anti-entropy synchronization complete")
		} else {
			log.Debug().Msg("anti-entropy complete with no synchronization")
		}
	}()

gossip:
	for {
		// Check to make sure the deadline isn't over.
		select {
		case <-ctx.Done():
			log.Debug().
				Err(ctx.Err()).
				Bool("phase1", !phase3).
				Bool("phase3", phase3).
				Msg("context canceled during gossip recv")
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
				log.Error().Err(err).Msg("gossip aborted early with error")
			}

			// The stream has been closed by the client, time to stop phase 3 (or phase 1 if there is an error)
			log.Trace().
				Bool("phase1", !phase3).
				Bool("phase3", phase3).
				Msg("remote phase 1/3 complete")
			return
		}

		// Handle the messages coming from the initiating replica
		switch sync.Status {
		case replica.Sync_CHECK:
			// If we're in phase 3 - meaning that the client sent a complete
			// message, then we ignore check messages to prevent checks bouncing
			// back and forth from poorly behaved clients.
			if phase3 {
				log.Warn().Msg("received check message in phase 3 from client")
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
					sender.Send(sync)
				} else {
					// This is an unhandled error, log and return error information.
					log.Warn().Err(err).
						Str("namespace", sync.Object.Namespace).
						Str("key", b64e(sync.Object.Key)).
						Msg("failed check sync on initiator: could not fetch object meta")

					// If there is an object error, send a response back to the
					// initiator's CHECK reqest so that it can be logged there as well.
					sync.Status = replica.Sync_ERROR
					sync.Error = err.Error()
					sender.Send(sync)
				}

				// CHECK is done here if there was an error.
				continue gossip
			}

			// Check which version is later, local or remote
			switch {
			case local.Version.IsLater(sync.Object.Version):
				// Send the local object back to the initiating replica as a repair
				sync.Status = replica.Sync_REPAIR
				sync.Object = local
				sync.Error = ""
				if ok := sender.Send(sync); ok {
					atomic.AddUint64(&updates, 1)
				}
			case sync.Object.Version.IsLater(local.Version):
				// Send a check request back to the initiating replica to fetch its version
				sender.Send(sync)
			default:
				// The versions are equal, do nothing
			}

		case replica.Sync_REPAIR:
			// Receiving a repaired object, check if it's still later than our local
			// replica's and if so, save it to disk with an update instead of a put.
			//
			// NOTE: honu.Update performs the version checking in a transaction.
			// TODO: record the update type in prometheus metrics.
			if _, err := r.db.Update(sync.Object, options.WithNamespace(sync.Object.Namespace)); err != nil {
				log.Error().Err(err).
					Str("namespace", sync.Object.Namespace).
					Str("key", b64e(sync.Object.Key)).
					Msg("could not update object from initiating replica")
				continue gossip
			}
			atomic.AddUint64(&repairs, 1)

		case replica.Sync_ERROR:
			// Something went wrong on the initiator, log and continue
			if sync.Object != nil {
				log.Warn().Str("error", sync.Error).
					Str("key", b64e(sync.Object.Key)).
					Str("namespace", sync.Object.Namespace).
					Msg("a replication error occurred")
			} else {
				log.Warn().Str("error", sync.Error).Msg("a replication error occurred")
			}

		case replica.Sync_COMPLETE:
			// Phase 1 is complete! Begin Phase 2 and Phase 3!
			// The phase3 bool causes this routine to stop saving seen messages and
			// the once sync ensures that multiple COMPLETE messages from the client
			// don't kick off more than one go routine.
			phase3 = true
			once.Do(func() {
				wg.Add(1)
				go r.remotePhase2(ctx, wg, log, &seen, sender, &updates)
			})

		default:
			log.Error().Str("status", sync.Status.String()).Msg("unhandled sync status")
		}
	}

	// Update Prometheus metrics
	metrics.PmAEPushes.WithLabelValues(r.conf.Name, r.conf.Region).Observe(float64(updates))
	metrics.PmAEPulls.WithLabelValues(r.conf.Name, r.conf.Region).Observe(float64(repairs))
	metrics.PmAEPushVSPull.Add(float64(updates))
	metrics.PmAEPushVSPull.Sub(float64(repairs))
}

// remotePhase2 is the counterpart to initiatorPhase1; the remote loops through all
// objects in its local database to check if there are any objects on the remote whose
// version vectors weren't seen during initiatorPhase1 - if so it means there is an
// object on the remote that the initiator hasn't seen before, so the remote sends a
// REPAIR message, pushing the object back. At the end of this go routine the remote
// sends a COMPLETE message, notifying the initiator that all phases of anti-entropy
// gossip are complete which allows the initiator to close the stream when ready.
// This go routine closes the sender channel when the phase is over because no more
// messages should be sent from the remote.
func (r *Service) remotePhase2(ctx context.Context, wg *sync.WaitGroup, log zerolog.Logger, seen *nsmap, sender *streamSender, updates *uint64) {
	// Ensure that this routine signals when it exists
	defer wg.Done()
	log.Trace().Msg("starting remote phase 2")

	// Ensure that we close the sending channel when this routine exits to prevent
	// deadlocks if this phase ends prematurely (e.g. the timeout expires).
	defer sender.Close()

	// Loop over all objects in all namespaces and determine what to push back
namespaces:
	for _, namespace := range r.replicatedNamespaces {
		iter, err := r.db.Iter(nil, options.WithNamespace(namespace))
		if err != nil {
			log.Error().Err(err).Str("namespace", namespace).Msg("could not iterate over namespace")
			continue namespaces
		}
		log.Trace().Str("namespace", namespace).Msg("sending namespace")

	objects:
		for iter.Next() {
			// Check if the context is done, and if so, break
			select {
			case <-ctx.Done():
				break namespaces
			default:
			}

			// Check if we've already seen this object in Phase 1 to deduplicate version
			// vector messages being sent back and forth. This check ensures that only
			// objects on the remote that are not on the initiator are pushed back.
			if seen.In(namespace, iter.Key()) {
				continue objects
			}

			// At this point, the remote has an object the initiator hasn't seen, send
			// a repair message to push the object back to the initiator (includes data).
			obj, err := iter.Object()
			if err != nil {
				log.Error().Err(err).
					Str("namespace", namespace).
					Str("key", b64e(iter.Key())).
					Msg("could not unmarshal honu metadata")
				continue objects
			}

			// Ensure that we can send, otherwise stop.
			if ok := sender.Send(&replica.Sync{Status: replica.Sync_REPAIR, Object: obj}); !ok {
				// NOTE: breaking objects loop not the namespaces loop or returning to
				// ensure that all iters get released. This may cause some additional
				// Send calls to happen, but the sender will simply ignore them.
				break objects
			}

			// Update the updates count so that we can log anti-entropy correctly.
			atomic.AddUint64(updates, 1)
		}

		if err = iter.Error(); err != nil {
			log.Error().Err(err).Str("namespace", namespace).Msg("could not iterate over namespace")
		}

		// Ensure the iterator is released, note even if break objects occurs, the
		// iterator should be released at this line since there is no return.
		iter.Release()
	}

	// Send a sync complete message to let the initiator know that the push phase is
	// complete and as soon as it is done processing repairs it can close the stream.
	// NOTE: this message MUST be sent, this function should not exit before this line.
	sender.Send(&replica.Sync{Status: replica.Sync_COMPLETE})
	log.Trace().Msg("phase 2 complete")
}

//===========================================================================
// Helper Methods
//===========================================================================

// Helper function to get the timestamp of last synchronization in a thread-safe manner
func (r *Service) LastSynchronization() string {
	r.RLock()
	defer r.RUnlock()
	if r.synchronized.IsZero() {
		return "never"
	}
	return r.synchronized.Format(time.RFC3339)
}

// Helper function to set the synchronized timestamp to now in a thread-safe manner
func (r *Service) synchronizedNow() {
	r.Lock()
	r.synchronized = time.Now()
	r.Unlock()
}
