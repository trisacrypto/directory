package api

import (
	"encoding/json"
	"os"
	"time"
)

// Credentials provides a basic interface for loading an access token from Auth0 into
// the BFF API Client. Credentials can be loaded from disk, generated, or directly
// authenticated with Auth0 via a machine-to-machine token.
//
// NOTE: this is a fundamentally different mechanism than the GDS Admin API credentials
// because the Admin API generates its own tokens, and Auth0 manages the BFF credentials.
type Credentials interface {
	AccessToken() (string, error)
}

// Check to ensure that the different types of Credentials implement the interface.
var (
	_ Credentials = Token("")
	_ Credentials = &LocalCredentials{}
	_ Credentials = &Auth0Token{}
)

// A Token is just the JWT base64 encoded token string that can be obtained from the
// Auth0 debugger or created in memory for tests using Token("mytoken"). Token
// implements the Credentials interface so it can be passed directly to the client.
type Token string

// Token implements the Credentials interface but performs limited validation on the string.
func (t Token) AccessToken() (string, error) {
	if string(t) == "" {
		return "", ErrInvalidCredentials
	}
	return string(t), nil
}

// LocalCredentials loads and saves the access token from disk.
type LocalCredentials struct {
	Path  string
	Token *Auth0Token
}

// Load the credentials from the path on disk.
func (t *LocalCredentials) Load() (err error) {
	if t.Path == "" {
		return ErrPathRequired
	}

	var f *os.File
	if f, err = os.Open(t.Path); err != nil {
		return err
	}
	defer f.Close()

	t.Token = &Auth0Token{}
	if err = json.NewDecoder(f).Decode(t.Token); err != nil {
		return err
	}
	return nil
}

// Dump the credentials to store them to the path on disk.
func (t *LocalCredentials) Dump() (err error) {
	if t.Path == "" {
		return ErrPathRequired
	}

	if t.Token == nil {
		return ErrInvalidCredentials
	}

	var f *os.File
	if f, err = os.OpenFile(t.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(t.Token); err != nil {
		return err
	}
	return nil
}

// AccessToken implements the Credentials interface by checking if the token has been
// loaded, and if not, it loads the token from disk. Validation is performed by the
// Auth0Token to ensure the credentials are valid and not expired.
func (t *LocalCredentials) AccessToken() (_ string, err error) {
	if t.Token == nil {
		if err = t.Load(); err != nil {
			return "", err
		}
	}
	return t.Token.AccessToken()
}

// Auth0Token is a JSON representation of the Token returned by Auth0
type Auth0Token struct {
	Token     string    `json:"access_token"`
	ExpiresIn int64     `json:"expires_in"`
	Scope     string    `json:"scope"`
	Type      string    `json:"token_type"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// AccessToken implements the Credentials interface and ensures that a token is present
// and that the token has not expired yet. It relies on the data structure from Auth0
// rather than parsing the jwt token directly.
func (t *Auth0Token) AccessToken() (_ string, err error) {
	if t.Token == "" {
		return "", ErrInvalidCredentials
	}

	if time.Now().After(t.ExpiresAt) {
		return "", ErrExpiredCredentials
	}

	return t.Token, nil
}
