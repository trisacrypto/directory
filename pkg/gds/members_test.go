package gds_test

import (
	pb "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
)

func (s *gdsTestSuite) TestNewMembers() {
	require := s.Require()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := pb.NewTRISAMembersClient(s.grpc.Conn)
	require.NotNil(client)
}

func (s *gdsTestSuite) TestServe() {}

func (s *gdsTestSuite) TestList() {
	s.LoadFullFixtures()
	require := s.Require()
	// ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := pb.NewTRISAMembersClient(s.grpc.Conn)
	require.NotNil(client)

	// Empty request should return an error
	// _, err := client.List(ctx, &pb.ListRequest{})
	// require.EqualError(err, "could not iterate over VASPs")
	// s.StatusError(err, codes.Internal, "could not iterate over directory service")
}
