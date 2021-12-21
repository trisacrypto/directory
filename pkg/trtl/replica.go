package trtl

import (
	"context"
	"math/rand"

	"github.com/rotationalio/honu"
	"github.com/rotationalio/honu/options"
	"github.com/rotationalio/honu/replica"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// A ReplicaService manages anti-entropy replication between peers.
type ReplicaService struct {
	conf   config.Config
	parent *Server
	db     *honu.DB
}

func NewReplicaService(s *Server) (*ReplicaService, error) {
	return &ReplicaService{
		parent: s, db: s.db,
	}, nil
}

// AntiEntropy is a service that periodically selects a remote peer to synchronize with
// via bilateral anti-entropy using the Gossip service. Jitter is applied to the
// interval between anti-entropy synchronizations to ensure that message traffic isn't
// bursty to disrupt normal messages to the GDS service.
// TODO: this background routine is currently untested.
func (*ReplicaService) AntiEntropy() {
	log.Warn().Msg("anti-entropy is not implemented; no anti-entropy is running")
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

			if peer.Id != r.conf.Replica.PID {
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
