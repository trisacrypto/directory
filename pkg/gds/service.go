package gds

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/certman"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
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

	// Configure Sentry for error and performance monitoring
	if conf.Sentry.UseSentry() {
		if err = sentry.Init(conf.Sentry); err != nil {
			return nil, err
		}
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

		// Ensure that the Members API is created even in maintenance mode, so that
		// maintenance status replies are sent.
		if s.members, err = NewMembers(s); err != nil {
			return nil, err
		}

		// Ensure that the Admin API is created even in maintenance mode, so that
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

	// Create the Email Manager with SendGrid API client
	if s.email, err = emails.New(conf.Email); err != nil {
		return nil, err
	}

	// Create secret manager and connect to backend vault service
	if s.secret, err = secrets.New(conf.Secrets); err != nil {
		return nil, err
	}

	// Create the certificate manager
	if s.certman, err = certman.New(conf.CertMan, s.db, s.secret, s.email); err != nil {
		return nil, err
	}

	// Start the activity publisher
	if err = activity.Start(conf.Activity); err != nil {
		return nil, err
	}

	// Initialize the gRPC API services at the very end.
	if s.gds, err = NewGDS(s); err != nil {
		return nil, err
	}

	if s.admin, err = NewAdmin(s); err != nil {
		return nil, err
	}

	if s.members, err = NewMembers(s); err != nil {
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
	db      store.Store
	gds     *GDS
	admin   *Admin
	members *Members
	conf    config.Config
	certman certman.Service
	email   *emails.EmailManager
	secret  *secrets.SecretManager
	wg      sync.WaitGroup
	echan   chan error
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
		// Start the TRISA Members service (does not have unavailable status)
		if err = s.members.Serve(); err != nil {
			return err
		}

		// These services should not run in maintenance mode
		// Start the certificate manager go routine process
		s.wg = sync.WaitGroup{}
		s.certman.Run(&s.wg)

		// Start the backup manager go routine process
		// TODO: Refactor to use the wait group and shutdown gracefully
		go s.BackupManager(nil)
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
	log.Info().Msg("gracefully shutting down GDS service")

	// Shutdown the TRISADirectory service gracefully
	if err = s.gds.Shutdown(); err != nil {
		sentry.Error(nil).Err(err).Msg("could not shutdown TRISADirectory service")
	}

	// Shutdown the DirectoryAdministration service gracefully
	if err = s.admin.Shutdown(); err != nil {
		sentry.Error(nil).Err(err).Msg("could not shutdown DirectoryAdministration service")
	}

	if !s.conf.Maintenance {
		// Shutdown the TRISA members service gracefully
		if err = s.members.Shutdown(); err != nil {
			sentry.Error(nil).Err(err).Msg("could not shutdown TRISAMembers service")
		}

		// Stop the certificate manager
		s.certman.Stop()

		// Wait for all go routines to finish
		s.wg.Wait()

		// Close the database correctly
		if err = s.db.Close(); err != nil {
			sentry.Error(nil).Err(err).Msg("could not shutdown database")
		}
	}

	// Flush alert messages to Sentry
	if s.conf.Sentry.UseSentry() {
		sentry.Flush(2 * time.Second)
	}

	log.Debug().Msg("successfully shutdown GDS service")
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

// GetMembers returns the Members gRPC server
func (s *Service) GetMembers() *Members {
	return s.members
}

// GetAdmin returns the Admin server
func (s *Service) GetAdmin() *Admin {
	return s.admin
}

// GetConf returns a copy of the current configuration
func (s *Service) GetConf() config.Config {
	return s.conf
}

// GetSecretManager returns the secret manager
func (s *Service) GetSecretManager() *secrets.SecretManager {
	return s.secret
}
