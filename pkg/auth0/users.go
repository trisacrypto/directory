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
	Password      string   `json:"password,omitempty"`
	Connection    string   `json:"connection,omitempty"`
	ClientID      string   `json:"client_id,omitempty"`
}

// Metadata is a variable JSON object that can store arbitrary values for the user.
type Metadata map[string]interface{}

// GetUser retrieves user information for the specified user ID.
// Although Auth0 allows including or excluding specific fields, this method does not
// allow the user to specify these queries and instead hard codes the exclusion of some
// fields that are not available in the user struct.
// See: https://auth0.com/docs/api/management/v2#!/Users/get_users_by_id
func (a *Auth0) GetUser(ctx context.Context, userID string) (user *User, err error) {
	// Perform pre-flight check on authenticated endpoint
	if err = a.Preflight(ctx); err != nil {
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

// UpdateUser via a patch method, which means only the information you want to change
// should be specified. Note that there are many rules about what data is required for
// updating. For example, to update email, username, or password you have to also supply
// the connection property. In general this method should be mostly used to update the
// app_metadata and user_metadata which uses a merge method rather than replacement.
// NOTE: the update cannot have ID set on it, it must be supplied as the userID argument.
// See: https://auth0.com/docs/api/management/v2#!/Users/patch_users_by_id
func (a *Auth0) UpdateUser(ctx context.Context, userID string, update *User) (updated *User, err error) {
	// A user ID is required to make the update
	if userID == "" {
		return nil, &APIError{StatusCode: 400, Status: "Invalid Request", Message: "A user ID is required to make an update request"}
	}

	if update.ID != "" {
		return nil, &APIError{StatusCode: 400, Status: "Invalid Update", Message: "A user ID cannot be specified on the updates sent to the server"}
	}

	// Perform pre-flight check on authenticated endpoint
	if err = a.Preflight(ctx); err != nil {
		return nil, err
	}

	// Create the request for the users update endpoint
	endpoint := a.Endpoint("/api/v2/users/%s", nil, userID)
	var req *http.Request
	if req, err = a.NewRequest(ctx, http.MethodPatch, endpoint, update); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = a.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	if err = json.NewDecoder(rep.Body).Decode(&updated); err != nil {
		return nil, err
	}
	return updated, nil
}
