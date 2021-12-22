package config

import (
	"errors"
	"time"

	"github.com/kelseyhightower/envconfig"
	honuconfig "github.com/rotationalio/honu/config"
	"github.com/rs/zerolog"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

// Config defines the struct that is expected to initialize the trtl server
// Note: because we need to validate the configuration, `config.New()`
// must be called to ensure that the `processed` is correctly set
type Config struct {
	Maintenance bool                `split_words:"true" default:"false"`
	BindAddr    string              `split_words:"true" default:":4436"`
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog  bool                `split_words:"true" default:"false"`
	Database    DatabaseConfig
	Replica     ReplicaConfig
	processed   bool
}

type DatabaseConfig struct {
	URL           string `split_words:"true" required:"true"`
	ReindexOnBoot bool   `split_words:"true" default:"false"`
}

type ReplicaConfig struct {
	Enabled        bool          `split_words:"true" default:"true"`
	PID            uint64        `split_words:"true" required:"false"`
	Region         string        `split_words:"true" required:"false"`
	Name           string        `split_words:"true" required:"false"`
	GossipInterval time.Duration `split_words:"true" default:"1m"`
	GossipSigma    time.Duration `split_words:"true" default:"5s"`
}

// New creates a new Config object, loading environment variables and defaults.
func New() (_ Config, err error) {
	var conf Config
	if err = envconfig.Process("trtl", &conf); err != nil {
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
func (c Config) GetHonuConfig() honuconfig.ReplicaConfig {
	return honuconfig.ReplicaConfig{
		Enabled:        true,
		BindAddr:       "",
		PID:            c.Replica.PID,
		Region:         c.Replica.Region,
		Name:           c.Replica.Name,
		GossipInterval: c.Replica.GossipInterval,
		GossipSigma:    c.Replica.GossipSigma,
	}
}

func (c Config) Validate() (err error) {
	// Validate config-specific constraints
	if err = c.Replica.Validate(); err != nil {
		return err
	}
	return nil
}

func (c ReplicaConfig) Validate() error {
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
