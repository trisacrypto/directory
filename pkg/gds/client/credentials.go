package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

// Login implements admin.Credentials so that the AdminProfile can provide access and refresh tokens to the client.
func (p *AdminProfile) Login(api admin.DirectoryAdministrationClient) (accessToken, refreshToken string, err error) {
	var cache *CredentialCache
	if cache, err = LoadCredentials(); err != nil {
		return "", "", fmt.Errorf("could not load credentials: %s", err)
	}

	var creds *Credentials
	if creds, err = cache.Get(p.Endpoint); err != nil {
		// Generate new access and refresh tokens using the token keys method
		// TODO: also allow CLI oauth2 login workflow
		if creds, err = p.GenerateTokens(api); err != nil {
			return "", "", err
		}

		// Save the credentials back to disk
		cache.Credentials[p.Endpoint] = creds
		if err = StoreCredentials(cache); err != nil {
			return "", "", fmt.Errorf("could not store credentials: %s", err)
		}
	}
	return creds.AccessToken, creds.RefreshToken, nil
}

// Refresh implements admin.Credentials so that the AdminProfile can reauthenticate
// access and refresh tokens and provide them to the client.
func (p *AdminProfile) Refresh(api admin.DirectoryAdministrationClient) (accessToken, refreshToken string, err error) {
	var cache *CredentialCache
	if cache, err = LoadCredentials(); err != nil {
		return "", "", fmt.Errorf("could not load credentials: %s", err)
	}

	var creds *Credentials
	if creds, err = cache.Get(p.Endpoint); err != nil {
		// Attempt to read the credentials off the api client
		creds = &Credentials{}
		if apiv2, ok := api.(*admin.APIv2); ok {
			creds.AccessToken, creds.RefreshToken = apiv2.Tokens()
		}
	}

	// Check if we can reauthenticate
	if err = creds.Validate(); err != nil {
		return "", "", fmt.Errorf("could not reauthenticate: %s", err)
	}

	// Attempt to reauthenticate
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rep *admin.AuthReply
	if rep, err = api.Reauthenticate(ctx, &admin.AuthRequest{Credential: creds.RefreshToken}); err != nil {
		return "", "", err
	}

	// Store the cached reply back to disk
	cache.Credentials[p.Endpoint] = &Credentials{AccessToken: rep.AccessToken, RefreshToken: rep.RefreshToken}
	if err = StoreCredentials(cache); err != nil {
		return "", "", fmt.Errorf("could not store credentials: %s", err)
	}

	return rep.AccessToken, rep.RefreshToken, nil
}

// Logout implements admin.Credentials so that the AdminProfile can remove access and
// refresh tokens on request from the client.
func (p *AdminProfile) Logout(api admin.DirectoryAdministrationClient) (err error) {
	// Load the cached credentials
	var cache *CredentialCache
	if cache, err = LoadCredentials(); err != nil {
		return fmt.Errorf("could not load credentials: %s", err)
	}

	// Delete the credentials from the cache
	delete(cache.Credentials, p.Endpoint)

	// Store the cached credentials back to disk
	if err = StoreCredentials(cache); err != nil {
		return fmt.Errorf("could not store credentials: %s", err)
	}
	return nil
}

// GenerateTokens creates a token manager to generate and save credentials
func (p *AdminProfile) GenerateTokens(api admin.DirectoryAdministrationClient) (creds *Credentials, err error) {
	if len(p.TokenKeys) == 0 {
		return nil, errors.New("invalid configuration: token keys are required for local key generation")
	}

	aud := p.Audience
	if aud == "" {
		aud = p.Endpoint
	}

	var tm *tokens.TokenManager
	if tm, err = tokens.New(p.TokenKeys, aud); err != nil {
		return nil, err
	}

	var accessToken, refreshToken *jwt.Token

	claims := map[string]interface{}{
		"hd":      "rotational.io",
		"email":   "admin@rotational.io",
		"name":    "GDS Admin CLI",
		"picture": "",
	}

	// Create the access and refresh tokens from the claims
	if accessToken, err = tm.CreateAccessToken(claims); err != nil {
		return nil, err
	}

	if refreshToken, err = tm.CreateRefreshToken(accessToken); err != nil {
		return nil, err
	}

	// Sign the tokens and return the response
	creds = new(Credentials)
	if creds.AccessToken, err = tm.Sign(accessToken); err != nil {
		return nil, err
	}
	if creds.RefreshToken, err = tm.Sign(refreshToken); err != nil {
		return nil, err
	}

	// Make sure ProtectAuthenticate is called so that we have csrf tokens
	apiv2, ok := api.(*admin.APIv2)
	if ok {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = apiv2.ProtectAuthenticate(ctx); err != nil {
			return nil, fmt.Errorf("could not get double cookie csrf tokens: %s", err)
		}
	}

	return creds, nil
}

const (
	credentialsJSON    = "admin_credentials.json"
	credentialsVersion = "v1"
)

type CredentialCache struct {
	Version     string                  `json:"version"`
	Credentials map[string]*Credentials `json:"credentials"`
}

type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Credentials returns the cached access and refresh tokens from disk.
func LoadCredentials() (cache *CredentialCache, err error) {
	folder := cfgd.QueryFolderContainsFile(credentialsJSON)
	if folder != nil {
		var data []byte
		if data, err = folder.ReadFile(credentialsJSON); err != nil {
			return nil, err
		}

		cache = &CredentialCache{}
		if err = json.Unmarshal(data, cache); err != nil {
			return nil, err
		}

		if cache.Version != credentialsVersion {
			return nil, fmt.Errorf("credentials file version %s mismatch with current version %s", cache.Version, credentialsVersion)
		}

		return cache, nil
	}

	// Return empty credentials if it cannot be loaded
	return &CredentialCache{Credentials: make(map[string]*Credentials)}, nil
}

// StoreCredentials saves and updates the cached access and refresh tokens back to disk.
func StoreCredentials(cache *CredentialCache) (err error) {
	folder := cfgd.QueryFolderContainsFile(credentialsJSON)
	if folder == nil {
		// Store the credentials in the same folder as the profiles
		if folder, err = GetProfilesFolder(); err != nil {
			return err
		}
	}

	// Ensure the cache is at the correct version
	cache.Version = credentialsVersion

	var data []byte
	if data, err = json.MarshalIndent(cache, "", "  "); err != nil {
		return err
	}

	return folder.WriteFile(credentialsJSON, data)
}

// Get credentials for the specified endpoint. Only returns valid credentials.
func (c *CredentialCache) Get(endpoint string) (creds *Credentials, err error) {
	var ok bool
	if creds, ok = c.Credentials[endpoint]; !ok {
		return nil, errors.New("credentials not found")
	}

	if err = creds.Validate(); err != nil {
		return nil, err
	}

	return creds, nil
}

// Validate that the have not expired, e.g. that the access token is not expired, or if
// it is that the refresh token has not expired. Does not check signatures or any other
// claims for validity.
func (c *Credentials) Validate() (err error) {
	if c.AccessToken == "" || c.RefreshToken == "" {
		return errors.New("credentials incomplete")
	}

	var accessClaims *tokens.Claims
	if token, _ := jwt.ParseWithClaims(c.AccessToken, accessClaims, nil); token == nil {
		return errors.New("could not parse access token")
	}

	now := time.Now().Unix()
	if accessClaims.ExpiresAt != 0 && now > accessClaims.ExpiresAt {
		// access token is expired, check if refresh is not expired
		var refreshClaims *tokens.Claims
		if token, _ := jwt.ParseWithClaims(c.RefreshToken, refreshClaims, nil); token == nil {
			return errors.New("could not parse refresh token")
		}

		if refreshClaims.ExpiresAt != 0 && now > refreshClaims.ExpiresAt {
			// refresh token is also expired
			return errors.New("tokens are both expired")
		}
	}

	// either access or refresh token is unexpired.
	return nil
}
