package trtl

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/rotationalio/honu"
	replication "github.com/rotationalio/honu/replica"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	prom "github.com/trisacrypto/directory/pkg/trtl/metrics"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"github.com/trisacrypto/directory/pkg/trtl/replica"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

// A Trtl server implements the following services
// 1. A database service for interacting with a database
// 2. A peers management service for interacting with remote peers
// 3. A replication service which implements auto-adapting anti-entropy replication.
type Server struct {
	srv     *grpc.Server         // The gRPC server that listens on its own independent port
	conf    config.Config        // Configuration for the trtl server
	db      *honu.DB             // Database connection for managing objects
	trtl    *TrtlService         // Service for interacting with a Honu database
	peers   *PeerService         // Service for managing remote peers
	replica *replica.Service     // Service that handles anti-entropy replication
	metrics *prom.MetricsService // Service for Prometheus metrics
	backup  *BackupManager       // Manages backups of the trtl database
	monitor *Monitor             // Monitors the storage usage of the trtl database
	started time.Time            // The timestamp that the server was started (for uptime)
	echan   chan error           // Channel for receiving errors from the gRPC server
}

// New creates a new trtl server given a configuration.
func New(conf config.Config) (s *Server, err error) {
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

	// Configure Sentry
	if conf.Sentry.UseSentry() {
		if err = sentry.Init(conf.Sentry); err != nil {
			return nil, err
		}
	}

	// Create the server and prepare to serve
	s = &Server{
		conf:  conf,
		echan: make(chan error, 1),
	}

	opts := make([]grpc.ServerOption, 0, 2)
	if !conf.MTLS.Insecure {
		// Add mTLS configuration if enabled
		tlsConf, err := conf.MTLS.ParseTLSConfig()
		if err != nil {
			return nil, fmt.Errorf("could not parse TLS config: %v", err)
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConf)))
	} else {
		log.Warn().Msg("trtl starting without mTLS enabled")
	}

	// NOTE: It appears this must happen outside the struct initialization of the Server
	// or else the UnaryInterceptor doesn't capture conf when it when it creates the closure
	opts = append(opts, grpc.ChainUnaryInterceptor(s.UnaryInterceptors()...))
	opts = append(opts, grpc.ChainStreamInterceptor(s.StreamInterceptors()...))
	s.srv = grpc.NewServer(opts...)

	// NOTE: if we are *not* in maintenance mode, we must open the database before we
	// initialize the Honu and Peer mgmt services, else we'll get panics on nil dbs
	// This is not the case for maintenance mode; which can proceed by initializing the
	// the services, allowing them to respond "unavailable", but which will prevent the
	// dbs in question from being engaged.
	if !s.conf.Maintenance {
		// Open a connection to the Honu wrapped database
		if s.db, err = honu.Open(conf.Database.URL, conf.GetHonuConfig()); err != nil {
			return nil, fmt.Errorf("honu error: %v", err)
		}

		// Initialize the backup manager
		if s.backup, err = NewBackupManager(s.conf.Backup, s.db); err != nil {
			return nil, err
		}

		// Initialize the database monitor
		if s.monitor, err = NewMonitor(s.conf.Metrics, s.db); err != nil {
			return nil, err
		}
	}

	// Initialize the Honu service
	if s.trtl, err = NewTrtlService(s); err != nil {
		return nil, err
	}
	pb.RegisterTrtlServer(s.srv, s.trtl)

	// Initialize the Peer Management service
	if s.peers, err = NewPeerService(s); err != nil {
		return nil, err
	}
	peers.RegisterPeerManagementServer(s.srv, s.peers)

	// Initialize the Replica service
	if s.replica, err = replica.New(s.conf, s.db, replicatedNamespaces); err != nil {
		return nil, err
	}
	replication.RegisterReplicationServer(s.srv, s.replica)

	// Initialize Metrics service for Prometheus
	if s.metrics, err = prom.New(conf.Metrics); err != nil {
		return nil, err
	}
	return s, nil
}

// Serve gRPC requests on the specified bind address.
func (t *Server) Serve() (err error) {
	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		t.echan <- t.Shutdown()
	}()

	// Run management routines only if we're not in maintenance mode
	if t.conf.Maintenance {
		log.Warn().Msg("starting trtl in maintenance mode")
	} else {
		// These services should not run in maintenance mode
		// Run the Gossip background routine - passing in the stop channel; it's a bit
		// awkward to create this channel here, but it follows the pattern for our other
		// background routines that need stop channels injected directly from tests.
		go t.replica.AntiEntropy(make(chan struct{}, 1))

		// Run the backup manager if enabled
		go t.backup.Run()

		// Run the monitor if enabeld
		go t.monitor.Run()
	}

	// If metrics are enabled, start Prometheus metrics server as separate go routine
	if t.conf.Metrics.Enabled {
		t.metrics.Serve()
	} else {
		log.Warn().Msg("trtl prometheus metrics server disabled")
	}

	// Listen for TCP requests
	var sock net.Listener
	if sock, err = net.Listen("tcp", t.conf.BindAddr); err != nil {
		log.Error().Err(err).Str("bindaddr", t.conf.BindAddr).Msg("could not listen on given bindaddr")
		return err
	}

	// Run the gRPC server
	go t.Run(sock)
	log.Info().Str("listen", t.conf.BindAddr).Msg("trtl server started")

	// Set the timestamp that the server was started now that it is booted up
	t.started = time.Now()

	// The server go routine is started so return nil error (any server errors will be
	// sent on the error channel).
	if err = <-t.echan; err != nil {
		return err
	}
	return nil
}

// Run the gRPC server. This method is extracted from the Serve function so that it can
// be run in its own go routine and to allow tests to Run a bufconn server without
// starting a live server with all of the various go routines and channels running.
func (t *Server) Run(sock net.Listener) {
	defer sock.Close()
	if err := t.srv.Serve(sock); err != nil {
		t.echan <- err
	}
}

// Shutdown the trtl server gracefully.
func (t *Server) Shutdown() (err error) {
	errs := make([]error, 0)
	log.Info().Msg("gracefully shutting down trtl server")
	t.srv.GracefulStop()

	// Shutdown the backup manager
	if t.conf.Backup.Enabled {
		if err = t.backup.Shutdown(); err != nil {
			log.Error().Err(err).Msg("could not shutdown backup manager")
			errs = append(errs, err)
		}
	}

	// Shutdown the Prometheus metrics server and the monitor
	if t.conf.Metrics.Enabled {
		if err = t.monitor.Shutdown(); err != nil {
			log.Error().Err(err).Msg("Could not shutdown database monitor")
			errs = append(errs, err)
		}

		if err = t.metrics.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("could not shutdown prometheus metrics server")
			errs = append(errs, err)
		}
	}

	// Stop the anti-entropy routine.
	if err = t.replica.Shutdown(); err != nil {
		log.Error().Err(err).Msg("could not shutdown anti-entropy routine")
		errs = append(errs, err)
	}

	// If we're in maintenance mode db will be nil, so check if available to avoid panic
	if t.db != nil {
		if err = t.db.Close(); err != nil {
			log.Error().Err(err).Msg("could not close database")
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		log.Debug().Msg("did not successfully shutdown trtl server")
		sentry.Error(nil).Errs(errs).Msg("trtl shutdown errors")
		return fmt.Errorf("%d shutdown errors occurred", len(errs))
	}

	log.Debug().Msg("successful shutdown of trtl server")
	return nil
}

//===========================================================================
// Accessors - used primarily for testing
//===========================================================================

// GetDB returns the underlying Honu database used by all sub-services.
func (t *Server) GetDB() *honu.DB {
	return t.db
}
