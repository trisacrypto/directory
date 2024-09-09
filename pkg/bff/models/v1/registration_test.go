package models_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff/models/v1"
	ivms101 "github.com/trisacrypto/trisa/pkg/ivms101"
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
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, nil, "", "Example Inc.", ValidationErrors{
			{Field: "website", Err: ErrMissingField.Error()},
			{Field: "business_category", Err: ErrMissingField.Error()},
			{Field: "established_on", Err: ErrMissingField.Error()},
		}},
		{"", pb.BusinessCategory_UNKNOWN_ENTITY, nil, "", "", ValidationErrors{
			{Field: "website", Err: ErrMissingField.Error()},
			{Field: "business_category", Err: ErrMissingField.Error()},
			{Field: "established_on", Err: ErrMissingField.Error()},
			{Field: "organization_name", Err: ErrMissingField.Error()},
		}},
		{"example.com", pb.BusinessCategory_GOVERNMENT_ENTITY, []string{}, " ", " ", ValidationErrors{
			{Field: "established_on", Err: ErrMissingField.Error()},
			{Field: "organization_name", Err: ErrMissingField.Error()},
		}},
		{" ", pb.BusinessCategory_GOVERNMENT_ENTITY, []string{"P2P", " "}, " ", " ", ValidationErrors{
			{Field: "website", Err: ErrMissingField.Error()},
			{Field: "vasp_categories", Err: ErrMissingField.Error(), Index: 1},
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

func TestValidateLegalPerson(t *testing.T) {
	testCases := []struct {
		entity *ivms101.LegalPerson
		errs   ValidationErrors
	}{
		{nil, ValidationErrors{{Field: FieldEntity, Err: ErrMissingField.Error()}}},
		{
			&ivms101.LegalPerson{},
			ValidationErrors{
				{Field: "entity.name", Err: "ivms101: missing name: this field is required"},
				{Field: FieldEntityGeographicAddresses, Err: ErrNoGeographicAddress.Error()},
				{Field: FieldEntityNationalIdentification, Err: ErrLegalNatIDRequired.Error()},
				{Field: FieldEntityCountryOfRegistration, Err: ErrMissingField.Error()},
			},
		},
		// Test C9 constraint is ignored but still return an error for missing RA
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
					RegistrationAuthority:  "RA777777",
				},
				CountryOfRegistration: "US",
			},
			ValidationErrors{
				{
					Field: "entity.nationalIdentification.registrationAuthority",
					Err:   "ivms101: invalid field nationalIdentification.registrationAuthority: registration authority not allowed for national identifier type code LEIX"},
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
			PrimaryRegulator:             "FinCEN",
			PrimaryNationalJurisdiction:  "USA",
			FinancialTransfersPermitted:  "Yes",
			HasRequiredRegulatoryProgram: "Yes",
		}, ValidationErrors{
			{Field: FieldTRIXOPrimaryNationalJurisdiction, Err: ErrInvalidCountry.Error()},
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
			{Field: FieldTRIXOFinancialTransfersPermitted, Err: ErrYesNoPartially.Error()},
		}},
		{&pb.TRIXOQuestionnaire{
			HasRequiredRegulatoryProgram: "NO",
			OtherJurisdictions: []*pb.Jurisdiction{
				{Country: "FR", RegulatorName: "AMF", LicenseNumber: "123"},
				{RegulatorName: "FinCEN", LicenseNumber: "456"},
				{Country: "US", LicenseNumber: "456"},
				{Country: "US"},
				{Country: "USA", RegulatorName: "FinCEN", LicenseNumber: "456"},
			},
		}, ValidationErrors{
			{Field: FieldTRIXOPrimaryNationalJurisdiction, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOPrimaryRegulator, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOFinancialTransfersPermitted, Err: ErrMissingField.Error()},
			{Field: FieldTRIXOOtherJurisdictionsCountry, Err: ErrMissingField.Error(), Index: 1},
			{Field: FieldTRIXOOtherJurisdictionsRegulatorName, Err: ErrMissingField.Error(), Index: 2},
			{Field: FieldTRIXOOtherJurisdictionsRegulatorName, Err: ErrMissingField.Error(), Index: 3},
			{Field: FieldTRIXOOtherJurisdictionsCountry, Err: ErrInvalidCountry.Error(), Index: 4},
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
			{Field: FieldTRIXOHasRequiredRegulatoryProgram, Err: ErrYesNo.Error()},
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
		// Missing phone number for technical contact is valid
		{&pb.Contact{Name: contact.Name, Email: contact.Email}, nil, contact, nil, nil},
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
		contactField string
		contact      *pb.Contact
		errs         ValidationErrors
	}{
		{FieldContactsAdministrative, &pb.Contact{Name: "L", Email: " leopold.wentzel@gmail.com ", Phone: "555-867-5309"}, ValidationErrors{
			{Field: FieldContactsAdministrativeName, Err: ErrTooShort.Error()},
		}},
		{FieldContactsAdministrative, &pb.Contact{Email: "not an email", Phone: " 555-867-5309 "}, ValidationErrors{
			{Field: FieldContactsAdministrativeName, Err: ErrMissingField.Error()},
			{Field: FieldContactsAdministrativeEmail, Err: ErrInvalidEmail.Error()},
		}},
		{FieldContactsAdministrative, &pb.Contact{Phone: "555-867-5309"}, ValidationErrors{
			{Field: FieldContactsAdministrativeName, Err: ErrMissingField.Error()},
			{Field: FieldContactsAdministrativeEmail, Err: ErrMissingField.Error()},
		}},
		{FieldContactsAdministrative, &pb.Contact{}, ValidationErrors{
			{Field: FieldContactsAdministrativeName, Err: ErrMissingField.Error()},
			{Field: FieldContactsAdministrativeEmail, Err: ErrMissingField.Error()},
			{Field: FieldContactsAdministrativePhone, Err: ErrMissingField.Error()},
		}},
		{FieldContactsAdministrative, &pb.Contact{Name: "Leopold Wentzel", Email: "leopold.wentzel@gmail.com", Phone: "555-867-5309"}, nil},
		{FieldContactsTechnical, &pb.Contact{Name: "Lt. Commander Data", Email: "data@enterpriseD.com"}, nil},
	}

	for i, tc := range testCases {
		err := ValidateContact(tc.contact, tc.contactField)
		if tc.errs == nil {
			require.NoError(t, err, "test case %d failed", i)
		} else {
			var verrs ValidationErrors
			require.ErrorAs(t, err, &verrs, "test case %d failed", i)
			require.Equal(t, tc.errs, verrs, "test case %d failed", i)
		}
	}
}

// Test validating the TRISA implementation details
func TestValidateTRISA(t *testing.T) {
	validNetwork := &NetworkDetails{
		Endpoint:   "main.trisa.io:443",
		CommonName: "main.trisa.io",
	}

	testCases := []struct {
		testnet *NetworkDetails
		mainnet *NetworkDetails
		errs    ValidationErrors
	}{
		{nil, nil, ValidationErrors{
			{Field: FieldTestNet, Err: ErrMissingTestNetOrMainNet.Error()},
			{Field: FieldMainNet, Err: ErrMissingTestNetOrMainNet.Error()},
		}},
		{&NetworkDetails{}, &NetworkDetails{}, ValidationErrors{
			{Field: FieldTestNet, Err: ErrMissingTestNetOrMainNet.Error()},
			{Field: FieldMainNet, Err: ErrMissingTestNetOrMainNet.Error()},
		}},
		{&NetworkDetails{CommonName: "test.trisa.io"}, validNetwork, ValidationErrors{
			{Field: FieldTestNetEndpoint, Err: ErrMissingField.Error()},
			{Field: FieldTestNetCommonName, Err: ErrCommonNameMismatch.Error()},
		}},
		{&NetworkDetails{Endpoint: "not an address", CommonName: "test.trisa.io"}, nil, ValidationErrors{
			{Field: FieldTestNetEndpoint, Err: ErrInvalidEndpoint.Error()},
			{Field: FieldTestNetCommonName, Err: ErrCommonNameMismatch.Error()},
		}},
		{&NetworkDetails{Endpoint: ":443", CommonName: "test.trisa.io"}, nil, ValidationErrors{
			{Field: FieldTestNetEndpoint, Err: ErrMissingHost.Error()},
			{Field: FieldTestNetCommonName, Err: ErrCommonNameMismatch.Error()},
		}},
		{&NetworkDetails{Endpoint: "test.trisa.io:", CommonName: "test.trisa.io"}, nil, ValidationErrors{
			{Field: FieldTestNetEndpoint, Err: ErrMissingPort.Error()},
		}},
		{&NetworkDetails{Endpoint: "test.trisa.io:foo", CommonName: "test.trisa.io"}, nil, ValidationErrors{
			{Field: FieldTestNetEndpoint, Err: ErrInvalidPort.Error()},
		}},
		{&NetworkDetails{Endpoint: "test.trisa.io:443"}, nil, ValidationErrors{
			{Field: FieldTestNetCommonName, Err: ErrMissingField.Error()},
		}},
		{&NetworkDetails{Endpoint: "test.trisa.io:443", CommonName: "*.trisa.io"}, nil, ValidationErrors{
			{Field: FieldTestNetCommonName, Err: ErrInvalidCommonName.Error()},
			{Field: FieldTestNetCommonName, Err: ErrCommonNameMismatch.Error()},
		}},
		{&NetworkDetails{Endpoint: "test.trisa.io:443", CommonName: "main.trisa.io"}, nil, ValidationErrors{
			{Field: FieldTestNetCommonName, Err: ErrCommonNameMismatch.Error()},
		}},
		{&NetworkDetails{Endpoint: "test.trisa.io:443", CommonName: "main.trisa.io", DnsNames: []string{"alt.trisa.io", "", "*.trisa.io", "https://trisa.io"}}, validNetwork, ValidationErrors{
			{Field: FieldTestNetCommonName, Err: ErrCommonNameMismatch.Error()},
			{Field: FieldTestNetDNSNames, Err: ErrMissingField.Error(), Index: 1},
			{Field: FieldTestNetDNSNames, Err: ErrInvalidCommonName.Error(), Index: 2},
			{Field: FieldTestNetDNSNames, Err: ErrInvalidCommonName.Error(), Index: 3},
		}},
		{&NetworkDetails{Endpoint: "test.trisa.io:443", CommonName: "test.trisa.io"}, &NetworkDetails{Endpoint: "main.trisa.io:443"}, ValidationErrors{
			{Field: FieldMainNetCommonName, Err: ErrMissingField.Error()},
		}},
		{validNetwork, validNetwork, ValidationErrors{
			{Field: FieldMainNetEndpoint, Err: ErrDuplicateEndpoint.Error()},
		}},
		{&NetworkDetails{Endpoint: "test.trisa.io:443", CommonName: "test.trisa.io"}, nil, nil},
		{nil, validNetwork, nil},
		{&NetworkDetails{Endpoint: "test.trisa.io:443", CommonName: "test.trisa.io"}, validNetwork, nil},
	}

	for i, tc := range testCases {
		form := RegistrationForm{
			Testnet: tc.testnet,
			Mainnet: tc.mainnet,
		}

		err := form.Validate(StepTRISA)
		if tc.errs == nil {
			require.NoError(t, err, "test case %d failed", i)
		} else {
			var verrs ValidationErrors
			require.ErrorAs(t, err, &verrs, "test case %d failed", i)
			require.Equal(t, tc.errs, verrs, "test case %d failed", i)
		}
	}
}

func loadJSONFixture(path string, v interface{}) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

// Test updating a registration form
func TestUpdateRegistrationForm(t *testing.T) {
	form := NewRegisterForm()

	// Load the registration form fixture
	update := &RegistrationForm{}
	err := loadJSONFixture("testdata/registration_form.json", update)
	require.NoError(t, err, "error loading registration form fixture")

	// An error should be returned if an unknown step is provided
	err = form.Update(update, "invalid")
	require.EqualError(t, err, "unknown step \"invalid\"", "error should be returned for unknown step")

	// Test updating the basic details step
	update.State.Current = 1
	err = form.Update(update, StepBasicDetails)
	require.NoError(t, err, "error updating basic details step")
	require.Equal(t, form.State, update.State, "state should be updated")

	// Test updating the legal person step
	update.State.Current = 2
	err = form.Update(update, StepLegalPerson)
	require.NoError(t, err, "error updating legal person step")
	require.Equal(t, form.State, update.State, "state should be updated")

	// Test updating the contacts step
	update.State.Current = 3
	err = form.Update(update, StepContacts)
	require.NoError(t, err, "error updating contacts step")
	require.Equal(t, form.State, update.State, "state should be updated")

	// Test updating the TRIXO step
	update.State.Current = 4
	err = form.Update(update, StepTRIXO)
	require.NoError(t, err, "error updating TRIXO step")
	require.Equal(t, form.State, update.State, "state should be updated")

	// Test updating the TRISA step
	update.State.Current = 5
	err = form.Update(update, StepTRISA)
	require.NoError(t, err, "error updating TRISA step")
	require.Equal(t, form.State, update.State, "state should be updated")

	// At this point the form should be fully updated
	require.True(t, proto.Equal(form, update), "form should be fully updated")

	// Test updating the entire form with no step
	form = NewRegisterForm()
	err = form.Update(update, StepNone)
	require.NoError(t, err, "error updating entire form")
	require.True(t, proto.Equal(form, update), "form should be fully updated")

	// Test updating the entire form with the all step
	form = NewRegisterForm()
	err = form.Update(update, StepAll)
	require.NoError(t, err, "error updating entire form")
	require.True(t, proto.Equal(form, update), "form should be fully updated")
}

// Test updating a form with validation errors
func TestUpdateRegistrationFormErrors(t *testing.T) {
	form := NewRegisterForm()

	// Load a registration form fixture that has validation errors
	update := &RegistrationForm{}
	err := loadJSONFixture("testdata/bad_registration_form.json", update)
	require.NoError(t, err, "error loading bad registration form fixture")

	// All the validation errors
	verrs := map[StepType]ValidationErrors{
		StepBasicDetails: {
			{Field: FieldBusinessCategory, Err: ErrMissingField.Error()},
		},
		StepLegalPerson: {
			{
				Field: "entity.country_of_registration",
				Err:   ErrMissingField.Error()},
		},
		StepContacts: {
			{Field: FieldContacts, Err: ErrMissingContact.Error()},
			{Field: FieldContactsAdministrative, Err: ErrMissingAdminOrTechnical.Error()},
			{Field: FieldContactsTechnical, Err: ErrMissingAdminOrTechnical.Error()},
		},
		StepTRIXO: {
			{Field: FieldTRIXOKYCThreshold, Err: ErrNegativeValue.Error()},
		},
		StepTRISA: {
			{Field: FieldTestNetCommonName, Err: ErrCommonNameMismatch.Error()},
		},
	}

	// Test updating the steps individually
	for step, verrs := range verrs {
		err = form.Update(update, step)
		require.Equal(t, verrs, err, "wrong validation errors for step %s", step)
	}

	// Test updating the entire form with no step
	allErrs := ValidationErrors{}
	for _, verrs := range verrs {
		allErrs = append(allErrs, verrs...)
	}
	form = NewRegisterForm()
	err = form.Update(update, StepNone)
	require.ElementsMatch(t, allErrs, err, "wrong validation errors for entire form")

	// Test updating the entire form with the all step
	form = NewRegisterForm()
	err = form.Update(update, StepAll)
	require.ElementsMatch(t, allErrs, err, "wrong validation errors for entire form")
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

func TestMarshalRegistrationFormStep(t *testing.T) {
	fixtureData, err := os.ReadFile("testdata/registration_form.json")
	require.NoError(t, err, "error reading default registration form fixture")

	form := &RegistrationForm{}
	err = form.UnmarshalJSON(fixtureData)
	require.NoError(t, err, "error marshaling registration form to JSON")

	t.Run("All", func(t *testing.T) {
		for _, step := range []StepType{StepNone, StepAll} {
			formData, err := form.MarshalJSON()
			require.NoError(t, err, "could not marshal form data")

			stepData, err := form.MarshalStepJSON(step)
			require.NoError(t, err, "could not marshal form step data for step %s", step)

			require.JSONEq(t, string(formData), string(stepData))
		}
	})

	makeStepTest := func(step StepType, keys ...string) func(t *testing.T) {
		// Ensure the keys always has the state
		keys = append(keys, FieldState)

		return func(t *testing.T) {
			data, err := form.MarshalStepJSON(step)
			require.NoError(t, err, "could not marshal json for step %s", step)

			var reply map[string]interface{}
			err = json.Unmarshal(data, &reply)
			require.NoError(t, err, "could not unmarshal json for step %s", step)

			require.Len(t, reply, len(keys), "expected reply to have expected number of keys for step %s", step)
			for _, key := range keys {
				require.Contains(t, reply, key, "expected reply to contain key for step %s", step)
			}
		}
	}

	t.Run("Basic", makeStepTest(StepBasicDetails, FieldWebsite, FieldBusinessCategory, FieldVASPCategories, FieldEstablishedOn, FieldOrganizationName))
	t.Run("Legal", makeStepTest(StepLegalPerson, FieldEntity))
	t.Run("Contacts", makeStepTest(StepContacts, FieldContacts))
	t.Run("TRIXO", makeStepTest(StepTRIXO, FieldTRIXO))
	t.Run("TRISA", makeStepTest(StepTRISA, FieldMainNet, FieldTestNet))
}
