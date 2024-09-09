package config_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/config"
)

var testEnv = map[string]string{
	"GDS_MAINTENANCE":                          "false",
	"GDS_DIRECTORY_ID":                         "testdirectory.org",
	"GDS_SECRET_KEY":                           "theeaglefliesatmidnight",
	"GDS_LOG_LEVEL":                            "debug",
	"GDS_CONSOLE_LOG":                          "true",
	"GDS_API_ENABLED":                          "true",
	"GDS_BIND_ADDR":                            ":443",
	"GDS_ADMIN_ENABLED":                        "true",
	"GDS_ADMIN_BIND_ADDR":                      ":444",
	"GDS_ADMIN_MODE":                           "debug",
	"GDS_ADMIN_TOKEN_KEYS":                     "1y9fT85qWaIvAAORW7DKxtpz9FB:testdata/key1.pem,1y9fVjaUlsVdFFDUWlvRq2PLkw3:testdata/key2.pem",
	"GDS_ADMIN_OAUTH_GOOGLE_AUDIENCE":          "abc-1234.example.fakegoogleusercontent.com",
	"GDS_ADMIN_OAUTH_AUTHORIZED_EMAIL_DOMAINS": "trisa.io,vaspdirectory.net,trisatest.net",
	"GDS_ADMIN_ALLOW_ORIGINS":                  "https://admin.trisatest.net",
	"GDS_ADMIN_COOKIE_DOMAIN":                  "admin.trisatest.net",
	"GDS_ADMIN_AUDIENCE":                       "https://api.admin.trisatest.net",
	"GDS_MEMBERS_ENABLED":                      "true",
	"GDS_MEMBERS_BIND_ADDR":                    ":445",
	"GDS_MEMBERS_INSECURE":                     "true",
	"GDS_MEMBERS_CERTS":                        "fixtures/creds/gds.gz",
	"GDS_MEMBERS_CERT_POOL":                    "fixtures/creds/pool.gz",
	"GDS_DATABASE_URL":                         "trtl://trtl.test:4436",
	"GDS_DATABASE_REINDEX_ON_BOOT":             "false",
	"GDS_DATABASE_INSECURE":                    "true",
	"GDS_DATABASE_CERT_PATH":                   "fixtures/creds/certs.pem",
	"GDS_DATABASE_POOL_PATH":                   "fixtures/creds/pool.zip",
	"SECTIGO_USERNAME":                         "foo",
	"SECTIGO_PASSWORD":                         "supersecret",
	"SECTIGO_PROFILE":                          "17",
	"SECTIGO_ENVIRONMENT":                      "staging",
	"SECTIGO_ENDPOINT":                         "https://cathy.io",
	"GDS_SERVICE_EMAIL":                        "test@example.com",
	"GDS_ADMIN_EMAIL":                          "admin@example.com",
	"SENDGRID_API_KEY":                         "bar1234",
	"GDS_VERIFY_CONTACT_URL":                   "http://localhost:3000/verify",
	"GDS_ADMIN_REVIEW_URL":                     "http://localhost:3001/vasps/",
	"GDS_EMAIL_TESTING":                        "true",
	"GDS_EMAIL_STORAGE":                        "fixtures/emails",
	"GDS_CERTMAN_ENABLED":                      "false",
	"GDS_CERTMAN_REQUEST_INTERVAL":             "60s",
	"GDS_CERTMAN_REISSUANCE_INTERVAL":          "90s",
	"GDS_CERTMAN_STORAGE":                      "fixtures/certs",
	"GDS_BACKUP_ENABLED":                       "true",
	"GDS_BACKUP_INTERVAL":                      "36h",
	"GDS_BACKUP_STORAGE":                       "fixtures/backups",
	"GDS_BACKUP_KEEP":                          "7",
	"GOOGLE_APPLICATION_CREDENTIALS":           "test.json",
	"GOOGLE_PROJECT_NAME":                      "test",
	"GDS_SECRETS_TESTING":                      "true",
	"GDS_SENTRY_DSN":                           "https://something.ingest.sentry.io",
	"GDS_SENTRY_ENVIRONMENT":                   "test",
	"GDS_SENTRY_RELEASE":                       "1.4",
	"GDS_SENTRY_DEBUG":                         "true",
	"GDS_SENTRY_TRACK_PERFORMANCE":             "true",
	"GDS_SENTRY_SAMPLE_RATE":                   "0.2",
	"GDS_ACTIVITY_ENABLED":                     "true",
	"GDS_ACTIVITY_TOPIC":                       "gds-activity",
	"GDS_ACTIVITY_NETWORK":                     "testnet",
	"GDS_ACTIVITY_AGGREGATION_WINDOW":          "10m",
	"GDS_ACTIVITY_ENSIGN_CLIENT_ID":            "client-id",
	"GDS_ACTIVITY_ENSIGN_CLIENT_SECRET":        "client-secret",
	"GDS_ACTIVITY_ENSIGN_ENDPOINT":             "api.ensign.world:443",
	"GDS_ACTIVITY_ENSIGN_AUTH_URL":             "https://auth.ensign.world",
	"GDS_ACTIVITY_ENSIGN_INSECURE":             "true",
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

	// Test configuration set from the environment
	require.Equal(t, false, conf.Maintenance)
	require.Equal(t, testEnv["GDS_DIRECTORY_ID"], conf.DirectoryID)
	require.Equal(t, testEnv["GDS_SECRET_KEY"], conf.SecretKey)
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.Equal(t, true, conf.ConsoleLog)
	require.Equal(t, true, conf.GDS.Enabled)
	require.Equal(t, testEnv["GDS_BIND_ADDR"], conf.GDS.BindAddr)
	require.Equal(t, true, conf.Admin.Enabled)
	require.Equal(t, testEnv["GDS_ADMIN_BIND_ADDR"], conf.Admin.BindAddr)
	require.Equal(t, testEnv["GDS_ADMIN_MODE"], conf.Admin.Mode)
	require.Len(t, conf.Admin.TokenKeys, 2)
	require.Equal(t, testEnv["GDS_ADMIN_OAUTH_GOOGLE_AUDIENCE"], conf.Admin.Oauth.GoogleAudience)
	require.Len(t, conf.Admin.Oauth.AuthorizedEmailDomains, 3)
	require.Len(t, conf.Admin.AllowOrigins, 1)
	require.Equal(t, testEnv["GDS_ADMIN_COOKIE_DOMAIN"], conf.Admin.CookieDomain)
	require.Equal(t, testEnv["GDS_ADMIN_AUDIENCE"], conf.Admin.Audience)
	require.True(t, conf.Members.Enabled)
	require.Equal(t, testEnv["GDS_MEMBERS_BIND_ADDR"], conf.Members.BindAddr)
	require.True(t, conf.Members.Insecure)
	require.Equal(t, testEnv["GDS_MEMBERS_CERTS"], conf.Members.Certs)
	require.Equal(t, testEnv["GDS_MEMBERS_CERT_POOL"], conf.Members.CertPool)
	require.Equal(t, testEnv["GDS_DATABASE_URL"], conf.Database.URL)
	require.Equal(t, false, conf.Database.ReindexOnBoot)
	require.Equal(t, true, conf.Database.Insecure)
	require.Equal(t, testEnv["GDS_DATABASE_CERT_PATH"], conf.Database.CertPath)
	require.Equal(t, testEnv["GDS_DATABASE_POOL_PATH"], conf.Database.PoolPath)
	require.Equal(t, testEnv["SECTIGO_USERNAME"], conf.CertMan.Sectigo.Username)
	require.Equal(t, testEnv["SECTIGO_PASSWORD"], conf.CertMan.Sectigo.Password)
	require.Equal(t, testEnv["SECTIGO_PROFILE"], conf.CertMan.Sectigo.Profile)
	require.Equal(t, testEnv["SECTIGO_ENVIRONMENT"], conf.CertMan.Sectigo.Environment)
	require.Equal(t, testEnv["SECTIGO_ENDPOINT"], conf.CertMan.Sectigo.Endpoint)
	require.Equal(t, testEnv["GDS_SERVICE_EMAIL"], conf.Email.ServiceEmail)
	require.Equal(t, testEnv["GDS_ADMIN_EMAIL"], conf.Email.AdminEmail)
	require.Equal(t, testEnv["SENDGRID_API_KEY"], conf.Email.SendGridAPIKey)
	require.Equal(t, testEnv["GDS_VERIFY_CONTACT_URL"], conf.Email.VerifyContactBaseURL)
	require.Equal(t, testEnv["GDS_ADMIN_REVIEW_URL"], conf.Email.AdminReviewBaseURL)
	require.Equal(t, testEnv["GDS_EMAIL_STORAGE"], conf.Email.Storage)
	require.True(t, conf.Email.Testing)
	require.Equal(t, testEnv["GDS_DIRECTORY_ID"], conf.Email.DirectoryID)
	require.False(t, conf.CertMan.Enabled)
	require.Equal(t, 1*time.Minute, conf.CertMan.RequestInterval)
	require.Equal(t, 90*time.Second, conf.CertMan.ReissuanceInterval)
	require.Equal(t, testEnv["GDS_CERTMAN_STORAGE"], conf.CertMan.Storage)
	require.Equal(t, testEnv["GDS_DIRECTORY_ID"], conf.CertMan.DirectoryID)
	require.Equal(t, true, conf.Backup.Enabled)
	require.Equal(t, 36*time.Hour, conf.Backup.Interval)
	require.Equal(t, testEnv["GDS_BACKUP_STORAGE"], conf.Backup.Storage)
	require.Equal(t, 7, conf.Backup.Keep)
	require.Equal(t, testEnv["GOOGLE_APPLICATION_CREDENTIALS"], conf.Secrets.Credentials)
	require.Equal(t, testEnv["GOOGLE_PROJECT_NAME"], conf.Secrets.Project)
	require.Equal(t, testEnv["GDS_SENTRY_DSN"], conf.Sentry.DSN)
	require.Equal(t, testEnv["GDS_SENTRY_ENVIRONMENT"], conf.Sentry.Environment)
	require.Equal(t, true, conf.Sentry.TrackPerformance)
	require.Equal(t, testEnv["GDS_SENTRY_RELEASE"], conf.Sentry.Release)
	require.Equal(t, true, conf.Sentry.Debug)
	require.Equal(t, .2, conf.Sentry.SampleRate)
	require.True(t, conf.Secrets.Testing)
	require.True(t, conf.Activity.Enabled)
	require.Equal(t, testEnv["GDS_ACTIVITY_TOPIC"], conf.Activity.Topic)
	require.Equal(t, 10*time.Minute, conf.Activity.AggregationWindow)
	require.Equal(t, testEnv["GDS_ACTIVITY_ENSIGN_CLIENT_ID"], conf.Activity.Ensign.ClientID)
	require.Equal(t, testEnv["GDS_ACTIVITY_ENSIGN_CLIENT_SECRET"], conf.Activity.Ensign.ClientSecret)
	require.Equal(t, testEnv["GDS_ACTIVITY_ENSIGN_ENDPOINT"], conf.Activity.Ensign.Endpoint)
	require.Equal(t, testEnv["GDS_ACTIVITY_ENSIGN_AUTH_URL"], conf.Activity.Ensign.AuthURL)
	require.Equal(t, true, conf.Activity.Ensign.Insecure)
}

func TestAuthorizedDomainsPreprocessing(t *testing.T) {
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

	// Set authorized domains to require processing
	os.Setenv("GDS_ADMIN_OAUTH_AUTHORIZED_EMAIL_DOMAINS", "EXAMPLE.com, spacedout.io ,'quotes.org', 'Abadcombo.TECH")

	conf, err := config.New()
	require.NoError(t, err)

	require.Len(t, conf.Admin.Oauth.AuthorizedEmailDomains, 4)
	require.Equal(t, "example.com", conf.Admin.Oauth.AuthorizedEmailDomains[0])
	require.Equal(t, "spacedout.io", conf.Admin.Oauth.AuthorizedEmailDomains[1])
	require.Equal(t, "quotes.org", conf.Admin.Oauth.AuthorizedEmailDomains[2])
	require.Equal(t, "abadcombo.tech", conf.Admin.Oauth.AuthorizedEmailDomains[3])
}

func TestRequiredConfig(t *testing.T) {
	t.Skip("test assumes that confire is processing required tags recursively, is it?")
	required := []string{
		"GDS_DATABASE_URL",
		"GDS_SECRET_KEY",
		"GDS_ADMIN_OAUTH_GOOGLE_AUDIENCE",
		"GDS_ADMIN_TOKEN_KEYS",
		"GDS_ADMIN_OAUTH_AUTHORIZED_EMAIL_DOMAINS",
		"GDS_MEMBERS_CERTS",
		"GDS_MEMBERS_CERT_POOL",
		"GDS_DATABASE_CERT_PATH",
		"GDS_DATABASE_POOL_PATH",
		"SENDGRID_API_KEY",
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

	// Admin verification is predicated on it being enabled
	os.Setenv("GDS_ADMIN_ENABLED", "true")

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
	require.Equal(t, testEnv["GDS_DATABASE_URL"], conf.Database.URL)
	require.Equal(t, testEnv["GDS_SECRET_KEY"], conf.SecretKey)
	require.Equal(t, testEnv["GDS_ADMIN_OAUTH_GOOGLE_AUDIENCE"], conf.Admin.Oauth.GoogleAudience)
	require.Len(t, conf.Admin.TokenKeys, 2)
	require.Len(t, conf.Admin.Oauth.AuthorizedEmailDomains, 3)
}

func TestEmailConfigValidation(t *testing.T) {
	conf := config.EmailConfig{
		VerifyContactBaseURL: "http://localhost:3000/verify",
		AdminReviewBaseURL:   "http://localhost:3001/vasps",
	}

	err := conf.Validate()
	require.EqualError(t, err, "invalid configuration: admin review base URL must end in a /")

	conf.AdminReviewBaseURL += "/"
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: sendgrid api key and service email are required")

	conf.SendGridAPIKey = "supersecretapikey"
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: sendgrid api key and service email are required")

	conf.ServiceEmail = "service@example.com"
	conf.Storage = "fixtures/emails"
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: email archiving is only supported in testing mode")

	conf.Storage = ""
	err = conf.Validate()
	require.NoError(t, err, "expected valid configuration in non-testing mode")

	conf.Testing = true
	err = conf.Validate()
	require.NoError(t, err, "expected valid configuration in testing mode")
}

func TestAdminConfigValidation(t *testing.T) {
	conf := config.AdminConfig{
		Mode: "invalid",
	}
	require.EqualError(t, conf.Validate(), fmt.Sprintf("%q is not a valid gin mode", conf.Mode))

	conf = config.AdminConfig{
		Mode:    gin.ReleaseMode,
		Enabled: true,
		Oauth: config.OauthConfig{
			GoogleAudience:         "http://localhost",
			AuthorizedEmailDomains: []string{"example.com"},
		},
	}
	require.EqualError(t, conf.Validate(), "invalid configuration: token keys required for enabled admin")

	conf.TokenKeys = map[string]string{"keyid": "path/to/key.pem"}
	require.NoError(t, conf.Validate())
}

func TestOauthConfigValidation(t *testing.T) {
	conf := config.OauthConfig{
		AuthorizedEmailDomains: []string{"example.com"},
	}
	require.EqualError(t, conf.Validate(), "invalid configuration: oauth audience required for enabled admin")

	conf.GoogleAudience = "http://localhost"
	conf.AuthorizedEmailDomains = make([]string, 0)
	require.EqualError(t, conf.Validate(), "invalid configuration: authorized email domains required for enabled admin")

	conf.AuthorizedEmailDomains = []string{"example.com"}
	require.NoError(t, conf.Validate())
}

func TestMembersConfigValidation(t *testing.T) {
	conf := config.MembersConfig{
		Insecure: true,
		Certs:    "",
		CertPool: "",
	}

	// If Insecure is set to true, certs and cert pool are not required.
	err := conf.Validate()
	require.NoError(t, err)

	// If Insecure is false, then the certs and cert pool are required.
	conf.Insecure = false
	err = conf.Validate()
	require.EqualError(t, err, "invalid configuration: serving mTLS requires the path to certs and the cert pool")
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
