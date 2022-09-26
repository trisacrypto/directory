package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/store/config"
)

func TestConfigValidation(t *testing.T) {
	conf := config.StoreConfig{
		URL:           "leveldb:///db",
		ReindexOnBoot: false,
		Insecure:      false,
		CertPath:      "",
		PoolPath:      "",
	}

	// Connecting to leveldb should not require extra validation
	err := conf.Validate()
	require.NoError(t, err, "could not validate leveldb configuration")

	// If Insecure is set to true, certs and cert pool are not required.
	conf.URL = "trtl://trtl.test.net:443"
	conf.Insecure = true
	err = conf.Validate()
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
