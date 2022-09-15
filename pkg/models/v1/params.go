package models

import (
	"fmt"

	"github.com/trisacrypto/directory/pkg/sectigo"
	"golang.org/x/exp/slices"
)

// GetCertificateRequestParams returns the params that must be submitted to Sectigo to
// create a certificate for the given profile. This method filters any extraneous
// parameters that cannot be submitted to Sectigo that may be on the CertificateRequest
// and ensures that the parameters are valid before returning the params. If the
// CertificateRequest params are missing required parameters then defaults are populated
// if they are available.

// NOTE: This method may update the certificate request but it is up to the caller to
// persist the changes back to the database.
func GetCertificateRequestParams(request *CertificateRequest, profile string) (params map[string]string, err error) {
	// Get the required parameters for the Sectigo profile
	// NOTE: this needs to happen before the request params are modified
	required, ok := sectigo.Profiles[profile]
	if !ok {
		return nil, fmt.Errorf("unknown profile: %s", profile)
	}

	if request.Params == nil {
		request.Params = make(map[string]string)
	}

	// Fill in default values for any missing parameters
	for key, val := range sectigo.Defaults {
		if request.Params[key] == "" {
			request.Params[key] = val
		}
	}

	// Ensure only the required parameters are returned since the params may contain
	// values for multiple profiles.
	params = make(map[string]string)
	for _, key := range required {
		params[key] = request.Params[key]
	}

	// Validate the parameters to ensure they can be submitted to Sectigo
	if err = ValidateCertificateRequestParams(params, profile); err != nil {
		return nil, err
	}

	return params, nil
}

// UpdateCertificateRequestParams updates the params map on the certificate request,
// adding the key value pair and overwriting the key if it already exists.
// NOTE: it is up to the caller to persist the request back to the database.
func UpdateCertificateRequestParams(request *CertificateRequest, key, val string) {
	if request.Params == nil {
		request.Params = make(map[string]string)
	}
	request.Params[key] = val
}

// ValidateCertificateRequestParams ensures that the parameters about to be sent to
// Sectigo are valid; e.g. that they contain only the required parameters for the
// profile and that all parameter values are populated.
func ValidateCertificateRequestParams(params map[string]string, profile string) (err error) {
	// Get the required parameters for the Sectigo profile
	required, ok := sectigo.Profiles[profile]
	if !ok {
		return fmt.Errorf("unknown profile: %s", profile)
	}

	// Check that all required parameters are in the params and that the params are
	// populated with a non-empty value.
	for _, key := range required {
		if val, ok := params[key]; !ok || val == "" {
			return fmt.Errorf("missing required parameter: %s", key)
		}
	}

	// Check that there are no parameters that are not in the required array
	// NOTE: expects that the required strings are sorted!
	for key := range params {
		if _, found := slices.BinarySearch(required, key); !found {
			return fmt.Errorf("params include extra profile parameter: %s", key)
		}
	}

	// TODO: Check to ensure the common name is in the dNSNAmes
	return nil
}
