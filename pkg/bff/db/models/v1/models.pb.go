// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: bff/models/v1/models.proto

package models

import (
	ivms101 "github.com/trisacrypto/trisa/pkg/ivms101"
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

// The Organization document contains VASP-specific information for a single VASP record
// in the directory service. This document differs in that it stores information
// relevant to the BFF and should not be used to duplicate storage in the directory.
type Organization struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// BFF Unique Identifier and Record Information
	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	LogoUrl string `protobuf:"bytes,3,opt,name=logo_url,json=logoUrl,proto3" json:"logo_url,omitempty"`
	// Directory Registrations for Lookups
	// TODO: populate these details in the Registration Endpoint
	Testnet *DirectoryRecord `protobuf:"bytes,10,opt,name=testnet,proto3" json:"testnet,omitempty"`
	Mainnet *DirectoryRecord `protobuf:"bytes,11,opt,name=mainnet,proto3" json:"mainnet,omitempty"`
	// Registration Form
	Registration *RegistrationForm `protobuf:"bytes,13,opt,name=registration,proto3" json:"registration,omitempty"`
	// Metadata as RFC3339Nano Timestamps
	Created  string `protobuf:"bytes,14,opt,name=created,proto3" json:"created,omitempty"`
	Modified string `protobuf:"bytes,15,opt,name=modified,proto3" json:"modified,omitempty"`
}

func (x *Organization) Reset() {
	*x = Organization{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bff_models_v1_models_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Organization) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Organization) ProtoMessage() {}

func (x *Organization) ProtoReflect() protoreflect.Message {
	mi := &file_bff_models_v1_models_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Organization.ProtoReflect.Descriptor instead.
func (*Organization) Descriptor() ([]byte, []int) {
	return file_bff_models_v1_models_proto_rawDescGZIP(), []int{0}
}

func (x *Organization) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Organization) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Organization) GetLogoUrl() string {
	if x != nil {
		return x.LogoUrl
	}
	return ""
}

func (x *Organization) GetTestnet() *DirectoryRecord {
	if x != nil {
		return x.Testnet
	}
	return nil
}

func (x *Organization) GetMainnet() *DirectoryRecord {
	if x != nil {
		return x.Mainnet
	}
	return nil
}

func (x *Organization) GetRegistration() *RegistrationForm {
	if x != nil {
		return x.Registration
	}
	return nil
}

func (x *Organization) GetCreated() string {
	if x != nil {
		return x.Created
	}
	return ""
}

func (x *Organization) GetModified() string {
	if x != nil {
		return x.Modified
	}
	return ""
}

// FormState contains the current state of an organization's registration form to
// enable a consistent user experience across multiple contexts.
type FormState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The current 1-indexed step of the form
	Current int32 `protobuf:"varint,1,opt,name=current,proto3" json:"current,omitempty"`
	// If set, the form is completely filled out and ready to be submitted
	ReadyToSubmit bool `protobuf:"varint,2,opt,name=ready_to_submit,json=readyToSubmit,proto3" json:"ready_to_submit,omitempty"`
	// The state of each step in the form
	Steps []*FormStep `protobuf:"bytes,3,rep,name=steps,proto3" json:"steps,omitempty"`
}

func (x *FormState) Reset() {
	*x = FormState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bff_models_v1_models_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FormState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FormState) ProtoMessage() {}

func (x *FormState) ProtoReflect() protoreflect.Message {
	mi := &file_bff_models_v1_models_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FormState.ProtoReflect.Descriptor instead.
func (*FormState) Descriptor() ([]byte, []int) {
	return file_bff_models_v1_models_proto_rawDescGZIP(), []int{1}
}

func (x *FormState) GetCurrent() int32 {
	if x != nil {
		return x.Current
	}
	return 0
}

func (x *FormState) GetReadyToSubmit() bool {
	if x != nil {
		return x.ReadyToSubmit
	}
	return false
}

func (x *FormState) GetSteps() []*FormStep {
	if x != nil {
		return x.Steps
	}
	return nil
}

// FormStep contains the state of a single step in an organization's registration form.
type FormStep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    int32  `protobuf:"varint,1,opt,name=key,proto3" json:"key,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *FormStep) Reset() {
	*x = FormStep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bff_models_v1_models_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FormStep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FormStep) ProtoMessage() {}

func (x *FormStep) ProtoReflect() protoreflect.Message {
	mi := &file_bff_models_v1_models_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FormStep.ProtoReflect.Descriptor instead.
func (*FormStep) Descriptor() ([]byte, []int) {
	return file_bff_models_v1_models_proto_rawDescGZIP(), []int{2}
}

func (x *FormStep) GetKey() int32 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *FormStep) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

// DirectoryRecord contains the information needed to lookup a VASP in a directory service.
type DirectoryRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                  string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RegisteredDirectory string `protobuf:"bytes,2,opt,name=registered_directory,json=registeredDirectory,proto3" json:"registered_directory,omitempty"`
	CommonName          string `protobuf:"bytes,3,opt,name=common_name,json=commonName,proto3" json:"common_name,omitempty"`
	// RFC 3339 timestamp -- if set, the form has been submitted without error
	Submitted string `protobuf:"bytes,15,opt,name=submitted,proto3" json:"submitted,omitempty"`
}

func (x *DirectoryRecord) Reset() {
	*x = DirectoryRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bff_models_v1_models_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirectoryRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirectoryRecord) ProtoMessage() {}

func (x *DirectoryRecord) ProtoReflect() protoreflect.Message {
	mi := &file_bff_models_v1_models_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirectoryRecord.ProtoReflect.Descriptor instead.
func (*DirectoryRecord) Descriptor() ([]byte, []int) {
	return file_bff_models_v1_models_proto_rawDescGZIP(), []int{3}
}

func (x *DirectoryRecord) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *DirectoryRecord) GetRegisteredDirectory() string {
	if x != nil {
		return x.RegisteredDirectory
	}
	return ""
}

func (x *DirectoryRecord) GetCommonName() string {
	if x != nil {
		return x.CommonName
	}
	return ""
}

func (x *DirectoryRecord) GetSubmitted() string {
	if x != nil {
		return x.Submitted
	}
	return ""
}

// RegistrationForm is an extension of the TRISA GDS RegistrationRequest with BFF fields.
type RegistrationForm struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Business information
	Website          string                   `protobuf:"bytes,1,opt,name=website,proto3" json:"website,omitempty"`
	BusinessCategory v1beta1.BusinessCategory `protobuf:"varint,2,opt,name=business_category,json=businessCategory,proto3,enum=trisa.gds.models.v1beta1.BusinessCategory" json:"business_category,omitempty"`
	VaspCategories   []string                 `protobuf:"bytes,3,rep,name=vasp_categories,json=vaspCategories,proto3" json:"vasp_categories,omitempty"`
	EstablishedOn    string                   `protobuf:"bytes,4,opt,name=established_on,json=establishedOn,proto3" json:"established_on,omitempty"`
	// IVMS 101 Legal Person record
	Entity *ivms101.LegalPerson `protobuf:"bytes,11,opt,name=entity,proto3" json:"entity,omitempty"`
	// Directory Record contacts
	Contacts *v1beta1.Contacts `protobuf:"bytes,12,opt,name=contacts,proto3" json:"contacts,omitempty"`
	// TRIXO Form
	Trixo *v1beta1.TRIXOQuestionnaire `protobuf:"bytes,13,opt,name=trixo,proto3" json:"trixo,omitempty"`
	// Network-specific information and submission details
	Testnet *NetworkDetails `protobuf:"bytes,14,opt,name=testnet,proto3" json:"testnet,omitempty"`
	Mainnet *NetworkDetails `protobuf:"bytes,15,opt,name=mainnet,proto3" json:"mainnet,omitempty"`
	// Current progress of the form for the frontend
	State *FormState `protobuf:"bytes,20,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *RegistrationForm) Reset() {
	*x = RegistrationForm{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bff_models_v1_models_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationForm) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationForm) ProtoMessage() {}

func (x *RegistrationForm) ProtoReflect() protoreflect.Message {
	mi := &file_bff_models_v1_models_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationForm.ProtoReflect.Descriptor instead.
func (*RegistrationForm) Descriptor() ([]byte, []int) {
	return file_bff_models_v1_models_proto_rawDescGZIP(), []int{4}
}

func (x *RegistrationForm) GetWebsite() string {
	if x != nil {
		return x.Website
	}
	return ""
}

func (x *RegistrationForm) GetBusinessCategory() v1beta1.BusinessCategory {
	if x != nil {
		return x.BusinessCategory
	}
	return v1beta1.BusinessCategory(0)
}

func (x *RegistrationForm) GetVaspCategories() []string {
	if x != nil {
		return x.VaspCategories
	}
	return nil
}

func (x *RegistrationForm) GetEstablishedOn() string {
	if x != nil {
		return x.EstablishedOn
	}
	return ""
}

func (x *RegistrationForm) GetEntity() *ivms101.LegalPerson {
	if x != nil {
		return x.Entity
	}
	return nil
}

func (x *RegistrationForm) GetContacts() *v1beta1.Contacts {
	if x != nil {
		return x.Contacts
	}
	return nil
}

func (x *RegistrationForm) GetTrixo() *v1beta1.TRIXOQuestionnaire {
	if x != nil {
		return x.Trixo
	}
	return nil
}

func (x *RegistrationForm) GetTestnet() *NetworkDetails {
	if x != nil {
		return x.Testnet
	}
	return nil
}

func (x *RegistrationForm) GetMainnet() *NetworkDetails {
	if x != nil {
		return x.Mainnet
	}
	return nil
}

func (x *RegistrationForm) GetState() *FormState {
	if x != nil {
		return x.State
	}
	return nil
}

// NetworkDetails contains directory-service specific submission information such as the
// certificate request and information about when the registration form was submitted.
type NetworkDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Certificate request information
	CommonName string   `protobuf:"bytes,1,opt,name=common_name,json=commonName,proto3" json:"common_name,omitempty"`
	Endpoint   string   `protobuf:"bytes,2,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	DnsNames   []string `protobuf:"bytes,3,rep,name=dns_names,json=dnsNames,proto3" json:"dns_names,omitempty"`
}

func (x *NetworkDetails) Reset() {
	*x = NetworkDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bff_models_v1_models_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkDetails) ProtoMessage() {}

func (x *NetworkDetails) ProtoReflect() protoreflect.Message {
	mi := &file_bff_models_v1_models_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkDetails.ProtoReflect.Descriptor instead.
func (*NetworkDetails) Descriptor() ([]byte, []int) {
	return file_bff_models_v1_models_proto_rawDescGZIP(), []int{5}
}

func (x *NetworkDetails) GetCommonName() string {
	if x != nil {
		return x.CommonName
	}
	return ""
}

func (x *NetworkDetails) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *NetworkDetails) GetDnsNames() []string {
	if x != nil {
		return x.DnsNames
	}
	return nil
}

// Announcements are made by network administrators to inform all TRISA members of
// important events, maintenance, or milestones. These are broadcast from the BFF so
// that all members receive the same announcement.
type Announcement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title    string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Body     string `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	PostDate string `protobuf:"bytes,4,opt,name=post_date,json=postDate,proto3" json:"post_date,omitempty"`
	Author   string `protobuf:"bytes,5,opt,name=author,proto3" json:"author,omitempty"`
	// Metadata as RFC3339Nano Timestamps
	Created  string `protobuf:"bytes,14,opt,name=created,proto3" json:"created,omitempty"`
	Modified string `protobuf:"bytes,15,opt,name=modified,proto3" json:"modified,omitempty"`
}

func (x *Announcement) Reset() {
	*x = Announcement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bff_models_v1_models_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Announcement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Announcement) ProtoMessage() {}

func (x *Announcement) ProtoReflect() protoreflect.Message {
	mi := &file_bff_models_v1_models_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Announcement.ProtoReflect.Descriptor instead.
func (*Announcement) Descriptor() ([]byte, []int) {
	return file_bff_models_v1_models_proto_rawDescGZIP(), []int{6}
}

func (x *Announcement) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Announcement) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Announcement) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *Announcement) GetPostDate() string {
	if x != nil {
		return x.PostDate
	}
	return ""
}

func (x *Announcement) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

func (x *Announcement) GetCreated() string {
	if x != nil {
		return x.Created
	}
	return ""
}

func (x *Announcement) GetModified() string {
	if x != nil {
		return x.Modified
	}
	return ""
}

// Announcements are stored in months to enable fast retrieval of the latest
// announcements in a specific time range without a reversal traversal of time-ordered
// anncouncement objects. Note that the annoucements are stored in a slice instead of
// a map to reduce data storage overhead. Accessing a specific announcement requires
// iterating over the annoucements, but the number of annoucements in a month should not
// be unbounded, so this cost is acceptable for data storage performance.
type AnnouncementMonth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Date          string          `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	Announcements []*Announcement `protobuf:"bytes,2,rep,name=announcements,proto3" json:"announcements,omitempty"`
	// Metadata as RFC3339Nano Timestamps
	Created  string `protobuf:"bytes,14,opt,name=created,proto3" json:"created,omitempty"`
	Modified string `protobuf:"bytes,15,opt,name=modified,proto3" json:"modified,omitempty"`
}

func (x *AnnouncementMonth) Reset() {
	*x = AnnouncementMonth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bff_models_v1_models_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnnouncementMonth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnnouncementMonth) ProtoMessage() {}

func (x *AnnouncementMonth) ProtoReflect() protoreflect.Message {
	mi := &file_bff_models_v1_models_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnnouncementMonth.ProtoReflect.Descriptor instead.
func (*AnnouncementMonth) Descriptor() ([]byte, []int) {
	return file_bff_models_v1_models_proto_rawDescGZIP(), []int{7}
}

func (x *AnnouncementMonth) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *AnnouncementMonth) GetAnnouncements() []*Announcement {
	if x != nil {
		return x.Announcements
	}
	return nil
}

func (x *AnnouncementMonth) GetCreated() string {
	if x != nil {
		return x.Created
	}
	return ""
}

func (x *AnnouncementMonth) GetModified() string {
	if x != nil {
		return x.Modified
	}
	return ""
}

var File_bff_models_v1_models_proto protoreflect.FileDescriptor

var file_bff_models_v1_models_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x62, 0x66, 0x66, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x76, 0x31, 0x2f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x62, 0x66,
	0x66, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x15, 0x69, 0x76, 0x6d,
	0x73, 0x31, 0x30, 0x31, 0x2f, 0x69, 0x76, 0x6d, 0x73, 0x31, 0x30, 0x31, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x25, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f, 0x67, 0x64, 0x73, 0x2f, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc2, 0x02, 0x0a, 0x0c, 0x4f, 0x72,
	0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19,
	0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x6f, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6c, 0x6f, 0x67, 0x6f, 0x55, 0x72, 0x6c, 0x12, 0x38, 0x0a, 0x07, 0x74, 0x65, 0x73,
	0x74, 0x6e, 0x65, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x62, 0x66, 0x66,
	0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x74, 0x65, 0x73, 0x74,
	0x6e, 0x65, 0x74, 0x12, 0x38, 0x0a, 0x07, 0x6d, 0x61, 0x69, 0x6e, 0x6e, 0x65, 0x74, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x62, 0x66, 0x66, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x6d, 0x61, 0x69, 0x6e, 0x6e, 0x65, 0x74, 0x12, 0x43, 0x0a,
	0x0c, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x62, 0x66, 0x66, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x6f, 0x72, 0x6d, 0x52, 0x0c, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x4a, 0x04, 0x08, 0x0c, 0x10, 0x0d, 0x22, 0x7c,
	0x0a, 0x09, 0x46, 0x6f, 0x72, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x72, 0x65, 0x61, 0x64, 0x79, 0x5f, 0x74,
	0x6f, 0x5f, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d,
	0x72, 0x65, 0x61, 0x64, 0x79, 0x54, 0x6f, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x12, 0x2d, 0x0a,
	0x05, 0x73, 0x74, 0x65, 0x70, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x62,
	0x66, 0x66, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6f, 0x72,
	0x6d, 0x53, 0x74, 0x65, 0x70, 0x52, 0x05, 0x73, 0x74, 0x65, 0x70, 0x73, 0x22, 0x34, 0x0a, 0x08,
	0x46, 0x6f, 0x72, 0x6d, 0x53, 0x74, 0x65, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x93, 0x01, 0x0a, 0x0f, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x31, 0x0a, 0x14, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x65, 0x64, 0x5f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64,
	0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x75,
	0x62, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73,
	0x75, 0x62, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x22, 0xa9, 0x04, 0x0a, 0x10, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x6d, 0x12, 0x18, 0x0a,
	0x07, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x12, 0x57, 0x0a, 0x11, 0x62, 0x75, 0x73, 0x69, 0x6e,
	0x65, 0x73, 0x73, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x2a, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x42, 0x75,
	0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x10,
	0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73, 0x73, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79,
	0x12, 0x27, 0x0a, 0x0f, 0x76, 0x61, 0x73, 0x70, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x69, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x76, 0x61, 0x73, 0x70, 0x43,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x73, 0x74,
	0x61, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x5f, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x65, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x64, 0x4f, 0x6e,
	0x12, 0x2c, 0x0a, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x69, 0x76, 0x6d, 0x73, 0x31, 0x30, 0x31, 0x2e, 0x4c, 0x65, 0x67, 0x61, 0x6c,
	0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x3e,
	0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x22, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x63, 0x74, 0x73, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x12, 0x42,
	0x0a, 0x05, 0x74, 0x72, 0x69, 0x78, 0x6f, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e,
	0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x54, 0x52, 0x49, 0x58, 0x4f, 0x51, 0x75,
	0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x6e, 0x61, 0x69, 0x72, 0x65, 0x52, 0x05, 0x74, 0x72, 0x69,
	0x78, 0x6f, 0x12, 0x37, 0x0a, 0x07, 0x74, 0x65, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x62, 0x66, 0x66, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x44, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x73, 0x52, 0x07, 0x74, 0x65, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x12, 0x37, 0x0a, 0x07, 0x6d,
	0x61, 0x69, 0x6e, 0x6e, 0x65, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x62,
	0x66, 0x66, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x07, 0x6d, 0x61, 0x69,
	0x6e, 0x6e, 0x65, 0x74, 0x12, 0x2e, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x14, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x62, 0x66, 0x66, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6f, 0x72, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x22, 0x6a, 0x0a, 0x0e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6e, 0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x64, 0x6e, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x73,
	0x22, 0xb3, 0x01, 0x0a, 0x0c, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x70,
	0x6f, 0x73, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x70, 0x6f, 0x73, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x6f,
	0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x6f,
	0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x22, 0xa0, 0x01, 0x0a, 0x11, 0x41, 0x6e, 0x6e, 0x6f, 0x75,
	0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65,
	0x12, 0x41, 0x0a, 0x0d, 0x61, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x62, 0x66, 0x66, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0d, 0x61, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0e,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x72, 0x69, 0x73, 0x61, 0x63, 0x72, 0x79,
	0x70, 0x74, 0x6f, 0x2f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x62, 0x66, 0x66, 0x2f, 0x64, 0x62, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f,
	0x76, 0x31, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_bff_models_v1_models_proto_rawDescOnce sync.Once
	file_bff_models_v1_models_proto_rawDescData = file_bff_models_v1_models_proto_rawDesc
)

func file_bff_models_v1_models_proto_rawDescGZIP() []byte {
	file_bff_models_v1_models_proto_rawDescOnce.Do(func() {
		file_bff_models_v1_models_proto_rawDescData = protoimpl.X.CompressGZIP(file_bff_models_v1_models_proto_rawDescData)
	})
	return file_bff_models_v1_models_proto_rawDescData
}

var file_bff_models_v1_models_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_bff_models_v1_models_proto_goTypes = []interface{}{
	(*Organization)(nil),               // 0: bff.models.v1.Organization
	(*FormState)(nil),                  // 1: bff.models.v1.FormState
	(*FormStep)(nil),                   // 2: bff.models.v1.FormStep
	(*DirectoryRecord)(nil),            // 3: bff.models.v1.DirectoryRecord
	(*RegistrationForm)(nil),           // 4: bff.models.v1.RegistrationForm
	(*NetworkDetails)(nil),             // 5: bff.models.v1.NetworkDetails
	(*Announcement)(nil),               // 6: bff.models.v1.Announcement
	(*AnnouncementMonth)(nil),          // 7: bff.models.v1.AnnouncementMonth
	(v1beta1.BusinessCategory)(0),      // 8: trisa.gds.models.v1beta1.BusinessCategory
	(*ivms101.LegalPerson)(nil),        // 9: ivms101.LegalPerson
	(*v1beta1.Contacts)(nil),           // 10: trisa.gds.models.v1beta1.Contacts
	(*v1beta1.TRIXOQuestionnaire)(nil), // 11: trisa.gds.models.v1beta1.TRIXOQuestionnaire
}
var file_bff_models_v1_models_proto_depIdxs = []int32{
	3,  // 0: bff.models.v1.Organization.testnet:type_name -> bff.models.v1.DirectoryRecord
	3,  // 1: bff.models.v1.Organization.mainnet:type_name -> bff.models.v1.DirectoryRecord
	4,  // 2: bff.models.v1.Organization.registration:type_name -> bff.models.v1.RegistrationForm
	2,  // 3: bff.models.v1.FormState.steps:type_name -> bff.models.v1.FormStep
	8,  // 4: bff.models.v1.RegistrationForm.business_category:type_name -> trisa.gds.models.v1beta1.BusinessCategory
	9,  // 5: bff.models.v1.RegistrationForm.entity:type_name -> ivms101.LegalPerson
	10, // 6: bff.models.v1.RegistrationForm.contacts:type_name -> trisa.gds.models.v1beta1.Contacts
	11, // 7: bff.models.v1.RegistrationForm.trixo:type_name -> trisa.gds.models.v1beta1.TRIXOQuestionnaire
	5,  // 8: bff.models.v1.RegistrationForm.testnet:type_name -> bff.models.v1.NetworkDetails
	5,  // 9: bff.models.v1.RegistrationForm.mainnet:type_name -> bff.models.v1.NetworkDetails
	1,  // 10: bff.models.v1.RegistrationForm.state:type_name -> bff.models.v1.FormState
	6,  // 11: bff.models.v1.AnnouncementMonth.announcements:type_name -> bff.models.v1.Announcement
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_bff_models_v1_models_proto_init() }
func file_bff_models_v1_models_proto_init() {
	if File_bff_models_v1_models_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bff_models_v1_models_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Organization); i {
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
		file_bff_models_v1_models_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FormState); i {
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
		file_bff_models_v1_models_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FormStep); i {
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
		file_bff_models_v1_models_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DirectoryRecord); i {
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
		file_bff_models_v1_models_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistrationForm); i {
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
		file_bff_models_v1_models_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkDetails); i {
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
		file_bff_models_v1_models_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Announcement); i {
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
		file_bff_models_v1_models_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnnouncementMonth); i {
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
			RawDescriptor: file_bff_models_v1_models_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bff_models_v1_models_proto_goTypes,
		DependencyIndexes: file_bff_models_v1_models_proto_depIdxs,
		MessageInfos:      file_bff_models_v1_models_proto_msgTypes,
	}.Build()
	File_bff_models_v1_models_proto = out.File
	file_bff_models_v1_models_proto_rawDesc = nil
	file_bff_models_v1_models_proto_goTypes = nil
	file_bff_models_v1_models_proto_depIdxs = nil
}
