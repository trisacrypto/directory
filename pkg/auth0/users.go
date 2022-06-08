package auth0

import (
	"context"
	"encoding/json"
	"net/http"
)

// User is a data struct that encapsulates the JSON response of an auth0 v2 user object.
// Note that the identities and multifactor arrays have been purposefully left off.
type User struct {
	ID            string   `json:"user_id,omitempty"`
	Email         string   `json:"email,omitempty"`
	EmailVerified bool     `json:"email_verified,omitempty"`
	Username      string   `json:"username,omitempty"`
	PhoneNumber   string   `json:"phone_number,omitempty"`
	PhoneVerified bool     `json:"phone_verified,omitempty"`
	CreatedAt     string   `json:"created_at,omitempty"`
	UpdatedAt     string   `json:"updated_at,omitempty"`
	AppMetadata   Metadata `json:"app_metadata,omitempty"`
	UserMetadata  Metadata `json:"user_metadata,omitempty"`
	Picture       string   `json:"picture,omitempty"`
	Name          string   `json:"name,omitempty"`
	Nickname      string   `json:"nickname,omitempty"`
	LastIP        string   `json:"last_ip,omitempty"`
	LastLogin     string   `json:"last_login,omitempty"`
	LoginsCount   int64    `json:"logins_count,omitempty"`
	Blocked       bool     `json:"blocked,omitempty"`
	GivenName     string   `json:"given_name,omitempty"`
	FamilyName    string   `json:"family_name,omitempty"`
}

// Metadata is a variable JSON object that can store arbitrary values for the user.
type Metadata map[string]interface{}

// GetUser retrieves user information for the specified user ID.
// Although Auth0 allows including or excluding specific fields, this method does not
// allow the user to specify these queries and instead hard codes the exclusion of some
// fields that are not available in the user struct.
func (a *Auth0) GetUser(ctx context.Context, userID string) (user *User, err error) {
	// Perform pre-flight check on authenticated endpoint
	if err = a.Preflight(); err != nil {
		return nil, err
	}

	// Manually exclude the identies and multifactor fields
	query := map[string]string{
		"fields":         "identities",
		"include_fields": "false",
	}

	// Create the request for the users endpoint
	endpoint := a.Endpoint("/api/v2/users/%s", query, userID)
	var req *http.Request
	if req, err = a.NewRequest(ctx, http.MethodGet, endpoint, nil); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = a.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = json.NewDecoder(rep.Body).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}
