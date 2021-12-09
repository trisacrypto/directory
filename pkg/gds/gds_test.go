package gds_test

import (
	"context"
	"time"

	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

// TestRegister tests that the Register RPC correcty registers a new VASP with GDS.
func (s *gdsTestSuite) TestRegister() {
	s.LoadEmptyFixtures()
	defer s.ResetEmptyFixtures()
	defer emails.PurgeMockEmails()
	require := s.Require()
	ctx := context.Background()
	refVASP := s.fixtures[vasps]["charliebank"].(*pb.VASP)

	// Start the gRPC client
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// Emails need to be filled in for a valid VASP registration. Need to make copies
	// of the contacts here to avoid modifying the fixtures for other tests.
	contacts := *refVASP.Contacts
	admin := *refVASP.Contacts.Administrative
	contacts.Administrative = &admin
	contacts.Administrative.Email = "admin@example.com"
	billing := *refVASP.Contacts.Billing
	contacts.Billing = &billing
	contacts.Billing.Email = "billing@example.com"
	legal := *refVASP.Contacts.Legal
	contacts.Legal = &legal
	contacts.Legal.Email = "legal@example.com"
	technical := *refVASP.Contacts.Technical
	contacts.Technical = &technical
	contacts.Technical.Email = "technical@example.com"

	// Common name is not parseable from the endpoint
	request := &api.RegisterRequest{
		Entity:           refVASP.Entity,
		Contacts:         &contacts,
		Website:          refVASP.Website,
		BusinessCategory: refVASP.BusinessCategory,
		VaspCategories:   refVASP.VaspCategories,
		EstablishedOn:    refVASP.EstablishedOn,
		Trixo:            refVASP.Trixo,
		TrisaEndpoint:    ":3000",
	}
	_, err := client.Register(ctx, request)
	require.Error(err)
	request.TrisaEndpoint = "trisatest.net"
	_, err = client.Register(ctx, request)
	require.Error(err)

	// VASP request is incomplete
	request.TrisaEndpoint = "trisatest.net:3000"
	request.Entity = nil
	_, err = client.Register(ctx, request)
	require.Error(err)

	// Successful VASP registration
	request.Entity = refVASP.Entity
	reply, err := client.Register(ctx, request)
	require.NoError(err)
	require.NotNil(reply)
	require.NotEmpty(reply.Id)
	require.Equal(s.svc.GetConf().DirectoryID, reply.RegisteredDirectory)
	require.Equal("trisatest.net", reply.CommonName)
	require.Equal(pb.VerificationState_SUBMITTED, reply.Status)
	require.Contains(reply.Message, "verification code has been sent")
	require.NotEmpty(reply.Pkcs12Password)
	// VASP should be in the database
	v, err := s.svc.GetStore().RetrieveVASP(reply.Id)
	require.NoError(err)
	require.Equal(reply.Id, v.Id)
	require.Equal(pb.VerificationState_SUBMITTED, v.VerificationStatus)
	// Emails should be sent to the contacts
	emails, err := models.GetEmailLog(v.Contacts.Administrative)
	require.NoError(err)
	require.Len(emails, 1)
	emails, err = models.GetEmailLog(v.Contacts.Billing)
	require.NoError(err)
	require.Len(emails, 1)
	emails, err = models.GetEmailLog(v.Contacts.Legal)
	require.NoError(err)
	require.Len(emails, 1)
	emails, err = models.GetEmailLog(v.Contacts.Technical)
	require.NoError(err)
	require.Len(emails, 1)
	// Certificate request should be created
	ids, err := models.GetCertReqIDs(v)
	require.NoError(err)
	require.Len(ids, 1)
	certReq, err := s.svc.GetStore().RetrieveCertReq(ids[0])
	require.NoError(err)
	require.Equal(v.Id, certReq.Vasp)
	require.Equal(v.CommonName, certReq.CommonName)
	require.Equal(models.CertificateRequestState_INITIALIZED, certReq.Status)
	// Audit log should contain SUBMITTED entry
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 1)
	require.Equal(pb.VerificationState_SUBMITTED, log[0].CurrentState)
	// Audit log prioritizes Technical contact as the source
	require.Equal(v.Contacts.Technical.Email, log[0].Source)
	// Should not be able to register an identical VASP
	_, err = client.Register(ctx, request)
	require.Error(err)
}

// TestLookup test that the Lookup RPC correctly returns details for a VASP.
func (s *gdsTestSuite) TestLookup() {
	s.LoadFullFixtures()
	require := s.Require()
	ctx := context.Background()

	charlieVASP := s.fixtures[vasps]["charliebank"].(*pb.VASP)

	// Start the gRPC client
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// Supplied VASP ID does not exist
	request := &api.LookupRequest{
		Id: "abc12345-41aa-11ec-9d29-acde48001122",
	}
	_, err := client.Lookup(ctx, request)
	require.Error(err)

	expected := &api.LookupReply{
		Id:                  charlieVASP.Id,
		RegisteredDirectory: charlieVASP.RegisteredDirectory,
		CommonName:          charlieVASP.CommonName,
		Endpoint:            charlieVASP.TrisaEndpoint,
		IdentityCertificate: charlieVASP.IdentityCertificate,
		Country:             charlieVASP.Entity.CountryOfRegistration,
		VerifiedOn:          charlieVASP.VerifiedOn,
		Name:                "CharlieBank",
	}

	// VASP exists in the database
	request.Id = charlieVASP.Id
	reply, err := client.Lookup(ctx, request)
	require.NoError(err)
	require.True(proto.Equal(expected, reply))
	// TODO: Check that a signing certificate is returned

	// Supplied Common Name does not exist
	request.Id = ""
	request.CommonName = "invalid.name"
	_, err = client.Lookup(ctx, request)
	require.Error(err)
}

// TestSearch tests that the Search RPC returns the correct search results.
func (s *gdsTestSuite) TestSearch() {
	s.LoadFullFixtures()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// No search criteria - should not return anything
	request := &api.SearchRequest{}
	reply, err := client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 0)

	// Search by name
	request.Name = []string{"CharlieBank"}
	reply, err = client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 1)
	charlieVASP := s.fixtures[vasps]["charliebank"].(*pb.VASP)
	require.Equal(charlieVASP.Id, reply.Results[0].Id)
	require.Equal(charlieVASP.RegisteredDirectory, reply.Results[0].RegisteredDirectory)
	require.Equal(charlieVASP.CommonName, reply.Results[0].CommonName)
	require.Equal(charlieVASP.TrisaEndpoint, reply.Results[0].Endpoint)

	// Multiple results
	request.Name = []string{"CharlieBank", "Delta Assets"}
	reply, err = client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 2)

	// Search by website
	request = &api.SearchRequest{
		Website: []string{"https://trisa.charliebank.io"},
	}
	reply, err = client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 1)

	// Filter by country
	request = &api.SearchRequest{
		Name:    []string{"CharlieBank"},
		Country: []string{charlieVASP.Entity.CountryOfRegistration},
	}
	reply, err = client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 1)

	// Filter by country - no results
	request = &api.SearchRequest{
		Name:    []string{"CharlieBank"},
		Country: []string{"US"},
	}
	reply, err = client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 0)

	// Filter by category
	request = &api.SearchRequest{
		Name:             []string{"CharlieBank"},
		BusinessCategory: []pb.BusinessCategory{charlieVASP.BusinessCategory},
	}
	reply, err = client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 1)

	// Filter by business category - no results
	request = &api.SearchRequest{
		Name:             []string{"CharlieBank"},
		BusinessCategory: []pb.BusinessCategory{pb.BusinessCategory_GOVERNMENT_ENTITY},
	}
	reply, err = client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 0)

	// Filter by VASP category
	request = &api.SearchRequest{
		Name:         []string{"CharlieBank"},
		VaspCategory: []string{"P2P"},
	}
	reply, err = client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 1)

	// Filter by VASP category - no results
	request = &api.SearchRequest{
		Name:         []string{"CharlieBank"},
		VaspCategory: []string{"Project"},
	}
	reply, err = client.Search(ctx, request)
	require.NoError(err)
	require.Empty(reply.Error)
	require.Len(reply.Results, 0)
}

// TestVerifyContact tests that the VerifyContact RPC correctly verifies the VASP
// against the token and sends verification emails to the admins.
func (s *gdsTestSuite) TestVerifyContact() {
	s.LoadFullFixtures()
	defer s.ResetFullFixtures()
	defer emails.PurgeMockEmails()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// VASP does not exist in the database
	request := &api.VerifyContactRequest{
		Id:    "abc12345-41aa-11ec-9d29-acde48001122",
		Token: "",
	}
	_, err := client.VerifyContact(ctx, request)
	require.Error(err)

	// Incorrect token - no verified contacts
	request.Id = s.fixtures[vasps]["charliebank"].(*pb.VASP).Id
	request.Token = "invalid"
	_, err = client.VerifyContact(ctx, request)
	require.Error(err)

	// TODO: Test previously verified contact - requires modifying the fixtures to
	// include a non-empty verification token

	// Successful verification
	request.Token = ""
	reply, err := client.VerifyContact(ctx, request)
	require.NoError(err)
	require.Nil(reply.Error)
	require.Equal(pb.VerificationState_PENDING_REVIEW, reply.Status)
	require.Contains(reply.Message, "successfully verified")

	// VASP on the database should be updated
	vasp, err := s.svc.GetStore().RetrieveVASP(request.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_PENDING_REVIEW, vasp.VerificationStatus)
	token, err := models.GetAdminVerificationToken(vasp)
	require.NoError(err)
	require.NotEmpty(token)

	// Email should be sent to the admins
	require.Len(emails.MockEmails, 1)

	// Audit log should contain new entries for contact verifications, EMAIL_VERIFIED,
	// PENDING_REVIEW, along with the intitial SUBMITTED.
	log, err := models.GetAuditLog(vasp)
	require.NoError(err)
	// Currently verifies all contacts because the fixtures all have the empty token.
	require.Len(log, 7)
	require.Equal(pb.VerificationState_SUBMITTED, log[0].CurrentState)
	require.Equal(pb.VerificationState_SUBMITTED, log[1].CurrentState)
	require.Equal(vasp.Contacts.Technical.Email, log[1].Source)
	require.Equal(pb.VerificationState_SUBMITTED, log[2].CurrentState)
	require.Equal(vasp.Contacts.Administrative.Email, log[2].Source)
	require.Equal(pb.VerificationState_EMAIL_VERIFIED, log[5].CurrentState)
	require.Equal(pb.VerificationState_PENDING_REVIEW, log[6].CurrentState)
}

// TestVerification tests that the Verification RPC returns the correct status
// information for a VASP.
func (s *gdsTestSuite) TestVerification() {
	s.LoadFullFixtures()
	require := s.Require()
	ctx := context.Background()

	charlieID := s.fixtures[vasps]["charliebank"].(*pb.VASP).Id

	// Start the gRPC client
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// The reference fixture doesn't contain the updated timestamp, so we retrieve the
	// real VASP object here for comparison purposes.
	vasp, err := s.svc.GetStore().RetrieveVASP(charlieID)
	require.NoError(err)

	// Supplied VASP ID does not exist
	request := &api.VerificationRequest{
		Id: "abc12345-41aa-11ec-9d29-acde48001122",
	}
	_, err = client.Verification(ctx, request)
	require.Error(err)

	expected := &api.VerificationReply{
		VerificationStatus: vasp.VerificationStatus,
		ServiceStatus:      vasp.ServiceStatus,
		VerifiedOn:         vasp.VerifiedOn,
		FirstListed:        vasp.FirstListed,
		LastUpdated:        vasp.LastUpdated,
	}

	// VASP exists in the database
	request.Id = charlieID
	reply, err := client.Verification(ctx, request)
	require.NoError(err)
	require.True(proto.Equal(expected, reply))

	// Supplied Common Name does not exist
	request.Id = ""
	request.CommonName = "invalid.name"
	_, err = client.Verification(ctx, request)
	require.Error(err)

	// No VASP ID or Common Name supplied
	request.Id = ""
	request.CommonName = ""
	_, err = client.Verification(ctx, request)
	require.Error(err)
}

// TestStatus tests that the Status RPC returns the correct status response.
func (s *gdsTestSuite) TestStatus() {
	s.LoadEmptyFixtures()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// Normal health check, not in maintenance mode.
	expectedNotBefore := time.Now().Add(30 * time.Minute)
	expectedNotAfter := time.Now().Add(60 * time.Minute)
	status, err := client.Status(ctx, &api.HealthCheck{})
	require.NoError(err)
	require.Equal(api.ServiceState_HEALTHY, status.Status)

	// Timestamps should be close to expected.
	notBefore, err := time.Parse(time.RFC3339, status.NotBefore)
	require.NoError(err)
	require.True(notBefore.Sub(expectedNotBefore) < time.Minute)
	notAfer, err := time.Parse(time.RFC3339, status.NotAfter)
	require.NoError(err)
	require.True(notAfer.Sub(expectedNotAfter) < time.Minute)
}

// TestStatus tests that the Status RPC returns the correct status response when in
// maintenance mode.
func (s *gdsTestSuite) TestStatusMaintenance() {
	conf := gds.MockConfig()
	conf.Maintenance = true
	s.SetConfig(conf)
	defer s.ResetConfig()
	s.LoadEmptyFixtures()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect())
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// Health check in maintenance mode.
	expectedNotBefore := time.Now().Add(30 * time.Minute)
	expectedNotAfter := time.Now().Add(60 * time.Minute)
	status, err := client.Status(ctx, &api.HealthCheck{})
	require.NoError(err)
	require.Equal(api.ServiceState_MAINTENANCE, status.Status)

	// Timestamps should be close to expected.
	notBefore, err := time.Parse(time.RFC3339, status.NotBefore)
	require.NoError(err)
	require.True(notBefore.Sub(expectedNotBefore) < time.Minute)
	notAfter, err := time.Parse(time.RFC3339, status.NotAfter)
	require.NoError(err)
	require.True(notAfter.Sub(expectedNotAfter) < time.Minute)
}
