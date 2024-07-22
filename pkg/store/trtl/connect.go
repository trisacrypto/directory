package trtl

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"

	"github.com/trisacrypto/directory/pkg/store/config"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func Connect(conf config.StoreConfig) (conn *grpc.ClientConn, err error) {
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
		if pool, err = sz.ReadPoolFile(conf.PoolPath); err != nil {
			return nil, err
		}

		var provider *trust.Provider
		if provider, err = sz.ReadFile(conf.CertPath); err != nil {
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

	// Connect the replica client
	if conn, err = grpc.NewClient(dsn.Host, opts...); err != nil {
		return nil, err
	}
	return conn, nil
}
