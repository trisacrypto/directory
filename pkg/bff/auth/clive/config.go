package clive

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

// Config stores the client ID and secrets for accessing auth0 in order to conduct
// "live" authentication on a CLI from the localhost.
type Config struct {
	Domain     string `envconfig:"AUTH0_DOMAIN"`
	Audience   string `envconfig:"AUTH0_AUDIENCE"`
	ClientID   string `envconfig:"AUTH0_CLIENT_ID"`
	TokenCache string `envconfig:"AUTH0_TOKEN_CACHE"`
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
		return errors.New("invalid configuration: missing $AUTH0_DOMAIN")
	case c.Audience == "":
		return errors.New("invalid configuration: missing $AUTH0_AUDIENCE")
	case c.ClientID == "":
		return errors.New("invalid configuration: missing $AUTH0_CLIENT_ID")
	case c.TokenCache == "":
		return errors.New("invalid configuration: missing path in $AUTH0_TOKEN_CACHE")
	default:
		return nil
	}
}

func (c Config) IsZero() bool {
	return c.Domain == "" && c.ClientID == "" && c.TokenCache == ""
}
