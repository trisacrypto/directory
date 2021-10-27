// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: gds/admin/v1/admin.proto

package admin

import (
	v1beta1 "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	v1beta11 "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
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

type ResendRequest_ResendType int32

const (
	ResendRequest_UNKNOWN        ResendRequest_ResendType = 0
	ResendRequest_VERIFY_CONTACT ResendRequest_ResendType = 1
	ResendRequest_REVIEW         ResendRequest_ResendType = 2
	ResendRequest_DELIVER_CERTS  ResendRequest_ResendType = 3
	ResendRequest_REJECTION      ResendRequest_ResendType = 4
)

// Enum value maps for ResendRequest_ResendType.
var (
	ResendRequest_ResendType_name = map[int32]string{
		0: "UNKNOWN",
		1: "VERIFY_CONTACT",
		2: "REVIEW",
		3: "DELIVER_CERTS",
		4: "REJECTION",
	}
	ResendRequest_ResendType_value = map[string]int32{
		"UNKNOWN":        0,
		"VERIFY_CONTACT": 1,
		"REVIEW":         2,
		"DELIVER_CERTS":  3,
		"REJECTION":      4,
	}
)

func (x ResendRequest_ResendType) Enum() *ResendRequest_ResendType {
	p := new(ResendRequest_ResendType)
	*p = x
	return p
}

func (x ResendRequest_ResendType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResendRequest_ResendType) Descriptor() protoreflect.EnumDescriptor {
	return file_gds_admin_v1_admin_proto_enumTypes[0].Descriptor()
}

func (ResendRequest_ResendType) Type() protoreflect.EnumType {
	return &file_gds_admin_v1_admin_proto_enumTypes[0]
}

func (x ResendRequest_ResendType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResendRequest_ResendType.Descriptor instead.
func (ResendRequest_ResendType) EnumDescriptor() ([]byte, []int) {
	return file_gds_admin_v1_admin_proto_rawDescGZIP(), []int{2, 0}
}

// Registration review requests are sent via email to the TRISA admin email address with
// a lightweight token for review. This endpoint allows administrators to submit a review
// determination back to the directory server.
type ReviewRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The ID of the VASP to perform the review for.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The verification token sent in the review email.
	// This token provides lightweight authentication but should be replaced with a more
	// robust authentication and authorization scheme.
	AdminVerificationToken string `protobuf:"bytes,2,opt,name=admin_verification_token,json=adminVerificationToken,proto3" json:"admin_verification_token,omitempty"`
	// If accept is false then the request will be rejected and a reject reason must be
	// specified. If it is true, then the certificate issuance process will begin.
	Accept       bool   `protobuf:"varint,3,opt,name=accept,proto3" json:"accept,omitempty"`
	RejectReason string `protobuf:"bytes,4,opt,name=reject_reason,json=rejectReason,proto3" json:"reject_reason,omitempty"`
}

func (x *ReviewRequest) Reset() {
	*x = ReviewRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_admin_v1_admin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReviewRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReviewRequest) ProtoMessage() {}

func (x *ReviewRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gds_admin_v1_admin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReviewRequest.ProtoReflect.Descriptor instead.
func (*ReviewRequest) Descriptor() ([]byte, []int) {
	return file_gds_admin_v1_admin_proto_rawDescGZIP(), []int{0}
}

func (x *ReviewRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ReviewRequest) GetAdminVerificationToken() string {
	if x != nil {
		return x.AdminVerificationToken
	}
	return ""
}

func (x *ReviewRequest) GetAccept() bool {
	if x != nil {
		return x.Accept
	}
	return false
}

func (x *ReviewRequest) GetRejectReason() string {
	if x != nil {
		return x.RejectReason
	}
	return ""
}

type ReviewReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If no error is specified, the verify email request was successful
	Error *v1beta1.Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	// The verification status of the VASP entity.
	Status  v1beta11.VerificationState `protobuf:"varint,2,opt,name=status,proto3,enum=trisa.gds.models.v1beta1.VerificationState" json:"status,omitempty"`
	Message string                     `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ReviewReply) Reset() {
	*x = ReviewReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_admin_v1_admin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReviewReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReviewReply) ProtoMessage() {}

func (x *ReviewReply) ProtoReflect() protoreflect.Message {
	mi := &file_gds_admin_v1_admin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReviewReply.ProtoReflect.Descriptor instead.
func (*ReviewReply) Descriptor() ([]byte, []int) {
	return file_gds_admin_v1_admin_proto_rawDescGZIP(), []int{1}
}

func (x *ReviewReply) GetError() *v1beta1.Error {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *ReviewReply) GetStatus() v1beta11.VerificationState {
	if x != nil {
		return x.Status
	}
	return v1beta11.VerificationState_NO_VERIFICATION
}

func (x *ReviewReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Resend requests allow extra attempts to resend emails to be made if they were not
// delivered or recieved the first time. This is a routine action that may need to be
// carried out from time to time.
type ResendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string                   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`                                                 // The ID of the VASP to resend emails for
	Type   ResendRequest_ResendType `protobuf:"varint,2,opt,name=type,proto3,enum=gds.admin.v1.ResendRequest_ResendType" json:"type,omitempty"` // The type of message to attempt to resend
	Reason string                   `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`                                         // If a rejection email, supply the reason for the rejection.
}

func (x *ResendRequest) Reset() {
	*x = ResendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_admin_v1_admin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResendRequest) ProtoMessage() {}

func (x *ResendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gds_admin_v1_admin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResendRequest.ProtoReflect.Descriptor instead.
func (*ResendRequest) Descriptor() ([]byte, []int) {
	return file_gds_admin_v1_admin_proto_rawDescGZIP(), []int{2}
}

func (x *ResendRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ResendRequest) GetType() ResendRequest_ResendType {
	if x != nil {
		return x.Type
	}
	return ResendRequest_UNKNOWN
}

func (x *ResendRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type ResendReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sent    int64  `protobuf:"varint,1,opt,name=sent,proto3" json:"sent,omitempty"`      // The number of emails sent
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"` // Any message from the server about status
}

func (x *ResendReply) Reset() {
	*x = ResendReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_admin_v1_admin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResendReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResendReply) ProtoMessage() {}

func (x *ResendReply) ProtoReflect() protoreflect.Message {
	mi := &file_gds_admin_v1_admin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResendReply.ProtoReflect.Descriptor instead.
func (*ResendReply) Descriptor() ([]byte, []int) {
	return file_gds_admin_v1_admin_proto_rawDescGZIP(), []int{3}
}

func (x *ResendReply) GetSent() int64 {
	if x != nil {
		return x.Sent
	}
	return 0
}

func (x *ResendReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type StatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NoRegistrations       bool `protobuf:"varint,1,opt,name=no_registrations,json=noRegistrations,proto3" json:"no_registrations,omitempty"`                     // Ignore counting the registration statuses
	NoCertificateRequests bool `protobuf:"varint,2,opt,name=no_certificate_requests,json=noCertificateRequests,proto3" json:"no_certificate_requests,omitempty"` // Ignore counting certificate request statuses
}

func (x *StatusRequest) Reset() {
	*x = StatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_admin_v1_admin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusRequest) ProtoMessage() {}

func (x *StatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gds_admin_v1_admin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusRequest.ProtoReflect.Descriptor instead.
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return file_gds_admin_v1_admin_proto_rawDescGZIP(), []int{4}
}

func (x *StatusRequest) GetNoRegistrations() bool {
	if x != nil {
		return x.NoRegistrations
	}
	return false
}

func (x *StatusRequest) GetNoCertificateRequests() bool {
	if x != nil {
		return x.NoCertificateRequests
	}
	return false
}

type StatusReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Registrations       map[string]int64 `protobuf:"bytes,1,rep,name=registrations,proto3" json:"registrations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	CertificateRequests map[string]int64 `protobuf:"bytes,2,rep,name=certificate_requests,json=certificateRequests,proto3" json:"certificate_requests,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *StatusReply) Reset() {
	*x = StatusReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_admin_v1_admin_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusReply) ProtoMessage() {}

func (x *StatusReply) ProtoReflect() protoreflect.Message {
	mi := &file_gds_admin_v1_admin_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusReply.ProtoReflect.Descriptor instead.
func (*StatusReply) Descriptor() ([]byte, []int) {
	return file_gds_admin_v1_admin_proto_rawDescGZIP(), []int{5}
}

func (x *StatusReply) GetRegistrations() map[string]int64 {
	if x != nil {
		return x.Registrations
	}
	return nil
}

func (x *StatusReply) GetCertificateRequests() map[string]int64 {
	if x != nil {
		return x.CertificateRequests
	}
	return nil
}

var File_gds_admin_v1_admin_proto protoreflect.FileDescriptor

var file_gds_admin_v1_admin_proto_rawDesc = []byte{
	0x0a, 0x18, 0x67, 0x64, 0x73, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x67, 0x64, 0x73, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f,
	0x67, 0x64, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f,
	0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x74, 0x72, 0x69, 0x73, 0x61,
	0x2f, 0x67, 0x64, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x96, 0x01, 0x0a, 0x0d, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x38, 0x0a, 0x18, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x76, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x16, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x61, 0x63,
	0x63, 0x65, 0x70, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x6a,
	0x65, 0x63, 0x74, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0xa0, 0x01, 0x0a, 0x0b, 0x52, 0x65,
	0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61,
	0x2e, 0x67, 0x64, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x43, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2b, 0x2e,
	0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xd0, 0x01, 0x0a,
	0x0d, 0x52, 0x65, 0x73, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3a,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x67,
	0x64, 0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x65,
	0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x6e, 0x64,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x22, 0x5b, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x65, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x12, 0x0a,
	0x0e, 0x56, 0x45, 0x52, 0x49, 0x46, 0x59, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x41, 0x43, 0x54, 0x10,
	0x01, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x10, 0x02, 0x12, 0x11, 0x0a,
	0x0d, 0x44, 0x45, 0x4c, 0x49, 0x56, 0x45, 0x52, 0x5f, 0x43, 0x45, 0x52, 0x54, 0x53, 0x10, 0x03,
	0x12, 0x0d, 0x0a, 0x09, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x04, 0x22,
	0x3b, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x65,
	0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x72, 0x0a, 0x0d,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a,
	0x10, 0x6e, 0x6f, 0x5f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x6e, 0x6f, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x36, 0x0a, 0x17, 0x6e, 0x6f, 0x5f, 0x63,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x15, 0x6e, 0x6f, 0x43, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73,
	0x22, 0xd2, 0x02, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x52, 0x0a, 0x0d, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0d, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x65, 0x0a, 0x14, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x32, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x43, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x13, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x1a, 0x40, 0x0a, 0x12, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x46, 0x0a,
	0x18, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0xe5, 0x01, 0x0a, 0x17, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x79, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x42, 0x0a, 0x06, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x12, 0x1b, 0x2e, 0x67, 0x64,
	0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x76, 0x69, 0x65,
	0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x65, 0x6e, 0x64, 0x12,
	0x1b, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x65, 0x73, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x67,
	0x64, 0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x65,
	0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x1b, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x39, 0x5a,
	0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x72, 0x69, 0x73,
	0x61, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x64, 0x73, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f,
	0x76, 0x31, 0x3b, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gds_admin_v1_admin_proto_rawDescOnce sync.Once
	file_gds_admin_v1_admin_proto_rawDescData = file_gds_admin_v1_admin_proto_rawDesc
)

func file_gds_admin_v1_admin_proto_rawDescGZIP() []byte {
	file_gds_admin_v1_admin_proto_rawDescOnce.Do(func() {
		file_gds_admin_v1_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_gds_admin_v1_admin_proto_rawDescData)
	})
	return file_gds_admin_v1_admin_proto_rawDescData
}

var file_gds_admin_v1_admin_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gds_admin_v1_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_gds_admin_v1_admin_proto_goTypes = []interface{}{
	(ResendRequest_ResendType)(0),   // 0: gds.admin.v1.ResendRequest.ResendType
	(*ReviewRequest)(nil),           // 1: gds.admin.v1.ReviewRequest
	(*ReviewReply)(nil),             // 2: gds.admin.v1.ReviewReply
	(*ResendRequest)(nil),           // 3: gds.admin.v1.ResendRequest
	(*ResendReply)(nil),             // 4: gds.admin.v1.ResendReply
	(*StatusRequest)(nil),           // 5: gds.admin.v1.StatusRequest
	(*StatusReply)(nil),             // 6: gds.admin.v1.StatusReply
	nil,                             // 7: gds.admin.v1.StatusReply.RegistrationsEntry
	nil,                             // 8: gds.admin.v1.StatusReply.CertificateRequestsEntry
	(*v1beta1.Error)(nil),           // 9: trisa.gds.api.v1beta1.Error
	(v1beta11.VerificationState)(0), // 10: trisa.gds.models.v1beta1.VerificationState
}
var file_gds_admin_v1_admin_proto_depIdxs = []int32{
	9,  // 0: gds.admin.v1.ReviewReply.error:type_name -> trisa.gds.api.v1beta1.Error
	10, // 1: gds.admin.v1.ReviewReply.status:type_name -> trisa.gds.models.v1beta1.VerificationState
	0,  // 2: gds.admin.v1.ResendRequest.type:type_name -> gds.admin.v1.ResendRequest.ResendType
	7,  // 3: gds.admin.v1.StatusReply.registrations:type_name -> gds.admin.v1.StatusReply.RegistrationsEntry
	8,  // 4: gds.admin.v1.StatusReply.certificate_requests:type_name -> gds.admin.v1.StatusReply.CertificateRequestsEntry
	1,  // 5: gds.admin.v1.DirectoryAdministration.Review:input_type -> gds.admin.v1.ReviewRequest
	3,  // 6: gds.admin.v1.DirectoryAdministration.Resend:input_type -> gds.admin.v1.ResendRequest
	5,  // 7: gds.admin.v1.DirectoryAdministration.Status:input_type -> gds.admin.v1.StatusRequest
	2,  // 8: gds.admin.v1.DirectoryAdministration.Review:output_type -> gds.admin.v1.ReviewReply
	4,  // 9: gds.admin.v1.DirectoryAdministration.Resend:output_type -> gds.admin.v1.ResendReply
	6,  // 10: gds.admin.v1.DirectoryAdministration.Status:output_type -> gds.admin.v1.StatusReply
	8,  // [8:11] is the sub-list for method output_type
	5,  // [5:8] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_gds_admin_v1_admin_proto_init() }
func file_gds_admin_v1_admin_proto_init() {
	if File_gds_admin_v1_admin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gds_admin_v1_admin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReviewRequest); i {
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
		file_gds_admin_v1_admin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReviewReply); i {
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
		file_gds_admin_v1_admin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResendRequest); i {
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
		file_gds_admin_v1_admin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResendReply); i {
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
		file_gds_admin_v1_admin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusRequest); i {
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
		file_gds_admin_v1_admin_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusReply); i {
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
			RawDescriptor: file_gds_admin_v1_admin_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gds_admin_v1_admin_proto_goTypes,
		DependencyIndexes: file_gds_admin_v1_admin_proto_depIdxs,
		EnumInfos:         file_gds_admin_v1_admin_proto_enumTypes,
		MessageInfos:      file_gds_admin_v1_admin_proto_msgTypes,
	}.Build()
	File_gds_admin_v1_admin_proto = out.File
	file_gds_admin_v1_admin_proto_rawDesc = nil
	file_gds_admin_v1_admin_proto_goTypes = nil
	file_gds_admin_v1_admin_proto_depIdxs = nil
}
