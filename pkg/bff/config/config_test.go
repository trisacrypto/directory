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
	"GDS_BFF_MAINTENANCE":      "false",
	"GDS_BFF_BIND_ADDR":        "8080",
	"GDS_BFF_MODE":             "debug",
	"GDS_BFF_LOG_LEVEL":        "debug",
	"GDS_BFF_CONSOLE_LOG":      "true",
	"GDS_BFF_TESTNET_INSECURE": "true",
	"GDS_BFF_TESTNET_ENDPOINT": "localhost:8443",
	"GDS_BFF_TESTNET_TIMEOUT":  "5s",
	"GDS_BFF_MAINNET_INSECURE": "true",
	"GDS_BFF_MAINNET_ENDPOINT": "localhost:8444",
	"GDS_BFF_MAINNET_TIMEOUT":  "3s",
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
	require.True(t, conf.TestNet.Insecure)
	require.Equal(t, testEnv["GDS_BFF_TESTNET_ENDPOINT"], conf.TestNet.Endpoint)
	require.Equal(t, 5*time.Second, conf.TestNet.Timeout)
	require.True(t, conf.MainNet.Insecure)
	require.Equal(t, testEnv["GDS_BFF_MAINNET_ENDPOINT"], conf.MainNet.Endpoint)
	require.Equal(t, 3*time.Second, conf.MainNet.Timeout)
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
