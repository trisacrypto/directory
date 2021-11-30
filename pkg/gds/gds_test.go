package gds_test

import (
	"context"
	"time"

	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// TestRegister tests that the Register RPC correcty registers a new VASP with GDS.
func (s *gdsTestSuite) TestRegister() {
	s.LoadEmptyFixtures()
	defer s.ResetEmptyFixtures()
	defer emails.PurgeMockEmails()
	require := s.Require()
	ctx := context.Background()
	refVASP := s.fixtures[vasps]["d9da630e-41aa-11ec-9d29-acde48001122"].(*pb.VASP)

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

	// Should not be able to register an identical VASP
	_, err = client.Register(ctx, request)
	require.Error(err)
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
	request.Id = "d9da630e-41aa-11ec-9d29-acde48001122"
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
}

// TestStatus tests that the Status RPC returns the correct status response.
func (s *gdsTestSuite) TestStatus() {
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
