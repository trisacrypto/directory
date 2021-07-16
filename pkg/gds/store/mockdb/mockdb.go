package mockdb

import (
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

var _ store.Store = &MockDB{}

// MockDB fulfills the store interface for testing.
type MockDB struct {
	OnCreate      func(v *pb.VASP) (string, error)
	CreateInvoked bool

	OnRetrieve      func(id string) (*pb.VASP, error)
	RetrieveInvoked bool

	OnUpdate      func(v *pb.VASP) error
	UpdateInvoked bool

	OnDestroy      func(id string) error
	DestroyInvoked bool

	OnSearch      func(query map[string]interface{}) ([]*pb.VASP, error)
	SearchInvoked bool

	OnClose      func() error
	CloseInvoked bool

	OnListCertRequests      func() ([]*models.CertificateRequest, error)
	ListCertRequestsInvoked bool

	OnGetCertRequest      func(id string) (*models.CertificateRequest, error)
	GetCertRequestInvoked bool

	OnSaveCertRequest      func(r *models.CertificateRequest) error
	SaveCertRequestInvoked bool

	OnDeleteCertRequest      func(id string) error
	DeleteCertRequestInvoked bool
}

func (m *MockDB) Create(v *pb.VASP) (string, error) {
	m.CreateInvoked = true
	return m.OnCreate(v)
}

func (m *MockDB) Retrieve(id string) (*pb.VASP, error) {
	m.RetrieveInvoked = true
	return m.OnRetrieve(id)
}

func (m *MockDB) Update(v *pb.VASP) error {
	m.UpdateInvoked = true
	return m.OnUpdate(v)
}

func (m *MockDB) Destroy(id string) error {
	m.DestroyInvoked = true
	return m.OnDestroy(id)
}

func (m *MockDB) Search(query map[string]interface{}) ([]*pb.VASP, error) {
	m.SearchInvoked = true
	return m.OnSearch(query)
}

func (m *MockDB) Close() error {
	m.CloseInvoked = true
	return m.OnClose()
}

func (m *MockDB) ListCertRequests() ([]*models.CertificateRequest, error) {
	m.ListCertRequestsInvoked = true
	return m.OnListCertRequests()
}

func (m *MockDB) GetCertRequest(id string) (*models.CertificateRequest, error) {
	m.GetCertRequestInvoked = true
	return m.OnGetCertRequest(id)
}

func (m *MockDB) SaveCertRequest(r *models.CertificateRequest) error {
	m.SaveCertRequestInvoked = true
	return m.OnSaveCertRequest(r)
}

func (m *MockDB) DeleteCertRequest(id string) error {
	m.DeleteCertRequestInvoked = true
	return m.OnDeleteCertRequest(id)
}
