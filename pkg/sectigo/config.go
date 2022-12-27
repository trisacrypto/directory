package sectigo

import (
	"errors"
	"fmt"
	"strings"
)

type Config struct {
	Username string `envconfig:"SECTIGO_USERNAME" required:"false"`
	Password string `envconfig:"SECTIGO_PASSWORD" required:"false"`
	Profile  string `envconfig:"SECTIGO_PROFILE" default:"CipherTrace EE"`
	Endpoint string `envconfig:"SECTIGO_ENDPOINT" required:"false"`
	Testing  bool   `envconfig:"SECTIGO_TESTING" default:"false"`
}

func (c Config) Validate() error {
	// Check valid certificate profiles
	if _, ok := Profiles[c.Profile]; !ok {
		return fmt.Errorf("%q is not a valid Sectigo profile name, specify one of %s", c.Profile, strings.Join(AllProfiles(), ", "))
	}

	// Can only specify an alternative Sectigo endpoint if in testing mode.
	if c.Endpoint != "" && !c.Testing {
		return errors.New("invalid configuration: cannot specify endpoint if not in testing mode")
	}
	return nil
}
