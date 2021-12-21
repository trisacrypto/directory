package gds_test

// func (s *gdsTestSuite) TestList() {
// // Remember to call SetupMembers after LoadFixtures!
// 	s.LoadFullFixtures()
// 	require := s.Require()
// 	// ctx := context.Background()

// 	// Start the gRPC client.
// 	require.NoError(s.grpc.Connect())
// 	defer s.grpc.Close()
// 	client := pb.NewTRISAMembersClient(s.grpc.Conn)
// 	require.NotNil(client)

// 	// Empty request should return an error
// 	// _, err := client.List(ctx, &pb.ListRequest{})
// 	// require.EqualError(err, "could not iterate over VASPs")
// 	// s.StatusError(err, codes.Internal, "could not iterate over directory service")
// }
