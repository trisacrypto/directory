package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"sync"

	"github.com/trisacrypto/directory/pkg/bff/config"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func Connect(conf config.DatabaseConfig) (db *DB, err error) {
	// Parse the URL to get the endpoint to the trtl server
	dsn, err := url.Parse(conf.URL)
	if err != nil {
		return nil, fmt.Errorf("could not parse dsn: %s", err)
	}

	var opts []grpc.DialOption
	if conf.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		var sz *trust.Serializer
		if sz, err = trust.NewSerializer(false); err != nil {
			return nil, err
		}

		var pool trust.ProviderPool
		if pool, err = sz.ReadPoolFile(conf.MTLS.PoolPath); err != nil {
			return nil, err
		}

		var provider *trust.Provider
		if provider, err = sz.ReadFile(conf.MTLS.CertPath); err != nil {
			return nil, err
		}

		var cert tls.Certificate
		if cert, err = provider.GetKeyPair(); err != nil {
			return nil, err
		}

		var certPool *x509.CertPool
		if certPool, err = pool.GetCertPool(false); err != nil {
			return nil, err
		}

		tlsConf := &tls.Config{
			ServerName:   dsn.Hostname(),
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
	}

	db = &DB{}
	if db.cc, err = grpc.Dial(dsn.Host, opts...); err != nil {
		return nil, err
	}

	db.trtl = trtl.NewTrtlClient(db.cc)
	return db, nil
}

func DirectConnect(cc *grpc.ClientConn) (db *DB, err error) {
	return &DB{
		cc:   cc,
		trtl: trtl.NewTrtlClient(cc),
	}, nil
}

// DB is a wrapper around a trtl client to provide BFF-specific database interactions
type DB struct {
	// Connection to the Trtl Database
	cc   *grpc.ClientConn
	trtl trtl.TrtlClient

	// Announcements collection and singleton helper
	announcements *Announcements
	muMakeAC      sync.Once
}

// Collection is an interface that identifies utilities that manage specific namespaces.
type Collection interface {
	Namespace() string
}

func (db *DB) Close() error {
	defer func() {
		db.cc = nil
		db.trtl = nil
	}()
	return db.cc.Close()
}

// Get is a high-level method for executing a fetch key request to a namespace in trtl.
func (db *DB) Get(ctx context.Context, key []byte, namespace string) (value []byte, err error) {
	req := &trtl.GetRequest{
		Key:       key,
		Namespace: namespace,
		Options: &trtl.Options{
			ReturnMeta: false,
		},
	}

	var rep *trtl.GetReply
	if rep, err = db.trtl.Get(ctx, req); err != nil {
		if serr, ok := status.FromError(err); ok {
			if serr.Code() == codes.NotFound {
				return nil, ErrNotFound
			}
		}
		return nil, err
	}
	return rep.Value, nil
}

// Put is a high-level method for executing a put value to key request to a namespace in trtl.
func (db *DB) Put(ctx context.Context, key, value []byte, namespace string) (err error) {
	req := &trtl.PutRequest{
		Key:       key,
		Value:     value,
		Namespace: namespace,
		Options: &trtl.Options{
			ReturnMeta: false,
		},
	}

	var rep *trtl.PutReply
	if rep, err = db.trtl.Put(ctx, req); err != nil {
		return err
	}

	if !rep.Success {
		return ErrUnsuccessfulPut
	}
	return nil
}

// Delete is a high-level method for executing a delete key-value request to a namespace in trtl.
func (db *DB) Delete(ctx context.Context, key []byte, namespace string) (err error) {
	req := &trtl.DeleteRequest{
		Key:       key,
		Namespace: namespace,
		Options: &trtl.Options{
			ReturnMeta: false,
		},
	}

	var rep *trtl.DeleteReply
	if rep, err = db.trtl.Delete(ctx, req); err != nil {
		return err
	}

	if !rep.Success {
		return ErrUnsuccessfulDelete
	}
	return nil
}
