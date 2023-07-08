package sectigo

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	envProduction = "production"
	envStaging    = "staging"
	envTesting    = "testing"
)

var environments = map[string]struct{}{
	envProduction: {},
	envStaging:    {},
	envTesting:    {},
}

type Config struct {
	Username    string `envconfig:"SECTIGO_USERNAME" required:"false"`
	Password    string `envconfig:"SECTIGO_PASSWORD" required:"false"`
	Profile     string `envconfig:"SECTIGO_PROFILE" default:"CipherTrace EE"`
	Endpoint    string `envconfig:"SECTIGO_ENDPOINT" default:"https://iot.sectigo.com"`
	Environment string `envconfig:"SECTIGO_ENVIRONMENT" default:"production"`
}

func (c Config) Validate() error {
	// Check valid certificate profiles
	if _, ok := Profiles[c.Profile]; !ok {
		return fmt.Errorf("%q is not a valid Sectigo profile name, specify one of %s", c.Profile, strings.Join(AllProfiles(), ", "))
	}

	// Check valid environments
	if _, ok := environments[c.GetEnvironment()]; !ok {
		return fmt.Errorf("%q is not a valid environment, specify production, staging, or testing", c.GetEnvironment())
	}

	// Check endpoint wrt to environment
	if c.Endpoint != "" {
		ep, err := url.Parse(c.Endpoint)
		if err != nil {
			return fmt.Errorf("could not parse endpoint as url: %w", err)
		}

		switch c.GetEnvironment() {
		case envProduction:
			// If we're in production, make sure we're using TLS and connecting to Sectgio
			if ep.Scheme != "https" {
				return fmt.Errorf("must use https in production, not %s", ep.Scheme)
			}

			host := ep.Hostname()
			if host != "iot.sectigo.com" {
				return fmt.Errorf("cannot connect to %s in production", host)
			}
		case envStaging:
			// If we're in staging mode, ensure we're not connecting to Sectigo
			host := ep.Hostname()
			if strings.HasSuffix(host, "sectigo.com") {
				return fmt.Errorf("cannot connect to %s in staging", host)
			}
		case envTesting:
			// If we're in testing mode, ensure we're connecting to localhost.
			host := ep.Hostname()
			if host != "localhost" && host != "127.0.0.1" {
				return fmt.Errorf("sectigo hostname must be set to localhost in testing mode")
			}
		}
	}
	return nil
}

func (c Config) GetEnvironment() string {
	return strings.ToLower(strings.TrimSpace(c.Environment))
}

func (c Config) Testing() bool {
	return c.GetEnvironment() == envTesting
}
