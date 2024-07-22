package mock

import (
	"context"
	"os"

	"github.com/trisacrypto/directory/pkg/store"
	trtlstore "github.com/trisacrypto/directory/pkg/store/trtl"
	"github.com/trisacrypto/directory/pkg/trtl"
	trtlmock "github.com/trisacrypto/directory/pkg/trtl/mock"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
)

func NewTrtl() (t *Trtl, err error) {
	t = &Trtl{}

	// Create a temporary directory for the testing database
	if t.path, err = os.MkdirTemp("", "trtldb-*"); err != nil {
		return nil, err
	}

	// Create the mock configuration
	conf := trtlmock.Config()
	conf.Database.URL = "leveldb:///" + t.path
	if conf, err = conf.Mark(); err != nil {
		return nil, err
	}

	// Start the Trtl server
	if t.srv, err = trtl.New(conf); err != nil {
		return nil, err
	}
	t.sock = bufconn.New("")
	go t.srv.Run(t.sock.Listener)
	return t, nil
}

type Trtl struct {
	path   string
	srv    *trtl.Server
	sock   *bufconn.GRPCListener
	client store.Store
}

func (t *Trtl) Client() (_ store.Store, err error) {
	if t.client == nil {
		// Connect to the running trtl server
		if err = t.sock.Connect(context.Background()); err != nil {
			return nil, err
		}

		if t.client, err = trtlstore.NewMock(t.sock.Conn); err != nil {
			return nil, err
		}
	}
	return t.client, nil
}

func (t *Trtl) Shutdown() {
	if t.client != nil {
		t.client.Close()
	}
	t.srv.Shutdown()
	t.sock.Release()
	os.RemoveAll(t.path)
}
