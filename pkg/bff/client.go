package bff

import (
	"context"
	"crypto/tls"

	"github.com/hashicorp/go-multierror"
	"github.com/trisacrypto/directory/pkg/bff/config"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// GlobalDirectoryClient is a unified interface to access multiple GDS services across
// multiple connections with different client interfaces.
type GlobalDirectoryClient interface {
	gds.TRISADirectoryClient
	members.TRISAMembersClient
}

// GDSClient is a unified client which contains sub-clients for interacting with the
// various GDS services. This helps reduce common client code when making parallel
// requests to both testnet and mainnet.
type GDSClient struct {
	directoryClient
	membersClient
}

type directoryClient struct {
	client gds.TRISADirectoryClient
	conn   *grpc.ClientConn
}

type membersClient struct {
	client members.TRISAMembersClient
	conn   *grpc.ClientConn
}

// Close the connection to both the TRISA directory service and the Members service.
func (c *GDSClient) Close() (err error) {
	if cerr := c.directoryClient.close(); cerr != nil {
		err = multierror.Append(err, cerr)
	}

	if cerr := c.membersClient.close(); cerr != nil {
		err = multierror.Append(err, cerr)
	}

	return err
}

// ConnectGDS creates a gRPC client to the TRISA Directory Service specified in the
// configuration using the provided dial options.
func (c *GDSClient) ConnectGDS(conf config.DirectoryConfig, opts ...grpc.DialOption) error {
	return c.directoryClient.connect(conf, opts...)
}

func (c *directoryClient) connect(conf config.DirectoryConfig, opts ...grpc.DialOption) (err error) {
	if len(opts) == 0 {
		opts = make([]grpc.DialOption, 0, 1)

		if conf.Insecure {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			config := &tls.Config{}
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
		}
	}

	// Connect the directory client (non-blocking)
	if c.conn, err = grpc.NewClient(conf.Endpoint, opts...); err != nil {
		return err
	}
	c.client = gds.NewTRISADirectoryClient(c.conn)

	return nil
}

func (c *directoryClient) close() error {
	defer func() {
		c.client = nil
		c.conn = nil
	}()
	return c.conn.Close()
}

// ConnectMembers creates a gRPC client to the TRISA Members Service specified in the
// configuration using the provided dial options.
func (c *GDSClient) ConnectMembers(conf config.MembersConfig, opts ...grpc.DialOption) error {
	return c.membersClient.connect(conf, opts...)
}

func (c *membersClient) connect(conf config.MembersConfig, opts ...grpc.DialOption) (err error) {
	if len(opts) == 0 {
		opts = make([]grpc.DialOption, 0, 1)

		if conf.MTLS.Insecure {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			var creds grpc.DialOption
			if creds, err = conf.MTLS.DialOption(conf.Endpoint); err != nil {
				return err
			}
			opts = append(opts, creds)
		}
	}

	// Connect the members client (non-blocking)
	if c.conn, err = grpc.NewClient(conf.Endpoint, opts...); err != nil {
		return err
	}
	c.client = members.NewTRISAMembersClient(c.conn)

	return nil
}

func (c *membersClient) close() error {
	defer func() {
		c.client = nil
		c.conn = nil
	}()
	return c.conn.Close()
}

// Compile time check that GDSClient implements the GlobalDirectoryClient interface.
var _ GlobalDirectoryClient = &GDSClient{}

// GDS methods
func (c *GDSClient) Lookup(ctx context.Context, in *gds.LookupRequest, opts ...grpc.CallOption) (*gds.LookupReply, error) {
	return c.directoryClient.client.Lookup(ctx, in, opts...)
}

func (c *GDSClient) Search(ctx context.Context, in *gds.SearchRequest, opts ...grpc.CallOption) (*gds.SearchReply, error) {
	return c.directoryClient.client.Search(ctx, in, opts...)
}

func (c *GDSClient) Register(ctx context.Context, in *gds.RegisterRequest, opts ...grpc.CallOption) (*gds.RegisterReply, error) {
	return c.directoryClient.client.Register(ctx, in, opts...)
}

func (c *GDSClient) VerifyContact(ctx context.Context, in *gds.VerifyContactRequest, opts ...grpc.CallOption) (*gds.VerifyContactReply, error) {
	return c.directoryClient.client.VerifyContact(ctx, in, opts...)
}

func (c *GDSClient) Verification(ctx context.Context, in *gds.VerificationRequest, opts ...grpc.CallOption) (*gds.VerificationReply, error) {
	return c.directoryClient.client.Verification(ctx, in, opts...)
}

func (c *GDSClient) Status(ctx context.Context, in *gds.HealthCheck, opts ...grpc.CallOption) (*gds.ServiceState, error) {
	return c.directoryClient.client.Status(ctx, in, opts...)
}

// Members methods
func (c *GDSClient) List(ctx context.Context, in *members.ListRequest, opts ...grpc.CallOption) (*members.ListReply, error) {
	return c.membersClient.client.List(ctx, in, opts...)
}

func (c *GDSClient) Summary(ctx context.Context, in *members.SummaryRequest, opts ...grpc.CallOption) (*members.SummaryReply, error) {
	return c.membersClient.client.Summary(ctx, in, opts...)
}

func (c *GDSClient) Details(ctx context.Context, in *members.DetailsRequest, opts ...grpc.CallOption) (*members.MemberDetails, error) {
	return c.membersClient.client.Details(ctx, in, opts...)
}
