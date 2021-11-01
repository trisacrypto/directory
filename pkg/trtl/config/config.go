package config

import (
	"errors"
	"time"

	"github.com/kelseyhightower/envconfig"
	honuconfig "github.com/rotationalio/honu/config"
)

type Config struct {
	Enabled        bool          `split_words:"true" default:"true"`
	BindAddr       string        `split_words:"true" default:":4435"`
	PID            uint64        `split_words:"true" required:"false"`
	Region         string        `split_words:"true" required:"false"`
	Name           string        `split_words:"true" required:"false"`
	GossipInterval time.Duration `split_words:"true" default:"1m"`
	GossipSigma    time.Duration `split_words:"true" default:"5s"`
	Database       DatabaseConfig
	processed      bool
}

type DatabaseConfig struct {
	URL           string `split_words:"true" required:"true"`
	ReindexOnBoot bool   `split_words:"true" default:"false"`
}

// New creates a new Config object, loading environment variables and defaults.
func New() (_ Config, err error) {
	var conf Config
	if err = envconfig.Process("trtl", &conf); err != nil {
		return Config{}, err
	}
	if err = conf.Validate(); err != nil {
		return Config{}, err
	}
	conf.processed = true
	return conf, nil
}

func (c Config) IsZero() bool {
	return !c.processed
}

// GetHonuConfig converts ReplicaConfig into honu's struct of the same name.
func (c Config) GetHonuConfig() honuconfig.ReplicaConfig {
	return honuconfig.ReplicaConfig{
		Enabled:        true,
		BindAddr:       "",
		PID:            c.PID,
		Region:         c.Region,
		Name:           c.Name,
		GossipInterval: c.GossipInterval,
		GossipSigma:    c.GossipSigma,
	}
}

func (c Config) Validate() error {
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
