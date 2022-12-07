package auth

import (
	"errors"

	"github.com/auth0/go-auth0/management"
)

// TODO: Should these be configurable?
const (
	// BFF Organization management
	ReadOrganizations   = "read:organizations"
	CreateOrganizations = "create:organizations"
	SwitchOrganizations = "switch:organizations"

	// Collaborators management
	ReadCollaborators   = "read:collaborators"
	UpdateCollaborators = "update:collaborators"

	// GDS Registration management
	ReadVASP   = "read:vasp"
	UpdateVASP = "update:vasp"

	// Posting announcements
	CreateAnnouncements = "create:announcements"

	// User roles
	LeaderRole       = "Organization Leader"
	CollaboratorRole = "Organization Collaborator"
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
