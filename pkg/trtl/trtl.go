package trtl

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/rotationalio/honu"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"google.golang.org/grpc"
)

// A Trtl server implements the following services
// 1. A database service for interacting with a database
// 2. A peers management service for interacting with remote peers
// 3. A replication service which implements auto-adapting anti-entropy replication.
type Server struct {
	srv     *grpc.Server    // The gRPC server that listens on its own independent port
	conf    config.Config   // Configuration for the trtl server
	db      *honu.DB        // Database connection for managing objects
	honu    *HonuService    // Service for interacting with a Honu database
	peers   *PeerService    // Service for managing remote peers
	replica *ReplicaService // Service that handles anti-entropy replication
	echan   chan error      // Channel for receiving errors from the gRPC server
}

// New creates a new trtl server given a configuration.
func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// TODO: manage logging

	// Create the server and prepare to serve
	s = &Server{
		conf:  conf,
		echan: make(chan error, 1),
		srv:   grpc.NewServer(grpc.UnaryInterceptor(s.interceptor)),
	}

	// TODO: check for maintenance mode

	// Everything that follows this comment assumes we're not in maintenance mode
	// Open a connection to the Honu wrapped database
	fmt.Println("here")
	if s.db, err = honu.Open(conf.Database.URL, conf.GetHonuConfig()); err != nil {
		return nil, fmt.Errorf("honu error: %v", err)
	}

	// Initialize the Honu service
	if s.honu, err = NewHonuService(s); err != nil {
		return nil, err
	}
	pb.RegisterTrtlServer(s.srv, s.honu)

	// Initialize the Peer Management service
	if s.peers, err = NewPeerService(s); err != nil {
		return nil, err
	}
	peers.RegisterPeerManagementServer(s.srv, s.peers)

	// TODO: initialize the Replica service

	return s, nil
}

// Serve gRPC requests on the specified bind address.
func (t *Server) Serve() (err error) {
	// TODO: change from enabled to maintenance mode
	if !t.conf.Enabled {
		log.Warn().Msg("trtl service is not enabled")
		return nil
	}

	// Catch OS signals for graceful shutdowns
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		t.echan <- t.Shutdown()
	}()

	// Run the Gossip background routine
	go t.replica.AntiEntropy()

	// Listen for TCP requests
	var sock net.Listener
	if sock, err = net.Listen("tcp", t.conf.BindAddr); err != nil {
		log.Error().Err(err).Str("bindaddr", t.conf.BindAddr).Msg("could not listen on given bindaddr")
		return err
	}

	// Run the gRPC server
	go func() {
		defer sock.Close()
		log.Info().Str("listen", t.conf.BindAddr).Msg("trtl server started")
		if err := t.srv.Serve(sock); err != nil {
			t.echan <- err
		}
	}()

	// The server go routine is started so return nil error (any server errors will be
	// sent on the error channel).
	if err = <-t.echan; err != nil {
		return err
	}
	return nil
}

// Shutdown the trtl server gracefully.
func (t *Server) Shutdown() (err error) {
	// TODO: collect multi errors to return after shutdown
	log.Info().Msg("gracefully shutting down trtl server")
	t.srv.GracefulStop()

	// TODO: Stop the anti-entropy routine.

	if err = t.db.Close(); err != nil {
		log.Error().Err(err).Msg("could not close database")
	}

	log.Debug().Msg("successful shutdown of trtl server")
	return nil
}
