package mockdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	bff "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/store/iterator"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

var db store.Store = &MockDB{}
var state = &MockState{
	VASPs: make(map[string]*pb.VASP),
	Keys:  []string{},
}

// MockState contains the current state of the MockDB for test verification.
type MockState struct {
	// in-memory database store
	VASPs map[string]*pb.VASP
	Keys  []string

	// keep track of store interface calls
	CloseInvoked                     bool
	CreateVASPInvoked                bool
	RetrieveVASPInvoked              bool
	UpdateVASPInvoked                bool
	DeleteVASPInvoked                bool
	ListVASPsInvoked                 bool
	SearchVASPsInvoked               bool
	CountVASPsInvoked                bool
	ListCertReqsInvoked              bool
	CreateCertReqInvoked             bool
	RetrieveCertReqInvoked           bool
	UpdateCertReqInvoked             bool
	DeleteCertReqInvoked             bool
	CountCertReqsInvoked             bool
	ListCertInvoked                  bool
	CreateCertInvoked                bool
	RetrieveCertInvoked              bool
	UpdateCertInvoked                bool
	DeleteCertInvoked                bool
	CountCertsInvoked                bool
	RetrieveAnnouncementMonthInvoked bool
	UpdateAnnouncementMonthInvoked   bool
	DeleteAnnouncementMonthInvoked   bool
	CountAnnouncementMonthsInvoked   bool
	RetrieveActivityMonthInvoked     bool
	UpdateActivityMonthInvoked       bool
	DeleteActivityMonthInvoked       bool
	CountActivityMonthsInvoked       bool
	ListOrganizationsInvoked         bool
	CreateOrganizationInvoked        bool
	RetrieveOrganizationInvoked      bool
	UpdateOrganizationInvoked        bool
	DeleteOrganizationInvoked        bool
	CountOrganizationsInvoked        bool
	ListContactsInvoked              bool
	CreateContactInvoked             bool
	RetrieveContactInvoked           bool
	UpdateContactInvoked             bool
	DeleteContactInvoked             bool
	CountContactsInvoked             bool
	ReindexInvoked                   bool
	BackupInvoked                    bool
}

func GetState() *MockState {
	return state
}

func ResetState() {
	state = &MockState{
		VASPs: make(map[string]*pb.VASP),
		Keys:  []string{},
	}
}

// MockDB fulfills the store interface for testing.
type MockDB struct {
	OnClose                     func() error
	OnCreateVASP                func(v *pb.VASP) (string, error)
	OnRetrieveVASP              func(id string) (*pb.VASP, error)
	OnUpdateVASP                func(v *pb.VASP) error
	OnDeleteVASP                func(id string) error
	OnListVASPs                 func() iterator.DirectoryIterator
	OnSearchVASPs               func(query map[string]interface{}) ([]*pb.VASP, error)
	OnCountVASPs                func(context.Context) (uint64, error)
	OnListCertReqs              func() iterator.CertificateRequestIterator
	OnCreateCertReq             func(r *models.CertificateRequest) (string, error)
	OnRetrieveCertReq           func(id string) (*models.CertificateRequest, error)
	OnUpdateCertReq             func(r *models.CertificateRequest) error
	OnDeleteCertReq             func(id string) error
	OnCountCertReqs             func(context.Context) (uint64, error)
	OnListCerts                 func() iterator.CertificateIterator
	OnCreateCert                func(c *models.Certificate) (string, error)
	OnRetrieveCert              func(id string) (*models.Certificate, error)
	OnUpdateCert                func(c *models.Certificate) error
	OnDeleteCert                func(id string) error
	OnCountCerts                func(context.Context) (uint64, error)
	OnRetrieveAnnouncementMonth func(date string) (*bff.AnnouncementMonth, error)
	OnUpdateAnnouncementMonth   func(o *bff.AnnouncementMonth) error
	OnDeleteAnnouncementMonth   func(date string) error
	OnCountAnnouncementMonths   func(context.Context) (uint64, error)
	OnRetrieveActivityMonth     func(date string) (*bff.ActivityMonth, error)
	OnUpdateActivityMonth       func(o *bff.ActivityMonth) error
	OnDeleteActivityMonth       func(date string) error
	OnCountActivityMonths       func(context.Context) (uint64, error)
	OnListOrganizations         func() iterator.OrganizationIterator
	OnCreateOrganization        func(o *bff.Organization) (string, error)
	OnRetrieveOrganization      func(id uuid.UUID) (*bff.Organization, error)
	OnUpdateOrganization        func(o *bff.Organization) error
	OnDeleteOrganization        func(id uuid.UUID) error
	OnCountOrganizations        func(context.Context) (uint64, error)
	OnListContacts              func() []*models.Contact
	OnCreateContact             func(c *models.Contact) (string, error)
	OnRetrieveContact           func(email string) (*models.Contact, error)
	OnUpdateContact             func(c *models.Contact) error
	OnDeleteContact             func(email string) error
	OnCountContacts             func(context.Context) (uint64, error)
	OnReindex                   func() error
	OnBackup                    func(string) error
}

func GetStore() store.Store {
	return db
}

func (m *MockDB) Close() error {
	state.CloseInvoked = true
	return m.OnClose()
}

func (m *MockDB) CreateVASP(_ context.Context, v *pb.VASP) (string, error) {
	state.CreateVASPInvoked = true
	if v.Id == "" {
		return "", errors.New("VASP must contain an ID")
	}
	if _, ok := state.VASPs[v.Id]; ok {
		return "", fmt.Errorf("VASP with ID %s already exists", v.Id)
	}
	state.VASPs[v.Id] = v
	state.Keys = append(state.Keys, v.Id)
	return v.Id, nil
}

func (m *MockDB) RetrieveVASP(_ context.Context, id string) (*pb.VASP, error) {
	state.RetrieveVASPInvoked = true
	if id == "" {
		return nil, errors.New("missing VASP ID")
	}

	v, ok := state.VASPs[id]
	if !ok {
		return nil, fmt.Errorf("VASP with ID %s not found", id)
	}
	return v, nil
}

func (m *MockDB) UpdateVASP(_ context.Context, v *pb.VASP) error {
	state.UpdateVASPInvoked = true
	return m.OnUpdateVASP(v)
}

func (m *MockDB) DeleteVASP(_ context.Context, id string) error {
	state.DeleteVASPInvoked = true
	return m.OnDeleteVASP(id)
}

func (m *MockDB) ListVASPs(_ context.Context) iterator.DirectoryIterator {
	state.ListVASPsInvoked = true
	return m.OnListVASPs()
}

func (m *MockDB) SearchVASPs(_ context.Context, query map[string]interface{}) ([]*pb.VASP, error) {
	state.SearchVASPsInvoked = true
	return m.OnSearchVASPs(query)
}

func (m *MockDB) CountVASPs(ctx context.Context) (uint64, error) {
	state.CountVASPsInvoked = true
	return m.OnCountVASPs(ctx)
}

func (m *MockDB) ListCertReqs(_ context.Context) iterator.CertificateRequestIterator {
	state.ListCertReqsInvoked = true
	return m.OnListCertReqs()
}

func (m *MockDB) CreateCertReq(_ context.Context, r *models.CertificateRequest) (string, error) {
	state.CreateCertReqInvoked = true
	return m.OnCreateCertReq(r)
}

func (m *MockDB) RetrieveCertReq(_ context.Context, id string) (*models.CertificateRequest, error) {
	state.RetrieveCertReqInvoked = true
	return m.OnRetrieveCertReq(id)
}

func (m *MockDB) UpdateCertReq(_ context.Context, r *models.CertificateRequest) error {
	state.UpdateCertReqInvoked = true
	return m.OnUpdateCertReq(r)
}

func (m *MockDB) DeleteCertReq(_ context.Context, id string) error {
	state.DeleteCertReqInvoked = true
	return m.OnDeleteCertReq(id)
}

func (m *MockDB) CountCertReqs(ctx context.Context) (uint64, error) {
	state.CountCertReqsInvoked = true
	return m.OnCountCertReqs(ctx)
}

func (m *MockDB) ListCerts(_ context.Context) iterator.CertificateIterator {
	state.ListCertInvoked = true
	return m.OnListCerts()
}

func (m *MockDB) CreateCert(_ context.Context, c *models.Certificate) (string, error) {
	state.CreateCertInvoked = true
	return m.OnCreateCert(c)
}

func (m *MockDB) RetrieveCert(_ context.Context, id string) (*models.Certificate, error) {
	state.RetrieveCertInvoked = true
	return m.OnRetrieveCert(id)
}

func (m *MockDB) UpdateCert(_ context.Context, c *models.Certificate) error {
	state.UpdateCertInvoked = true
	return m.OnUpdateCert(c)
}

func (m *MockDB) DeleteCert(_ context.Context, id string) error {
	state.DeleteCertInvoked = true
	return m.OnDeleteCert(id)
}

func (m *MockDB) CountCerts(ctx context.Context) (uint64, error) {
	state.CountCertsInvoked = true
	return m.OnCountCerts(ctx)
}

func (m *MockDB) RetrieveAnnouncementMonth(_ context.Context, date string) (*bff.AnnouncementMonth, error) {
	state.RetrieveAnnouncementMonthInvoked = true
	return m.OnRetrieveAnnouncementMonth(date)
}

func (m *MockDB) UpdateAnnouncementMonth(_ context.Context, o *bff.AnnouncementMonth) error {
	state.UpdateAnnouncementMonthInvoked = true
	return m.OnUpdateAnnouncementMonth(o)
}

func (m *MockDB) DeleteAnnouncementMonth(_ context.Context, date string) error {
	state.DeleteAnnouncementMonthInvoked = true
	return m.OnDeleteAnnouncementMonth(date)
}

func (m *MockDB) CountAnnouncementMonths(ctx context.Context) (uint64, error) {
	state.CountAnnouncementMonthsInvoked = true
	return m.OnCountAnnouncementMonths(ctx)
}

func (m *MockDB) RetrieveActivityMonth(_ context.Context, date string) (*bff.ActivityMonth, error) {
	state.RetrieveActivityMonthInvoked = true
	return m.OnRetrieveActivityMonth(date)
}

func (m *MockDB) UpdateActivityMonth(_ context.Context, o *bff.ActivityMonth) error {
	state.UpdateActivityMonthInvoked = true
	return m.OnUpdateActivityMonth(o)
}

func (m *MockDB) DeleteActivityMonth(_ context.Context, date string) error {
	state.DeleteActivityMonthInvoked = true
	return m.OnDeleteActivityMonth(date)
}

func (m *MockDB) CountActivityMonth(ctx context.Context) (uint64, error) {
	state.CountActivityMonthsInvoked = true
	return m.OnCountActivityMonths(ctx)
}

func (m *MockDB) ListOrganizations(_ context.Context) iterator.OrganizationIterator {
	state.ListOrganizationsInvoked = true
	return m.OnListOrganizations()
}

func (m *MockDB) CreateOrganization(_ context.Context, o *bff.Organization) (string, error) {
	state.CreateOrganizationInvoked = true
	return m.OnCreateOrganization(o)
}

func (m *MockDB) RetrieveOrganization(_ context.Context, id uuid.UUID) (*bff.Organization, error) {
	state.RetrieveOrganizationInvoked = true
	return m.OnRetrieveOrganization(id)
}

func (m *MockDB) UpdateOrganization(_ context.Context, o *bff.Organization) error {
	state.UpdateOrganizationInvoked = true
	return m.OnUpdateOrganization(o)
}

func (m *MockDB) DeleteOrganization(_ context.Context, id uuid.UUID) error {
	state.DeleteOrganizationInvoked = true
	return m.OnDeleteOrganization(id)
}

func (m *MockDB) CountOrganizations(ctx context.Context) (uint64, error) {
	state.CountOrganizationsInvoked = true
	return m.OnCountOrganizations(ctx)
}

func (m *MockDB) ListContacts(_ context.Context) []*models.Contact {
	state.ListContactsInvoked = true
	return m.OnListContacts()
}

func (m *MockDB) CreateContact(_ context.Context, c *models.Contact) (string, error) {
	state.CreateContactInvoked = true
	return m.OnCreateContact(c)
}

func (m *MockDB) RetrieveContact(_ context.Context, email string) (*models.Contact, error) {
	state.RetrieveCertInvoked = true
	return m.OnRetrieveContact(email)
}

func (m *MockDB) UpdateContact(_ context.Context, c *models.Contact) error {
	state.UpdateCertInvoked = true
	return m.OnUpdateContact(c)
}

func (m *MockDB) DeleteContact(_ context.Context, email string) error {
	state.DeleteContactInvoked = true
	return m.OnDeleteContact(email)
}

func (m *MockDB) CountContacts(ctx context.Context) (uint64, error) {
	state.CountContactsInvoked = true
	return m.OnCountContacts(ctx)
}

func (m *MockDB) Reindex() error {
	state.ReindexInvoked = true
	return m.OnReindex()
}

func (m *MockDB) Backup(path string) error {
	state.BackupInvoked = true
	return m.OnBackup(path)
}
