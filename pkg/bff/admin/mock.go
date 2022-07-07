package admin

import (
	"net/http/httptest"

	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

// NewMock creates a new mock admin client with an embedded token manager and httptest
// client for testing.
func NewMock(tm *tokens.TokenManager, srv *httptest.Server) (admin.DirectoryAdministrationClient, error) {
	creds := &Credentials{
		tm: tm,
	}

	return admin.New(srv.URL, creds, srv.Client())
}
