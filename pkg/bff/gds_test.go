package bff_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/mock"
	records "github.com/trisacrypto/directory/pkg/bff/models/v1"
	models "github.com/trisacrypto/directory/pkg/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
	"github.com/trisacrypto/directory/pkg/utils/wire"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *bffTestSuite) TestVASPNames() {
	require := s.Require()
	ctx := context.Background()

	// Test that no names are returned when there are no VASPs
	rep, err := s.client.LookupAutocomplete(ctx)
	require.NoError(err, "error calling names endpoint")
	require.Empty(rep, "expected no VASP names to be returned")

	// Test that a single name is returned when there is one VASP in testnet
	testnetVASP := &pb.VASP{}
	testnetFixture := filepath.Join("testdata", "testnet", "vasp.json")
	require.NoError(loadFixture(testnetFixture, testnetVASP))
	testnetVASP.VerificationStatus = pb.VerificationState_VERIFIED
	_, err = s.TestNetDB().CreateVASP(ctx, testnetVASP)
	require.NoError(err, "error creating VASP fixture in database")
	rep, err = s.client.LookupAutocomplete(ctx)
	require.NoError(err, "error calling names endpoint")
	require.Equal(map[string]string{"Alice VASP, Inc.": "alice", "alice": "alice"}, rep, "wrong names returned from names endpoint")

	// Test that unverified VASPs are not returned
	mainnetVASP := &pb.VASP{}
	mainnetFixture := filepath.Join("testdata", "mainnet", "vasp.json")
	require.NoError(loadFixture(mainnetFixture, mainnetVASP))
	_, err = s.MainNetDB().CreateVASP(ctx, mainnetVASP)
	require.NoError(err, "error creating VASP fixture in database")
	rep, err = s.client.LookupAutocomplete(ctx)
	require.NoError(err, "error calling names endpoint")
	require.Equal(map[string]string{"Alice VASP, Inc.": "alice", "alice": "alice"}, rep, "wrong names returned from names endpoint")

	// Test that duplicate names are not returned
	mainnetVASP.VerificationStatus = pb.VerificationState_VERIFIED
	require.NoError(s.MainNetDB().UpdateVASP(ctx, mainnetVASP), "error updating VASP fixture in database")
	rep, err = s.client.LookupAutocomplete(ctx)
	require.NoError(err, "error calling names endpoint")
	require.Equal(map[string]string{"Alice VASP, Inc.": "alice", "alice": "alice"}, rep, "wrong names returned from names endpoint")

	// Test names returned from both testnet and mainnet
	testnetVASP.Id = uuid.New().String()
	testnetVASP.CommonName = "testnet.bob.vaspbot.net"
	testnetVASP.Entity.Name.NameIdentifiers[0].LegalPersonName = "Bob VASP, Inc."
	_, err = s.TestNetDB().CreateVASP(ctx, testnetVASP)
	require.NoError(err, "error creating VASP fixture in database")
	mainnetVASP.Id = uuid.New().String()
	mainnetVASP.CommonName = "mainnet.charlie.vaspbot.net"
	mainnetVASP.Entity.Name.NameIdentifiers[0].LegalPersonName = "Charlie VASP, Inc."
	_, err = s.MainNetDB().CreateVASP(ctx, mainnetVASP)
	require.NoError(err, "error creating VASP fixture in database")

	expected := map[string]string{
		"Alice VASP, Inc.": "alice", "alice": "alice",
		"Bob VASP, Inc.": "testnet.bob.vaspbot.net", "testnet.bob.vaspbot.net": "testnet.bob.vaspbot.net",
		"Charlie VASP, Inc.": "mainnet.charlie.vaspbot.net", "mainnet.charlie.vaspbot.net": "mainnet.charlie.vaspbot.net",
	}

	rep, err = s.client.LookupAutocomplete(ctx)
	require.NoError(err, "error calling names endpoint")
	require.Equal(expected, rep, "wrong names returned from names endpoint")
}

func (s *bffTestSuite) TestLoadRegisterForm() {
	require := s.Require()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := s.client.LoadRegistrationForm(context.TODO(), nil)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint must have the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions from claims")
	_, err = s.client.LoadRegistrationForm(context.TODO(), nil)
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID and the server must not panic if it does not
	claims.Permissions = []string{auth.ReadVASP}
	require.NoError(s.SetClientCredentials(claims), "could not create token without organizationID from claims")

	_, err = s.client.LoadRegistrationForm(context.TODO(), nil)
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")

	_, err = s.client.LoadRegistrationForm(context.TODO(), nil)
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Create organization in the database, but without registration form.
	// An empty registration form should be returned without panic.
	org := &records.Organization{}
	_, err = s.DB().CreateOrganization(context.Background(), org)
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.DB().DeleteOrganization(context.Background(), org.UUID())
	}()

	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")

	// Should return an error if the step is not a valid step
	_, err = s.client.LoadRegistrationForm(context.TODO(), &api.RegistrationFormParams{Step: "invalid"})
	s.requireError(err, http.StatusBadRequest, "unknown registration form step \"invalid\"", "expected error when step is invalid")

	out, err := s.client.LoadRegistrationForm(context.TODO(), nil)
	require.NoError(err, "expected no error when no form data is stored")
	require.NotNil(out.Form, "expected empty registration form when no form data is stored")

	form := out.Form
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

	err = s.DB().UpdateOrganization(context.Background(), org)
	require.NoError(err, "could not update organization in database")

	out, err = s.client.LoadRegistrationForm(context.TODO(), nil)
	require.NoError(err, "expected no error when form data is available")
	require.NotNil(out.Form, "expected completed registration form when form data is available")
	require.True(proto.Equal(out.Form, org.Registration), "expected completed registration form when form data is available")

	// Test loading individual steps of the registration form
	// Basic Details
	full := out.Form
	params := &api.RegistrationFormParams{
		Step: api.StepBasicDetails,
	}
	out, err = s.client.LoadRegistrationForm(context.Background(), params)
	require.NoError(err, "expected no error when form data is available")
	require.Equal(params.Step, out.Step, "expected returned step to be the same as the requested step")
	require.Nil(out.Errors, "expected no validation errors for basic details step")
	require.NotNil(out.Form, "expected returned form to not be nil")
	require.Equal(full.State, out.Form.State, "expected returned form to have the same state as the full form")
	require.Equal(full.Website, out.Form.Website, "expected returned form to have the same website as the full form")
	require.Equal(full.BusinessCategory, out.Form.BusinessCategory, "expected returned form to have the same business category as the full form")
	require.Equal(full.VaspCategories, out.Form.VaspCategories, "expected returned form to have the same vasp categories as the full form")
	require.Equal(full.EstablishedOn, out.Form.EstablishedOn, "expected returned form to have the same established on date as the full form")
	require.Equal(full.OrganizationName, out.Form.OrganizationName, "expected returned form to have the same organization name as the full form")
	require.Nil(out.Form.Entity, "expected returned form to not have legal person data")
	require.Nil(out.Form.Contacts, "expected returned form to not have contacts data")
	require.Nil(out.Form.Trixo, "expected returned form to not have trixo data")
	require.Nil(out.Form.Testnet, "expected returned form to not have testnet data")
	require.Nil(out.Form.Mainnet, "expected returned form to not have mainnet data")

	// Legal Person
	params.Step = api.StepLegalPerson
	out, err = s.client.LoadRegistrationForm(context.Background(), params)
	require.NoError(err, "expected no error when form data is available")
	require.Equal(params.Step, out.Step, "expected returned step to be the same as the requested step")

	fmt.Printf("%s\n", out.Errors)
	require.Nil(out.Errors, "expected no validation errors for legal person step")
	require.NotNil(out.Form, "expected returned form to not be nil")
	require.Equal(full.State, out.Form.State, "expected returned form to have the same state as the full form")
	require.Equal(full.Entity, out.Form.Entity, "expected returned form to have the same legal person data as the full form")
	require.Empty(out.Form.Website, "expected returned form to not have basic details data")
	require.Nil(out.Form.Contacts, "expected returned form to not have contacts data")
	require.Nil(out.Form.Trixo, "expected returned form to not have trixo data")
	require.Nil(out.Form.Testnet, "expected returned form to not have testnet data")
	require.Nil(out.Form.Mainnet, "expected returned form to not have mainnet data")

	// Contacts
	params.Step = api.StepContacts
	out, err = s.client.LoadRegistrationForm(context.Background(), params)
	require.NoError(err, "expected no error when form data is available")
	require.Equal(params.Step, out.Step, "expected returned step to be the same as the requested step")
	require.Nil(out.Errors, "expected no validation errors for contacts step")
	require.NotNil(out.Form, "expected returned form to not be nil")
	require.Equal(full.State, out.Form.State, "expected returned form to have the same state as the full form")
	require.Equal(full.Contacts, out.Form.Contacts, "expected returned form to have the same contacts data as the full form")
	require.Empty(out.Form.Website, "expected returned form to not have basic details data")
	require.Nil(out.Form.Entity, "expected returned form to not have legal person data")
	require.Nil(out.Form.Trixo, "expected returned form to not have trixo data")
	require.Nil(out.Form.Testnet, "expected returned form to not have testnet data")
	require.Nil(out.Form.Mainnet, "expected returned form to not have mainnet data")

	// TRIXO
	params.Step = api.StepTRIXO
	out, err = s.client.LoadRegistrationForm(context.Background(), params)
	require.NoError(err, "expected no error when form data is available")
	require.Equal(params.Step, out.Step, "expected returned step to be the same as the requested step")
	require.Nil(out.Errors, "expected no validation errors for trixo step")
	require.NotNil(out.Form, "expected returned form to not be nil")
	require.Equal(full.State, out.Form.State, "expected returned form to have the same state as the full form")
	require.Equal(full.Trixo, out.Form.Trixo, "expected returned form to have the same trixo data as the full form")
	require.Empty(out.Form.Website, "expected returned form to not have basic details data")
	require.Nil(out.Form.Entity, "expected returned form to not have legal person data")
	require.Nil(out.Form.Contacts, "expected returned form to not have contacts data")
	require.Nil(out.Form.Testnet, "expected returned form to not have testnet data")
	require.Nil(out.Form.Mainnet, "expected returned form to not have mainnet data")

	// TRISA
	params.Step = api.StepTRISA
	out, err = s.client.LoadRegistrationForm(context.Background(), params)
	require.NoError(err, "expected no error when form data is available")
	require.Equal(params.Step, out.Step, "expected returned step to be the same as the requested step")
	require.Nil(out.Errors, "expected no validation errors for trisa step")
	require.NotNil(out.Form, "expected returned form to not be nil")
	require.Equal(full.State, out.Form.State, "expected returned form to have the same state as the full form")
	require.Equal(full.Testnet, out.Form.Testnet, "expected returned form to have the same testnet data as the full form")
	require.Equal(full.Mainnet, out.Form.Mainnet, "expected returned form to have the same mainnet data as the full form")
	require.Empty(out.Form.Website, "expected returned form to not have basic details data")
	require.Nil(out.Form.Entity, "expected returned form to not have legal person data")
	require.Nil(out.Form.Contacts, "expected returned form to not have contacts data")
	require.Nil(out.Form.Trixo, "expected returned form to not have trixo data")

	// Load a form with validation errors into the database
	org.Registration = &records.RegistrationForm{}
	err = loadFixture("testdata/bad_registration_form.pb.json", org.Registration)
	require.NoError(err, "could not load registration form fixture")
	require.False(proto.Equal(form, org.Registration), "expected fixture to not be empty")

	err = s.DB().UpdateOrganization(context.Background(), org)
	require.NoError(err, "could not update organization in database")

	// All the validation errors in the fixture
	verrs := map[api.RegistrationFormStep][]*api.FieldValidationError{
		api.StepBasicDetails: {{Field: records.FieldWebsite, Error: records.ErrMissingField.Error()}},
		api.StepLegalPerson: {
			{Field: "", Error: records.ErrMissingField.Error()},
			{Field: "", Error: "ivms101: invalid field nationalIdentification.nationalIdentifier: invalid LEIX: invalid checksum"},
		},
		api.StepContacts: {{Field: records.FieldContactsTechnicalEmail, Error: records.ErrMissingField.Error()}},
		api.StepTRIXO:    {{Field: records.FieldTRIXOPrimaryNationalJurisdiction, Error: records.ErrMissingField.Error()}},
		api.StepTRISA:    {{Field: records.FieldTestNetCommonName, Error: records.ErrMissingField.Error()}},
	}

	// Test that errors are only returned for the requested step
	for step, expected := range verrs {
		params.Step = step
		out, err = s.client.LoadRegistrationForm(context.Background(), params)
		require.NoError(err, "expected no error when form data is available")
		require.Equal(params.Step, out.Step, "expected returned step to be the same as the requested step")
		require.Equal(expected, out.Errors, "expected returned validation errors to match the fixture errors for step %s", step)
		require.NotNil(out.Form, "expected returned form to not be nil")
	}

	// Test that all errors are returned when the step is not specified
	allErrs := []*api.FieldValidationError{}
	for _, verr := range verrs {
		allErrs = append(allErrs, verr...)
	}
	out, err = s.client.LoadRegistrationForm(context.Background(), nil)
	require.NoError(err, "expected no error when form data is available")
	require.ElementsMatch(allErrs, out.Errors, "expected all the validation errors in the form to be returned")
	require.NotNil(out.Form, "expected returned form to not be nil")
}

func (s *bffTestSuite) TestSaveRegisterForm() {
	require := s.Require()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Load registration forms fixture
	form := &api.RegistrationForm{
		Form: &records.RegistrationForm{},
	}
	err := loadFixture("testdata/registration_form.pb.json", form.Form)
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
	claims.Permissions = []string{auth.UpdateVASP}
	require.NoError(s.SetClientCredentials(claims), "could not create token without organizationID from claims")
	_, err = s.client.SaveRegistrationForm(context.TODO(), form)
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	_, err = s.client.SaveRegistrationForm(context.TODO(), form)
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Create an organization in the database that does not contain a registration form
	org := &records.Organization{}
	_, err = s.DB().CreateOrganization(context.Background(), org)
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.DB().DeleteOrganization(context.Background(), org.UUID())
	}()

	// Create valid credentials for the remaining tests
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")

	// Should return an error if the form step is not a valid step
	_, err = s.client.SaveRegistrationForm(context.TODO(), &api.RegistrationForm{Step: "invalid"})
	s.requireError(err, http.StatusBadRequest, "unknown registration form step \"invalid\"", "expected error when step is invalid")

	// Providing no form should return an error
	_, err = s.client.SaveRegistrationForm(context.TODO(), &api.RegistrationForm{Form: nil})
	s.requireError(err, http.StatusBadRequest, "no form was provided", "expected error when form is not provided")

	// Should be able to save the fixture form
	reply, err := s.client.SaveRegistrationForm(context.TODO(), form)
	require.NoError(err, "should not receive an error when saving a registration form")
	require.NotNil(reply, "uploaded form should be returned when a non-empty registration form is saved")
	require.NotEmpty(reply.Form.State.Started, "expected form started timestamp to be set")
	reply.Form.State.Started = ""
	require.True(proto.Equal(form.Form, reply.Form), "expected returned registration form to match uploaded form")

	org, err = s.DB().RetrieveOrganization(context.Background(), org.UUID())
	require.NoError(err, "could not retrieve updated org from database")
	require.NotEmpty(org.Registration.State.Started, "expected registration form started timestamp to be populated")
	org.Registration.State.Started = ""
	require.True(proto.Equal(org.Registration, form.Form), "expected form saved in database to match form uploaded")

	// Reset the form in the database
	org.Registration = records.NewRegisterForm()
	err = s.DB().UpdateOrganization(context.Background(), org)
	require.NoError(err, "could not update org in database")

	// Test saving steps of a registration form one by one
	// Basic details
	defaultForm := records.NewRegisterForm()
	partial := &api.RegistrationForm{}
	partial.Step = api.StepBasicDetails
	partial.Form = records.NewRegisterForm()
	partial.Form.Website = form.Form.Website
	partial.Form.BusinessCategory = form.Form.BusinessCategory
	partial.Form.VaspCategories = form.Form.VaspCategories
	partial.Form.EstablishedOn = form.Form.EstablishedOn
	partial.Form.OrganizationName = form.Form.OrganizationName
	reply, err = s.client.SaveRegistrationForm(context.TODO(), partial)
	require.NoError(err, "should not receive an error when saving a partial registration form")
	require.NotNil(reply, "uploaded form should be returned when a non-empty registration form is saved")
	require.Nil(reply.Errors, "expected no errors when saving a partially valid registration form")

	org, err = s.DB().RetrieveOrganization(context.Background(), org.UUID())
	require.NoError(err, "could not retrieve updated org from database")
	require.NotEmpty(org.Registration.State.Started, "expected registration form started timestamp to be populated")
	require.Equal(partial.Form.Website, org.Registration.Website, "expected form saved in database to match partial form uploaded")
	require.Equal(partial.Form.BusinessCategory, org.Registration.BusinessCategory, "expected form field in database to match partial form uploaded")
	require.Equal(partial.Form.VaspCategories, org.Registration.VaspCategories, "expected form field in database to match partial form uploaded")
	require.Equal(partial.Form.EstablishedOn, org.Registration.EstablishedOn, "expected form field in database to match partial form uploaded")
	require.Equal(partial.Form.OrganizationName, org.Registration.OrganizationName, "expected form field in database to match partial form uploaded")
	require.Equal(defaultForm.Entity, org.Registration.Entity, "expected form field in database to match default form")
	require.Equal(defaultForm.Contacts, org.Registration.Contacts, "expected form field in database to match default form")
	require.Equal(defaultForm.Trixo, org.Registration.Trixo, "expected form field in database to match default form")
	require.Equal(defaultForm.Testnet, org.Registration.Testnet, "expected form field in database to match default form")
	require.Equal(defaultForm.Mainnet, org.Registration.Mainnet, "expected form field in database to match default form")

	// Legal Person
	partial.Step = api.StepLegalPerson
	partial.Form = &records.RegistrationForm{
		Entity: form.Form.Entity,
	}
	reply, err = s.client.SaveRegistrationForm(context.TODO(), partial)
	require.NoError(err, "should not receive an error when saving a partial registration form")
	require.NotNil(reply, "uploaded form should be returned when a non-empty registration form is saved")
	require.Nil(reply.Errors, "expected no errors when saving a partially valid registration form")

	org, err = s.DB().RetrieveOrganization(context.Background(), org.UUID())
	require.NoError(err, "could not retrieve updated org from database")
	require.True(proto.Equal(partial.Form.Entity, org.Registration.Entity), "expected form saved in database to match partial form uploaded")
	require.True(proto.Equal(defaultForm.Contacts, org.Registration.Contacts), "expected form field in database to match default form")
	require.True(proto.Equal(defaultForm.Trixo, org.Registration.Trixo), "expected form field in database to match default form")
	require.True(proto.Equal(defaultForm.Testnet, org.Registration.Testnet), "expected form field in database to match default form")
	require.True(proto.Equal(defaultForm.Mainnet, org.Registration.Mainnet), "expected form field in database to match default form")

	// Contacts
	partial.Step = api.StepContacts
	partial.Form = &records.RegistrationForm{
		Contacts: form.Form.Contacts,
	}
	reply, err = s.client.SaveRegistrationForm(context.TODO(), partial)
	require.NoError(err, "should not receive an error when saving a partial registration form")
	require.NotNil(reply, "uploaded form should be returned when a non-empty registration form is saved")
	require.Nil(reply.Errors, "expected no errors when saving a partially valid registration form")

	org, err = s.DB().RetrieveOrganization(context.Background(), org.UUID())
	require.NoError(err, "could not retrieve updated org from database")
	require.True(proto.Equal(partial.Form.Contacts, org.Registration.Contacts), "expected form saved in database to match partial form uploaded")
	require.True(proto.Equal(defaultForm.Trixo, org.Registration.Trixo), "expected form field in database to match default form")
	require.True(proto.Equal(defaultForm.Testnet, org.Registration.Testnet), "expected form field in database to match default form")
	require.True(proto.Equal(defaultForm.Mainnet, org.Registration.Mainnet), "expected form field in database to match default form")

	// TRIXO
	partial.Step = api.StepTRIXO
	partial.Form = &records.RegistrationForm{
		Trixo: form.Form.Trixo,
	}
	reply, err = s.client.SaveRegistrationForm(context.TODO(), partial)
	require.NoError(err, "should not receive an error when saving a partial registration form")
	require.NotNil(reply, "uploaded form should be returned when a non-empty registration form is saved")
	require.Nil(reply.Errors, "expected no errors when saving a partially valid registration form")

	org, err = s.DB().RetrieveOrganization(context.Background(), org.UUID())
	require.NoError(err, "could not retrieve updated org from database")
	require.True(proto.Equal(partial.Form.Trixo, org.Registration.Trixo), "expected form saved in database to match partial form uploaded")
	require.True(proto.Equal(defaultForm.Testnet, org.Registration.Testnet), "expected form field in database to match default form")
	require.True(proto.Equal(defaultForm.Mainnet, org.Registration.Mainnet), "expected form field in database to match default form")

	// TRISA
	partial.Step = api.StepTRISA
	partial.Form = &records.RegistrationForm{
		Testnet: form.Form.Testnet,
		Mainnet: form.Form.Mainnet,
	}
	reply, err = s.client.SaveRegistrationForm(context.TODO(), partial)
	require.NoError(err, "should not receive an error when saving a partial registration form")
	require.NotNil(reply, "uploaded form should be returned when a non-empty registration form is saved")
	require.Nil(reply.Errors, "expected no errors when saving a partially valid registration form")

	// Ensure that the complete form is now in the database
	org, err = s.DB().RetrieveOrganization(context.Background(), org.UUID())
	require.NoError(err, "could not retrieve updated org from database")
	require.NotEmpty(org.Registration.State.Started, "expected registration start time to be set")
	org.Registration.State.Started = ""
	require.True(proto.Equal(form.Form, org.Registration), "expected entire form in database to match the fixture form")

	// Load a form fixture with validation errors
	err = loadFixture("testdata/bad_registration_form.pb.json", form.Form)
	require.NoError(err, "could not load registration form fixture")

	// All the validation errors in the fixture
	verrs := map[api.RegistrationFormStep][]*api.FieldValidationError{
		api.StepBasicDetails: {{Field: records.FieldWebsite, Error: records.ErrMissingField.Error()}},
		api.StepLegalPerson: {
			{Field: "", Error: records.ErrMissingField.Error()},
			{Field: records.FieldEntity, Error: ""},
		},
		api.StepContacts: {{Field: records.FieldContactsTechnicalEmail, Error: records.ErrMissingField.Error()}},
		api.StepTRIXO:    {{Field: records.FieldTRIXOPrimaryNationalJurisdiction, Error: records.ErrMissingField.Error()}},
		api.StepTRISA:    {{Field: records.FieldTestNetCommonName, Error: records.ErrMissingField.Error()}},
	}

	// Saving each step of the form should only return the validation errors for that step
	testCases := []struct {
		step api.RegistrationFormStep
		form *records.RegistrationForm
		errs []*api.FieldValidationError
	}{
		{api.StepBasicDetails, &records.RegistrationForm{
			Website:          form.Form.Website,
			BusinessCategory: form.Form.BusinessCategory,
			VaspCategories:   form.Form.VaspCategories,
			EstablishedOn:    form.Form.EstablishedOn,
			OrganizationName: form.Form.OrganizationName,
		}, []*api.FieldValidationError{{Field: records.FieldWebsite, Error: records.ErrMissingField.Error()}}},
		{api.StepLegalPerson, &records.RegistrationForm{
			Entity: form.Form.Entity,
		}, []*api.FieldValidationError{{Field: "", Error: records.ErrMissingField.Error()},
			{Field: records.FieldEntity, Error: ""}}},
		{api.StepContacts, &records.RegistrationForm{
			Contacts: form.Form.Contacts,
		}, []*api.FieldValidationError{{Field: records.FieldContactsTechnicalEmail, Error: records.ErrMissingField.Error()}}},
		{api.StepTRIXO, &records.RegistrationForm{
			Trixo: form.Form.Trixo,
		}, []*api.FieldValidationError{{Field: records.FieldTRIXOPrimaryNationalJurisdiction, Error: records.ErrMissingField.Error()}}},
		{api.StepTRISA, &records.RegistrationForm{
			Testnet: form.Form.Testnet,
			Mainnet: form.Form.Mainnet,
		}, []*api.FieldValidationError{{Field: records.FieldTestNetCommonName, Error: records.ErrMissingField.Error()}}},
	}

	for _, tc := range testCases {
		partial.Step = tc.step
		partial.Form = tc.form
		reply, err = s.client.SaveRegistrationForm(context.Background(), partial)
		require.NoError(err, "should not receive an error when saving a partial registration form")
		require.NotNil(reply, "uploaded form should be returned when a non-empty registration form is saved")
		require.Equal(verrs[tc.step], reply.Errors, "wrong errors when saving form step %s", tc.step)
	}
}

func (s *bffTestSuite) TestResetRegisterForm() {
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
	_, err = s.client.ResetRegistrationForm(context.TODO(), nil)
	s.requireError(err, http.StatusForbidden, "csrf verification failed for request", "expected error when request is not CSRF protected")

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
	_, err = s.client.ResetRegistrationForm(context.TODO(), nil)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.ResetRegistrationForm(context.TODO(), nil)
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID and the server must not panic if it does not
	claims.Permissions = []string{auth.UpdateVASP}
	require.NoError(s.SetClientCredentials(claims), "could not create token without organizationID from claims")
	_, err = s.client.ResetRegistrationForm(context.TODO(), nil)
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	_, err = s.client.ResetRegistrationForm(context.TODO(), nil)
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Create an organization in the database
	org := &records.Organization{}
	_, err = s.DB().CreateOrganization(context.Background(), org)
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.DB().DeleteOrganization(context.Background(), org.UUID())
	}()

	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")

	// Should return an error if the step is not a valid step
	_, err = s.client.ResetRegistrationForm(context.TODO(), &api.RegistrationFormParams{Step: "invalid"})
	s.requireError(err, http.StatusBadRequest, "unknown registration form step \"invalid\"", "expected error when step is invalid")

	// Load the registration form on the organization
	org.Registration = form
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization with registration form")

	// Test deleting the entire form
	defaultForm := records.NewRegisterForm()
	rep, err := s.client.ResetRegistrationForm(context.Background(), nil)
	require.NoError(err, "should not receive an error when deleting a registration form")
	require.Nil(rep.Errors, "should not receive any validation errors when deleting a registration form")
	require.True(proto.Equal(defaultForm, rep.Form), "default form should be returned when a registration form is deleted")

	// Load the complete form back on the organization
	org.Registration = form
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization with registration form")

	// Test deleting specific steps in the form
	params := &api.RegistrationFormParams{Step: api.StepBasicDetails}
	rep, err = s.client.ResetRegistrationForm(context.Background(), params)
	require.NoError(err, "should not receive an error when deleting a registration form")
	require.Nil(rep.Errors, "should not receive any validation errors when deleting a registration form")
	require.Equal(defaultForm.Website, rep.Form.Website, "website should be reset when basic details are deleted")
	require.Equal(defaultForm.BusinessCategory, rep.Form.BusinessCategory, "business category should be reset when basic details are deleted")
	require.Equal(defaultForm.VaspCategories, rep.Form.VaspCategories, "vasp categories should be reset when basic details are deleted")
	require.Equal(defaultForm.EstablishedOn, rep.Form.EstablishedOn, "established on should be reset when basic details are deleted")
	require.Equal(defaultForm.OrganizationName, rep.Form.OrganizationName, "organization name should be reset when basic details are deleted")
	require.Nil(rep.Form.Entity, "entity should not be returned on basic details delete")
	require.Nil(rep.Form.Contacts, "contacts should not be returned on basic details delete")
	require.Nil(rep.Form.Trixo, "trixo should not be returned on basic details delete")
	require.Nil(rep.Form.Testnet, "testnet should not be returned on basic details delete")
	require.Nil(rep.Form.Mainnet, "mainnet should not be returned on basic details delete")

	params.Step = api.StepLegalPerson
	rep, err = s.client.ResetRegistrationForm(context.Background(), params)
	require.NoError(err, "should not receive an error when deleting a registration form")
	require.Nil(rep.Errors, "should not receive any validation errors when deleting a registration form")
	require.Empty(rep.Form.Website, "website should not be returned on legal person delete")
	require.Equal(defaultForm.Entity, rep.Form.Entity, "entity should be reset when legal person is deleted")
	require.Nil(rep.Form.Contacts, "contacts should not be returned on legal person delete")
	require.Nil(rep.Form.Trixo, "trixo should not be returned on legal person delete")
	require.Nil(rep.Form.Testnet, "testnet should not be returned on legal person delete")
	require.Nil(rep.Form.Mainnet, "mainnet should not be returned on legal person delete")

	params.Step = api.StepContacts
	rep, err = s.client.ResetRegistrationForm(context.Background(), params)
	require.NoError(err, "should not receive an error when deleting a registration form")
	require.Nil(rep.Errors, "should not receive any validation errors when deleting a registration form")
	require.Empty(rep.Form.Website, "website should not be returned on contacts delete")
	require.Empty(rep.Form.Entity, "entity should not be returned on contacts delete")
	require.Equal(defaultForm.Contacts, rep.Form.Contacts, "contacts should be reset when contacts are deleted")
	require.Nil(rep.Form.Trixo, "trixo should not be returned on contacts delete")
	require.Nil(rep.Form.Testnet, "testnet should not be returned on contacts delete")
	require.Nil(rep.Form.Mainnet, "mainnet should not be returned on contacts delete")

	params.Step = api.StepTRIXO
	rep, err = s.client.ResetRegistrationForm(context.Background(), params)
	require.NoError(err, "should not receive an error when deleting a registration form")
	require.Nil(rep.Errors, "should not receive any validation errors when deleting a registration form")
	require.Empty(rep.Form.Website, "website should not be returned on trixo delete")
	require.Empty(rep.Form.Entity, "entity should not be returned on trixo delete")
	require.Empty(rep.Form.Contacts, "contacts should not be returned on trixo delete")
	require.Equal(defaultForm.Trixo, rep.Form.Trixo, "trixo should be reset when trixo is deleted")
	require.Nil(rep.Form.Testnet, "testnet should not be returned on trixo delete")
	require.Nil(rep.Form.Mainnet, "mainnet should not be returned on trixo delete")

	params.Step = api.StepTRISA
	rep, err = s.client.ResetRegistrationForm(context.Background(), params)
	require.NoError(err, "should not receive an error when deleting a registration form")
	require.Nil(rep.Errors, "should not receive any validation errors when deleting a registration form")
	require.Empty(rep.Form.Website, "website should not be returned on trisa delete")
	require.Empty(rep.Form.Entity, "entity should not be returned on trisa delete")
	require.Empty(rep.Form.Contacts, "contacts should not be returned on trisa delete")
	require.Empty(rep.Form.Trixo, "trixo should not be returned on trisa delete")
	require.Equal(defaultForm.Testnet, rep.Form.Testnet, "testnet should be reset when trisa is deleted")
	require.Equal(defaultForm.Mainnet, rep.Form.Mainnet, "mainnet should be reset when trisa is deleted")

	// At this point the form in the database should match the default
	org, err = s.DB().RetrieveOrganization(context.Background(), org.UUID())
	require.NoError(err, "should not receive an error when retrieving the organization")
	require.True(proto.Equal(defaultForm, org.Registration), "registration form should match the default form")
}

func (s *bffTestSuite) TestSubmitRegistration() {
	var err error
	require := s.Require()
	defer s.ResetDB()

	// Test setup: create an organization with a valid registration form that has not
	// been submitted yet - at the end of the test both mainnet and testnet should be
	// submitted and the response from the directory updated on the organization.
	org := &records.Organization{}
	_, err = s.DB().CreateOrganization(context.Background(), org)
	require.NoError(err, "could not create organization in the database")

	// Save the registration form fixture on the organization
	org.Registration = &records.RegistrationForm{}
	require.NoError(loadFixture("testdata/registration_form.pb.json", org.Registration), "could not load registration form from the fixtures")
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization with registration form")

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
		claims.Permissions = []string{auth.UpdateVASP}
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

		// Load the valid register reply response
		reply := &gds.RegisterReply{}
		require.NoError(loadPBFixture(fmt.Sprintf("testdata/%s/register_reply.json", network), reply), "could not load register reply fixture for %s", network)

		// Test a valid register reply
		mgds.OnRegister = func(ctx context.Context, req *gds.RegisterRequest) (*gds.RegisterReply, error) {
			// TODO: We could do deeper validation, but we have GDS tests for that.
			// This makes sure that all the fields in the form are being passed to GDS.
			if req.Entity == nil {
				return nil, status.Error(codes.InvalidArgument, "legal person entity missing in request")
			}

			if req.Contacts == nil {
				return nil, status.Error(codes.InvalidArgument, "contacts missing in request")
			}

			if req.TrisaEndpoint == "" {
				return nil, status.Error(codes.InvalidArgument, "trisa endpoint missing in request")
			}

			if req.CommonName == "" {
				return nil, status.Error(codes.InvalidArgument, "common name missing in request")
			}

			if req.Website != org.Registration.Website {
				return nil, status.Error(codes.InvalidArgument, "wrong website in request")
			}

			if req.BusinessCategory != org.Registration.BusinessCategory {
				return nil, status.Error(codes.InvalidArgument, "wrong business category in request")
			}

			if len(req.VaspCategories) != len(org.Registration.VaspCategories) {
				return nil, errors.New("wrong vasp categories in request")
			}

			if req.EstablishedOn == "" {
				return nil, status.Error(codes.InvalidArgument, "established on date missing in request")
			}

			if req.Trixo == nil {
				return nil, status.Error(codes.InvalidArgument, "trixo questionnaire missing in request")
			}

			// Send the register reply back
			return reply, nil
		}

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
	org, err = s.DB().RetrieveOrganization(context.Background(), org.UUID())
	require.NoError(err, "could not update organization from the database")

	require.NotNil(org.Testnet, "missing testnet directory record after registration")
	require.Equal(org.Testnet.Id, "6041571e-09b4-47e7-870a-723f8032cd6c", "incorrect testnet directory id")
	require.Equal(org.Testnet.RegisteredDirectory, "trisatest.net", "incorrect testnet registered directory ")
	require.Equal(org.Testnet.CommonName, "test.trisa.example.ua", "incorrect testnet directory common name")
	require.NotEmpty(org.Testnet.Submitted, "expected testnet submitted timestamp stored in database")

	require.NotNil(org.Mainnet, "missing mainnet directory record after registration")
	require.Equal(org.Mainnet.Id, "5bafb054-5868-439e-9b3c-75db91810714", "incorrect mainnet directory id")
	require.Equal(org.Mainnet.RegisteredDirectory, "vaspdirectory.net", "incorrect mainnet registered directory ")
	require.Equal(org.Mainnet.CommonName, "trisa.example.ua", "incorrect mainnet directory common name")
	require.NotEmpty(org.Mainnet.Submitted, "expected mainnet submitted timestamp stored in database")

	// User metadata should be updated with the directory IDs
	appdata := &auth.AppMetadata{}
	require.NoError(appdata.Load(s.auth.GetUserAppMetadata()))
	require.Equal("6041571e-09b4-47e7-870a-723f8032cd6c", appdata.VASPs.TestNet, "incorrect testnet directory id in user metadata")
	require.Equal("5bafb054-5868-439e-9b3c-75db91810714", appdata.VASPs.MainNet, "incorrect mainnet directory id in user metadata")
}

func (s *bffTestSuite) TestSubmitRegistrationNotReady() {
	require := s.Require()

	// Ensure that a bad argument error is returned if the registration form is not
	// ready to submit. Create an organization that has a registration form without
	// network details and valid claims to access the record.
	org := &records.Organization{}
	_, err := s.DB().CreateOrganization(context.Background(), org)
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.DB().DeleteOrganization(context.Background(), org.UUID())
	}()

	// Ensure the registration is not ready to submit by removing mainnet and testnet
	org.Registration = &records.RegistrationForm{}
	require.NoError(loadFixture("testdata/registration_form.pb.json", org.Registration), "could not load registration form from the fixtures")
	org.Registration.Mainnet = nil
	org.Registration.Testnet = nil
	require.False(org.Registration.ReadyToSubmit("both"), "registration should not be ready to submit")

	// Save the registration form fixture on the organization
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization with registration form")

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
	org := &records.Organization{}
	_, err := s.DB().CreateOrganization(context.Background(), org)
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.DB().DeleteOrganization(context.Background(), org.UUID())
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
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization with registration form")

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
				Status:  pb.VerificationState_PENDING_REVIEW,
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

func (s *bffTestSuite) TestCertificates() {
	require := s.Require()
	defer s.ResetTestNetDB()
	defer s.ResetMainNetDB()

	// Load fixtures for testing
	uniform := &models.Certificate{}
	uniformFixture := filepath.Join("testdata", "testnet", "certs", "uniform.json")
	require.NoError(loadFixture(uniformFixture, uniform), "could not load uniform certificate fixture")
	testnetVASP := &pb.VASP{}
	testnetVASPFixture := filepath.Join("testdata", "testnet", "vasp.json")
	require.NoError(loadFixture(testnetVASPFixture, testnetVASP), "could not load testnet VASP fixture")

	victor := &models.Certificate{}
	victorFixture := filepath.Join("testdata", "mainnet", "certs", "victor.json")
	require.NoError(loadFixture(victorFixture, victor), "could not load victor certificate fixture")
	zulu := &models.Certificate{}
	zuluFixture := filepath.Join("testdata", "mainnet", "certs", "zulu.json")
	require.NoError(loadFixture(zuluFixture, zulu), "could not load zulu certificate fixture")
	mainnetVASP := &pb.VASP{}
	mainnetVASPFixture := filepath.Join("testdata", "mainnet", "vasp.json")
	require.NoError(loadFixture(mainnetVASPFixture, mainnetVASP), "could not load mainnet VASP fixture")

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		VASPs:       map[string]string{},
	}

	// Endpoint must be authenticated
	_, err := s.client.Certificates(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.Certificates(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Set valid credentials for the remainder of the tests
	claims.Permissions = []string{auth.ReadVASP}
	claims.VASPs["testnet"] = testnetVASP.Id
	claims.VASPs["mainnet"] = mainnetVASP.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token from valid credentials")

	// Test error message is populated when only testnet returns an error
	_, err = s.MainNetDB().CreateVASP(context.Background(), mainnetVASP)
	require.NoError(err, "could not create mainnet VASP")
	reply, err := s.client.Certificates(context.TODO())
	require.NoError(err, "expected no error when only testnet returns an error")
	require.Empty(reply.TestNet)
	require.Empty(reply.MainNet)
	require.NotEmpty(reply.Error.TestNet, "expected error message when only testnet returns an error")
	require.Empty(reply.Error.MainNet, "expected no error when mainnet returns a valid response")

	// Test error message is populated when only mainnet returns an error
	_, err = s.TestNetDB().CreateVASP(context.Background(), testnetVASP)
	require.NoError(err, "could not create testnet VASP")
	require.NoError(s.MainNetDB().DeleteVASP(context.Background(), mainnetVASP.Id), "could not delete VASP from mainnet database")
	reply, err = s.client.Certificates(context.TODO())
	require.NoError(err, "expected no error when only mainnet returns an error")
	require.Empty(reply.TestNet)
	require.Empty(reply.MainNet)
	require.NotEmpty(reply.Error.MainNet, "expected error message when only mainnet returns an error")
	require.Empty(reply.Error.TestNet, "expected no error when testnet returns a valid response")

	// Test empty results are returned even if there is no mainnet registration
	delete(claims.VASPs, "mainnet")
	require.NoError(s.SetClientCredentials(claims), "could not create token from valid credentials")
	reply, err = s.client.Certificates(context.TODO())
	require.NoError(err, "could not retrieve certificates")
	require.Empty(reply.TestNet, "expected no testnet certificates")
	require.Empty(reply.MainNet, "expected no mainnet certificates")
	require.Empty(reply.Error, "expected no errors")

	// Create certificate fixtures in the databases
	require.NoError(s.TestNetDB().UpdateCert(context.Background(), uniform), "could not create uniform certificate")
	require.NoError(models.AppendCertID(testnetVASP, uniform.Id), "could not append testnet certificate ID to VASP")
	require.NoError(s.TestNetDB().UpdateVASP(context.Background(), testnetVASP), "could not update testnet VASP")

	require.NoError(s.MainNetDB().UpdateCert(context.Background(), victor), "could not create victor certificate")
	require.NoError(models.AppendCertID(mainnetVASP, victor.Id), "could not append mainnet certificate ID to VASP")
	_, err = s.MainNetDB().CreateVASP(context.Background(), mainnetVASP)
	require.NoError(err, "could not create mainnet VASP")

	require.NoError(s.MainNetDB().UpdateCert(context.Background(), zulu), "could not create zulu certificate")
	require.NoError(models.AppendCertID(mainnetVASP, zulu.Id), "could not append mainnet certificate ID to VASP")
	require.NoError(s.MainNetDB().UpdateVASP(context.Background(), mainnetVASP), "could not update mainnet VASP")

	// Test certificates are returned from both testnet and mainnet
	claims.VASPs["mainnet"] = mainnetVASP.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token from valid credentials")
	reply, err = s.client.Certificates(context.TODO())
	require.NoError(err, "could not retrieve certificates")
	require.Empty(reply.Error, "expected no errors")
	require.Len(reply.TestNet, 1, "wrong number of testnet certificates")
	require.Len(reply.MainNet, 2, "wrong number of mainnet certificates")

	// Verify the testnet certificate fields
	expected := uniform.Details
	actual := reply.TestNet[0]
	require.Equal(uniform.Id, actual.SerialNumber, "expected testnet certificate serial to match")
	require.Equal(expected.NotBefore, actual.IssuedAt, "expected testnet certificate issued date to match")
	require.Equal(expected.NotAfter, actual.ExpiresAt, "expected testnet certificate expiration date to match")
	require.False(actual.Revoked, "expected testnet certificate to not be revoked")
	details, err := wire.Rewire(expected)
	require.NoError(err, "could not rewire uniform certificate details")
	require.Equal(details, actual.Details, "expected mainnet certificate details to match")

	// Both mainnet certificates should be returned
	require.Len(reply.MainNet, 2, "expected two mainnet certificates")
	for _, actual := range reply.MainNet {
		var expected *pb.Certificate
		switch actual.SerialNumber {
		case victor.Id:
			expected = victor.Details
		case zulu.Id:
			expected = zulu.Details
		default:
			require.Fail("unexpected mainnet certificate serial number", actual.SerialNumber)
		}

		// Compare the certificate data in the API response to the fixture certificate data
		require.Equal(expected.NotBefore, actual.IssuedAt, fmt.Sprintf("mainnet certificate %s issued date did not match", expected.SerialNumber))
		require.Equal(expected.NotAfter, actual.ExpiresAt, fmt.Sprintf("mainnet certificate %s expiration date did not match", expected.SerialNumber))
		require.Equal(expected.Revoked, actual.Revoked, fmt.Sprintf("mainnet certificate %s revoked bool did not match", expected.SerialNumber))
		details, err = wire.Rewire(expected)
		require.NoError(err, "could not rewire mainnet certificate details")
		require.Equal(details, actual.Details, fmt.Sprintf("mainnet certificate %s details did not match", expected.SerialNumber))
	}
}

func (s *bffTestSuite) TestAttention() {
	require := s.Require()
	defer s.ResetDB()
	defer s.ResetTestNetDB()
	defer s.ResetMainNetDB()

	// Load fixtures for testing
	testnetVASP := &pb.VASP{}
	mainnetVASP := &pb.VASP{}
	testnetFixture := filepath.Join("testdata", "testnet", "vasp.json")
	mainnetFixture := filepath.Join("testdata", "mainnet", "vasp.json")
	require.NoError(loadFixture(testnetFixture, testnetVASP))
	require.NoError(loadFixture(mainnetFixture, mainnetVASP))

	// Create an organization in the database with no registration form
	org := &records.Organization{}
	_, err := s.DB().CreateOrganization(context.Background(), org)
	require.NoError(err, "could not create organization in the database")

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		VASPs:       map[string]string{},
	}

	// Endpoint must be authenticated
	_, err = s.client.Attention(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.Attention(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{auth.ReadVASP}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")
	_, err = s.client.Attention(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	_, err = s.client.Attention(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Start registration message should be returned when there is no registration form
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	expected := &api.AttentionMessage{
		Message:  bff.StartRegistration,
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_START_REGISTRATION.String(),
	}
	reply, err := s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected start registration message")
	require.Equal(expected, reply.Messages[0], "expected start registration message")

	// Start registration message should still be returned if the registration form state is empty
	org.Registration = &records.RegistrationForm{}
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization in the database")
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected start registration message")
	require.Equal(expected, reply.Messages[0], "expected start registration message")

	// Start registration message should still be returned if the registration form has not been started
	org.Registration = records.NewRegisterForm()
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization in the database")
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected start registration message")
	require.Equal(expected, reply.Messages[0], "expected start registration message")

	// Complete registration message should be returned when the registration form has been started but not submitted
	org.Registration.State.Started = time.Now().Format(time.RFC3339)
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization in the database")
	expected = &api.AttentionMessage{
		Message:  bff.CompleteRegistration,
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_COMPLETE_REGISTRATION.String(),
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected complete registration message")
	require.Equal(expected, reply.Messages[0], "expected complete registration message")

	// Submit mainnet message should be returned when the registration form has been submitted only to testnet
	org.Testnet = &records.DirectoryRecord{
		Submitted: time.Now().Format(time.RFC3339),
	}
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization in the database")
	expected = &api.AttentionMessage{
		Message:  bff.SubmitMainnet,
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_SUBMIT_MAINNET.String(),
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected submit mainnet message")
	require.Equal(expected, reply.Messages[0], "expected submit mainnet message")

	// Submit testnet message should be returned when the registration form has been submitted only to mainnet
	org.Testnet.Submitted = ""
	org.Mainnet = &records.DirectoryRecord{
		Submitted: time.Now().Format(time.RFC3339),
	}
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization in the database")
	submitTestnet := &api.AttentionMessage{
		Message:  bff.SubmitTestnet,
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_SUBMIT_TESTNET.String(),
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 1, "expected submit testnet message")
	require.Equal(submitTestnet, reply.Messages[0], "expected submit testnet message")

	// Test an error is returned when VASP does not exist in testnet
	claims.VASPs["testnet"] = testnetVASP.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	_, err = s.client.Attention(context.TODO())
	s.requireError(err, http.StatusInternalServerError, storeerrors.ErrEntityNotFound.Error(), "expected error when VASP does not exist in testnet")

	// Test an error is returned when VASP does not exist in mainnet
	claims.VASPs["testnet"] = ""
	claims.VASPs["mainnet"] = mainnetVASP.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	_, err = s.client.Attention(context.TODO())
	s.requireError(err, http.StatusInternalServerError, storeerrors.ErrEntityNotFound.Error(), "expected error when VASP does not exist in mainnet")

	// Verify emails message should be returned when the VASP has been submitted but
	// emails are not yet verified
	mainnetVASP.VerificationStatus = pb.VerificationState_SUBMITTED
	_, err = s.MainNetDB().CreateVASP(context.Background(), mainnetVASP)
	require.NoError(err, "could not create VASP in the mainnet database")
	verifyMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.VerifyEmails, "MainNet"),
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_VERIFY_EMAILS.String(),
	}
	messages := []*api.AttentionMessage{
		submitTestnet,
		verifyMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Registration pending message should be returned when the VASP has been submitted
	// and is pending email verification
	mainnetVASP.VerificationStatus = pb.VerificationState_PENDING_REVIEW
	require.NoError(s.MainNetDB().UpdateVASP(context.Background(), mainnetVASP), "could not update VASP in the database")
	pendingMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.RegistrationPending, "MainNet"),
		Severity: records.AttentionSeverity_INFO.String(),
		Action:   records.AttentionAction_NO_ACTION.String(),
	}
	messages = []*api.AttentionMessage{
		submitTestnet,
		pendingMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Registration approved message should be returned when the VASP is verified
	mainnetVASP.VerificationStatus = pb.VerificationState_VERIFIED
	require.NoError(s.MainNetDB().UpdateVASP(context.Background(), mainnetVASP), "could not update VASP in the database")
	approvedMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.RegistrationApproved, "MainNet"),
		Severity: records.AttentionSeverity_SUCCESS.String(),
		Action:   records.AttentionAction_NO_ACTION.String(),
	}
	messages = []*api.AttentionMessage{
		submitTestnet,
		approvedMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Rejected message should be returned when the VASP state is rejected
	mainnetVASP.VerificationStatus = pb.VerificationState_REJECTED
	require.NoError(s.MainNetDB().UpdateVASP(context.Background(), mainnetVASP), "could not update VASP in the database")
	rejectMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.RegistrationRejected, "MainNet"),
		Severity: records.AttentionSeverity_ALERT.String(),
		Action:   records.AttentionAction_CONTACT_SUPPORT.String(),
	}
	messages = []*api.AttentionMessage{
		submitTestnet,
		rejectMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Revoked message should be returned when the certificate is revoked
	mainnetVASP.VerificationStatus = pb.VerificationState_VERIFIED
	mainnetVASP.IdentityCertificate = &pb.Certificate{
		Revoked: true,
	}
	require.NoError(s.MainNetDB().UpdateVASP(context.Background(), mainnetVASP), "could not update VASP in the database")
	revokedMainnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.CertificateRevoked, "MainNet"),
		Severity: records.AttentionSeverity_ALERT.String(),
		Action:   records.AttentionAction_CONTACT_SUPPORT.String(),
	}
	messages = []*api.AttentionMessage{
		submitTestnet,
		revokedMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Configure testnet fixture with expired certificate
	claims.VASPs["testnet"] = "alice0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	org.Testnet.Submitted = time.Now().Format(time.RFC3339)
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization in the database")
	expires := time.Now().AddDate(0, 0, 28)
	testnetVASP.VerificationStatus = pb.VerificationState_VERIFIED
	testnetVASP.IdentityCertificate.Revoked = false
	testnetVASP.IdentityCertificate.NotAfter = expires.Format(time.RFC3339)

	// Expired message should be returned when the certificate is expired
	_, err = s.TestNetDB().CreateVASP(context.Background(), testnetVASP)
	require.NoError(err, "could not create VASP in the testnet database")
	expiredTestnet := &api.AttentionMessage{
		Message:  fmt.Sprintf(bff.RenewCertificate, "TestNet", expires.Format("January 2, 2006")),
		Severity: records.AttentionSeverity_WARNING.String(),
		Action:   records.AttentionAction_RENEW_CERTIFICATE.String(),
	}
	messages = []*api.AttentionMessage{
		expiredTestnet,
		revokedMainnet,
	}
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Len(reply.Messages, 2, "wrong number of messages returned")
	require.ElementsMatch(messages, reply.Messages, "wrong messages returned")

	// Should return 204 when there are no attention messages
	claims.VASPs["testnet"] = ""
	claims.VASPs["mainnet"] = ""
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	reply, err = s.client.Attention(context.TODO())
	require.NoError(err, "received error from attention endpoint")
	require.Nil(reply, "expected nil reply")
}

func (s *bffTestSuite) TestRegistrationStatus() {
	require := s.Require()
	defer s.ResetDB()

	// Create an organization in the database with no directory records
	org := &records.Organization{}
	_, err := s.DB().CreateOrganization(context.Background(), org)
	require.NoError(err, "could not create organization in the database")

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
		VASPs:       map[string]string{},
	}

	// Endpoint must be authenticated
	_, err = s.client.RegistrationStatus(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.RegistrationStatus(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{auth.ReadVASP}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")
	_, err = s.client.RegistrationStatus(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic and should return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	_, err = s.client.RegistrationStatus(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when claims are valid but no organization is in the database")

	// Should return an empty response when there are no directory records
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid claims")
	reply, err := s.client.RegistrationStatus(context.TODO())
	require.NoError(err, "received error from registration status endpoint")
	require.Empty(reply, "expected empty response when there are no directory records")

	// Should return only the testnet timestamp when testnet registration has been submitted
	org.Testnet = &records.DirectoryRecord{
		Submitted: time.Now().Format(time.RFC3339),
	}
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization in the database")
	reply, err = s.client.RegistrationStatus(context.TODO())
	require.NoError(err, "received error from registration status endpoint")
	require.Equal(org.Testnet.Submitted, reply.TestNetSubmitted, "expected testnet timestamp to be returned")
	require.Empty(reply.MainNetSubmitted, "expected mainnet timestamp to be empty")

	// Should return only the mainnet timestamp when mainnet registration has been submitted
	org.Testnet.Submitted = ""
	org.Mainnet = &records.DirectoryRecord{
		Submitted: time.Now().Format(time.RFC3339),
	}
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization in the database")
	reply, err = s.client.RegistrationStatus(context.TODO())
	require.NoError(err, "received error from registration status endpoint")
	require.Equal(org.Mainnet.Submitted, reply.MainNetSubmitted, "expected mainnet timestamp to be returned")
	require.Empty(reply.TestNetSubmitted, "expected testnet timestamp to be empty")

	// Should return both timestamps when both registrations have been submitted
	org.Testnet.Submitted = time.Now().Format(time.RFC3339)
	org.Mainnet.Submitted = time.Now().Format(time.RFC3339)
	require.NoError(s.DB().UpdateOrganization(context.Background(), org), "could not update organization in the database")
	reply, err = s.client.RegistrationStatus(context.TODO())
	require.NoError(err, "received error from registration status endpoint")
	require.Equal(org.Testnet.Submitted, reply.TestNetSubmitted, "expected testnet timestamp to be returned")
	require.Equal(org.Mainnet.Submitted, reply.MainNetSubmitted, "expected mainnet timestamp to be returned")
}
