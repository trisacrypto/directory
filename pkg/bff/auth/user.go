package auth

import (
	"errors"

	"github.com/auth0/go-auth0/management"
)

// UserDisplayName is a helper to get the user's display name from the Auth0 user
// record. This should be used when the backend needs to retrieve a user-facing display
// name for the user and returns an error if no name is available.
func UserDisplayName(user *management.User) (string, error) {
	if user == nil {
		return "", errors.New("user record is nil")
	}

	// Prefer the user's actual name if available
	switch {
	case user.Name != nil && *user.Name != "":
		return *user.Name, nil
	case user.Email != nil && *user.Email != "":
		return *user.Email, nil
	default:
		return "", errors.New("user record has no name or email address")
	}
}
