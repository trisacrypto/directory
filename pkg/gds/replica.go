package gds

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
	"github.com/trisacrypto/directory/pkg/gds/store/leveldb"
	"google.golang.org/grpc"
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

	// TODO: right now replica only works with LevelDB need to adapt the Store to work
	// for other store types such as sqlite3.
	var ok bool
	if r.db, ok = svc.db.(*leveldb.Store); !ok {
		return nil, fmt.Errorf("replica currently only works with leveldb, not %T", r.db)
	}

	// Initialize the gRPC server
	r.srv = grpc.NewServer(grpc.UnaryInterceptor(svc.serverInterceptor))
	global.RegisterReplicationServer(r.srv, r)
	return r, nil
}

// Replica implements the ReplicationServer as defined by the v1 or later GDS protocol
// buffers. This service is a machine-to-machine implementation that allows GDS to be
// globally distributed by implementing auto-adapting anti-entropy.
type Replica struct {
	global.UnimplementedReplicationServer
	peers.UnimplementedPeerManagementServer
	svc  *Service              // The parent Service the replica uses to interact with other components
	srv  *grpc.Server          // The gRPC server that listens on its own independent port
	conf *config.ReplicaConfig // The replica specific configuration (alias to r.svc.conf.Replica)
	db   *leveldb.Store        // Database connection for managing objects (alias to s.svc.db)
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

// During gossip, the initiating replica sends a randomly selected remote peer the
// version vectors of all objects it currently stores. The remote peer should
// respond with updates that correspond to more recent versions of the objects. The
// remote peer can than also make a reciprocal request for updates by sending the
// set of versions requested that were more recent on the initiating replica, and
// use a partial flag to indicate that it is requesting specific versions. This
// mechanism implements bilateral anti-entropy: a push and pull gossip.
func (r *Replica) Gossip(ctx context.Context, in *global.VersionVectors) (out *global.Updates, err error) {
	log.Debug().
		Bool("partial", in.Partial).
		Int("nobjects", len(in.Objects)).
		Strs("namespaces", in.Namespaces).
		Msg("incoming anti-entropy")

	// // TODO: don't use leveldb forever
	// ldb := r.db.DB()

	// out = &global.Updates{
	// 	Objects: make([]*global.Object, 0),
	// }

	// // Step 1: determine if any of the incoming objects have a later version locally
	// // or if the remote version is later than our local version.
	// seen := make(map[string]struct{})
	// fetch := make(map[string]struct{})

	// // Step 1a: determine the namespaces to iterate over
	// namespaces := in.Namespaces
	// if len(namespaces) == 0 {
	// 	namespaces = allNamespaces
	// }

	// // Step

	// // Step 2: if not partial, determine if any new objects exist locally to send back
	// // to the remote.

	// // Step 3: request updates if any via a partial request back to the remote.
	// // NOTE: do not send a request to a partial update (assuming we're at the end of
	// // bilateral anti-entropy) to prevent possible infinite recursion.

	// if in.Partial {
	// 	// Only consider the objects sent in the version vectors
	// } else {
	// 	// Consider all objects
	// }

	return &global.Updates{}, nil
}

// TODO: what are we doing with Status?
// GetPeers queries the data store to determine which peers it contains, and returns them
func (r *Replica) GetPeers(ctx context.Context, in *peers.PeersFilter) (out *peers.PeersList, err error) {

	// TODO: Not sure what StatusOnly is for
	if in.StatusOnly {
		return nil, errors.New("StatusOnly not supported yet")
	}

	// Initialize var for candidate peers
	var peers []*peers.Peer

	// Get all the peers
	if peers, err = r.db.ListPeers(); err != nil {
		// TODO: Not sure what error we want here
		return nil, errors.New("No peers retrieved from the database")
	}

	// If there is no region filter on the request, return all peers
	if in.Region == nil {
		out.Peers = peers
	}

	// Otherwise, use the regions to determin which peers to keep
	for _, peer := range peers {
		for _, region := range in.Region {
			if peer.Region == region {
				out.Peers = append(out.Peers, peer)
			}
		}
	}

	return out, nil
}

func (r *Replica) AddPeers(ctx context.Context, in *peers.Peer) (out *peers.PeersStatus, err error) {
	return nil, nil
}

func (r *Replica) RmPeers(ctx context.Context, in *peers.Peer) (out *peers.PeersStatus, err error) {
	return nil, nil
}
