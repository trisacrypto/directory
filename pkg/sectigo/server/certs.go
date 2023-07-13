package server

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/trisa/pkg/trisa/keys/signature"
)

func (s *Server) CreateCerts(params Params, profile string, batchID int) {
	// Wait 5 minutes to simulate real processing times
	time.Sleep(5 * time.Minute)

	sub := pkix.Name{
		CommonName:   params.Get(sectigo.ParamCommonName, ""),
		Organization: []string{params.Get(sectigo.ParamOrganizationName, "Staging")},
		Locality:     []string{params.Get(sectigo.ParamLocalityName, "Queenstown")},
		Province:     []string{params.Get(sectigo.ParamStateOrProvinceName, "MD")},
		Country:      []string{params.Get(sectigo.ParamCountryName, "US")},
	}

	// TODO: handle dNSNames
	certs, err := s.certs.Issue(sub)
	if err != nil {
		sentry.Error(nil).Err(err).Msg("could not issue certs")
		if err = s.store.RejectBatch(batchID, err.Error()); err != nil {
			sentry.Error(nil).Err(err).Msg("could not reject batch")
		}
		return
	}

	if err := s.store.AddCert(batchID, certs); err != nil {
		sentry.Error(nil).Err(err).Msg("could not add certs to store")
		return
	}

	log.Info().Str("common_name", sub.CommonName).Msg("certificates issued")
}

type Certs struct {
	cacrt *x509.Certificate
	cakey crypto.PrivateKey
}

// TODO: create a different CA for each profile
func NewCerts(conf Config) (certs *Certs, err error) {
	certs = &Certs{}
	if certs.cacrt, certs.cakey, err = conf.CA(); err != nil {
		return nil, err
	}
	return certs, nil
}

// TODO: handle dNSNames
func (c *Certs) Issue(sub pkix.Name) (_ []byte, err error) {
	var (
		key    *rsa.PrivateKey
		keyID  string
		signed []byte
		chain  bytes.Buffer
	)

	// Create the certificate
	cert := &x509.Certificate{
		SerialNumber: SerialNumber(),
		Subject:      sub,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 1, 0),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{sub.CommonName},
	}

	if key, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return nil, err
	}

	if keyID, err = signature.New(&key.PublicKey); err != nil {
		return nil, err
	}
	cert.SubjectKeyId = []byte(keyID)

	if signed, err = x509.CreateCertificate(rand.Reader, cert, c.cacrt, &key.PublicKey, c.cakey); err != nil {
		return nil, err
	}

	// Encode the certificate chain
	if err = pem.Encode(&chain, &pem.Block{Type: "CERTIFICATE", Bytes: signed}); err != nil {
		return nil, err
	}

	if err = pem.Encode(&chain, &pem.Block{Type: "CERTIFICATE", Bytes: c.cacrt.Raw}); err != nil {
		return nil, err
	}

	if err = pem.Encode(&chain, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		return nil, err
	}

	return chain.Bytes(), nil
}

func InitCA(commonName string) (cert *x509.Certificate, priv crypto.PrivateKey, err error) {
	// Generate a new self-signed certificate to issue certs
	template := &x509.Certificate{
		SerialNumber: SerialNumber(),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"TRISA", "Rotational"},
			Country:      []string{"US"},
			Province:     []string{"MD"},
			Locality:     []string{"Queenstown"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		DNSNames:              []string{commonName},
	}

	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return nil, nil, err
	}

	var signed []byte
	if signed, err = x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key); err != nil {
		return nil, nil, err
	}

	if cert, err = x509.ParseCertificate(signed); err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}

func SerialNumber() *big.Int {
	sn := make([]byte, 16)
	rand.Read(sn)

	i := &big.Int{}
	return i.SetBytes(sn)
}

type Params map[string]string

func (p Params) Get(name, defaultValue string) string {
	if val, ok := p[name]; ok {
		return val
	}
	return defaultValue
}
