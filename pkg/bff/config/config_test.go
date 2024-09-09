package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

var testEnv = map[string]string{
	"GDS_BFF_MAINTENANCE":                   "false",
	"GDS_BFF_BIND_ADDR":                     "8080",
	"GDS_BFF_MODE":                          "debug",
	"GDS_BFF_LOG_LEVEL":                     "debug",
	"GDS_BFF_CONSOLE_LOG":                   "true",
	"GDS_BFF_ALLOW_ORIGINS":                 "https://vaspdirectory.net",
	"GDS_BFF_COOKIE_DOMAIN":                 "vaspdirectory.net",
	"GDS_BFF_LOGIN_URL":                     "https://vaspdirectory.net/auth/login",
	"GDS_BFF_REGISTER_URL":                  "https://vaspdirectory.net/auth/register",
	"GDS_BFF_SERVE_DOCS":                    "true",
	"GDS_BFF_AUTH0_DOMAIN":                  "example.auth0.com",
	"GDS_BFF_AUTH0_ISSUER":                  "https://auth.example.com",
	"GDS_BFF_AUTH0_AUDIENCE":                "https://vaspdirectory.net",
	"GDS_BFF_AUTH0_PROVIDER_CACHE":          "10m",
	"GDS_BFF_AUTH0_CLIENT_ID":               "exampleid",
	"GDS_BFF_AUTH0_CLIENT_SECRET":           "supersecretsquirrel",
	"GDS_BFF_AUTH0_TESTING":                 "true",
	"GDS_BFF_TESTNET_DATABASE_URL":          "trtl://trtl.testnet:4436",
	"GDS_BFF_TESTNET_DATABASE_INSECURE":     "true",
	"GDS_BFF_TESTNET_DATABASE_CERT_PATH":    "fixtures/creds/testnet/trtl/certs.pem",
	"GDS_BFF_TESTNET_DATABASE_POOL_PATH":    "fixtures/creds/testnet/trtl/pool.zip",
	"GDS_BFF_TESTNET_DIRECTORY_INSECURE":    "true",
	"GDS_BFF_TESTNET_DIRECTORY_ENDPOINT":    "localhost:8443",
	"GDS_BFF_TESTNET_DIRECTORY_TIMEOUT":     "5s",
	"GDS_BFF_TESTNET_MEMBERS_MTLS_INSECURE": "true",
	"GDS_BFF_TESTNET_MEMBERS_ENDPOINT":      "localhost:9443",
	"GDS_BFF_TESTNET_MEMBERS_TIMEOUT":       "5s",
	"GDS_BFF_TESTNET_MEMBERS_CERT_PATH":     "fixtures/members/creds/testnet/certs.pem",
	"GDS_BFF_TESTNET_MEMBERS_POOL_PATH":     "fixtures/members/creds/testnet/pool.zip",
	"GDS_BFF_MAINNET_DATABASE_URL":          "trtl://trtl.mainnet:4436",
	"GDS_BFF_MAINNET_DATABASE_INSECURE":     "true",
	"GDS_BFF_MAINNET_DATABASE_CERT_PATH":    "fixtures/creds/mainnet/trtl/certs.pem",
	"GDS_BFF_MAINNET_DATABASE_POOL_PATH":    "fixtures/creds/mainnet/trtl/pool.zip",
	"GDS_BFF_MAINNET_DIRECTORY_INSECURE":    "true",
	"GDS_BFF_MAINNET_DIRECTORY_ENDPOINT":    "localhost:8444",
	"GDS_BFF_MAINNET_DIRECTORY_TIMEOUT":     "3s",
	"GDS_BFF_MAINNET_MEMBERS_MTLS_INSECURE": "true",
	"GDS_BFF_MAINNET_MEMBERS_ENDPOINT":      "localhost:9444",
	"GDS_BFF_MAINNET_MEMBERS_TIMEOUT":       "3s",
	"GDS_BFF_MAINNET_MEMBERS_CERT_PATH":     "fixtures/members/creds/mainnet/certs.pem",
	"GDS_BFF_MAINNET_MEMBERS_POOL_PATH":     "fixtures/members/creds/mainnet/pool.zip",
	"GDS_BFF_DATABASE_URL":                  "trtl://trtl.test:4436",
	"GDS_BFF_DATABASE_REINDEX_ON_BOOT":      "false",
	"GDS_BFF_DATABASE_INSECURE":             "true",
	"GDS_BFF_DATABASE_CERT_PATH":            "fixtures/creds/certs.pem",
	"GDS_BFF_DATABASE_POOL_PATH":            "fixtures/creds/pool.zip",
	"GDS_BFF_SERVICE_EMAIL":                 "test@example.com",
	"SENDGRID_API_KEY":                      "foo1234",
	"GDS_BFF_EMAIL_TESTING":                 "true",
	"GDS_BFF_EMAIL_STORAGE":                 "fixtures/emails",
	"GDS_BFF_SENTRY_DSN":                    "https://something.ingest.sentry.io",
	"GDS_BFF_SENTRY_ENVIRONMENT":            "test",
	"GDS_BFF_SENTRY_RELEASE":                "1.4",
	"GDS_BFF_SENTRY_DEBUG":                  "true",
	"GDS_BFF_SENTRY_TRACK_PERFORMANCE":      "true",
	"GDS_BFF_SENTRY_SAMPLE_RATE":            "0.2",
	"GDS_BFF_USER_CACHE_ENABLED":            "true",
	"GDS_BFF_USER_CACHE_EXPIRATION":         "10h",
	"GDS_BFF_USER_CACHE_SIZE":               "1000",
	"GDS_BFF_ACTIVITY_ENABLED":              "true",
	"GDS_BFF_ACTIVITY_TOPIC":                "network-activity",
	"GDS_BFF_ACTIVITY_NETWORK":              "testnet",
	"GDS_BFF_ACTIVITY_ENSIGN_CLIENT_ID":     "client-id",
	"GDS_BFF_ACTIVITY_ENSIGN_CLIENT_SECRET": "client-secret",
	"GDS_BFF_ACTIVITY_ENSIGN_ENDPOINT":      "api.ensign.world:443",
	"GDS_BFF_ACTIVITY_ENSIGN_AUTH_URL":      "https://auth.ensign.world",
	"GDS_BFF_ACTIVITY_ENSIGN_INSECURE":      "true",
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
	require.Equal(t, testEnv["GDS_BFF_LOGIN_URL"], conf.LoginURL)
	require.Equal(t, testEnv["GDS_BFF_REGISTER_URL"], conf.RegisterURL)
	require.True(t, conf.ServeDocs)
	require.Equal(t, testEnv["GDS_BFF_AUTH0_DOMAIN"], conf.Auth0.Domain)
	require.Equal(t, testEnv["GDS_BFF_AUTH0_ISSUER"], conf.Auth0.Issuer)
	require.Equal(t, testEnv["GDS_BFF_AUTH0_AUDIENCE"], conf.Auth0.Audience)
	require.Equal(t, testEnv["GDS_BFF_AUTH0_CLIENT_ID"], conf.Auth0.ClientID)
	require.Equal(t, testEnv["GDS_BFF_AUTH0_CLIENT_SECRET"], conf.Auth0.ClientSecret)
	require.True(t, conf.Auth0.Testing)
	require.Equal(t, 10*time.Minute, conf.Auth0.ProviderCache)
	require.Equal(t, testEnv["GDS_BFF_TESTNET_DATABASE_URL"], conf.TestNet.Database.URL)
	require.True(t, conf.TestNet.Database.Insecure)
	require.Equal(t, testEnv["GDS_BFF_TESTNET_DATABASE_CERT_PATH"], conf.TestNet.Database.CertPath)
	require.Equal(t, testEnv["GDS_BFF_TESTNET_DATABASE_POOL_PATH"], conf.TestNet.Database.PoolPath)
	require.True(t, conf.TestNet.Directory.Insecure)
	require.True(t, conf.TestNet.Members.MTLS.Insecure)
	require.Equal(t, testEnv["GDS_BFF_TESTNET_DIRECTORY_ENDPOINT"], conf.TestNet.Directory.Endpoint)
	require.Equal(t, testEnv["GDS_BFF_TESTNET_MEMBERS_ENDPOINT"], conf.TestNet.Members.Endpoint)
	require.Equal(t, 5*time.Second, conf.TestNet.Directory.Timeout)
	require.Equal(t, 5*time.Second, conf.TestNet.Members.Timeout)
	require.Equal(t, testEnv["GDS_BFF_TESTNET_MEMBERS_MTLS_CERT_PATH"], conf.TestNet.Members.MTLS.CertPath)
	require.Equal(t, testEnv["GDS_BFF_TESTNET_MEMBERS_MTLS_POOL_PATH"], conf.TestNet.Members.MTLS.PoolPath)
	require.Equal(t, testEnv["GDS_BFF_MAINNET_DATABASE_URL"], conf.MainNet.Database.URL)
	require.True(t, conf.MainNet.Database.Insecure)
	require.Equal(t, testEnv["GDS_BFF_MAINNET_DATABASE_CERT_PATH"], conf.MainNet.Database.CertPath)
	require.Equal(t, testEnv["GDS_BFF_MAINNET_DATABASE_POOL_PATH"], conf.MainNet.Database.PoolPath)
	require.True(t, conf.MainNet.Directory.Insecure)
	require.True(t, conf.MainNet.Members.MTLS.Insecure)
	require.Equal(t, testEnv["GDS_BFF_MAINNET_DIRECTORY_ENDPOINT"], conf.MainNet.Directory.Endpoint)
	require.Equal(t, testEnv["GDS_BFF_MAINNET_MEMBERS_ENDPOINT"], conf.MainNet.Members.Endpoint)
	require.Equal(t, 3*time.Second, conf.MainNet.Directory.Timeout)
	require.Equal(t, 3*time.Second, conf.MainNet.Members.Timeout)
	require.Equal(t, testEnv["GDS_BFF_MAINNET_MEMBERS_MTLS_CERT_PATH"], conf.MainNet.Members.MTLS.CertPath)
	require.Equal(t, testEnv["GDS_BFF_MAINNET_MEMBERS_MTLS_POOL_PATH"], conf.MainNet.Members.MTLS.PoolPath)
	require.Equal(t, testEnv["GDS_BFF_DATABASE_URL"], conf.Database.URL)
	require.Equal(t, false, conf.Database.ReindexOnBoot)
	require.Equal(t, true, conf.Database.Insecure)
	require.Equal(t, testEnv["GDS_BFF_DATABASE_CERT_PATH"], conf.Database.CertPath)
	require.Equal(t, testEnv["GDS_BFF_DATABASE_POOL_PATH"], conf.Database.PoolPath)
	require.Equal(t, testEnv["GDS_BFF_SERVICE_EMAIL"], conf.Email.ServiceEmail)
	require.Equal(t, testEnv["SENDGRID_API_KEY"], conf.Email.SendGridAPIKey)
	require.True(t, conf.Email.Testing)
	require.Equal(t, testEnv["GDS_BFF_EMAIL_STORAGE"], conf.Email.Storage)
	require.Equal(t, testEnv["GDS_BFF_SENTRY_DSN"], conf.Sentry.DSN)
	require.Equal(t, testEnv["GDS_BFF_SENTRY_ENVIRONMENT"], conf.Sentry.Environment)
	require.Equal(t, testEnv["GDS_BFF_SENTRY_RELEASE"], conf.Sentry.Release)
	require.True(t, conf.UserCache.Enabled)
	require.Equal(t, 10*time.Hour, conf.UserCache.Expiration)
	require.Equal(t, uint(1000), conf.UserCache.Size)
	require.Equal(t, true, conf.Sentry.Debug)
	require.Equal(t, true, conf.Sentry.TrackPerformance)
	require.Equal(t, 0.2, conf.Sentry.SampleRate)
	require.True(t, conf.Activity.Enabled)
	require.Equal(t, testEnv["GDS_BFF_ACTIVITY_TOPIC"], conf.Activity.Topic)
	require.Equal(t, testEnv["GDS_BFF_ACTIVITY_ENSIGN_CLIENT_ID"], conf.Activity.Ensign.ClientID)
	require.Equal(t, testEnv["GDS_BFF_ACTIVITY_ENSIGN_CLIENT_SECRET"], conf.Activity.Ensign.ClientSecret)
	require.Equal(t, testEnv["GDS_BFF_ACTIVITY_ENSIGN_ENDPOINT"], conf.Activity.Ensign.Endpoint)
	require.Equal(t, testEnv["GDS_BFF_ACTIVITY_ENSIGN_AUTH_URL"], conf.Activity.Ensign.AuthURL)
	require.Equal(t, true, conf.Activity.Ensign.Insecure)
}

func TestRequiredConfig(t *testing.T) {
	t.Skip("test assumes that confire is processing required tags recursively, is it?")
	required := []string{
		"GDS_BFF_LOGIN_URL",
		"GDS_BFF_REGISTER_URL",
		"GDS_BFF_AUTH0_DOMAIN",
		"GDS_BFF_AUTH0_AUDIENCE",
		"GDS_BFF_AUTH0_CLIENT_ID",
		"GDS_BFF_AUTH0_CLIENT_SECRET",
		"GDS_BFF_TESTNET_DATABASE_URL",
		"GDS_BFF_TESTNET_DIRECTORY_ENDPOINT",
		"GDS_BFF_TESTNET_MEMBERS_ENDPOINT",
		"GDS_BFF_MAINNET_DATABASE_URL",
		"GDS_BFF_MAINNET_DIRECTORY_ENDPOINT",
		"GDS_BFF_MAINNET_MEMBERS_ENDPOINT",
		"GDS_BFF_DATABASE_URL",
		"SENDGRID_API_KEY",
	}

	// Insecure must be true if no mTLS certs are provided
	os.Setenv("GDS_BFF_TESTNET_MEMBERS_MTLS_INSECURE", "true")
	os.Setenv("GDS_BFF_MAINNET_MEMBERS_MTLS_INSECURE", "true")
	os.Setenv("GDS_BFF_DATABASE_INSECURE", "true")

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
		fmt.Println(err)
		require.Errorf(t, err, "expected %q to be required but no error occurred", envvar)
	}

	// Test required configuration
	require.Equal(t, testEnv["GDS_BFF_DATABASE_URL"], conf.Database.URL)
}

func TestAuthConfig(t *testing.T) {
	conf := config.AuthConfig{
		Domain:        "example.auth0.com",
		Audience:      "https://vaspdirectory.net",
		ProviderCache: 0,
		Testing:       true,
	}

	// Ensure that a provider cache is required
	require.EqualError(t, conf.Validate(), "invalid configuration: auth0 provider cache duration should be longer than 0")

	// Ensure that client Id and secret are required when not testing
	conf.ProviderCache = 5 * time.Minute
	conf.Testing = false
	require.EqualError(t, conf.Validate(), "invalid configuration: auth0 client id is required in production")
	conf.ClientID = "exampleid"
	require.EqualError(t, conf.Validate(), "invalid configuration: auth0 client secret is required in production")
	conf.ClientSecret = "supersecretpassword"
	require.NoError(t, conf.Validate(), "could not validate auth config")

	// Test Domain only configuration (default config)
	conf.Domain = "example.auth0.com"
	url, err := conf.IssuerURL()
	require.NoError(t, err, "could not parse issuer url")
	require.Equal(t, "https://example.auth0.com/", url.String())

	// Test empty domain invalid configuration
	conf.Domain = ""
	_, err = conf.IssuerURL()
	require.EqualError(t, err, "invalid configuration: auth0 domain must be configured")
	require.Error(t, conf.Validate())

	// Test issuer url formatting returns an invalid configuration error
	for _, scheme := range []string{"", "http://", "https://"} {
		for _, suffix := range []string{"", "/"} {
			if scheme == "" && suffix == "" {
				continue
			}

			conf.Domain = scheme + "example.auth0.com" + suffix
			_, err := conf.IssuerURL()
			require.EqualError(t, err, "invalid configuration: auth0 domain must not be a url or have a trailing slash")
			require.Error(t, conf.Validate())
		}
	}

	// Ensure that if issuer is set it is returned instead of the domain
	conf.Issuer = "https://auth.example.com/"
	u, err := conf.IssuerURL()
	require.NoError(t, err, "could not parse issuer string")
	require.Equal(t, conf.Issuer, u.String())
}

func TestMembersConfigValidation(t *testing.T) {
	conf := config.MembersConfig{
		Endpoint: "https://example.com",
		MTLS: config.MTLSConfig{
			Insecure: false,
			CertPath: "",
			PoolPath: "",
		},
	}

	conf.MTLS.Insecure = true
	err := conf.Validate()
	require.NoError(t, err)

	// If Insecure is false, then the certs and cert pool are required.
	conf.MTLS.Insecure = false
	err = conf.Validate()
	require.EqualError(t, err, "invalid members configuration: connecting over mTLS requires certs and cert pool")

	conf.MTLS.CertPath = "fixtures/certs.pem"
	err = conf.Validate()
	require.EqualError(t, err, "invalid members configuration: connecting over mTLS requires certs and cert pool")

	conf.MTLS.PoolPath = "fixtures/pool.zip"
	err = conf.Validate()
	require.NoError(t, err, "expected valid configuration")
}

func TestEmailConfigValidation(t *testing.T) {
	conf := config.EmailConfig{}
	err := conf.Validate()
	require.EqualError(t, err, "invalid configuration: sendgrid api key and service email are required")

	conf.SendGridAPIKey = "supersecretapikey"
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: sendgrid api key and service email are required")

	conf.ServiceEmail = "service@example.com"
	conf.Storage = "fixtures/emails"
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: email archiving is only supported in testing mode")

	conf.Testing = true
	err = conf.Validate()
	require.NoError(t, err, "expected valid configuration")
}

func TestCacheConfigValidation(t *testing.T) {
	conf := config.CacheConfig{
		Size:    100,
		Enabled: true,
	}
	err := conf.Validate()
	require.EqualError(t, err, "invalid configuration: cache expiration must be greater than 0")

	conf.Expiration = time.Hour
	conf.Size = 0
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: cache size must be greater than 0")

	conf.Enabled = false
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
