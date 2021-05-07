package sectigo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredsCopy(t *testing.T) {
	api, err := New("foo", "supersecret")
	require.NoError(t, err)

	// Ensure that creds are copied and are not the same object
	creds := api.Creds()
	require.NotEqual(t, &api.creds, &creds)

	require.Equal(t, api.creds.Username, creds.Username)
	creds.Username = "superbunny"
	require.NotEqual(t, api.creds.Username, creds.Username)
	require.Equal(t, api.creds.Username, "foo")
}
