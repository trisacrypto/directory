package config

import (
	"errors"
	"strings"
)

type StoreConfig struct {
	URL           string `split_words:"true" required:"true"`
	ReindexOnBoot bool   `split_words:"true" default:"false"`
	Insecure      bool   `split_words:"true" default:"false"`
	CertPath      string `split_words:"true"`
	PoolPath      string `split_words:"true"`
}

func (c StoreConfig) Validate() error {
	// If the insecure flag isn't set then we must have certs when connecting to trtl.
	if strings.HasPrefix(c.URL, "trtl://") && !c.Insecure {
		if c.CertPath == "" || c.PoolPath == "" {
			return errors.New("invalid configuration: connecting to trtl over mTLS requires certs and cert pool")
		}
	}
	return nil
}
