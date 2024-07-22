package mock

import (
	"context"
	"fmt"
	"os"

	"github.com/rotationalio/honu/replica"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/directory/pkg/trtl/peers/v1"
	"github.com/trisacrypto/directory/pkg/utils/bufconn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	GetRPC      = "trtl.v1.Trtl/Get"
	PutRPC      = "trtl.v1.Trtl/Put"
	DeleteRPC   = "trtl.v1.Trtl/Delete"
	IterRPC     = "trtl.v1.Trtl/Iter"
	BatchRPC    = "trtl.v1.Trtl/Batch"
	CursorRPC   = "trtl.v1.Trtl/Cursor"
	SyncRPC     = "trtl.v1.Trtl/Sync"
	StatusRPC   = "trtl.v1.Trtl/Status"
	GetPeersRPC = "trtl.peers.v1.PeerManagement/GetPeers"
	AddPeersRPC = "trtl.peers.v1.PeerManagement/AddPeers"
	RmPeersRPC  = "trtl.peers.v1.PeerManagement/RmPeers"
	GossipRPC   = "honu.replica.v1.Replication/Gossip"
)

// New creates a new mock RemoteTrtl. If bufnet is nil, one is created for the user.
// The gRPC server defaults to insecure which can be overridden by passing in
// ServerOptions with configured TLS.
func New(bufnet *bufconn.GRPCListener, opts ...grpc.ServerOption) *RemoteTrtl {
	if bufnet == nil {
		bufnet = bufconn.New("")
	}

	if len(opts) == 0 {
		opts = make([]grpc.ServerOption, 0)
		opts = append(opts, grpc.Creds(insecure.NewCredentials()))
	}

	remote := &RemoteTrtl{
		bufnet: bufnet,
		srv:    grpc.NewServer(opts...),
		Calls:  make(map[string]int),
	}

	pb.RegisterTrtlServer(remote.srv, remote)
	peers.RegisterPeerManagementServer(remote.srv, remote)
	go remote.srv.Serve(remote.bufnet.Listener)
	return remote
}

// RemoteTrtl implements a mock gRPC server for testing client connections to other
// trtl peers. The desired response of the remote peer can be set by external callers
// using the OnRPC functions or the WithFixture or WithError functions. The Calls map
// can be used to count the number of times the remote peer RPC was called.
type RemoteTrtl struct {
	pb.UnimplementedTrtlServer
	peers.UnimplementedPeerManagementServer
	replica.UnimplementedReplicationServer
	bufnet     *bufconn.GRPCListener
	srv        *grpc.Server
	Calls      map[string]int
	OnGet      func(context.Context, *pb.GetRequest) (*pb.GetReply, error)
	OnPut      func(context.Context, *pb.PutRequest) (*pb.PutReply, error)
	OnDelete   func(context.Context, *pb.DeleteRequest) (*pb.DeleteReply, error)
	OnIter     func(context.Context, *pb.IterRequest) (*pb.IterReply, error)
	OnBatch    func(pb.Trtl_BatchServer) error
	OnCursor   func(*pb.CursorRequest, pb.Trtl_CursorServer) error
	OnSync     func(pb.Trtl_SyncServer) error
	OnStatus   func(context.Context, *pb.HealthCheck) (*pb.ServerStatus, error)
	OnGetPeers func(context.Context, *peers.PeersFilter) (*peers.PeersList, error)
	OnAddPeers func(context.Context, *peers.Peer) (*peers.PeersStatus, error)
	OnRmPeers  func(context.Context, *peers.Peer) (*peers.PeersStatus, error)
	OnGossip   func(replica.Replication_GossipServer) error
}

func (s *RemoteTrtl) Channel() *bufconn.GRPCListener {
	return s.bufnet
}

func (s *RemoteTrtl) DBClient(opts ...grpc.DialOption) (_ pb.TrtlClient, err error) {
	if err = s.connect(opts...); err != nil {
		return nil, err
	}
	return pb.NewTrtlClient(s.bufnet.Conn), nil
}

func (s *RemoteTrtl) PeersClient(opts ...grpc.DialOption) (_ peers.PeerManagementClient, err error) {
	if err = s.connect(opts...); err != nil {
		return nil, err
	}
	return peers.NewPeerManagementClient(s.bufnet.Conn), nil
}

func (s *RemoteTrtl) ReplicaClient(opts ...grpc.DialOption) (_ replica.ReplicationClient, err error) {
	if err = s.connect(opts...); err != nil {
		return nil, err
	}
	return replica.NewReplicationClient(s.bufnet.Conn), nil
}

func (s *RemoteTrtl) connect(opts ...grpc.DialOption) (err error) {
	// Note: The bufconn package only supports one client connection at a time and is
	// not thread-safe.
	if s.bufnet.Conn != nil {
		s.bufnet.Close()
	}
	if err = s.bufnet.Connect(context.Background(), opts...); err != nil {
		return err
	}
	return nil
}

func (s *RemoteTrtl) CloseClient() {
	if s.bufnet.Conn != nil {
		s.bufnet.Close()
	}
}

func (s *RemoteTrtl) Shutdown() {
	s.srv.GracefulStop()
	s.CloseClient()
	s.bufnet.Release()
}

func (s *RemoteTrtl) Reset() {
	for key := range s.Calls {
		s.Calls[key] = 0
	}

	s.OnGet = nil
	s.OnPut = nil
	s.OnDelete = nil
	s.OnIter = nil
	s.OnBatch = nil
	s.OnCursor = nil
	s.OnSync = nil
	s.OnStatus = nil
	s.OnGetPeers = nil
	s.OnAddPeers = nil
	s.OnRmPeers = nil
	s.OnGossip = nil
}

// UseFixture loadsa a JSON fixture from disk (usually in a testdata folder) to use as
// the protocol buffer response to the specified RPC, simplifying handler mocking.
func (s *RemoteTrtl) UseFixture(rpc, path string) (err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return fmt.Errorf("could not read fixture: %v", err)
	}

	jsonpb := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	switch rpc {
	case GetRPC:
		out := &pb.GetReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
			return out, nil
		}
	case PutRPC:
		out := &pb.PutReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
			return out, nil
		}
	case DeleteRPC:
		out := &pb.DeleteReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
			return out, nil
		}
	case IterRPC:
		out := &pb.IterReply{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnIter = func(ctx context.Context, in *pb.IterRequest) (*pb.IterReply, error) {
			return out, nil
		}
	case BatchRPC:
		return fmt.Errorf("cannot use fixture for Batch RPC, instead set OnBatchRPC directly")
	case CursorRPC:
		return fmt.Errorf("cannot use fixture for Cursor RPC, instead set OnCursorRPC directly")
	case SyncRPC:
		return fmt.Errorf("cannot use fixture for Sync RPC, instead set OnSyncRPC directly")
	case StatusRPC:
		out := &pb.ServerStatus{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnStatus = func(ctx context.Context, in *pb.HealthCheck) (*pb.ServerStatus, error) {
			return out, nil
		}
	case GetPeersRPC:
		out := &peers.PeersList{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnGetPeers = func(ctx context.Context, in *peers.PeersFilter) (*peers.PeersList, error) {
			return out, nil
		}
	case AddPeersRPC:
		out := &peers.PeersStatus{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnAddPeers = func(ctx context.Context, in *peers.Peer) (*peers.PeersStatus, error) {
			return out, nil
		}
	case RmPeersRPC:
		out := &peers.PeersStatus{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnRmPeers = func(ctx context.Context, in *peers.Peer) (*peers.PeersStatus, error) {
			return out, nil
		}
	case GossipRPC:
		return fmt.Errorf("cannot use fixture for Gossip RPC, instead set OnGossip directly")
	default:
		return fmt.Errorf("unknown RPC %q", rpc)
	}

	return nil
}

// UseError allows you to specify a gRPC status error to return from the specified RPC.
func (s *RemoteTrtl) UseError(rpc string, code codes.Code, msg string) error {
	switch rpc {
	case GetRPC:
		s.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
			return nil, status.Error(code, msg)
		}
	case PutRPC:
		s.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
			return nil, status.Error(code, msg)
		}
	case DeleteRPC:
		s.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
			return nil, status.Error(code, msg)
		}
	case IterRPC:
		s.OnIter = func(ctx context.Context, in *pb.IterRequest) (*pb.IterReply, error) {
			return nil, status.Error(code, msg)
		}
	case BatchRPC:
		s.OnBatch = func(pb.Trtl_BatchServer) error {
			return status.Error(code, msg)
		}
	case CursorRPC:
		s.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
			return status.Error(code, msg)
		}
	case SyncRPC:
		s.OnSync = func(pb.Trtl_SyncServer) error {
			return status.Error(code, msg)
		}
	case StatusRPC:
		s.OnStatus = func(ctx context.Context, in *pb.HealthCheck) (*pb.ServerStatus, error) {
			return nil, status.Error(code, msg)
		}
	case GetPeersRPC:
		s.OnGetPeers = func(ctx context.Context, in *peers.PeersFilter) (*peers.PeersList, error) {
			return nil, status.Error(code, msg)
		}
	case AddPeersRPC:
		s.OnAddPeers = func(ctx context.Context, in *peers.Peer) (*peers.PeersStatus, error) {
			return nil, status.Error(code, msg)
		}
	case RmPeersRPC:
		s.OnRmPeers = func(ctx context.Context, in *peers.Peer) (*peers.PeersStatus, error) {
			return nil, status.Error(code, msg)
		}
	case GossipRPC:
		s.OnGossip = func(replica.Replication_GossipServer) error {
			return status.Error(code, msg)
		}
	default:
		return fmt.Errorf("unknown RPC %q", rpc)
	}
	return nil
}

func (s *RemoteTrtl) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	s.Calls[GetRPC]++
	return s.OnGet(ctx, in)
}

func (s *RemoteTrtl) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
	s.Calls[PutRPC]++
	return s.OnPut(ctx, in)
}

func (s *RemoteTrtl) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
	s.Calls[DeleteRPC]++
	return s.OnDelete(ctx, in)
}

func (s *RemoteTrtl) Iter(ctx context.Context, in *pb.IterRequest) (*pb.IterReply, error) {
	s.Calls[IterRPC]++
	return s.OnIter(ctx, in)
}

func (s *RemoteTrtl) Batch(stream pb.Trtl_BatchServer) error {
	s.Calls[BatchRPC]++
	return s.OnBatch(stream)
}

func (s *RemoteTrtl) Cursor(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
	s.Calls[CursorRPC]++
	return s.OnCursor(in, stream)
}

func (s *RemoteTrtl) Sync(stream pb.Trtl_SyncServer) error {
	s.Calls[SyncRPC]++
	return s.OnSync(stream)
}

func (s *RemoteTrtl) Status(ctx context.Context, in *pb.HealthCheck) (*pb.ServerStatus, error) {
	s.Calls[StatusRPC]++
	return s.OnStatus(ctx, in)
}

func (s *RemoteTrtl) GetPeers(ctx context.Context, in *peers.PeersFilter) (*peers.PeersList, error) {
	s.Calls[GetPeersRPC]++
	return s.OnGetPeers(ctx, in)
}

func (s *RemoteTrtl) AddPeers(ctx context.Context, in *peers.Peer) (*peers.PeersStatus, error) {
	s.Calls[AddPeersRPC]++
	return s.OnAddPeers(ctx, in)
}

func (s *RemoteTrtl) RmPeers(ctx context.Context, in *peers.Peer) (*peers.PeersStatus, error) {
	s.Calls[RmPeersRPC]++
	return s.OnRmPeers(ctx, in)
}

func (s *RemoteTrtl) Gossip(stream replica.Replication_GossipServer) error {
	s.Calls[GossipRPC]++
	return s.OnGossip(stream)
}
