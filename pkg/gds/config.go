package trisads

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

// Settings uses envconfig to load required settings from the environment and
// validate them in preparation for running the TRISA Directory Service.
type Settings struct {
	BindAddr        string          `envconfig:"TRISADS_BIND_ADDR" default:":4433"`
	DatabaseDSN     string          `envconfig:"TRISADS_DATABASE" required:"true"`
	SectigoUsername string          `envconfig:"SECTIGO_USERNAME" required:"false"`
	SectigoPassword string          `envconfig:"SECTIGO_PASSWORD" required:"false"`
	SendGridAPIKey  string          `envconfig:"SENDGRID_API_KEY" required:"false"`
	ServiceEmail    string          `envconfig:"TRISADS_SERVICE_EMAIL" default:"admin@vaspdirectory.net"`
	AdminEmail      string          `envconfig:"TRISADS_ADMIN_EMAIL" default:"admin@trisa.io"`
	LogLevel        LogLevelDecoder `envconfig:"TRISADS_LOG_LEVEL" default:"info"`
	DirectoryID     string          `envconfig:"TRISADS_DIRECTORY_ID" default:"vaspdirectory.net"`
	SecretKey       string          `envconfig:"TRISADS_SECRET_KEY" required:"true"`
	CertManInterval time.Duration   `envconfig:"TRISADS_CERTMAN_INTERVAL" default:"10m"`
	CertManStorage  string          `envconfig:"TRISADS_CERTS_STORE" required:"false"`
}

// Config creates a new settings object, loading environment variables and defaults.
func Config() (_ *Settings, err error) {
	var conf Settings
	if err = envconfig.Process("trisads", &conf); err != nil {
		return nil, err
	}
	return &conf, nil
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
