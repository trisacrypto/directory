package gds

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
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
	// TODO: identify remote peer via context
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

	// Step 1a: Loop over all incoming objects
incomingLoop:
	for _, remoteObj := range in.Objects {
		// Step 1b: Fetch data from database (determining if the data is not found)
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

		// Step 1c: Get the localObj representation based on the object namespace
		// TODO: create an interface that returns the object representation
		var localObj *global.Object
		switch remoteObj.Namespace {
		case store.NamespaceVASPs:
			vasp := &pb.VASP{}
			if err = proto.Unmarshal(data, vasp); err != nil {
				// This is an unhandled error; log it and return unavailable
				log.Error().Err(err).Msg("could not unmarshal VASP from leveldb")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}

			// VASP specific retrieval of object metadata
			if localObj, _, err = models.GetMetadata(vasp); err != nil {
				// This is an unhandled error; log it and return unavailable
				log.Error().Err(err).Msg("could not get metadata from VASP")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}

			// Ensure vasp is stored on the object in case it is sent back to remote
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

			// Certreq specific method to retrieve metadata
			localObj = certreq.Metadata

			// Ensure certreq is stored on the object in case it is sent back to remote
			if localObj.Data, err = anypb.New(certreq); err != nil {
				log.Error().Err(err).Msg("could not marshal CertificateRequest Any back to remote replica")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}
		default:
			// Log error but continue processing, foreign namespace will be ignored
			log.Warn().Str("namespace", remoteObj.Namespace).Msg("unknown namespace")
			continue incomingLoop
		}

		// Step 1d: Check which version is later, local or remote.
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
	// to the remote (if partial, is likely the pull phase of bilateral anti-entropy)
	var nLocalObjs uint64
	if !in.Partial {
		// Step 2a: loop over all keys in the database, ignoring any that have already been seen
		iter := ldb.NewIterator(nil, nil)
	outgoingLoop:
		for iter.Next() {
			nLocalObjs++
			key := string(iter.Key())
			if _, ok := seen[key]; ok {
				// We've already handled this key in the incomingLoop, ignore
				continue outgoingLoop
			}

			// Step 2b: if this key hasn't been seen then it is a new local key that
			// needs to be pushed back to the remote replica. Load the object from
			// the database and add to outgoing objects.
			data := iter.Value()
			prefix := strings.Split(key, "::")[0]

			// TODO: create an interface that handles the object
			var localObj *global.Object
			switch prefix {
			case store.NamespaceVASPs:
				vasp := &pb.VASP{}
				if err = proto.Unmarshal(data, vasp); err != nil {
					// This is an unhandled error; log it and return unavailable
					log.Error().Err(err).Msg("could not unmarshal VASP from leveldb")
					return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
				}

				// VASP specific retrieval of object metadata
				if localObj, _, err = models.GetMetadata(vasp); err != nil {
					// This is an unhandled error; log it and return unavailable
					log.Error().Err(err).Msg("could not get metadata from VASP")
					return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
				}

				// Ensure vasp is stored on the object in case it is sent back to remote
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

				// Certreq specific method to retrieve metadata
				localObj = certreq.Metadata

				// Ensure certreq is stored on the object in case it is sent back to remote
				if localObj.Data, err = anypb.New(certreq); err != nil {
					log.Error().Err(err).Msg("could not marshal CertificateRequest Any back to remote replica")
					return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
				}
			case store.NamespaceIndices:
				// Ignore indices without warning and don't count as part of local objects
				nLocalObjs--
				continue outgoingLoop
			default:
				// Log error but continue processing, foreign namespace will be ignored
				log.Warn().Str("namespace", prefix).Msg("unknown namespace")
				continue outgoingLoop
			}

			// Step 2c: Add the local object to send back in the response
			out.Objects = append(out.Objects, localObj)
		}

		// Step 2d: Cleanup after database iteration
		if err = iter.Error(); err != nil {
			iter.Release()
			return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
		}
		iter.Release()
	}

	// Step 3: request updates if any via a partial request back to the remote.
	// NOTE: do not send a request to a partial update (assuming we're at the end of
	// bilateral anti-entropy) to prevent possible infinite recursion.
	// TODO: for this step to work, the remote peer must be identified.
	// TODO: implement step 3 after remote peer can be identified
	if len(fetch) > 0 {
		log.Warn().Int("fetch", len(fetch)).Msg("remote peer not identified, cannot pull objects")
	}

	log.Info().
		Bool("partial", in.Partial).
		Int("nLocal", int(nLocalObjs)).
		Int("nRemote", len(in.Objects)).
		Int("nRepairRemote", len(out.Objects)).
		Int("nRepairLocal", len(fetch)).
		Msg("anti-entropy gossip request handled")

	return out, nil
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

	// Initialize var for candidate peers
	var ps []*peers.Peer

	// Get all the peers (necessary for both list and status-only)
	if ps, err = r.db.ListPeers(); err != nil {
		log.Error().Err(err).Msg("unable to retrieve peers from the database")
		return nil, status.Error(codes.FailedPrecondition, "error reading from database")
	}

	// Get an overall replica count - we need this regardless
	// TODO: delete self from the list?
	out = &peers.PeersList{
		Peers:  make([]*peers.Peer, 0, len(ps)),
		Status: &peers.PeersStatus{},
	}

	out.Status.NetworkSize = int64(len(ps))

	for _, peer := range ps {
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
	return out, nil
}
