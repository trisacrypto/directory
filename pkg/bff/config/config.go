package config

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

// Config uses envconfig to load the required settings from the environment, parse and
// validate them in preparation for running the GDS BFF API service.
type Config struct {
	Maintenance bool                `split_words:"true" default:"false"`
	BindAddr    string              `split_words:"true" default:":4437"`
	Mode        string              `split_words:"true" default:"release"`
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog  bool                `split_words:"true" default:"false"`
	TestNet     DirectoryConfig
	MainNet     DirectoryConfig
	processed   bool
}

// DirectoryConfig is a generic configuration for connecting to a GDS service.
type DirectoryConfig struct {
	Insecure bool          `split_words:"true" default:"true"`
	Endpoint string        `split_words:"true" required:"true"`
	Timeout  time.Duration `split_words:"true" default:"10s"`
}

// New creates a new Config object from environment variables prefixed with GDS_BFF.
func New() (conf Config, err error) {
	if err = envconfig.Process("gds_bff", &conf); err != nil {
		return Config{}, err
	}

	// Validate the configuration
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

// Validate the config to make sure that it is usable to run the GDS BFF server.
func (c Config) Validate() error {
	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return fmt.Errorf("%q is not a valid gin mode", c.Mode)
	}
	return nil
}
