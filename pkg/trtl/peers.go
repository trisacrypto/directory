package trtl

import (
	context "context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/store/iterator"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type TemporaryPeerStore interface {
	CreatePeer(peer *peers.Peer) (string, error)
	DeletePeer(id string) error
	ListPeers() iterator.Iterator
	AllPeers() ([]*peers.Peer, error)
}

type notImplementedStore struct{}

func (s *notImplementedStore) CreatePeer(peer *peers.Peer) (string, error) {
	return "", errors.New("not implemented")
}
func (s *notImplementedStore) DeletePeer(id string) error   { return errors.New("not implemented") }
func (s *notImplementedStore) ListPeers() iterator.Iterator { return nil }
func (s *notImplementedStore) AllPeers() ([]*peers.Peer, error) {
	return nil, errors.New("not implemented")
}

// A PeerService implements the RPCs for managing remote peers.
type PeerService struct {
	peers.UnimplementedPeerManagementServer
	parent *Server
	store  TemporaryPeerStore
}

func NewPeerService(s *Server) (*PeerService, error) {
	return &PeerService{
		parent: s,
		store:  &notImplementedStore{},
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

// AddPeers adds a peer and returns a report of the status of all peers in the network
func (p *PeerService) AddPeers(ctx context.Context, in *peers.Peer) (out *peers.PeersStatus, err error) {
	// CreatePeer handles possibility of an already-existing or previously deleted peer
	if _, err := p.store.CreatePeer(in); err != nil {
		log.Error().Err(err).Msg("unable to add peer")
		return nil, status.Error(codes.InvalidArgument, "invalid peer; could not be added")
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
	// TODO: check what kind of errors delete peer returns.
	if err := p.store.DeletePeer(in.Key()); err != nil {
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
		Peers:  make([]*peers.Peer, 0),
		Status: &peers.PeersStatus{},
	}

	// Iterate over all the peers (necessary for both list and status-only)
	// TODO: filter self from the list?
	ps := p.store.ListPeers()
	defer ps.Release()

	// TODO: Implement peer iteration
	/*
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
		}*/

	if err = ps.Error(); err != nil {
		log.Error().Err(err).Msg("unable to retrieve peers from the database")
		return nil, status.Error(codes.FailedPrecondition, "error reading from database")
	}
	return out, nil
}
