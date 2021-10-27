package config

import (
	"errors"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type ReplicaConfig struct {
	Enabled        bool          `split_words:"true" default:"true"`
	BindAddr       string        `split_words:"true" default:":4435"`
	PID            uint64        `split_words:"true" required:"false"`
	Region         string        `split_words:"true" required:"false"`
	Name           string        `split_words:"true" required:"false"`
	GossipInterval time.Duration `split_words:"true" default:"1m"`
	GossipSigma    time.Duration `split_words:"true" default:"5s"`
}

// New creates a new ReplicaConfig object, loading environment variables and defaults.
func New() (_ ReplicaConfig, err error) {
	var conf ReplicaConfig
	if err = envconfig.Process("trtl", &conf); err != nil {
		return ReplicaConfig{}, err
	}
	if err = conf.Validate(); err != nil {
		return ReplicaConfig{}, err
	}
	return conf, nil
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
