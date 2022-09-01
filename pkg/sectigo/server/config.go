package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

// Configure the server in a lightweight fashion by fetching environment variables.
type Config struct {
	BindAddr   string              `split_words:"true" default:":8831"`
	Mode       string              `split_words:"true" default:"release"`
	LogLevel   logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog bool                `split_words:"true" default:"false"`
	processed  bool
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

func (c Config) Validate() error {
	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return fmt.Errorf("%q is not a valid gin mode", c.Mode)
	}
	return nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

func (c Config) IsZero() bool {
	return !c.processed
}
