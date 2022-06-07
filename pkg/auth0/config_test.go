package auth0_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/auth0"
)

func TestConfigValidate(t *testing.T) {
	conf := auth0.Config{}
	require.Error(t, conf.Validate(), "zero-valued conf should be invalid")

	// Build a valid configuration from scratch
	require.EqualError(t, conf.Validate(), "invalid configuration: missing auth0 domain", "config without a domain should be invalid")
	conf.Domain = "example.auth0.com"

	// Testing should only require the domain
	conf.Testing = true
	require.NoError(t, conf.Validate(), "testing configuration with domain should be valid")

	// Without testing, client ID and secret are required
	conf.Testing = false
	require.EqualError(t, conf.Validate(), "invalid configuration: missing auth0 client id", "non-testing config without a client id should be invalid")
	conf.ClientID = "exampleid"
	require.EqualError(t, conf.Validate(), "invalid configuration: missing auth0 client secret", "non-testing config without a client secret should be invalid")
	conf.ClientSecret = "examplesecret"

	// Configuration should now be valid
	require.NoError(t, conf.Validate(), "complete configuration should be valid")
}

func TestZeroConfig(t *testing.T) {
	require.True(t, auth0.Config{}.IsZero(), "blank config should be zero-valued")

	require.False(t, auth0.Config{Domain: "foo"}.IsZero(), "config with domain should not be zero-valued")
	require.False(t, auth0.Config{ClientID: "foo"}.IsZero(), "config with client id should not be zero-valued")
	require.False(t, auth0.Config{ClientSecret: "foo"}.IsZero(), "config with client secret should not be zero-valued")
	require.False(t, auth0.Config{Testing: true}.IsZero(), "config marked as testing should not be zero-valued")
}

func TestConfigBaseURL(t *testing.T) {
	conf := auth0.Config{Domain: "example.auth0.com", Testing: true}
	u := conf.BaseURL()
	require.Equal(t, "http://example.auth0.com", u.String(), "base url should be http in testing mode")

	conf.Testing = false
	u = conf.BaseURL()
	require.Equal(t, "https://example.auth0.com", u.String(), "base url should be https in non-testing mode")
}
