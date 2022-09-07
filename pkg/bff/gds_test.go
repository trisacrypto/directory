package bff_test

import (
	"context"
	"fmt"
	"net/http"

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
	s.requireError(err, http.StatusBadRequest, "must provide either uuid or common_name in query params", "expected a 400 error with no params")

	// Provide some params
	params.CommonName = "api.alice.vaspbot.net"

	// Test NotFound returns a 404
	require.NoError(s.testnet.gds.UseError(mock.LookupRPC, codes.NotFound, "testnet not found"))
	require.NoError(s.mainnet.gds.UseError(mock.LookupRPC, codes.NotFound, "mainnet not found"))
	_, err = s.client.Lookup(context.TODO(), params)
	s.requireError(err, http.StatusNotFound, "no results returned for query", "expected a 404 error when both GDSes return not found")

	// Test InternalError when both GDSes return Unavailable
	require.NoError(s.testnet.gds.UseError(mock.LookupRPC, codes.Unavailable, "testnet cannot connect"))
	require.NoError(s.mainnet.gds.UseError(mock.LookupRPC, codes.Unavailable, "mainnet cannot connect"))
	_, err = s.client.Lookup(context.TODO(), params)
	s.requireError(err, http.StatusInternalServerError, "unable to execute Lookup request", "expected a 500 error when both GDSes return unavailable")

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
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint must have the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions from claims")
	_, err = s.client.LoadRegistrationForm(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID and the server must not panic if it does not
	claims.Permissions = []string{"read:vasp"}
	require.NoError(s.SetClientCredentials(claims), "could not create token without organizationID from claims")

	_, err = s.client.LoadRegistrationForm(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")

	_, err = s.client.LoadRegistrationForm(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Create organization in the database, but without registration form.
	// An empty registration form should be returned without panic.
	org, err := s.db.Organizations().Create(context.TODO())
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.Organizations().Delete(context.TODO(), org.Id)
	}()

	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")

	form, err := s.client.LoadRegistrationForm(context.TODO())
	require.NoError(err, "expected no error when no form data is stored")
	require.NotNil(form, "expected empty registration form when no form data is stored")
	require.NotNil(form.State, "expected form state to be populated")
	require.Equal(int32(1), form.State.Current, "expected initial form step to be 1")
	require.False(form.State.ReadyToSubmit, "expected form state to be not ready to submit")
	require.Len(form.State.Steps, 1, "expected 1 step in initial form state")
	require.Equal("progress", form.State.Steps[0].Status, "expected first form step to be in progress")
	require.Empty(form.State.Started, "expected form started timestamp to be empty")

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

func (s *bffTestSuite) TestSaveRegisterForm() {
	require := s.Require()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Load registration forms fixture
	form := &records.RegistrationForm{}
	err := loadFixture("testdata/registration_form.pb.json", form)
	require.NoError(err, "could not load registration form fixture")

	// Endpoint requires CSRF protection
	_, err = s.client.SaveRegistrationForm(context.TODO(), form)
	s.requireError(err, http.StatusForbidden, "csrf verification failed for request", "expected error when request is not CSRF protected")

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
	_, err = s.client.SaveRegistrationForm(context.TODO(), form)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.SaveRegistrationForm(context.TODO(), form)
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID and the server must panic if it does not
	claims.Permissions = []string{"update:vasp"}
	require.NoError(s.SetClientCredentials(claims), "could not create token without organizationID from claims")
	_, err = s.client.SaveRegistrationForm(context.TODO(), form)
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	_, err = s.client.SaveRegistrationForm(context.TODO(), form)
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Create an organization in the database that does not contain a registration form
	org, err := s.db.Organizations().Create(context.TODO())
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.Organizations().Delete(context.TODO(), org.Id)
	}()

	// Create valid credentials for the remaining tests
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")

	// Should be able to save an empty registration form
	reply, err := s.client.SaveRegistrationForm(context.TODO(), &records.RegistrationForm{})
	require.NoError(err, "should not receive an error when saving an empty registration form")
	require.Nil(reply, "should receive 204 No Content when saving an empty registration form")

	// Empty registration form should be saved in the database
	org, err = s.db.Organizations().Retrieve(context.TODO(), org.Id)
	require.NoError(err, "could not retrieve organization from database")
	require.True(proto.Equal(org.Registration, &records.RegistrationForm{}), "expected empty registration form")

	// Should be able to save the fixture form
	reply, err = s.client.SaveRegistrationForm(context.TODO(), form)
	require.NoError(err, "should not receive an error when saving a registration form")
	require.NotNil(reply, "uploaded form should be returned when a non-empty registration form is saved")
	require.NotEmpty(reply.State.Started, "expected form started timestamp to be set")
	reply.State.Started = ""
	require.True(proto.Equal(form, reply), "expected returned registration form to match uploaded form")

	org, err = s.db.Organizations().Retrieve(context.TODO(), org.Id)
	require.NoError(err, "could not retrieve updated org from database")
	require.NotEmpty(org.Registration.State.Started, "expected registration form started timestamp to be populated")
	org.Registration.State.Started = ""
	require.True(proto.Equal(org.Registration, form), "expected form saved in database to match form uploaded")

	// Should be able to "clear" a registration by saving an empty registration form
	reply, err = s.client.SaveRegistrationForm(context.TODO(), &records.RegistrationForm{})
	require.NoError(err, "should not receive an error when saving an empty registration form")
	require.Nil(reply, "should receive 204 No Content when saving an empty registration form")

	org, err = s.db.Organizations().Retrieve(context.TODO(), org.Id)
	require.NoError(err, "could not retrieve updated org from database")
	require.False(proto.Equal(org.Registration, form), "expected form saved in database to be cleared")
}

func (s *bffTestSuite) TestSubmitRegistration() {
	var err error
	require := s.Require()

	// Test setup: create an organization with a valid registration form that has not
	// been submitted yet - at the end of the test both mainnet and testnet should be
	// submitted and the response from the directory updated on the organization.
	org, err := s.db.Organizations().Create(context.TODO())
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.Organizations().Delete(context.TODO(), org.Id)
	}()

	// Save the registration form fixture on the organization
	org.Registration = &records.RegistrationForm{}
	require.NoError(loadFixture("testdata/registration_form.pb.json", org.Registration), "could not load registration form from the fixtures")
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization with registration form")

	// Test both the testnet and the mainnet registration
	for _, network := range []string{"testnet", "mainnet"} {
		// Create new claims and unset CSRF protection to ensure that the endpoint
		// permissions tests are checked for both testnet and mainnet.
		s.client.(*api.APIv1).SetCredentials(nil)
		s.client.(*api.APIv1).SetCSRFProtect(false)
		claims := &authtest.Claims{
			Email:       "leopold.wentzel@gmail.com",
			Permissions: []string{"read:nothing"},
		}

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

		// Endpoint should require CSRF protection
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		s.requireError(err, http.StatusForbidden, "csrf verification failed for request", "expected error when request is not CSRF protected")

		// Endpoint must be authenticated
		require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

		// Endpoint requires the update:vasp permission
		require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

		// Claims must have an organization ID and the server must panic if it does not
		claims.Permissions = []string{"update:vasp"}
		require.NoError(s.SetClientCredentials(claims), "could not create token without organizationID from claims")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

		// Create valid claims but no record in the database - should not panic and should return an error
		claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
		require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

		// From this point on submit valid claims and test responses from GDS
		// NOTE: for registration form validation see TestSubmitRegistrationNotReady
		claims.OrgID = org.Id
		require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")

		// Test GDS returns Invalid Argument Error
		mgds.UseError(mock.RegisterRPC, codes.InvalidArgument, "the TRISA endpoint is not valid")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		s.requireError(err, http.StatusBadRequest, "the TRISA endpoint is not valid")
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test GDS returns Already Exists error
		mgds.UseError(mock.RegisterRPC, codes.AlreadyExists, "this VASP is already registered")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		s.requireError(err, http.StatusBadRequest, "this VASP is already registered")
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test GDS returns Aborted error
		mgds.UseError(mock.RegisterRPC, codes.Aborted, "a conflict occurred")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		s.requireError(err, http.StatusConflict, "a conflict occurred")
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test GDS returns Timeout error
		mgds.UseError(mock.RegisterRPC, codes.DeadlineExceeded, "deadline exceeded")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		s.requireError(err, http.StatusInternalServerError, fmt.Sprintf("could not register with %s", network))
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test GDS returnsFailedPrecondition error
		mgds.UseError(mock.RegisterRPC, codes.FailedPrecondition, "couldn't access database")
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		expectedCalls[network]++
		s.requireError(err, http.StatusInternalServerError, fmt.Sprintf("could not register with %s", network))
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)

		// Test a valid register reply
		// TODO: we should validate what is being sent to the GDS server
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
		require.NotEmpty(rep.CommonName, "the common name was not returned from the server")
		require.Equal(rep.Status, "PENDING_REVIEW", "the verification status was not returned by the server")
		require.Equal(rep.Message, "thank you for registering", "a message was not returned from the server")
		require.Equal(rep.PKCS12Password, "supersecret", "a pkcs12 password was not returned from the server")

		// Test that a post to an incorrect network returns an error.
		_, err = s.client.SubmitRegistration(context.TODO(), "notanetwork")
		s.requireError(err, http.StatusNotFound, "network should be either testnet or mainnet")
		require.Equal(expectedCalls["testnet"], s.testnet.gds.Calls[mock.RegisterRPC], "check testnet calls during %s testing", network)
		require.Equal(expectedCalls["mainnet"], s.mainnet.gds.Calls[mock.RegisterRPC], "check mainnet calls during %s testing", network)
	}

	// Ensure that the directory record is stored on the database after registration
	org, err = s.db.Organizations().Retrieve(context.TODO(), org.Id)
	require.NoError(err, "could not update organization from the database")

	require.NotNil(org.Testnet, "missing testnet directory record after registration")
	require.Equal(org.Testnet.Id, "6041571e-09b4-47e7-870a-723f8032cd6c", "incorrect testnet directory id")
	require.Equal(org.Testnet.RegisteredDirectory, "trisatest.net", "incorrect testnet registerd directory ")
	require.Equal(org.Testnet.CommonName, "test.trisa.example.ua", "incorrect testnet directory common name")
	require.NotEmpty(org.Testnet.Submitted, "expected testnet submitted timestamp stored in database")

	require.NotNil(org.Mainnet, "missing mainnet directory record after registration")
	require.Equal(org.Mainnet.Id, "5bafb054-5868-439e-9b3c-75db91810714", "incorrect mainnet directory id")
	require.Equal(org.Mainnet.RegisteredDirectory, "vaspdirectory.net", "incorrect mainnet registerd directory ")
	require.Equal(org.Mainnet.CommonName, "trisa.example.ua", "incorrect mainnet directory common name")
	require.NotEmpty(org.Mainnet.Submitted, "expected mainnet submitted timestamp stored in database")
}

func (s *bffTestSuite) TestSubmitRegistrationNotReady() {
	require := s.Require()

	// Ensure that a bad argument error is returned if the registration form is not
	// ready to submit. Create an organization that has a registration form without
	// network details and valid claims to access the record.
	org, err := s.db.Organizations().Create(context.TODO())
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.Organizations().Delete(context.TODO(), org.Id)
	}()

	// Ensure the registration is not ready to submit by removing mainnet and testnet
	org.Registration = &records.RegistrationForm{}
	require.NoError(loadFixture("testdata/registration_form.pb.json", org.Registration), "could not load registration form from the fixtures")
	org.Registration.Mainnet = nil
	org.Registration.Testnet = nil
	require.False(org.Registration.ReadyToSubmit("both"), "registration should not be ready to submit")

	// Save the registration form fixture on the organization
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization with registration form")

	// Create authenticated user context
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"update:vasp"},
		OrgID:       org.Id,
	}
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	require.NoError(s.SetClientCSRFProtection(), "could not set CSRF protection on client")

	// Expect 400 error for both mainnet and testnet
	for _, network := range []string{"testnet", "mainnet"} {
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		s.requireError(err, http.StatusBadRequest, "registration form is not ready to submit", "expected error when registration form is not ready to submit")
	}

	// While we're here, also test that we receive a 404 for a bad network
	_, err = s.client.SubmitRegistration(context.TODO(), "notanetwork")
	s.requireError(err, http.StatusNotFound, "network should be either testnet or mainnet", "expected error when submitting registration to incorrect network name")
}

func (s *bffTestSuite) TestCannotResubmitRegistration() {
	require := s.Require()

	// Ensure that a conflict error is returned if the registration form has already
	// been ready to submitted. Create an organization that has directory records for
	// both networks and valid claims to access the record.
	org, err := s.db.Organizations().Create(context.TODO())
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.Organizations().Delete(context.TODO(), org.Id)
	}()

	// Ensure the registration is not ready to submit by removing mainnet and testnet
	org.Testnet = &records.DirectoryRecord{
		Id:                  uuid.NewString(),
		CommonName:          "test.trisa.example.com",
		RegisteredDirectory: "trisatest.net",
		Submitted:           "2022-02-21T15:32:31Z",
	}
	org.Mainnet = &records.DirectoryRecord{
		Id:                  uuid.NewString(),
		CommonName:          "trisa.example.com",
		RegisteredDirectory: "vaspdirectory.net",
		Submitted:           "2022-02-23T09:51:15Z",
	}

	// Save the registration form fixture on the organization
	require.NoError(s.db.Organizations().Update(context.TODO(), org), "could not update organization with registration form")

	// Create authenticated user context
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"update:vasp"},
		OrgID:       org.Id,
	}
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	require.NoError(s.SetClientCSRFProtection(), "could not set CSRF protection on client")

	// Expect 400 error for both mainnet and testnet
	for _, network := range []string{"testnet", "mainnet"} {
		_, err = s.client.SubmitRegistration(context.TODO(), network)
		s.requireError(err, http.StatusConflict, fmt.Sprintf("registration form has already been submitted to the %s", network), "expected error when registration form has already been submitted")
	}
}

func (s *bffTestSuite) TestVerifyEmail() {
	require := s.Require()
	params := &api.VerifyContactParams{}

	// Test Bad Request (no parameters)
	_, err := s.client.VerifyContact(context.TODO(), params)
	s.requireError(err, http.StatusBadRequest, "must provide vaspID, token, and registered_directory in query parameters", "expected a 400 error with no params")

	// Test Bad Request (only vaspID specified)
	params.ID = uuid.NewString()
	_, err = s.client.VerifyContact(context.TODO(), params)
	s.requireError(err, http.StatusBadRequest, "must provide vaspID, token, and registered_directory in query parameters", "expected a 400 error with no params")

	// Test Bad Request (only vaspID and token specified)
	params.Token = "abcdefghijklmnopqrstuvwxyz"
	_, err = s.client.VerifyContact(context.TODO(), params)
	s.requireError(err, http.StatusBadRequest, "must provide vaspID, token, and registered_directory in query parameters", "expected a 400 error with no params")

	// Test Bad Request (only unknown registered_directory specified)
	params.Directory = "equitylo.rd"
	_, err = s.client.VerifyContact(context.TODO(), params)
	s.requireError(err, http.StatusBadRequest, "unknown registered directory")

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
		s.requireError(err, http.StatusBadRequest, "incorrect vasp id")
		require.Equal(expectedCalls["trisatest.net"], s.testnet.gds.Calls[mock.VerifyContactRPC], "check testnet calls during %s testing", directory)
		require.Equal(expectedCalls["vaspdirectory.net"], s.mainnet.gds.Calls[mock.VerifyContactRPC], "check mainnet calls during %s testing", directory)

		// Test not found error
		mgds.UseError(mock.VerifyContactRPC, codes.NotFound, "could not lookup contact with token")
		_, err = s.client.VerifyContact(context.TODO(), params)
		expectedCalls[directory]++
		s.requireError(err, http.StatusNotFound, "could not lookup contact with token")
		require.Equal(expectedCalls["trisatest.net"], s.testnet.gds.Calls[mock.VerifyContactRPC], "check testnet calls during %s testing", directory)
		require.Equal(expectedCalls["vaspdirectory.net"], s.mainnet.gds.Calls[mock.VerifyContactRPC], "check mainnet calls during %s testing", directory)

		// Test aborted error
		mgds.UseError(mock.VerifyContactRPC, codes.Aborted, "could not update verification status")
		_, err = s.client.VerifyContact(context.TODO(), params)
		expectedCalls[directory]++
		s.requireError(err, http.StatusConflict, "could not update verification status")
		require.Equal(expectedCalls["trisatest.net"], s.testnet.gds.Calls[mock.VerifyContactRPC], "check testnet calls during %s testing", directory)
		require.Equal(expectedCalls["vaspdirectory.net"], s.mainnet.gds.Calls[mock.VerifyContactRPC], "check mainnet calls during %s testing", directory)

		// Test failed precondition error
		mgds.UseError(mock.VerifyContactRPC, codes.FailedPrecondition, "something went wrong")
		_, err = s.client.VerifyContact(context.TODO(), params)
		expectedCalls[directory]++
		s.requireError(err, http.StatusInternalServerError, "something went wrong")
		require.Equal(expectedCalls["trisatest.net"], s.testnet.gds.Calls[mock.VerifyContactRPC], "check testnet calls during %s testing", directory)
		require.Equal(expectedCalls["vaspdirectory.net"], s.mainnet.gds.Calls[mock.VerifyContactRPC], "check mainnet calls during %s testing", directory)

		// Test internal error
		mgds.UseError(mock.VerifyContactRPC, codes.FailedPrecondition, "boom hiss")
		_, err = s.client.VerifyContact(context.TODO(), params)
		expectedCalls[directory]++
		s.requireError(err, http.StatusInternalServerError, "boom hiss")
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
