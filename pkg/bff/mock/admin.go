package mock

import (
	"net/http/httptest"

	"github.com/trisacrypto/directory/pkg/bff/admin"
	"github.com/trisacrypto/directory/pkg/gds"
	apiv2 "github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"google.golang.org/grpc"
)

func NewAdmin(trtlConn *grpc.ClientConn) (a *Admin, err error) {
	a = &Admin{}

	serviceConfig := gds.MockConfig()
	serviceConfig.Members.Enabled = false
	var svc *gds.Service
	if svc, err = gds.NewMock(serviceConfig, trtlConn); err != nil {
		return nil, err
	}

	a.srv = httptest.NewTLSServer(svc.GetAdmin().GetRouter())
	if a.client, err = admin.NewMock(svc.GetAdmin().GetTokenManager(), a.srv); err != nil {
		return nil, err
	}

	return a, nil
}

type Admin struct {
	srv    *httptest.Server
	client apiv2.DirectoryAdministrationClient
}

func (a *Admin) Shutdown() {
	a.srv.Close()
}

func (a *Admin) Client() apiv2.DirectoryAdministrationClient {
	return a.client
}
