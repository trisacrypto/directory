package sectigo

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/shibukawa/configdir"
	"gopkg.in/yaml.v2"
)

// Environment variables that are loaded into credentials.
const (
	UsernameEnv = "SECTIGO_USERNAME"
	PasswordEnv = "SECTIGO_PASSWORD"
	ProfileEnv  = "SECTIGO_PROFILE"
)

// Cache directory configuration
const (
	vendorName       = "trisa"
	applicationName  = "sectigo"
	credentialsCache = "credentials.yaml"
)

// Credentials stores login and authentication information to connect to the Sectigo API.
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

// CredentialsManager implements an interface for managing a cached credential store.
// Its primary purpose is to cache access and refresh tokens to prevent multiple logins
// accross different API commands and to store user authentication data or to fetch it
// from the environment. It also provides helper methods for determining when tokens are
// expired by reading the JWT data that has been returned.
type CredentialsManager struct {
	creds *Credentials
}

// CredentialsClient is an interface which can be implemented by either CredentialsManager
// or a mock for testing or emulating the loading and storing of Credential objects.
type CredentialsClient interface {
	Creds() Credentials
	Load(username, password string) error
	Dump() (path string, err error)
	Update(accessToken, refreshToken string) error
	Check() error
	Valid() bool
	Current() bool
	Refreshable() bool
	Clear()
	CacheFile() string
}

// Return a copy of the underyling Creditentials object.
func (c *CredentialsManager) Creds() Credentials {
	return *c.creds
}

// Load initializes a Credentials object. If the username and password are specified,
// they are populated into the credentials, otherwise they are fetched from the
// $SECTIGO_USERNAME and $SECTIGO_PASSWORD environment variables. Access and refresh
// tokens are loaded from an application and OS-specific configuration file if available.
// This method is best effort and does not return intermediate errors. It will return
// an error if the credentials are empty after being loaded.
func (c *CredentialsManager) Load(username, password string) (err error) {
	// Load credentials from the environment
	var ok bool
	if username == "" {
		if username, ok = os.LookupEnv(UsernameEnv); ok {
			c.creds.Username = username
		}
	} else {
		c.creds.Username = username
	}

	if password == "" {
		if password, ok = os.LookupEnv(PasswordEnv); ok {
			c.creds.Password = password
		}
	} else {
		c.creds.Password = password
	}

	// Load tokens from the cache file, stored in an OS-specific application cache, e.g.
	// usually $HOME/.cache or $HOME/Library/Caches for a specific user.
	c.creds.cache = configdir.New(vendorName, applicationName).QueryCacheFolder()
	if c.creds.cache.Exists(credentialsCache) {
		data, _ := c.creds.cache.ReadFile(credentialsCache)
		yaml.Unmarshal(data, &c.creds)
	}

	// Check tokens to ensure they have not expired or are still refreshable.
	if err = c.Check(); err != nil {
		// Tokens are not valid, clear the cache.
		c.Clear()
		c.Dump()
	}

	// Ensure that some credentials are available
	if (c.creds.Username != "" || c.creds.Password != "") && (c.creds.Username == "" || c.creds.Password == "") {
		return ErrCredentialsMismatch
	}
	if c.creds.Username == "" && c.creds.Password == "" && c.creds.AccessToken == "" && c.creds.RefreshToken == "" {
		return ErrNoCredentials
	}

	return nil
}

// Dump the credentials to a local cache file, usually $HOME/.cache or
// $HOME/Library/Caches for a specific user.
func (c *CredentialsManager) Dump() (path string, err error) {
	var data []byte
	if data, err = yaml.Marshal(&c.creds); err != nil {
		return "", err
	}

	// Attempt storage to user folder
	if err = c.creds.cache.WriteFile(credentialsCache, data); err != nil {
		return "", err
	}

	return filepath.Join(c.creds.cache.Path, credentialsCache), nil
}

// Update the credentials with new access and refresh tokens. Credentials are checked
// and if they're ok they are dumped to the cache on disk.
func (c *CredentialsManager) Update(accessToken, refreshToken string) (err error) {
	var atc, rtc *jwt.StandardClaims
	if atc, err = parseToken(accessToken); err != nil {
		return fmt.Errorf("could not parse access token: %s", err)
	}

	if rtc, err = parseToken(refreshToken); err != nil {
		return fmt.Errorf("could not parse refresh token: %s", err)
	}

	c.creds.AccessToken = accessToken
	c.creds.RefreshToken = refreshToken
	c.creds.Subject = atc.Subject
	c.creds.IssuedAt = time.Unix(atc.IssuedAt, 0)
	c.creds.ExpiresAt = time.Unix(atc.ExpiresAt, 0)
	c.creds.RefreshBy = time.Unix(rtc.ExpiresAt, 0)

	if rtc.NotBefore > 0 {
		c.creds.NotBefore = time.Unix(rtc.NotBefore, 0)
	} else {
		c.creds.NotBefore = time.Unix(rtc.IssuedAt, 0)
	}

	if err = c.Check(); err != nil {
		c.Clear()
		c.Dump()
		return err
	}

	// If cache dump errors, do nothing - just keep going without cache
	c.Dump()
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
func (c *CredentialsManager) Check() (err error) {
	// If no access tokens are available, then skip the check.
	if c.creds.AccessToken == "" && c.creds.RefreshToken == "" {
		return nil
	}

	// Both access and refresh tokens are required.
	if (c.creds.AccessToken != "" || c.creds.RefreshToken != "") && (c.creds.AccessToken == "" || c.creds.RefreshToken == "") {
		return ErrTokensMismatch
	}

	// If the current time is after the refresh by, the tokens are expired.
	if !c.Current() {
		return ErrTokensExpired
	}
	return nil
}

// Valid returns true if the access tokens are unexpired.
func (c *CredentialsManager) Valid() bool {
	if c.creds.AccessToken != "" {
		return time.Now().Before(c.creds.ExpiresAt)
	}
	return false
}

// Current returns true if the refresh tokens are unexpired.
func (c *CredentialsManager) Current() bool {
	if c.creds.AccessToken != "" && c.creds.RefreshToken != "" {
		return time.Now().Before(c.creds.RefreshBy)
	}
	return false
}

// Refreshable returns true if the current time is after NotBefore and before RefreshBy.
func (c *CredentialsManager) Refreshable() bool {
	if c.creds.RefreshToken != "" {
		now := time.Now()
		return now.After(c.creds.NotBefore) && now.Before(c.creds.RefreshBy)
	}
	return false
}

// Clear the access and refresh tokens and reset all timestamps.
func (c *CredentialsManager) Clear() {
	zeroTime := time.Time{}

	c.creds.AccessToken = ""
	c.creds.RefreshToken = ""
	c.creds.Subject = ""
	c.creds.IssuedAt = zeroTime
	c.creds.ExpiresAt = zeroTime
	c.creds.NotBefore = zeroTime
	c.creds.RefreshBy = zeroTime
}

// CacheFile returns the path to the credentials cache if it exists.
func (c *CredentialsManager) CacheFile() string {
	if c.creds.cache.Exists(credentialsCache) {
		return filepath.Join(c.creds.cache.Path, credentialsCache)
	}
	return ""
}
