package config_test

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/trtl/config"
)

func TestReplicaConfigureNoModification(t *testing.T) {
	s1 := makeRandomTestStrategy()
	s2 := makeRandomTestStrategy()

	conf := config.ReplicaConfig{
		Enabled:        true,
		PID:            53,
		Region:         "ibelin",
		Name:           "balian",
		GossipInterval: 10 * time.Minute,
		GossipSigma:    1500 * time.Second,
	}

	other, err := conf.Configure(s1, s2)
	require.NoError(t, err, "test strategies returned an error")
	require.NotEqual(t, conf, other, "other should be modified")

	// Assert the original configuration has not been modified
	require.True(t, conf.Enabled)
	require.Equal(t, uint64(53), conf.PID)
	require.Equal(t, "ibelin", conf.Region)
	require.Equal(t, "balian", conf.Name)
	require.Equal(t, 10*time.Minute, conf.GossipInterval)
	require.Equal(t, 1500*time.Second, conf.GossipSigma)
}

func TestReplicaConfigurationErrors(t *testing.T) {
	s1 := makeErroringStrategy()   // should generate an error
	s2 := makeErrorBreakStrategy() // should break
	s3 := makeRandomTestStrategy() // should never get to this

	conf := config.ReplicaConfig{
		Enabled:        true,
		PID:            53,
		Region:         "ibelin",
		Name:           "balian",
		GossipInterval: 10 * time.Minute,
		GossipSigma:    1500 * time.Second,
	}

	// Assert that s1 returns an error
	_, err := conf.Configure(s1)
	require.Error(t, err, "s1 should return an error")

	// Assert a strategy break error stops processing
	other, err := conf.Configure(s1, s2, s3)
	require.NoError(t, err, "expected nil error for StrategyBreak")

	// Assert the returned configuration didn't make it to s3
	require.True(t, other.Enabled)
	require.Equal(t, uint64(53), other.PID)
	require.Equal(t, "ibelin", other.Region)
	require.Equal(t, "balian", other.Name)
	require.Equal(t, 10*time.Minute, other.GossipInterval)
	require.Equal(t, 1500*time.Second, other.GossipSigma)

	// Assert that errors are passed through
	_, err = conf.Configure(s1, s3)
	require.Error(t, err, "error was not passed through")
}

func TestMultiErrorStrategy(t *testing.T) {
	s1 := config.MultiError(makeErroringStrategy())
	s2 := config.MultiError(makeRandomTestStrategy())
	s3 := config.MultiError(makeErroringStrategy())

	conf := config.ReplicaConfig{
		Enabled:        true,
		PID:            53,
		Region:         "ibelin",
		Name:           "balian",
		GossipInterval: 10 * time.Minute,
		GossipSigma:    1500 * time.Second,
	}

	_, err := conf.Configure(s1, s2, s3)
	require.Error(t, err, "a multi error should have been returned")
	require.IsType(t, &multierror.Error{}, err)

	merr := err.(*multierror.Error)
	require.Equal(t, 2, merr.Len())
}

func TestHostnamePID(t *testing.T) {
	// NOTE: cannot actually test os.Hostname() since this is platform dependent, this
	// test primarily tests different hostnames to ensure they are parsed correctly.
	tc := []struct {
		hostname string
		name     string
		pid      uint64
		err      bool
	}{
		{"trtl", "", 0, true},
		{"ninja trtl-23", "", 0, true},
		{"ninja_trtl_212", "", 0, true},
		{"trtl-23", "trtl-63", 63, false},
		{"ninja-trtl-0", "ninja-trtl-40", 40, false},
		{"ninja_trtl-13920", "ninja_trtl-13960", 13960, false},
		{"trtl.us.testnet.directory-18", "trtl.us.testnet.directory-58", 58, false},
	}

	conf := config.ReplicaConfig{
		Enabled:        true,
		PID:            40,
		Region:         "ibelin",
		Name:           "balian",
		GossipInterval: 10 * time.Minute,
		GossipSigma:    1500 * time.Second,
	}

	for i, tt := range tc {
		other, err := conf.Configure(config.HostnamePID(tt.hostname))
		if tt.err {
			require.Error(t, err, "test case %d expected an error, none was returned", i)
		} else {
			require.NoError(t, err, "test case %d errored unexpectedly", i)
			require.Equal(t, tt.name, other.Name, "test case %d did not parse correctly", i)
			require.Equal(t, tt.pid, other.PID, "test case %d did not parse correctly", i)
		}
	}
}

func TestFilePID(t *testing.T) {
	conf := config.ReplicaConfig{
		Enabled:        true,
		PID:            40,
		Region:         "ibelin",
		Name:           "balian",
		GossipInterval: 10 * time.Minute,
		GossipSigma:    1500 * time.Second,
	}

	other, err := conf.Configure(config.FilePID("testdata/test.pid"))
	require.NoError(t, err, "could not read pid file")
	require.Equal(t, uint64(158), other.PID)
	require.Equal(t, "balian-158", other.Name)

	_, err = conf.Configure(config.FilePID("testdata/bad.pid"))
	require.Error(t, err, "bad pid did not error")
}

func TestJSONConfig(t *testing.T) {
	conf := config.ReplicaConfig{
		Enabled:        true,
		PID:            48,
		Region:         "ibelin",
		Name:           "balian",
		GossipInterval: 10 * time.Minute,
		GossipSigma:    1500 * time.Second,
	}

	// Ensure JSON config does not work if an error happens ahead of it
	other, err := conf.Configure(makeErroringStrategy(), config.JSONConfig("testdata/replicas.json"))
	require.Error(t, err)
	require.Equal(t, "balian", other.Name)

	// Ensure JSON config skips without error if the file doesn't exist
	other, err = conf.Configure(config.JSONConfig("testdata/missing.json"))
	require.NoError(t, err, "error occurred instead of skipping")
	require.Equal(t, "balian", other.Name)

	// Ensure JSON config processes real data and processing doesn't continue
	other, err = conf.Configure(config.JSONConfig("testdata/replicas.json"), makeRandomTestStrategy())
	require.NoError(t, err, "could not process json config")
	require.False(t, other.Enabled)
	require.Equal(t, uint64(48), other.PID)
	require.Equal(t, "queens", other.Region)
	require.Equal(t, "leonardo", other.Name)
	require.Equal(t, 18*time.Minute, other.GossipInterval)
	require.Equal(t, 3*time.Minute, other.GossipSigma)

}

func makeRandomTestStrategy() config.ReplicaStrategy {
	return func(in config.ReplicaConfig, err error) (config.ReplicaConfig, error) {
		randval := rand.Int63n(10000) + 1
		in.Enabled = randval%2 == 0
		in.Name = fmt.Sprintf("testing-%05X", randval)
		in.Region = fmt.Sprintf("place-%05X", randval/2)
		in.GossipInterval = time.Duration(rand.Int63n(int64(24 * time.Hour)))
		in.GossipSigma = time.Duration(rand.Int63n(int64(4 * time.Hour)))
		return in, err
	}
}

func makeErrorBreakStrategy() config.ReplicaStrategy {
	return func(in config.ReplicaConfig, err error) (config.ReplicaConfig, error) {
		if err != nil {
			return in, config.ErrStrategyBreak
		}
		return in, err
	}
}

func makeErroringStrategy() config.ReplicaStrategy {
	return func(in config.ReplicaConfig, err error) (config.ReplicaConfig, error) {
		return in, errors.New("something bad happened")
	}
}
