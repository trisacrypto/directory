package gds

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	pb "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewMembers creates a new Member server derived from a parent Service.
func NewMembers(svc *Service) (members *Members, err error) {
	members = &Members{
		svc:  svc,
		conf: &svc.conf.Members,
		db:   svc.db,
	}

	// Initialize the gRPC server
	members.srv = grpc.NewServer(grpc.UnaryInterceptor(svc.serverInterceptor))
	pb.RegisterTRISAMembersServer(members.srv, members)
	return members, nil
}

// Members implements the TRISAMembers service as defined by the experimental v1alpha1
// protocol buffers in the GDS repository. This service is intended to be an mTLS
// authenticated service (which is why it is separate from the GDS service) that is used
// directly by TRISA members to facilitate p2p exchanges and GDS lookups.
//
// NOTE: this is a prototype service, this service may eventually be moved into the GDS
// specification in trisacrypto/trisa.
type Members struct {
	pb.UnimplementedTRISAMembersServer
	svc  *Service              // The parent Service GDS uses to interact with other components
	srv  *grpc.Server          // The gRPC server that listens on its own independent port
	conf *config.MembersConfig // The GDS service specific configuration (helper alias to s.svc.conf.Members)
	db   store.Store           // Database connection for loading objects (helper alias to s.svc.db)
}

// Serve gRPC requests on the specified address.
func (s *Members) Serve() (err error) {
	if !s.conf.Enabled {
		log.Warn().Msg("trisa members service is not enabled")
		return nil
	}

	// This service must not run if we're in maintenance mode
	if s.svc.conf.Maintenance {
		return errors.New("cannot serve Members server in maintenace mode")
	}

	// Listen for TCP requests on the specified address and port
	var sock net.Listener
	if sock, err = net.Listen("tcp", s.conf.BindAddr); err != nil {
		return fmt.Errorf("could not listen on %q", s.conf.BindAddr)
	}

	// Run the server
	go s.Run(sock)
	log.Info().Str("listen", s.conf.BindAddr).Str("version", pkg.Version()).Msg("trisa members server started")

	// Now that the go routine is started return nil, meaning the service has started
	// successfully with no problems.
	return nil
}

// Run the gRPC server. This method is extracted from the Serve function so that it can
// be run in its own go routine and to allow tests to Run a bufconn server without
// starting a live server with all of the various go routines and channels running.
func (s *Members) Run(sock net.Listener) {
	defer sock.Close()
	if err := s.srv.Serve(sock); err != nil {
		s.svc.echan <- err
	}
}

// Shutdown the TRISA Members Service gracefully
func (s *Members) Shutdown() (err error) {
	log.Debug().Msg("gracefully shutting down TRISA Members server")
	s.srv.GracefulStop()
	log.Debug().Msg("successful shutdown of TRISA Members server")
	return nil
}

// List all verified VASP members in the Directory Service. This RPC returns an
// abbreviated listing of VASP details intended to facilitate p2p exchanges or more
// detailed lookups against the Directory Service. The response is paginated. If there
// are more results than the specified page size, then the reply will include a next
// page token. That token can be used to fetch the next page so long as the parameters
// of the original request are not modified (e.g. any filters or pagination parameters).
// See https://cloud.google.com/apis/design/design_patterns#list_pagination for more.
func (s *Members) List(ctx context.Context, in *pb.ListRequest) (out *pb.ListReply, err error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented quite yet")
}
