package mock

import (
	"context"
	"fmt"
	"os"

	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	RegisterRPC      = "Register"
	LookupRPC        = "Lookup"
	SearchRPC        = "Search"
	VerificationRPC  = "Verification"
	VerifyContactRPC = "VerifyContact"
	StatusRPC        = "Status"
)

func NewGDS(conf config.DirectoryConfig) (g *GDS, err error) {
	g = &GDS{
		srv:   grpc.NewServer(),
		sock:  bufconn.New(""),
		Calls: make(map[string]int),
	}

	gds.RegisterTRISADirectoryServer(g.srv, g)
	go g.srv.Serve(g.sock.Listener)
	return g, nil
}

// GDS implements a mock gRPC server that listens on a buffconc and registers the
// TRISA directory service. The RPC methods are able to be set from individual tests
// so that the user can specify the return of the RPC in order to test specific
// functionality. The mock allows us to test dual GDS (TestNet and MainNet) handling
// from the BFF without having to set up a whole microservices architecture with data
// storage and managing fixtures.
//
// To set the response of the mock for a particular test, update the GDS OnRPC method.
// e.g. to mock the Register RPC, set the OnRegister method. The number of calls to the
// RPC will be recorded to verifiy that the service is being called correctly. Use the
// Reset method to remove all RPC handlers and set calls back to 0.
//
// NOTE: this mock is not safe for concurrent use.
// NOTE: if the OnRPC function is not set, the test will panic
type GDS struct {
	gds.UnimplementedTRISADirectoryServer
	sock            *bufconn.GRPCListener
	srv             *grpc.Server
	client          gds.TRISADirectoryClient
	Calls           map[string]int
	OnRegister      func(context.Context, *gds.RegisterRequest) (*gds.RegisterReply, error)
	OnLookup        func(context.Context, *gds.LookupRequest) (*gds.LookupReply, error)
	OnSearch        func(context.Context, *gds.SearchRequest) (*gds.SearchReply, error)
	OnVerification  func(context.Context, *gds.VerificationRequest) (*gds.VerificationReply, error)
	OnVerifyContact func(context.Context, *gds.VerifyContactRequest) (*gds.VerifyContactReply, error)
	OnStatus        func(context.Context, *gds.HealthCheck) (*gds.ServiceState, error)
}

func (g *GDS) Client() (client gds.TRISADirectoryClient, err error) {
	if g.client == nil {
		if err = g.sock.Connect(context.Background()); err != nil {
			return nil, err
		}
		g.client = gds.NewTRISADirectoryClient(g.sock.Conn)
	}
	return g.client, nil
}

func (g *GDS) DialOpts() (opts []grpc.DialOption) {
	return []grpc.DialOption{
		grpc.WithContextDialer(g.sock.Dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func (g *GDS) Shutdown() {
	// Close the connection that any clients may have opened
	if g.sock.Conn != nil {
		g.sock.Close()
		g.sock.Conn = nil
	}

	// Stop the gRPC server
	g.srv.GracefulStop()

	// Release the buffcon
	if g.sock != nil {
		g.sock.Release()
		g.sock = nil
	}
}

func (g *GDS) Reset() {
	// Set the calls to 0
	for key := range g.Calls {
		g.Calls[key] = 0
	}

	// Reset all of the gRPC methods to ensure that an RPC from a previous test doesn't
	// interfere with the operation of a current test.
	g.OnRegister = nil
	g.OnLookup = nil
	g.OnSearch = nil
	g.OnVerification = nil
	g.OnVerifyContact = nil
	g.OnStatus = nil
}

// UseFixture allows you to specify a JSON fixture that is loaded from disk as the
// protocol buffer response for the specified RPC, simplifying the handler mocking.
func (g *GDS) UseFixture(rpc, path string) (err error) {
	// Read the fixture data from disk
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return fmt.Errorf("could not read fixture data: %s", err)
	}

	// Create a protobuf JSON unmarshaler
	jsonpb := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	switch rpc {
	case RegisterRPC:
		out := &gds.RegisterReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %s", out, err)
		}
		g.OnRegister = func(context.Context, *gds.RegisterRequest) (*gds.RegisterReply, error) {
			return out, nil
		}
	case LookupRPC:
		out := &gds.LookupReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %s", out, err)
		}
		g.OnLookup = func(context.Context, *gds.LookupRequest) (*gds.LookupReply, error) {
			return out, nil
		}
	case SearchRPC:
		out := &gds.SearchReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %s", out, err)
		}
		g.OnSearch = func(context.Context, *gds.SearchRequest) (*gds.SearchReply, error) {
			return out, nil
		}
	case VerificationRPC:
		out := &gds.VerificationReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %s", out, err)
		}
		g.OnVerification = func(context.Context, *gds.VerificationRequest) (*gds.VerificationReply, error) {
			return out, nil
		}
	case VerifyContactRPC:
		out := &gds.VerifyContactReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %s", out, err)
		}
		g.OnVerifyContact = func(context.Context, *gds.VerifyContactRequest) (*gds.VerifyContactReply, error) {
			return out, nil
		}
	case StatusRPC:
		out := &gds.ServiceState{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %s", out, err)
		}
		g.OnStatus = func(context.Context, *gds.HealthCheck) (*gds.ServiceState, error) {
			return out, nil
		}
	default:
		return fmt.Errorf("unknown rpc %q", rpc)
	}
	return nil
}

// UseError allows you to specify a gRPC status error to return from the specified RPC.
func (g *GDS) UseError(rpc string, code codes.Code, msg string) error {
	switch rpc {
	case RegisterRPC:
		g.OnRegister = func(context.Context, *gds.RegisterRequest) (*gds.RegisterReply, error) {
			return nil, status.Error(code, msg)
		}
	case LookupRPC:
		g.OnLookup = func(context.Context, *gds.LookupRequest) (*gds.LookupReply, error) {
			return nil, status.Error(code, msg)
		}
	case SearchRPC:
		g.OnSearch = func(context.Context, *gds.SearchRequest) (*gds.SearchReply, error) {
			return nil, status.Error(code, msg)
		}
	case VerificationRPC:
		g.OnVerification = func(context.Context, *gds.VerificationRequest) (*gds.VerificationReply, error) {
			return nil, status.Error(code, msg)
		}
	case VerifyContactRPC:
		g.OnVerifyContact = func(context.Context, *gds.VerifyContactRequest) (*gds.VerifyContactReply, error) {
			return nil, status.Error(code, msg)
		}
	case StatusRPC:
		g.OnStatus = func(context.Context, *gds.HealthCheck) (*gds.ServiceState, error) {
			return nil, status.Error(code, msg)
		}
	default:
		return fmt.Errorf("unknown rpc %q", rpc)
	}
	return nil
}

func (g *GDS) Register(ctx context.Context, in *gds.RegisterRequest) (*gds.RegisterReply, error) {
	g.Calls[RegisterRPC]++
	return g.OnRegister(ctx, in)
}

func (g *GDS) Lookup(ctx context.Context, in *gds.LookupRequest) (*gds.LookupReply, error) {
	g.Calls[LookupRPC]++
	return g.OnLookup(ctx, in)
}

func (g *GDS) Search(ctx context.Context, in *gds.SearchRequest) (*gds.SearchReply, error) {
	g.Calls[SearchRPC]++
	return g.OnSearch(ctx, in)
}

func (g *GDS) Verification(ctx context.Context, in *gds.VerificationRequest) (*gds.VerificationReply, error) {
	g.Calls[VerificationRPC]++
	return g.OnVerification(ctx, in)
}

func (g *GDS) VerifyContact(ctx context.Context, in *gds.VerifyContactRequest) (*gds.VerifyContactReply, error) {
	g.Calls[VerifyContactRPC]++
	return g.OnVerifyContact(ctx, in)
}

func (g *GDS) Status(ctx context.Context, in *gds.HealthCheck) (*gds.ServiceState, error) {
	g.Calls[StatusRPC]++
	return g.OnStatus(ctx, in)
}
