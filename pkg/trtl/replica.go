package trtl

import (
	"context"
	"encoding/base64"
	"math/rand"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/rotationalio/honu/replica"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/jitter"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// A ReplicaService manages anti-entropy replication between peers.
type ReplicaService struct {
	replica.UnimplementedReplicationServer
	parent *Server
	conf   config.ReplicaConfig
	db     *honu.DB
	aestop chan struct{}
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
	// Create the anti-entropy ticker and store the channel for shutdown
	ticker := jitter.New(r.conf.GossipInterval, r.conf.GossipSigma)
	r.aestop = stop

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
		logctx := log.With().Uint64("peer_id", peer.Id).Str("peer_addr", peer.Addr).Str("peer", peer.Name).Str("service", "anti-entropy").Logger()

		// Ensure we can dial the client before we prepare the version vector.
		// TODO: add mTLS to peer-to-peer connection
		// TODO: better initialization of gossip client and connection management.
		var (
			cc  *grpc.ClientConn
			err error
		)

		if cc, err = grpc.Dial(peer.Addr, grpc.WithInsecure(), grpc.WithBlock()); err != nil {
			logctx.Warn().Err(err).Msg("could not dial remote peer")
			continue bayou
		}

		client := replica.NewReplicationClient(cc)
		logctx.Debug().Msg("dialed remote peer")

		// Prepare version vector to send to the remote peer.
		// Note that this is a full request because the bilateral exchange is being initiated
		versions := &replica.VersionVectors{
			Objects:    make([]*object.Object, 0),
			Partial:    false,
			Namespaces: replicatedNamespaces,
		}

		// Access the objects in the object-store by namespace
	namespaces:
		for _, namespace := range versions.Namespaces {
			iter, err := r.db.Iter(nil, options.WithNamespace(namespace))
			if err != nil {
				log.Error().Err(err).Str("namespace", namespace).Msg("could not iterate over namespace")
				continue namespaces
			}

		namespace:
			for iter.Next() {
				// Load the object metadata without the data itself, otherwise anti-
				// entropy would exchange way more data than required, putting pressure
				// on pod memory and increasing our cloud bill.
				obj, err := iter.Object()
				if err != nil {
					log.Error().Err(err).
						Str("namespace", namespace).
						Str("key", b64e(iter.Key())).
						Msg("could not unmarshal honu metadata")
					continue namespace
				}

				// Remove the data from the object
				obj.Data = nil
				versions.Objects = append(versions.Objects, obj)
			}

			if err = iter.Error(); err != nil {
				log.Error().Err(err).Str("namespace", namespace).Msg("could not iterate over namespace")
			}
			iter.Release()
		}

		// Ensure we send the request even if we have no local versions to retrieve
		// any versions that might be on the remote peer replica.
		logctx.Debug().
			Int("versions", len(versions.Objects)).
			Int("namespaces", len(versions.Namespaces)).
			Msg("sending version vector to remote peer")

		// Perform the gossip request
		var updates *replica.Updates
		ctx, cancel := context.WithTimeout(context.Background(), r.conf.GossipInterval)
		if updates, err = client.Gossip(ctx, versions); err != nil {
			cancel()
			logctx.Error().Err(err).Msg("could not gossip with remote peer")
			continue bayou
		}
		cancel()

		// Repair local database as last step
		var nUpdates uint64
		for _, obj := range updates.Objects {
			if err = r.db.Update(obj, options.WithNamespace(obj.Namespace)); err != nil {
				logctx.Error().Err(err).Str("namespace", obj.Namespace).Str("key", b64e(obj.Key)).Msg("could not update object from remote peer")
			} else {
				nUpdates++
			}
		}

		// Log success if any objects were synchronized
		if nUpdates > 0 {
			logctx.Info().
				Int("versions", len(versions.Objects)).
				Int("namespaces", len(versions.Namespaces)).
				Uint64("updates", nUpdates).
				Msg("anti-entropy synchronization complete")
		} else {
			logctx.Trace().Msg("anti-entropy complete with no synchronization")
		}
	}
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

// During gossip, the initiating replica sends a randomly selected remote peer the
// version vectors of all objects it currently stores. The remote peer should
// respond with updates that correspond to more recent versions of the objects. The
// remote peer can than also make a reciprocal request for updates by sending the
// set of versions requested that were more recent on the initiating replica, and
// use a partial flag to indicate that it is requesting specific versions. This
// mechanism implements bilateral anti-entropy: a push and pull gossip.
func (r *ReplicaService) Gossip(ctx context.Context, in *replica.VersionVectors) (out *replica.Updates, err error) {
	return nil, status.Error(codes.Unimplemented, "this replica does not yet implement gossip")
}

func (r *ReplicaService) Shutdown() error {
	if r.aestop != nil {
		r.aestop <- struct{}{}
	}
	return nil
}

func b64e(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}
