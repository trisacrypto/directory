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
	trtlstore "github.com/trisacrypto/directory/pkg/gds/store/trtl"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/mock"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"google.golang.org/grpc"
)

// NewMock creates and returns a mocked Service for testing, using values provided in
// the config. The config should contain values specific to testing as the mock method
// only mocks at the top level of the service, lower level mocks such as mocking the
// secret manager or email service must be implemented with configuration. Use
// MockConfig to ensure a configuration is generated that fully mocks the service.
func NewMock(conf config.Config, trtlConn *grpc.ClientConn) (s *Service, err error) {
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

	if trtlConn != nil {
		// The Trtl store mock requires a bufconn connection
		if svc.db, err = trtlstore.NewMock(trtlConn); err != nil {
			return nil, err
		}
	} else {
		if svc.db, err = store.Open(conf.Database); err != nil {
			return nil, err
		}
	}

	if svc.gds, err = NewGDS(svc); err != nil {
		return nil, err
	}

	if svc.members, err = NewMembers(svc); err != nil {
		return nil, err
	}

	// TODO: Should we be using NewAdmin() here?
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

	if conf.Sectigo.Testing {
		if err = mock.Start(conf.Sectigo.Profile); err != nil {
			return nil, err
		}
	}

	if svc.certs, err = sectigo.New(conf.Sectigo); err != nil {
		return nil, err
	}

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
				GoogleAudience:         "http://localhost",
				AuthorizedEmailDomains: []string{"gds.dev"},
			},
			TokenKeys: nil,
		},
		Members: config.MembersConfig{
			Enabled:  true,
			Insecure: true,
		},
		Database: config.DatabaseConfig{
			URL:           "leveldb:///testdata/testdb",
			ReindexOnBoot: false,
		},
		Sectigo: sectigo.Config{
			Username: "foo",
			Password: "supersecretsquirrel",
			Profile:  "CipherTrace EE",
			Testing:  true,
		},
		Email: config.EmailConfig{
			ServiceEmail:         "GDS <service@gds.dev>",
			AdminEmail:           "GDS Admin <admin@gds.dev>",
			SendGridAPIKey:       "notarealsendgridapikey",
			DirectoryID:          "gds.dev",
			VerifyContactBaseURL: "https://gds.dev/verify",
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
