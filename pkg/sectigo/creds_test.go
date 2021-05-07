package sectigo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shibukawa/configdir"
	"github.com/stretchr/testify/require"
)

var (
	testAccessToken  = "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiIvYWNjb3VudC80Mi91c2VyLzQyIiwic2NvcGVzIjpbIlJPTEVfVVNFUiJdLCJmaXJzdC1sb2dpbiI6ZmFsc2UsImlzcyI6Imh0dHBzOi8vaW90LnNlY3RpZ28uY29tLyIsImlhdCI6MTYwNjU3MjE4OSwiZXhwIjoxNjA2NTczMDg5fQ.opGnWDYXU_1AlGCSMaZORVO7BKVR-0z5fXlsQUYhlcXxZX0Ma1kXgzUJou218iTX7pFB_38pMA6UUyE3Lpz2XQ"
	testRefreshToken = "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiIvYWNjb3VudC80Mi91c2VyLzQyIiwic2NvcGVzIjpbIlJPTEVfUkVGUkVTSF9UT0tFTiJdLCJpc3MiOiJodHRwczovL2lvdC5zZWN0aWdvLmNvbS8iLCJqdGkiOiI1YTdjOThkYS05ZjA2LTRhMzYtOTBiNy04YmNhYmEwOTFlMTMiLCJpYXQiOjE2MDY1NzIxODksImV4cCI6MTYwNjU3OTM4OX0.9gV4WT1lxbitIhgD0vwyst6eF5XVs4MjIM33fpKbAUddzH6wgMgugKjC9i1ByX-P3lx0I0Zz7r3NOC3sgwzY7g"
)

func TestCredentials(t *testing.T) {
	// Skip test if a cache file already exists
	if err := checkCache(); err != nil {
		t.Skipf(err.Error())
	}

	// Ensure the environment is setup for the test and cleaned up afterward
	os.Setenv(UsernameEnv, "foo")
	os.Setenv(PasswordEnv, "secretz")
	defer os.Clearenv()

	// Load credentials from the environment with no cache
	creds := new(Credentials)
	require.NoError(t, creds.Load("", ""))

	require.Equal(t, "foo", creds.Username)
	require.Equal(t, "secretz", creds.Password)
	require.Zero(t, creds.AccessToken)
	require.Zero(t, creds.RefreshToken)

	// Set expired access and refresh tokens
	require.Error(t, creds.Update(testAccessToken, testRefreshToken))
	require.NoError(t, refreshTokens())

	// Set valid access and refresh tokens
	require.NoError(t, creds.Update(testAccessToken, testRefreshToken))
	require.NotZero(t, creds.AccessToken)
	require.NotZero(t, creds.RefreshToken)
	require.NotZero(t, creds.Subject)
	require.NotZero(t, creds.IssuedAt)
	require.NotZero(t, creds.ExpiresAt)
	require.NotZero(t, creds.NotBefore)
	require.NotZero(t, creds.RefreshBy)
	require.True(t, creds.Valid())
	require.True(t, creds.Current())

	// Load credentials from user supplied values and cached tokens
	require.NoError(t, creds.Load("teller", "tigerpaw"))

	require.Equal(t, "teller", creds.Username)
	require.Equal(t, "tigerpaw", creds.Password)
}

func checkCache() (err error) {
	cdir := configdir.New(vendorName, applicationName).QueryCacheFolder()
	if cdir.Exists(credentialsCache) {
		return fmt.Errorf("credentials already exists at %s", filepath.Join(cdir.Path, credentialsCache))
	}
	return nil
}

func refreshTokens() (err error) {
	signKey := []byte("supersecret")
	claims := apiClaims{
		jwt.StandardClaims{
			Subject:   "/account/42/user/42",
			Issuer:    "https://iot.sectigo.com/",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
		[]string{"ROLE_USER"},
		false,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	if testAccessToken, err = token.SignedString(signKey); err != nil {
		return err
	}

	// Create refresh token
	claims.Scopes = []string{"ROLE_REFRESH_TOKEN"}
	claims.StandardClaims.Id = "5a7c98da-9f06-4a36-90b7-8bcaba091e13"
	claims.StandardClaims.ExpiresAt = time.Now().Add(2 * time.Hour).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	if testRefreshToken, err = token.SignedString(signKey); err != nil {
		return err
	}

	return nil
}

type apiClaims struct {
	jwt.StandardClaims
	Scopes     []string `json:"scopes"`
	FirstLogin bool     `json:"first-login"`
}
