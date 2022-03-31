// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: gds/members/v1alpha1/members.proto

package members

import (
	v1beta1 "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
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

// ListRequest manages paginating the VASP listing. If there are more results than the
// specified page size, then the ListReply will return a page token; that token can be
// used to fetch the next page so long as the parameters of the original request are not
// modified (e.g. any filters or pagination parameters).
// See https://cloud.google.com/apis/design/design_patterns#list_pagination for more.
type ListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageSize  int32  `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`   // specify the number of results per page, cannot change between page requests (default 100)
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"` // specify the page token to fetch the next page of results
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_members_v1alpha1_members_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gds_members_v1alpha1_members_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRequest.ProtoReflect.Descriptor instead.
func (*ListRequest) Descriptor() ([]byte, []int) {
	return file_gds_members_v1alpha1_members_proto_rawDescGZIP(), []int{0}
}

func (x *ListRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

// ListReply returns an abbreviated listing of VASP details intended to facilitate p2p
// key exchanges or more detailed lookups against the Directory Service.
type ListReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vasps         []*VASPMember `protobuf:"bytes,1,rep,name=vasps,proto3" json:"vasps,omitempty"`                                        // a list of VASP information for the requested page
	NextPageToken string        `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"` // if specified, another page of results exists
}

func (x *ListReply) Reset() {
	*x = ListReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_members_v1alpha1_members_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListReply) ProtoMessage() {}

func (x *ListReply) ProtoReflect() protoreflect.Message {
	mi := &file_gds_members_v1alpha1_members_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListReply.ProtoReflect.Descriptor instead.
func (*ListReply) Descriptor() ([]byte, []int) {
	return file_gds_members_v1alpha1_members_proto_rawDescGZIP(), []int{1}
}

func (x *ListReply) GetVasps() []*VASPMember {
	if x != nil {
		return x.Vasps
	}
	return nil
}

func (x *ListReply) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

// VASPMember is a lightweight data structure containing enough information to
// facilitate p2p exchanges or more detailed lookups against the Directory Service.
type VASPMember struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The uniquely identifying components of the VASP in the directory service
	Id                  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RegisteredDirectory string `protobuf:"bytes,2,opt,name=registered_directory,json=registeredDirectory,proto3" json:"registered_directory,omitempty"`
	CommonName          string `protobuf:"bytes,3,opt,name=common_name,json=commonName,proto3" json:"common_name,omitempty"`
	// Address to connect to the remote VASP on to perform a TRISA request
	Endpoint string `protobuf:"bytes,4,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// Extra details used to faciliate searches and matching
	Name             string                   `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Website          string                   `protobuf:"bytes,6,opt,name=website,proto3" json:"website,omitempty"`
	Country          string                   `protobuf:"bytes,7,opt,name=country,proto3" json:"country,omitempty"`
	BusinessCategory v1beta1.BusinessCategory `protobuf:"varint,8,opt,name=business_category,json=businessCategory,proto3,enum=trisa.gds.models.v1beta1.BusinessCategory" json:"business_category,omitempty"`
	VaspCategories   []string                 `protobuf:"bytes,9,rep,name=vasp_categories,json=vaspCategories,proto3" json:"vasp_categories,omitempty"`
	VerifiedOn       string                   `protobuf:"bytes,10,opt,name=verified_on,json=verifiedOn,proto3" json:"verified_on,omitempty"`
}

func (x *VASPMember) Reset() {
	*x = VASPMember{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_members_v1alpha1_members_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VASPMember) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VASPMember) ProtoMessage() {}

func (x *VASPMember) ProtoReflect() protoreflect.Message {
	mi := &file_gds_members_v1alpha1_members_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VASPMember.ProtoReflect.Descriptor instead.
func (*VASPMember) Descriptor() ([]byte, []int) {
	return file_gds_members_v1alpha1_members_proto_rawDescGZIP(), []int{2}
}

func (x *VASPMember) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *VASPMember) GetRegisteredDirectory() string {
	if x != nil {
		return x.RegisteredDirectory
	}
	return ""
}

func (x *VASPMember) GetCommonName() string {
	if x != nil {
		return x.CommonName
	}
	return ""
}

func (x *VASPMember) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *VASPMember) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *VASPMember) GetWebsite() string {
	if x != nil {
		return x.Website
	}
	return ""
}

func (x *VASPMember) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *VASPMember) GetBusinessCategory() v1beta1.BusinessCategory {
	if x != nil {
		return x.BusinessCategory
	}
	return v1beta1.BusinessCategory(0)
}

func (x *VASPMember) GetVaspCategories() []string {
	if x != nil {
		return x.VaspCategories
	}
	return nil
}

func (x *VASPMember) GetVerifiedOn() string {
	if x != nil {
		return x.VerifiedOn
	}
	return ""
}

var File_gds_members_v1alpha1_members_proto protoreflect.FileDescriptor

var file_gds_members_v1alpha1_members_proto_rawDesc = []byte{
	0x0a, 0x22, 0x67, 0x64, 0x73, 0x2f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x2f, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x25, 0x74, 0x72, 0x69, 0x73,
	0x61, 0x2f, 0x67, 0x64, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x49, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x6b, 0x0a, 0x09,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x36, 0x0a, 0x05, 0x76, 0x61, 0x73,
	0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e,
	0x56, 0x41, 0x53, 0x50, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x05, 0x76, 0x61, 0x73, 0x70,
	0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74,
	0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xf7, 0x02, 0x0a, 0x0a, 0x56, 0x41,
	0x53, 0x50, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x31, 0x0a, 0x14, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x5f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x65, 0x64, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x77,
	0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72,
	0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x57, 0x0a, 0x11, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x5f, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2a, 0x2e, 0x74, 0x72,
	0x69, 0x73, 0x61, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x42, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x43,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x10, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73,
	0x73, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x27, 0x0a, 0x0f, 0x76, 0x61, 0x73,
	0x70, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x09, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0e, 0x76, 0x61, 0x73, 0x70, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69,
	0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x5f, 0x6f,
	0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65,
	0x64, 0x4f, 0x6e, 0x32, 0x5c, 0x0a, 0x0c, 0x54, 0x52, 0x49, 0x53, 0x41, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x73, 0x12, 0x4c, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x21, 0x2e, 0x67, 0x64,
	0x73, 0x2e, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f,
	0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x00, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x74, 0x72, 0x69, 0x73, 0x61, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2f, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x64, 0x73, 0x2f, 0x6d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x73, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gds_members_v1alpha1_members_proto_rawDescOnce sync.Once
	file_gds_members_v1alpha1_members_proto_rawDescData = file_gds_members_v1alpha1_members_proto_rawDesc
)

func file_gds_members_v1alpha1_members_proto_rawDescGZIP() []byte {
	file_gds_members_v1alpha1_members_proto_rawDescOnce.Do(func() {
		file_gds_members_v1alpha1_members_proto_rawDescData = protoimpl.X.CompressGZIP(file_gds_members_v1alpha1_members_proto_rawDescData)
	})
	return file_gds_members_v1alpha1_members_proto_rawDescData
}

var file_gds_members_v1alpha1_members_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_gds_members_v1alpha1_members_proto_goTypes = []interface{}{
	(*ListRequest)(nil),           // 0: gds.members.v1alpha1.ListRequest
	(*ListReply)(nil),             // 1: gds.members.v1alpha1.ListReply
	(*VASPMember)(nil),            // 2: gds.members.v1alpha1.VASPMember
	(v1beta1.BusinessCategory)(0), // 3: trisa.gds.models.v1beta1.BusinessCategory
}
var file_gds_members_v1alpha1_members_proto_depIdxs = []int32{
	2, // 0: gds.members.v1alpha1.ListReply.vasps:type_name -> gds.members.v1alpha1.VASPMember
	3, // 1: gds.members.v1alpha1.VASPMember.business_category:type_name -> trisa.gds.models.v1beta1.BusinessCategory
	0, // 2: gds.members.v1alpha1.TRISAMembers.List:input_type -> gds.members.v1alpha1.ListRequest
	1, // 3: gds.members.v1alpha1.TRISAMembers.List:output_type -> gds.members.v1alpha1.ListReply
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_gds_members_v1alpha1_members_proto_init() }
func file_gds_members_v1alpha1_members_proto_init() {
	if File_gds_members_v1alpha1_members_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gds_members_v1alpha1_members_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRequest); i {
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
		file_gds_members_v1alpha1_members_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListReply); i {
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
		file_gds_members_v1alpha1_members_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VASPMember); i {
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
			RawDescriptor: file_gds_members_v1alpha1_members_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gds_members_v1alpha1_members_proto_goTypes,
		DependencyIndexes: file_gds_members_v1alpha1_members_proto_depIdxs,
		MessageInfos:      file_gds_members_v1alpha1_members_proto_msgTypes,
	}.Build()
	File_gds_members_v1alpha1_members_proto = out.File
	file_gds_members_v1alpha1_members_proto_rawDesc = nil
	file_gds_members_v1alpha1_members_proto_goTypes = nil
	file_gds_members_v1alpha1_members_proto_depIdxs = nil
}
