package bff_test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	records "github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func (s *bffTestSuite) TestLookup() {
	require := s.Require()
	params := &api.LookupParams{}

	// Test Bad Request (no parameters)
	_, err := s.client.Lookup(context.TODO(), params)
	require.EqualError(err, "[400] must provide either uuid or common_name in query params", "expected a 400 error with no params")

	// Provide some params
	params.CommonName = "api.alice.vaspbot.net"

	// Test NotFound returns a 404
	require.NoError(s.testnet.gds.UseError(mock.LookupRPC, codes.NotFound, "testnet not found"))
	require.NoError(s.mainnet.gds.UseError(mock.LookupRPC, codes.NotFound, "mainnet not found"))
	_, err = s.client.Lookup(context.TODO(), params)
	require.EqualError(err, "[404] no results returned for query", "expected a 404 error when both GDSes return not found")

	// Test InternalError when both GDSes return Unavailable
	require.NoError(s.testnet.gds.UseError(mock.LookupRPC, codes.Unavailable, "testnet cannot connect"))
	require.NoError(s.mainnet.gds.UseError(mock.LookupRPC, codes.Unavailable, "mainnet cannot connect"))
	_, err = s.client.Lookup(context.TODO(), params)
	require.EqualError(err, "[500] unable to execute Lookup request", "expected a 500 error when both GDSes return unavailable")

	// Test one result from TestNet
	require.NoError(s.testnet.gds.UseFixture(mock.LookupRPC, "testdata/testnet/lookup_reply.json"))
	require.NoError(s.mainnet.gds.UseError(mock.LookupRPC, codes.NotFound, "mainnet not found"))
	rep, err := s.client.Lookup(context.TODO(), params)
	require.NoError(err, "could not fetch expected result from testnet")
	require.Nil(rep.MainNet, "expected no mainnet result back from server")
	require.NotEmpty(rep.TestNet, "expected testnet result from server")
	require.Equal("6a57fea4-8fb7-42f3-bf0c-55fecccd2e53", rep.TestNet["id"])

	// Test one result from MainNet
	require.NoError(s.testnet.gds.UseError(mock.LookupRPC, codes.NotFound, "testnet not found"))
	require.NoError(s.mainnet.gds.UseFixture(mock.LookupRPC, "testdata/mainnet/lookup_reply.json"))
	rep, err = s.client.Lookup(context.TODO(), params)
	require.NoError(err, "could not fetch expected result from mainnet")
	require.Nil(rep.TestNet, "expected no testnet result back from server")
	require.NotEmpty(rep.MainNet, "expected mainnet result from server")
	require.Equal("ca0cff66-719f-4a62-8086-be953699b27d", rep.MainNet["id"])

	// Test results from both TestNet and MainNet
	require.NoError(s.testnet.gds.UseFixture(mock.LookupRPC, "testdata/testnet/lookup_reply.json"))
	require.NoError(s.mainnet.gds.UseFixture(mock.LookupRPC, "testdata/mainnet/lookup_reply.json"))
	rep, err = s.client.Lookup(context.TODO(), params)
	require.NoError(err, "could not fetch expected result from mainnet and testnet")
	require.NotEmpty(rep.MainNet, "expected mainnet result from server")
	require.NotEmpty(rep.TestNet, "expected testnet result from server")
}

func (s *bffTestSuite) TestLoadRegisterForm() {
	require := s.Require()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := s.client.LoadRegistrationForm(context.TODO())
	require.EqualError(err, "[401] this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint must have the read:vasp permission
	token, err := s.auth.NewTokenWithClaims(claims)
	require.NoError(err, "could not create token with incorrect permissions from claims")
	s.client.SetCredentials(api.Token(token))

	_, err = s.client.LoadRegistrationForm(context.TODO())
	require.EqualError(err, "[401] user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID and the server must not panic if it does not
	claims.Permissions = []string{"read:vasp"}
	token, err = s.auth.NewTokenWithClaims(claims)
	require.NoError(err, "could not create token without organizationID from claims")
	s.client.SetCredentials(api.Token(token))

	_, err = s.client.LoadRegistrationForm(context.TODO())
	require.EqualError(err, "[400] missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	token, err = s.auth.NewTokenWithClaims(claims)
	require.NoError(err, "could not create token with valid claims")
	s.client.SetCredentials(api.Token(token))

	_, err = s.client.LoadRegistrationForm(context.TODO())
	require.EqualError(err, "[404] no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Create organization in the database, but without registration form.
	// An empty registration form should be returned without panic.
	org, err := s.db.Organizations().Create(context.TODO())
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.Organizations().Delete(context.TODO(), org.Id)
	}()

	claims.OrgID = org.Id
	token, err = s.auth.NewTokenWithClaims(claims)
	require.NoError(err, "could not create token with valid claims")
	s.client.SetCredentials(api.Token(token))

	form, err := s.client.LoadRegistrationForm(context.TODO())
	require.NoError(err, "expected no error when no form data is stored")
	require.NotNil(form, "expected empty registration form when no form data is stored")

	// Load a registration form from fixtures and store it in the database
	org.Registration = &records.RegistrationForm{}
	err = loadFixture("testdata/registration_form.pb.json", org.Registration)
	require.NoError(err, "could not load registration form fixture")
	require.False(proto.Equal(form, org.Registration), "expected fixture to not be empty")

	err = s.db.Organizations().Update(context.TODO(), org)
	require.NoError(err, "could not update organization in database")

	form, err = s.client.LoadRegistrationForm(context.TODO())
	require.NoError(err, "expected no error when form data is available")
	require.NotNil(form, "expected completed registration form when form data is available")
	require.True(proto.Equal(form, org.Registration), "expected completed registration form when form data is available")
}

func (s *bffTestSuite) TestSubmitRegistration() {
	s.T().Skip() // Needs authtest setup
	var err error
	require := s.Require()

	// Test both the testnet and the mainnet registration
	for _, network := range []string{"testnet", "mainnet"} {
		// Identify the mock being used in this loop
		var mgds *mock.GDS
		switch network {
		case "testnet":
			mgds = s.testnet.gds
		case "mainnet":
			mgds = s.mainnet.gds
		}

		// Reset the calls on the mocks to ensure the correct mock GDS is being called
		expectedCalls := make(map[string]int)
		s.testnet.gds.Reset()
		s.mainnet.gds.Reset()

		// Test Invalid Argument Error
		mgds.UseError(mock.RegisterRPC, codes.InvalidArgument, "the TRISA endpoint is not valid")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		require.EqualError(err, "[400] the TRISA endpoint is not valid")
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test Already Exists error
		mgds.UseError(mock.RegisterRPC, codes.AlreadyExists, "this VASP is already registered")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		require.EqualError(err, "[400] this VASP is already registered")
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test Aborted error
		mgds.UseError(mock.RegisterRPC, codes.Aborted, "a conflict occurred")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		require.EqualError(err, "[409] a conflict occurred")
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test Timeout error
		mgds.UseError(mock.RegisterRPC, codes.DeadlineExceeded, "deadline exceeded")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		require.EqualError(err, fmt.Sprintf("[500] could not register with %s", network))
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test FailedPrecondition error
		mgds.UseError(mock.RegisterRPC, codes.FailedPrecondition, "couldn't access database")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		require.EqualError(err, fmt.Sprintf("[500] could not register with %s", network))
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test a valid register reply
		err = mgds.UseFixture(mock.RegisterRPC, fmt.Sprintf("testdata/%s/register_reply.json", network))
		require.NoError(err, "could not load register reply fixture")

		rep, err := s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		require.NoError(err, "could not make register call with valid payload")
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Check the register response is valid
		require.Empty(rep.Error, "an error message was returned from the server")
		require.NotEmpty(rep.Id, "the ID was not returned from the server")
		require.NotEmpty(rep.RegisteredDirectory, "the registered directory was not returned from the server")
		require.Equal(rep.CommonName, "trisa.example.com", "the common name was not returned from the server")
		require.Equal(rep.Status, "PENDING_REVIEW", "the verification status was not returned by the server")
		require.Equal(rep.Message, "thank you for registering", "a message was not returned from the server")
		require.Equal(rep.PKCS12Password, "supersecret", "a pkcs12 password was not returned from the server")

		// Test that a post to an incorrect network returns an error.
		_, err = s.client.SubmitRegistration(context.TODO(), "notanetwork")
		require.EqualError(err, "[404] network should be either testnet or mainnet")
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)
	}
}

func (s *bffTestSuite) TestVerifyEmail() {
	require := s.Require()
	params := &api.VerifyContactParams{}

	// Test Bad Request (no parameters)
	_, err := s.client.VerifyContact(context.TODO(), params)
	require.EqualError(err, "[400] must provide vaspID, token, and registered_directory in query parameters", "expected a 400 error with no params")

	// Test Bad Request (only vaspID specified)
	params.ID = uuid.NewString()
	_, err = s.client.VerifyContact(context.TODO(), params)
	require.EqualError(err, "[400] must provide vaspID, token, and registered_directory in query parameters", "expected a 400 error with no params")

	// Test Bad Request (only vaspID and token specified)
	params.Token = "abcdefghijklmnopqrstuvwxyz"
	_, err = s.client.VerifyContact(context.TODO(), params)
	require.EqualError(err, "[400] must provide vaspID, token, and registered_directory in query parameters", "expected a 400 error with no params")

	// Test Bad Request (only unknown registered_directory specified)
	params.Directory = "equitylo.rd"
	_, err = s.client.VerifyContact(context.TODO(), params)
	require.EqualError(err, "[400] unknown registered directory")

	// Assert that to this point no GDS method has been called
	require.Empty(s.testnet.gds.Calls[mock.VerifyContactRPC], "expected no testnet calls")
	require.Empty(s.mainnet.gds.Calls[mock.VerifyContactRPC], "expected no mainnet calls")

	// Test good requests to the registered directory
	for i, directory := range []string{"trisatest.net", "vaspdirectory.net"} {
		params.Directory = directory

		// Identify the mock being used in this loop
		var mgds *mock.GDS
		switch i {
		case 0:
			mgds = s.testnet.gds
		case 1:
			mgds = s.mainnet.gds
		}

		// Reset the calls on the mocks to ensure the correct mock GDS is being called
		expectedCalls := make(map[string]int)
		s.testnet.gds.Reset()
		s.mainnet.gds.Reset()

		// Test invalid argument error
		mgds.UseError(mock.VerifyContactRPC, codes.InvalidArgument, "incorrect vasp id")
		_, err = s.client.VerifyContact(context.TODO(), params)
		expectedCalls[directory]++
		require.EqualError(err, "[400] incorrect vasp id")
		require.Equal(expectedCalls["trisatest.net"], s.testnet.gds.Calls[mock.VerifyContactRPC], "check testnet calls during %s testing", directory)
		require.Equal(expectedCalls["vaspdirectory.net"], s.mainnet.gds.Calls[mock.VerifyContactRPC], "check mainnet calls during %s testing", directory)

		// Test not found error
		mgds.UseError(mock.VerifyContactRPC, codes.NotFound, "could not lookup contact with token")
		_, err = s.client.VerifyContact(context.TODO(), params)
		expectedCalls[directory]++
		require.EqualError(err, "[404] could not lookup contact with token")
		require.Equal(expectedCalls["trisatest.net"], s.testnet.gds.Calls[mock.VerifyContactRPC], "check testnet calls during %s testing", directory)
		require.Equal(expectedCalls["vaspdirectory.net"], s.mainnet.gds.Calls[mock.VerifyContactRPC], "check mainnet calls during %s testing", directory)

		// Test aborted error
		mgds.UseError(mock.VerifyContactRPC, codes.Aborted, "could not update verification status")
		_, err = s.client.VerifyContact(context.TODO(), params)
		expectedCalls[directory]++
		require.EqualError(err, "[409] could not update verification status")
		require.Equal(expectedCalls["trisatest.net"], s.testnet.gds.Calls[mock.VerifyContactRPC], "check testnet calls during %s testing", directory)
		require.Equal(expectedCalls["vaspdirectory.net"], s.mainnet.gds.Calls[mock.VerifyContactRPC], "check mainnet calls during %s testing", directory)

		// Test failed precondition error
		mgds.UseError(mock.VerifyContactRPC, codes.FailedPrecondition, "something went wrong")
		_, err = s.client.VerifyContact(context.TODO(), params)
		expectedCalls[directory]++
		require.EqualError(err, "[500] something went wrong")
		require.Equal(expectedCalls["trisatest.net"], s.testnet.gds.Calls[mock.VerifyContactRPC], "check testnet calls during %s testing", directory)
		require.Equal(expectedCalls["vaspdirectory.net"], s.mainnet.gds.Calls[mock.VerifyContactRPC], "check mainnet calls during %s testing", directory)

		// Test internal error
		mgds.UseError(mock.VerifyContactRPC, codes.FailedPrecondition, "boom hiss")
		_, err = s.client.VerifyContact(context.TODO(), params)
		expectedCalls[directory]++
		require.EqualError(err, "[500] boom hiss")
		require.Equal(expectedCalls["trisatest.net"], s.testnet.gds.Calls[mock.VerifyContactRPC], "check testnet calls during %s testing", directory)
		require.Equal(expectedCalls["vaspdirectory.net"], s.mainnet.gds.Calls[mock.VerifyContactRPC], "check mainnet calls during %s testing", directory)

		// Test a valid verify email response
		mgds.OnVerifyContact = func(ctx context.Context, in *gds.VerifyContactRequest) (out *gds.VerifyContactReply, err error) {
			return &gds.VerifyContactReply{
				Status:  models.VerificationState_PENDING_REVIEW,
				Message: "thank you for verifying your contact information",
			}, nil
		}

		rep, err := s.client.VerifyContact(context.TODO(), params)
		require.NoError(err, "unexpected error during good request")
		require.NotEmpty(rep, "empty response returned")
		require.Empty(rep.Error, "an error message was returned by the bff")
		require.Equal("PENDING_REVIEW", rep.Status)
		require.Equal("thank you for verifying your contact information", rep.Message)
	}
}
