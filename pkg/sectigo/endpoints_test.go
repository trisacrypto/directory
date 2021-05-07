package sectigo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUrlFor(t *testing.T) {
	// No params
	require.Equal(t, "https://iot.sectigo.com/api/v1/certificates/find", urlFor(findCertificateEP))

	// Test params
	require.Equal(t, "https://iot.sectigo.com/api/v1/organizations/42/authority/24", urlFor(authorityDetailEP, 42, 24))

	// Test copy
	require.Equal(t, "https://iot.sectigo.com/api/v1/organizations/95/authority/14", urlFor(authorityDetailEP, 95, 14))

	// Test panic with unknown endpoint
	require.Panics(t, func() { urlFor("foo") })
}
