package sectigo

import (
	"fmt"
	"strings"
)

type Config struct {
	Username string `envconfig:"SECTIGO_USERNAME" required:"false"`
	Password string `envconfig:"SECTIGO_PASSWORD" required:"false"`
	Profile  string `envconfig:"SECTIGO_PROFILE" default:"CipherTrace EE"`
	Testing  bool   `envconfig:"SECTIGO_TESTING" default:"false"`
}

func (c Config) Validate() error {
	// Check valid certificate profiles
	validProfile := false
	for _, profile := range AllProfiles {
		if profile == c.Profile {
			validProfile = true
			break
		}
	}

	if !validProfile {
		return fmt.Errorf("%q is not a valid Sectigo profile name, specify one of %s", c.Profile, strings.Join(AllProfiles[:], ", "))
	}
	return nil
}
