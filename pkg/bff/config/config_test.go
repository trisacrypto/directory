package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

var testEnv = map[string]string{
	"GDS_BFF_MAINTENANCE":              "false",
	"GDS_BFF_BIND_ADDR":                "8080",
	"GDS_BFF_MODE":                     "debug",
	"GDS_BFF_LOG_LEVEL":                "debug",
	"GDS_BFF_CONSOLE_LOG":              "true",
	"GDS_BFF_ALLOW_ORIGINS":            "https://vaspdirectory.net",
	"GDS_BFF_COOKIE_DOMAIN":            "vaspdirectory.net",
	"GDS_BFF_AUTH0_ISSUER":             "example.auth0.com",
	"GDS_BFF_AUTH0_AUDIENCE":           "https://vaspdirectory.net",
	"GDS_BFF_AUTH0_PROVIDER_CACHE":     "10m",
	"GDS_BFF_TESTNET_INSECURE":         "true",
	"GDS_BFF_TESTNET_ENDPOINT":         "localhost:8443",
	"GDS_BFF_TESTNET_TIMEOUT":          "5s",
	"GDS_BFF_MAINNET_INSECURE":         "true",
	"GDS_BFF_MAINNET_ENDPOINT":         "localhost:8444",
	"GDS_BFF_MAINNET_TIMEOUT":          "3s",
	"GDS_BFF_DATABASE_URL":             "trtl://trtl.test:4436",
	"GDS_BFF_DATABASE_REINDEX_ON_BOOT": "false",
	"GDS_BFF_DATABASE_INSECURE":        "true",
	"GDS_BFF_DATABASE_CERT_PATH":       "fixtures/creds/certs.pem",
	"GDS_BFF_DATABASE_POOL_PATH":       "fixtures/creds/pool.zip",
	"GDS_BFF_SENTRY_DSN":               "https://something.ingest.sentry.io",
	"GDS_BFF_SENTRY_ENVIRONMENT":       "test",
	"GDS_BFF_SENTRY_RELEASE":           "1.4",
	"GDS_BFF_SENTRY_DEBUG":             "true",
	"GDS_BFF_SENTRY_TRACK_PERFORMANCE": "true",
	"GDS_BFF_SENTRY_SAMPLE_RATE":       "0.2",
}

func TestConfig(t *testing.T) {
	// Set required environment variables and cleanup after
	prevEnv := curEnv()
	t.Cleanup(func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	})
	setEnv()

	conf, err := config.New()
	require.NoError(t, err)
	require.False(t, conf.IsZero())

	require.False(t, conf.Maintenance)
	require.Equal(t, testEnv["GDS_BFF_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, testEnv["GDS_BFF_MODE"], conf.Mode)
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.Len(t, conf.AllowOrigins, 1)
	require.Equal(t, testEnv["GDS_BFF_COOKIE_DOMAIN"], conf.CookieDomain)
	require.Equal(t, testEnv["GDS_BFF_AUTH0_ISSUER"], conf.Auth0.Issuer)
	require.Equal(t, testEnv["GDS_BFF_AUTH0_AUDIENCE"], conf.Auth0.Audience)
	require.Equal(t, 10*time.Minute, conf.Auth0.ProviderCache)
	require.True(t, conf.TestNet.Insecure)
	require.Equal(t, testEnv["GDS_BFF_TESTNET_ENDPOINT"], conf.TestNet.Endpoint)
	require.Equal(t, 5*time.Second, conf.TestNet.Timeout)
	require.True(t, conf.MainNet.Insecure)
	require.Equal(t, testEnv["GDS_BFF_MAINNET_ENDPOINT"], conf.MainNet.Endpoint)
	require.Equal(t, 3*time.Second, conf.MainNet.Timeout)
	require.Equal(t, testEnv["GDS_BFF_DATABASE_URL"], conf.Database.URL)
	require.Equal(t, false, conf.Database.ReindexOnBoot)
	require.Equal(t, true, conf.Database.Insecure)
	require.Equal(t, testEnv["GDS_BFF_DATABASE_CERT_PATH"], conf.Database.CertPath)
	require.Equal(t, testEnv["GDS_BFF_DATABASE_POOL_PATH"], conf.Database.PoolPath)
	require.Equal(t, testEnv["GDS_BFF_SENTRY_DSN"], conf.Sentry.DSN)
	require.Equal(t, testEnv["GDS_BFF_SENTRY_ENVIRONMENT"], conf.Sentry.Environment)
	require.Equal(t, testEnv["GDS_BFF_SENTRY_RELEASE"], conf.Sentry.Release)
	require.Equal(t, true, conf.Sentry.Debug)
	require.Equal(t, true, conf.Sentry.TrackPerformance)
	require.Equal(t, 0.2, conf.Sentry.SampleRate)
}

func TestRequiredConfig(t *testing.T) {
	required := []string{
		"GDS_BFF_AUTH0_ISSUER",
		"GDS_BFF_AUTH0_AUDIENCE",
		"GDS_BFF_TESTNET_ENDPOINT",
		"GDS_BFF_MAINNET_ENDPOINT",
		"GDS_BFF_DATABASE_URL",
		"GDS_BFF_DATABASE_CERT_PATH",
		"GDS_BFF_DATABASE_POOL_PATH",
	}

	// Collect required environment variables and cleanup after
	prevEnv := curEnv(required...)
	cleanup := func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
	t.Cleanup(cleanup)

	// Ensure that we've captured the complete set of required environment variables
	setEnv(required...)
	conf, err := config.New()
	require.NoError(t, err)

	// Ensure that each environment variable is required
	for _, envvar := range required {
		// Add all environment variables but the current one
		for _, key := range required {
			if key == envvar {
				os.Unsetenv(key)
			} else {
				setEnv(key)
			}
		}

		_, err := config.New()
		require.Errorf(t, err, "expected %q to be required but no error occurred", envvar)
	}

	// Test required configuration
	require.Equal(t, testEnv["GDS_BFF_DATABASE_URL"], conf.Database.URL)
}

func TestAuthConfig(t *testing.T) {
	conf := config.AuthConfig{
		Issuer:        "example.auth0.com",
		Audience:      "https://vaspdirectory.net",
		ProviderCache: 0,
	}

	// Ensure that a provider cache is required
	require.EqualError(t, conf.Validate(), "invalid configuration: auth0 provider cache duration should be longer than 0")

	conf.ProviderCache = 5 * time.Minute
	require.NoError(t, conf.Validate(), "could not validate auth config")

	// Test Domain only configuration (default config)
	conf.Issuer = "example.auth0.com"
	url, err := conf.IssuerURL()
	require.NoError(t, err, "could not parse issuer url")
	require.Equal(t, "https://example.auth0.com/", url.String())

	// Test empty domain invalid configuration
	conf.Issuer = ""
	_, err = conf.IssuerURL()
	require.EqualError(t, err, "invalid configuration: auth0 domain must be configured")
	require.Error(t, conf.Validate())

	// Test URL isuser without trailing slash
	conf.Issuer = "https://example.auth0.com"
	url, err = conf.IssuerURL()
	require.NoError(t, err, "could not parse issuer url")
	require.Equal(t, "https://example.auth0.com/", url.String())

	// Test URL isuser without trailing slash
	conf.Issuer = "https://example.auth0.com/"
	url, err = conf.IssuerURL()
	require.NoError(t, err, "could not parse issuer url")
	require.Equal(t, "https://example.auth0.com/", url.String())

	// Test HTTP URL isuser
	conf.Issuer = "http://example.auth0.com/"
	url, err = conf.IssuerURL()
	require.NoError(t, err, "could not parse issuer url")
	require.Equal(t, "http://example.auth0.com/", url.String())

	// Test invalid url with scheme
	conf.Issuer = "https:///foo"
	_, err = conf.IssuerURL()
	require.EqualError(t, err, "invalid configuration: could not parse issuer url")
	require.Error(t, conf.Validate())

	// Test invalid url with no domain
	conf.Issuer = "https:///"
	_, err = conf.IssuerURL()
	require.EqualError(t, err, "invalid configuration: could not parse issuer url")
	require.Error(t, conf.Validate())

	// Test invalid url with extended path
	// TODO: should this be invalid?
	conf.Issuer = "https://example.auth0.com/foo"
	_, err = conf.IssuerURL()
	require.EqualError(t, err, "invalid configuration: could not parse issuer url")
	require.Error(t, conf.Validate())
}

func TestDatabaseConfigValidation(t *testing.T) {
	conf := config.DatabaseConfig{
		URL:           "trtl://trtl.test.net:443",
		ReindexOnBoot: false,
		Insecure:      false,
		CertPath:      "",
		PoolPath:      "",
	}

	conf.Insecure = true
	err := conf.Validate()
	require.NoError(t, err)

	// If Insecure is false, then the certs and cert pool are required.
	conf.Insecure = false
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: connecting to trtl over mTLS requires certs and cert pool")

	conf.CertPath = "fixtures/certs.pem"
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: connecting to trtl over mTLS requires certs and cert pool")

	conf.PoolPath = "fixtures/pool.zip"
	err = conf.Validate()
	require.NoError(t, err, "expected valid configuration")
}

// Returns the current environment for the specified keys, or if no keys are specified
// then returns the current environment for all keys in testEnv.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)
	if len(keys) > 0 {
		for _, envvar := range keys {
			if val, ok := os.LookupEnv(envvar); ok {
				env[envvar] = val
			}
		}
	} else {
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variable from the testEnv, if no keys are specified, then sets
// all environment variables from the test env.
func setEnv(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := testEnv[key]; ok {
				os.Setenv(key, val)
			}
		}
	} else {
		for key, val := range testEnv {
			os.Setenv(key, val)
		}
	}
}
