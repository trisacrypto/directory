package models_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff/models/v1"
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

	for i, tc := range testCases {
		form := RegistrationForm{}
		form.Website = tc.website
		form.BusinessCategory = tc.businessCategory
		form.VaspCategories = tc.vaspCategories
		form.EstablishedOn = tc.establishedOn
		form.OrganizationName = tc.orgName

		err := form.Validate(StepBasicDetails)
		if tc.errs == nil {
			require.NoError(t, err, "test case %d failed", i)
		} else {
			var verrs ValidationErrors
			require.ErrorAs(t, err, &verrs, "test case %d failed", i)
			require.Equal(t, tc.errs, verrs, "test case %d failed", i)
		}
	}
}

// Test validating the trixo questionnaire step of the registration form
func TestValidateTRIXO(t *testing.T) {
	testCases := []struct {
		trixo *pb.TRIXOQuestionnaire
		errs  ValidationErrors
	}{
		{nil, ValidationErrors{
			{Field: FieldTRIXO, Err: ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  "Yes",
			HasRequiredRegulatoryProgram: "Yes",
		}, ValidationErrors{
			{Field: FieldTRIXOPrimaryNationalJurisdiction, Err: ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			FinancialTransfersPermitted:  "No",
			HasRequiredRegulatoryProgram: "No",
		}, ValidationErrors{
			{Field: FieldTRIXOPrimaryNationalJurisdiction, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOPrimaryRegulator, Err: ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			HasRequiredRegulatoryProgram: "yes",
		}, ValidationErrors{
			{Field: FieldTRIXOPrimaryNationalJurisdiction, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOPrimaryRegulator, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOFinancialTransfersPermitted, Err: ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			FinancialTransfersPermitted:  "idk",
			HasRequiredRegulatoryProgram: "YES",
		}, ValidationErrors{
			{Field: FieldTRIXOPrimaryNationalJurisdiction, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOPrimaryRegulator, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOFinancialTransfersPermitted, Err: ErrInvalidField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			HasRequiredRegulatoryProgram: "NO",
			OtherJurisdictions: []*pb.Jurisdiction{
				{Country: "FR", RegulatorName: "AMF", LicenseNumber: "123"},
				{RegulatorName: "FinCEN", LicenseNumber: "456"},
				{Country: "US", LicenseNumber: "456"},
				{Country: "US"},
			},
		}, ValidationErrors{
			{Field: FieldTRIXOPrimaryNationalJurisdiction, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOPrimaryRegulator, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOFinancialTransfersPermitted, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOOtherJurisdictionsCountry, Err: ErrMissingField.Error(), Index: 1},
			{Field: FieldTRIXOOtherJurisdictionsRegulatorName, Err: ErrMissingField.Error(), Index: 2},
			{Field: FieldTRIXOOtherJurisdictionsRegulatorName, Err: ErrMissingField.Error(), Index: 3},
			{Field: FieldTRIXOOtherJurisdictionsLicenseNumber, Err: ErrMissingField.Error(), Index: 3},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction: "US",
			PrimaryRegulator:            "FinCEN",
			FinancialTransfersPermitted: "Yes",
		}, ValidationErrors{
			{Field: FieldTRIXOHasRequiredRegulatoryProgram, Err: ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction:  "US",
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  "Yes",
			HasRequiredRegulatoryProgram: "idk",
		}, ValidationErrors{
			{Field: FieldTRIXOHasRequiredRegulatoryProgram, Err: ErrInvalidField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction:  "US",
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  " Yes ",
			HasRequiredRegulatoryProgram: " Yes ",
			ConductsCustomerKyc:          true,
			KycThreshold:                 -1,
		}, ValidationErrors{
			{Field: FieldTRIXOKYCThreshold, Err: ErrNegativeValue.Error()},
			{Field: FieldTRIXOKYCThresholdCurrency, Err: ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction:  "US",
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  "Partially",
			HasRequiredRegulatoryProgram: "Yes",
			MustComplyTravelRule:         true,
			ComplianceThreshold:          -1,
		}, ValidationErrors{
			{Field: FieldTRIXOApplicableRegulations, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOComplianceThreshold, Err: ErrNegativeValue.Error()},
			{Field: FieldTRIXOComplianceThresholdCurrency, Err: ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction:  "US",
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  "Partially",
			HasRequiredRegulatoryProgram: "Yes",
			MustComplyTravelRule:         true,
			ApplicableRegulations:        []string{"Reg1", ""},
			ComplianceThreshold:          -1,
			ComplianceThresholdCurrency:  "USD",
		}, ValidationErrors{
			{Field: FieldTRIXOApplicableRegulations, Err: ErrMissingField.Error(), Index: 1},
			{Field: FieldTRIXOComplianceThreshold, Err: ErrNegativeValue.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction:  "US",
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  "yes",
			HasRequiredRegulatoryProgram: "yes",
			ConductsCustomerKyc:          true,
			KycThreshold:                 1000,
			KycThresholdCurrency:         "USD",
			MustComplyTravelRule:         true,
			ApplicableRegulations:        []string{"Reg1", "Reg2"},
			ComplianceThreshold:          1000,
			ComplianceThresholdCurrency:  "USD",
		}, nil},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction:  "US",
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  "Yes",
			HasRequiredRegulatoryProgram: "Yes",
			KycThreshold:                 -1,
			ComplianceThreshold:          -1,
		}, nil},
	}

	for i, tc := range testCases {
		form := RegistrationForm{
			Trixo: tc.trixo,
		}

		err := form.Validate(StepTRIXO)
		if tc.errs == nil {
			require.NoError(t, err, "test case %d failed", i)
		} else {
			var verrs ValidationErrors
			require.ErrorAs(t, err, &verrs, "test case %d failed", i)
			require.Equal(t, tc.errs, verrs, "test case %d failed", i)
		}
	}
}

// Test validating the contacts step of the registration form
func TestValidateContacts(t *testing.T) {
	// A single error should be returned for nil contacts
	form := RegistrationForm{}
	expected := ValidationErrors{
		{Field: FieldContacts, Err: ErrMissingField.Error()},
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
			{Field: FieldContacts, Err: ErrNoContacts.Error()},
		}},
		// Only technical provided should nominate admin/legal contact to be populated
		{contact, nil, nil, nil, ValidationErrors{
			{Field: FieldContacts, Err: ErrMissingContact.Error()},
			{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrLegal.Error()},
			{Field: FieldContactsLegal, Err: ErrMissingAdminOrLegal.Error()},
		}},
		// Only legal provided should nominate admin/technical contact to be populated
		{nil, nil, contact, nil, ValidationErrors{
			{Field: FieldContacts, Err: ErrMissingContact.Error()},
			{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrTechnical.Error()},
			{Field: FieldContactsTechnical, Err: ErrMissingAdminOrTechnical.Error()},
		}},
		// Only billing provided should nominate admin/technical/legal contact to be populated
		{nil, nil, nil, contact, ValidationErrors{
			{Field: FieldContacts, Err: ErrMissingContact.Error()},
			{Field: FieldContactsAdministrative, Err: ErrMissingContact.Error()},
			{Field: FieldContactsTechnical, Err: ErrMissingContact.Error()},
			{Field: FieldContactsLegal, Err: ErrMissingContact.Error()},
		}},
		// Legal and billing provided should nominate admin/technical contact to be populated
		{nil, nil, contact, contact, ValidationErrors{
			{Field: FieldContacts, Err: ErrMissingContact.Error()},
			{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrTechnical.Error()},
			{Field: FieldContactsTechnical, Err: ErrMissingAdminOrTechnical.Error()},
		}},
		// Technical and billing provided should nominate admin/legal contact to be populated
		{contact, nil, nil, contact, ValidationErrors{
			{Field: FieldContacts, Err: ErrMissingContact.Error()},
			{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrLegal.Error()},
			{Field: FieldContactsLegal, Err: ErrMissingAdminOrLegal.Error()},
		}},
		{missingEmail, contact, contact, contact, ValidationErrors{
			{Field: FieldContactsTechnicalEmail, Err: ErrMissingField.Error()},
		}},
		{missingEmail, missingEmail, contact, contact, ValidationErrors{
			{Field: FieldContactsTechnicalEmail, Err: ErrMissingField.Error()},
			{Field: FieldContactsAdministrativeEmail, Err: ErrMissingField.Error()},
		}},
		{missingEmail, missingEmail, missingEmail, contact, ValidationErrors{
			{Field: FieldContactsTechnicalEmail, Err: ErrMissingField.Error()},
			{Field: FieldContactsAdministrativeEmail, Err: ErrMissingField.Error()},
			{Field: FieldContactsLegalEmail, Err: ErrMissingField.Error()},
		}},
		{missingEmail, missingEmail, missingEmail, missingEmail, ValidationErrors{
			{Field: FieldContactsTechnicalEmail, Err: ErrMissingField.Error()},
			{Field: FieldContactsAdministrativeEmail, Err: ErrMissingField.Error()},
			{Field: FieldContactsLegalEmail, Err: ErrMissingField.Error()},
			{Field: FieldContactsBillingEmail, Err: ErrMissingField.Error()},
		}},
		// Only admin provided is valid
		{nil, contact, nil, nil, nil},
		// Admin and billing provided is valid
		{nil, contact, nil, contact, nil},
		// Technical and legal provided is valid
		{contact, nil, contact, nil, nil},
		// Providing all contacts is valid
		{contact, contact, contact, contact, nil},
	}

	for i, tc := range testCases {
		form := RegistrationForm{
			Contacts: &pb.Contacts{
				Technical:      tc.technical,
				Administrative: tc.admin,
				Legal:          tc.legal,
				Billing:        tc.billing,
			},
		}

		err := form.Validate(StepContacts)
		if tc.errs == nil {
			require.NoError(t, err, "test case %d failed", i)
		} else {
			var verrs ValidationErrors
			require.ErrorAs(t, err, &verrs, "test case %d failed", i)
			require.Equal(t, tc.errs, verrs, "test case %d failed", i)
		}
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

	for i, tc := range testCases {
		err := ValidateContact(tc.contact, "admin")
		if tc.errs == nil {
			require.NoError(t, err, "test case %d failed", i)
		} else {
			var verrs ValidationErrors
			require.ErrorAs(t, err, &verrs, "test case %d failed", i)
			require.Equal(t, tc.errs, verrs, "test case %d failed", i)
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
