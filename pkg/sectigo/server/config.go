package server

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/trisa/pkg/trust"
)

// Configure the server in a lightweight fashion by fetching environment variables.
type Config struct {
	BindAddr   string              `split_words:"true" default:":8831"`
	Mode       string              `split_words:"true" default:"release"`
	LogLevel   logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog bool                `split_words:"true" default:"false"`
	CAPath     string              `split_words:"true"`
	Auth       AuthConfig
	processed  bool
}

type AuthConfig struct {
	Username string   `required:"true"`
	Password string   `required:"true"`
	Issuer   string   `default:"https://cathy.test-net.io"`
	Subject  string   `default:"/account/42/user/staging"`
	Scopes   []string `default:"ROLE_USER"`
	Secret   string   `required:"false"`
}

func NewConfig() (conf Config, err error) {
	if err = envconfig.Process("sias", &conf); err != nil {
		return Config{}, err
	}

	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

func (c Config) Mark() (Config, error) {
	if err := c.Validate(); err != nil {
		return c, err
	}
	c.processed = true
	return c, nil
}

func (c Config) Validate() error {
	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return fmt.Errorf("%q is not a valid gin mode", c.Mode)
	}

	if err := c.Auth.Validate(); err != nil {
		return err
	}

	return nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

func (c Config) IsZero() bool {
	return !c.processed
}

func (c Config) CA() (cert *x509.Certificate, priv crypto.PrivateKey, err error) {
	// Load the CA from disk if the path is specified
	if c.CAPath != "" {
		var sz *trust.Serializer
		if sz, err = trust.NewSerializer(false); err != nil {
			return nil, nil, err
		}

		var ca *trust.Provider
		if ca, err = sz.ReadFile(c.CAPath); err != nil {
			return nil, nil, err
		}

		var catls tls.Certificate
		if catls, err = ca.GetKeyPair(); err != nil {
			return nil, nil, err
		}

		if cert, err = x509.ParseCertificate(catls.Certificate[0]); err != nil {
			return nil, nil, err
		}

		return cert, catls.PrivateKey, nil
	}

	// Generate a new self-signed certificate to issue certs
	template := &x509.Certificate{
		SerialNumber: SerialNumber(),
		Subject: pkix.Name{
			CommonName:   "trisa.dev",
			Organization: []string{"TRISA", "Rotational"},
			Country:      []string{"USA"},
			Province:     []string{"MD"},
			Locality:     []string{"Queenstown"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		DNSNames:              []string{"trisa.dev"},
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

func (c AuthConfig) Validate() error {
	if c.Secret != "" {
		if secret, err := base64.StdEncoding.DecodeString(c.Secret); err != nil {
			return fmt.Errorf("invalid configuration: cannot parse base64 secret: %w", err)
		} else {
			if len(secret) != 32 {
				return fmt.Errorf("invalid configuration: secret must be 32 bytes")
			}
		}
	}
	return nil
}

func (c AuthConfig) ParseSecret() []byte {
	if c.Secret != "" {
		if secret, err := base64.StdEncoding.DecodeString(c.Secret); err == nil {
			return secret
		}
		log.Warn().Msg("could not decode base64 secret -- using random secret instead")
	}

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		panic(err)
	}
	return secret
}
