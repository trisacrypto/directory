package gds

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"

	"github.com/rotationalio/honu/replica"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"google.golang.org/grpc"
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

// NewReplica creates a new GDS replica server derived from a parent Service.
func NewReplica(svc *Service) (r *Replica, err error) {
	r = &Replica{
		svc:  svc,
		conf: &svc.conf.Replica,
	}

	// TODO: Check if the database Store is an Honu DB, if not then the Replica cannot Gossip.

	// Initialize the gRPC server
	r.db = svc.db
	r.srv = grpc.NewServer(grpc.UnaryInterceptor(svc.serverInterceptor))
	replica.RegisterReplicationServer(r.srv, r)
	return r, nil
}

// Replica implements the ReplicationServer as defined by the v1 or later GDS protocol
// buffers. This service is a machine-to-machine implementation that allows GDS to be
// globally distributed by implementing auto-adapting anti-entropy.
type Replica struct {
	replica.UnimplementedReplicationServer
	peers.UnimplementedPeerManagementServer
	svc  *Service              // The parent Service the replica uses to interact with other components
	srv  *grpc.Server          // The gRPC server that listens on its own independent port
	conf *config.ReplicaConfig // The replica specific configuration (alias to r.svc.conf.Replica)
	db   store.Store           // Database connection for managing objects (alias to s.svc.db)
}

// Serve gRPC requests on the specified bind address.
func (r *Replica) Serve() (err error) {
	// This service should not be started in maintenance mode.
	if r.svc.conf.Maintenance {
		return errors.New("could not start replication service in maintenance mode")
	}

	if !r.conf.Enabled {
		log.Warn().Msg("replication service is not enabled")
		return nil
	}

	// Listen for TCP requests on the specified address and port
	var sock net.Listener
	if sock, err = net.Listen("tcp", r.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on %q", r.conf.BindAddr)
	}

	// Run the server
	go func() {
		defer sock.Close()
		log.Info().
			Str("listen", r.conf.BindAddr).
			Str("version", pkg.Version()).
			Msg("replication service started")

		if err := r.srv.Serve(sock); err != nil {
			r.svc.echan <- err
		}
	}()

	// Run the Gossip background routine
	go r.AntiEntropy()

	// The server go routine is started so return nil error (any server errors will be
	// sent on the error channel).
	return nil
}

// Shutdown the Replication Service gracefully
func (r *Replica) Shutdown() error {
	log.Debug().Msg("gracefully shutting down replication server")
	r.srv.GracefulStop()
	log.Debug().Msg("successful shutdown of replica server")
	return nil
}

// AntiEntropy is a service that periodically selects a remote peer to synchronize with
// via bilateral anti-entropy using the Gossip service. Jitter is applied to the
// interval between anti-entropy synchronizations to ensure that message traffic isn't
// bursty to disrupt normal messages to the GDS service.
// TODO: this background routine is currently untested.
func (r *Replica) AntiEntropy() {
	log.Warn().Msg("anti-entropy is not implemented; no anti-entropy is running")
}

// SelectPeer randomly that is not self to perform anti-entropy with. If a peer
// cannot be selected, then nil is returned.
func (r *Replica) SelectPeer() (peer *peers.Peer) {
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
func (r *Replica) Gossip(ctx context.Context, in *replica.VersionVectors) (out *replica.Updates, err error) {
	return nil, status.Error(codes.Unimplemented, "this replica does not yet implement gossip")
}

// GetPeers queries the data store to determine which peers it contains, and returns them
func (r *Replica) GetPeers(ctx context.Context, in *peers.PeersFilter) (out *peers.PeersList, err error) {

	if out, err = r.peerStatus(ctx, in); err != nil {
		// peerStatus returns status error and does logging
		return nil, err
	}

	return out, nil
}

// AddPeers adds a peer and returns a report of the status of all peers in the network
func (r *Replica) AddPeers(ctx context.Context, in *peers.Peer) (out *peers.PeersStatus, err error) {
	// CreatePeer handles possibility of an already-existing or previously deleted peer
	if _, err := r.db.CreatePeer(in); err != nil {
		log.Error().Err(err).Msg("unable to add peer")
		return nil, status.Error(codes.InvalidArgument, "invalid peer; could not be added")
	}

	// Assuming we don't need all the Peer details in this case
	ftr := &peers.PeersFilter{
		StatusOnly: true,
	}
	if pl, err := r.peerStatus(ctx, ftr); err != nil {
		return nil, err
	} else {
		out = pl.Status
	}
	return out, nil
}

func (r *Replica) RmPeers(ctx context.Context, in *peers.Peer) (out *peers.PeersStatus, err error) {
	if err := r.db.DeletePeer(in.Key()); err != nil {
		log.Error().Err(err).Msg("unable to remove peer")
		return nil, status.Error(codes.InvalidArgument, "invalid peer; could not be removed")
	}

	// Assuming we don't need all the Peer details in this case
	ftr := &peers.PeersFilter{
		StatusOnly: true,
	}
	if pl, err := r.peerStatus(ctx, ftr); err != nil {
		return nil, err
	} else {
		out = pl.Status
	}
	return out, nil
}

// Helper to get the peer network status
func (r *Replica) peerStatus(ctx context.Context, in *peers.PeersFilter) (out *peers.PeersList, err error) {
	// Create the response
	out = &peers.PeersList{
		Peers:  make([]*peers.Peer, 0),
		Status: &peers.PeersStatus{},
	}

	// Iterate over all the peers (necessary for both list and status-only)
	// TODO: filter self from the list?
	ps := r.db.ListPeers()
	defer ps.Release()

	for ps.Next() {
		peer := ps.Peer()
		if peer == nil {
			continue
		}

		out.Status.NetworkSize++
		out.Status.Regions[peer.Region]++

		// If it's not a status only, get the details for each Peer
		if !in.StatusOnly {
			// If we've been asked to filter by region
			if len(in.Region) > 0 {
				for _, region := range in.Region {
					if peer.Region == region {
						out.Peers = append(out.Peers, peer)
					}
				}
			} else {
				// Otherwise don't filter and keep all the Peers
				out.Peers = append(out.Peers, peer)
			}
		}
	}

	if err = ps.Error(); err != nil {
		log.Error().Err(err).Msg("unable to retrieve peers from the database")
		return nil, status.Error(codes.FailedPrecondition, "error reading from database")
	}
	return out, nil
}
