package gds

import (
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	"github.com/trisacrypto/directory/pkg/gds/logger"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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
	zerolog.SetGlobalLevel(zerolog.Level(conf.LogLevel))

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
		return s, nil
	}

	// Everything that follows here assumes we're not in maintenance mode.
	if s.db, err = store.Open(conf.Database); err != nil {
		return nil, err
	}

	// If replication is enabled, then add a global version manager to the store.
	if s.conf.Replica.Enabled {
		var versionManager *global.VersionManager
		if versionManager, err = global.New(s.conf.Replica); err != nil {
			return nil, err
		}
		if err = s.db.WithVersionManager(versionManager); err != nil {
			return nil, err
		}
	}

	// Create the Sectigo API client
	if s.certs, err = sectigo.New(conf.Sectigo.Username, conf.Sectigo.Password); err != nil {
		return nil, err
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
	if s.secret, err = NewSecretManager(conf.Secrets); err != nil {
		return s, nil
	}

	// Initialize the gRPC API services at the very end.
	if s.gds, err = NewGDS(s); err != nil {
		return nil, err
	}

	if s.admin, err = NewAdmin(s); err != nil {
		return nil, err
	}

	if s.replica, err = NewReplica(s); err != nil {
		return nil, err
	}

	return s, nil
}

// Service defines the entirety of the TRISA Global Directory Service including the GDS
// server that handles TRISA requests, the Admin server that handles administrative
// interactions, the Replica server that performs anti-entropy, as well as the smaller
// routines and managers to handle email, secrets, backups, and certificates. E.g. this
// is the parent service that coordinates all subservices.
type Service struct {
	db      store.Store
	gds     *GDS
	admin   *Admin
	replica *Replica
	conf    config.Config
	certs   *sectigo.Sectigo
	email   *emails.EmailManager
	secret  *SecretManager
	echan   chan error
}

// Serve GRPC requests on the specified address.
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
		// Start the certificate manager go routine process
		go s.CertManager()

		// Start the backup manager go routine process
		go s.BackupManager()

		// Start the admin service
		if err = s.admin.Serve(); err != nil {
			return err
		}

		// Start the replica service
		if err = s.replica.Serve(); err != nil {
			return err
		}
	}

	// The only service that runs in maintenance mode is the TRISADirectoryService,
	// which responds to Status RPC requests that the server is in maintenance mode.
	if err = s.gds.Serve(); err != nil {
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

	if !s.conf.Maintenance {
		// Shutdown the DirectoryAdministration service gracefully
		if err = s.admin.Shutdown(); err != nil {
			log.Error().Err(err).Msg("could not shutdown DirectoryAdministration service")
		}

		// Shutdown the ReplicationServer gracefully
		if err = s.replica.Shutdown(); err != nil {
			log.Error().Err(err).Msg("could not shutdown Replication service")
		}

		// Close the database correctly
		if err = s.db.Close(); err != nil {
			log.Error().Err(err).Msg("could not shutdown database")
		}
	}

	return nil
}
