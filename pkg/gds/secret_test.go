package gds_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/config"
)

// Three environment variables are required to run this test;
// The first can be included when the tests are run, e.g.
// GDS_TEST_GOOGLE_SECRETS=1 go test ./pkg/gds/secret_test.go
// The others, $GOOGLE_APPLICATION_CREDENTIALS and $GOOGLE_PROJECT_NAME
// are both required to be a valid google secret manager service
// credentials JSON (absolute path), and a valid google secret manager
// project name.
// Note: tests execute against live secret manager API, so use caution!
func TestSecrets(t *testing.T) {
	if os.Getenv("GDS_TEST_GOOGLE_SECRETS") == "" {
		t.Skip("skip Google SecretManager API connection test")
	}
	testConf := config.SecretsConfig{
		Credentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		Project:     os.Getenv("GOOGLE_PROJECT_NAME"),
	}

	// create secret manager
	sm, err := gds.NewSecretManager(testConf, "foo")
	require.NoError(t, err)

	// create a secret
	require.NoError(t, sm.CreateSecret(context.Background(), "test"))

	// create a version
	require.NoError(t, sm.AddSecretVersion(context.Background(), "test", []byte("test payload")))

	// access version
	result, err := sm.GetLatestVersion(context.Background(), "test")
	require.Equal(t, result, []byte("test payload"))
	require.NoError(t, err)

	// delete the secret
	require.NoError(t, sm.DeleteSecret(context.Background(), "test"))

}
