package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/trtl/config"
)

var testEnv = map[string]string{
	"TRTL_MAINTENANCE":              "true",
	"TRTL_BIND_ADDR":                ":445",
	"TRTL_LOG_LEVEL":                "debug",
	"TRTL_CONSOLE_LOG":              "true",
	"TRTL_DATABASE_URL":             "leveldb:///fixtures/db",
	"TRTL_DATABASE_REINDEX_ON_BOOT": "true",
	"TRTL_REPLICA_ENABLED":          "true",
	"TRTL_REPLICA_PID":              "8",
	"TRTL_REPLICA_NAME":             "mitchell",
	"TRTL_REPLICA_REGION":           "us-east-1c",
	"TRTL_REPLICA_GOSSIP_INTERVAL":  "30m",
	"TRTL_REPLICA_GOSSIP_SIGMA":     "3m",
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

	// Load configuration from the environment
	conf, err := config.New()
	require.NoError(t, err)

	// Test configuration set from the environment
	require.True(t, conf.Maintenance)
	require.Equal(t, testEnv["TRTL_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.True(t, conf.ConsoleLog)
	require.Equal(t, testEnv["TRTL_DATABASE_URL"], conf.Database.URL)
	require.True(t, conf.Database.ReindexOnBoot)
	require.True(t, conf.Replica.Enabled)
	require.Equal(t, uint64(8), conf.Replica.PID)
	require.Equal(t, testEnv["TRTL_REPLICA_NAME"], conf.Replica.Name)
	require.Equal(t, testEnv["TRTL_REPLICA_REGION"], conf.Replica.Region)
	require.Equal(t, 30*time.Minute, conf.Replica.GossipInterval)
	require.Equal(t, 3*time.Minute, conf.Replica.GossipSigma)
}

func TestRequiredConfig(t *testing.T) {
	required := []string{
		"TRTL_DATABASE_URL",
		"TRTL_REPLICA_PID",
		"TRTL_REPLICA_REGION",
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

	// Replica verification is predicated on its being enabled
	os.Setenv("GDS_REPLICA_ENABLED", "true")

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
	require.Equal(t, testEnv["TRTL_DATABASE_URL"], conf.Database.URL)
	require.True(t, conf.Replica.Enabled)
	require.Equal(t, uint64(8), conf.Replica.PID)
	require.Equal(t, testEnv["TRTL_REPLICA_REGION"], conf.Replica.Region)
}

func TestValidateReplicaConfig(t *testing.T) {
	// Replica config should only be validated when it is enabled
	conf := &config.ReplicaConfig{}
	require.NoError(t, conf.Validate())

	conf.Enabled = true
	require.EqualError(t, conf.Validate(), "invalid configuration: PID required for enabled replica")

	// PID should be required
	conf.PID = 32
	require.EqualError(t, conf.Validate(), "invalid configuration: region required for enabled replica")

	// Region should be required
	conf.Region = "raleigh"
	require.EqualError(t, conf.Validate(), "invalid configuration: specify non-zero gossip interval and sigma")

	// Gossip Interval should be required
	conf.GossipInterval = time.Second * 30
	require.EqualError(t, conf.Validate(), "invalid configuration: specify non-zero gossip interval and sigma")

	// Gossip Sigma should be required
	conf.GossipSigma = time.Millisecond * 500
	require.NoError(t, conf.Validate())
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
