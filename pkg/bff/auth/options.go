package auth

import (
	"net/http"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
)

// WithHTTPClient configures the authentication provider to use the specified client.
// This is used in tests to configure the client to use a localhost TLS httptest server.
// This option should NOT be used in production.
//
// NOTE: this has been added to the jwks code but not tagged yet. Once the library gets
// updated we can remove this function and use their implementation.
// https://github.com/auth0/go-jwt-middleware/blob/master/jwks/provider.go#L55
func WithHTTPClient(client *http.Client) jwks.ProviderOption {
	return func(p *jwks.Provider) {
		p.Client = client
	}
}
