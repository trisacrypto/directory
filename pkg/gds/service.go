package gds

import (
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/sectigo/mock"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

func init() {
	// Initialize zerolog with GCP logging requirements
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg

	// Add the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
}

// New creates a TRISA Directory Service with the specified configuration and prepares
// it to listen for and serve GRPC requests.
func New(conf config.Config) (s *Service, err error) {
	// Load the default configuration from the environment
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Set the global level
	zerolog.SetGlobalLevel(conf.GetLogLevel())

	// Set human readable logging if specified
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Create the server and prepare to serve
	s = &Service{conf: conf, echan: make(chan error, 1)}

	// Stop configuration at this point for maintenance mode (no error)
	if s.conf.Maintenance {
		// Ensure that the GDS service is created even in maintenance mode, so that
		// maintenance status replies are sent.
		if s.gds, err = NewGDS(s); err != nil {
			return nil, err
		}

		// Ensure that the Admin API is created even in maintenace mode, so that
		// maintenance status replies are sent.
		if s.admin, err = NewAdmin(s); err != nil {
			return nil, err
		}

		return s, nil
	}

	// Everything that follows here assumes we're not in maintenance mode.
	if s.db, err = store.Open(conf.Database); err != nil {
		return nil, err
	}

	// Create the Sectigo API client
	if s.certs, err = sectigo.New(conf.Sectigo); err != nil {
		return nil, err
	}

	// Start mocked Sectigo server if testing is enabled
	if conf.Sectigo.Testing {
		if s.mock, err = mock.New(); err != nil {
			return nil, err
		}
	}

	// Ensure the certificate storage can be reached
	if _, err = s.getCertStorage(); err != nil {
		return nil, err
	}

	// Create the Email Manager with SendGrid API client
	if s.email, err = emails.New(conf.Email); err != nil {
		return nil, err
	}

	// Create secret manager and connect to backend vault service
	if s.secret, err = secrets.New(conf.Secrets); err != nil {
		return nil, err
	}

	// Initialize the gRPC API services at the very end.
	if s.gds, err = NewGDS(s); err != nil {
		return nil, err
	}

	if s.admin, err = NewAdmin(s); err != nil {
		return nil, err
	}

	return s, nil
}

// Service defines the entirety of the TRISA Global Directory Service including the GDS
// server that handles TRISA requests, the Admin server that handles administrative
// interactions, as well as the smaller routines and managers to handle email, secrets,
// backups, and certificates.
// E.g. this is the parent service that coordinates all subservices.
type Service struct {
	db     store.Store
	gds    *GDS
	admin  *Admin
	conf   config.Config
	certs  *sectigo.Sectigo
	mock   *mock.Server
	email  *emails.EmailManager
	secret *secrets.SecretManager
	echan  chan error
}

// Serve GRPC requests on the specified addresses and all internal servers.
func (s *Service) Serve() (err error) {
	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.echan <- s.Shutdown()
	}()

	// Run management routines only if we're not in maintenance mode
	if s.conf.Maintenance {
		log.Warn().Msg("starting server in maintenance mode")
	} else {
		// These services should not run in maintenance mode
		// Start the certificate manager go routine process
		go s.CertManager()

		// Start the backup manager go routine process
		go s.BackupManager()
	}

	// The TRISADirectoryService service can run in maintenance mode
	if err = s.gds.Serve(); err != nil {
		return err
	}

	// The DirectoryAdministrationService can run in maintenance mode
	if err = s.admin.Serve(); err != nil {
		return err
	}

	// Listen for any errors that might have occurred and wait for all go routines to finish
	if err = <-s.echan; err != nil {
		return err
	}
	return nil
}

// Shutdown the TRISA Directory Service gracefully
func (s *Service) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down")

	// Shutdown the TRISADirectory service gracefully
	if err = s.gds.Shutdown(); err != nil {
		log.Error().Err(err).Msg("could not shutdown TRISADirectory service")
	}

	// Shutdown the DirectoryAdministration service gracefully
	if err = s.admin.Shutdown(); err != nil {
		log.Error().Err(err).Msg("could not shutdown DirectoryAdministration service")
	}

	// Shutdown the Sectigo mock server
	if s.mock != nil {
		s.mock.Close()
	}

	if !s.conf.Maintenance {
		// Close the database correctly
		if err = s.db.Close(); err != nil {
			log.Error().Err(err).Msg("could not shutdown database")
		}
	}

	return nil
}

//===========================================================================
// Accessors - used primarily for testing
//===========================================================================

// GetStore returns the underlying database store used by all sub-services.
func (s *Service) GetStore() store.Store {
	return s.db
}

// GetGDS returns the GDS gRPC server
func (s *Service) GetGDS() *GDS {
	return s.gds
}

// GetAdmin returns the Admin server
func (s *Service) GetAdmin() *Admin {
	return s.admin
}

// GetConf returns a copy of the current configuration
func (s *Service) GetConf() config.Config {
	return s.conf
}
