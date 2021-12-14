package config_test

import (
	"os"
	"testing"
	"time"

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
	"GDS_DATABASE_URL":                         "fixtures/db",
	"GDS_DATABASE_REINDEX_ON_BOOT":             "false",
	"SECTIGO_USERNAME":                         "foo",
	"SECTIGO_PASSWORD":                         "supersecret",
	"SECTIGO_PROFILE":                          "17",
	"GDS_SERVICE_EMAIL":                        "test@example.com",
	"GDS_ADMIN_EMAIL":                          "admin@example.com",
	"SENDGRID_API_KEY":                         "bar1234",
	"GDS_VERIFY_CONTACT_URL":                   "http://localhost:3000/verify-contact",
	"GDS_ADMIN_REVIEW_URL":                     "http://localhost:3001/vasps/",
	"GDS_EMAIL_TESTING":                        "true",
	"GDS_CERTMAN_INTERVAL":                     "60s",
	"GDS_CERTMAN_STORAGE":                      "fixtures/certs",
	"GDS_BACKUP_ENABLED":                       "true",
	"GDS_BACKUP_INTERVAL":                      "36h",
	"GDS_BACKUP_STORAGE":                       "fixtures/backups",
	"GDS_BACKUP_KEEP":                          "7",
	"GOOGLE_APPLICATION_CREDENTIALS":           "test.json",
	"GOOGLE_PROJECT_NAME":                      "test",
	"GDS_SECRETS_TESTING":                      "true",
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
	require.Equal(t, testEnv["SECTIGO_USERNAME"], conf.Sectigo.Username)
	require.Equal(t, testEnv["SECTIGO_PASSWORD"], conf.Sectigo.Password)
	require.Equal(t, testEnv["SECTIGO_PROFILE"], conf.Sectigo.Profile)
	require.Equal(t, testEnv["GDS_SERVICE_EMAIL"], conf.Email.ServiceEmail)
	require.Equal(t, testEnv["GDS_ADMIN_EMAIL"], conf.Email.AdminEmail)
	require.Equal(t, testEnv["SENDGRID_API_KEY"], conf.Email.SendGridAPIKey)
	require.Equal(t, testEnv["GDS_VERIFY_CONTACT_URL"], conf.Email.VerifyContactBaseURL)
	require.Equal(t, testEnv["GDS_ADMIN_REVIEW_URL"], conf.Email.AdminReviewBaseURL)
	require.True(t, conf.Email.Testing)
	require.Equal(t, testEnv["GDS_DIRECTORY_ID"], conf.Email.DirectoryID)
	require.Equal(t, 1*time.Minute, conf.CertMan.Interval)
	require.Equal(t, testEnv["GDS_CERTMAN_STORAGE"], conf.CertMan.Storage)
	require.Equal(t, true, conf.Backup.Enabled)
	require.Equal(t, 36*time.Hour, conf.Backup.Interval)
	require.Equal(t, testEnv["GDS_BACKUP_STORAGE"], conf.Backup.Storage)
	require.Equal(t, 7, conf.Backup.Keep)
	require.Equal(t, testEnv["GOOGLE_APPLICATION_CREDENTIALS"], conf.Secrets.Credentials)
	require.Equal(t, testEnv["GOOGLE_PROJECT_NAME"], conf.Secrets.Project)
	require.True(t, conf.Secrets.Testing)
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
	required := []string{
		"GDS_DATABASE_URL",
		"GDS_SECRET_KEY",
		"GDS_ADMIN_OAUTH_GOOGLE_AUDIENCE",
		"GDS_ADMIN_TOKEN_KEYS",
		"GDS_ADMIN_OAUTH_AUTHORIZED_EMAIL_DOMAINS",
		"GDS_MEMBERS_CERTS",
		"GDS_MEMBERS_CERT_POOL",
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
		VerifyContactBaseURL: "http://localhost:3000/verify-contact",
		AdminReviewBaseURL:   "http://localhost:3001/vasps",
	}

	err := conf.Validate()
	require.EqualError(t, err, "admin review base URL must end in a /")

	conf.AdminReviewBaseURL += "/"
	err = conf.Validate()
	require.NoError(t, err)
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
