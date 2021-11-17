package sectigo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/mock"
)

func TestEndpoint(t *testing.T) {
	// No params
	ep, err := Endpoint(FindCertificateEP)
	require.NoError(t, err)
	require.Equal(t, "https://iot.sectigo.com/api/v1/certificates/find", ep.String())

	// Test params
	ep, err = Endpoint(AuthorityDetailEP, 42, 24)
	require.NoError(t, err)
	require.Equal(t, "https://iot.sectigo.com/api/v1/organizations/42/authority/24", ep.String())

	// Test copy
	ep, err = Endpoint(AuthorityDetailEP, 95, 14)
	require.NoError(t, err)
	require.Equal(t, "https://iot.sectigo.com/api/v1/organizations/95/authority/14", ep.String())

	// Test panic with unknown endpoint'
	ep, err = Endpoint("foo")
	require.EqualError(t, err, `no endpoint named "foo"`)
	require.Nil(t, ep)
}

func TestModifyBaseURL(t *testing.T) {
	defer ResetBaseURL() // Ensure Base URL is always reset

	// Check original url
	ep, err := Endpoint(FindCertificateEP)
	require.NoError(t, err)
	require.Equal(t, "https://iot.sectigo.com/api/v1/certificates/find", ep.String())

	// Create a mock Server
	m, err := mock.New()
	require.NoError(t, err)
	defer m.Close()

	// Check that mock Server updated the sectigo URL
	ep, err = Endpoint(FindCertificateEP)
	require.NoError(t, err)
	require.Equal(t, "http", ep.Scheme)
	require.Contains(t, ep.Host, "127.0.0.1")
	require.Equal(t, "/api/v1/certificates/find", ep.Path)

	// Reset Base URL
	ResetBaseURL()
	ep, err = Endpoint(FindCertificateEP)
	require.NoError(t, err)
	require.Equal(t, "https://iot.sectigo.com/api/v1/certificates/find", ep.String())
}
