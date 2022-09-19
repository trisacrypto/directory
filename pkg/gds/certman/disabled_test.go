package certman_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/certman"
	"github.com/trisacrypto/directory/pkg/gds/config"
)

func TestDisabled(t *testing.T) {
	// Should be able to create a new disabled certman with no dependencies
	service, err := certman.New(config.CertManConfig{Enabled: false}, nil, nil, nil)
	require.NoError(t, err, "could not create a disabled certman service with no config")
	require.IsType(t, &certman.Disabled{}, service, "expected the service to be disabled")

	// Should be able to run and shutdown the service without blocking
	var wg sync.WaitGroup
	err = service.Run(&wg)
	require.NoError(t, err, "could not run the disabled service")

	// If the test times out it means that disabled isn't properly managing the wait group
	require.NotPanics(t, service.Stop, "stop should not panic")
	wg.Wait()

	require.NotPanics(t, service.CertManager, "cert manager should not panic")
	require.NotPanics(t, service.HandleCertificateRequests, "handle certificate requests should not panic")
}
