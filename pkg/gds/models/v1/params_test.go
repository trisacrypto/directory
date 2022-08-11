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
	require.Equal(t, map[string]string{}, request.Params, "expected request params map to be initialized")

	// If the profile is known but there are no request params then empty values should be returned
	expected := map[string]string{
		sectigo.ParamCommonName: "",
		sectigo.ParamDNSNames:   "",
		sectigo.ParamPassword:   "",
	}
	params, err = models.GetCertificateRequestParams(request, sectigo.ProfileCipherTraceEE)
	require.NoError(t, err, "GetCertificateRequestParams returned an error")
	require.Equal(t, expected, params, "wrong params returned")
	require.Equal(t, map[string]string{}, request.Params, "expected request params map to be unchanged")

	// If there are different request params then only the required ones should be returned
	request.Params["foo"] = "bar"
	params, err = models.GetCertificateRequestParams(request, sectigo.ProfileCipherTraceEE)
	require.NoError(t, err, "GetCertificateRequestParams returned an error")
	require.Equal(t, expected, params, "wrong params returned")

	// If there are some matching request params then only the required ones should be returned
	request.Params[sectigo.ParamCommonName] = "alice.example.com"
	expected[sectigo.ParamCommonName] = "alice.example.com"
	params, err = models.GetCertificateRequestParams(request, sectigo.ProfileCipherTraceEE)
	require.NoError(t, err, "GetCertificateRequestParams returned an error")
	require.Equal(t, expected, params, "wrong params returned")
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
	request := &models.CertificateRequest{}
	expected := map[string]string{
		sectigo.ParamOrganizationName:    "TRISA Member VASP",
		sectigo.ParamLocalityName:        "Menlo Park",
		sectigo.ParamStateOrProvinceName: "California",
		sectigo.ParamCountryName:         "US",
	}
	err := models.ValidateCertificateRequestParams(request, "invalid")
	require.EqualError(t, err, "unknown profile: invalid", "expected error when profile is unknown")
	require.Equal(t, expected, request.Params, "expected request params map to be initialized")

	// If required params are missing then an error should be returned
	err = models.ValidateCertificateRequestParams(request, sectigo.ProfileCipherTraceEE)
	require.ErrorContains(t, err, "missing required parameter", "expected error when required parameter is missing")
	require.Equal(t, expected, request.Params, "expected request params map to be unchanged")

	// Test defaults are filled in but the required params are still missing
	err = models.ValidateCertificateRequestParams(request, sectigo.ProfileCipherTraceEE)
	require.ErrorContains(t, err, "missing required parameter", "expected error when required parameter is missing")
	require.Equal(t, expected, request.Params, "expected request params map to be populated with defaults")

	// If required params are present then the vaildation should succeed
	request.Params = map[string]string{
		sectigo.ParamCommonName: "alice.example.com",
		sectigo.ParamDNSNames:   "alice.example.com\nalice.us.example.com",
		sectigo.ParamPassword:   "password",
	}
	expected[sectigo.ParamCommonName] = "alice.example.com"
	expected[sectigo.ParamDNSNames] = "alice.example.com\nalice.us.example.com"
	expected[sectigo.ParamPassword] = "password"
	err = models.ValidateCertificateRequestParams(request, sectigo.ProfileCipherTraceEE)
	require.NoError(t, err, "ValidateCertificateRequestParams returned an error")
	require.Equal(t, expected, request.Params, "expected request params map to be unchanged")

	// Test existing values are not overwritten by the defaults
	request.Params = map[string]string{
		sectigo.ParamOrganizationName: "Alice VASP",
		sectigo.ParamCommonName:       "alice.example.com",
		sectigo.ParamDNSNames:         "alice.example.com\nalice.us.example.com",
		sectigo.ParamPassword:         "password",
	}
	expected[sectigo.ParamOrganizationName] = "Alice VASP"
	err = models.ValidateCertificateRequestParams(request, sectigo.ProfileCipherTraceEndEntityCertificate)
	require.NoError(t, err, "ValidateCertificateRequestParams returned an error")
	require.Equal(t, expected, request.Params, "expected request params map to be populated with defaults")
}
