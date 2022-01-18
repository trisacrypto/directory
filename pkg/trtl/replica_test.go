package trtl_test

import (
	"context"

	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
)

func (s *trtlTestSuite) TestSelectPeer() {
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client and replica service
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	replica, err := trtl.NewReplicaService(s.trtl)
	require.NoError(err)
	client := peers.NewPeerManagementClient(s.grpc.Conn)

	// Add a peer to the database to represent "self"
	self := &peers.Peer{
		Id:     s.conf.Replica.PID,
		Addr:   "localhost",
		Name:   s.conf.Replica.Name,
		Region: s.conf.Replica.Region,
	}
	_, err = client.AddPeers(ctx, self)
	require.NoError(err)

	// Initially there should be no peers to select from
	require.Nil(replica.SelectPeer())

	// Add a peer to the network
	peer := &peers.Peer{
		Id:     1,
		Addr:   "tc1.trtl.dev:443",
		Name:   "alpha",
		Region: "tauceti",
	}
	_, err = client.AddPeers(ctx, peer)
	require.NoError(err)

	// Select the only peer available
	selected := replica.SelectPeer()
	require.NotNil(selected)
	require.Equal(peer.Id, selected.Id)
	require.Equal(peer.Addr, selected.Addr)
	require.Equal(peer.Name, selected.Name)
	require.Equal(peer.Region, selected.Region)

	// Add some more peers to the network
	network := []*peers.Peer{
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
		_, err = client.AddPeers(ctx, peer)
		require.NoError(err)
	}

	// Select a peer that is not self
	selected = replica.SelectPeer()
	require.NotNil(selected)
	require.NotEqual(peer.Id, self.Id)
}
