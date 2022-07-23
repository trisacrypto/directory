package trtl_test

import (
	"context"

	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
)

// TestPeers is a unified function that performs multiple interactions with the Peers
// server, getting, adding, and removing peers to exercise all Peers server methods
// without requiring a dedicated test fixture.
func (s *trtlTestSuite) TestPeers() {
	// Setup the test and create a new Peers client
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := peers.NewPeerManagementClient(s.grpc.Conn)

	// GetPeers should return nothing
	rep, err := client.GetPeers(ctx, &peers.PeersFilter{})
	require.NoError(err, "could not get empty peers")
	require.Len(rep.Peers, 0, "unexpected peers returned, have the fixtures changed?")
	require.Equal(int64(0), rep.Status.NetworkSize, "unexpected network size returned")
	require.Len(rep.Status.Regions, 0, "unexpected regions mapping returned")

	// Add several peers in different regions
	network := []*peers.Peer{
		{
			Id:     1,
			Addr:   "tc1.trtl.dev:443",
			Name:   "alpha",
			Region: "tauceti",
		},
		{
			Id:     2,
			Addr:   "tc2.trtl.dev:443",
			Name:   "bravo",
			Region: "tauceti",
		},
		{
			Id:     3,
			Addr:   "ei.trtl.dev:443",
			Name:   "charlie",
			Region: "epsiloniridini",
		},
		{
			Id:     4,
			Addr:   "wolf.trtl.dev:443",
			Name:   "delta",
			Region: "wolf359",
		},
	}

	for _, peer := range network {
		rep, err := client.AddPeers(ctx, peer)
		require.NoError(err, "could not add peer to server")
		require.NotEmpty(rep.NetworkSize)
		require.NotEmpty(rep.Regions)
	}

	// GetPeers should return the complete network
	rep, err = client.GetPeers(ctx, &peers.PeersFilter{})
	require.NoError(err, "could not get full network of peers")
	require.Len(rep.Peers, 4, "peers not returned after they were added")
	require.Equal(int64(4), rep.Status.NetworkSize, "unexpected network size returned")
	require.Len(rep.Status.Regions, 3, "unexpected regions mapping returned")

	// GetPeers should filter by region
	// Note: filters apply to peers list only
	rep, err = client.GetPeers(ctx, &peers.PeersFilter{Region: []string{"tauceti"}})
	require.NoError(err, "could not get filtered network of peers")
	require.Len(rep.Peers, 2, "peers list not filtered correctly")
	require.Equal(int64(4), rep.Status.NetworkSize, "unexpected network size returned")
	require.Len(rep.Status.Regions, 3, "unexpected regions mapping returned")

	// GetPeers should filter by multiple regions
	// Note: filters apply to peers list only
	rep, err = client.GetPeers(ctx, &peers.PeersFilter{Region: []string{"tauceti", "wolf359"}})
	require.NoError(err, "could not get multiple filtered network of peers")
	require.Len(rep.Peers, 3, "peers list not filtered correctly")
	require.Equal(int64(4), rep.Status.NetworkSize, "unexpected network size returned")
	require.Len(rep.Status.Regions, 3, "unexpected regions mapping returned")

	// GetPeers should filter by multiple duplicate regions
	// Note: filters apply to peers list only
	rep, err = client.GetPeers(ctx, &peers.PeersFilter{Region: []string{"wolf359", "wolf359", "wolf359"}})
	require.NoError(err, "could not get multiple duplicates filtered network of peers")
	require.Len(rep.Peers, 1, "peers list not filtered correctly")
	require.Equal(int64(4), rep.Status.NetworkSize, "unexpected network size returned")
	require.Len(rep.Status.Regions, 3, "unexpected regions mapping returned")

	// GetPeers should return no peers if no region filter matches
	// Note: filters apply to peers list only
	rep, err = client.GetPeers(ctx, &peers.PeersFilter{Region: []string{"not a star"}})
	require.NoError(err, "could not get no region match filtered network of peers")
	require.Len(rep.Peers, 0, "peers list not filtered correctly")
	require.Equal(int64(4), rep.Status.NetworkSize, "unexpected network size returned")
	require.Len(rep.Status.Regions, 3, "unexpected regions mapping returned")

	// GetPeers should return status only
	rep, err = client.GetPeers(ctx, &peers.PeersFilter{StatusOnly: true})
	require.NoError(err, "could not get status only")
	require.Len(rep.Peers, 0, "peers not returned after they were added")
	require.Equal(int64(4), rep.Status.NetworkSize, "unexpected network size returned")
	require.Len(rep.Status.Regions, 3, "unexpected regions mapping returned")

	// Should be able to update a peer by ID
	network[3].Addr = "wolf.trtl.star:433"
	network[3].Name = "wolfstar"
	_, err = client.AddPeers(ctx, network[3])
	require.NoError(err, "could not update peer")

	// Get peers to check the update
	rep, err = client.GetPeers(ctx, &peers.PeersFilter{Region: []string{"wolf359"}})
	require.NoError(err, "could not get filtered network of updated peers")
	require.Len(rep.Peers, 1, "peers list not filtered correctly")

	require.Equal(network[3].Id, rep.Peers[0].Id)
	require.Equal(network[3].Addr, rep.Peers[0].Addr)
	require.Equal(network[3].Name, rep.Peers[0].Name)
	require.Equal(network[3].Region, rep.Peers[0].Region)

	// Should be able to delete peers
	rmrep, err := client.RmPeers(ctx, network[0])
	require.NoError(err, "could not delete peers")
	require.Equal(int64(3), rmrep.NetworkSize)
	require.Len(rmrep.Regions, 3)

	// GetPeers to check the deletion
	rep, err = client.GetPeers(ctx, &peers.PeersFilter{})
	require.NoError(err, "could not get network of peers with peer removed")
	require.Len(rep.Peers, 3, "peers not deleted correctly")
	require.Equal(int64(3), rep.Status.NetworkSize, "unexpected network size returned")
	require.Len(rep.Status.Regions, 3, "unexpected regions mapping returned")
}
