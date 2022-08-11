package models

import (
	"fmt"

	"github.com/trisacrypto/directory/pkg/sectigo"
)

// GetCertificateRequestParams returns a subset of the params map on the certificate
// request containing only the required parameters for the given profile.
func GetCertificateRequestParams(request *CertificateRequest, profile string) (params map[string]string, err error) {
	if request.Params == nil {
		request.Params = make(map[string]string)
	}

	var required map[string]struct{}
	if required, err = GetProfileParams(profile); err != nil {
		return nil, err
	}

	params = make(map[string]string)
	for key := range required {
		params[key] = request.Params[key]
	}
	return params, nil
}

// UpdateCertificateRequestParams updates the params map on the certificate request,
// adding the key value pair and overwriting the key if it already exists.
func UpdateCertificateRequestParams(request *CertificateRequest, key, val string) {
	if request.Params == nil {
		request.Params = make(map[string]string)
	}
	request.Params[key] = val
}

// ValidateCertificateRequestParams updates the certificate request params with default
// values if they are missing and checks to make sure that all required values are
// present. If a required value is missing or a default parameter cannot be computed
// then an error is returned. This should only be called right before submitting the
// certificate request to the CA because it modifies the params map on the request
// record and could create a condition where a change is not captured in the request.
func ValidateCertificateRequestParams(request *CertificateRequest, profile string) (err error) {
	if request.Params == nil {
		request.Params = make(map[string]string)
	}

	// Fill in default values for any missing parameters
	for key, val := range sectigo.Defaults {
		if request.Params[key] == "" {
			request.Params[key] = val
		}
	}

	// Check that all required parameters are present
	var required map[string]struct{}
	if required, err = GetProfileParams(profile); err != nil {
		return err
	}
	for key := range required {
		if request.Params[key] == "" {
			return fmt.Errorf("missing required parameter: %s", key)
		}
	}
	return nil
}

// GetProfileParams returns a new params map containing only the required parameters
// for the given profile.
func GetProfileParams(profile string) (params map[string]struct{}, err error) {
	var (
		required []string
		ok       bool
	)

	if required, ok = sectigo.Profiles[profile]; !ok {
		return nil, fmt.Errorf("unknown profile: %s", profile)
	}

	params = make(map[string]struct{})
	for _, key := range required {
		params[key] = struct{}{}
	}
	return params, nil
}
