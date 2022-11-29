package server_test

import (
	"crypto/x509/pkix"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/sectigo/server"
)

func TestCertIssuer(t *testing.T) {
	certs, err := server.NewCerts(server.Config{})
	require.NoError(t, err, "could not create cert issuer")

	sub := pkix.Name{
		CommonName:   "test.example.com",
		Organization: []string{"Test Organization"},
		Country:      []string{"US"},
	}

	data, err := certs.Issue(sub)
	require.NoError(t, err, "could not issue certs")
	require.NotEmpty(t, data, "no certs were returned")
}

func TestSerialNumber(t *testing.T) {
	sn := server.SerialNumber()
	sns := fmt.Sprintf("%X", sn)
	require.Len(t, sns, 32)

	sn2 := server.SerialNumber()
	sns2 := fmt.Sprintf("%X", sn2)
	require.NotEqual(t, sns, sns2)
}
