package models_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

func TestStepType(t *testing.T) {
	testCases := []struct {
		s        string
		expected StepType
	}{
		{"", StepNone},
		{"  ", StepNone},
		{"\n", StepNone},
		{"all", StepAll},
		{"ALL", StepAll},
		{"  All  ", StepAll},
		{"basic", StepBasicDetails},
		{"BASIC   ", StepBasicDetails},
		{"Basic", StepBasicDetails},
		{"legal", StepLegalPerson},
		{" LEGal\n", StepLegalPerson},
		{"LEGAL\t", StepLegalPerson},
		{"contacts", StepContacts},
		{"trixo", StepTRIXO},
		{"trisa", StepTRISA},
	}

	for _, tc := range testCases {
		actual, err := ParseStepType(tc.s)
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
		actual, err := ParseStepType(tc)
		require.Error(t, err)
		require.True(t, strings.HasPrefix(err.Error(), "unknown registration form step"))
		require.Equal(t, StepNone, actual)
	}

}

// Test validating the basic details step of the registration form
func TestValidateBasicDetails(t *testing.T) {
	testCases := []struct {
		website          string
		businessCategory pb.BusinessCategory
		vaspCategories   []string
		establishedOn    string
		orgName          string
		errs             ValidationErrors
	}{
		{"", pb.BusinessCategory_BUSINESS_ENTITY, []string{"P2P"}, "2021-01-01", "Example, Inc.", ValidationErrors{
			{Field: "website", Err: ErrMissingField.Error()},
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, []string{"P2P"}, "2021-01-01", "Example, Inc.", ValidationErrors{
			{Field: "website", Err: ErrMissingField.Error()},
			{Field: "business_category", Err: ErrMissingField.Error()},
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, nil, "2021-01-01", "Example, Inc.", ValidationErrors{
			{Field: "website", Err: ErrMissingField.Error()},
			{Field: "business_category", Err: ErrMissingField.Error()},
			{Field: "vasp_categories", Err: ErrMissingField.Error()},
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, nil, "", "Example Inc.", ValidationErrors{
			{Field: "website", Err: ErrMissingField.Error()},
			{Field: "business_category", Err: ErrMissingField.Error()},
			{Field: "vasp_categories", Err: ErrMissingField.Error()},
			{Field: "established_on", Err: ErrMissingField.Error()},
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, nil, "", "", ValidationErrors{
			{Field: "website", Err: ErrMissingField.Error()},
			{Field: "business_category", Err: ErrMissingField.Error()},
			{Field: "vasp_categories", Err: ErrMissingField.Error()},
			{Field: "established_on", Err: ErrMissingField.Error()},
			{Field: "organization_name", Err: ErrMissingField.Error()},
		}},
		{"example.com", pb.BusinessCategory_GOVERNMENT_ENTITY, []string{}, " ", " ", ValidationErrors{
			{Field: "vasp_categories", Err: ErrMissingField.Error()},
			{Field: "established_on", Err: ErrMissingField.Error()},
			{Field: "organization_name", Err: ErrMissingField.Error()},
		}},
		{" ", pb.BusinessCategory_GOVERNMENT_ENTITY, []string{}, " ", " ", ValidationErrors{
			{Field: "website", Err: ErrMissingField.Error()},
			{Field: "vasp_categories", Err: ErrMissingField.Error()},
			{Field: "established_on", Err: ErrMissingField.Error()},
			{Field: "organization_name", Err: ErrMissingField.Error()},
		}},
		{"example.com", pb.BusinessCategory_GOVERNMENT_ENTITY, []string{"Exchange"}, "2021-01-01", "Example, Inc.", nil},
	}

	for _, tc := range testCases {
		form := RegistrationForm{}
		form.Website = tc.website
		form.BusinessCategory = tc.businessCategory
		form.VaspCategories = tc.vaspCategories
		form.EstablishedOn = tc.establishedOn
		form.OrganizationName = tc.orgName

		errs := form.Validate(StepBasicDetails)
		if tc.errs == nil {
			require.NoError(t, errs)
		} else {
			var verrs ValidationErrors
			require.ErrorAs(t, errs, &verrs)
			require.Equal(t, tc.errs, verrs)
		}
	}
}

func TestValidateLegalPerson(t *testing.T) {
	testCases := []struct {
		entity *ivms101.LegalPerson
		errs   ValidationErrors
	}{
		{nil, ValidationErrors{{Field: FieldEntity, Err: ErrMissingField.Error()}}},
		{
			&ivms101.LegalPerson{},
			ValidationErrors{
				{Field: FieldEntityName, Err: ErrMissingField.Error()},
				{Field: FieldEntityGeographicAddresses, Err: ErrNoGeographicAddress.Error()},
				{Field: FieldEntityNationalIdentification, Err: ErrLegalNatIDRequired.Error()},
				{Field: FieldEntityCountryOfRegistration, Err: ErrMissingField.Error()},
				{Field: FieldEntity, Err: "one or more legal person name identifiers is required"},
			},
		},
		{
			&ivms101.LegalPerson{
				Name: &ivms101.LegalPersonName{
					NameIdentifiers: []*ivms101.LegalPersonNameId{
						{
							LegalPersonName:               "Wayne Enterprises, LTD",
							LegalPersonNameIdentifierType: ivms101.LegalPersonLegal,
						},
					},
				},
				GeographicAddresses: []*ivms101.Address{
					{
						AddressType: ivms101.AddressTypeBusiness,
						AddressLine: []string{
							"1 Wayne Tower",
							"Gotham City, NJ 08302",
						},
						Country: "US",
					},
				},
				NationalIdentification: &ivms101.NationalIdentification{
					NationalIdentifier:     "ZGWO00PIA5JMETFLPG72",
					NationalIdentifierType: ivms101.NationalIdentifierLEIX,
				},
				CountryOfRegistration: "US",
			},
			nil,
		},
	}

	for i, tc := range testCases {
		form := &RegistrationForm{Entity: tc.entity}
		err := form.Validate(StepLegalPerson)

		if len(tc.errs) > 0 {
			var valid ValidationErrors
			require.ErrorAs(t, err, &valid, "expected validation errors in test case %d", i)
			require.Len(t, valid, len(tc.errs), "expected same number of validation errors in test case %d", i)
			require.Equal(t, tc.errs, valid, "expected same validation errors in test case %d", i)
		} else {
			require.NoError(t, err, "expected fully valid entity on test case %d", i)
		}
	}
}

// Test that the registration form marshals and unmarshals correctly to and from JSON
func TestMarshalRegistrationForm(t *testing.T) {
	// Load the JSON fixture
	fixtureData, err := os.ReadFile("testdata/default_registration_form.json")
	require.NoError(t, err, "error reading default registration form fixture")

	// Default form should be marshaled correctly
	form := NewRegisterForm()
	data, err := json.Marshal(form)
	require.NoError(t, err, "error marshaling registration form to JSON")
	require.JSONEq(t, string(fixtureData), string(data), "default registration form does not match fixture")

	// Default form should be unmarshaled correctly
	result := &RegistrationForm{}
	require.NoError(t, json.Unmarshal(data, result), "error unmarshaling registration form from JSON")
	require.True(t, proto.Equal(form, result), "registration form should be unmarshaled correctly")

	// Modified form should be marshaled correctly
	form.Contacts.Administrative.Email = "admin@example.com"
	data, err = json.Marshal(form)
	require.NoError(t, err, "error marshaling registration form to JSON")

	// Modified form should be unmarshaled correctly
	result = &RegistrationForm{}
	require.NoError(t, json.Unmarshal(data, result), "error unmarshaling registration form from JSON")
	require.True(t, proto.Equal(form, result), "registration form should be unmarshaled correctly")
}

// Test validating the contacts step of the registration form
func TestValidateContacts(t *testing.T) {
	// A single error should be returned for nil contacts
	form := RegistrationForm{}
	expected := ValidationErrors{
		{Field: "contacts", Err: ErrMissingField.Error()},
	}
	verrs := form.ValidateContacts()
	require.Equal(t, expected, verrs, "expected a single error for missing contacts field")

	errs := form.Validate(StepContacts)
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

	// Test that contacts are valided if not empty
	testCases := []struct {
		technical *pb.Contact
		admin     *pb.Contact
		legal     *pb.Contact
		billing   *pb.Contact
		errs      ValidationErrors
	}{
		// No contacts provided
		{nil, nil, &pb.Contact{}, nil, ValidationErrors{
			{Field: "contacts", Err: ErrNoContacts.Error()},
		}},
		// Only technical provided should nominate admin/legal contact to be populated
		{contact, nil, nil, nil, ValidationErrors{
			{Field: "contacts", Err: ErrMissingContact.Error()},
			{Field: "contacts.administrative", Err: ErrMissingAdminOrLegal.Error()},
			{Field: "contacts.legal", Err: ErrMissingAdminOrLegal.Error()},
		}},
		// Only legal provided should nominate admin/technical contact to be populated
		{nil, nil, contact, nil, ValidationErrors{
			{Field: "contacts", Err: ErrMissingContact.Error()},
			{Field: "contacts.administrative", Err: ErrMissingAdminOrTechnical.Error()},
			{Field: "contacts.technical", Err: ErrMissingAdminOrTechnical.Error()},
		}},
		// Only billing provided should nominate admin/technical/legal contact to be populated
		{nil, nil, nil, contact, ValidationErrors{
			{Field: "contacts", Err: ErrMissingContact.Error()},
			{Field: "contacts.administrative", Err: ErrMissingContact.Error()},
			{Field: "contacts.technical", Err: ErrMissingContact.Error()},
			{Field: "contacts.legal", Err: ErrMissingContact.Error()},
		}},
		// Legal and billing provided should nominate admin/technical contact to be populated
		{nil, nil, contact, contact, ValidationErrors{
			{Field: "contacts", Err: ErrMissingContact.Error()},
			{Field: "contacts.administrative", Err: ErrMissingAdminOrTechnical.Error()},
			{Field: "contacts.technical", Err: ErrMissingAdminOrTechnical.Error()},
		}},
		// Technical and billing provided should nominate admin/legal contact to be populated
		{contact, nil, nil, contact, ValidationErrors{
			{Field: "contacts", Err: ErrMissingContact.Error()},
			{Field: "contacts.administrative", Err: ErrMissingAdminOrLegal.Error()},
			{Field: "contacts.legal", Err: ErrMissingAdminOrLegal.Error()},
		}},
		{missingEmail, contact, contact, contact, ValidationErrors{
			{Field: "contacts.technical.email", Err: ErrMissingField.Error()},
		}},
		{missingEmail, missingEmail, contact, contact, ValidationErrors{
			{Field: "contacts.technical.email", Err: ErrMissingField.Error()},
			{Field: "contacts.administrative.email", Err: ErrMissingField.Error()},
		}},
		{missingEmail, missingEmail, missingEmail, contact, ValidationErrors{
			{Field: "contacts.technical.email", Err: ErrMissingField.Error()},
			{Field: "contacts.administrative.email", Err: ErrMissingField.Error()},
			{Field: "contacts.legal.email", Err: ErrMissingField.Error()},
		}},
		{missingEmail, missingEmail, missingEmail, missingEmail, ValidationErrors{
			{Field: "contacts.technical.email", Err: ErrMissingField.Error()},
			{Field: "contacts.administrative.email", Err: ErrMissingField.Error()},
			{Field: "contacts.legal.email", Err: ErrMissingField.Error()},
			{Field: "contacts.billing.email", Err: ErrMissingField.Error()},
		}},
		// Only admin provided is valid
		{nil, contact, nil, nil, nil},
		// Admin and billing provided is valid
		{nil, contact, nil, contact, nil},
		// Technical and legal provided is valid
		{contact, contact, contact, contact, nil},
	}

	for _, tc := range testCases {
		form := RegistrationForm{
			Contacts: &pb.Contacts{
				Technical:      tc.technical,
				Administrative: tc.admin,
				Legal:          tc.legal,
				Billing:        tc.billing,
			},
		}

		verrs := form.ValidateContacts()
		require.Equal(t, tc.errs, verrs)

		errs := form.Validate(StepContacts)
		require.Equal(t, tc.errs, errs)
	}
}

// Test validating a single contact
func TestValidateContact(t *testing.T) {
	testCases := []struct {
		contact *pb.Contact
		errs    ValidationErrors
	}{
		{&pb.Contact{Name: "L", Email: " leopold.wentzel@gmail.com ", Phone: "555-867-5309"}, ValidationErrors{
			{Field: "admin.name", Err: ErrTooShort.Error()},
		}},
		{&pb.Contact{Email: "not an email", Phone: " 555-867-5309 "}, ValidationErrors{
			{Field: "admin.name", Err: ErrMissingField.Error()},
			{Field: "admin.email", Err: ErrInvalidEmail.Error()},
		}},
		{&pb.Contact{Phone: "555-867-5309"}, ValidationErrors{
			{Field: "admin.name", Err: ErrMissingField.Error()},
			{Field: "admin.email", Err: ErrMissingField.Error()},
		}},
		{&pb.Contact{}, ValidationErrors{
			{Field: "admin.name", Err: ErrMissingField.Error()},
			{Field: "admin.email", Err: ErrMissingField.Error()},
			{Field: "admin.phone", Err: ErrMissingField.Error()},
		}},
		{&pb.Contact{Name: "Leopold Wentzel", Email: "leopold.wentzel@gmail.com", Phone: "555-867-5309"}, nil},
	}

	for _, tc := range testCases {
		errs := ValidateContact(tc.contact, "admin")
		require.Equal(t, tc.errs, errs)
	}
}
