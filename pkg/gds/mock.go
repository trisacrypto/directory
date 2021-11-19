package gds

import (
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

// NewMock creates and returns a mocked Service for testing, using values provided in
// the config.
func NewMock(conf config.Config) (s *Service, err error) {
	svc := &Service{
		conf: conf,
	}
	if svc.email, err = emails.New(conf.Email); err != nil {
		return nil, err
	}
	if svc.db, err = store.Open(conf.Database); err != nil {
		return nil, err
	}

	admin := &Admin{
		svc:  svc,
		conf: &svc.conf.Admin,
		db:   svc.db,
	}
	if admin.tokens, err = tokens.MockTokenManager(); err != nil {
		return nil, err
	}
	svc.admin = admin
	return svc, nil
}

// NewMockedAdmin creates and returns a mocked Admin for testing, using values provided
// in the config.
func NewMockedAdmin(conf config.Config) (admin *Admin, db store.Store, tokens *tokens.TokenManager, err error) {
	var svc *Service
	if svc, err = NewMock(conf); err != nil {
		return nil, nil, nil, err
	}
	return svc.admin, svc.admin.db, svc.admin.tokens, nil
}
