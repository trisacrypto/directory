package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

// Config uses envconfig to load required settings from the environment and validate
// them in preparation for running the TRISA Global Directory Service.
type Config struct {
	DirectoryID string          `split_words:"true" default:"vaspdirectory.net"`
	SecretKey   string          `split_words:"true" required:"true"`
	Maintenance bool            `split_words:"true" default:"false"`
	LogLevel    LogLevelDecoder `split_words:"true" default:"info"`
	ConsoleLog  bool            `split_words:"true" default:"false"`
	GDS         GDSConfig
	Admin       AdminConfig
	Replica     ReplicaConfig
	Database    DatabaseConfig
	Sectigo     SectigoConfig
	Email       EmailConfig
	CertMan     CertManConfig
	Backup      BackupConfig
	Secrets     SecretsConfig
	processed   bool
}

type GDSConfig struct {
	Enabled  bool   `envconfig:"GDS_API_ENABLED" default:"true"`
	BindAddr string `envconfig:"GDS_BIND_ADDR" default:":4433"`
}

type AdminConfig struct {
	Enabled           bool     `split_words:"true" default:"true"`
	BindAddr          string   `split_words:"true" default:":4434"`
	Mode              string   `split_words:"true" default:"release"`
	Audience          string   `split_words:"true"`
	AuthorizedDomains []string `split_words:"true"`

	// TokenKeys are the paths to RSA JWT signing keys in PEM encoded format. The
	// environment variable should be a comma separated list of keyid:path/to/key.pem
	// Multiple keys are used in order to rotate keys regularly; keyids therefore must
	// be sortable; in general we prefer to use ksuid for key ids.
	TokenKeys map[string]string `split_words:"true"`
}

type ReplicaConfig struct {
	Enabled        bool          `split_words:"true" default:"true"`
	BindAddr       string        `split_words:"true" default:":4435"`
	PID            uint64        `split_words:"true" required:"false"`
	Region         string        `split_words:"true" required:"false"`
	Name           string        `split_words:"true" required:"false"`
	GossipInterval time.Duration `split_words:"true" default:"1m"`
	GossipSigma    time.Duration `split_words:"true" default:"5s"`
}

type DatabaseConfig struct {
	URL           string `split_words:"true" required:"true"`
	ReindexOnBoot bool   `split_words:"true" default:"false"`
}

type SectigoConfig struct {
	Username string `envconfig:"SECTIGO_USERNAME" required:"false"`
	Password string `envconfig:"SECTIGO_PASSWORD" required:"false"`
	Profile  string `envconfig:"SECTIGO_PROFILE" default:"CipherTrace EE"`
}

type EmailConfig struct {
	ServiceEmail         string `envconfig:"GDS_SERVICE_EMAIL" default:"TRISA Directory Service <admin@vaspdirectory.net>"`
	AdminEmail           string `envconfig:"GDS_ADMIN_EMAIL" default:"TRISA Admins <admin@trisa.io>"`
	SendGridAPIKey       string `envconfig:"SENDGRID_API_KEY" required:"false"`
	DirectoryID          string `envconfig:"GDS_DIRECTORY_ID" default:"vaspdirectory.net"`
	VerifyContactBaseURL string `envconfig:"GDS_VERIFY_CONTACT_URL" default:"https://vaspdirectory.net/verify-contact"`
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
	Testing     bool   `split_words:"true" default:"false"`
}

// New creates a new Config object, loading environment variables and defaults.
func New() (_ Config, err error) {
	var conf Config
	if err = envconfig.Process("gds", &conf); err != nil {
		return Config{}, err
	}

	// Validate config-specific constraints
	if err = conf.Admin.Validate(); err != nil {
		return Config{}, err
	}

	if err = conf.Replica.Validate(); err != nil {
		return Config{}, err
	}

	if err = conf.Sectigo.Validate(); err != nil {
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

func (c AdminConfig) Validate() error {
	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return fmt.Errorf("%q is not a valid gin mode", c.Mode)
	}

	if c.Enabled {
		// Check configurations that are only required if the admin API is enabled
		if c.Audience == "" {
			return errors.New("invalid configuration: audience required for enabled admin")
		}

		if len(c.AuthorizedDomains) == 0 {
			return errors.New("invalid configuration: authorized domains required for enabled admin")
		}

		if len(c.TokenKeys) == 0 {
			return errors.New("invalid configuration: token keys required for enabled admin")
		}
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

func (c SectigoConfig) Validate() error {
	// Check valid certificate profiles
	validProfile := false
	for _, profile := range sectigo.AllProfiles {
		if profile == c.Profile {
			validProfile = true
			break
		}
	}

	if !validProfile {
		return fmt.Errorf("%q is not a valid Sectigo profile name, specify one of %s", c.Profile, strings.Join(sectigo.AllProfiles[:], ", "))
	}
	return nil
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
