package models_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	ivms101 "github.com/trisacrypto/trisa/pkg/ivms101"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func TestOrganizationKey(t *testing.T) {
	org := &models.Organization{}
	require.Equal(t, uuid.Nil[:], org.Key(), "expected nil uuid when organization id is empty string")

	uu := uuid.New()
	org.Id = uu.String()
	require.Equal(t, uu[:], org.Key(), "expected key to be uuid bytes")

	require.Panics(t, func() {
		org.Id = "notauuid"
		org.Key()
	}, "if the organization id is not a uuid string, expect a panic")
}

func TestParseOrgID(t *testing.T) {
	example := uuid.New()

	testCases := []struct {
		expected uuid.UUID
		input    interface{}
		err      error
	}{
		{example, example.String(), nil},       // parse string
		{example, example[:], nil},             // parse bytes
		{example, example, nil},                // parse uuid
		{uuid.Nil, 14, models.ErrInvalidOrgID}, // unknown type
	}

	for i, tc := range testCases {
		uu, err := models.ParseOrgID(tc.input)
		if tc.err != nil {
			require.Equal(t, uuid.Nil, uu, "expected nil uuid in test case %d", i)
			require.ErrorIs(t, err, tc.err, "expected error in test case %d", i)
		} else {
			require.NoError(t, err, "expected no error to occur in test case %d", i)
			require.Equal(t, tc.expected, uu, "unexpected org id returned when parsed in test case %d", i)
		}
	}
}

func TestReadyToSubmit(t *testing.T) {
	testCases := []struct {
		r        *models.RegistrationForm
		assert   require.BoolAssertionFunc
		networks []string
		message  string
	}{
		{
			r:        &models.RegistrationForm{},
			assert:   require.False,
			networks: []string{"testnet", "mainnet", "all", "both", ""},
			message:  "empty registration form should not be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				Entity:   &ivms101.LegalPerson{CountryOfRegistration: "GY"},
				Contacts: &gds.Contacts{Technical: &gds.Contact{Email: "jdoe@example.com"}},
				Trixo:    &gds.TRIXOQuestionnaire{PrimaryNationalJurisdiction: "GY"},
				Testnet:  &models.NetworkDetails{CommonName: "test.trisa.example.com"},
				Mainnet:  &models.NetworkDetails{CommonName: "trisa.example.com"},
			},
			assert:   require.False,
			networks: []string{"testnet", "mainnet", "all", "both", ""},
			message:  "registration form without vasp categories should not be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				VaspCategories: []string{"P2P", "other"},
				Contacts:       &gds.Contacts{Technical: &gds.Contact{Email: "jdoe@example.com"}},
				Trixo:          &gds.TRIXOQuestionnaire{PrimaryNationalJurisdiction: "GY"},
				Testnet:        &models.NetworkDetails{CommonName: "test.trisa.example.com"},
				Mainnet:        &models.NetworkDetails{CommonName: "trisa.example.com"},
			},
			assert:   require.False,
			networks: []string{"testnet", "mainnet", "all", "both", ""},
			message:  "registration form without entity should not be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				VaspCategories: []string{"P2P", "other"},
				Entity:         &ivms101.LegalPerson{CountryOfRegistration: "GY"},
				Trixo:          &gds.TRIXOQuestionnaire{PrimaryNationalJurisdiction: "GY"},
				Testnet:        &models.NetworkDetails{CommonName: "test.trisa.example.com"},
				Mainnet:        &models.NetworkDetails{CommonName: "trisa.example.com"},
			},
			assert:   require.False,
			networks: []string{"testnet", "mainnet", "all", "both", ""},
			message:  "registration form without contacts should not be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				VaspCategories: []string{"P2P", "other"},
				Entity:         &ivms101.LegalPerson{CountryOfRegistration: "GY"},
				Contacts:       &gds.Contacts{Technical: &gds.Contact{Email: "jdoe@example.com"}},
				Testnet:        &models.NetworkDetails{CommonName: "test.trisa.example.com"},
				Mainnet:        &models.NetworkDetails{CommonName: "trisa.example.com"},
			},
			assert:   require.False,
			networks: []string{"testnet", "mainnet", "all", "both", ""},
			message:  "registration form without trixo should not be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				VaspCategories: []string{"P2P", "other"},
				Entity:         &ivms101.LegalPerson{CountryOfRegistration: "GY"},
				Contacts:       &gds.Contacts{Technical: &gds.Contact{Email: "jdoe@example.com"}},
				Trixo:          &gds.TRIXOQuestionnaire{PrimaryNationalJurisdiction: "GY"},
			},
			assert:   require.False,
			networks: []string{"testnet", "mainnet", "all", "both", ""},
			message:  "registration form without network details should not be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				VaspCategories: []string{"P2P", "other"},
				Entity:         &ivms101.LegalPerson{CountryOfRegistration: "GY"},
				Contacts:       &gds.Contacts{Technical: &gds.Contact{Email: "jdoe@example.com"}},
				Trixo:          &gds.TRIXOQuestionnaire{PrimaryNationalJurisdiction: "GY"},
				Testnet:        &models.NetworkDetails{CommonName: "test.trisa.example.com"},
			},
			assert:   require.False,
			networks: []string{"mainnet", "all", "both", ""},
			message:  "registration form with testnet network details should not be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				VaspCategories: []string{"P2P", "other"},
				Entity:         &ivms101.LegalPerson{CountryOfRegistration: "GY"},
				Contacts:       &gds.Contacts{Technical: &gds.Contact{Email: "jdoe@example.com"}},
				Trixo:          &gds.TRIXOQuestionnaire{PrimaryNationalJurisdiction: "GY"},
				Testnet:        &models.NetworkDetails{CommonName: "test.trisa.example.com"},
			},
			assert:   require.True,
			networks: []string{"testnet"},
			message:  "registration form with testnet network details should be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				VaspCategories: []string{"P2P", "other"},
				Entity:         &ivms101.LegalPerson{CountryOfRegistration: "GY"},
				Contacts:       &gds.Contacts{Technical: &gds.Contact{Email: "jdoe@example.com"}},
				Trixo:          &gds.TRIXOQuestionnaire{PrimaryNationalJurisdiction: "GY"},
				Mainnet:        &models.NetworkDetails{CommonName: "trisa.example.com"},
			},
			assert:   require.False,
			networks: []string{"testnet", "all", "both", ""},
			message:  "registration form with mainnet network details should not be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				VaspCategories: []string{"P2P", "other"},
				Entity:         &ivms101.LegalPerson{CountryOfRegistration: "GY"},
				Contacts:       &gds.Contacts{Technical: &gds.Contact{Email: "jdoe@example.com"}},
				Trixo:          &gds.TRIXOQuestionnaire{PrimaryNationalJurisdiction: "GY"},
				Mainnet:        &models.NetworkDetails{CommonName: "trisa.example.com"},
			},
			assert:   require.True,
			networks: []string{"mainnet"},
			message:  "registration form with mainnet network details should be ready to submit",
		},
		{
			r: &models.RegistrationForm{
				VaspCategories: []string{"P2P", "other"},
				Entity:         &ivms101.LegalPerson{CountryOfRegistration: "GY"},
				Contacts:       &gds.Contacts{Technical: &gds.Contact{Email: "jdoe@example.com"}},
				Trixo:          &gds.TRIXOQuestionnaire{PrimaryNationalJurisdiction: "GY"},
				Testnet:        &models.NetworkDetails{CommonName: "test.trisa.example.com"},
				Mainnet:        &models.NetworkDetails{CommonName: "trisa.example.com"},
			},
			assert:   require.True,
			networks: []string{"testnet", "mainnet", "all", "both", ""},
			message:  "registration form with all network details should be ready to submit",
		},
	}

	for _, tc := range testCases {
		for _, network := range tc.networks {
			tc.assert(t, tc.r.ReadyToSubmit(network), "%s (network %q)", tc.message, network)
		}
	}
}
