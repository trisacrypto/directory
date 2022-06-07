package auth0

import (
	"errors"
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

const (
	contentType = "application/json"
	userAgent   = "TRISA GDS Auth0 Client v1.0"
)

// Config stores the client ID and secrets for accessing auth0. It is possible to load
// the configuration from the environment, but generally speaking the configuration is
// created from process-specifc environment configs like the one in the BFF config.
type Config struct {
	Domain       string `envconfig:"AUTH0_DOMAIN"`
	ClientID     string `envconfig:"AUTH0_CLIENT_ID"`
	ClientSecret string `envconfig:"AUTH0_CLIENT_SECRET"`
	TokenCache   string `envconfig:"AUTH0_TOKEN_CACHE"`
	Testing      bool   `envconfig:"AUTH0_TESTING"`
}

func NewConfig() (conf Config, err error) {
	if err = envconfig.Process("", &conf); err != nil {
		return Config{}, err
	}

	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	return conf, nil
}

func (c Config) Validate() error {
	if c.Domain == "" {
		return errors.New("invalid configuration: missing auth0 domain")
	}

	if !c.Testing {
		if c.ClientID == "" {
			return errors.New("invalid configuration: missing auth0 client id")
		}

		if c.ClientSecret == "" {
			return errors.New("invalid configuration: missing auth0 client secret")
		}
	}
	return nil
}

func (c Config) IsZero() bool {
	return c.Domain == "" && c.ClientID == "" && c.ClientSecret == "" && !c.Testing
}

func (c Config) BaseURL() *url.URL {
	if c.Testing {
		return &url.URL{Scheme: "http", Host: c.Domain}
	}
	return &url.URL{Scheme: "https", Host: c.Domain}
}
