package activity

import (
	"time"

	"github.com/trisacrypto/directory/pkg/utils/ensign"
)

type Config struct {
	Enabled           bool          `split_words:"true" default:"false"`
	Topic             string        `split_words:"true"`
	Network           Network       `split_words:"true"`
	AggregationWindow time.Duration `split_words:"true" default:"5m"`
	Testing           bool          `split_words:"true" default:"false"`
	Ensign            ensign.Config
}

func (c Config) Validate() (err error) {
	if c.Enabled {
		if c.Topic == "" {
			return ErrMissingTopic
		}

		if err = c.Ensign.Validate(); err != nil {
			return err
		}
	}

	return nil
}
