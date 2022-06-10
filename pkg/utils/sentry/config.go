package sentry

import (
	"errors"
	"fmt"

	"github.com/trisacrypto/directory/pkg"
)

// Sentry configuration
type Config struct {
	DSN              string  `split_words:"true"`
	Environment      string  `split_words:"true"`
	Release          string  `split_words:"true"`
	TrackPerformance bool    `split_words:"true" default:"false"`
	SampleRate       float64 `split_words:"true" default:"1.0"`
	Debug            bool    `default:"false"`
}

// Returns True if Sentry is enabled.
func (c Config) UseSentry() bool {
	return c.DSN != ""
}

// Returns True if performance tracking is enabled.
func (c Config) UsePerformanceTracking() bool {
	return c.UseSentry() && c.TrackPerformance
}

func (c Config) Validate() error {
	// If Sentry is enabled then the envionment must be set.
	if c.UseSentry() && c.Environment == "" {
		return errors.New("invalid configuration: envrionment must be configured when Sentry is enabled")
	}
	return nil
}

// Get the configured version string or the current semantic version if not configured.
func (c Config) GetRelease() string {
	if c.Release == "" {
		return fmt.Sprintf("gds@%s", pkg.Version())
	}
	return c.Release
}
