package models_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

func TestStepType(t *testing.T) {
	testCases := []struct {
		s        string
		expected models.StepType
	}{
		{"", models.StepNone},
		{"  ", models.StepNone},
		{"\n", models.StepNone},
		{"all", models.StepAll},
		{"ALL", models.StepAll},
		{"  All  ", models.StepAll},
		{"basic", models.StepBasicDetails},
		{"BASIC   ", models.StepBasicDetails},
		{"Basic", models.StepBasicDetails},
		{"legal", models.StepLegalPerson},
		{" LEGal\n", models.StepLegalPerson},
		{"LEGAL\t", models.StepLegalPerson},
		{"contacts", models.StepContacts},
		{"trixo", models.StepTRIXO},
		{"trisa", models.StepTRISA},
	}

	for _, tc := range testCases {
		actual, err := models.ParseStepType(tc.s)
		require.NoError(t, err, "expected valid step type to be parsed")
		require.Equal(t, tc.expected, actual)
		require.Equal(t, tc.expected.String(), actual.String())
	}

	invalidTestCases := []string{
		"foo",
		"FOO",
		"  not a real step",
		"TRISAXXAA",
	}

	for _, tc := range invalidTestCases {
		actual, err := models.ParseStepType(tc)
		require.Error(t, err)
		require.True(t, strings.HasPrefix(err.Error(), "unknown registration form step"))
		require.Equal(t, models.StepNone, actual)
	}

}

// Test that the registration form marshals and unmarshals correctly to and from JSON
func TestMarshalRegistrationForm(t *testing.T) {
	// Load the JSON fixture
	fixtureData, err := os.ReadFile("testdata/default_registration_form.json")
	require.NoError(t, err, "error reading default registration form fixture")

	// Default form should be marshaled correctly
	form := models.NewRegisterForm()
	data, err := json.Marshal(form)
	require.NoError(t, err, "error marshaling registration form to JSON")
	require.JSONEq(t, string(fixtureData), string(data), "default registration form does not match fixture")

	// Default form should be unmarshaled correctly
	result := &models.RegistrationForm{}
	require.NoError(t, json.Unmarshal(data, result), "error unmarshaling registration form from JSON")
	require.True(t, proto.Equal(form, result), "registration form should be unmarshaled correctly")

	// Modified form should be marshaled correctly
	form.Contacts.Administrative.Email = "admin@example.com"
	data, err = json.Marshal(form)
	require.NoError(t, err, "error marshaling registration form to JSON")

	// Modified form should be unmarshaled correctly
	result = &models.RegistrationForm{}
	require.NoError(t, json.Unmarshal(data, result), "error unmarshaling registration form from JSON")
	require.True(t, proto.Equal(form, result), "registration form should be unmarshaled correctly")
}

// Test validating the contacts step of the registration form
func TestValidateContacts(t *testing.T) {
	// A single error should be returned for nil contacts
	form := models.RegistrationForm{}
	expected := models.ValidationErrors{
		{Field: "contacts", Err: models.ErrMissingField.Error()},
	}
	verrs := form.ValidateContacts()
	require.Equal(t, expected, verrs, "expected a single error for missing contacts field")

	errs := form.Validate(models.StepContacts)
	require.Equal(t, expected, errs, "expected a single error for missing contacts field")

	contact := &pb.Contact{
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
		Phone: "555-867-5309",
	}

	missingEmail := &pb.Contact{
		Name:  contact.Name,
		Phone: contact.Phone,
	}

	// Test that all contacts are validated
	testCases := []struct {
		technical *pb.Contact
		admin     *pb.Contact
		legal     *pb.Contact
		billing   *pb.Contact
		errs      models.ValidationErrors
	}{
		{nil, nil, nil, nil, models.ValidationErrors{
			{Field: "contacts.technical", Err: models.ErrMissingField.Error()},
			{Field: "contacts.administrative", Err: models.ErrMissingField.Error()},
			{Field: "contacts.legal", Err: models.ErrMissingField.Error()},
			{Field: "contacts.billing", Err: models.ErrMissingField.Error()},
		}},
		{missingEmail, contact, contact, contact, models.ValidationErrors{
			{Field: "contacts.technical.email", Err: models.ErrMissingField.Error()},
		}},
		{missingEmail, missingEmail, contact, contact, models.ValidationErrors{
			{Field: "contacts.technical.email", Err: models.ErrMissingField.Error()},
			{Field: "contacts.administrative.email", Err: models.ErrMissingField.Error()},
		}},
		{missingEmail, missingEmail, missingEmail, contact, models.ValidationErrors{
			{Field: "contacts.technical.email", Err: models.ErrMissingField.Error()},
			{Field: "contacts.administrative.email", Err: models.ErrMissingField.Error()},
			{Field: "contacts.legal.email", Err: models.ErrMissingField.Error()},
		}},
		{missingEmail, missingEmail, missingEmail, missingEmail, models.ValidationErrors{
			{Field: "contacts.technical.email", Err: models.ErrMissingField.Error()},
			{Field: "contacts.administrative.email", Err: models.ErrMissingField.Error()},
			{Field: "contacts.legal.email", Err: models.ErrMissingField.Error()},
			{Field: "contacts.billing.email", Err: models.ErrMissingField.Error()},
		}},
		{contact, contact, contact, contact, nil},
	}

	for _, tc := range testCases {
		form := models.RegistrationForm{
			Contacts: &pb.Contacts{
				Technical:      tc.technical,
				Administrative: tc.admin,
				Legal:          tc.legal,
				Billing:        tc.billing,
			},
		}

		verrs := form.ValidateContacts()
		require.Equal(t, tc.errs, verrs)

		errs := form.Validate(models.StepContacts)
		require.Equal(t, tc.errs, errs)
	}
}

// Test validating a single contact
func TestValidateContact(t *testing.T) {
	testCases := []struct {
		contact *pb.Contact
		errs    models.ValidationErrors
	}{
		{nil, models.ValidationErrors{{Field: "admin", Err: models.ErrMissingField.Error()}}},
		{&pb.Contact{Email: "leopold.wentzel@gmail.com", Phone: "555-867-5309"}, models.ValidationErrors{
			{Field: "admin.name", Err: models.ErrMissingField.Error()},
		}},
		{&pb.Contact{Name: "L", Email: " leopold.wentzel@gmail.com ", Phone: "555-867-5309"}, models.ValidationErrors{
			{Field: "admin.name", Err: models.ErrTooShort.Error()},
		}},
		{&pb.Contact{Email: "not an email", Phone: " 555-867-5309 "}, models.ValidationErrors{
			{Field: "admin.name", Err: models.ErrMissingField.Error()},
			{Field: "admin.email", Err: models.ErrInvalidEmail.Error()},
		}},
		{&pb.Contact{Phone: "555-867-5309"}, models.ValidationErrors{
			{Field: "admin.name", Err: models.ErrMissingField.Error()},
			{Field: "admin.email", Err: models.ErrMissingField.Error()},
		}},
		{&pb.Contact{}, models.ValidationErrors{
			{Field: "admin.name", Err: models.ErrMissingField.Error()},
			{Field: "admin.email", Err: models.ErrMissingField.Error()},
			{Field: "admin.phone", Err: models.ErrMissingField.Error()},
		}},
		{&pb.Contact{Name: "Leopold Wentzel", Email: "leopold.wentzel@gmail.com", Phone: "555-867-5309"}, nil},
	}

	for _, tc := range testCases {
		errs := models.ValidateContact(tc.contact, "admin")
		require.Equal(t, tc.errs, errs)
	}
}