package sectigo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredsCopy(t *testing.T) {
	api, err := New("foo", "supersecret", "CipherTrace EE")
	require.NoError(t, err)

	// Ensure that creds are copied and are not the same object
	creds := api.Creds()
	require.NotEqual(t, &api.creds, &creds)

	require.Equal(t, api.creds.Username, creds.Username)
	creds.Username = "superbunny"
	require.NotEqual(t, api.creds.Username, creds.Username)
	require.Equal(t, api.creds.Username, "foo")
}

func TestAuthenticate(t *testing.T) {
	api, err := New("foo", "supersecret", "CipherTrace EE")
	require.NoError(t, err)

	s, err := NewMockServer(nil)
	require.NoError(t, err)
	defer s.server.Close()

	err = api.Authenticate()
	require.NoError(t, err)
}
