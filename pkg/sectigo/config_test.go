package sectigo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

const ValidProfile = sectigo.ProfileCipherTraceEE

func TestConfigValidation(t *testing.T) {
	validTestCases := []sectigo.Config{
		{Profile: sectigo.ProfileCipherTraceEE, Environment: "production"},
		{Profile: sectigo.ProfileIDCipherTraceEE, Environment: "staging"},
		{Profile: sectigo.ProfileCipherTraceEndEntityCertificate, Environment: "testing"},
		{Profile: sectigo.ProfileIDCipherTraceEndEntityCertificate, Environment: "staging"},
		{Profile: ValidProfile, Environment: "production", Endpoint: "https://iot.sectigo.com"},
		{Profile: ValidProfile, Environment: "staging", Endpoint: "http://cathy:8831"},
		{Profile: ValidProfile, Environment: "staging", Endpoint: "https://cathy.test-net.io"},
		{Profile: ValidProfile, Environment: "testing", Endpoint: "https://localhost:4432"},
		{Profile: ValidProfile, Environment: "testing", Endpoint: "https://127.0.0.1:4432"},
	}

	for i, conf := range validTestCases {
		require.NoError(t, conf.Validate(), "expected test case %d to be valid", i)
	}

	invalidTestCases := []struct {
		conf sectigo.Config
		errs string
	}{
		{sectigo.Config{Profile: "invalid profile"}, "is not a valid Sectigo profile name"},
		{sectigo.Config{Profile: ValidProfile, Environment: "foo"}, "is not a valid environment"},
		{sectigo.Config{Profile: ValidProfile, Environment: "production", Endpoint: "http://iot.sectigo.com"}, "must use https in production"},
		{sectigo.Config{Profile: ValidProfile, Environment: "production", Endpoint: "ftps://iot.sectigo.com"}, "must use https in production"},
		{sectigo.Config{Profile: ValidProfile, Environment: "production", Endpoint: "https://cathy.test-net.io"}, "cannot connect to cathy.test-net.io in production"},
		{sectigo.Config{Profile: ValidProfile, Environment: "staging", Endpoint: "http://iot.sectigo.com"}, "cannot connect to iot.sectigo.com in staging"},
		{sectigo.Config{Profile: ValidProfile, Environment: "testing", Endpoint: "http://iot.sectigo.com"}, "sectigo hostname must be set to localhost in testing mode"},
	}

	for i, tc := range invalidTestCases {
		err := tc.conf.Validate()
		require.ErrorContains(t, err, tc.errs, "expected test case %d to be invalid", i)
	}
}
