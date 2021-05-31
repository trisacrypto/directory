package gds_test

import (
	"context"
	"crypto/x509"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/trisa/pkg/trust"
)

// Six environment variables are required to run this test:
// $GDS_TEST_GOOGLE_SECRETS
//     Required; can be included when the tests are run, e.g.
//     GDS_TEST_GOOGLE_SECRETS=1 go test ./pkg/gds/secret_test.go
// $GOOGLE_APPLICATION_CREDENTIALS
//     Required; must be a valid google secret manager service
//     credentials JSON (absolute path)
// $GOOGLE_PROJECT_NAME
//     Required; must be a valid google secret manager project name
// $TRISA_TEST_PASSWORD
//     Required; must be a valid pkcs12password that matches the TRISA_TEST_CERT
// $TRISA_TEST_CERT
//     Required; must be a valid TRISA certificate that matches the TRISA_TEST_PASSWORD
// $TRISA_TEST_FILE
//     Required; absolute path for an intermediate tempfile to write the retrieved cert
//     TODO: Temp file write & delete can be removed once trust serializer can unzip raw bytes, see:
//     https://github.com/trisacrypto/trisa/issues/51
// Note: tests execute against live secret manager API, so use caution!
func TestSecrets(t *testing.T) {
	if os.Getenv("GDS_TEST_GOOGLE_SECRETS") == "" {
		t.Skip("skip Google SecretManager API connection test")
	}
	testConf := config.SecretsConfig{
		Credentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		Project:     os.Getenv("GOOGLE_PROJECT_NAME"),
	}

	// Prep test fixtures
	testRequestId := "123"
	testSecretType := "testCert"
	testContext := context.Background()
	tempFile := os.Getenv("TRISA_TEST_FILE")
	testPassword := os.Getenv("TRISA_TEST_PASSWORD")
	testPayload, err := ioutil.ReadFile(os.Getenv("TRISA_TEST_CERT"))
	require.NoError(t, err)

	// Create secret manager
	sm, err := gds.NewSecretManager(testConf)
	require.NoError(t, err)

	// Create a secret for the testSecretType
	require.NoError(t, sm.With(testRequestId).CreateSecret(testContext, testSecretType))

	// Create a version to hold the testPayload
	require.NoError(t, sm.With(testRequestId).AddSecretVersion(testContext, testSecretType, testPayload))

	// Retrieve the testPayload
	testResult, err := sm.With(testRequestId).GetLatestVersion(testContext, testSecretType)
	require.Equal(t, testResult, testPayload)
	require.NoError(t, err)

	// Create a serializer to read in the retrieved payload
	var archive *trust.Serializer
	archive, err = trust.NewSerializer(true, testPassword, trust.CompressionZIP)
	require.NoError(t, err)

	// Write the cert to a temp file
	require.NoError(t, ioutil.WriteFile(tempFile, testResult, 0777))

	// Create provider to read in the bytes of the zipfile
	var provider *trust.Provider
	provider, err = archive.ReadFile(tempFile)
	require.NoError(t, err)

	// Delete the temporary file now that we're done with it
	require.NoError(t, os.Remove(tempFile))

	// Verify that the leaves of the retrieved cert can be extracted
	var cert *x509.Certificate
	cert, err = provider.GetLeafCertificate()
	require.NoError(t, err)
	require.NotNil(t, cert)

	// Delete the secret
	require.NoError(t, sm.With(testRequestId).DeleteSecret(testContext, testSecretType))
}
