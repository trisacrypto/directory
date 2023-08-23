package ensign

import (
	sdk "github.com/rotationalio/go-ensign"
)

// Config defines common configuration for Ensign clients.
type Config struct {
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	Endpoint     string `split_words:"true" default:"ensign.rotational.app:443"`
	AuthURL      string `split_words:"true" default:"https://auth.rotational.app"`
	Insecure     bool   `split_words:"true" default:"false"`
	Testing      bool   `split_words:"true" default:"false"`
}

func (c Config) Validate() error {
	if c.ClientID == "" {
		return ErrMissingClientID
	}

	if c.ClientSecret == "" {
		return ErrMissingClientSecret
	}

	if c.Endpoint == "" {
		return ErrMissingEndpoint
	}

	if c.AuthURL == "" {
		return ErrMissingAuthURL
	}

	return nil
}

func (c Config) ClientOptions() []sdk.Option {
	return []sdk.Option{
		sdk.WithCredentials(c.ClientID, c.ClientSecret),
		sdk.WithEnsignEndpoint(c.Endpoint, c.Insecure),
		sdk.WithAuthenticator(c.AuthURL, false),
	}
}

// Create an Ensign client from the configuration.
func (c Config) Client() (*sdk.Client, error) {
	return sdk.New(c.ClientOptions()...)
}
