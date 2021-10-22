package admin_test

import (
	"crypto/rand"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

type MockCredentials struct {
	Calls       map[string]int
	tokenSecret []byte
}

func (c *MockCredentials) Login(api DirectoryAdministrationClient) (accessToken, refreshToken string, err error) {
	now := time.Now()
	if accessToken, err = c.make(now.Add(1 * time.Second)); err != nil {
		return "", "", err
	}

	if refreshToken, err = c.make(now.Add(2 * time.Second)); err != nil {
		return "", "", err
	}

	c.incr("Login")
	return accessToken, refreshToken, nil
}

func (c *MockCredentials) Refresh(api DirectoryAdministrationClient) (accessToken, refreshToken string, err error) {
	now := time.Now()
	if accessToken, err = c.make(now.Add(1 * time.Second)); err != nil {
		return "", "", err
	}

	if refreshToken, err = c.make(now.Add(2 * time.Second)); err != nil {
		return "", "", err
	}

	c.incr("Refresh")
	return accessToken, refreshToken, nil
}

func (c *MockCredentials) Logout(api DirectoryAdministrationClient) (err error) {

	c.incr("Logout")
	return nil
}

func (c *MockCredentials) make(expires time.Time) (tks string, err error) {
	if len(c.tokenSecret) == 0 {
		c.tokenSecret = make([]byte, 32)
		if _, err := rand.Read(c.tokenSecret); err != nil {
			return "", err
		}
	}

	claims := &tokens.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.tokenSecret)
}

func (c *MockCredentials) incr(method string) {
	if c.Calls == nil {
		c.Calls = make(map[string]int)
	}
	c.Calls[method] += 1
}

// Ensure the API implments the Service interface.
var _ Credentials = &MockCredentials{}

func TestMockCredentials(t *testing.T) {
	creds := &MockCredentials{}

	atks, rtks, err := creds.Login(nil)
	require.NoError(t, err)
	require.NotEmpty(t, atks)
	require.NotEmpty(t, rtks)

	require.Len(t, creds.Calls, 1)
	require.Contains(t, creds.Calls, "Login")
	require.Equal(t, creds.Calls["Login"], 1)

	atks, rtks, err = creds.Refresh(nil)
	require.NoError(t, err)
	require.NotEmpty(t, atks)
	require.NotEmpty(t, rtks)

	require.Len(t, creds.Calls, 2)
	require.Contains(t, creds.Calls, "Refresh")
	require.Equal(t, creds.Calls["Refresh"], 1)

	err = creds.Logout(nil)
	require.NoError(t, err)

	require.Len(t, creds.Calls, 3)
	require.Contains(t, creds.Calls, "Logout")
	require.Equal(t, creds.Calls["Logout"], 1)

	err = creds.Logout(nil)
	require.NoError(t, err)

	require.Len(t, creds.Calls, 3)
	require.Contains(t, creds.Calls, "Logout")
	require.Equal(t, creds.Calls["Logout"], 2)

}
