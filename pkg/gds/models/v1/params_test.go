package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

// Test retrieving params from a certificate request model.
func TestGetCertificateRequestParams(t *testing.T) {
	// If the profile is unknown then an error should be returned
	request := &models.CertificateRequest{}
	params, err := models.GetCertificateRequestParams(request, "invalid")
	require.EqualError(t, err, "unknown profile: invalid", "expected error when profile is unknown")
	require.Nil(t, params, "expected nil params when profile is unknown")
	require.Nil(t, request.Params, "expected request params not be modified when profile is unknown")

	testCases := []struct {
		profile  string
		expected map[string]string
	}{
		{
			sectigo.ProfileCipherTraceEE,
			map[string]string{
				"commonName":     "api.alice.vaspbot.net",
				"dNSName":        "api.alice.vaspbot.net\napi.alice.vaspbot.dev",
				"pkcs12Password": "supersecretsquirrel",
			},
		},
		{
			sectigo.ProfileCipherTraceEndEntityCertificate,
			map[string]string{
				"commonName":          "api.bob.vaspbot.net",
				"dNSName":             "api.bob.vaspbot.net\napi.bob.vaspbot.dev",
				"pkcs12Password":      "supersecreteagle",
				"organizationName":    "Bob VASP, PTE",
				"localityName":        "Scarborough",
				"stateOrProvinceName": "North Yorkshire",
				"countryName":         "UK",
			},
		},
	}

	for _, tc := range testCases {
		// Without common name, dns names or pkcs12 password, the params will be invalid, but it
		// should be populated with default values before validation.
		request = &models.CertificateRequest{}
		_, err := models.GetCertificateRequestParams(request, tc.profile)
		require.ErrorContains(t, err, "missing required parameter", "expected invalid params when missing common name or pkcs12 password")
		require.NotNil(t, request.Params, "expected params to be populated with default values")
		require.Equal(t, sectigo.Defaults, request.Params, "expected defaults to be copied into the params")

		// Add common name, dns names, and pkcs12 password
		models.UpdateCertificateRequestParams(request, sectigo.ParamCommonName, "example.com")
		models.UpdateCertificateRequestParams(request, sectigo.ParamDNSNames, "example.com\nsub.example.com")
		models.UpdateCertificateRequestParams(request, sectigo.ParamPassword, "supersecret")

		// Params  should now be valid with default values
		params, err := models.GetCertificateRequestParams(request, tc.profile)
		require.NoError(t, err, "expected params to be valid with defaults, common name, and pkcs12 password")
		require.NotEmpty(t, params, "expected params to be populated with default values")

		// Check params and request.Params consistency
		for _, key := range sectigo.Profiles[tc.profile] {
			require.Contains(t, params, key)
			require.Contains(t, request.Params, key)
			require.NotEmpty(t, params[key])
			require.Equal(t, params[key], request.Params[key])
		}

		// When the request has parameters already set, they should be returned
		request = &models.CertificateRequest{Params: make(map[string]string)}

		// Ensure the params are copied from expected so that the pointer isn't modified.
		for key, val := range tc.expected {
			request.Params[key] = val
		}

		params, err = models.GetCertificateRequestParams(request, tc.profile)
		require.NoError(t, err, "should be able to get params when set on certificate request")
		require.Equal(t, tc.expected, params, "expected original parameters to not be modified")
	}
}

// Test updating params on a certificate request model.
func TestUpdateCertificateRequestParams(t *testing.T) {
	// Test updating a nil map
	request := &models.CertificateRequest{}
	models.UpdateCertificateRequestParams(request, "foo", "bar")
	require.Equal(t, map[string]string{"foo": "bar"}, request.Params, "expected request params map to be updated")

	// Test updating an existing key
	models.UpdateCertificateRequestParams(request, "foo", "rebar")
	require.Equal(t, map[string]string{"foo": "rebar"}, request.Params, "expected request params map to be updated")

	// Test adding a new key
	models.UpdateCertificateRequestParams(request, "bar", "baz")
	require.Equal(t, map[string]string{"foo": "rebar", "bar": "baz"}, request.Params, "expected request params map to be updated")
}

// Test validating parameters on a certificate request model.
func TestValidateCertificateRequestParams(t *testing.T) {

	// If an unknown profile is specified then an error should be returned
	params := make(map[string]string)
	err := models.ValidateCertificateRequestParams(params, "invalid")
	require.EqualError(t, err, "unknown profile: invalid", "expected error when profile is unknown")

	// Ensure that all required params for each profile are required
	for profile, required := range sectigo.Profiles {
		params := make(map[string]string)
		err = models.ValidateCertificateRequestParams(params, profile)
		require.ErrorContains(t, err, "missing required parameter", "expected error when params is empty for profile %q", profile)

		// Keep adding required params one by one until the last parameter
		for _, key := range required[:len(required)-1] {
			params[key] = "notanemptyvalue"
			err = models.ValidateCertificateRequestParams(params, profile)
			require.ErrorContains(t, err, "missing required parameter", "expected error when params is missing keys for profile %q", profile)
		}

		// Params should be valid when they are completed
		// TODO: if common name validation happens this will need to be updated.
		params[required[len(required)-1]] = "lastnotanemptyvalue"
		err = models.ValidateCertificateRequestParams(params, profile)
		require.NoError(t, err, "complete params should return no error")

		// Adding a non required parameter should cause validation to fail
		params["notvalid"] = "notaparam"
		err = models.ValidateCertificateRequestParams(params, profile)
		require.ErrorContains(t, err, "extra profile parameter: notvalid", "expected invalid params when extra keys are included")
	}
}
