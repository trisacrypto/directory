package authtest

import (
	"errors"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

// Config stores the client ID and secrets for accessing auth0 in order to conduct
// "live" tests against our actual development auth0 tenant. If this config is zero or
// invalid then the live tests should be skipped.
type Config struct {
	Domain       string `envconfig:"AUTH0_DOMAIN"`
	ClientID     string `envconfig:"AUTH0_CLIENT_ID"`
	ClientSecret string `envconfig:"AUTH0_CLIENT_SECRET"`
	TokenCache   string `envconfig:"AUTH0_TOKEN_CACHE"`
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
	switch {
	case c.Domain == "":
		return errors.New("invalid configuration: missing auth0 domain")
	case c.ClientID == "":
		return errors.New("invalid configuration: missing auth0 client id")
	case c.ClientSecret == "":
		return errors.New("invalid configuration: missing auth0 client secret")
	case c.TokenCache == "":
		return errors.New("invalid configuration: missing path to token cache")
	default:
		return nil
	}
}

func (c Config) IsZero() bool {
	return c.Domain == "" && c.ClientID == "" && c.ClientSecret == "" && c.TokenCache == ""
}

func (c Config) AuthConfig() config.AuthConfig {
	return config.AuthConfig{
		Domain:        c.Domain,
		Audience:      "https://bff.trisa.directory",
		ProviderCache: 1 * time.Minute,
		ClientID:      c.ClientID,
		ClientSecret:  c.ClientSecret,
		Testing:       true,
	}
}
