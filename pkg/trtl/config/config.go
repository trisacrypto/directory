package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"

	"github.com/kelseyhightower/envconfig"
	honuconfig "github.com/rotationalio/honu/config"
	"github.com/rs/zerolog"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/trisa/pkg/trust"
)

// Config defines the struct that is expected to initialize the trtl server
// Note: because we need to validate the configuration, `config.New()`
// must be called to ensure that the `processed` is correctly set
type Config struct {
	Maintenance     bool                  `split_words:"true" default:"false"`
	BindAddr        string                `split_words:"true" default:":4436"`
	LogLevel        logger.LevelDecoder   `split_words:"true" default:"info"`
	ConsoleLog      bool                  `split_words:"true" default:"false"`
	Metrics         MetricsConfig         `split_words:"true"`
	Database        DatabaseConfig        `split_words:"true"`
	Replica         ReplicaConfig         `split_words:"true"`
	ReplicaStrategy ReplicaStrategyConfig `split_words:"true"`
	MTLS            MTLSConfig            `split_words:"true"`
	Backup          BackupConfig          `split_words:"true"`
	Sentry          sentry.Config         `split_words:"true"`
	processed       bool
}

type MetricsConfig struct {
	Addr    string `split_words:"true" default:":7777"`
	Enabled bool   `split_words:"true" default:"true"`
}

type DatabaseConfig struct {
	URL           string `split_words:"true" required:"true"`
	ReindexOnBoot bool   `split_words:"true" default:"false"`
}

type ReplicaConfig struct {
	Enabled        bool          `split_words:"true" default:"true" json:"enabled"`
	PID            uint64        `split_words:"true" required:"false" json:"pid"`
	Region         string        `split_words:"true" required:"false" json:"region"`
	Name           string        `split_words:"true" required:"false" json:"name"`
	GossipInterval time.Duration `split_words:"true" default:"1m" json:"gossip_interval"`
	GossipSigma    time.Duration `split_words:"true" default:"5s" json:"gossip_sigma"`
}

type ReplicaStrategyConfig struct {
	HostnamePID bool   `split_words:"true" default:"false"` // Set to true to use HostnamePID
	Hostname    string `envconfig:"TRTL_REPLICA_HOSTNAME"`  // configure the hostname from the environment
	FilePID     string `split_words:"true"`                 // Set to the PID filename path to use FilePID
	JSONConfig  string `split_words:"true"`                 // Set to the config map path to use JSONConfig
}

type MTLSConfig struct {
	Insecure  bool   `envconfig:"TRTL_INSECURE" default:"false"`
	ChainPath string `split_words:"true" required:"false"`
	CertPath  string `split_words:"true" required:"false"`

	// Cache loaded cert pool and certificate on config for reuse without reloading
	pool *x509.CertPool
	cert tls.Certificate
}

type BackupConfig struct {
	Enabled  bool          `split_words:"true" default:"false"`
	Interval time.Duration `split_words:"true" default:"24h"`
	Storage  string        `split_words:"true" required:"false"`
	Keep     int           `split_words:"true" default:"1"`
}

// New creates a new Config object, loading environment variables and defaults.
func New() (_ Config, err error) {
	var conf Config
	if err = envconfig.Process("trtl", &conf); err != nil {
		return Config{}, err
	}

	// Process replica strategies before configuration validation
	if conf.Replica, err = conf.Replica.Configure(conf.ReplicaStrategy.Strategies()...); err != nil {
		return Config{}, err
	}

	// Validate config-specific constraints
	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

func (c Config) IsZero() bool {
	return !c.processed
}

// Mark a manually constructed as processed as long as it is validated.
func (c Config) Mark() (Config, error) {
	if err := c.Validate(); err != nil {
		return c, err
	}
	c.processed = true
	return c, nil
}

// GetHonuConfig converts ReplicaConfig into honu's struct of the same name.
func (c Config) GetHonuConfig() honuconfig.Option {
	return honuconfig.WithReplica(honuconfig.ReplicaConfig{
		PID:    c.Replica.PID,
		Region: c.Replica.Region,
		Name:   c.Replica.Name,
	})
}

func (c Config) Validate() (err error) {
	// Validate config-specific constraints
	if err = c.Replica.Validate(); err != nil {
		return err
	}
	if err = c.MTLS.Validate(); err != nil {
		return err
	}
	return nil
}

func (c *ReplicaConfig) Validate() error {
	if c.Enabled {
		if c.PID == 0 {
			return errors.New("invalid configuration: PID required for enabled replica")
		}

		if c.Region == "" {
			return errors.New("invalid configuration: region required for enabled replica")
		}

		if c.GossipInterval == time.Duration(0) || c.GossipSigma == time.Duration(0) {
			return errors.New("invalid configuration: specify non-zero gossip interval and sigma")
		}
	}
	return nil
}

// Strategies extracts the replica configuration strategies from the configuration
func (c *ReplicaStrategyConfig) Strategies() (strategies []ReplicaStrategy) {
	strategies = make([]ReplicaStrategy, 0)
	if c.HostnamePID {
		strategies = append(strategies, HostnamePID(c.Hostname))
	}

	if c.FilePID != "" {
		strategies = append(strategies, FilePID(c.FilePID))
	}

	if c.JSONConfig != "" {
		strategies = append(strategies, JSONConfig(c.JSONConfig))
	}

	return strategies
}

func (c *MTLSConfig) Validate() error {
	if c.Insecure {
		return nil
	}

	if c.ChainPath == "" || c.CertPath == "" {
		return errors.New("invalid configuration: specify chain and cert paths")
	}

	return nil
}

func (c *MTLSConfig) ParseTLSConfig() (_ *tls.Config, err error) {
	if c.Insecure {
		return nil, errors.New("cannot create TLS configuration in insecure mode")
	}

	var certPool *x509.CertPool
	if certPool, err = c.GetCertPool(); err != nil {
		return nil, err
	}

	var cert tls.Certificate
	if cert, err = c.GetCert(); err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  certPool,
	}, nil
}

func (c *MTLSConfig) GetCertPool() (_ *x509.CertPool, err error) {
	if c.pool == nil {
		if err = c.load(); err != nil {
			return nil, err
		}
	}
	return c.pool, nil
}

func (c *MTLSConfig) GetCert() (_ tls.Certificate, err error) {
	if len(c.cert.Certificate) == 0 {
		if err = c.load(); err != nil {
			return c.cert, err
		}
	}
	return c.cert, nil
}

func (c *MTLSConfig) load() (err error) {
	var sz *trust.Serializer
	if sz, err = trust.NewSerializer(false); err != nil {
		return err
	}

	var pool trust.ProviderPool
	if pool, err = sz.ReadPoolFile(c.ChainPath); err != nil {
		return err
	}

	var provider *trust.Provider
	if provider, err = sz.ReadFile(c.CertPath); err != nil {
		return err
	}

	if c.pool, err = pool.GetCertPool(false); err != nil {
		return err
	}

	if c.cert, err = provider.GetKeyPair(); err != nil {
		return err
	}

	return nil
}
