package trtl_test

import (
	"context"
	"crypto/x509/pkix"
	"net"

	"github.com/trisacrypto/directory/pkg/trtl"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
)

// TestPeerFromTLS tests that the peerFromTLS function is able to extract peer info
// from TLS within the gRPC context.
func (s *trtlTestSuite) TestPeerFromTLS() {
	require := s.Require()

	// TODO: We probably want to test more variations of the certificate.
	expected := &trtl.PeerInfo{
		Name: &pkix.Name{
			Country:       []string{"US"},
			Organization:  []string{"TRISA Development Client"},
			Locality:      []string{"Raleigh"},
			Province:      []string{"North Carolina"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
			SerialNumber:  "", // TODO: Serial number was not parsed?
			CommonName:    clientTarget,
		},
		DNSNames:    []string{clientTarget},
		IPAddresses: []net.IP{},
	}

	// Get the dial options for the client connection
	opts, err := s.loadClientCredentials()
	require.NoError(err, "could not load client credentials")

	// Create a database client to the remote trtl peer
	dbClient, err := s.remote.DBClient(opts...)
	require.NoError(err, "could not create trtl database client")
	defer s.remote.CloseClient()

	// Create peer info that the RPC can write to
	info := &trtl.PeerInfo{}

	// Configure the remote peer to extract the peer info from the context
	s.remote.OnGet = func(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
		// Extract the peer info from the context
		info, err = trtl.PeerFromTLS(ctx)
		require.NoError(err, "could not extract peer info from TLS")

		return &pb.GetReply{}, nil
	}

	// Make a request to the remote peer to trigger the test
	_, err = dbClient.Get(context.Background(), &pb.GetRequest{})
	require.NoError(err, "could not make database request to remote peer")
	s.checkPeerInfo(expected, info)

	// Create a peers client to the remote trtl peer
	peersClient, err := s.remote.PeersClient(opts...)
	require.NoError(err, "could not create trtl peers client")

	// Configure the remote peer to extract the peer info from the context
	s.remote.OnAddPeers = func(ctx context.Context, req *peers.Peer) (*peers.PeersStatus, error) {
		// Extract the peer info from the context
		info, err = trtl.PeerFromTLS(ctx)
		require.NoError(err, "could not extract peer info from TLS")

		return &peers.PeersStatus{}, nil
	}

	// Make a request to the remote peer to trigger the test
	_, err = peersClient.AddPeers(context.Background(), &peers.Peer{})
	require.NoError(err, "could not make peers request to remote peer")
	s.checkPeerInfo(expected, info)

	// TODO: Do we need to test the replica service?
}

// checkPeerInfo checks that the peer info matches the expected one
func (s *trtlTestSuite) checkPeerInfo(expected *trtl.PeerInfo, actual *trtl.PeerInfo) {
	require := s.Require()

	// Remove the RDNSequence fields for comparison
	actual.Name.Names = nil

	require.Equal(expected, actual, "peer info does not match expected")
}
