package trtl

import (
	"context"
	"math/rand"

	"github.com/rotationalio/honu/replica"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	vaspPrefix    = "vasps::"
	certreqPrefix = "certreqs::"
	peerPrefix    = "peers::"
)

var (
	// All of the namespaces that are allowed for replication.
	AllNamespaces = []string{vaspPrefix, certreqPrefix, peerPrefix}
)

// A ReplicaService manages anti-entropy replication between peers.
type ReplicaService struct {
	conf config.Config
	db   store.Store
}

func NewReplicaService(db store.Store, conf config.Config) *ReplicaService {
	return &ReplicaService{
		db:   db,
		conf: conf,
	}
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
	peers, err := r.db.ListPeers().All()
	if err != nil {
		log.Error().Err(err).Msg("could not fetch peers from database")
		return nil
	}

	if len(peers) > 1 {
		// 10 attempts to select a random peer that is not self.
		for i := 0; i < 10; i++ {
			peer = peers[rand.Intn(len(peers))]
			if peer.Id != r.conf.PID {
				return peer
			}
		}
		log.Warn().Int("nPeers", len(peers)).Msg("could not select peer after 10 attempts")
	}

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
