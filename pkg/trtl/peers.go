package trtl

import (
	context "context"
	"errors"
	"time"

	"github.com/rotationalio/honu"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/object"
	"github.com/rotationalio/honu/options"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// A PeerService implements the RPCs for managing remote peers.
type PeerService struct {
	peers.UnimplementedPeerManagementServer
	parent *Server
	db     *honu.DB
}

func NewPeerService(s *Server) (*PeerService, error) {
	return &PeerService{
		parent: s,
		db:     s.db,
	}, nil
}

// GetPeers queries the data store to determine which peers it contains, and returns them
func (p *PeerService) GetPeers(ctx context.Context, in *peers.PeersFilter) (out *peers.PeersList, err error) {
	if out, err = p.peerStatus(ctx, in); err != nil {
		// peerStatus returns status error and does logging
		return nil, err
	}

	return out, nil
}

// AddPeers adds a peer and returns a report of the status of all peers in the network.
func (p *PeerService) AddPeers(ctx context.Context, in *peers.Peer) (out *peers.PeersStatus, err error) {
	// Compute the key from the Peer method to ensure consistent key access.
	// If the key is empty, that means the peer ID == 0; return an error.
	key := in.Key()
	if key == "" {
		return nil, status.Error(codes.InvalidArgument, "a valid peer ID is required")
	}

	// Check if the peer is in the database; if it is unmarshal it and update it
	if data, err := p.db.Get([]byte(key), options.WithNamespace(NamespacePeers)); err != nil {
		if !errors.Is(err, engine.ErrNotFound) {
			log.Error().Err(err).Msg("could not access database for peers lookup")
			return nil, status.Error(codes.FailedPrecondition, "could not get peer from database")
		}

		current := new(peers.Peer)
		if err = proto.Unmarshal(data, current); err != nil {
			log.Error().Err(err).Str("key", key).Msg("could not unmarshal peer from database")
			return nil, status.Error(codes.FailedPrecondition, "could not get peer from database")
		}

		// Update the incoming peer with data from the previous peer
		// TODO: merge empty fields from current peer into incoming peer then validate.
		in.Created = current.Created

	}

	// TODO: validate other Peer fields
	in.Modified = time.Now().Format(time.RFC3339)

	// Insert the peer into the database
	var value []byte
	if value, err = proto.Marshal(in); err != nil {
		log.Error().Err(err).Msg("could not marshal peer protocol buffers")
		return nil, status.Error(codes.FailedPrecondition, "could not marshal peer protocol buffers")
	}

	if _, err = p.db.Put([]byte(key), value, options.WithNamespace(NamespacePeers)); err != nil {
		log.Error().Err(err).Msg("could not put peer to database")
		return nil, status.Error(codes.FailedPrecondition, "could not insert peer into database")
	}

	// Assuming we don't need all the Peer details in this case
	ftr := &peers.PeersFilter{
		StatusOnly: true,
	}
	if pl, err := p.peerStatus(ctx, ftr); err != nil {
		return nil, err
	} else {
		out = pl.Status
	}
	return out, nil
}

func (p *PeerService) RmPeers(ctx context.Context, in *peers.Peer) (out *peers.PeersStatus, err error) {
	if _, err = p.db.Delete([]byte(in.Key()), options.WithNamespace(NamespacePeers)); err != nil {
		log.Error().Err(err).Msg("unable to remove peer")
		return nil, status.Error(codes.InvalidArgument, "invalid peer; could not be removed")
	}

	// Assuming we don't need all the Peer details in this case
	ftr := &peers.PeersFilter{
		StatusOnly: true,
	}
	if pl, err := p.peerStatus(ctx, ftr); err != nil {
		return nil, err
	} else {
		out = pl.Status
	}
	return out, nil
}

// Helper to get the peer network status
func (p *PeerService) peerStatus(ctx context.Context, in *peers.PeersFilter) (out *peers.PeersList, err error) {
	// Create the response
	out = &peers.PeersList{
		Peers: make([]*peers.Peer, 0),
		Status: &peers.PeersStatus{
			Regions:             make(map[string]int64),
			LastSynchronization: p.parent.replica.lastSynchronization(),
		},
	}

	// Iterate over all the peers (necessary for both list and status-only)
	// TODO: filter self from the list?
	ps, err := p.db.Iter(nil, options.WithNamespace(NamespacePeers))
	if err != nil {
		return nil, err
	}
	defer ps.Release()

	for ps.Next() {
		// Skip tombstones
		// TODO: why is Honu returning tombstones in iter?
		var obj *object.Object
		if obj, err = ps.Object(); err != nil {
			log.Error().Err(err).Str("key", string(ps.Key())).Msg("could not retrieve object from db")
			continue
		}

		if obj.Tombstone() {
			continue
		}

		peer := new(peers.Peer)
		if err = proto.Unmarshal(obj.Data, peer); err != nil {
			log.Warn().Err(err).Str("key", string(ps.Key())).Msg("could not unmarshal peer")
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
						// Append peer to peers then break out of the regions loop in
						// case the user has specified multiple duplicate regions.
						out.Peers = append(out.Peers, peer)
						break
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
