package gds

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc"
)

// NewMock creates and returns a mocked Service for testing, using values provided in
// the config. The config should contain values specific to testing as the mock method
// only mocks at the top level of the service, lower level mocks such as mocking the
// secret manager or email service must be implemented with configuration. Use
// MockConfig to ensure a configuration is generated that fully mocks the service.
func NewMock(conf config.Config) (s *Service, err error) {
	// Set the global level
	zerolog.SetGlobalLevel(conf.GetLogLevel())

	// Set human readable logging if specified
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	svc := &Service{
		conf: conf,
	}
	if svc.email, err = emails.New(conf.Email); err != nil {
		return nil, err
	}
	if svc.secret, err = secrets.NewMock(conf.Secrets); err != nil {
		return nil, err
	}
	if svc.db, err = store.Open(conf.Database); err != nil {
		return nil, err
	}

	gds := &GDS{
		svc:  svc,
		conf: &svc.conf.GDS,
		db:   svc.db,
	}
	gds.srv = grpc.NewServer(grpc.UnaryInterceptor(svc.serverInterceptor))
	api.RegisterTRISADirectoryServer(gds.srv, gds)
	svc.gds = gds

	admin := &Admin{
		svc:  svc,
		conf: &svc.conf.Admin,
		db:   svc.db,
	}
	if admin.tokens, err = tokens.MockTokenManager(); err != nil {
		return nil, err
	}
	gin.SetMode(admin.conf.Mode)
	admin.router = gin.New()
	if err = admin.setupRoutes(); err != nil {
		return nil, err
	}
	svc.admin = admin
	return svc, nil
}

// MockConfig returns a configuration that ensures the service will operate in a fully
// mocked way with all testing parameters set correctly. The config is returned directly
// for required modifications, such as pointing the database path to a fixtures path.
func MockConfig() config.Config {
	conf := config.Config{
		DirectoryID: "gds.dev",
		SecretKey:   "supersecretsquirrel",
		Maintenance: false,
		LogLevel:    logger.LevelDecoder(zerolog.WarnLevel),
		ConsoleLog:  true,
		GDS: config.GDSConfig{
			Enabled:  false,
			BindAddr: "",
		},
		Admin: config.AdminConfig{
			Enabled:      false,
			BindAddr:     "",
			Mode:         gin.TestMode,
			AllowOrigins: []string{"http://127.0.0.1", "http://localhost", "http://admin.gds.dev"},
			CookieDomain: "admin.gds.dev",
			Audience:     "http://api.admin.gds.dev",
			Oauth: config.OauthConfig{
				GoogleAudience:         "4284607864536-notarealgoogleaudience.apps.googleusercontent.com",
				AuthorizedEmailDomains: []string{"gds.dev"},
			},
			TokenKeys: nil,
		},
		Database: config.DatabaseConfig{
			URL:           "leveldb:///testdata/testdb",
			ReindexOnBoot: false,
		},
		Sectigo: config.SectigoConfig{
			Username: "foo",
			Password: "supersecretsquirrel",
			Profile:  "CipherTrace EE",
		},
		Email: config.EmailConfig{
			ServiceEmail:         "GDS <service@gds.dev>",
			AdminEmail:           "GDS Admin <admin@gds.dev>",
			SendGridAPIKey:       "notarealsendgridapikey",
			DirectoryID:          "gds.dev",
			VerifyContactBaseURL: "https://gds.dev/verify-contact",
			AdminReviewBaseURL:   "https://admin.gds.dev/vasps/",
			Testing:              true,
		},
		CertMan: config.CertManConfig{
			Interval: 24 * time.Hour,
			Storage:  "testdata/certs",
		},
		Backup: config.BackupConfig{
			Enabled:  false,
			Interval: 24 * time.Hour,
			Storage:  "testdata/backups",
			Keep:     1,
		},
		Secrets: config.SecretsConfig{
			Credentials: "",
			Project:     "",
			Testing:     true,
		},
	}

	var err error
	if conf, err = conf.Mark(); err != nil {
		panic(err)
	}
	return conf
}
