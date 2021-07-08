package mockdb

import (
	models "github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// MockDB fulfills the store interface for testing.
type MockDB struct {
	OnCreate      func(v *pb.VASP) (string, error)
	CreateInvoked bool

	OnRetrieve      func(id string) (*pb.VASP, error)
	RetrieveInvoked bool

	OnRetrieveAll      func(opts *models.RetrieveAllOpts, c chan *pb.VASP) error
	RetrieveAllInvoked bool

	OnUpdate      func(v *pb.VASP) error
	UpdateInvoked bool

	OnUpdateStatus      func(id string, status int32) error
	UpdateStatusInvoked bool

	OnDestroy      func(id string) error
	DestroyInvoked bool

	OnSearch      func(query map[string]interface{}) ([]*pb.VASP, error)
	SearchInvoked bool
}

func (m *MockDB) Create(v *pb.VASP) (string, error) {
	m.CreateInvoked = true
	return m.OnCreate(v)
}
func (m *MockDB) Retrieve(id string) (*pb.VASP, error) {
	m.RetrieveInvoked = true
	return m.OnRetrieve(id)
}
func (m *MockDB) RetrieveAll(opts *models.RetrieveAllOpts, c chan *pb.VASP) error {
	m.RetrieveAllInvoked = true
	return m.OnRetrieveAll(opts, c)
}
func (m *MockDB) Update(v *pb.VASP) error {
	m.UpdateInvoked = true
	return m.OnUpdate(v)
}
func (m *MockDB) UpdateStatus(id string, status int32) error {
	m.UpdateStatusInvoked = true
	return m.OnUpdateStatus(id, status)
}
func (m *MockDB) Destroy(id string) error {
	m.DestroyInvoked = true
	return m.OnDestroy(id)
}
func (m *MockDB) Search(query map[string]interface{}) ([]*pb.VASP, error) {
	m.SearchInvoked = true
	return m.OnSearch(query)
}
