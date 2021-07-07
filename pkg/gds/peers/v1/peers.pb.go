// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: gds/peers/v1/peers.proto

package peers

import (
	v1 "github.com/trisacrypto/directory/pkg/gds/global/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Peer contains metadata about how to connect to remote peers in the directory service
// network. This message services as a data-transfer and exchange mechanism for dynamic
// networks with changing membership.
type Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`        // the process id of the peer must be unique in the network; used for distributed versions
	Addr   string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`     // the network address to connect to the peer on (don't forget the port!)
	Name   string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`     // optional - a unique, human readable name for the peer
	Region string `protobuf:"bytes,4,opt,name=region,proto3" json:"region,omitempty"` // optional? - the region the peer is running in
	// Logging information timestamps
	Created  string `protobuf:"bytes,9,opt,name=created,proto3" json:"created,omitempty"`
	Modified string `protobuf:"bytes,10,opt,name=modified,proto3" json:"modified,omitempty"`
	Deleted  string `protobuf:"bytes,11,opt,name=deleted,proto3" json:"deleted,omitempty"`
	// extra information that might be relevant to process-specific functions; e.g. for
	// specific clouds or data that's been parsed (optional).
	Extra    map[string]string `protobuf:"bytes,14,rep,name=extra,proto3" json:"extra,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Metadata *v1.Object        `protobuf:"bytes,15,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *Peer) Reset() {
	*x = Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_peers_v1_peers_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Peer) ProtoMessage() {}

func (x *Peer) ProtoReflect() protoreflect.Message {
	mi := &file_gds_peers_v1_peers_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Peer.ProtoReflect.Descriptor instead.
func (*Peer) Descriptor() ([]byte, []int) {
	return file_gds_peers_v1_peers_proto_rawDescGZIP(), []int{0}
}

func (x *Peer) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Peer) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Peer) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Peer) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Peer) GetCreated() string {
	if x != nil {
		return x.Created
	}
	return ""
}

func (x *Peer) GetModified() string {
	if x != nil {
		return x.Modified
	}
	return ""
}

func (x *Peer) GetDeleted() string {
	if x != nil {
		return x.Deleted
	}
	return ""
}

func (x *Peer) GetExtra() map[string]string {
	if x != nil {
		return x.Extra
	}
	return nil
}

func (x *Peer) GetMetadata() *v1.Object {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// Used to filter the peers that are returned. If no filters are specified then all
// known peers on the remote replica are returned.
type PeersFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Region     []string `protobuf:"bytes,1,rep,name=region,proto3" json:"region,omitempty"`                            // Specify the region(s) to return the peers for. Only effects PeersList not PeersStatus
	StatusOnly bool     `protobuf:"varint,2,opt,name=status_only,json=statusOnly,proto3" json:"status_only,omitempty"` // Return only the peers status, not a list of peers.
}

func (x *PeersFilter) Reset() {
	*x = PeersFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_peers_v1_peers_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeersFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeersFilter) ProtoMessage() {}

func (x *PeersFilter) ProtoReflect() protoreflect.Message {
	mi := &file_gds_peers_v1_peers_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeersFilter.ProtoReflect.Descriptor instead.
func (*PeersFilter) Descriptor() ([]byte, []int) {
	return file_gds_peers_v1_peers_proto_rawDescGZIP(), []int{1}
}

func (x *PeersFilter) GetRegion() []string {
	if x != nil {
		return x.Region
	}
	return nil
}

func (x *PeersFilter) GetStatusOnly() bool {
	if x != nil {
		return x.StatusOnly
	}
	return false
}

// Returns the list of peers currently known to the replica and its peer management status.
type PeersList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Peers  []*Peer      `protobuf:"bytes,1,rep,name=peers,proto3" json:"peers,omitempty"`
	Status *PeersStatus `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *PeersList) Reset() {
	*x = PeersList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_peers_v1_peers_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeersList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeersList) ProtoMessage() {}

func (x *PeersList) ProtoReflect() protoreflect.Message {
	mi := &file_gds_peers_v1_peers_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeersList.ProtoReflect.Descriptor instead.
func (*PeersList) Descriptor() ([]byte, []int) {
	return file_gds_peers_v1_peers_proto_rawDescGZIP(), []int{2}
}

func (x *PeersList) GetPeers() []*Peer {
	if x != nil {
		return x.Peers
	}
	return nil
}

func (x *PeersList) GetStatus() *PeersStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

// A response to a peer management command that describes the current state of the network.
type PeersStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetworkSize int64            `protobuf:"varint,1,opt,name=network_size,json=networkSize,proto3" json:"network_size,omitempty"`                                                              // The total number of peers known to the replica (including itself)
	Regions     map[string]int64 `protobuf:"bytes,2,rep,name=regions,proto3" json:"regions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"` // The number of peers known to the replica per known region
}

func (x *PeersStatus) Reset() {
	*x = PeersStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_peers_v1_peers_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeersStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeersStatus) ProtoMessage() {}

func (x *PeersStatus) ProtoReflect() protoreflect.Message {
	mi := &file_gds_peers_v1_peers_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeersStatus.ProtoReflect.Descriptor instead.
func (*PeersStatus) Descriptor() ([]byte, []int) {
	return file_gds_peers_v1_peers_proto_rawDescGZIP(), []int{3}
}

func (x *PeersStatus) GetNetworkSize() int64 {
	if x != nil {
		return x.NetworkSize
	}
	return 0
}

func (x *PeersStatus) GetRegions() map[string]int64 {
	if x != nil {
		return x.Regions
	}
	return nil
}

var File_gds_peers_v1_peers_proto protoreflect.FileDescriptor

var file_gds_peers_v1_peers_proto_rawDesc = []byte{
	0x0a, 0x18, 0x67, 0x64, 0x73, 0x2f, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x65, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x67, 0x64, 0x73, 0x2e,
	0x70, 0x65, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1a, 0x67, 0x64, 0x73, 0x2f, 0x67, 0x6c,
	0x6f, 0x62, 0x61, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc8, 0x02, 0x0a, 0x04, 0x50, 0x65, 0x65, 0x72, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x6f, 0x64, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x6f, 0x64, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x33, 0x0a,
	0x05, 0x65, 0x78, 0x74, 0x72, 0x61, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67,
	0x64, 0x73, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65, 0x65, 0x72,
	0x2e, 0x45, 0x78, 0x74, 0x72, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x65, 0x78, 0x74,
	0x72, 0x61, 0x12, 0x31, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61,
	0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x38, 0x0a, 0x0a, 0x45, 0x78, 0x74, 0x72, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x46, 0x0a, 0x0b, 0x50, 0x65, 0x65, 0x72, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x4f, 0x6e, 0x6c, 0x79, 0x22, 0x68, 0x0a, 0x09, 0x50, 0x65, 0x65, 0x72, 0x73,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x52, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x12, 0x31,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x67, 0x64, 0x73, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65,
	0x65, 0x72, 0x73, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0xae, 0x01, 0x0a, 0x0b, 0x50, 0x65, 0x65, 0x72, 0x73, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x53, 0x69, 0x7a, 0x65, 0x12, 0x40, 0x0a, 0x07, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x70, 0x65, 0x65, 0x72,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x73, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x72,
	0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x32, 0xcb, 0x01, 0x0a, 0x0e, 0x50, 0x65, 0x65, 0x72, 0x4d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x40, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x50, 0x65, 0x65, 0x72,
	0x73, 0x12, 0x19, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x65, 0x65, 0x72, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x1a, 0x17, 0x2e, 0x67,
	0x64, 0x73, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65, 0x65, 0x72,
	0x73, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x08, 0x41, 0x64, 0x64, 0x50, 0x65,
	0x65, 0x72, 0x73, 0x12, 0x12, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x1a, 0x19, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x70, 0x65,
	0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x73, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x07, 0x52, 0x6d, 0x50, 0x65, 0x65, 0x72, 0x73, 0x12,
	0x12, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x65, 0x65, 0x72, 0x1a, 0x19, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x73, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00,
	0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74,
	0x72, 0x69, 0x73, 0x61, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x64, 0x73, 0x2f, 0x70, 0x65, 0x65,
	0x72, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x65, 0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_gds_peers_v1_peers_proto_rawDescOnce sync.Once
	file_gds_peers_v1_peers_proto_rawDescData = file_gds_peers_v1_peers_proto_rawDesc
)

func file_gds_peers_v1_peers_proto_rawDescGZIP() []byte {
	file_gds_peers_v1_peers_proto_rawDescOnce.Do(func() {
		file_gds_peers_v1_peers_proto_rawDescData = protoimpl.X.CompressGZIP(file_gds_peers_v1_peers_proto_rawDescData)
	})
	return file_gds_peers_v1_peers_proto_rawDescData
}

var file_gds_peers_v1_peers_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_gds_peers_v1_peers_proto_goTypes = []interface{}{
	(*Peer)(nil),        // 0: gds.peers.v1.Peer
	(*PeersFilter)(nil), // 1: gds.peers.v1.PeersFilter
	(*PeersList)(nil),   // 2: gds.peers.v1.PeersList
	(*PeersStatus)(nil), // 3: gds.peers.v1.PeersStatus
	nil,                 // 4: gds.peers.v1.Peer.ExtraEntry
	nil,                 // 5: gds.peers.v1.PeersStatus.RegionsEntry
	(*v1.Object)(nil),   // 6: gds.global.v1.Object
}
var file_gds_peers_v1_peers_proto_depIdxs = []int32{
	4, // 0: gds.peers.v1.Peer.extra:type_name -> gds.peers.v1.Peer.ExtraEntry
	6, // 1: gds.peers.v1.Peer.metadata:type_name -> gds.global.v1.Object
	0, // 2: gds.peers.v1.PeersList.peers:type_name -> gds.peers.v1.Peer
	3, // 3: gds.peers.v1.PeersList.status:type_name -> gds.peers.v1.PeersStatus
	5, // 4: gds.peers.v1.PeersStatus.regions:type_name -> gds.peers.v1.PeersStatus.RegionsEntry
	1, // 5: gds.peers.v1.PeerManagement.GetPeers:input_type -> gds.peers.v1.PeersFilter
	0, // 6: gds.peers.v1.PeerManagement.AddPeers:input_type -> gds.peers.v1.Peer
	0, // 7: gds.peers.v1.PeerManagement.RmPeers:input_type -> gds.peers.v1.Peer
	2, // 8: gds.peers.v1.PeerManagement.GetPeers:output_type -> gds.peers.v1.PeersList
	3, // 9: gds.peers.v1.PeerManagement.AddPeers:output_type -> gds.peers.v1.PeersStatus
	3, // 10: gds.peers.v1.PeerManagement.RmPeers:output_type -> gds.peers.v1.PeersStatus
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_gds_peers_v1_peers_proto_init() }
func file_gds_peers_v1_peers_proto_init() {
	if File_gds_peers_v1_peers_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gds_peers_v1_peers_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Peer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gds_peers_v1_peers_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeersFilter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gds_peers_v1_peers_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeersList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gds_peers_v1_peers_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeersStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gds_peers_v1_peers_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gds_peers_v1_peers_proto_goTypes,
		DependencyIndexes: file_gds_peers_v1_peers_proto_depIdxs,
		MessageInfos:      file_gds_peers_v1_peers_proto_msgTypes,
	}.Build()
	File_gds_peers_v1_peers_proto = out.File
	file_gds_peers_v1_peers_proto_rawDesc = nil
	file_gds_peers_v1_peers_proto_goTypes = nil
	file_gds_peers_v1_peers_proto_depIdxs = nil
}
