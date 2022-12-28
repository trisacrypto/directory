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
	"TRTL_METRICS_ADDR":             ":9090",
	"TRTL_METRICS_ENABLED":          "true",
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
	"TRTL_INSECURE":                 "true",
	"TRTL_MTLS_CHAIN_PATH":          "fixtures/certs/chain.pem",
	"TRTL_MTLS_CERT_PATH":           "fixtures/certs/cert.pem",
	"TRTL_BACKUP_ENABLED":           "true",
	"TRTL_BACKUP_INTERVAL":          "1h",
	"TRTL_BACKUP_STORAGE":           "fixtures/backups",
	"TRTL_BACKUP_KEEP":              "7",
	"TRTL_SENTRY_DSN":               "https://something.ingest.sentry.io",
	"TRTL_SENTRY_ENVIRONMENT":       "test",
	"TRTL_SENTRY_RELEASE":           "1.4",
	"TRTL_SENTRY_DEBUG":             "true",
	"TRTL_SENTRY_TRACK_PERFORMANCE": "true",
	"TRTL_SENTRY_SAMPLE_RATE":       "0.2",
}

var strategyEnv = map[string]string{
	"TRTL_REPLICA_STRATEGY_HOSTNAME_PID": "true",
	"TRTL_REPLICA_HOSTNAME":              "pizza-36",
	"TRTL_REPLICA_STRATEGY_FILE_PID":     "",
	"TRTL_REPLICA_STRATEGY_JSON_CONFIG":  "testdata/replicas.json",
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
	require.Equal(t, testEnv["TRTL_METRICS_ADDR"], conf.Metrics.Addr)
	require.True(t, conf.Metrics.Enabled)
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
	require.True(t, conf.MTLS.Insecure)
	require.Equal(t, testEnv["TRTL_MTLS_CHAIN_PATH"], conf.MTLS.ChainPath)
	require.Equal(t, testEnv["TRTL_MTLS_CERT_PATH"], conf.MTLS.CertPath)
	require.True(t, conf.Backup.Enabled)
	require.Equal(t, 1*time.Hour, conf.Backup.Interval)
	require.Equal(t, testEnv["TRTL_BACKUP_STORAGE"], conf.Backup.Storage)
	require.Equal(t, 7, conf.Backup.Keep)
	require.Equal(t, testEnv["TRTL_SENTRY_DSN"], conf.Sentry.DSN)
	require.Equal(t, testEnv["TRTL_SENTRY_ENVIRONMENT"], conf.Sentry.Environment)
	require.Equal(t, testEnv["TRTL_SENTRY_RELEASE"], conf.Sentry.Release)
	require.Equal(t, true, conf.Sentry.Debug)
	require.Equal(t, true, conf.Sentry.TrackPerformance)
	require.Equal(t, 0.2, conf.Sentry.SampleRate)
}

func TestRequiredConfig(t *testing.T) {
	required := []string{
		"TRTL_DATABASE_URL",
		"TRTL_REPLICA_PID",
		"TRTL_REPLICA_REGION",
		"TRTL_MTLS_CHAIN_PATH",
		"TRTL_MTLS_CERT_PATH",
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

	// Replica verification is predicated on it being enabled
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

func TestValidateMTLSConfig(t *testing.T) {
	// MTLS config should only be validated when insecure is false
	conf := &config.MTLSConfig{
		Insecure: true,
	}
	require.NoError(t, conf.Validate())

	// Both ChainPath and CertPath are required
	conf.Insecure = false
	require.Error(t, conf.Validate())

	conf.ChainPath = "/path/to/chain"
	require.Error(t, conf.Validate())

	conf = &config.MTLSConfig{
		CertPath: "/path/to/cert",
	}
	require.Error(t, conf.Validate())

	conf = &config.MTLSConfig{
		ChainPath: "/path/to/chain",
		CertPath:  "/path/to/cert",
	}
	require.NoError(t, conf.Validate())
}

func TestKubernetesStatefulSetStrategy(t *testing.T) {
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

	// Set the replica strategy environment variables
	for key, val := range strategyEnv {
		os.Setenv(key, val)
	}

	// Load configuration from the environment
	conf, err := config.New()
	require.NoError(t, err)

	// Check that the configuration was loaded from the strategy
	require.True(t, conf.Replica.Enabled)
	require.Equal(t, uint64(44), conf.Replica.PID)
	require.Equal(t, "brooklyn", conf.Replica.Region)
	require.Equal(t, "donatello", conf.Replica.Name)
	require.Equal(t, 21*time.Minute, conf.Replica.GossipInterval)
	require.Equal(t, 1500*time.Millisecond, conf.Replica.GossipSigma)
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
