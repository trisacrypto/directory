package mock

import (
	"context"
	"fmt"
	"os"

	"github.com/trisacrypto/directory/pkg/bff/config"
	members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	ListRPC    = "List"
	SummaryRPC = "Summary"
	DetailsRPC = "Details"
)

func NewMembers(conf config.MembersConfig) (m *Members, err error) {
	m = &Members{
		srv:   grpc.NewServer(),
		sock:  bufconn.New(""),
		Calls: make(map[string]int),
	}

	members.RegisterTRISAMembersServer(m.srv, m)
	go m.srv.Serve(m.sock.Listener)
	return m, nil
}

// Memberss implements a mock gRPC server that listens on a buffcon and registers the
// TRISA members service. The RPC methods are able to be set from individual tests
// so that the user can specify the return of the RPC in order to test specific
// functionality. The mock allows us to test dual members (TestNet and MainNet) handling
// from the BFF without having to set up a whole microservices architecture with data
// storage and managing fixtures.
//
// To set the response of the mock for a particular test, update the Members OnRPC method.
// e.g. to mock the Summary RPC, set the OnSummary method. The number of calls to the
// RPC will be recorded to verifiy that the service is being called correctly. Use the
// Reset method to remove all RPC handlers and set calls back to 0.
//
// NOTE: this mock is not safe for concurrent use.
// NOTE: if the OnRPC function is not set, the test will panic
type Members struct {
	members.UnimplementedTRISAMembersServer
	sock      *bufconn.GRPCListener
	srv       *grpc.Server
	client    members.TRISAMembersClient
	Calls     map[string]int
	OnList    func(context.Context, *members.ListRequest) (*members.ListReply, error)
	OnSummary func(context.Context, *members.SummaryRequest) (*members.SummaryReply, error)
	OnDetails func(context.Context, *members.DetailsRequest) (*members.MemberDetails, error)
}

func (g *Members) Client() (client members.TRISAMembersClient, err error) {
	if g.client == nil {
		if err = g.sock.Connect(context.Background()); err != nil {
			return nil, err
		}
		g.client = members.NewTRISAMembersClient(g.sock.Conn)
	}
	return g.client, nil
}

func (g *Members) DialOpts() (opts []grpc.DialOption) {
	return []grpc.DialOption{
		grpc.WithContextDialer(g.sock.Dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func (g *Members) Shutdown() {
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

func (m *Members) Reset() {
	// Set the calls to 0
	for key := range m.Calls {
		m.Calls[key] = 0
	}

	// Reset all of the gRPC methods to ensure that an RPC from a previous test doesn't
	// interfere with the operation of a current test.
	m.OnList = nil
	m.OnSummary = nil
}

// UseFixture allows you to specify a JSON fixture that is loaded from disk as the
// protocol buffer response for the specified RPC, simplifying the handler mocking.
func (m *Members) UseFixture(rpc, path string) (err error) {
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
	case ListRPC:
		out := &members.ListReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %s", out, err)
		}
		m.OnList = func(context.Context, *members.ListRequest) (*members.ListReply, error) {
			return out, nil
		}
	case SummaryRPC:
		out := &members.SummaryReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %s", out, err)
		}
		m.OnSummary = func(context.Context, *members.SummaryRequest) (*members.SummaryReply, error) {
			return out, nil
		}
	case DetailsRPC:
		out := &members.MemberDetails{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %s", out, err)
		}
		m.OnDetails = func(context.Context, *members.DetailsRequest) (*members.MemberDetails, error) {
			return out, nil
		}
	default:
		return fmt.Errorf("unknown rpc %q", rpc)
	}
	return nil
}

// UseError allows you to specify a gRPC status error to return from the specified RPC.
func (m *Members) UseError(rpc string, code codes.Code, msg string) error {
	switch rpc {
	case ListRPC:
		m.OnList = func(context.Context, *members.ListRequest) (*members.ListReply, error) {
			return nil, status.Error(code, msg)
		}
	case SummaryRPC:
		m.OnSummary = func(context.Context, *members.SummaryRequest) (*members.SummaryReply, error) {
			return nil, status.Error(code, msg)
		}
	case DetailsRPC:
		m.OnDetails = func(context.Context, *members.DetailsRequest) (*members.MemberDetails, error) {
			return nil, status.Error(code, msg)
		}
	default:
		return fmt.Errorf("unknown rpc %q", rpc)
	}
	return nil
}

func (m *Members) List(ctx context.Context, in *members.ListRequest) (*members.ListReply, error) {
	m.Calls[ListRPC]++
	return m.OnList(ctx, in)
}

func (m *Members) Summary(ctx context.Context, in *members.SummaryRequest) (*members.SummaryReply, error) {
	m.Calls[SummaryRPC]++
	return m.OnSummary(ctx, in)
}

func (m *Members) Details(ctx context.Context, in *members.DetailsRequest) (*members.MemberDetails, error) {
	m.Calls[DetailsRPC]++
	return m.OnDetails(ctx, in)
}
