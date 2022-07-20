package sentry_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

func TestSentryConfigValidation(t *testing.T) {
	conf := sentry.Config{
		DSN:         "",
		Environment: "",
		Release:     "1.4",
		Debug:       true,
	}

	// If DSN is empty, then Sentry is not enabled
	err := conf.Validate()
	require.NoError(t, err)

	// If Sentry is enabled, then the environment is required
	conf.DSN = "https://something.ingest.sentry.io"
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: envrionment must be configured when Sentry is enabled")

	conf.Environment = "test"
	err = conf.Validate()
	require.NoError(t, err, "expected valid configuration")
}
