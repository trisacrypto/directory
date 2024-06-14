package mock

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	bff "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/store/iterator"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

var (
	ErrNoMock = errors.New("no mock handler implemented")
)

const (
	Close                     = "Close"
	Reindex                   = "Reindex"
	Backup                    = "Backup"
	ListVASPs                 = "ListVASPs"
	SearchVASPs               = "SearchVASPs"
	CreateVASP                = "CreateVASP"
	RetrieveVASP              = "RetrieveVASP"
	UpdateVASP                = "UpdateVASP"
	DeleteVASP                = "DeleteVASP"
	CountVASPs                = "CountVASPs"
	ListCertReqs              = "ListCertReqs"
	CreateCertReq             = "CreateCertReq"
	RetrieveCertReq           = "RetrieveCertReq"
	UpdateCertReq             = "UpdateCertReq"
	DeleteCertReq             = "DeleteCertReq"
	CountCertReqs             = "CountCertReqs"
	ListCerts                 = "ListCerts"
	CreateCert                = "CreateCert"
	RetrieveCert              = "RetrieveCert"
	UpdateCert                = "UpdateCert"
	DeleteCert                = "DeleteCert"
	CountCerts                = "CountCerts"
	RetrieveAnnouncementMonth = "RetrieveAnnouncementMonth"
	UpdateAnnouncementMonth   = "UpdateAnnouncementMonth"
	DeleteAnnouncementMonth   = "DeleteAnnouncementMonth"
	CountAnnouncementMonths   = "CountAnnouncementMonths"
	RetrieveActivityMonth     = "RetrieveActivityMonth"
	UpdateActivityMonth       = "UpdateActivityMonth"
	DeleteActivityMonth       = "DeleteActivityMonth"
	CountActivityMonth        = "CountActivityMonth"
	ListOrganizations         = "ListOrganizations"
	CreateOrganization        = "CreateOrganization"
	RetrieveOrganization      = "RetrieveOrganization"
	UpdateOrganization        = "UpdateOrganization"
	DeleteOrganization        = "DeleteOrganization"
	CountOrganizations        = "CountOrganizations"
	ListContacts              = "ListContacts"
	CreateContact             = "CreateContact"
	RetrieveContact           = "RetrieveContact"
	UpdateContact             = "UpdateContact"
	DeleteContact             = "DeleteContact"
	CountContacts             = "CountContacts"
	ListEmails                = "ListEmails"
	CreateEmail               = "CreateEmail"
	RetrieveEmail             = "RetrieveEmail"
	UpdateEmail               = "UpdateEmail"
	DeleteEmail               = "DeleteEmail"
	CountEmails               = "CountEmails"
	VASPContacts              = "VASPContacts"
	RetrieveVASPContacts      = "RetrieveVASPContacts"
	UpdateVASPContacts        = "UpdateVASPContacts"
)

// Store implements the store.Store interface for testing purposes.
type Store struct {
	sync.RWMutex
	calls                       map[string]int
	OnClose                     func() error
	OnReindex                   func() error
	OnBackup                    func(string) error
	OnListVASPs                 func(ctx context.Context) iterator.DirectoryIterator
	OnSearchVASPs               func(ctx context.Context, query map[string]interface{}) ([]*pb.VASP, error)
	OnCreateVASP                func(ctx context.Context, v *pb.VASP) (string, error)
	OnRetrieveVASP              func(ctx context.Context, id string) (*pb.VASP, error)
	OnUpdateVASP                func(ctx context.Context, v *pb.VASP) error
	OnDeleteVASP                func(ctx context.Context, id string) error
	OnCountVASPs                func(ctx context.Context) (uint64, error)
	OnListCertReqs              func(ctx context.Context) iterator.CertificateRequestIterator
	OnCreateCertReq             func(ctx context.Context, r *models.CertificateRequest) (string, error)
	OnRetrieveCertReq           func(ctx context.Context, id string) (*models.CertificateRequest, error)
	OnUpdateCertReq             func(ctx context.Context, r *models.CertificateRequest) error
	OnDeleteCertReq             func(ctx context.Context, id string) error
	OnCountCertReqs             func(context.Context) (uint64, error)
	OnListCerts                 func(ctx context.Context) iterator.CertificateIterator
	OnCreateCert                func(ctx context.Context, c *models.Certificate) (string, error)
	OnRetrieveCert              func(ctx context.Context, id string) (*models.Certificate, error)
	OnUpdateCert                func(ctx context.Context, c *models.Certificate) error
	OnDeleteCert                func(ctx context.Context, id string) error
	OnCountCerts                func(context.Context) (uint64, error)
	OnRetrieveAnnouncementMonth func(ctx context.Context, date string) (*bff.AnnouncementMonth, error)
	OnUpdateAnnouncementMonth   func(ctx context.Context, m *bff.AnnouncementMonth) error
	OnDeleteAnnouncementMonth   func(ctx context.Context, date string) error
	OnCountAnnouncementMonths   func(context.Context) (uint64, error)
	OnRetrieveActivityMonth     func(ctx context.Context, date string) (*bff.ActivityMonth, error)
	OnUpdateActivityMonth       func(ctx context.Context, m *bff.ActivityMonth) error
	OnDeleteActivityMonth       func(ctx context.Context, date string) error
	OnCountActivityMonth        func(context.Context) (uint64, error)
	OnListOrganizations         func(ctx context.Context) iterator.OrganizationIterator
	OnCreateOrganization        func(ctx context.Context, o *bff.Organization) (string, error)
	OnRetrieveOrganization      func(ctx context.Context, id uuid.UUID) (*bff.Organization, error)
	OnUpdateOrganization        func(ctx context.Context, o *bff.Organization) error
	OnDeleteOrganization        func(ctx context.Context, id uuid.UUID) error
	OnCountOrganizations        func(context.Context) (uint64, error)
	OnListContacts              func(ctx context.Context) []*models.Contact
	OnCreateContact             func(ctx context.Context, c *models.Contact) (string, error)
	OnRetrieveContact           func(ctx context.Context, email string) (*models.Contact, error)
	OnUpdateContact             func(ctx context.Context, c *models.Contact) error
	OnDeleteContact             func(ctx context.Context, email string) error
	OnCountContacts             func(context.Context) (uint64, error)
	OnListEmails                func(ctx context.Context) iterator.EmailIterator
	OnCreateEmail               func(ctx context.Context, c *models.Email) (string, error)
	OnRetrieveEmail             func(ctx context.Context, email string) (*models.Email, error)
	OnUpdateEmail               func(ctx context.Context, c *models.Email) error
	OnDeleteEmail               func(ctx context.Context, email string) error
	OnCountEmails               func(context.Context) (uint64, error)
	OnVASPContacts              func(ctx context.Context, vasp *pb.VASP) (*models.Contacts, error)
	OnRetrieveVASPContacts      func(ctx context.Context, vaspID string) (*models.Contacts, error)
	OnUpdateVASPContacts        func(ctx context.Context, vaspID string, contacts *models.Contacts) error
}

func Open() (*Store, error) {
	return &Store{}, nil
}

func (s *Store) Calls(call string) int {
	s.RLock()
	defer s.RUnlock()
	if s.calls == nil {
		return 0
	}
	return s.calls[call]
}

func (s *Store) Invoked(call string) bool {
	return s.Calls(call) > 0
}

func (s *Store) UseError(call string, err error) {
	// TODO: complete the rest of the calls as needed for tests.
	switch call {
	case Close:
		s.OnClose = func() error { return err }
	case Reindex:
		s.OnReindex = func() error { return err }
	case Backup:
		s.OnBackup = func(s string) error { return err }
	case ListVASPs:
		s.OnListVASPs = func(ctx context.Context) iterator.DirectoryIterator { return BadDirectoryIterator(err) }
	case SearchVASPs:
		s.OnSearchVASPs = func(ctx context.Context, query map[string]interface{}) ([]*pb.VASP, error) { return nil, err }
	case CreateVASP:
		s.OnCreateVASP = func(ctx context.Context, v *pb.VASP) (string, error) { return "", err }
	case RetrieveVASP:
		s.OnRetrieveVASP = func(ctx context.Context, id string) (*pb.VASP, error) { return nil, err }
	case UpdateVASP:
		s.OnUpdateVASP = func(ctx context.Context, v *pb.VASP) error { return err }
	case DeleteVASP:
		s.OnDeleteVASP = func(ctx context.Context, id string) error { return err }
	case CountVASPs:
		s.OnCountVASPs = func(ctx context.Context) (uint64, error) { return 0, err }
	default:
		panic(fmt.Errorf("unknown call %q", call))
	}
}

func (s *Store) Reset() {
	s.Lock()
	defer s.Unlock()
	for key := range s.calls {
		s.calls[key] = 0
	}

	s.OnClose = nil
	s.OnReindex = nil
	s.OnBackup = nil
	s.OnListVASPs = nil
	s.OnSearchVASPs = nil
	s.OnCreateVASP = nil
	s.OnRetrieveVASP = nil
	s.OnUpdateVASP = nil
	s.OnDeleteVASP = nil
	s.OnCountVASPs = nil
	s.OnListCertReqs = nil
	s.OnCreateCertReq = nil
	s.OnRetrieveCertReq = nil
	s.OnUpdateCertReq = nil
	s.OnDeleteCertReq = nil
	s.OnCountCertReqs = nil
	s.OnListCerts = nil
	s.OnCreateCert = nil
	s.OnRetrieveCert = nil
	s.OnUpdateCert = nil
	s.OnDeleteCert = nil
	s.OnCountCerts = nil
	s.OnRetrieveAnnouncementMonth = nil
	s.OnUpdateAnnouncementMonth = nil
	s.OnDeleteAnnouncementMonth = nil
	s.OnCountAnnouncementMonths = nil
	s.OnRetrieveActivityMonth = nil
	s.OnUpdateActivityMonth = nil
	s.OnDeleteActivityMonth = nil
	s.OnCountActivityMonth = nil
	s.OnListOrganizations = nil
	s.OnCreateOrganization = nil
	s.OnRetrieveOrganization = nil
	s.OnUpdateOrganization = nil
	s.OnDeleteOrganization = nil
	s.OnCountOrganizations = nil
	s.OnListContacts = nil
	s.OnCreateContact = nil
	s.OnRetrieveContact = nil
	s.OnUpdateContact = nil
	s.OnDeleteContact = nil
	s.OnCountContacts = nil
	s.OnListEmails = nil
	s.OnCreateEmail = nil
	s.OnRetrieveEmail = nil
	s.OnUpdateEmail = nil
	s.OnDeleteEmail = nil
	s.OnCountEmails = nil
	s.OnVASPContacts = nil
	s.OnRetrieveVASPContacts = nil
	s.OnUpdateVASPContacts = nil
}

func (s *Store) Close() error {
	s.incrCalls(Close)
	if s.OnClose != nil {
		return s.OnClose()
	}
	panic(ErrNoMock)
}

func (s *Store) Reindex() error {
	s.incrCalls(Reindex)
	if s.OnReindex != nil {
		return s.OnReindex()
	}
	panic(ErrNoMock)
}

func (s *Store) Backup(path string) error {
	s.incrCalls(Backup)
	if s.OnBackup != nil {
		return s.OnBackup(path)
	}
	panic(ErrNoMock)
}

func (s *Store) ListVASPs(ctx context.Context) iterator.DirectoryIterator {
	s.incrCalls(ListVASPs)
	if s.OnListVASPs != nil {
		return s.OnListVASPs(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) SearchVASPs(ctx context.Context, query map[string]interface{}) ([]*pb.VASP, error) {
	s.incrCalls(SearchVASPs)
	if s.OnSearchVASPs != nil {
		return s.OnSearchVASPs(ctx, query)
	}
	panic(ErrNoMock)
}

func (s *Store) CreateVASP(ctx context.Context, v *pb.VASP) (string, error) {
	s.incrCalls(CreateVASP)
	if s.OnCreateVASP != nil {
		return s.OnCreateVASP(ctx, v)
	}
	panic(ErrNoMock)
}

func (s *Store) RetrieveVASP(ctx context.Context, id string) (*pb.VASP, error) {
	s.incrCalls(RetrieveVASP)
	if s.OnRetrieveVASP != nil {
		return s.OnRetrieveVASP(ctx, id)
	}
	panic(ErrNoMock)
}

func (s *Store) UpdateVASP(ctx context.Context, v *pb.VASP) error {
	s.incrCalls(UpdateVASP)
	if s.OnUpdateVASP != nil {
		return s.OnUpdateVASP(ctx, v)
	}
	panic(ErrNoMock)
}

func (s *Store) DeleteVASP(ctx context.Context, id string) error {
	s.incrCalls(DeleteVASP)
	if s.OnDeleteVASP != nil {
		return s.OnDeleteVASP(ctx, id)
	}
	panic(ErrNoMock)
}

func (s *Store) CountVASPs(ctx context.Context) (uint64, error) {
	s.incrCalls(CountVASPs)
	if s.OnCountVASPs != nil {
		return s.OnCountVASPs(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) ListCertReqs(ctx context.Context) iterator.CertificateRequestIterator {
	s.incrCalls(ListCertReqs)
	if s.OnListCertReqs != nil {
		return s.OnListCertReqs(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) CreateCertReq(ctx context.Context, r *models.CertificateRequest) (string, error) {
	s.incrCalls(CreateCertReq)
	if s.OnCreateCertReq != nil {
		return s.OnCreateCertReq(ctx, r)
	}
	panic(ErrNoMock)
}

func (s *Store) RetrieveCertReq(ctx context.Context, id string) (*models.CertificateRequest, error) {
	s.incrCalls(RetrieveCertReq)
	if s.OnRetrieveCertReq != nil {
		return s.OnRetrieveCertReq(ctx, id)
	}
	panic(ErrNoMock)
}

func (s *Store) UpdateCertReq(ctx context.Context, r *models.CertificateRequest) error {
	s.incrCalls(UpdateCertReq)
	if s.OnUpdateCertReq != nil {
		return s.OnUpdateCertReq(ctx, r)
	}
	panic(ErrNoMock)
}

func (s *Store) DeleteCertReq(ctx context.Context, id string) error {
	s.incrCalls(DeleteCertReq)
	if s.OnDeleteCertReq != nil {
		return s.OnDeleteCertReq(ctx, id)
	}
	panic(ErrNoMock)
}

func (s *Store) CountCertReqs(ctx context.Context) (uint64, error) {
	s.incrCalls(CountCertReqs)
	if s.OnCountCertReqs != nil {
		return s.OnCountCertReqs(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) ListCerts(ctx context.Context) iterator.CertificateIterator {
	s.incrCalls(ListCerts)
	if s.OnListCerts != nil {
		return s.OnListCerts(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) CreateCert(ctx context.Context, c *models.Certificate) (string, error) {
	s.incrCalls(CreateCert)
	if s.OnCreateCert != nil {
		return s.OnCreateCert(ctx, c)
	}
	panic(ErrNoMock)
}

func (s *Store) RetrieveCert(ctx context.Context, id string) (*models.Certificate, error) {
	s.incrCalls(RetrieveCert)
	if s.OnRetrieveCert != nil {
		return s.OnRetrieveCert(ctx, id)
	}
	panic(ErrNoMock)
}

func (s *Store) UpdateCert(ctx context.Context, c *models.Certificate) error {
	s.incrCalls(UpdateCert)
	if s.OnUpdateCert != nil {
		return s.OnUpdateCert(ctx, c)
	}
	panic(ErrNoMock)
}

func (s *Store) DeleteCert(ctx context.Context, id string) error {
	s.incrCalls(DeleteCert)
	if s.OnDeleteCert != nil {
		return s.OnDeleteCert(ctx, id)
	}
	panic(ErrNoMock)
}

func (s *Store) CountCerts(ctx context.Context) (uint64, error) {
	s.incrCalls(CountCerts)
	if s.OnCountCerts != nil {
		return s.OnCountCerts(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) RetrieveAnnouncementMonth(ctx context.Context, date string) (*bff.AnnouncementMonth, error) {
	s.incrCalls(RetrieveAnnouncementMonth)
	if s.OnRetrieveAnnouncementMonth != nil {
		return s.OnRetrieveAnnouncementMonth(ctx, date)
	}
	panic(ErrNoMock)
}

func (s *Store) UpdateAnnouncementMonth(ctx context.Context, m *bff.AnnouncementMonth) error {
	s.incrCalls(UpdateAnnouncementMonth)
	if s.OnUpdateAnnouncementMonth != nil {
		return s.OnUpdateAnnouncementMonth(ctx, m)
	}
	panic(ErrNoMock)
}

func (s *Store) DeleteAnnouncementMonth(ctx context.Context, date string) error {
	s.incrCalls(DeleteAnnouncementMonth)
	if s.OnDeleteAnnouncementMonth != nil {
		return s.OnDeleteAnnouncementMonth(ctx, date)
	}
	panic(ErrNoMock)
}

func (s *Store) CountAnnouncementMonths(ctx context.Context) (uint64, error) {
	s.incrCalls(CountAnnouncementMonths)
	if s.OnCountAnnouncementMonths != nil {
		return s.OnCountAnnouncementMonths(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) RetrieveActivityMonth(ctx context.Context, date string) (*bff.ActivityMonth, error) {
	s.incrCalls(RetrieveActivityMonth)
	if s.OnRetrieveActivityMonth != nil {
		return s.OnRetrieveActivityMonth(ctx, date)
	}
	panic(ErrNoMock)
}

func (s *Store) UpdateActivityMonth(ctx context.Context, m *bff.ActivityMonth) error {
	s.incrCalls(UpdateActivityMonth)
	if s.OnUpdateActivityMonth != nil {
		return s.OnUpdateActivityMonth(ctx, m)
	}
	panic(ErrNoMock)
}

func (s *Store) DeleteActivityMonth(ctx context.Context, date string) error {
	s.incrCalls(DeleteActivityMonth)
	if s.OnDeleteActivityMonth != nil {
		return s.OnDeleteActivityMonth(ctx, date)
	}
	panic(ErrNoMock)
}

func (s *Store) CountActivityMonth(ctx context.Context) (uint64, error) {
	s.incrCalls(CountActivityMonth)
	if s.OnCountActivityMonth != nil {
		return s.OnCountActivityMonth(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) ListOrganizations(ctx context.Context) iterator.OrganizationIterator {
	s.incrCalls(ListOrganizations)
	if s.OnListOrganizations != nil {
		return s.OnListOrganizations(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) CreateOrganization(ctx context.Context, o *bff.Organization) (string, error) {
	s.incrCalls(CreateOrganization)
	if s.OnCreateOrganization != nil {
		return s.OnCreateOrganization(ctx, o)
	}
	panic(ErrNoMock)
}

func (s *Store) RetrieveOrganization(ctx context.Context, id uuid.UUID) (*bff.Organization, error) {
	s.incrCalls(RetrieveOrganization)
	if s.OnRetrieveOrganization != nil {
		return s.OnRetrieveOrganization(ctx, id)
	}
	panic(ErrNoMock)
}

func (s *Store) UpdateOrganization(ctx context.Context, o *bff.Organization) error {
	s.incrCalls(UpdateOrganization)
	if s.OnUpdateOrganization != nil {
		return s.OnUpdateOrganization(ctx, o)
	}
	panic(ErrNoMock)
}

func (s *Store) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	s.incrCalls(DeleteOrganization)
	if s.OnDeleteOrganization != nil {
		return s.OnDeleteOrganization(ctx, id)
	}
	panic(ErrNoMock)
}

func (s *Store) CountOrganizations(ctx context.Context) (uint64, error) {
	s.incrCalls(CountOrganizations)
	if s.OnCountOrganizations != nil {
		return s.OnCountOrganizations(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) ListContacts(ctx context.Context) []*models.Contact {
	s.incrCalls(ListContacts)
	if s.OnListContacts != nil {
		return s.OnListContacts(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) CreateContact(ctx context.Context, c *models.Contact) (string, error) {
	s.incrCalls(CreateContact)
	if s.OnCreateContact != nil {
		return s.OnCreateContact(ctx, c)
	}
	panic(ErrNoMock)
}

func (s *Store) RetrieveContact(ctx context.Context, email string) (*models.Contact, error) {
	s.incrCalls(RetrieveContact)
	if s.OnRetrieveContact != nil {
		return s.OnRetrieveContact(ctx, email)
	}
	panic(ErrNoMock)
}

func (s *Store) UpdateContact(ctx context.Context, c *models.Contact) error {
	s.incrCalls(UpdateContact)
	if s.OnUpdateContact != nil {
		return s.OnUpdateContact(ctx, c)
	}
	panic(ErrNoMock)
}

func (s *Store) DeleteContact(ctx context.Context, email string) error {
	s.incrCalls(DeleteContact)
	if s.OnDeleteContact != nil {
		return s.OnDeleteContact(ctx, email)
	}
	panic(ErrNoMock)
}

func (s *Store) CountContacts(ctx context.Context) (uint64, error) {
	s.incrCalls(CountContacts)
	if s.OnCountContacts != nil {
		return s.OnCountContacts(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) ListEmails(ctx context.Context) iterator.EmailIterator {
	s.incrCalls(ListEmails)
	if s.OnListEmails != nil {
		return s.OnListEmails(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) CreateEmail(ctx context.Context, c *models.Email) (string, error) {
	s.incrCalls(CreateEmail)
	if s.OnCreateEmail != nil {
		return s.OnCreateEmail(ctx, c)
	}
	panic(ErrNoMock)
}

func (s *Store) RetrieveEmail(ctx context.Context, email string) (*models.Email, error) {
	s.incrCalls(RetrieveEmail)
	if s.OnRetrieveEmail != nil {
		return s.OnRetrieveEmail(ctx, email)
	}
	panic(ErrNoMock)
}

func (s *Store) UpdateEmail(ctx context.Context, c *models.Email) error {
	s.incrCalls(UpdateEmail)
	if s.OnUpdateEmail != nil {
		return s.OnUpdateEmail(ctx, c)
	}
	panic(ErrNoMock)
}

func (s *Store) DeleteEmail(ctx context.Context, email string) error {
	s.incrCalls(DeleteEmail)
	if s.OnDeleteEmail != nil {
		return s.OnDeleteEmail(ctx, email)
	}
	panic(ErrNoMock)
}

func (s *Store) CountEmails(ctx context.Context) (uint64, error) {
	s.incrCalls(CountEmails)
	if s.OnCountEmails != nil {
		return s.OnCountEmails(ctx)
	}
	panic(ErrNoMock)
}

func (s *Store) VASPContacts(ctx context.Context, vasp *pb.VASP) (*models.Contacts, error) {
	s.incrCalls(VASPContacts)
	if s.OnVASPContacts != nil {
		return s.OnVASPContacts(ctx, vasp)
	}
	panic(ErrNoMock)
}

func (s *Store) RetrieveVASPContacts(ctx context.Context, vaspID string) (*models.Contacts, error) {
	s.incrCalls(RetrieveVASPContacts)
	if s.OnRetrieveVASPContacts != nil {
		return s.OnRetrieveVASPContacts(ctx, vaspID)
	}
	panic(ErrNoMock)
}

func (s *Store) UpdateVASPContacts(ctx context.Context, vaspID string, contacts *models.Contacts) error {
	s.incrCalls(UpdateVASPContacts)
	if s.OnUpdateVASPContacts != nil {
		return s.OnUpdateVASPContacts(ctx, vaspID, contacts)
	}
	panic(ErrNoMock)
}

func (s *Store) incrCalls(call string) {
	s.Lock()
	defer s.Unlock()
	if s.calls == nil {
		s.calls = make(map[string]int)
	}
	s.calls[call]++
}
