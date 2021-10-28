package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/trtl/config"
)

var testEnv = map[string]string{
	"TRTL_ENABLED":                 "true",
	"TRTL_REPLICA_BIND_ADDR":       ":445",
	"TRTL_REPLICA_PID":             "8",
	"TRTL_REPLICA_NAME":            "mitchell",
	"TRTL_REPLICA_REGION":          "us-east-1c",
	"TRTL_REPLICA_GOSSIP_INTERVAL": "30m",
	"TRTL_REPLICA_GOSSIP_SIGMA":    "3m",
}

func TestConfig(t *testing.T) {
	// TODO: Fix this test
	t.Skip("TestConfig is erroring: PID required for enabled replica")

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
	require.Equal(t, true, conf.Enabled)
	require.Equal(t, testEnv["TRTL_REPLICA_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, uint64(8), conf.PID)
	require.Equal(t, testEnv["TRTL_REPLICA_NAME"], conf.Name)
	require.Equal(t, testEnv["TRTL_REPLICA_REGION"], conf.Region)
	require.Equal(t, 30*time.Minute, conf.GossipInterval)
	require.Equal(t, 3*time.Minute, conf.GossipSigma)
}

func TestRequiredConfig(t *testing.T) {
	// TODO: Fix this test
	t.Skip("TestRequiredConfig is erroring: PID required for enabled replica")

	required := []string{
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
	require.True(t, conf.Enabled)
	require.Equal(t, uint64(8), conf.PID)
	require.Equal(t, testEnv["TRTL_REPLICA_REGION"], conf.Region)

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
