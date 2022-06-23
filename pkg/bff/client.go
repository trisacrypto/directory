package bff

import (
	"context"

	"github.com/trisacrypto/directory/pkg/bff/config"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GlobalDirectoryClient is a unified interface which can be implemented to create a
// unified client to multiple GDS services.
type GlobalDirectoryClient interface {
	gds.TRISADirectoryClient
	members.TRISAMembersClient
}

// GDSClient is a unified client which contains sub-clients for interacting with the
// directory service and members service. This helps reduce common client code when
// making parallel requests to both the TestNet and MainNet.
type GDSClient struct {
	gds         gds.TRISADirectoryClient
	members     members.TRISAMembersClient
	gdsConn     *grpc.ClientConn
	membersConn *grpc.ClientConn
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
func (c *GDSClient) ConnectMembers(conf config.DirectoryConfig, opts ...grpc.DialOption) (err error) {
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

func (c *GDSClient) List(ctx context.Context, in *members.ListRequest, opts ...grpc.CallOption) (*members.ListReply, error) {
	return c.members.List(ctx, in, opts...)
}

func (c *GDSClient) Summary(ctx context.Context, in *members.SummaryRequest, opts ...grpc.CallOption) (*members.SummaryReply, error) {
	return c.members.Summary(ctx, in, opts...)
}
