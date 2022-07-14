package mock

import (
	"net/http/httptest"

	"github.com/trisacrypto/directory/pkg/bff/admin"
	apiv2 "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

// New creates a new mock admin client with an embedded token manager and httptest
// client for testing.
func New(tm *tokens.TokenManager, srv *httptest.Server) (_ apiv2.DirectoryAdministrationClient, err error) {
	var creds apiv2.Credentials
	if creds, err = admin.NewCredentialsFromTokens(tm); err != nil {
		return nil, err
	}

	return apiv2.New(srv.URL, creds, apiv2.WithClient(srv.Client()))
}
