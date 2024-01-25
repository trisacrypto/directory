package gds

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/config"
	api "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"github.com/trisacrypto/directory/pkg/gds/pagination/v1"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/mtls"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultPageSize = 100
)

// NewMembers creates a new Member server derived from a parent Service.
func NewMembers(svc *Service) (members *Members, err error) {
	members = &Members{
		svc:  svc,
		conf: &svc.conf.Members,
		db:   svc.db,
	}

	// Attempt to load and parse the TRISA certificates for server-size mTLS
	var sz *trust.Serializer
	if sz, err = trust.NewSerializer(false); err != nil {
		return nil, err
	}

	// Initialize mTLS for the server if configured
	opts := make([]grpc.ServerOption, 0, 2)
	if !members.conf.Insecure {
		// Read the certificates issued by the directory service to run the directory service
		if members.mtlsCerts, err = sz.ReadFile(members.conf.Certs); err != nil {
			return nil, fmt.Errorf("could not load members certs and private key: %s", err)
		}

		// Read the trust pool that was issued by the directory service (public CA keys)
		if members.trustPool, err = sz.ReadPoolFile(members.conf.CertPool); err != nil {
			return nil, fmt.Errorf("could not load members public cert pool: %s", err)
		}

		// Create TLS credentials for the server
		var creds grpc.ServerOption
		if creds, err = mtls.ServerCreds(members.mtlsCerts, members.trustPool); err != nil {
			return nil, fmt.Errorf("could not create mTLS creds: %s", err)
		}
		opts = append(opts, creds)
	} else {
		sentry.Warn(nil).Msg("creating insecure trisa members server")
	}

	// Add the unary interceptor to the gRPC server
	opts = append(opts, grpc.ChainUnaryInterceptor(svc.UnaryInterceptors()...))
	opts = append(opts, grpc.ChainStreamInterceptor(svc.StreamInterceptors()...))

	// Initialize the gRPC server
	members.srv = grpc.NewServer(opts...)
	api.RegisterTRISAMembersServer(members.srv, members)
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
	api.UnimplementedTRISAMembersServer
	svc       *Service              // The parent Service GDS uses to interact with other components
	srv       *grpc.Server          // The gRPC server that listens on its own independent port
	conf      *config.MembersConfig // The GDS service specific configuration (helper alias to s.svc.conf.Members)
	db        store.Store           // Database connection for loading objects (helper alias to s.svc.db)
	mtlsCerts *trust.Provider       // Server certificate and private keys for server-auth
	trustPool trust.ProviderPool    // Cert pool for client-side authentication
}

// Serve gRPC requests on the specified address.
func (s *Members) Serve() (err error) {
	if !s.conf.Enabled {
		log.Warn().Msg("trisa members service is not enabled")
		return nil
	}

	// This service must not run if we're in maintenance mode
	if s.svc.conf.Maintenance {
		return errors.New("cannot serve Members server in maintenance mode")
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
func (s *Members) List(ctx context.Context, in *api.ListRequest) (out *api.ListReply, err error) {
	// Use default page size if one isn't specified
	if in.PageSize == 0 {
		in.PageSize = defaultPageSize
	}

	// If a page cursor is provided, load it - otherwise create a cursor for iteration
	cursor := &pagination.PageCursor{}
	if in.PageToken != "" {
		if err = cursor.Load(in.PageToken); err != nil {
			sentry.Warn(ctx).Err(err).Msg("invalid page token on list request")
			return nil, status.Error(codes.InvalidArgument, "invalid page token")
		}

		// Validate the request has not changed
		if cursor.PageSize != in.PageSize {
			log.Debug().Int32("cursor", cursor.PageSize).Int32("opts", in.PageSize).Msg("invalid members list request: mismatched page size")
			return nil, status.Error(codes.InvalidArgument, "page size cannot change between requests")
		}

	} else {
		// Update the cursor with the input request
		cursor.PageSize = in.PageSize
	}

	// Create response
	out = &api.ListReply{
		Vasps: make([]*api.VASPMember, 0, cursor.PageSize),
	}

	// Create the VASPs iterator to begin collecting validated VASPs data
	iter := s.db.ListVASPs(ctx)
	defer iter.Release()

	// If necessary, seek to the next key specified by the cursor.
	if cursor.NextVasp != "" {
		// If iter.SeekId() returns false (e.g. seek did not find the specified key) then
		// iter.Next() should also return false, so it isn't necessary to check the return.
		// NOTE: next key must be deleted after it's used for seeking so that the last
		// page doesn't retain the old key and loop forever.
		iter.SeekId(cursor.NextVasp)
		cursor.NextVasp = ""

		// Because we're going to be calling Next, we need to back up one key to ensure
		// that we start on the right key in the for loop.
		iter.Prev()
	}

	// Iterate over VASPs, collecting the one we're looking for.
	for iter.Next() {
		// Check if we're done iterating - if so and there is next data, there is
		// another page, so create the page token to return it.
		if len(out.Vasps) == int(cursor.PageSize) {
			cursor.NextVasp = iter.Id()
			break
		}

		// Collect the VASP from the iterator
		var vasp *pb.VASP
		if vasp, err = iter.VASP(); err != nil {
			sentry.Error(ctx).Err(err).Msg("could not parse VASP from database")
			continue
		}

		// Skip any VASPs that are not verified yet
		if vasp.VerificationStatus != pb.VerificationState_VERIFIED {
			continue
		}

		out.Vasps = append(out.Vasps, GetVASPMember(vasp))
	}

	if err = iter.Error(); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not iterate over VASPs")
		return nil, status.Error(codes.Internal, "could not iterate over directory service")
	}

	// Check if there is a next page cursor
	if cursor.NextVasp != "" {
		if out.NextPageToken, err = cursor.Dump(); err != nil {
			sentry.Error(ctx).Err(err).Msg("could not serialize next page token on vasp member list")
			return nil, status.Error(codes.Internal, "could not create next page token")
		}
	}

	// Request Complete
	log.Info().Int("count", len(out.Vasps)).Bool("has_next_page", out.NextPageToken != "").Msg("vasp member list complete")
	return out, nil
}

// Summary returns a summary of the VASP members in the Directory Service.
// Note: Any VASP can call this endpoint with any VASP ID, therefore we need to avoid
// returning sensitive VASP details here such as IVMS info.
func (s *Members) Summary(ctx context.Context, in *api.SummaryRequest) (out *api.SummaryReply, err error) {
	// Get the since timestamp from the request
	var since time.Time
	if in.Since != "" {
		if since, err = time.Parse(time.RFC3339, in.Since); err != nil {
			sentry.Warn(ctx).Str("since", in.Since).Err(err).Msg("invalid since timestamp on summary request")
			return nil, status.Error(codes.InvalidArgument, "since must be a valid RFC3339 timestamp")
		}

		if since.After(time.Now()) {
			sentry.Warn(ctx).Str("since", in.Since).Err(err).Msg("since timestamp must be in the past")
			return nil, status.Error(codes.InvalidArgument, "since timestamp must be in the past")
		}
	} else {
		// Default to one month ago for new members
		since = time.Now().Add(-time.Hour * 24 * 30)
	}

	// Create response
	out = &api.SummaryReply{}

	if in.MemberId != "" {
		// Fetch the requested VASP if provided
		var vasp *pb.VASP
		if vasp, err = s.db.RetrieveVASP(ctx, in.MemberId); err != nil {
			sentry.Warn(ctx).Err(err).Str("vasp_id", in.MemberId).Msg("VASP not found")
			return nil, status.Error(codes.NotFound, "requested VASP not found")
		}

		// Add the VASP member details to the response
		out.MemberInfo = GetVASPMember(vasp)
	}

	// Create the VASPs iterator
	iter := s.db.ListVASPs(ctx)
	defer iter.Release()

	// Iterate over VASPs
	for iter.Next() {
		// Collect the VASP from the iterator
		var vasp *pb.VASP
		if vasp, err = iter.VASP(); err != nil {
			sentry.Error(ctx).Err(err).Msg("could not parse VASP from database")
			continue
		}

		// Only count verified VASPs
		if vasp.VerificationStatus == pb.VerificationState_VERIFIED {
			// Parse verified timestamp
			var verified time.Time
			if verified, err = time.Parse(time.RFC3339, vasp.VerifiedOn); err != nil {
				sentry.Error(ctx).Err(err).Str("vasp_id", vasp.Id).Msg("could not parse verified timestamp")
				continue
			}

			// Count if the VASP is a new member
			if verified.After(since) {
				out.NewMembers++
			}

			// Increment the total VASPs count
			out.Vasps++
		}

		// Count issued certificates
		if vasp.IdentityCertificate != nil && !vasp.IdentityCertificate.Revoked {
			out.CertificatesIssued++
		}
	}

	if err = iter.Error(); err != nil {
		sentry.Error(ctx).Err(err).Msg("could not iterate over VASPs")
		return nil, status.Error(codes.Internal, "could not iterate over directory service")
	}

	return out, nil
}

// Details returns the details of the VASP member using the provided ID.
func (s *Members) Details(ctx context.Context, in *api.DetailsRequest) (out *api.MemberDetails, err error) {
	// Fetch the requested VASP if provided
	var vasp *pb.VASP
	if vasp, err = s.db.RetrieveVASP(ctx, in.MemberId); err != nil {
		sentry.Warn(ctx).Err(err).Str("vasp_id", in.MemberId).Msg("VASP not found")
		return nil, status.Error(codes.NotFound, "requested VASP not found")
	}

	// Construct the member response
	out = &api.MemberDetails{
		MemberSummary: GetVASPMember(vasp),
	}

	// Add the IVMS101 legal person to the response
	if vasp.Entity == nil {
		sentry.Warn(ctx).Str("vasp_id", vasp.Id).Msg("VASP is missing legal person")
	}
	out.LegalPerson = vasp.Entity

	// Add the TRIXO form data to the response
	if vasp.Trixo == nil {
		sentry.Warn(ctx).Str("vasp_id", vasp.Id).Msg("VASP is missing TRIXO form data")
	}
	out.Trixo = vasp.Trixo

	// Add the contacts data to the response
	if vasp.Contacts == nil {
		sentry.Warn(ctx).Str("vasp_id", vasp.Id).Msg("VASP is missing contacts data")
	}

	// NOTE: ensure that extra and person are null on return.
	// We only want to supply name, email, and phone to other VASPs.
	out.Contacts = &pb.Contacts{}
	if vasp.Contacts.Administrative != nil {
		out.Contacts.Administrative = &pb.Contact{
			Name:  vasp.Contacts.Administrative.Name,
			Email: vasp.Contacts.Administrative.Email,
			Phone: vasp.Contacts.Administrative.Phone,
		}
	}

	if vasp.Contacts.Technical != nil {
		out.Contacts.Technical = &pb.Contact{
			Name:  vasp.Contacts.Technical.Name,
			Email: vasp.Contacts.Technical.Email,
			Phone: vasp.Contacts.Technical.Phone,
		}
	}

	if vasp.Contacts.Billing != nil {
		out.Contacts.Billing = &pb.Contact{
			Name:  vasp.Contacts.Billing.Name,
			Email: vasp.Contacts.Billing.Email,
			Phone: vasp.Contacts.Billing.Phone,
		}
	}

	if vasp.Contacts.Legal != nil {
		out.Contacts.Legal = &pb.Contact{
			Name:  vasp.Contacts.Legal.Name,
			Email: vasp.Contacts.Legal.Email,
			Phone: vasp.Contacts.Legal.Phone,
		}
	}
	return out, nil
}

// GetVASPMember is a helper function to construct a VASPMember from a VASP record.
func GetVASPMember(vasp *pb.VASP) *api.VASPMember {
	var err error

	info := &api.VASPMember{
		Id:                  vasp.Id,
		RegisteredDirectory: vasp.RegisteredDirectory,
		CommonName:          vasp.CommonName,
		Endpoint:            vasp.TrisaEndpoint,
		Website:             vasp.Website,
		BusinessCategory:    vasp.BusinessCategory,
		VaspCategories:      vasp.VaspCategories,
		VerifiedOn:          vasp.VerifiedOn,
		Status:              vasp.VerificationStatus,
		FirstListed:         vasp.FirstListed,
		LastUpdated:         vasp.LastUpdated,
	}

	// Try to add the name information
	if info.Name, err = vasp.Name(); err != nil {
		sentry.Error(nil).Err(err).Str("vasp_id", vasp.Id).Msg("could not retrieve VASP name from record")
	}

	// Add the country information if available
	if vasp.Entity != nil {
		info.Country = vasp.Entity.CountryOfRegistration
	}

	return info
}
