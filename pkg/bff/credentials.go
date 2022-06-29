package bff

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

// Credentials implements the admin.Credentials interface to provide access and refresh
// tokens to the GDSClient.
type Credentials struct {
	tm      *tokens.TokenManager
	access  string
	refresh string
}

// NewCredentials creates a new Credentials object which can generate new access and
// refresh tokens for authentication based on the provided token keys and audience.
func NewCredentials(conf config.AdminConfig) (creds *Credentials, err error) {
	creds = &Credentials{}

	aud := conf.Audience
	if aud == "" {
		aud = conf.Endpoint
	}

	if creds.tm, err = tokens.New(conf.TokenKeys, conf.Audience); err != nil {
		return creds, err
	}

	return creds, nil
}

// Load signs new access and refresh tokens using the token manager.
func (c *Credentials) Load(api admin.DirectoryAdministrationClient) (err error) {
	var accessToken, refreshToken *jwt.Token

	claims := map[string]interface{}{
		"hd":      "rotational.io",
		"email":   "bff@rotational.io",
		"name":    "GDS BFF",
		"picture": "",
	}

	// Create the access and refresh tokens from the claims
	if accessToken, err = c.tm.CreateAccessToken(claims); err != nil {
		return err
	}

	if refreshToken, err = c.tm.CreateRefreshToken(accessToken); err != nil {
		return err
	}

	// Sign the tokens
	if c.access, err = c.tm.Sign(accessToken); err != nil {
		return err
	}

	if c.refresh, err = c.tm.Sign(refreshToken); err != nil {
		return err
	}

	// Call ProtectAuthenticate to get the csrf tokens
	apiv2, ok := api.(*admin.APIv2)
	if ok {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = apiv2.ProtectAuthenticate(ctx); err != nil {
			return fmt.Errorf("could not get double cookie csrf tokens: %s", err)
		}
	}

	return nil
}

// Validate that the access and refresh tokens have not expired.
func (c *Credentials) Validate() (err error) {
	if c.access == "" || c.refresh == "" {
		return errors.New("credentials incomplete")
	}

	accessClaims := new(tokens.Claims)
	if token, _ := jwt.ParseWithClaims(c.access, accessClaims, nil); token == nil {
		return errors.New("could not parse access token")
	}

	now := time.Now()
	if !accessClaims.ExpiresAt.IsZero() && accessClaims.ExpiresAt.Before(now) {
		// access token is expired, check if refresh is not expired
		refreshClaims := new(tokens.Claims)
		if token, _ := jwt.ParseWithClaims(c.refresh, refreshClaims, nil); token == nil {
			return errors.New("could not parse refresh token")
		}

		if !refreshClaims.ExpiresAt.IsZero() && refreshClaims.ExpiresAt.Before(now) {
			// refresh token is also expired
			return errors.New("tokens are both expired")
		}
	}

	// either access or refresh token is unexpired.
	return nil
}

// Login creates new access and refresh tokens if they don't already exist.
func (c *Credentials) Login(api admin.DirectoryAdministrationClient) (accessToken, refreshToken string, err error) {
	if c.access == "" || c.refresh == "" {
		if err = c.Load(api); err != nil {
			return "", "", err
		}
	}

	return c.access, c.refresh, nil
}

// Refresh reauthenticates with the server to get new access and refresh tokens if they
// have not expired.
func (c *Credentials) Refresh(api admin.DirectoryAdministrationClient) (accessToken, refreshToken string, err error) {
	// Validate that the tokens are not expired
	if err = c.Validate(); err != nil {
		return "", "", err
	}

	// Attempt to reauthenticate
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rep *admin.AuthReply
	if rep, err = api.Reauthenticate(ctx, &admin.AuthRequest{Credential: c.refresh}); err != nil {
		return "", "", err
	}

	c.access = rep.AccessToken
	c.refresh = rep.RefreshToken

	return c.access, c.refresh, nil
}

// Logout deletes the access and refresh tokens.
func (c *Credentials) Logout(api admin.DirectoryAdministrationClient) (err error) {
	c.access = ""
	c.refresh = ""
	return nil
}
