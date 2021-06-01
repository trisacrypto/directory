package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

// Config uses envconfig to load required settings from the environment and validate
// them in preparation for running the TRISA Global Directory Service.
type Config struct {
	BindAddr    string          `split_words:"true" default:":4433"`
	DirectoryID string          `split_words:"true" default:"vaspdirectory.net"`
	SecretKey   string          `split_words:"true" required:"true"`
	DatabaseURL string          `split_words:"true" required:"true"`
	Maintenance bool            `split_words:"true" default:"false"`
	LogLevel    LogLevelDecoder `split_words:"true" default:"info"`
	Sectigo     SectigoConfig
	Email       EmailConfig
	CertMan     CertManConfig
	Backup      BackupConfig
	Secrets     SecretsConfig
	processed   bool
}

type SectigoConfig struct {
	Username string `envconfig:"SECTIGO_USERNAME" required:"false"`
	Password string `envconfig:"SECTIGO_PASSWORD" required:"false"`
}

type EmailConfig struct {
	ServiceEmail   string `envconfig:"GDS_SERVICE_EMAIL" default:"TRISA Directory Service <admin@vaspdirectory.net>"`
	AdminEmail     string `envconfig:"GDS_ADMIN_EMAIL" default:"TRISA Admins <admin@trisa.io>"`
	SendGridAPIKey string `envconfig:"SENDGRID_API_KEY" required:"false"`
}

type CertManConfig struct {
	Interval time.Duration `split_words:"true" default:"10m"`
	Storage  string        `split_words:"true" required:"false"`
}

type BackupConfig struct {
	Enabled  bool          `split_words:"true" default:"false"`
	Interval time.Duration `split_words:"true" default:"24h"`
	Storage  string        `split_words:"true" required:"false"`
	Keep     int           `split_words:"true" default:"1"`
}

type SecretsConfig struct {
	Credentials string `envconfig:"GOOGLE_APPLICATION_CREDENTIALS" required:"false"`
	Project     string `envconfig:"GOOGLE_PROJECT_NAME" required:"false"`
}

// New creates a new Config object, loading environment variables and defaults.
func New() (_ Config, err error) {
	var conf Config
	if err = envconfig.Process("gds", &conf); err != nil {
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

// LogLevelDecoder deserializes the log level from a config string.
type LogLevelDecoder zerolog.Level

// Decode implements envconfig.Decoder
func (ll *LogLevelDecoder) Decode(value string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "panic":
		*ll = LogLevelDecoder(zerolog.PanicLevel)
	case "fatal":
		*ll = LogLevelDecoder(zerolog.FatalLevel)
	case "error":
		*ll = LogLevelDecoder(zerolog.ErrorLevel)
	case "warn":
		*ll = LogLevelDecoder(zerolog.WarnLevel)
	case "info":
		*ll = LogLevelDecoder(zerolog.InfoLevel)
	case "debug":
		*ll = LogLevelDecoder(zerolog.DebugLevel)
	case "trace":
		*ll = LogLevelDecoder(zerolog.TraceLevel)
	default:
		return fmt.Errorf("unknown log level %q", value)
	}
	return nil
}
