package gds_test

import (
	"context"

	pb "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"google.golang.org/grpc/codes"
)

func (s *gdsTestSuite) TestList() {
	// Remember to call SetupMembers after LoadFixtures!
	s.LoadFullFixtures()
	s.SetupMembers()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := pb.NewTRISAMembersClient(s.grpc.Conn)
	require.NotNil(client)

	// Test defaults, expecting 5 results
	out, err := client.List(ctx, &pb.ListRequest{})
	require.NoError(err, "default list request failed")
	require.Len(out.Vasps, 5, "unexpected vasp count from List; have the fixtures changed?")
	require.Empty(out.NextPageToken, "a next page token was returned for a one page response")

	// Test with a page size, expecting 1 result and a next page token
	out, err = client.List(ctx, &pb.ListRequest{PageSize: 1})
	require.NoError(err, "default list request failed")
	require.Len(out.Vasps, 1, "too many vasps returned from List")
	require.NotEmpty(out.NextPageToken, "no next page token was returned")

	// Test invalid page cursor
	_, err = client.List(ctx, &pb.ListRequest{PageToken: "123"})
	require.Error(err)
	s.StatusError(err, codes.InvalidArgument, "invalid page token")

	// Changing page size between requests results in an error
	token := "CAISLHBlb3BsZTo6NDZlNzg5MTctOGQyMC00N2MwLWIwZDEtZTUyMDQxNDlhOTM2"
	_, err = client.List(ctx, &pb.ListRequest{PageToken: token, PageSize: 27})
	s.StatusError(err, codes.InvalidArgument, "page size cannot change between requests")

}
