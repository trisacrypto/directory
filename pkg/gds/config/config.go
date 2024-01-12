package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/confire"
	"github.com/rs/zerolog"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/store/config"
	"github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

// Config uses envconfig to load required settings from the environment and validate
// them in preparation for running the TRISA Global Directory Service.
type Config struct {
	DirectoryID string              `split_words:"true" default:"vaspdirectory.net"`
	SecretKey   string              `split_words:"true" required:"true"`
	Maintenance bool                `split_words:"true" default:"false"`
	LogLevel    logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog  bool                `split_words:"true" default:"false"`
	GDS         GDSConfig
	Admin       AdminConfig
	Members     MembersConfig
	Database    config.StoreConfig
	Email       EmailConfig
	CertMan     CertManConfig
	Backup      BackupConfig
	Secrets     SecretsConfig
	Sentry      sentry.Config
	Activity    activity.Config
	processed   bool
}

type GDSConfig struct {
	Enabled  bool   `envconfig:"GDS_API_ENABLED" default:"true"`
	BindAddr string `envconfig:"GDS_BIND_ADDR" default:":4433"`
}

type AdminConfig struct {
	Enabled      bool     `split_words:"true" default:"true"`
	BindAddr     string   `split_words:"true" default:":4434"`
	Mode         string   `split_words:"true" default:"release"`
	AllowOrigins []string `split_words:"true" default:"http://localhost,http://localhost:3000,http://localhost:3001"`
	CookieDomain string   `split_words:"true"`
	Audience     string   `split_words:"true"`
	Oauth        OauthConfig

	// TokenKeys are the paths to RSA JWT signing keys in PEM encoded format. The
	// environment variable should be a comma separated list of keyid:path/to/key.pem
	// Multiple keys are used in order to rotate keys regularly; keyids therefore must
	// be sortable; in general we prefer to use ksuid for key ids.
	TokenKeys map[string]string `split_words:"true"`
}

type OauthConfig struct {
	GoogleAudience         string   `split_words:"true"`
	AuthorizedEmailDomains []string `split_words:"true"`
}

type MembersConfig struct {
	Enabled  bool   `split_words:"true" default:"true"`
	BindAddr string `split_words:"true" default:":4435"`
	Insecure bool   `split_words:"true" default:"false"`
	Certs    string `split_words:"true"`
	CertPool string `split_words:"true"`
}

type EmailConfig struct {
	ServiceEmail         string `envconfig:"GDS_SERVICE_EMAIL" default:"TRISA Directory Service <admin@vaspdirectory.net>"`
	AdminEmail           string `envconfig:"GDS_ADMIN_EMAIL" default:"TRISA Admins <admin@trisa.io>"`
	SendGridAPIKey       string `envconfig:"SENDGRID_API_KEY" required:"false"`
	DirectoryID          string `envconfig:"GDS_DIRECTORY_ID" default:"vaspdirectory.net"`
	VerifyContactBaseURL string `envconfig:"GDS_VERIFY_CONTACT_URL" default:"https://vaspdirectory.net/verify"`
	AdminReviewBaseURL   string `envconfig:"GDS_ADMIN_REVIEW_URL" default:"https://admin.vaspdirectory.net/vasps/"`
	Testing              bool   `split_words:"true" default:"false"`
	Storage              string `split_words:"true" default:""`
}

type CertManConfig struct {
	Enabled            bool          `split_words:"true" default:"true"`
	RequestInterval    time.Duration `split_words:"true" default:"10m"`
	ReissuanceInterval time.Duration `split_words:"true" default:"24h"`
	Storage            string        `split_words:"true" required:"false"`
	DirectoryID        string        `envconfig:"GDS_DIRECTORY_ID" default:"vaspdirectory.net"`
	Sectigo            sectigo.Config
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
func New() (conf Config, err error) {
	// Load and validate the configuration from the environment.
	if err = confire.Process("gds", &conf); err != nil {
		return Config{}, err
	}

	// Preprocess authorized domains
	for i, domain := range conf.Admin.Oauth.AuthorizedEmailDomains {
		conf.Admin.Oauth.AuthorizedEmailDomains[i] = strings.ToLower(strings.Trim(strings.TrimSpace(domain), "\"'"))
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

func (c Config) Validate() (err error) {
	if err = c.GDS.Validate(); err != nil {
		return err
	}

	if err = c.Admin.Validate(); err != nil {
		return err
	}

	if err = c.Members.Validate(); err != nil {
		return err
	}

	if err = c.Database.Validate(); err != nil {
		return err
	}

	if err = c.Email.Validate(); err != nil {
		return err
	}

	if err = c.CertMan.Validate(); err != nil {
		return err
	}

	return nil
}

func (c GDSConfig) Validate() error {
	if c.Enabled {
		if c.BindAddr == "" {
			return errors.New("invalid configuration: bind addr is required for enabled GDS")
		}
	}

	return nil
}

func (c AdminConfig) Validate() error {
	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return fmt.Errorf("%q is not a valid gin mode", c.Mode)
	}

	if c.Enabled {
		if err := c.Oauth.Validate(); err != nil {
			return err
		}

		if len(c.TokenKeys) == 0 {
			return errors.New("invalid configuration: token keys required for enabled admin")
		}
	}

	return nil
}

func (c OauthConfig) Validate() error {
	// Check configurations that are only required if the admin API is enabled
	if c.GoogleAudience == "" {
		return errors.New("invalid configuration: oauth audience required for enabled admin")
	}

	if len(c.AuthorizedEmailDomains) == 0 {
		return errors.New("invalid configuration: authorized email domains required for enabled admin")
	}

	return nil
}

func (c MembersConfig) Validate() error {
	// If the insecure flag isn't set then we must have certs.
	if !c.Insecure {
		if c.Certs == "" || c.CertPool == "" {
			return errors.New("invalid configuration: serving mTLS requires the path to certs and the cert pool")
		}
	}

	return nil
}

func (c EmailConfig) Validate() error {
	if c.AdminReviewBaseURL != "" && !strings.HasSuffix(c.AdminReviewBaseURL, "/") {
		return errors.New("invalid configuration: admin review base URL must end in a /")
	}

	if !c.Testing {
		if c.SendGridAPIKey == "" || c.ServiceEmail == "" {
			return errors.New("invalid configuration: sendgrid api key and service email are required")
		}

		if c.Storage != "" {
			return errors.New("invalid configuration: email archiving is only supported in testing mode")
		}
	}

	return nil
}

func (c CertManConfig) Validate() (err error) {
	if err = c.Sectigo.Validate(); err != nil {
		return err
	}

	return nil
}
