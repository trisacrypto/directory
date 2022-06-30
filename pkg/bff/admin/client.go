package admin

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

const (
	HD    = "hd"
	EMAIL = "bff@rotational.io"
	NAME  = "GDS BFF"
)

// New creates a new DirectoryAdministrationClient which uses its own self signed
// credentials to make authenticated requests to a GDS Admin service.
func New(conf config.AdminConfig) (_ admin.DirectoryAdministrationClient, err error) {
	aud := conf.Audience
	if aud == "" {
		aud = conf.Endpoint
	}

	var creds *Credentials
	if creds, err = NewCredentials(conf.TokenKeys, conf.Audience); err != nil {
		return nil, err
	}

	return admin.New(conf.Endpoint, creds)
}

// NewCredentials creates a new Credentials object with the given token keys and
// audience which can generate self signed access tokens for authenticated requests.
func NewCredentials(tokenKeys map[string]string, audience string) (creds *Credentials, err error) {
	creds = &Credentials{}

	if creds.tm, err = tokens.New(tokenKeys, audience); err != nil {
		return nil, err
	}

	return creds, nil
}

// Credentials implements the admin.Credentials interface to provide access tokens to
// authenticated requests.
type Credentials struct {
	tm     *tokens.TokenManager
	access string
}

// Generate signs a new access token using the token manager if a valid one does not
// already exist.
func (c *Credentials) Generate() (err error) {
	if c.Valid() {
		return nil
	}

	claims := map[string]interface{}{
		"hd":      HD,
		"email":   EMAIL,
		"name":    NAME,
		"picture": "",
	}

	// Create the access token from the claims
	var token *jwt.Token
	if token, err = c.tm.CreateAccessToken(claims); err != nil {
		return err
	}

	// Sign the token
	if c.access, err = c.tm.Sign(token); err != nil {
		return err
	}

	// Call ProtectAuthenticate to get the csrf tokens
	/*
		apiv2, ok := api.(*admin.APIv2)
		if ok {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err = apiv2.ProtectAuthenticate(ctx); err != nil {
				return fmt.Errorf("could not get double cookie csrf tokens: %s", err)
			}
		}*/

	return nil
}

// Check if the access token is valid and not expired.
func (c *Credentials) Valid() bool {
	if c.access == "" {
		return false
	}

	accessClaims := new(tokens.Claims)
	if token, _ := jwt.ParseWithClaims(c.access, accessClaims, nil); token == nil {
		log.Error().Msg("could not parse access token")
		return false
	}
	return !accessClaims.ExpiresAt.IsZero() && accessClaims.ExpiresAt.After(time.Now())
}

// Login creates new access and refresh tokens if they don't already exist.
func (c *Credentials) Login(api admin.DirectoryAdministrationClient) (accessToken, _ string, err error) {
	if err = c.Generate(); err != nil {
		return "", "", err
	}

	return c.access, "", nil
}

// Refresh
func (c *Credentials) Refresh(api admin.DirectoryAdministrationClient) (accessToken, _ string, err error) {
	if err = c.Generate(); err != nil {
		return "", "", err
	}

	return c.access, "", nil
}

// Logout deletes the access token.
func (c *Credentials) Logout(api admin.DirectoryAdministrationClient) (err error) {
	c.access = ""
	return nil
}
