package sectigo

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shibukawa/configdir"
	"gopkg.in/yaml.v2"
)

// Environment variables that are loaded into credentials.
const (
	UsernameEnv = "SECTIGO_USERNAME"
	PasswordEnv = "SECTIGO_PASSWORD"
)

// Cache directory configuration
const (
	vendorName       = "trisa"
	applicationName  = "sectigo"
	credentialsCache = "credentials.yaml"
)

// Credentials stores login and authentication information to connect to the Sectigo API.
// Its primary purpose is to cache access and refresh tokens to prevent multiple logins
// accross different API commands and to store user authentication data or to fetch it
// from the environment. It also provides helper methods for determining when tokens are
// expired by reading the JWT data that has been returned.
type Credentials struct {
	Username     string            `yaml:"-" json:"-"`              // Username is fetched from environment or supplied by user (not stored in cache)
	Password     string            `yaml:"-" json:"-"`              // Password is fetched from environment or supplied by user (not stored in cache)
	AccessToken  string            `yaml:"access_token,omitempty"`  // Temporary bearer token to authenticate API calls; issued on login. Expires after 10 minutes.
	RefreshToken string            `yaml:"refresh_token,omitempty"` // Temporary refresh token to acquire a new access token without reauthentication.
	Subject      string            `yaml:"subject,omitempty"`       // The account and user detail endpoint, e.g. /account/:id/user/:id
	IssuedAt     time.Time         `yaml:"issued_at,omitempty"`     // The timestamp the tokens were issued at
	ExpiresAt    time.Time         `yaml:"expires_at,omitempty"`    // When the access token expires and needs to be refreshed
	NotBefore    time.Time         `yaml:"not_before,omitempty"`    // The earliest timestamp that tokens can be refreshed
	RefreshBy    time.Time         `yaml:"refresh_by,omitempty"`    // The latest timestamp that tokens can be refreshed
	cache        *configdir.Config `yaml:"-"`                       // The cache directory the credentials are loaded and dumped to
}

// Load initializes a Credentials object. If the username and password are specified,
// they are populated into the credentials, otherwise they are fetched from the
// $SECTIGO_USERNAME and $SECTIGO_PASSWORD environment variables. Access and refresh
// tokens are loaded from an application and OS-specific configuration file if available.
// This method is best effort and does not return intermediate errors. It will return
// an error if the credentials are empty after being loaded.
func (creds *Credentials) Load(username, password string) (err error) {
	// Load credentials from the environment
	var ok bool
	if username == "" {
		if username, ok = os.LookupEnv(UsernameEnv); ok {
			creds.Username = username
		}
	} else {
		creds.Username = username
	}

	if password == "" {
		if password, ok = os.LookupEnv(PasswordEnv); ok {
			creds.Password = password
		}
	} else {
		creds.Password = password
	}

	// Load tokens from the cache file, stored in an OS-specific application cache, e.g.
	// usually $HOME/.cache or $HOME/Library/Caches for a specific user.
	creds.cache = configdir.New(vendorName, applicationName).QueryCacheFolder()
	if creds.cache.Exists(credentialsCache) {
		data, _ := creds.cache.ReadFile(credentialsCache)
		yaml.Unmarshal(data, &creds)
	}

	// Check tokens to ensure they have not expired or are still refreshable.
	if err = creds.Check(); err != nil {
		// Tokens are not valid, clear the cache.
		creds.Clear()
		creds.Dump()
	}

	// Ensure that some credentials are available
	if (creds.Username != "" || creds.Password != "") && (creds.Username == "" || creds.Password == "") {
		return ErrCredentialsMismatch
	}
	if creds.Username == "" && creds.Password == "" && creds.AccessToken == "" && creds.RefreshToken == "" {
		return ErrNoCredentials
	}

	return nil
}

// Dump the credentials to a local cache file, usually $HOME/.cache or
// $HOME/Library/Caches for a specific user.
func (creds *Credentials) Dump() (path string, err error) {
	var data []byte
	if data, err = yaml.Marshal(&creds); err != nil {
		return "", err
	}

	// Attempt storage to user folder
	if err = creds.cache.WriteFile(credentialsCache, data); err != nil {
		return "", err
	}

	return filepath.Join(creds.cache.Path, credentialsCache), nil
}

// Update the credentials with new access and refresh tokens. Credentials are checked
// and if they're ok they are dumped to the cache on disk.
func (creds *Credentials) Update(accessToken, refreshToken string) (err error) {
	var atc, rtc *jwt.StandardClaims
	if atc, err = parseToken(accessToken); err != nil {
		return fmt.Errorf("could not parse access token: %s", err)
	}

	if rtc, err = parseToken(refreshToken); err != nil {
		return fmt.Errorf("could not parse refresh token: %s", err)
	}

	creds.AccessToken = accessToken
	creds.RefreshToken = refreshToken
	creds.Subject = atc.Subject
	creds.IssuedAt = time.Unix(atc.IssuedAt, 0)
	creds.ExpiresAt = time.Unix(atc.ExpiresAt, 0)
	creds.RefreshBy = time.Unix(rtc.ExpiresAt, 0)

	if rtc.NotBefore > 0 {
		creds.NotBefore = time.Unix(rtc.NotBefore, 0)
	} else {
		creds.NotBefore = time.Unix(rtc.IssuedAt, 0)
	}

	if err = creds.Check(); err != nil {
		creds.Clear()
		creds.Dump()
		return err
	}

	// If cache dump errors, do nothing - just keep going without cache
	creds.Dump()
	return nil
}

func parseToken(tks string) (_ *jwt.StandardClaims, err error) {
	claims := &jwt.StandardClaims{}
	if _, _, err = new(jwt.Parser).ParseUnverified(tks, claims); err != nil {
		return nil, err
	}

	if claims.IssuedAt == 0 || claims.ExpiresAt == 0 {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// Check reteurns an error if the access and refresh tokens are expired, clearing the
// tokens from the struct. It does not raise an error if no tokens are available.
func (creds *Credentials) Check() (err error) {
	// If no access tokens are available, then skip the check.
	if creds.AccessToken == "" && creds.RefreshToken == "" {
		return nil
	}

	// Both access and refresh tokens are required.
	if (creds.AccessToken != "" || creds.RefreshToken != "") && (creds.AccessToken == "" || creds.RefreshToken == "") {
		return ErrTokensMismatch
	}

	// If the current time is after the refresh by, the tokens are expired.
	if !creds.Current() {
		return ErrTokensExpired
	}
	return nil
}

// Valid returns true if the access tokens are unexpired.
func (creds *Credentials) Valid() bool {
	if creds.AccessToken != "" {
		return time.Now().Before(creds.ExpiresAt)
	}
	return false
}

// Current returns true if the refresh tokens are unexpired.
func (creds *Credentials) Current() bool {
	if creds.AccessToken != "" && creds.RefreshToken != "" {
		return time.Now().Before(creds.RefreshBy)
	}
	return false
}

// Refreshable returns true if the current time is after NotBefore and before RefreshBy.
func (creds *Credentials) Refreshable() bool {
	if creds.RefreshToken != "" {
		now := time.Now()
		return now.After(creds.NotBefore) && now.Before(creds.RefreshBy)
	}
	return false
}

// Clear the access and refresh tokens and reset all timestamps.
func (creds *Credentials) Clear() {
	zeroTime := time.Time{}

	creds.AccessToken = ""
	creds.RefreshToken = ""
	creds.Subject = ""
	creds.IssuedAt = zeroTime
	creds.ExpiresAt = zeroTime
	creds.NotBefore = zeroTime
	creds.RefreshBy = zeroTime
}

// CacheFile returns the path to the credentials cache if it exists.
func (creds *Credentials) CacheFile() string {
	if creds.cache.Exists(credentialsCache) {
		return filepath.Join(creds.cache.Path, credentialsCache)
	}
	return ""
}
