package secrets_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
)

// Test that the CreateToken function creates a valid token of the given length.
func TestCreateToken(t *testing.T) {
	// Negative or zero length should return an empty string
	require.Equal(t, "", secrets.CreateToken(-1))
	require.Equal(t, "", secrets.CreateToken(0))

	// Returned token should not contain any unexpected characters
	token := secrets.CreateToken(20)
	require.True(t, secrets.ValidateToken(token), "token %s contains invalid characters", token)

	// Successive calls should return different tokens
	nextToken := secrets.CreateToken(20)
	require.True(t, secrets.ValidateToken(nextToken), "token %s contains invalid characters", nextToken)
	require.NotEqual(t, token, nextToken, "CreateToken returned the same token twice: %s", token)
}
