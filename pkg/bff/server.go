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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/config"
	docs "github.com/trisacrypto/directory/pkg/bff/docs"
	"github.com/trisacrypto/directory/pkg/bff/emails"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils/cache"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"google.golang.org/grpc"
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
	if !s.conf.Maintenance {
		if s.conf.Mode != gin.TestMode {
			if s.testnetDB, err = store.Open(conf.TestNet.Database); err != nil {
				return nil, err
			}

			if s.mainnetDB, err = store.Open(conf.MainNet.Database); err != nil {
				return nil, err
			}

			if s.testnetGDS, err = ConnectGDS(s.conf.TestNet); err != nil {
				return nil, fmt.Errorf("could not connect to testnet: %w", err)
			}

			if s.mainnetGDS, err = ConnectGDS(s.conf.MainNet); err != nil {
				return nil, fmt.Errorf("could not connect to mainnet: %w", err)
			}

			if s.db, err = store.Open(s.conf.Database); err != nil {
				return nil, fmt.Errorf("could not connect to trtl database: %w", err)
			}
			log.Debug().Str("dsn", s.conf.Database.URL).Bool("insecure", s.conf.Database.Insecure).Msg("connected to trtl database")

			// Initialize Activity Subscriber if enabled
			if s.conf.Activity.Enabled {
				if s.activity, err = NewActivitySubscriber(s.conf.Activity, s.db); err != nil {
					return nil, fmt.Errorf("could not create activity subscriber: %w", err)
				}
			}
		}

		if s.email, err = emails.New(conf.Email); err != nil {
			return nil, fmt.Errorf("could not connect to email service: %w", err)
		}

		if s.auth0, err = auth.NewManagementClient(s.conf.Auth0); err != nil {
			return nil, fmt.Errorf("could not connect to auth0 management api: %w", err)
		}

		// Initialize the user cache or use a no-op cache if disabled
		if s.users, err = cache.New(s.conf.UserCache); err != nil {
			return nil, fmt.Errorf("could not initialize user cache: %w", err)
		}

		log.Debug().Str("domain", s.conf.Auth0.Domain).Msg("connected to auth0")
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

// ConnectGDS creates a unified client to the TRISA Directory Service and TRISA
// members service specified in the configuration. This method is used to connect to
// both the TestNet and the MainNet so we can maintain separate clients for each.
func ConnectGDS(conf config.NetworkConfig) (_ GlobalDirectoryClient, err error) {
	client := &GDSClient{}

	if err = client.ConnectGDS(conf.Directory); err != nil {
		return nil, fmt.Errorf("could not connect to directory service: %s", err)
	}

	if conf.Members.MTLS.Insecure {
		if err = client.ConnectMembers(conf.Members); err != nil {
			return nil, fmt.Errorf("could not connect to insecure members service: %s", err)
		}
	} else {
		var mtls grpc.DialOption
		if mtls, err = conf.Members.MTLS.DialOption(conf.Members.Endpoint); err != nil {
			return nil, fmt.Errorf("could not create dial option for mTLS: %s", err)
		}

		if err = client.ConnectMembers(conf.Members, mtls); err != nil {
			return nil, fmt.Errorf("could not connect to members service with mTLS: %s", err)
		}
	}

	log.Info().
		Str("directory", conf.Directory.Endpoint).
		Bool("directory_insecure", conf.Directory.Insecure).
		Str("members", conf.Members.Endpoint).
		Bool("members_insecure", conf.Members.MTLS.Insecure).
		Msg("connected to the GDS")
	return client, nil
}

type Server struct {
	sync.RWMutex
	conf       config.Config
	srv        *http.Server
	router     *gin.Engine
	testnetDB  store.Store
	mainnetDB  store.Store
	testnetGDS GlobalDirectoryClient
	mainnetGDS GlobalDirectoryClient
	db         store.Store
	auth0      *management.Management
	email      *emails.EmailManager
	users      cache.Cache
	activity   *ActivitySubscriber
	started    time.Time
	healthy    bool
	url        string
	echan      chan error
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

	// Start the activity subscriber
	if s.conf.Activity.Enabled {
		if err = s.activity.Run(&sync.WaitGroup{}); err != nil {
			return fmt.Errorf("could not start activity subscriber: %w", err)
		}
	}

	// Create a socket to listen on so that we can infer the final URL (e.g. if the
	// BindAddr is 127.0.0.1:0 for testing, a random port will be assigned, manually
	// creating the listener will allow us to determine which port).
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on bind addr %s: %w", s.conf.BindAddr, err)
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

	// Stop the activity subscriber
	if s.activity != nil {
		s.activity.Stop()
	}

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
				sentry.Error(nil).Err(err).Msg("could not shutdown trtl db connection")
			}
		}
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

// @title BFF API
// @version 1.0
// @description BFF server which supports the GDS user frontend
// @BasePath /v1
func (s *Server) setupRoutes() (err error) {
	var (
		authenticator gin.HandlerFunc
		sentryRecover gin.HandlerFunc
		tags          gin.HandlerFunc
		tracing       gin.HandlerFunc
		bffTags       map[string]string
	)

	// Instantiate authentication middleware
	if authenticator, err = auth.Authenticate(s.conf.Auth0); err != nil {
		return err
	}

	// Instantiate user info middleware
	var userinfo gin.HandlerFunc
	if userinfo, err = auth.UserInfo(s.conf.Auth0); err != nil {
		return err
	}

	if s.conf.Sentry.UseSentry() {
		sentryRecover = sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
		})

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
		sentryRecover,

		// Add searchable tags to the sentry context.
		tags,

		// Tracing helps us with our peformance metrics and should be as early in the
		// chain as possible. It is after recovery to ensure trace panics recover.
		tracing,

		// CORS configuration allows the front-end to make cross-origin requests.
		cors.New(cors.Config{
			AllowOrigins:     s.conf.AllowOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-TOKEN", "sentry-trace", "baggage"},
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
		v1.GET("/lookup/autocomplete", s.LookupAutocomplete)
		v1.GET("/verify", s.VerifyContact)
		v1.POST("/users/login", userinfo, s.Login)
		v1.GET("/users/roles", s.ListUserRoles)
		v1.GET("/network/activity", s.NetworkActivity)

		// Authenticated routes
		// Logged-in user queries and updates
		v1.GET("/users/organization", auth.Authorize(auth.ReadOrganizations), s.UserOrganization)
		v1.PATCH("/users", userinfo, s.UpdateUser)

		// Organizations resource to manage user-ui teams
		organizations := v1.Group("/organizations")
		{
			organizations.GET("", auth.Authorize(auth.ReadOrganizations), userinfo, s.ListOrganizations)
			organizations.POST("", auth.DoubleCookie(), auth.Authorize(auth.CreateOrganizations), userinfo, s.CreateOrganization)
			organizations.DELETE("/:orgID", auth.DoubleCookie(), auth.Authorize(auth.DeleteOrganizations), userinfo, s.DeleteOrganization)
			organizations.PATCH("/:orgID", auth.DoubleCookie(), auth.Authorize(auth.UpdateOrganizations), userinfo, s.PatchOrganization)
		}

		// Collaborators resource to manage user-ui teams inside of organizations
		collaborators := v1.Group("/collaborators")
		{
			collaborators.GET("", auth.Authorize(auth.ReadCollaborators), s.ListCollaborators)
			collaborators.POST("", auth.DoubleCookie(), auth.Authorize(auth.UpdateCollaborators), userinfo, s.AddCollaborator)
			collaborators.POST("/:collabID", auth.DoubleCookie(), auth.Authorize(auth.UpdateCollaborators), s.UpdateCollaboratorRoles)
			collaborators.DELETE("/:collabID", auth.DoubleCookie(), auth.Authorize(auth.UpdateCollaborators), s.DeleteCollaborator)
		}

		// The register endpoint sends the VASP registration form to the GDS server to
		// register the VASP as a GDS TestNet or MainNet member.
		register := v1.Group("/register")
		{
			register.GET("", auth.Authorize(auth.ReadVASP), s.LoadRegisterForm)
			register.PUT("", auth.DoubleCookie(), auth.Authorize(auth.UpdateVASP), s.SaveRegisterForm)
			register.DELETE("", auth.DoubleCookie(), auth.Authorize(auth.UpdateVASP), s.ResetRegisterForm)
			register.POST("/:network", auth.DoubleCookie(), auth.Authorize(auth.UpdateVASP), userinfo, s.SubmitRegistration)
		}

		// Certificates is a resource to allow VASP members to perform certificate
		// self-service and manage the certificates issued to them.
		certificates := v1.Group("/certificates")
		{
			certificates.GET("", auth.Authorize(auth.ReadVASP), s.Certificates)
		}

		// The members resource describes verified VASPs and is only available to other
		// verified VASPs (e.g. they cannot use this endpoint during registration).
		members := v1.Group("/members", s.CheckVerification)
		{
			members.GET("", auth.Authorize(auth.ReadVASP), s.MemberList)
			members.GET("/:vaspID", auth.Authorize(auth.ReadVASP), s.MemberDetail)
		}

		// Announcements allows TRISA admins to post announcements to logged in users.
		announcements := v1.Group("/announcements")
		{
			announcements.GET("", auth.Authorize(auth.ReadVASP), s.Announcements)
			announcements.POST("", auth.DoubleCookie(), auth.Authorize("create:announcements"), s.MakeAnnouncement)
		}

		// The following are one-off endpoints that provide information to front-end
		// components at different stages in the VASP registration process.
		v1.GET("/registration", auth.Authorize(auth.ReadVASP), s.RegistrationStatus)
		v1.GET("/overview", auth.Authorize(auth.ReadVASP), s.Overview)
		v1.GET("/attention", auth.Authorize(auth.ReadVASP), s.Attention)
	}

	// NotFound and NotAllowed routes
	s.router.NoRoute(api.NotFound)
	s.router.NoMethod(api.NotAllowed)

	if s.conf.ServeDocs {
		docs.SwaggerInfo.BasePath = "/v1"
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return nil
}

//===========================================================================
// Accessors - used primarily for testing
//===========================================================================

// SetTestNetDB allows tests to set the testnet database client to a mock client
func (s *Server) SetTestNetDB(testnet store.Store) {
	s.testnetDB = testnet
}

// SetMainNetDB allows tests to set the mainnet database client to a mock client
func (s *Server) SetMainNetDB(mainnet store.Store) {
	s.mainnetDB = mainnet
}

// SetGDSClients allows tests to set a bufconn client to a mock GDS server.
func (s *Server) SetGDSClients(testnet, mainnet *GDSClient) {
	s.testnetGDS = testnet
	s.mainnetGDS = mainnet
}

// SetDB allows tests to set a bufconn client to a mock trtl server.
func (s *Server) SetDB(db store.Store) {
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

// GetURL returns the URL that the server can be reached if it has been started. This
// accessor is primarily used to create a test client.
func (s *Server) GetURL() string {
	s.RLock()
	defer s.RUnlock()
	return s.url
}
