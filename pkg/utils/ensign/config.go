package ensign

import (
	"errors"

	sdk "github.com/rotationalio/go-ensign"
)

// Config defines common configuration for Ensign clients.
type Config struct {
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	Endpoint     string `default:"ensign.rotational.app:443"`
	AuthURL      string `split_words:"true" default:"https://auth.rotational.app"`
	Insecure     bool   `default:"false"`
	Testing      bool   `default:"false"`
}

// Validate that the ensign config is ready for connection.
func (c Config) Validate() (err error) {
	if c.Testing {
		return nil
	}

	if c.ClientID == "" {
		err = errors.Join(err, ErrMissingClientID)
	}

	if c.ClientSecret == "" {
		err = errors.Join(err, ErrMissingClientSecret)
	}

	if c.Endpoint == "" {
		err = errors.Join(err, ErrMissingEndpoint)
	}

	if c.AuthURL == "" {
		err = errors.Join(err, ErrMissingAuthURL)
	}

	return err
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
