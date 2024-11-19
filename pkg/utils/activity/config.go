package activity

import (
	"errors"
	"time"

	"github.com/trisacrypto/directory/pkg/utils/ensign"
)

type Config struct {
	Enabled           bool          `default:"false"`
	Topic             string        `required:"false"`
	Network           Network       `required:"false"`
	AggregationWindow time.Duration `split_words:"true" default:"5m"`
	Testing           bool          `default:"false"`
	Ensign            ensign.Config
}

func (c Config) Validate() (err error) {
	if c.Enabled {
		if c.Topic == "" {
			err = errors.Join(err, ErrMissingTopic)
		}

		if verr := c.Network.IsValid(); verr != nil {
			err = errors.Join(err, verr)
		}

		if verr := c.Ensign.IsValid(); verr != nil {
			err = errors.Join(err, verr)
		}
	}
	return err
}
