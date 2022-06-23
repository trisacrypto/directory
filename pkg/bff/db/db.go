package db

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"

	"github.com/trisacrypto/directory/pkg/bff/config"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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
	cc   *grpc.ClientConn
	trtl trtl.TrtlClient
}

func (db *DB) Close() error {
	return db.cc.Close()
}
