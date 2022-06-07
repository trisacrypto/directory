package auth0

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type Credentials struct {
	AccessToken string    `json:"access_token,omitempty"`
	ExpiresIn   int       `json:"expires_in,omitempty"`
	Scope       string    `json:"scope,omitempty"`
	TokenType   string    `json:"token_type,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// Valid returns true if there is an access token and it has not expired yet.
func (c *Credentials) Valid() bool {
	if c.AccessToken != "" {
		return c.GetExpiresAt().After(time.Now())
	}
	return false
}

// ExpiresAt returns the timestamp of expiration less 5 seconds as a timeout buffer.
// If there is no requested timestamp or expires in value then a zero-valued timestamp
// is returned to ensure that all checks show the token as expired.
func (c *Credentials) GetExpiresAt() time.Time {
	if c.ExpiresAt.IsZero() {
		if c.ExpiresIn != 0 && !c.CreatedAt.IsZero() {
			seconds := time.Duration(c.ExpiresIn) * time.Second
			c.ExpiresAt = c.CreatedAt.Add(seconds)
		}
	}
	return c.ExpiresAt
}

// Load the credentials from JSON either from a file or an http response. If the created
// at timestamp is zero (e.g. loaded from an http response) then it is set to now. The
// load method should be the primary way of creating credentials.
func (c *Credentials) Load(r io.Reader) error {
	if err := json.NewDecoder(r).Decode(c); err != nil {
		return err
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	return nil
}

// LoadFrom a path on disk.
func (c *Credentials) LoadFrom(path string) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()
	return c.Load(f)
}

// LoadCache loads credentials from a cache file. If the cache file doesn't exist or
// the credentials are invalid, no error is returned, the credentials are just zero.
// All other errors are returned (e.g. permission errors on the file).
func (c *Credentials) LoadCache(path string) (err error) {
	if path != "" {
		if _, err = os.Stat(path); !os.IsNotExist(err) {
			if err = c.LoadFrom(path); err != nil {
				return err
			}
		}
	}

	if !c.Valid() {
		// reset the token to a zero valued state if credentials are invalid
		c.Reset()
	}
	return nil
}

// Dump the credentials as nicely formated JSON data.
func (c *Credentials) Dump(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}

// DumpTo a path on disk.
func (c *Credentials) DumpTo(path string) (err error) {
	// Before dumping, ensure that expires at is set so that it is serialized.
	c.GetExpiresAt()

	var f *os.File
	if f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return err
	}
	defer f.Close()
	return c.Dump(f)
}

// DumpCache attempts to cache the credentials, but only does so if the credentials are
// valid. This won't error if the cache cannot be created, e.g. because of permissions.
func (c *Credentials) DumpCache(path string) {
	if path != "" && c.Valid() {
		c.DumpTo(path)
	}
}

// Reset the credentials to a zero valued struct to force re-authentication.
func (c *Credentials) Reset() {
	c.AccessToken = ""
	c.ExpiresIn = 0
	c.Scope = ""
	c.TokenType = ""
	c.CreatedAt = time.Time{}
	c.ExpiresAt = time.Time{}
}
