package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/trisa/pkg/trust"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "certs"
	app.Version = pkg.Version()
	app.Usage = "local disk certificate management and test CA"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:      "decrypt",
			Usage:     "decrypt a PKCS12 zip file from Sectigo and save as gzip provider",
			ArgsUsage: "src",
			Category:  "certs",
			Action:    decrypt,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "o, outpath",
					Usage: "location to write the decrypted PEM file out to",
				},
				cli.StringFlag{
					Name:  "p, password",
					Usage: "the PKCS12 password to decrypt the file with",
				},
			},
		},
		{
			Name:      "pool",
			Usage:     "create a trust chain pool from certificates on disk",
			ArgsUsage: "src [src ...]",
			Category:  "certs",
			Action:    pool,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "o, outpath",
					Usage: "location to write the chain pool out to",
					Value: "trisa.zip",
				},
			},
		},
		{
			Name:     "init",
			Usage:    "create CA certs and keys if they do not exist",
			Category: "CA",
			Action:   initCA,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "c, ca",
					Usage: "location to write the CA certificates to",
					Value: "fixtures/certs/ca.gz",
				},
				cli.BoolFlag{
					Name:  "f, force",
					Usage: "overwrite keys even if they already exist",
				},
			},
		},
		{
			Name:     "issue",
			Usage:    "issue certs signed by the local CA certificates",
			Category: "CA",
			Action:   issue,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "c, ca",
					Usage: "location to read the CA certificates from",
					Value: "fixtures/certs/ca.gz",
				},
				cli.StringFlag{
					Name:  "o, outpath",
					Usage: "location to write the issued certificates to",
					Value: "",
				},
				cli.StringFlag{
					Name:  "P, password",
					Usage: "PKCS12 encrypt the output file if specified",
				},
				cli.StringFlag{
					Name:  "n, name",
					Usage: "common name for the certificate",
				},
				cli.StringFlag{
					Name:  "O, organization",
					Usage: "name of organization to issue certificates for",
				},
				cli.StringFlag{
					Name:  "C, country",
					Usage: "country of the organization",
				},
				cli.StringFlag{
					Name:  "p, province",
					Usage: "province or state of the organization",
				},
				cli.StringFlag{
					Name:  "l, locality",
					Usage: "locality or city of the organization",
				},
				cli.StringFlag{
					Name:  "a, address",
					Usage: "streed address of the organization",
				},
				cli.StringFlag{
					Name:  "z, postcode",
					Usage: "postal code of the organization",
				},
			},
		},
	}

	app.Run(os.Args)
}

func decrypt(c *cli.Context) (err error) {
	if c.NArg() != 1 {
		return cli.NewExitError("specify one source PKCS12 file", 1)
	}

	path := c.Args()[0]
	var outpath, password string
	if password = c.String("password"); password == "" {
		return cli.NewExitError("specify password to decrypt", 1)
	}
	if outpath = c.String("outpath"); outpath == "" {
		outpath = strings.TrimSuffix(path, filepath.Ext(path)) + ".gz"
	}

	var sz *trust.Serializer
	if sz, err = trust.NewSerializer(true, password); err != nil {
		return cli.NewExitError(err, 1)
	}

	var provider *trust.Provider
	if provider, err = sz.ReadFile(path); err != nil {
		return cli.NewExitError(err, 1)
	}

	if sz, err = trust.NewSerializer(false); err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = sz.WriteFile(provider, outpath); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func pool(c *cli.Context) (err error) {
	if c.NArg() < 1 {
		return cli.NewExitError("no providers specified", 1)
	}

	pool := trust.NewPool()
	for _, path := range c.Args() {
		var (
			sz       *trust.Serializer
			provider *trust.Provider
		)
		if sz, err = trust.NewSerializer(false); err != nil {
			return cli.NewExitError(err, 1)
		}
		if provider, err = sz.ReadFile(path); err != nil {
			return cli.NewExitError(err, 1)
		}
		pool.Add(provider)
	}

	var sz *trust.Serializer
	if sz, err = trust.NewSerializer(false); err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = sz.WritePoolFile(pool, c.String("outpath")); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func initCA(c *cli.Context) (err error) {
	force := c.Bool("force")
	path := c.String("ca")

	if !force {
		if _, err = os.Stat(path); err == nil {
			return cli.NewExitError("certificate file already exists", 1)
		}
	}

	// Create the CA certificate
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1942),
		Subject: pkix.Name{
			CommonName:    "trisa.dev",
			Organization:  []string{"TRISA", "Rotational Labs"},
			Country:       []string{"US"},
			Province:      []string{"MD"},
			Locality:      []string{"Queenstown"},
			StreetAddress: []string{"215 Alynn Way"},
			PostalCode:    []string{"21658"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		DNSNames:              []string{"trisa.dev"},
	}

	// Create private key
	var priv *rsa.PrivateKey
	if priv, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return cli.NewExitError(fmt.Errorf("could not create private key: %s", err), 1)
	}

	pub := &priv.PublicKey
	signed, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("create CA certs failed: %s", err), 1)
	}

	var chain bytes.Buffer
	if err = pem.Encode(&chain, &pem.Block{Type: "CERTIFICATE", Bytes: signed}); err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = pem.Encode(&chain, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return cli.NewExitError(err, 1)
	}

	var provider *trust.Provider
	if provider, err = trust.New(chain.Bytes()); err != nil {
		return cli.NewExitError(err, 1)
	}

	var sz *trust.Serializer
	if sz, err = trust.NewSerializer(false); err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = sz.WriteFile(provider, path); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func issue(c *cli.Context) (err error) {
	if c.String("name") == "" {
		return cli.NewExitError("please specify the common name for the cert", 1)
	}

	if c.String("organization") == "" {
		return cli.NewExitError("specify the name of the organization", 1)
	}

	// Load the CA certificates
	var ca *trust.Provider
	if caPath := c.String("ca"); caPath != "" {
		var sz *trust.Serializer
		if sz, err = trust.NewSerializer(false); err != nil {
			return cli.NewExitError(err, 1)
		}
		if ca, err = sz.ReadFile(caPath); err != nil {
			return cli.NewExitError(err, 1)
		}
	} else {
		return cli.NewExitError("specify path to CA certs", 1)
	}

	var catls tls.Certificate
	if catls, err = ca.GetKeyPair(); err != nil {
		return cli.NewExitError(err, 1)
	}

	var cacert *x509.Certificate
	if cacert, err = x509.ParseCertificate(catls.Certificate[0]); err != nil {
		return cli.NewExitError(err, 1)
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1945),
		Subject: pkix.Name{
			CommonName:    c.String("name"),
			Organization:  []string{c.String("organization")},
			Country:       []string{c.String("country")},
			Province:      []string{c.String("province")},
			Locality:      []string{c.String("locality")},
			StreetAddress: []string{c.String("address")},
			PostalCode:    []string{c.String("postcode")},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 5, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{c.String("name")},
	}

	// Create private key
	var priv *rsa.PrivateKey
	if priv, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return cli.NewExitError(fmt.Errorf("could not create private key: %s", err), 1)
	}

	pub := &priv.PublicKey
	signed, err := x509.CreateCertificate(rand.Reader, cert, cacert, pub, catls.PrivateKey)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("create CA certs failed: %s", err), 1)
	}

	// Encode the certificate
	var chain bytes.Buffer
	if err = pem.Encode(&chain, &pem.Block{Type: "CERTIFICATE", Bytes: signed}); err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = pem.Encode(&chain, &pem.Block{Type: "CERTIFICATE", Bytes: cacert.Raw}); err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = pem.Encode(&chain, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return cli.NewExitError(err, 1)
	}

	var (
		provider *trust.Provider
		sz       *trust.Serializer
		password string
		outpath  string
	)

	if provider, err = trust.New(chain.Bytes()); err != nil {
		return cli.NewExitError(err, 1)
	}

	if password = c.String("password"); password != "" {
		if sz, err = trust.NewSerializer(true, password); err != nil {
			return cli.NewExitError(err, 1)
		}
	} else {
		if sz, err = trust.NewSerializer(false); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	if outpath = c.String("outpath"); outpath == "" {
		outpath = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(c.String("name")), " ", "_"))
		outpath += ".gz"
	}

	if err = sz.WriteFile(provider, outpath); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}
