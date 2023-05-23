package models_test

import (
	"encoding/json"
	"fmt"
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

// Test validating the basic details step of the registration form
func TestValidateBasicDetails(t *testing.T) {
	testCases := []struct {
		website          string
		businessCategory pb.BusinessCategory
		vaspCategories   []string
		establishedOn    string
		orgName          string
		errs             models.ValidationErrors
	}{
		{"", pb.BusinessCategory_BUSINESS_ENTITY, []string{"P2P"}, "2021-01-01", "Example, Inc.", models.ValidationErrors{
			{Field: "website", Err: models.ErrMissingField.Error()},
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, []string{"P2P"}, "2021-01-01", "Example, Inc.", models.ValidationErrors{
			{Field: "website", Err: models.ErrMissingField.Error()},
			{Field: "business_category", Err: models.ErrMissingField.Error()},
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, nil, "2021-01-01", "Example, Inc.", models.ValidationErrors{
			{Field: "website", Err: models.ErrMissingField.Error()},
			{Field: "business_category", Err: models.ErrMissingField.Error()},
			{Field: "vasp_categories", Err: models.ErrMissingField.Error()},
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, nil, "", "Example Inc.", models.ValidationErrors{
			{Field: "website", Err: models.ErrMissingField.Error()},
			{Field: "business_category", Err: models.ErrMissingField.Error()},
			{Field: "vasp_categories", Err: models.ErrMissingField.Error()},
			{Field: "established_on", Err: models.ErrMissingField.Error()},
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, nil, "", "", models.ValidationErrors{
			{Field: "website", Err: models.ErrMissingField.Error()},
			{Field: "business_category", Err: models.ErrMissingField.Error()},
			{Field: "vasp_categories", Err: models.ErrMissingField.Error()},
			{Field: "established_on", Err: models.ErrMissingField.Error()},
			{Field: "organization_name", Err: models.ErrMissingField.Error()},
		}},
		{"example.com", pb.BusinessCategory_GOVERNMENT_ENTITY, []string{}, " ", " ", models.ValidationErrors{
			{Field: "vasp_categories", Err: models.ErrMissingField.Error()},
			{Field: "established_on", Err: models.ErrMissingField.Error()},
			{Field: "organization_name", Err: models.ErrMissingField.Error()},
		}},
		{" ", pb.BusinessCategory_GOVERNMENT_ENTITY, []string{}, " ", " ", models.ValidationErrors{
			{Field: "website", Err: models.ErrMissingField.Error()},
			{Field: "vasp_categories", Err: models.ErrMissingField.Error()},
			{Field: "established_on", Err: models.ErrMissingField.Error()},
			{Field: "organization_name", Err: models.ErrMissingField.Error()},
		}},
		{"example.com", pb.BusinessCategory_GOVERNMENT_ENTITY, []string{"Exchange"}, "2021-01-01", "Example, Inc.", nil},
	}

	for _, tc := range testCases {
		form := models.RegistrationForm{}
		form.Website = tc.website
		form.BusinessCategory = tc.businessCategory
		form.VaspCategories = tc.vaspCategories
		form.EstablishedOn = tc.establishedOn
		form.OrganizationName = tc.orgName

		verrs := form.ValidateBasicDetails()
		fmt.Println(tc)
		require.Equal(t, tc.errs, verrs)

		errs := form.Validate(models.StepBasicDetails)
		require.Equal(t, tc.errs, errs)
	}
}

// Test validating the trixo questionnaire step of the registration form
func TestValidateTRIXO(t *testing.T) {
	testCases := []struct {
		trixo *pb.TRIXOQuestionnaire
		errs  models.ValidationErrors
	}{
		{nil, models.ValidationErrors{
			{Field: "trixo", Err: models.ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  "Yes",
			HasRequiredRegulatoryProgram: "Yes",
		}, models.ValidationErrors{
			{Field: "trixo.primary_national_jurisdiction", Err: models.ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			FinancialTransfersPermitted:  "No",
			HasRequiredRegulatoryProgram: "No",
		}, models.ValidationErrors{
			{Field: "trixo.primary_national_jurisdiction", Err: models.ErrMissingField.Error()},
			{Field: "trixo.primary_regulator", Err: models.ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			HasRequiredRegulatoryProgram: "yes",
		}, models.ValidationErrors{
			{Field: "trixo.primary_national_jurisdiction", Err: models.ErrMissingField.Error()},
			{Field: "trixo.primary_regulator", Err: models.ErrMissingField.Error()},
			{Field: "trixo.financial_transfers_permitted", Err: models.ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			FinancialTransfersPermitted:  "idk",
			HasRequiredRegulatoryProgram: "YES",
		}, models.ValidationErrors{
			{Field: "trixo.primary_national_jurisdiction", Err: models.ErrMissingField.Error()},
			{Field: "trixo.primary_regulator", Err: models.ErrMissingField.Error()},
			{Field: "trixo.financial_transfers_permitted", Err: models.ErrInvalidField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			HasRequiredRegulatoryProgram: "NO",
			OtherJurisdictions: []*pb.Jurisdiction{
				{Country: "FR", RegulatorName: "AMF", LicenseNumber: "123"},
				{RegulatorName: "FinCEN", LicenseNumber: "456"},
				{Country: "US", LicenseNumber: "456"},
				{Country: "US"},
			},
		}, models.ValidationErrors{
			{Field: "trixo.primary_national_jurisdiction", Err: models.ErrMissingField.Error()},
			{Field: "trixo.primary_regulator", Err: models.ErrMissingField.Error()},
			{Field: "trixo.financial_transfers_permitted", Err: models.ErrMissingField.Error()},
			{Field: "trixo.other_jurisdictions.country", Err: models.ErrMissingField.Error(), Index: 1},
			{Field: "trixo.other_jurisdictions.regulator_name", Err: models.ErrMissingField.Error(), Index: 2},
			{Field: "trixo.other_jurisdictions.regulator_name", Err: models.ErrMissingField.Error(), Index: 3},
			{Field: "trixo.other_jurisdictions.license_number", Err: models.ErrMissingField.Error(), Index: 3},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction: "US",
			PrimaryRegulator:            "FinCEN",
			FinancialTransfersPermitted: "Yes",
		}, models.ValidationErrors{
			{Field: "trixo.has_required_regulatory_program", Err: models.ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction:  "US",
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  "Yes",
			HasRequiredRegulatoryProgram: "idk",
		}, models.ValidationErrors{
			{Field: "trixo.has_required_regulatory_program", Err: models.ErrInvalidField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction:  "US",
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  " Yes ",
			HasRequiredRegulatoryProgram: " Yes ",
			ConductsCustomerKyc:          true,
			KycThreshold:                 -1,
		}, models.ValidationErrors{
			{Field: "trixo.kyc_threshold", Err: models.ErrNegativeValue.Error()},
			{Field: "trixo.kyc_threshold_currency", Err: models.ErrMissingField.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			PrimaryNationalJurisdiction:  "US",
			PrimaryRegulator:             "FinCEN",
			FinancialTransfersPermitted:  "Partially",
			HasRequiredRegulatoryProgram: "Yes",
			MustComplyTravelRule:         true,
			ComplianceThreshold:          -1,
		}, models.ValidationErrors{
			{Field: "trixo.applicable_regulations", Err: models.ErrMissingField.Error()},
			{Field: "trixo.compliance_threshold", Err: models.ErrNegativeValue.Error()},
			{Field: "trixo.compliance_threshold_currency", Err: models.ErrMissingField.Error()},
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
		}, models.ValidationErrors{
			{Field: "trixo.applicable_regulations", Err: models.ErrMissingField.Error(), Index: 1},
			{Field: "trixo.compliance_threshold", Err: models.ErrNegativeValue.Error()},
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

	for _, tc := range testCases {
		form := models.RegistrationForm{
			Trixo: tc.trixo,
		}
		fmt.Println(tc.trixo)

		verrs := form.ValidateTRIXO()
		require.Equal(t, tc.errs, verrs)

		errs := form.Validate(models.StepTRIXO)
		require.Equal(t, tc.errs, errs)
	}
}
