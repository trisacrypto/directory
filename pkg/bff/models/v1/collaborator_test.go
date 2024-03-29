package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
)

func TestCollaboratorKeys(t *testing.T) {
	// Should return an empty key when there is no email address
	collab := &models.Collaborator{}
	require.Empty(t, collab.Key(), "expected empty key when email address is missing")

	// Should return a base64-encoded string that is url-safe and does not contain padding
	// Note: This fixture generates a key with / and + characters using standard encoding
	collab = &models.Collaborator{
		Email: "fsdfsfds@local.dev",
	}
	reserved := []string{"+", "/", "="}
	key := collab.Key()
	require.Len(t, key, 22, "expected key to be 22 characters long")
	for _, r := range reserved {
		require.NotContains(t, key, r, "key contains reserved character")
	}

	// Different email addresses should return different keys
	collab.Email = "alice@example.com"
	require.Len(t, collab.Key(), 22, "expected key to be 22 characters long")
	require.NotEqual(t, collab.Key(), key, "expected different key for different email address")
}

func TestValidateCollaborator(t *testing.T) {
	// Collaborator must have an email address
	collab := &models.Collaborator{}
	require.EqualError(t, collab.Validate(), "collaborator is missing email address", "expected error when email address is missing")

	// Collaborator must not have mismatched ID and email
	collab = &models.Collaborator{
		Email: "alice@example.com",
		Id:    "badhash",
	}
	require.EqualError(t, collab.Validate(), "collaborator has invalid id", "expected error when id is invalid")

	// ID should be populated if the email address is set
	collab.Id = ""
	require.NoError(t, collab.Validate(), "expected no error when email address is set")
	require.NotEmpty(t, collab.Id, "expected id to be populated when email address is set")

	// No error should be returned when all fields are valid
	collab = &models.Collaborator{
		Email: "alice@example.com",
		Id:    "wWD4zGmk8L8rA2J1I1PQYA",
	}
	require.NoError(t, collab.Validate(), "expected no error when all fields are valid")
}
