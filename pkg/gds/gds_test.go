package gds_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/emails/mock"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// StatusError is a helper assertion function that checks a gRPC status error
func (s *gdsTestSuite) StatusError(err error, code codes.Code, theError string) {
	require := s.Require()
	require.Error(err, "no status error returned")

	var serr *status.Status
	serr, ok := status.FromError(err)
	require.True(ok, "error is not a grpc status error")
	require.Equal(code, serr.Code(), "status code does not match")
	require.Equal(theError, serr.Message(), "status error message does not match")
}

// TestRegister tests that the Register RPC correctly registers a new VASP with GDS.
func (s *gdsTestSuite) TestRegister() {
	// Load the fixtures and start the GDS server
	s.LoadEmptyFixtures()
	s.SetupGDS()
	defer s.ResetFixtures()
	defer s.fixtures.LoadReferenceFixtures()
	defer mock.PurgeEmails()
	require := s.Require()
	ctx := context.Background()
	charlie, _, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err)

	// Start the gRPC client
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// Emails need to be filled in for a valid VASP registration. Note: This modifies
	// the contacts on the original fixtures so LoadReferenceFixtures() must be
	// deferred in order to restore them before the next test.
	contacts := charlie.Contacts
	if contacts.Administrative == nil {
		contacts.Administrative = &pb.Contact{}
	}
	contacts.Administrative.Name = "Admin Person"
	contacts.Administrative.Email = "admin@example.com"

	if contacts.Billing == nil {
		contacts.Billing = &pb.Contact{}
	}
	contacts.Billing.Name = "Billing Person"
	contacts.Billing.Email = "billingandlegal@example.com"

	if contacts.Legal == nil {
		contacts.Legal = &pb.Contact{}
	}
	contacts.Legal.Name = "Legal Person"
	contacts.Legal.Email = "billingandlegal@example.com"

	if contacts.Technical == nil {
		contacts.Technical = &pb.Contact{}
	}
	contacts.Technical.Name = "Technical Person"
	contacts.Technical.Email = "technical@example.com"

	// Request contains an invalid endpoint
	request := &api.RegisterRequest{
		Entity:           charlie.Entity,
		Contacts:         contacts,
		Website:          charlie.Website,
		BusinessCategory: charlie.BusinessCategory,
		VaspCategories:   charlie.VaspCategories,
		EstablishedOn:    charlie.EstablishedOn,
		Trixo:            charlie.Trixo,
	}
	_, err = client.Register(ctx, request)
	require.Error(err)
	request.TrisaEndpoint = "http://trisatest.net:443"
	_, err = client.Register(ctx, request)
	require.Error(err)
	request.TrisaEndpoint = "grpc://:443"
	_, err = client.Register(ctx, request)
	require.Error(err)
	request.TrisaEndpoint = ":443"
	_, err = client.Register(ctx, request)
	require.Error(err)
	request.TrisaEndpoint = "trisatest.net"
	_, err = client.Register(ctx, request)
	require.Error(err)
	request.TrisaEndpoint = "trisatest.net:443/"
	_, err = client.Register(ctx, request)
	require.Error(err)
	request.TrisaEndpoint = "trisatest.net:443/path"
	_, err = client.Register(ctx, request)
	require.Error(err)

	// Request contains an invalid common name
	request.TrisaEndpoint = "trisatest.net:443"
	request.CommonName = "http://trisatest.net"
	_, err = client.Register(ctx, request)
	require.Error(err)
	request.CommonName = ":443"
	_, err = client.Register(ctx, request)
	require.Error(err)
	request.CommonName = "trisatest.net:443"
	_, err = client.Register(ctx, request)
	require.Error(err)
	request.CommonName = "trisatest.net/path"
	_, err = client.Register(ctx, request)
	require.Error(err)

	// Request contains no entity
	request.CommonName = ""
	request.Entity = nil
	_, err = client.Register(ctx, request)
	require.Error(err)

	// Successful VASP registration
	request.Entity = charlie.Entity
	sent := time.Now()
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
	v, err := s.svc.GetStore().RetrieveVASP(context.Background(), reply.Id)
	require.NoError(err)
	require.Equal(reply.Id, v.Id)
	require.Equal(pb.VerificationState_SUBMITTED, v.VerificationStatus)
	// Certificate request should be created
	ids, err := models.GetCertReqIDs(v)
	require.NoError(err)
	require.Len(ids, 1)
	certReq, err := s.svc.GetStore().RetrieveCertReq(context.Background(), ids[0])
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

	// Emails should be sent to all unique contacts
	messages := []*emails.EmailMeta{
		{
			Contact:   v.Contacts.Administrative,
			To:        v.Contacts.Administrative.Email,
			From:      s.svc.GetConf().Email.ServiceEmail,
			Subject:   emails.VerifyContactRE,
			Reason:    "verify_contact",
			Timestamp: sent,
		},
		{
			Contact:   v.Contacts.Legal,
			To:        v.Contacts.Legal.Email,
			From:      s.svc.GetConf().Email.ServiceEmail,
			Subject:   emails.VerifyContactRE,
			Reason:    "verify_contact",
			Timestamp: sent,
		},
		{
			Contact:   v.Contacts.Technical,
			To:        v.Contacts.Technical.Email,
			From:      s.svc.GetConf().Email.ServiceEmail,
			Subject:   emails.VerifyContactRE,
			Reason:    "verify_contact",
			Timestamp: sent,
		},
	}
	emails.CheckEmails(s.T(), messages)
}

func (s *gdsTestSuite) TestRegisterAlreadyVerified() {
	s.T().Skip("requires updates to fixtures")

	// Load the fixtures and start the GDS server
	s.LoadSmallFixtures()
	s.SetupGDS()
	defer s.ResetFixtures()
	defer s.fixtures.LoadReferenceFixtures()
	defer mock.PurgeEmails()
	require := s.Require()
	ctx := context.Background()
	charlie, _, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err)

	// Ensure that the contact fixtures were loaded
	adam, err := s.fixtures.GetContact("adam@example.com")
	require.NoError(err, "missing contact fixture for adam@example.com")
	require.False(adam.Verified, "expected contact adam@example.com to be unverified, have the fixtures changed?")
	bruce, err := s.fixtures.GetContact("bruce@example.com")
	require.NoError(err, "missing contact fixture for bruce@example.com")
	require.True(bruce.Verified, "expected contact bruce@example.com to be verified, have the fixtures changed?")

	// Start the gRPC client
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// Emails need to be filled in for a valid VASP registration. Note: This modifies
	// the contacts on the original fixtures so LoadReferenceFixtures() must be
	// deferred in order to restore them before the next test.
	contacts := charlie.Contacts
	if contacts.Administrative == nil {
		contacts.Administrative = &pb.Contact{}
	}
	contacts.Administrative.Name = adam.Name
	contacts.Administrative.Email = adam.Email

	if contacts.Billing == nil {
		contacts.Billing = &pb.Contact{}
	}
	contacts.Billing.Name = bruce.Name
	contacts.Billing.Email = bruce.Email

	if contacts.Legal == nil {
		contacts.Legal = &pb.Contact{}
	}
	contacts.Legal.Name = "Legal Person"
	contacts.Legal.Email = "billingandlegal@example.com"

	if contacts.Technical == nil {
		contacts.Technical = &pb.Contact{}
	}
	contacts.Technical.Name = "Technical Person"
	contacts.Technical.Email = "technical@example.com"

	// Valid register request with a verified contact
	request := &api.RegisterRequest{
		Entity:           charlie.Entity,
		Contacts:         contacts,
		Website:          charlie.Website,
		BusinessCategory: charlie.BusinessCategory,
		VaspCategories:   charlie.VaspCategories,
		EstablishedOn:    charlie.EstablishedOn,
		Trixo:            charlie.Trixo,
		TrisaEndpoint:    "trisatest.net:443",
	}
	sent := time.Now()
	reply, err := client.Register(ctx, request)
	require.NoError(err)
	require.NotNil(reply)
	require.NotEmpty(reply.Id)
	require.Equal(s.svc.GetConf().DirectoryID, reply.RegisteredDirectory)
	require.Equal("trisatest.net", reply.CommonName)
	require.Equal(pb.VerificationState_PENDING_REVIEW, reply.Status)
	require.Contains(reply.Message, "verification code has been sent")
	require.NotEmpty(reply.Pkcs12Password)
	// VASP should be in the database
	v, err := s.svc.GetStore().RetrieveVASP(context.Background(), reply.Id)
	require.NoError(err)
	require.Equal(reply.Id, v.Id)
	require.Equal(pb.VerificationState_PENDING_REVIEW, v.VerificationStatus)
	// Certificate request should be created
	ids, err := models.GetCertReqIDs(v)
	require.NoError(err)
	require.Len(ids, 1)
	certReq, err := s.svc.GetStore().RetrieveCertReq(context.Background(), ids[0])
	require.NoError(err)
	require.Equal(v.Id, certReq.Vasp)
	require.Equal(v.CommonName, certReq.CommonName)
	require.Equal(models.CertificateRequestState_INITIALIZED, certReq.Status)
	// Audit log should contain SUBMITTED, EMAIL_VERIFIED, and PENDING_REVIEW
	log, err := models.GetAuditLog(v)
	require.NoError(err)
	require.Len(log, 3)
	require.Equal(pb.VerificationState_SUBMITTED, log[0].CurrentState)
	// Audit log prioritizes Technical contact as the submission source
	require.Equal(v.Contacts.Technical.Email, log[0].Source)
	// The verified email should be used for the other two entries
	require.Equal(pb.VerificationState_EMAIL_VERIFIED, log[1].CurrentState)
	require.Equal(bruce.Email, log[1].Source)
	require.Equal(pb.VerificationState_PENDING_REVIEW, log[2].CurrentState)
	require.Equal(bruce.Email, log[2].Source)

	// Should not be able to register an identical VASP
	_, err = client.Register(ctx, request)
	require.Error(err)

	// Emails should be sent to all unverified contacts
	messages := []*emails.EmailMeta{
		{
			Contact:   v.Contacts.Administrative,
			To:        v.Contacts.Administrative.Email,
			From:      s.svc.GetConf().Email.ServiceEmail,
			Subject:   emails.VerifyContactRE,
			Reason:    "verify_contact",
			Timestamp: sent,
		},
		{
			Contact:   v.Contacts.Legal,
			To:        v.Contacts.Legal.Email,
			From:      s.svc.GetConf().Email.ServiceEmail,
			Subject:   emails.VerifyContactRE,
			Reason:    "verify_contact",
			Timestamp: sent,
		},
		{
			Contact:   v.Contacts.Technical,
			To:        v.Contacts.Technical.Email,
			From:      s.svc.GetConf().Email.ServiceEmail,
			Subject:   emails.VerifyContactRE,
			Reason:    "verify_contact",
			Timestamp: sent,
		},
		{
			To:        s.svc.GetConf().Email.AdminEmail,
			From:      s.svc.GetConf().Email.ServiceEmail,
			Subject:   emails.ReviewRequestRE,
			Reason:    "review_request",
			Timestamp: sent,
		},
	}
	emails.CheckEmails(s.T(), messages)
}

// TestLookup test that the Lookup RPC correctly returns details for a VASP.
func (s *gdsTestSuite) TestLookup() {
	// Load the fixtures and start the GDS server
	s.LoadFullFixtures()
	s.SetupGDS()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	s.Run("IDNotFound", func() {
		// Supplied VASP ID does not exist
		require := s.Require()
		request := &api.LookupRequest{
			Id: "abc12345-41aa-11ec-9d29-acde48001122",
		}
		_, err := client.Lookup(ctx, request)
		require.EqualError(err, "rpc error: code = NotFound desc = could not find VASP by ID")
	})

	s.Run("CommonNameNotFound", func() {
		// Supplied Common Name does not exist
		require := s.Require()
		request := &api.LookupRequest{
			CommonName: "invalid.name",
		}
		_, err := client.Lookup(ctx, request)
		require.EqualError(err, "rpc error: code = NotFound desc = could not find VASP by common name")
	})

	s.Run("VerifiedRequired", func() {
		// The VASP must be verified or it is not found
		require := s.Require()
		charlieVASP, _, err := s.fixtures.GetVASP("charliebank")
		require.NoError(err)
		require.NotEqual(charlieVASP.VerificationStatus, pb.VerificationState_VERIFIED, "expected fixture to not be in verified state")

		request := &api.LookupRequest{
			Id: charlieVASP.Id,
		}
		_, err = client.Lookup(ctx, request)
		require.EqualError(err, "rpc error: code = NotFound desc = no VASP record available")
	})

	s.Run("Found", func() {
		// Expect lookup to succeed
		require := s.Require()
		hotelVASP, _, err := s.fixtures.GetVASP("hotel")
		require.NoError(err)
		require.Equal(hotelVASP.VerificationStatus, pb.VerificationState_VERIFIED, "expected fixture to be in verified state")

		expected := &api.LookupReply{
			Id:                  hotelVASP.Id,
			RegisteredDirectory: hotelVASP.RegisteredDirectory,
			CommonName:          hotelVASP.CommonName,
			Endpoint:            hotelVASP.TrisaEndpoint,
			SigningCertificate:  hotelVASP.SigningCertificates[len(hotelVASP.SigningCertificates)-1],
			IdentityCertificate: hotelVASP.IdentityCertificate,
			Country:             hotelVASP.Entity.CountryOfRegistration,
			VerifiedOn:          hotelVASP.VerifiedOn,
			Name:                "Hotel Corp",
		}

		// VASP exists in the database
		request := &api.LookupRequest{
			Id: hotelVASP.Id,
		}
		reply, err := client.Lookup(ctx, request)
		require.NoError(err)
		require.Equal(expected.Name, reply.Name)
		require.True(proto.Equal(expected, reply))
	})
}

// TestSearch tests that the Search RPC returns the correct search results.
func (s *gdsTestSuite) TestSearch() {
	// Load the fixtures and start the GDS server
	s.LoadFullFixtures()
	s.SetupGDS()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// Collect some fixtures
	charlieVASP, _, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err)
	require.NotEqual(charlieVASP.VerificationStatus, pb.VerificationState_VERIFIED, "expected fixture to not be in verified state")

	hotelVASP, _, err := s.fixtures.GetVASP("hotel")
	require.NoError(err)
	require.Equal(hotelVASP.VerificationStatus, pb.VerificationState_VERIFIED, "expected fixture to be in verified state")

	novemberVASP, _, err := s.fixtures.GetVASP("novembercash")
	require.NoError(err)
	require.Equal(novemberVASP.VerificationStatus, pb.VerificationState_VERIFIED, "expected fixture to be in verified state")

	s.Run("Empty", func() {
		// No search criteria - should not return anything
		require := s.Require()
		request := &api.SearchRequest{}
		reply, err := client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 0)
	})

	s.Run("ByName", func() {
		// Search by name
		require := s.Require()
		request := &api.SearchRequest{
			Name: []string{"Hotel Corp"},
		}

		reply, err := client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 1)

		require.Equal(hotelVASP.Id, reply.Results[0].Id)
		require.Equal(hotelVASP.RegisteredDirectory, reply.Results[0].RegisteredDirectory)
		require.Equal(hotelVASP.CommonName, reply.Results[0].CommonName)
		require.Equal(hotelVASP.TrisaEndpoint, reply.Results[0].Endpoint)
	})

	s.Run("Verified", func() {
		// Only verified results are returned
		require := s.Require()
		request := &api.SearchRequest{
			Name: []string{"CharlieBank", "Hotel Corp"},
		}

		reply, err := client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 1)
		require.Equal(hotelVASP.Id, reply.Results[0].Id)
	})

	s.Run("FuzzyPrefix", func() {
		// Fuzzy search by case-insensitive prefix
		require := s.Require()
		request := &api.SearchRequest{
			Name: []string{"NOV"},
		}

		reply, err := client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 1)
		require.Equal(novemberVASP.Id, reply.Results[0].Id)

		// Prefix search must have at least three characters
		request.Name = []string{"no"}
		reply, err = client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 0)
	})

	s.Run("MultipleResults", func() {
		// Multiple results
		require := s.Require()
		request := &api.SearchRequest{
			Name: []string{"Hotel Corp", "November"},
		}

		reply, err := client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 2)
	})

	s.Run("Website", func() {
		// Search by website
		require := s.Require()
		request := &api.SearchRequest{
			Website: []string{"https://trisa.hotel.io"},
		}
		reply, err := client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 1)
	})

	s.Run("Country", func() {
		// Filter by country
		require := s.Require()
		request := &api.SearchRequest{
			Name:    []string{"Hotel Corp"},
			Country: []string{hotelVASP.Entity.CountryOfRegistration},
		}
		reply, err := client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 1)

		// Filter by country - no results
		request = &api.SearchRequest{
			Name:    []string{"Hotel Corp"},
			Country: []string{"GY"},
		}
		reply, err = client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 0)
	})

	s.Run("Category", func() {
		// Filter by category
		require := s.Require()
		request := &api.SearchRequest{
			Name:             []string{"Hotel Corp"},
			BusinessCategory: []pb.BusinessCategory{hotelVASP.BusinessCategory},
		}
		reply, err := client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 1)

		// Filter by business category - no results
		request = &api.SearchRequest{
			Name:             []string{"Hotel Corp"},
			BusinessCategory: []pb.BusinessCategory{pb.BusinessCategory_GOVERNMENT_ENTITY},
		}
		reply, err = client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 0)
	})

	s.Run("VASPCategory", func() {
		// Filter by VASP category
		require := s.Require()
		request := &api.SearchRequest{
			Name:         []string{"Hotel Corp"},
			VaspCategory: []string{"Mixer"},
		}
		reply, err := client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 1)

		// Filter by VASP category - no results
		request = &api.SearchRequest{
			Name:         []string{"Hotel Corp"},
			VaspCategory: []string{"Project"},
		}
		reply, err = client.Search(ctx, request)
		require.NoError(err)
		require.Empty(reply.Error)
		require.Len(reply.Results, 0)
	})
}

// TestVerifyContact tests that the VerifyContact RPC correctly verifies the VASP
// against the token and sends verification emails to the admins.
func (s *gdsTestSuite) TestVerifyContact() {
	// Load the fixtures and start the GDS server
	s.LoadFullFixtures()
	s.SetupGDS()
	defer s.ResetFixtures()
	defer s.fixtures.LoadReferenceFixtures()
	defer mock.PurgeEmails()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	charlie, contacts, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err)

	// Cannot verify contact without a token
	request := &api.VerifyContactRequest{
		Id: charlie.Id,
	}
	_, err = client.VerifyContact(ctx, request)
	require.Error(err)

	// VASP does not exist in the database
	request = &api.VerifyContactRequest{
		Id:    "abc12345-41aa-11ec-9d29-acde48001122",
		Token: "administrative_token",
	}
	_, err = client.VerifyContact(ctx, request)
	require.Error(err)

	// Incorrect token - no verified contacts
	request.Id = charlie.Id
	request.Token = "invalid"
	_, err = client.VerifyContact(ctx, request)
	require.Error(err)

	iter := contacts.NewIterator()
	for iter.Next() {
		contact := iter.Contact()
		s.svc.GetStore().CreateContact(ctx, &models.Contact{
			Email: contact.Email.Email,
			Token: "administrative_token",
		})
	}

	// Successful verification
	request.Token = "administrative_token"
	sent := time.Now()
	reply, err := client.VerifyContact(ctx, request)
	require.NoError(err)
	require.Nil(reply.Error)
	require.Equal(pb.VerificationState_PENDING_REVIEW, reply.Status)
	require.Contains(reply.Message, "successfully verified")
	// VASP on the database should be updated
	vasp, err := s.svc.GetStore().RetrieveVASP(context.Background(), request.Id)
	require.NoError(err)
	require.Equal(pb.VerificationState_PENDING_REVIEW, vasp.VerificationStatus)
	token, err := models.GetAdminVerificationToken(vasp)
	require.NoError(err)
	require.NotEmpty(token)

	// Attempt to verify an already verified contact
	request.Token = "administrative_token"
	_, err = client.VerifyContact(ctx, request)
	require.NoError(err)

	// Check audit log entries
	log, err := models.GetAuditLog(vasp)
	require.NoError(err)
	require.Len(log, 4)
	// Pre-existing entry for SUBMITTED
	require.Equal(pb.VerificationState_SUBMITTED, log[0].CurrentState)
	// Administrative contact verified
	require.Equal(pb.VerificationState_SUBMITTED, log[1].CurrentState)
	require.Equal(vasp.Contacts.Administrative.Email, log[1].Source)
	// State of the VASP changes to EMAIL_VERIFIED then PENDING_REVIEW
	require.Equal(pb.VerificationState_EMAIL_VERIFIED, log[2].CurrentState)
	require.Equal(pb.VerificationState_PENDING_REVIEW, log[3].CurrentState)

	// Only one email should be sent to the admins
	messages := []*emails.EmailMeta{
		{
			To:        s.svc.GetConf().Email.AdminEmail,
			From:      s.svc.GetConf().Email.ServiceEmail,
			Subject:   emails.ReviewRequestRE,
			Timestamp: sent,
		},
	}
	emails.CheckEmails(s.T(), messages)
}

// TestVerification tests that the Verification RPC returns the correct status
// information for a VASP.
func (s *gdsTestSuite) TestVerification() {
	// Load the fixtures and start the GDS server
	s.LoadFullFixtures()
	s.SetupGDS()
	require := s.Require()
	ctx := context.Background()

	charlie, _, err := s.fixtures.GetVASP("charliebank")
	require.NoError(err)

	// Start the gRPC client
	require.NoError(s.grpc.Connect(ctx))
	defer s.grpc.Close()
	client := api.NewTRISADirectoryClient(s.grpc.Conn)

	// The reference fixture doesn't contain the updated timestamp, so we retrieve the
	// real VASP object here for comparison purposes.
	vasp, err := s.svc.GetStore().RetrieveVASP(context.Background(), charlie.Id)
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
	request.Id = charlie.Id
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
	// Load the fixtures and start the GDS server
	s.LoadEmptyFixtures()
	s.SetupGDS()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect(ctx))
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
	notAfter, err := time.Parse(time.RFC3339, status.NotAfter)
	require.NoError(err)
	require.True(notAfter.Sub(expectedNotAfter) < time.Minute)
}

// TestStatus tests that the Status RPC returns the correct status response when in
// maintenance mode.
func (s *gdsTestSuite) TestStatusMaintenance() {
	conf := gds.MockConfig()
	conf.Maintenance = true
	s.SetConfig(conf)
	defer s.ResetConfig()

	// Load the fixtures and start the GDS server
	s.LoadEmptyFixtures()
	s.SetupGDS()
	require := s.Require()
	ctx := context.Background()

	// Start the gRPC client.
	require.NoError(s.grpc.Connect(ctx))
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

// Test Common Name Validation
func TestValidateCommonName(t *testing.T) {
	var (
		noMatch     = errors.New("common name does not match domain name regular expression")
		empty       = errors.New("common name should not be empty")
		noWildcards = errors.New("wildcards are not allowed in TRISA common names")
	)

	testCases := []struct {
		input    string
		expected error
	}{
		{"trisa.example.com", nil},
		{"subdomain.trisa-testnet.example.com", nil},
		{"example.com", nil},
		{"localhost", nil},
		{"", empty},
		{"*.example.com", noWildcards},
		{"-foo.example.com", noMatch},
		{"https://trisa.example.com", noMatch},
		{"trisa.example.com:443", noMatch},
		{"  trisa.example.com   ", noMatch},
	}

	for _, tc := range testCases {
		if tc.expected == nil {
			require.NoError(t, utils.ValidateCommonName(tc.input), "could not validate %q", tc.input)
		} else {
			err := utils.ValidateCommonName(tc.input)
			require.EqualError(t, err, tc.expected.Error(), "%q was not invalid", tc.input)
		}
	}
}
