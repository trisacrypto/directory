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
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/store/leveldb"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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

	// TODO: don't use leveldb for v1.1; this is just for prototype purposes.
	ldb := r.db.DB()

	out = &global.Updates{
		Objects: make([]*global.Object, 0),
	}

	// Step 1: determine if any of the incoming objects have a later version locally
	// or if the remote version is later than our local version.
	seen := make(map[string]struct{})
	fetch := make(map[string]struct{})

	// Loop over all incoming objects
incomingLoop:
	for _, remoteObj := range in.Objects {
		// Fetch data from database (determining if the data is not found)
		// TODO: move object fetching logic back to Store; ensure tombstones are included.
		var data []byte
		if data, err = ldb.Get([]byte(remoteObj.Key), nil); err != nil {
			if errors.Is(err, leveldb.ErrEntityNotFound) {
				// This exists on the remote, but not locally; so add to fetch.
				if !in.Partial {
					fetch[remoteObj.Key] = struct{}{}
				}
				continue incomingLoop
			} else {
				// This is an unhandled error; log it and return replica requires repair
				log.Error().Err(err).Msg("could not get key from leveldb")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}
		}

		// Get the localObj representation
		var localObj *global.Object
		switch remoteObj.Namespace {
		case store.NamespaceVASPs:
			vasp := &pb.VASP{}
			if err = proto.Unmarshal(data, vasp); err != nil {
				// This is an unhandled error; log it and return unavailable
				log.Error().Err(err).Msg("could not unmarshal VASP from leveldb")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}

			if localObj, _, err = models.GetMetadata(vasp); err != nil {
				// This is an unhandled error; log it and return unavailable
				log.Error().Err(err).Msg("could not get metadata from VASP")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}

			if localObj.Data, err = anypb.New(vasp); err != nil {
				log.Error().Err(err).Msg("could not marshal VASP Any back to remote replica")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}

		case store.NamespaceCertReqs:
			certreq := &models.CertificateRequest{}
			if err = proto.Unmarshal(data, certreq); err != nil {
				// This is an unhandled error; log it and return unavailable
				log.Error().Err(err).Msg("could not unmarshal CertificateRequest from leveldb")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}
			localObj = certreq.Metadata
			if localObj.Data, err = anypb.New(certreq); err != nil {
				log.Error().Err(err).Msg("could not marshal CertificateRequest Any back to remote replica")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}
		default:
			// Log error but continue processing, foreign namespace will be ignored
			log.Warn().Str("namespace", remoteObj.Namespace).Msg("unknown namespace")
			continue incomingLoop
		}

		// Check which version is later, local or remote.
		switch {
		case localObj.Version.IsLater(remoteObj.Version):
			// Send the local object back to the remote in the updates
			out.Objects = append(out.Objects, localObj)
		case remoteObj.Version.IsLater(localObj.Version):
			// Mark the remoteObj to be fetched
			fetch[remoteObj.Key] = struct{}{}
		default:
			// The versions are equal; do nothing
		}

		// Add the remoteObj to seen to make sure that we do not handle it in the next phase
		seen[remoteObj.Key] = struct{}{}
	}

	// Step 2: if not partial, determine if any new objects exist locally to send back
	// to the remote.

	// Step 3: request updates if any via a partial request back to the remote.
	// NOTE: do not send a request to a partial update (assuming we're at the end of
	// bilateral anti-entropy) to prevent possible infinite recursion.

	if in.Partial {
		// Only consider the objects sent in the version vectors
		log.Debug().Msg("sending updates")
	} else {
		// Consider all objects
		log.Debug().Msg("determining what is available locally and not on the remote")
	}

	return &global.Updates{}, nil
}
