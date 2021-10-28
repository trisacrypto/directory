package trtl

import (
	"net"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/store"
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
	conf    *config.Config  // Configuration for the trtl server
	db      store.Store     // Database connection for managing objects (alias to s.svc.db)
	honu    *HonuService    // Service for interacting with a Honu database
	peer    *PeerService    // Service for managing remote peers
	replica *ReplicaService // Service that handles anti-entropy replication
	echan   chan error      // Channel for receiving errors from the gRPC server
}

// New creates a new trtl server given a configuration.
func New(db store.Store, conf config.Config) (s *Server, err error) {
	s = &Server{
		conf:    &conf,
		db:      db,
		honu:    NewHonuService(),
		peer:    NewPeerService(db),
		replica: NewReplicaService(db, conf),
	}

	// TODO: Check if the database Store is an Honu DB, if not then the Replica cannot Gossip.

	// Initialize the gRPC server
	s.srv = grpc.NewServer(grpc.UnaryInterceptor(s.interceptor))
	pb.RegisterTrtlServer(s.srv, s.honu)
	peers.RegisterPeerManagementServer(s.srv, s.peer)
	return s, nil
}

// Serve gRPC requests on the specified bind address.
func (t *Server) Serve() (err error) {
	if !t.conf.Enabled {
		log.Warn().Msg("trtl service is not enabled")
		return nil
	}

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
	return nil
}

// Shutdown the trtl server gracefully.
func (t *Server) Shutdown() (err error) {
	log.Debug().Msg("gracefully shutting down trtl server")
	t.srv.GracefulStop()
	// TODO: Also need a way to stop the anti-entropy routine.
	log.Debug().Msg("successful shutdown of trtl server")
	return nil
}
