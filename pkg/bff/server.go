package bff

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/auth0/go-auth0/management"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/bff/db"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
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

// New creates a new BFF server from the specified configuration.
func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment if config is empty
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

	// Connect to the TestNet and MainNet directory services and database if we're not
	// in maintenance or testing mode (in testing mode, the connection will be manual).
	if !s.conf.Maintenance && s.conf.Mode != gin.TestMode {
		if s.testnet.gds, err = ConnectGDS(conf.TestNet); err != nil {
			return nil, fmt.Errorf("could not connect to the TestNet directory service: %s", err)
		}
		log.Debug().Str("endpoint", conf.TestNet.Endpoint).Str("network", "testnet").Msg("connected to GDS")

		if s.mainnet.gds, err = ConnectGDS(conf.MainNet); err != nil {
			return nil, fmt.Errorf("could not connect to the MainNet: %s", err)
		}
		log.Debug().Str("endpoint", conf.MainNet.Endpoint).Str("network", "mainnet").Msg("connected to GDS")

		if s.testnet.members, err = ConnectMembers(conf.TestNet); err != nil {
			return nil, fmt.Errorf("could not connect to the TestNet members service: %s", err)
		}
		log.Debug().Str("endpoint", conf.TestNet.Endpoint).Str("network", "testnet").Msg("connected to members")

		if s.mainnet.members, err = ConnectMembers(conf.MainNet); err != nil {
			return nil, fmt.Errorf("could not connect to the MainNet members service: %s", err)
		}
		log.Debug().Str("endpoint", conf.MainNet.Endpoint).Str("network", "mainnet").Msg("connected to members")

		if s.db, err = db.Connect(s.conf.Database); err != nil {
			return nil, fmt.Errorf("could not connect to trtl database: %s", err)
		}
		log.Debug().Str("dsn", s.conf.Database.URL).Bool("insecure", s.conf.Database.Insecure).Msg("connected to trtl database")

		if s.auth0, err = management.New(s.conf.Auth0.Domain, s.conf.Auth0.ClientCredentials()); err != nil {
			return nil, fmt.Errorf("could not connect to auth0 management api: %s", err)
		}
		log.Debug().Str("domain", s.conf.Auth0.Domain).Msg("connected to auth0")
	}

	if conf.Auth0.Testing {
		// Setup the auth mock server
		if s.mockAuth, err = authtest.New(); err != nil {
			return nil, fmt.Errorf("could not create auth mock server: %s", err)
		}
	}

	// Create the router
	gin.SetMode(conf.Mode)
	s.router = gin.New()
	if err = s.setupRoutes(); err != nil {
		return nil, err
	}

	// Create the http server
	s.srv = &http.Server{
		Addr:         s.conf.BindAddr,
		Handler:      s.router,
		ErrorLog:     nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return s, nil
}

type NetworkClient struct {
	gds     gds.TRISADirectoryClient
	members members.TRISAMembersClient
}

type Server struct {
	sync.RWMutex
	conf     config.Config
	srv      *http.Server
	router   *gin.Engine
	testnet  NetworkClient
	mainnet  NetworkClient
	db       *db.DB
	mockAuth *authtest.Server
	auth0    *management.Management
	started  time.Time
	healthy  bool
	url      string
	echan    chan error
}

// Serve API requests on the specified address.
func (s *Server) Serve() (err error) {
	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.echan <- s.Shutdown()
	}()

	// Set the health of the service to true unless we're in maintenance mode.
	// The server should still start so that it can return unavailable to requests.
	s.SetHealth(!s.conf.Maintenance)
	if s.conf.Maintenance {
		log.Warn().Msg("starting server in maintenance mode")
	}

	// Create a socket to listen on so that we can infer the final URL (e.g. if the
	// BindAddr is 127.0.0.1:0 for testing, a random port will be assigned, manually
	// creating the listener will allow us to determine which port).
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		s.echan <- err
	}

	// Set the URL from the listener
	s.SetURL("http://" + sock.Addr().String())
	s.started = time.Now()

	// Listen for HTTP requests on the specified address and port
	go func() {
		if err = s.srv.Serve(sock); err != nil && err != http.ErrServerClosed {
			s.echan <- err
		}
	}()

	log.Info().
		Str("listen", s.url).
		Str("version", pkg.Version()).
		Msg("gds bff server started")

	// Listen for any errors that might have occurred and wait for all go routines to stop
	if err = <-s.echan; err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down")

	// Flush the Sentry log before shutting down
	defer sentry.Flush(2 * time.Second)

	s.SetHealth(false)
	s.srv.SetKeepAlivesEnabled(false)

	// Require shutdown in 30 seconds without blocking
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = s.srv.Shutdown(ctx); err != nil {
		return err
	}

	// Shut down maintenance mode systems
	if !s.conf.Maintenance {
		if s.db != nil {
			if err = s.db.Close(); err != nil {
				log.Error().Err(err).Msg("could not shutdown trtl db connection")
			}
		}
	}

	if s.mockAuth != nil {
		s.mockAuth.Close()
	}

	log.Debug().Msg("successfully shutdown server")
	return nil
}

func (s *Server) SetHealth(health bool) {
	s.Lock()
	s.healthy = health
	s.Unlock()
	log.Debug().Bool("healthy", health).Msg("server health set")
}

func (s *Server) SetURL(url string) {
	s.Lock()
	s.url = url
	s.Unlock()
	log.Debug().Str("url", url).Msg("server url set")
}

func (s *Server) setupRoutes() (err error) {
	var (
		authenticator gin.HandlerFunc
		tags          gin.HandlerFunc
		tracing       gin.HandlerFunc
		bffTags       map[string]string
	)

	if s.conf.Auth0.Testing {
		// Instantiate the mock authentication middleware
		if authenticator, err = auth.Authenticate(s.mockAuth.Config(), auth.WithHTTPClient(s.mockAuth.Client())); err != nil {
			return err
		}
	} else {
		// Instantiate authentication middleware
		if authenticator, err = auth.Authenticate(s.conf.Auth0); err != nil {
			return err
		}
	}

	// Instantiate user info middleware
	var userinfo gin.HandlerFunc
	if userinfo, err = auth.UserInfo(s.conf.Auth0); err != nil {
		return err
	}

	if s.conf.Sentry.UseSentry() {
		bffTags = map[string]string{"service": "bff"}
		tags = sentry.UseTags(bffTags)
	}

	if s.conf.Sentry.UsePerformanceTracking() {
		tracing = sentry.TrackPerformance(bffTags)
	}

	// Application Middleware
	// NOTE: ordering is very important to how middleware is handled.
	middlewares := []gin.HandlerFunc{
		// Logging should be outside so we can record the complete latency of requests.
		// NOTE: logging panics will not recover.
		logger.GinLogger("bff"),

		// Panic recovery middleware; note: gin middleware needs to be added before sentry
		gin.Recovery(),
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
		}),

		// Add searchable tags to the sentry context.
		tags,

		// Tracing helps us with our peformance metrics and should be as early in the
		// chain as possible. It is after recovery to ensure trace panics recover.
		tracing,

		// CORS configuration allows the front-end to make cross-origin requests.
		cors.New(cors.Config{
			AllowOrigins:     s.conf.AllowOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-TOKEN", "sentry-trace"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),

		// Maintenance mode handling - does not require authentication.
		s.Available(),

		// Authentication happens as late as possible; all middleware after this should
		// require a user context; if it doesn't, it should come before authentication.
		authenticator,
	}

	// Add the middleware to the router
	for _, middleware := range middlewares {
		if middleware != nil {
			s.router.Use(middleware)
		}
	}

	// Add the v1 API routes
	v1 := s.router.Group("/v1")
	{
		// Heartbeat route (no authentication required)
		v1.GET("/status", s.Status)

		// GDS public routes (no authentication required)
		v1.GET("/lookup", s.Lookup)
		v1.POST("/register/:network", s.Register)
		v1.GET("/verify", s.VerifyContact)
		v1.POST("/users/login", userinfo, s.Login)

		// Members authenticated routes
		v1.GET("/overview", auth.Authorize("read:vasp"), s.Overview)
		v1.GET("/certificates", auth.Authorize("read:vasp"), s.Certificates)
	}

	// NotFound and NotAllowed routes
	s.router.NoRoute(api.NotFound)
	s.router.NoMethod(api.NotAllowed)
	return nil
}

//===========================================================================
// Accessors - used primarily for testing
//===========================================================================

// SetGDSClients allows tests to set a bufconn client to a mock GDS server.
func (s *Server) SetGDSClients(testnet, mainnet gds.TRISADirectoryClient) {
	s.testnet.gds = testnet
	s.mainnet.gds = mainnet
}

// SetMembersClients allows tests to set a bufconn client to a mock members server.
func (s *Server) SetMembersClients(testnet, mainnet members.TRISAMembersClient) {
	s.testnet.members = testnet
	s.mainnet.members = mainnet
}

// SetDB allows tests to set a bufconn client to a mock trtl server.
func (s *Server) SetDB(db *db.DB) {
	s.db = db
}

// GetConf returns a copy of the current configuration.
func (s *Server) GetConf() config.Config {
	return s.conf
}

// GetRouter returns the Gin API router for testing purposes.
func (s *Server) GetRouter() http.Handler {
	return s.router
}

// GetTestNet returns the TestNet directory client for testing purposes.
func (s *Server) GetTestNetGDS() gds.TRISADirectoryClient {
	return s.testnet.gds
}

// GetMainNet returns the MainNet directory client for testing purposes.
func (s *Server) GetMainNetGDS() gds.TRISADirectoryClient {
	return s.mainnet.gds
}

// GetTestNet returns the TestNet members client for testing purposes.
func (s *Server) GetTestNetMembers() members.TRISAMembersClient {
	return s.testnet.members
}

// GetMainNet returns the MainNet members client for testing purposes.
func (s *Server) GetMainNetMembers() members.TRISAMembersClient {
	return s.mainnet.members
}

func (s *Server) GetMockAuth() *authtest.Server {
	return s.mockAuth
}

// GetURL returns the URL that the server can be reached if it has been started. This
// accessor is primarily used to create a test client.
func (s *Server) GetURL() string {
	s.RLock()
	defer s.RUnlock()
	return s.url
}
