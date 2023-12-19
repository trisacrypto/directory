package gds_test

import (
	"context"

	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func (s *gdsTestSuite) TestMembersList() {
	// Remember to call SetupMembers after LoadFixtures!
	s.LoadFullFixtures()
	s.SetupMembers()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := members.NewTRISAMembersClient(s.grpc.Conn)
	require.NotNil(client)

	// Test defaults, expecting 5 results
	out, err := client.List(ctx, &members.ListRequest{})
	require.NoError(err, "default list request failed")
	require.Len(out.Vasps, 5, "unexpected vasp count from List; have the fixtures changed?")
	require.Empty(out.NextPageToken, "a next page token was returned for a one page response")

	// Test with a page size, expecting 1 result and a next page token
	out, err = client.List(ctx, &members.ListRequest{PageSize: 1})
	require.NoError(err, "default list request failed")
	require.Len(out.Vasps, 1, "too many vasps returned from List")
	require.NotEmpty(out.NextPageToken, "no next page token was returned")

	// Test invalid page cursor
	_, err = client.List(ctx, &members.ListRequest{PageToken: "123"})
	require.Error(err)
	s.StatusError(err, codes.InvalidArgument, "invalid page token")

	// Changing page size between requests results in an error
	token := "CAISLHBlb3BsZTo6NDZlNzg5MTctOGQyMC00N2MwLWIwZDEtZTUyMDQxNDlhOTM2"
	_, err = client.List(ctx, &members.ListRequest{PageToken: token, PageSize: 27})
	s.StatusError(err, codes.InvalidArgument, "page size cannot change between requests")
}

func (s *gdsTestSuite) TestMembersSummary() {
	s.LoadFullFixtures()
	s.SetupMembers()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := members.NewTRISAMembersClient(s.grpc.Conn)
	require.NotNil(client)

	// Test with default parameters
	out, err := client.Summary(ctx, &members.SummaryRequest{})
	require.NoError(err, "default summary request failed")
	require.Nil(out.MemberInfo, "VASP should not be populated if not provided in request")
	require.Equal(int32(5), out.Vasps, "unexpected total vasp count from Summary; have the fixtures changed?")
	require.Equal(int32(5), out.CertificatesIssued, "unexpected certificates issued count from Summary; have the fixtures changed?")
	require.Equal(int32(0), out.NewMembers, "unexpected new members count from Summary; have the fixtures changed?")

	// Test retrieving VASP details
	charlie, _, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err)
	name, err := charlie.Name()
	require.NoError(err)
	details := &members.VASPMember{
		Id:                  charlie.Id,
		RegisteredDirectory: charlie.RegisteredDirectory,
		CommonName:          charlie.CommonName,
		Endpoint:            charlie.TrisaEndpoint,
		Name:                name,
		Website:             charlie.Website,
		Country:             charlie.Entity.CountryOfRegistration,
		BusinessCategory:    charlie.BusinessCategory,
		VaspCategories:      charlie.VaspCategories,
		VerifiedOn:          charlie.VerifiedOn,
		Status:              charlie.VerificationStatus,
		FirstListed:         charlie.FirstListed,
		LastUpdated:         charlie.LastUpdated,
	}
	out, err = client.Summary(ctx, &members.SummaryRequest{
		MemberId: charlie.Id,
	})
	require.NoError(err, "summary request with VASP failed")
	require.True(proto.Equal(details, out.MemberInfo), "VASP details mismatch")

	// Test with a non-existent VASP
	_, err = client.Summary(ctx, &members.SummaryRequest{
		MemberId: "invalid",
	})
	s.StatusError(err, codes.NotFound, "requested VASP not found")

	// Test with a since timestamp
	out, err = client.Summary(ctx, &members.SummaryRequest{
		Since: "2021-06-01T00:00:00Z",
	})
	require.NoError(err, "summary request with since failed")
	require.Equal(int32(5), out.Vasps, "unexpected total vasp count from Summary; have the fixtures changed?")
	require.Equal(int32(5), out.CertificatesIssued, "unexpected certificates issued count from Summary; have the fixtures changed?")
	require.Equal(int32(5), out.NewMembers, "unexpected new members count from Summary; have the fixtures changed?")

	// Test with an invalid since timestamp
	_, err = client.Summary(ctx, &members.SummaryRequest{
		Since: "not a timestamp",
	})
	s.StatusError(err, codes.InvalidArgument, "since must be a valid RFC3339 timestamp")

	// Test with an out of range timestamp
	_, err = client.Summary(ctx, &members.SummaryRequest{
		Since: "2063-04-05T00:00:00Z",
	})
	s.StatusError(err, codes.InvalidArgument, "since timestamp must be in the past")
}

func (s *gdsTestSuite) TestMembersDetails() {
	s.LoadFullFixtures()
	s.SetupMembers()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := members.NewTRISAMembersClient(s.grpc.Conn)
	require.NotNil(client)

	// Test with a non-existent VASP
	_, err := client.Details(ctx, &members.DetailsRequest{
		MemberId: "invalid",
	})
	s.StatusError(err, codes.NotFound, "requested VASP not found")

	// Test with a valid VASP
	charlie, _, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err, "could not get charliebank VASP")
	name, err := charlie.Name()
	require.NoError(err)
	details := &members.VASPMember{
		Id:                  charlie.Id,
		RegisteredDirectory: charlie.RegisteredDirectory,
		CommonName:          charlie.CommonName,
		Endpoint:            charlie.TrisaEndpoint,
		Name:                name,
		Website:             charlie.Website,
		Country:             charlie.Entity.CountryOfRegistration,
		BusinessCategory:    charlie.BusinessCategory,
		VaspCategories:      charlie.VaspCategories,
		VerifiedOn:          charlie.VerifiedOn,
		Status:              charlie.VerificationStatus,
		FirstListed:         charlie.FirstListed,
		LastUpdated:         charlie.LastUpdated,
	}
	out, err := client.Details(ctx, &members.DetailsRequest{
		MemberId: charlie.Id,
	})
	require.NoError(err, "details request with VASP failed")
	require.True(proto.Equal(details, out.MemberSummary), "VASP details mismatch")
	require.True(proto.Equal(charlie.Entity, out.LegalPerson), "VASP legal person mismatch")
	require.True(proto.Equal(charlie.Trixo, out.Trixo), "VASP trixo form mismatch")

	// Check contacts return is correct
	contacts := []struct {
		expected *pb.Contact
		actual   *pb.Contact
		name     string
	}{
		{charlie.Contacts.Administrative, out.Contacts.Administrative, "administrative"},
		{charlie.Contacts.Technical, out.Contacts.Technical, "technical"},
		{charlie.Contacts.Legal, out.Contacts.Legal, "legal"},
		{charlie.Contacts.Billing, out.Contacts.Billing, "billing"},
	}

	nContacts := 0
	for _, contact := range contacts {
		if contact.expected == nil {
			require.Nil(contact.actual, "expected nil entry for %s contact", contact.name)
			continue
		}

		require.NotNil(contact.actual, "expected non-nil entry for %s contact", contact.name)
		nContacts++

		// Ensure that private information is not returned
		require.Nil(contact.actual.Person)
		require.Nil(contact.actual.Extra)

		// Ensure that contact matches
		require.Equal(contact.expected.Name, contact.actual.Name, "mismatch name for %s contact", contact.name)
		require.Equal(contact.expected.Email, contact.actual.Email, "mismatch email for %s contact", contact.name)
		require.Equal(contact.expected.Phone, contact.actual.Phone, "mismatch phone for %s contact", contact.name)
	}

	require.Greater(nContacts, 0, "charlie fixture has no contact data")
}
