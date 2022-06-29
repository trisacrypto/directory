package bff

import (
	"context"

	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GlobalDirectoryClient is a unified interface which can be implemented to create a
// unified client to multiple GDS services.
type GlobalDirectoryClient interface {
	admin.DirectoryAdministrationClient
	gds.TRISADirectoryClient
	members.TRISAMembersClient
}

// GDSClient is a unified client which contains sub-clients for interacting with the
// various GDS services. This helps reduce common client code when making parallel
// request to both testnet and mainnet.
type GDSClient struct {
	admin       admin.DirectoryAdministrationClient
	creds       *Credentials
	gds         gds.TRISADirectoryClient
	members     members.TRISAMembersClient
	gdsConn     *grpc.ClientConn
	membersConn *grpc.ClientConn
}

// ConnectAdmin creates a DirectoryAdministrationClient to the GDS Admin service.
func (c *GDSClient) ConnectAdmin(conf config.AdminConfig) (err error) {
	if c.creds, err = NewCredentials(conf); err != nil {
		return err
	}

	if c.admin, err = admin.New(conf.Endpoint, c.creds); err != nil {
		return err
	}

	return nil
}

// ConnectGDS creates a gRPC client to the TRISA Directory Service specified in the
// configuration using the provided dial options.
func (c *GDSClient) ConnectGDS(conf config.DirectoryConfig, opts ...grpc.DialOption) (err error) {
	if len(opts) == 0 {
		opts = make([]grpc.DialOption, 0)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	// Connect the directory client (non-blocking)
	if c.gdsConn, err = grpc.DialContext(ctx, conf.Endpoint, opts...); err != nil {
		return err
	}
	c.gds = gds.NewTRISADirectoryClient(c.gdsConn)

	return nil
}

// ConnectMembers creates a gRPC client to the TRISA Members Service specified in the
// configuration using the provided dial options.
// TODO: Connect using mTLS.
func (c *GDSClient) ConnectMembers(conf config.MembersConfig, opts ...grpc.DialOption) (err error) {
	if len(opts) == 0 {
		opts = make([]grpc.DialOption, 0)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	// Connect the members client (non-blocking)
	if c.membersConn, err = grpc.DialContext(ctx, conf.Endpoint, opts...); err != nil {
		return err
	}
	c.members = members.NewTRISAMembersClient(c.membersConn)

	return nil
}

// Compile time check that GDSClient implements the GlobalDirectoryClient interface.
var _ GlobalDirectoryClient = &GDSClient{}

// Admin methods
func (c *GDSClient) Login(ctx context.Context) (err error) {
	return c.admin.Login(ctx)
}

func (c *GDSClient) Refresh(ctx context.Context) (err error) {
	return c.admin.Refresh(ctx)
}

func (c *GDSClient) Logout(ctx context.Context) (err error) {
	return c.admin.Logout(ctx)
}

func (c *GDSClient) AdminStatus(ctx context.Context) (out *admin.StatusReply, err error) {
	return c.admin.AdminStatus(ctx)
}

func (c *GDSClient) Authenticate(ctx context.Context, in *admin.AuthRequest) (out *admin.AuthReply, err error) {
	return c.admin.Authenticate(ctx, in)
}

func (c *GDSClient) Reauthenticate(ctx context.Context, in *admin.AuthRequest) (out *admin.AuthReply, err error) {
	return c.admin.Reauthenticate(ctx, in)
}

func (c *GDSClient) AdminSummary(ctx context.Context) (out *admin.SummaryReply, err error) {
	return c.admin.AdminSummary(ctx)
}

func (c *GDSClient) Autocomplete(ctx context.Context) (out *admin.AutocompleteReply, err error) {
	return c.admin.Autocomplete(ctx)
}

func (c *GDSClient) ReviewTimeline(ctx context.Context, params *admin.ReviewTimelineParams) (out *admin.ReviewTimelineReply, err error) {
	return c.admin.ReviewTimeline(ctx, params)
}

func (c *GDSClient) ListVASPs(ctx context.Context, params *admin.ListVASPsParams) (out *admin.ListVASPsReply, err error) {
	return c.admin.ListVASPs(ctx, params)
}

func (c *GDSClient) RetrieveVASP(ctx context.Context, id string) (out *admin.RetrieveVASPReply, err error) {
	return c.admin.RetrieveVASP(ctx, id)
}

func (c *GDSClient) UpdateVASP(ctx context.Context, in *admin.UpdateVASPRequest) (out *admin.UpdateVASPReply, err error) {
	return c.admin.UpdateVASP(ctx, in)
}

func (c *GDSClient) DeleteVASP(ctx context.Context, id string) (out *admin.Reply, err error) {
	return c.admin.DeleteVASP(ctx, id)
}

func (c *GDSClient) ReplaceContact(ctx context.Context, in *admin.ReplaceContactRequest) (out *admin.Reply, err error) {
	return c.admin.ReplaceContact(ctx, in)
}

func (c *GDSClient) DeleteContact(ctx context.Context, vaspID string, kind string) (out *admin.Reply, err error) {
	return c.admin.DeleteContact(ctx, vaspID, kind)
}

func (c *GDSClient) CreateReviewNote(ctx context.Context, in *admin.ModifyReviewNoteRequest) (out *admin.ReviewNote, err error) {
	return c.admin.CreateReviewNote(ctx, in)
}

func (c *GDSClient) ListReviewNotes(ctx context.Context, id string) (out *admin.ListReviewNotesReply, err error) {
	return c.admin.ListReviewNotes(ctx, id)
}

func (c *GDSClient) UpdateReviewNote(ctx context.Context, in *admin.ModifyReviewNoteRequest) (out *admin.ReviewNote, err error) {
	return c.admin.UpdateReviewNote(ctx, in)
}

func (c *GDSClient) DeleteReviewNote(ctx context.Context, vaspID string, noteID string) (out *admin.Reply, err error) {
	return c.admin.DeleteReviewNote(ctx, vaspID, noteID)
}

func (c *GDSClient) ReviewToken(ctx context.Context, vaspID string) (out *admin.ReviewTokenReply, err error) {
	return c.admin.ReviewToken(ctx, vaspID)
}

func (c *GDSClient) Review(ctx context.Context, in *admin.ReviewRequest) (out *admin.ReviewReply, err error) {
	return c.admin.Review(ctx, in)
}

func (c *GDSClient) Resend(ctx context.Context, in *admin.ResendRequest) (out *admin.ResendReply, err error) {
	return c.admin.Resend(ctx, in)
}

// GDS methods
func (c *GDSClient) Lookup(ctx context.Context, in *gds.LookupRequest, opts ...grpc.CallOption) (*gds.LookupReply, error) {
	return c.gds.Lookup(ctx, in, opts...)
}

func (c *GDSClient) Search(ctx context.Context, in *gds.SearchRequest, opts ...grpc.CallOption) (*gds.SearchReply, error) {
	return c.gds.Search(ctx, in, opts...)
}

func (c *GDSClient) Register(ctx context.Context, in *gds.RegisterRequest, opts ...grpc.CallOption) (*gds.RegisterReply, error) {
	return c.gds.Register(ctx, in, opts...)
}

func (c *GDSClient) VerifyContact(ctx context.Context, in *gds.VerifyContactRequest, opts ...grpc.CallOption) (*gds.VerifyContactReply, error) {
	return c.gds.VerifyContact(ctx, in, opts...)
}

func (c *GDSClient) Verification(ctx context.Context, in *gds.VerificationRequest, opts ...grpc.CallOption) (*gds.VerificationReply, error) {
	return c.gds.Verification(ctx, in, opts...)
}

func (c *GDSClient) Status(ctx context.Context, in *gds.HealthCheck, opts ...grpc.CallOption) (*gds.ServiceState, error) {
	return c.gds.Status(ctx, in, opts...)
}

// Members methods
func (c *GDSClient) List(ctx context.Context, in *members.ListRequest, opts ...grpc.CallOption) (*members.ListReply, error) {
	return c.members.List(ctx, in, opts...)
}

func (c *GDSClient) Summary(ctx context.Context, in *members.SummaryRequest, opts ...grpc.CallOption) (*members.SummaryReply, error) {
	return c.members.Summary(ctx, in, opts...)
}
