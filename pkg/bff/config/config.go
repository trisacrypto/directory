package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/trisa/pkg/trisa/mtls"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
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
	CookieDomain string              `split_words:"true"`
	Auth0        AuthConfig
	TestNet      NetworkConfig
	MainNet      NetworkConfig
	Database     DatabaseConfig
	Sentry       sentry.Config
	processed    bool
}

// AuthConfig handles Auth0 configuration and authentication
type AuthConfig struct {
	Domain        string        `split_words:"true" required:"true"`
	Audience      string        `split_words:"true" required:"true"`
	ProviderCache time.Duration `split_words:"true" default:"5m"`
	ClientID      string        `split_words:"true"`
	ClientSecret  string        `split_words:"true"`
	Testing       bool          `split_words:"true" default:"false"` // If true a mock authenticator is used for testing
}

// NetworkConfig contains sub configurations for connecting to specific GDS and members
// services.
type NetworkConfig struct {
	Admin     AdminConfig
	Directory DirectoryConfig
	Members   MembersConfig
}

// DirectoryConfig is a generic configuration for connecting to a GDS service.
type DirectoryConfig struct {
	Insecure bool          `split_words:"true" default:"true"`
	Endpoint string        `split_words:"true" required:"true"`
	Timeout  time.Duration `split_words:"true" default:"10s"`
}

// AdminConfig is a configuration for connecting to an Admin service.
type AdminConfig struct {
	Endpoint string `split_words:"true" required:"true"`
	// Audience and TokenKeys should match the Admin server configuration, since the
	// BFF uses these parameters to sign its own JWT tokens.
	Audience  string            `split_words:"true"`
	TokenKeys map[string]string `split_words:"true" required:"true"`
	Testing   bool              `split_words:"true" default:"false"`
}

// MembersConfig is a configuration for connecting to a members service.
type MembersConfig struct {
	Endpoint string        `split_words:"true" required:"true"`
	Timeout  time.Duration `split_words:"true" default:"10s"`
	MTLS     MTLSConfig
}

type DatabaseConfig struct {
	URL           string `split_words:"true" required:"true"`
	ReindexOnBoot bool   `split_words:"true" default:"false"`
	MTLS          MTLSConfig
}

type MTLSConfig struct {
	Insecure bool   `split_words:"true"`
	CertPath string `split_words:"true"`
	PoolPath string `split_words:"true"`
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
func (c Config) Validate() (err error) {
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

	if err = c.Sentry.Validate(); err != nil {
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

func (c DatabaseConfig) Validate() error {
	if err := c.MTLS.Validate(); err != nil {
		return fmt.Errorf("invalid database configuration: %w", err)
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
	return management.WithClientCredentials(c.ClientID, c.ClientSecret)
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
