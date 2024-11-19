package ensign_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils/ensign"
)

func TestValidate(t *testing.T) {
	config := ensign.Config{
		ClientSecret: "client-secret",
		Endpoint:     "ensign.rotational.app:443",
		AuthURL:      "https://auth.rotational.app",
	}

	// Should error if client id is missing.
	require.ErrorIs(t, config.IsValid(), ensign.ErrMissingClientID, "expected missing client id error")

	// Should error if client secret is missing.
	config.ClientID = "client-id"
	config.ClientSecret = ""
	require.ErrorIs(t, config.IsValid(), ensign.ErrMissingClientSecret, "expected missing client secret error")

	// Should error if endpoint is missing.
	config.ClientSecret = "client-secret"
	config.Endpoint = ""
	require.ErrorIs(t, config.IsValid(), ensign.ErrMissingEndpoint, "expected missing endpoint error")

	// Should error if auth url is missing.
	config.Endpoint = "ensign.rotational.app:443"
	config.AuthURL = ""
	require.ErrorIs(t, config.IsValid(), ensign.ErrMissingAuthURL, "expected missing auth url error")

	// Should not error if all required fields are present.
	config.AuthURL = "https://auth.rotational.app"
	require.NoError(t, config.IsValid(), "expected no error for valid configuration")
}
