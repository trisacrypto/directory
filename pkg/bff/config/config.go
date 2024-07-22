package config

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/rotationalio/confire"
	"github.com/rs/zerolog"
	"github.com/trisacrypto/directory/pkg/store/config"
	"github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/ensign"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/trisa/pkg/trisa/mtls"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
)

const (
	TestNet = "testnet"
	MainNet = "mainnet"
)

// Config uses envconfig to load the required settings from the environment, parse and
// validate them in preparation for running the GDS BFF API service.
type Config struct {
	Maintenance  bool                `split_words:"true" default:"false"`
	BindAddr     string              `split_words:"true" default:":4437"`
	Mode         string              `split_words:"true" default:"release"`
	LogLevel     logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog   bool                `split_words:"true" default:"false"`
	AllowOrigins []string            `split_words:"true" default:"http://localhost,http://localhost:3000,http://localhost:3003"`
	RegisterURL  string              `split_words:"true" required:"true"` // Trailing slash is not allowed
	LoginURL     string              `split_words:"true" required:"true"` // Trailing slash is not allowed
	CookieDomain string              `split_words:"true"`
	ServeDocs    bool                `split_words:"true" default:"false"`
	UserCache    CacheConfig         `split_words:"true"`
	Auth0        AuthConfig
	TestNet      NetworkConfig
	MainNet      NetworkConfig
	Database     config.StoreConfig
	Email        EmailConfig
	Sentry       sentry.Config
	Activity     activity.Config
	processed    bool
}

// AuthConfig handles Auth0 configuration and authentication
type AuthConfig struct {
	Domain        string        `split_words:"true" required:"true"`
	Issuer        string        `split_words:"true" required:"false"` // Set to the custom domain if enabled in Auth0 (ensure trailing slash is set if required!)
	Audience      string        `split_words:"true" required:"true"`
	ProviderCache time.Duration `split_words:"true" default:"5m"`
	ClientID      string        `split_words:"true"`
	ClientSecret  string        `split_words:"true"`
	Testing       bool          `split_words:"true" default:"false"` // If true a mock authenticator is used for testing
}

// NetworkConfig contains sub configurations for connecting to specific GDS and members
// services.
type NetworkConfig struct {
	Database  config.StoreConfig
	Directory DirectoryConfig
	Members   MembersConfig
}

// DirectoryConfig is a generic configuration for connecting to a GDS service.
type DirectoryConfig struct {
	Insecure bool          `split_words:"true" default:"true"`
	Endpoint string        `split_words:"true" required:"true"`
	Timeout  time.Duration `split_words:"true" default:"10s"`
}

// MembersConfig is a configuration for connecting to a members service.
type MembersConfig struct {
	Endpoint string        `split_words:"true" required:"true"`
	Timeout  time.Duration `split_words:"true" default:"10s"`
	MTLS     MTLSConfig
}

type MTLSConfig struct {
	Insecure bool   `split_words:"true"`
	CertPath string `split_words:"true"`
	PoolPath string `split_words:"true"`
}

// EmailConfig defines how emails are sent from the BFF.
type EmailConfig struct {
	ServiceEmail   string `envconfig:"GDS_BFF_SERVICE_EMAIL" default:"TRISA Directory Service <admin@vaspdirectory.net>"`
	SendGridAPIKey string `envconfig:"SENDGRID_API_KEY" required:"false"`
	Testing        bool   `split_words:"true" default:"false"`
	Storage        string `split_words:"true" default:""`
}

type CacheConfig struct {
	Enabled    bool          `split_words:"true" default:"false"`
	Size       uint          `split_words:"true" default:"16384"`
	Expiration time.Duration `split_words:"true" default:"8h"`
}

type ActivityConfig struct {
	Enabled bool   `split_words:"true" default:"false"`
	Topic   string `split_words:"true"`
	Ensign  ensign.Config
}

// New creates a new Config object from environment variables prefixed with GDS_BFF.
func New() (conf Config, err error) {
	// Load and validate the configuration from the environment.
	if err = confire.Process("gds_bff", &conf); err != nil {
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
func (c Config) Validate() (err error) {
	if err = validateURL(c.LoginURL); err != nil {
		return fmt.Errorf("invalid configuration: invalid login url: %w", err)
	}

	if err = validateURL(c.RegisterURL); err != nil {
		return fmt.Errorf("invalid configuration: invalid register url: %w", err)
	}

	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return fmt.Errorf("%q is not a valid gin mode", c.Mode)
	}

	if err = c.Auth0.Validate(); err != nil {
		return err
	}

	if err = c.TestNet.Validate(); err != nil {
		return err
	}

	if err = c.Database.Validate(); err != nil {
		return err
	}

	if err = c.Email.Validate(); err != nil {
		return err
	}

	if err = c.Sentry.Validate(); err != nil {
		return err
	}

	if err = c.UserCache.Validate(); err != nil {
		return err
	}

	return nil
}

func (c NetworkConfig) Validate() error {
	if err := c.Members.Validate(); err != nil {
		return err
	}
	return nil
}

func (c MembersConfig) Validate() error {
	if err := c.MTLS.Validate(); err != nil {
		return fmt.Errorf("invalid members configuration: %w", err)
	}
	return nil
}

func (c AuthConfig) Validate() error {
	if _, err := c.IssuerURL(); err != nil {
		return err
	}

	if c.ProviderCache == 0 {
		return errors.New("invalid configuration: auth0 provider cache duration should be longer than 0")
	}

	// If testing is false then the client id and secret are required
	if !c.Testing {
		if c.ClientID == "" {
			return errors.New("invalid configuration: auth0 client id is required in production")
		}

		if c.ClientSecret == "" {
			return errors.New("invalid configuration: auth0 client secret is required in production")
		}
	}

	return nil
}

func (c AuthConfig) IssuerURL() (u *url.URL, err error) {
	// Use the configured issuer if its set.
	if c.Issuer != "" {
		if u, err = url.Parse(c.Issuer); err != nil {
			return nil, errors.New("invalid configuration: could not parse issuer")
		}
		return u, nil
	}

	if c.Domain == "" {
		return nil, errors.New("invalid configuration: auth0 domain must be configured")
	}

	// Do not allow the domain to be a URL -- this is a very basic check
	if strings.HasSuffix(c.Domain, "/") || strings.HasPrefix(c.Domain, "http://") || strings.HasPrefix(c.Domain, "https://") {
		return nil, errors.New("invalid configuration: auth0 domain must not be a url or have a trailing slash")
	}

	// Default to the HTTPS scheme and reparse domain only configuration.
	if u, err = url.Parse("https://" + c.Domain + "/"); err != nil {
		return nil, errors.New("invalid configuration: specify auth0 domain of the configured tenant")
	}
	return u, nil
}

func (c AuthConfig) ClientCredentials() management.Option {
	return management.WithClientCredentials(context.Background(), c.ClientID, c.ClientSecret)
}

func (c MTLSConfig) Validate() error {
	if !c.Insecure {
		if c.CertPath == "" || c.PoolPath == "" {
			return errors.New("connecting over mTLS requires certs and cert pool")
		}
	}
	return nil
}

// DialOption returns a configured dial option which can be directly used in a
// grpc.Dial or grpc.DialContext call to connect using mTLS.
func (c MTLSConfig) DialOption(endpoint string) (opt grpc.DialOption, err error) {
	var (
		sz    *trust.Serializer
		certs *trust.Provider
		pool  trust.ProviderPool
	)

	if sz, err = trust.NewSerializer(false); err != nil {
		return nil, err
	}

	if certs, err = sz.ReadFile(c.CertPath); err != nil {
		return nil, err
	}

	if pool, err = sz.ReadPoolFile(c.PoolPath); err != nil {
		return nil, err
	}

	if opt, err = mtls.ClientCreds(endpoint, certs, pool); err != nil {
		return nil, err
	}

	return opt, nil
}

func (c EmailConfig) Validate() error {
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

func (c CacheConfig) Validate() error {
	if c.Enabled {
		if c.Size == 0 {
			return errors.New("invalid configuration: cache size must be greater than 0")
		}

		if c.Expiration == 0 {
			return errors.New("invalid configuration: cache expiration must be greater than 0")
		}
	}
	return nil
}

func validateURL(path string) (err error) {
	if path == "" {
		return errors.New("url is empty")
	}

	// URL should not have a trailing slash
	if strings.HasSuffix(path, "/") {
		return errors.New("url must not have a trailing slash")
	}

	if _, err = url.Parse(path); err != nil {
		return errors.New("url is not parseable")
	}

	return nil
}
