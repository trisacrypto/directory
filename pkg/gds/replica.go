package gds

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	"github.com/trisacrypto/directory/pkg/gds/jitter"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/store/wire"
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

	// Check if the database Store is an ObjectStore, if not then the Replica cannot Gossip.
	if _, ok := svc.db.(global.ObjectStore); !ok {
		return nil, fmt.Errorf("replica %T does not implement global.ObjectStore", r.db)
	}

	// Initialize the gRPC server
	r.db = svc.db
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
	// Create the anti-entropy ticker
	ticker := jitter.New(r.conf.GossipInterval, r.conf.GossipSigma)

	// Run anti-entropy forever
bayou:
	for {
		<-ticker.C // Block until the next anti-entropy synchronization
		log.Debug().Msg("starting anti-entropy")

		// Randomly select a remote peer to synchronize with, continuing if we cannot
		// select a peer or no remote peers exist yet.
		var peer *peers.Peer
		if peer = r.selectPeer(); peer == nil {
			log.Debug().Msg("no remote peer available, skipping synchronization")
			continue bayou
		}

		// Ensure we can dial the client before we prepare the version vector
		// TODO: better initialization of gossip client and connection management
		cc, err := grpc.Dial(peer.Addr, grpc.WithInsecure())
		if err != nil {
			log.Error().Err(err).Str("addr", peer.Addr).Msg("could not dial remote peer")
		}
		client := global.NewReplicationClient(cc)
		log.Debug().Str("addr", peer.Addr).Str("peer", peer.String()).Msg("dialed remote peer")

		// Perepare a version vector to send to the remote peer
		// Because this is the initiation of anti-entropy this is not a partial request.
		versions := &global.VersionVectors{
			Objects:    make([]*global.Object, 0),
			Partial:    false,
			Namespaces: global.Namespaces[:],
		}

		// Access the objects in the object-store by namespace
		db := r.db.(global.ObjectStore)

		for _, ns := range versions.Namespaces {
			iter := db.Iter(ns)
		namespace:
			for iter.Next() {
				// Load the object metadata without the data itself, otherwise anti-entropy
				// would exchange way more data than required, putting pressure on memory.
				obj, err := iter.Object(false)
				if err != nil {
					log.Error().Err(err).Str("namespace", ns).Msg("could not unmarshal object")
					continue namespace
				}
				versions.Objects = append(versions.Objects, obj)
			}

			if err := iter.Error(); err != nil {
				log.Error().Err(err).Str("namespace", ns).Msg("could not iterate over object namespace")
			}
			iter.Release()
		}

		// Ensure that we send the request even if we have no local versions, to
		// retrieve any versions that might be on the remote peer.
		log.Debug().
			Int("versions", len(versions.Objects)).
			Int("namespaces", len(versions.Namespaces)).
			Msg("sending version vector to remote peer")

		// Perform the Gossip request
		var updates *global.Updates
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		if updates, err = client.Gossip(ctx, versions); err != nil {
			cancel()
			log.Error().Err(err).Str("peer", peer.String()).Msg("could not gossip with remote peer")
		}
		cancel()

		// Repair local database as last step
		for _, obj := range updates.Objects {
			if err = db.Put(obj); err != nil {
				log.Error().Err(err).Str("namespace", obj.Namespace).Msg("could not update local store")
			}
		}

		// Log success if any objects where synchronized
		if len(updates.Objects) > 0 {
			log.Info().
				Str("peer", peer.String()).
				Int("versions", len(versions.Objects)).
				Int("namespaces", len(versions.Namespaces)).
				Int("updates", len(updates.Objects)).
				Msg("anti-entropy synchronization complete")
		} else {
			log.Debug().Msg("anti-entropy complete with no synchronization")
		}
	}
}

// Randomly select a replica that is not self to perform anti-entropy with. If a peer
// cannot be selected, then nil is returned.
func (r *Replica) selectPeer() (peer *peers.Peer) {
	// Select a random peer that is not self to perform anti entropy with.
	peers, err := r.db.ListPeers()
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
func (r *Replica) Gossip(ctx context.Context, in *global.VersionVectors) (out *global.Updates, err error) {
	// TODO: identify remote peer via context
	log.Debug().
		Bool("partial", in.Partial).
		Int("nobjects", len(in.Objects)).
		Strs("namespaces", in.Namespaces).
		Msg("incoming anti-entropy")

	// Get the object store
	db := r.db.(global.ObjectStore)

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
		var localObj *global.Object
		if localObj, err = db.Get(remoteObj.Namespace, remoteObj.Key, true); err != nil {
			if errors.Is(err, wire.ErrObjectNotFound) {
				// This exists on the remote, but not locally; so add to fetch.
				if !in.Partial {
					fetch[remoteObj.Key] = struct{}{}
				}
				continue incomingLoop
			} else if errors.Is(err, wire.ErrCannotReplicate) {
				log.Warn().
					Str("namespace", remoteObj.Namespace).
					Str("key", remoteObj.Key).
					Msg("received known object namespace that should not be replicated")
				continue incomingLoop
			} else {
				// This is an unhandled error; log it and return replica requires repair
				log.Error().Err(err).Msg("could not get key from object store")
				return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
			}
		}

		// Step 1c: Check which version is later, local or remote.
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
		iter := db.Iter("")
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
			var localObj *global.Object
			if localObj, err = iter.Object(true); err != nil {
				if errors.Is(err, wire.ErrCannotReplicate) {
					// Ignore objects that cannot be replicated without warning and
					// don't count as part of local objects
					nLocalObjs--
					continue outgoingLoop
				} else {
					log.Error().
						Err(err).
						Str("namespace", strings.Split(iter.Key(), "::")[0]).
						Msg("could not unmarshal object from database")
					return nil, status.Error(codes.FailedPrecondition, "replica requires repair")
				}
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
