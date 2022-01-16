package replica

import (
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rotationalio/honu"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/rotationalio/honu/replica"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/config"
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
	if r.aestop != nil {
		r.aestop <- struct{}{}
	}
	return nil
}

//===========================================================================
// Gossip (server-side) Methods
//===========================================================================

// Gossip implements bilateral anti-entropy: during a Gossip session the initiating
// replica pushes updates to the remote peer and pulls requested changes. Using
// bidirectional streaming, the initiating peer sends data-less sync messages with
// the versions of objects it stores locally. The remote replica then responds with
// data if its local version is later or sends a sync message back requesting the
// data from the initiating replica if its local version is earlier. No exchange
// occurs if both replicas have the same version. At the end of a gossip session,
// both replicas should have synchronized and have identical underlying data stores.
func (r *Service) Gossip(stream replica.Replication_GossipServer) (err error) {
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
			for _, namespace := range r.replicatedNamespaces {
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
	return nil
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
