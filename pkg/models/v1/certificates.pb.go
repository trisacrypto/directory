// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: gds/models/v1/certificates.proto

package models

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

type CertificateRequestState int32

const (
	CertificateRequestState_INITIALIZED     CertificateRequestState = 0
	CertificateRequestState_READY_TO_SUBMIT CertificateRequestState = 1
	CertificateRequestState_PROCESSING      CertificateRequestState = 2
	CertificateRequestState_DOWNLOADING     CertificateRequestState = 3
	CertificateRequestState_DOWNLOADED      CertificateRequestState = 4
	CertificateRequestState_COMPLETED       CertificateRequestState = 5
	CertificateRequestState_CR_REJECTED     CertificateRequestState = 6
	CertificateRequestState_CR_ERRORED      CertificateRequestState = 7
)

// Enum value maps for CertificateRequestState.
var (
	CertificateRequestState_name = map[int32]string{
		0: "INITIALIZED",
		1: "READY_TO_SUBMIT",
		2: "PROCESSING",
		3: "DOWNLOADING",
		4: "DOWNLOADED",
		5: "COMPLETED",
		6: "CR_REJECTED",
		7: "CR_ERRORED",
	}
	CertificateRequestState_value = map[string]int32{
		"INITIALIZED":     0,
		"READY_TO_SUBMIT": 1,
		"PROCESSING":      2,
		"DOWNLOADING":     3,
		"DOWNLOADED":      4,
		"COMPLETED":       5,
		"CR_REJECTED":     6,
		"CR_ERRORED":      7,
	}
)

func (x CertificateRequestState) Enum() *CertificateRequestState {
	p := new(CertificateRequestState)
	*p = x
	return p
}

func (x CertificateRequestState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CertificateRequestState) Descriptor() protoreflect.EnumDescriptor {
	return file_gds_models_v1_certificates_proto_enumTypes[0].Descriptor()
}

func (CertificateRequestState) Type() protoreflect.EnumType {
	return &file_gds_models_v1_certificates_proto_enumTypes[0]
}

func (x CertificateRequestState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CertificateRequestState.Descriptor instead.
func (CertificateRequestState) EnumDescriptor() ([]byte, []int) {
	return file_gds_models_v1_certificates_proto_rawDescGZIP(), []int{0}
}

type CertificateState int32

const (
	CertificateState_ISSUED  CertificateState = 0
	CertificateState_EXPIRED CertificateState = 1
	CertificateState_REVOKED CertificateState = 2
)

// Enum value maps for CertificateState.
var (
	CertificateState_name = map[int32]string{
		0: "ISSUED",
		1: "EXPIRED",
		2: "REVOKED",
	}
	CertificateState_value = map[string]int32{
		"ISSUED":  0,
		"EXPIRED": 1,
		"REVOKED": 2,
	}
)

func (x CertificateState) Enum() *CertificateState {
	p := new(CertificateState)
	*p = x
	return p
}

func (x CertificateState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CertificateState) Descriptor() protoreflect.EnumDescriptor {
	return file_gds_models_v1_certificates_proto_enumTypes[1].Descriptor()
}

func (CertificateState) Type() protoreflect.EnumType {
	return &file_gds_models_v1_certificates_proto_enumTypes[1]
}

func (x CertificateState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CertificateState.Descriptor instead.
func (CertificateState) EnumDescriptor() ([]byte, []int) {
	return file_gds_models_v1_certificates_proto_rawDescGZIP(), []int{1}
}

// Certificate requests are maintained separately from the VASP record since they should
// not be replicated. E.g. every directory process is responsible for certificate
// issuance and only public keys and certificate metadata should be exchanged between
// directories.
type CertificateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A unique identifier generated by the directory service, should be a globally
	// unique identifier generated by the replica specified in requesting_replica.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// VASP information for the request
	Vasp       string `protobuf:"bytes,2,opt,name=vasp,proto3" json:"vasp,omitempty"`
	CommonName string `protobuf:"bytes,3,opt,name=common_name,json=commonName,proto3" json:"common_name,omitempty"`
	// Request pipeline status
	Status CertificateRequestState `protobuf:"varint,4,opt,name=status,proto3,enum=gds.models.v1.CertificateRequestState" json:"status,omitempty"`
	// Sectigo create single certificate batch metadata
	AuthorityId  int64  `protobuf:"varint,5,opt,name=authority_id,json=authorityId,proto3" json:"authority_id,omitempty"`
	BatchId      int64  `protobuf:"varint,6,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
	BatchName    string `protobuf:"bytes,7,opt,name=batch_name,json=batchName,proto3" json:"batch_name,omitempty"`
	BatchStatus  string `protobuf:"bytes,8,opt,name=batch_status,json=batchStatus,proto3" json:"batch_status,omitempty"`
	OrderNumber  int64  `protobuf:"varint,9,opt,name=order_number,json=orderNumber,proto3" json:"order_number,omitempty"`
	CreationDate string `protobuf:"bytes,10,opt,name=creation_date,json=creationDate,proto3" json:"creation_date,omitempty"`
	Profile      string `protobuf:"bytes,11,opt,name=profile,proto3" json:"profile,omitempty"`
	RejectReason string `protobuf:"bytes,12,opt,name=reject_reason,json=rejectReason,proto3" json:"reject_reason,omitempty"`
	// Generic parameters used for making requests to certificate issuing services
	Params map[string]string `protobuf:"bytes,13,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Optional list of alternate dns names in addition to the common name
	DnsNames []string `protobuf:"bytes,14,rep,name=dns_names,json=dnsNames,proto3" json:"dns_names,omitempty"`
	// Logging information timestamps
	Created  string `protobuf:"bytes,15,opt,name=created,proto3" json:"created,omitempty"`
	Modified string `protobuf:"bytes,16,opt,name=modified,proto3" json:"modified,omitempty"`
	// Log of historical request states
	AuditLog []*CertificateRequestLogEntry `protobuf:"bytes,17,rep,name=audit_log,json=auditLog,proto3" json:"audit_log,omitempty"`
	// The certificate ID downloaded from the request, if completed successfully
	Certificate string `protobuf:"bytes,18,opt,name=certificate,proto3" json:"certificate,omitempty"`
}

func (x *CertificateRequest) Reset() {
	*x = CertificateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_models_v1_certificates_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CertificateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CertificateRequest) ProtoMessage() {}

func (x *CertificateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gds_models_v1_certificates_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CertificateRequest.ProtoReflect.Descriptor instead.
func (*CertificateRequest) Descriptor() ([]byte, []int) {
	return file_gds_models_v1_certificates_proto_rawDescGZIP(), []int{0}
}

func (x *CertificateRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CertificateRequest) GetVasp() string {
	if x != nil {
		return x.Vasp
	}
	return ""
}

func (x *CertificateRequest) GetCommonName() string {
	if x != nil {
		return x.CommonName
	}
	return ""
}

func (x *CertificateRequest) GetStatus() CertificateRequestState {
	if x != nil {
		return x.Status
	}
	return CertificateRequestState_INITIALIZED
}

func (x *CertificateRequest) GetAuthorityId() int64 {
	if x != nil {
		return x.AuthorityId
	}
	return 0
}

func (x *CertificateRequest) GetBatchId() int64 {
	if x != nil {
		return x.BatchId
	}
	return 0
}

func (x *CertificateRequest) GetBatchName() string {
	if x != nil {
		return x.BatchName
	}
	return ""
}

func (x *CertificateRequest) GetBatchStatus() string {
	if x != nil {
		return x.BatchStatus
	}
	return ""
}

func (x *CertificateRequest) GetOrderNumber() int64 {
	if x != nil {
		return x.OrderNumber
	}
	return 0
}

func (x *CertificateRequest) GetCreationDate() string {
	if x != nil {
		return x.CreationDate
	}
	return ""
}

func (x *CertificateRequest) GetProfile() string {
	if x != nil {
		return x.Profile
	}
	return ""
}

func (x *CertificateRequest) GetRejectReason() string {
	if x != nil {
		return x.RejectReason
	}
	return ""
}

func (x *CertificateRequest) GetParams() map[string]string {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *CertificateRequest) GetDnsNames() []string {
	if x != nil {
		return x.DnsNames
	}
	return nil
}

func (x *CertificateRequest) GetCreated() string {
	if x != nil {
		return x.Created
	}
	return ""
}

func (x *CertificateRequest) GetModified() string {
	if x != nil {
		return x.Modified
	}
	return ""
}

func (x *CertificateRequest) GetAuditLog() []*CertificateRequestLogEntry {
	if x != nil {
		return x.AuditLog
	}
	return nil
}

func (x *CertificateRequest) GetCertificate() string {
	if x != nil {
		return x.Certificate
	}
	return ""
}

// Certificate embeds a TRISA Certificate into a record that can be stored in the
// database for certificate management.
type Certificate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A unique identifier generated by the directory service for storage
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// CertificateRequest that this certificate is associated with
	Request string `protobuf:"bytes,2,opt,name=request,proto3" json:"request,omitempty"`
	// VASP that this certificate belongs to
	Vasp string `protobuf:"bytes,3,opt,name=vasp,proto3" json:"vasp,omitempty"`
	// Current status of the certificate
	Status CertificateState `protobuf:"varint,4,opt,name=status,proto3,enum=gds.models.v1.CertificateState" json:"status,omitempty"`
	// Certificate details
	Details *v1beta1.Certificate `protobuf:"bytes,5,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *Certificate) Reset() {
	*x = Certificate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_models_v1_certificates_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Certificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Certificate) ProtoMessage() {}

func (x *Certificate) ProtoReflect() protoreflect.Message {
	mi := &file_gds_models_v1_certificates_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Certificate.ProtoReflect.Descriptor instead.
func (*Certificate) Descriptor() ([]byte, []int) {
	return file_gds_models_v1_certificates_proto_rawDescGZIP(), []int{1}
}

func (x *Certificate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Certificate) GetRequest() string {
	if x != nil {
		return x.Request
	}
	return ""
}

func (x *Certificate) GetVasp() string {
	if x != nil {
		return x.Vasp
	}
	return ""
}

func (x *Certificate) GetStatus() CertificateState {
	if x != nil {
		return x.Status
	}
	return CertificateState_ISSUED
}

func (x *Certificate) GetDetails() *v1beta1.Certificate {
	if x != nil {
		return x.Details
	}
	return nil
}

// CertificateRequestLogEntry contains information about the state of a certificate request.
type CertificateRequestLogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// RFC3339 timestamp
	Timestamp string `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Previous request state (handled internally) and current request state
	PreviousState CertificateRequestState `protobuf:"varint,2,opt,name=previous_state,json=previousState,proto3,enum=gds.models.v1.CertificateRequestState" json:"previous_state,omitempty"`
	CurrentState  CertificateRequestState `protobuf:"varint,3,opt,name=current_state,json=currentState,proto3,enum=gds.models.v1.CertificateRequestState" json:"current_state,omitempty"`
	// Description of the current state of the certificate request
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// Email address of the Admin who made the state change, "automated" if the state
	// change happened automatically
	Source string `protobuf:"bytes,5,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *CertificateRequestLogEntry) Reset() {
	*x = CertificateRequestLogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gds_models_v1_certificates_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CertificateRequestLogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CertificateRequestLogEntry) ProtoMessage() {}

func (x *CertificateRequestLogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_gds_models_v1_certificates_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CertificateRequestLogEntry.ProtoReflect.Descriptor instead.
func (*CertificateRequestLogEntry) Descriptor() ([]byte, []int) {
	return file_gds_models_v1_certificates_proto_rawDescGZIP(), []int{2}
}

func (x *CertificateRequestLogEntry) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *CertificateRequestLogEntry) GetPreviousState() CertificateRequestState {
	if x != nil {
		return x.PreviousState
	}
	return CertificateRequestState_INITIALIZED
}

func (x *CertificateRequestLogEntry) GetCurrentState() CertificateRequestState {
	if x != nil {
		return x.CurrentState
	}
	return CertificateRequestState_INITIALIZED
}

func (x *CertificateRequestLogEntry) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CertificateRequestLogEntry) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

var File_gds_models_v1_certificates_proto protoreflect.FileDescriptor

var file_gds_models_v1_certificates_proto_rawDesc = []byte{
	0x0a, 0x20, 0x67, 0x64, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0d, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76,
	0x31, 0x1a, 0x21, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f, 0x67, 0x64, 0x73, 0x2f, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x63, 0x61, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f, 0x67, 0x64, 0x73, 0x2f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdf, 0x05, 0x0a, 0x12,
	0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x61, 0x73, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x76, 0x61, 0x73, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x3e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x61, 0x74, 0x63, 0x68,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x62, 0x61, 0x74, 0x63,
	0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x6a,
	0x65, 0x63, 0x74, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x45,
	0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d,
	0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6e, 0x73, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x64, 0x6e, 0x73, 0x4e, 0x61, 0x6d,
	0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0f, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x46, 0x0a, 0x09, 0x61, 0x75, 0x64, 0x69,
	0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x18, 0x11, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x67, 0x64,
	0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x61, 0x75, 0x64, 0x69, 0x74, 0x4c, 0x6f, 0x67,
	0x12, 0x20, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18,
	0x12, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x1a, 0x39, 0x0a, 0x0b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xc5, 0x01,
	0x0a, 0x0b, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x61, 0x73, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x76, 0x61, 0x73, 0x70, 0x12, 0x37, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x67, 0x64,
	0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x3f, 0x0a, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x67, 0x64,
	0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2e, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x07, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x90, 0x02, 0x0a, 0x1a, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x4d, 0x0a, 0x0e, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x67, 0x64, 0x73,
	0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x0d, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x4b, 0x0a, 0x0d, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2a, 0xa0, 0x01, 0x0a, 0x17, 0x43, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x49, 0x54, 0x49, 0x41, 0x4c, 0x49,
	0x5a, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x54,
	0x4f, 0x5f, 0x53, 0x55, 0x42, 0x4d, 0x49, 0x54, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x52,
	0x4f, 0x43, 0x45, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x4f,
	0x57, 0x4e, 0x4c, 0x4f, 0x41, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x44,
	0x4f, 0x57, 0x4e, 0x4c, 0x4f, 0x41, 0x44, 0x45, 0x44, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x43,
	0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x05, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x52,
	0x5f, 0x52, 0x45, 0x4a, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x06, 0x12, 0x0e, 0x0a, 0x0a, 0x43,
	0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x45, 0x44, 0x10, 0x07, 0x2a, 0x38, 0x0a, 0x10, 0x43,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x0a, 0x0a, 0x06, 0x49, 0x53, 0x53, 0x55, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x45,
	0x58, 0x50, 0x49, 0x52, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x56, 0x4f,
	0x4b, 0x45, 0x44, 0x10, 0x02, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x72, 0x69, 0x73, 0x61, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2f,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gds_models_v1_certificates_proto_rawDescOnce sync.Once
	file_gds_models_v1_certificates_proto_rawDescData = file_gds_models_v1_certificates_proto_rawDesc
)

func file_gds_models_v1_certificates_proto_rawDescGZIP() []byte {
	file_gds_models_v1_certificates_proto_rawDescOnce.Do(func() {
		file_gds_models_v1_certificates_proto_rawDescData = protoimpl.X.CompressGZIP(file_gds_models_v1_certificates_proto_rawDescData)
	})
	return file_gds_models_v1_certificates_proto_rawDescData
}

var file_gds_models_v1_certificates_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_gds_models_v1_certificates_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_gds_models_v1_certificates_proto_goTypes = []interface{}{
	(CertificateRequestState)(0),       // 0: gds.models.v1.CertificateRequestState
	(CertificateState)(0),              // 1: gds.models.v1.CertificateState
	(*CertificateRequest)(nil),         // 2: gds.models.v1.CertificateRequest
	(*Certificate)(nil),                // 3: gds.models.v1.Certificate
	(*CertificateRequestLogEntry)(nil), // 4: gds.models.v1.CertificateRequestLogEntry
	nil,                                // 5: gds.models.v1.CertificateRequest.ParamsEntry
	(*v1beta1.Certificate)(nil),        // 6: trisa.gds.models.v1beta1.Certificate
}
var file_gds_models_v1_certificates_proto_depIdxs = []int32{
	0, // 0: gds.models.v1.CertificateRequest.status:type_name -> gds.models.v1.CertificateRequestState
	5, // 1: gds.models.v1.CertificateRequest.params:type_name -> gds.models.v1.CertificateRequest.ParamsEntry
	4, // 2: gds.models.v1.CertificateRequest.audit_log:type_name -> gds.models.v1.CertificateRequestLogEntry
	1, // 3: gds.models.v1.Certificate.status:type_name -> gds.models.v1.CertificateState
	6, // 4: gds.models.v1.Certificate.details:type_name -> trisa.gds.models.v1beta1.Certificate
	0, // 5: gds.models.v1.CertificateRequestLogEntry.previous_state:type_name -> gds.models.v1.CertificateRequestState
	0, // 6: gds.models.v1.CertificateRequestLogEntry.current_state:type_name -> gds.models.v1.CertificateRequestState
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_gds_models_v1_certificates_proto_init() }
func file_gds_models_v1_certificates_proto_init() {
	if File_gds_models_v1_certificates_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gds_models_v1_certificates_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CertificateRequest); i {
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
		file_gds_models_v1_certificates_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Certificate); i {
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
		file_gds_models_v1_certificates_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CertificateRequestLogEntry); i {
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
			RawDescriptor: file_gds_models_v1_certificates_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gds_models_v1_certificates_proto_goTypes,
		DependencyIndexes: file_gds_models_v1_certificates_proto_depIdxs,
		EnumInfos:         file_gds_models_v1_certificates_proto_enumTypes,
		MessageInfos:      file_gds_models_v1_certificates_proto_msgTypes,
	}.Build()
	File_gds_models_v1_certificates_proto = out.File
	file_gds_models_v1_certificates_proto_rawDesc = nil
	file_gds_models_v1_certificates_proto_goTypes = nil
	file_gds_models_v1_certificates_proto_depIdxs = nil
}
