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
	"GDS_BIND_ADDR":        ":443",
	"GDS_DATABASE_URL":     "fixtures/db",
	"SECTIGO_USERNAME":     "foo",
	"SECTIGO_PASSWORD":     "supersecret",
	"SENDGRID_API_KEY":     "bar1234",
	"GDS_SERVICE_EMAIL":    "test@example.com",
	"GDS_ADMIN_EMAIL":      "admin@example.com",
	"GDS_LOG_LEVEL":        "debug",
	"GDS_DIRECTORY_ID":     "testdirectory.org",
	"GDS_SECRET_KEY":       "theeaglefliesatmidnight",
	"GDS_CERTMAN_INTERVAL": "60s",
	"GDS_CERTMAN_STORAGE":  "fixtures/certs",
	"GDS_BACKUP_ENABLED":   "true",
	"GDS_BACKUP_INTERVAL":  "36h",
	"GDS_BACKUP_STORAGE":   "fixtures/backups",
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
	require.Equal(t, testEnv["GDS_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, testEnv["GDS_DATABASE_URL"], conf.DatabaseURL)
	require.Equal(t, testEnv["SECTIGO_USERNAME"], conf.Sectigo.Username)
	require.Equal(t, testEnv["SECTIGO_PASSWORD"], conf.Sectigo.Password)
	require.Equal(t, testEnv["SENDGRID_API_KEY"], conf.SendGridAPIKey)
	require.Equal(t, testEnv["GDS_SERVICE_EMAIL"], conf.ServiceEmail)
	require.Equal(t, testEnv["GDS_ADMIN_EMAIL"], conf.AdminEmail)
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.Equal(t, testEnv["GDS_DIRECTORY_ID"], conf.DirectoryID)
	require.Equal(t, testEnv["GDS_SECRET_KEY"], conf.SecretKey)
	require.Equal(t, 1*time.Minute, conf.CertMan.Interval)
	require.Equal(t, testEnv["GDS_CERTMAN_STORAGE"], conf.CertMan.Storage)
	require.Equal(t, true, conf.Backup.Enabled)
	require.Equal(t, 36*time.Hour, conf.Backup.Interval)
	require.Equal(t, testEnv["GDS_BACKUP_STORAGE"], conf.Backup.Storage)
}

func TestRequiredConfig(t *testing.T) {
	// Set required environment variables and cleanup after
	prevEnv := curEnv("GDS_DATABASE_URL", "GDS_SECRET_KEY")
	t.Cleanup(func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	})

	_, err := config.New()
	require.Error(t, err)
	setEnv("GDS_DATABASE_URL", "GDS_SECRET_KEY")

	conf, err := config.New()
	require.NoError(t, err)

	// Test required configuration
	require.Equal(t, testEnv["GDS_DATABASE_URL"], conf.DatabaseURL)
	require.Equal(t, testEnv["GDS_SECRET_KEY"], conf.SecretKey)
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
