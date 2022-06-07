package auth0_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/auth0"
)

func TestCredentials(t *testing.T) {
	// Create a set of credentials from JSON example returned from auth0 docs
	// See: https://auth0.com/docs/secure/tokens/access-tokens/get-management-api-access-tokens-for-production
	creds := &auth0.Credentials{}
	err := creds.LoadFrom("testdata/example_token.json")
	require.NoError(t, err, "could not load credentials")

	// Ensure that created at is not zero when loaded
	require.NotEmpty(t, creds.CreatedAt, "created at should be now when loaded from auth0")
	require.Empty(t, creds.ExpiresAt, "expires at should not be set when loaded")

	// Ensure that the credentials are valid (e.g. they should expire 24 hours from now)
	require.True(t, creds.Valid(), "creds not valid when loaded from auth0")
	require.NotEmpty(t, creds.ExpiresAt, "valid creds cache expires at on the credentials")

	// Test JSON serialization and deserialization
	path := filepath.Join(t.TempDir(), "credentials.json")
	err = creds.DumpTo(path)
	require.NoError(t, err, "could not dump credentials to temp dir")

	compat := &auth0.Credentials{}
	err = compat.LoadFrom(path)
	require.NoError(t, err, "cannot load dumped credentials from temp dir")

	require.Equal(t, creds.AccessToken, compat.AccessToken, "expected serialized credentials to match original")
	require.Equal(t, creds.ExpiresIn, compat.ExpiresIn, "expected serialized credentials to match original")
	require.Equal(t, creds.Scope, compat.Scope, "expected serialized credentials to match original")
	require.Equal(t, creds.TokenType, compat.TokenType, "expected serialized credentials to match original")
	require.True(t, creds.CreatedAt.Equal(compat.CreatedAt), "expected serialized credentials to match original")
	require.True(t, creds.ExpiresAt.Equal(compat.ExpiresAt), "expected serialized credentials to match original")
}

func TestInvalidCredentials(t *testing.T) {
	creds := &auth0.Credentials{}
	require.False(t, creds.Valid(), "creds should not be valid without an access token")

	creds.AccessToken = "foo"
	require.False(t, creds.Valid(), "creds should not be valid without a created at or expires in")
	require.Zero(t, creds.ExpiresAt, "expired at should remain zero-valued when invalid")

	creds.ExpiresIn = 86400
	require.False(t, creds.Valid(), "creds should not be valid without a created at")
	require.Zero(t, creds.ExpiresAt, "expired at should remain zero-valued when invalid")

	creds.CreatedAt, _ = time.Parse(time.RFC3339, "1985-03-21T14:31:21Z")
	require.False(t, creds.Valid(), "creds should not be valid if they are expired")
	require.NotZero(t, creds.ExpiresAt, "once created and and expires in are set, expired at should be cached")

	creds.ExpiresAt, _ = time.Parse(time.RFC3339, "1999-02-14T19:12:48Z")
	require.False(t, creds.Valid(), "creds should not be valid if expired at is set directly")
}

func TestCredentialsCache(t *testing.T) {
	// Load Cache should not fail when the cache doesn't exist
	creds := &auth0.Credentials{}
	err := creds.LoadCache("testdata/path/to/nowhere.json")
	require.NoError(t, err, "expected no error when loading non-existing cache")
	require.Empty(t, creds, "expected creds to be zero-valued after cache miss")

	// Load Cache should ignore empty string
	require.NoError(t, creds.LoadCache(""), "cache should ignore empty string path")

	// Should be able to load from actual cache file
	err = creds.LoadCache("testdata/example_token.json")
	require.NoError(t, err, "could not load cache")
	require.NotEmpty(t, creds, "expected credentials to be loaded")
	require.True(t, creds.Valid(), "expected credentials to be valid")

	// Should not error when dumping to a cache file that does not exist
	creds.DumpCache("testdata/path/to/nowhere.json")
	require.NoFileExists(t, "testdata/path/to/nowhere.json", "unexpected cache file")

	// Dump Cache should ignore empty string
	creds.DumpCache("")
}
