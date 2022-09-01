/*
Package server implements a lightweight Sectigo mock server that can be used in staging
to issue mock certificates and perform integration tests. This server implements a
subset of the Sectigo IoT API that targets GDS-specific usage. All state is held
in-memory and is periodically flushed so this service should not be relied on for
anything other than staging and systems integration tests.
*/
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
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

type Server struct {
	sync.RWMutex
	conf    Config
	srv     *http.Server
	router  *gin.Engine
	started time.Time
	healthy bool
	url     string
	echan   chan error
}

// New is the primary entry point to creating a Sectigo Integration API Server.
func New(conf Config) (s *Server, err error) {
	// Load config from environment if an empty config is passed in.
	if conf.IsZero() {
		if conf, err = NewConfig(); err != nil {
			return nil, err
		}
	}

	// Manage logging for debugging in k8s.
	zerolog.SetGlobalLevel(conf.GetLogLevel())
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Prepare to serve the Sectigo API
	s = &Server{
		conf:  conf,
		echan: make(chan error, 1),
	}

	// Create the router for handling HTTP requests
	gin.SetMode(conf.Mode)
	s.router = gin.New()
	if err = s.setupRoutes(); err != nil {
		return nil, err
	}

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

func (s *Server) Serve() (err error) {
	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.echan <- s.Shutdown()
	}()

	// Create a socket to listen on so that we can infer the final URL.
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		return fmt.Errorf("format strincould not listen on bind addr %s: %s", s.conf.BindAddr, err)
	}

	// Update the server's internal state
	s.Lock()
	s.healthy = true
	s.url = "http://" + sock.Addr().String()
	s.started = time.Now()
	s.Unlock()

	// Listen for HTTP requests on the specified address and port
	go func() {
		if err = s.srv.Serve(sock); err != nil && err != http.ErrServerClosed {
			s.echan <- err
		}
	}()

	log.Info().
		Str("listen", s.url).
		Str("version", pkg.Version()).
		Msg("sectigo integration api server started")

	// Fatal errors or stop signals should be sent on the error chan.
	if err = <-s.echan; err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() error {
	log.Info().Msg("gracefully shutting down sectigo integration api server")

	s.SetHealth(false)
	s.srv.SetKeepAlivesEnabled(false)

	// Require shutdown in 30 seconds without blocking
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Debug().Msg("successfully shutdown sectigo integration api server")
	return nil
}

func (s *Server) SetHealth(health bool) {
	s.Lock()
	s.healthy = health
	s.Unlock()
	log.Debug().Bool("healthy", health).Msg("server health set")
}

func (s *Server) setupRoutes() (err error) {
	// Application Middleware
	// NOTE: ordering is very important to how middlware is handled
	middlewares := []gin.HandlerFunc{
		// Logging should be outside so we can record the complete latency of requests.
		// NOTE: logging panics will not recover.
		logger.GinLogger("sias"),

		// Panic recovery middleware
		gin.Recovery(),

		// Maintenance mode handling - does not require authentication
		s.Available(),
	}

	for _, middleware := range middlewares {
		if middleware != nil {
			s.router.Use(middleware)
		}
	}

	// Heartbeat route (no authentication required)
	s.router.GET("/status", s.Status)

	s.router.NoRoute(s.NotFound)
	s.router.NoMethod(s.NotAllowed)
	return nil
}
