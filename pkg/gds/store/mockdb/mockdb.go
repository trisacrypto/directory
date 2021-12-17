package mockdb

import (
	"errors"
	"fmt"

	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/store/iterator"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

var db store.Store = &MockDB{}
var state = &MockState{
	VASPs: make(map[string]pb.VASP),
	Keys:  []string{},
}

// MockState contains the current state of the MockDB for test verification.
type MockState struct {
	// in-memory database store
	VASPs map[string]pb.VASP
	Keys  []string

	// keep track of store interface calls
	CloseInvoked           bool
	CreateVASPInvoked      bool
	RetrieveVASPInvoked    bool
	UpdateVASPInvoked      bool
	DeleteVASPInvoked      bool
	ListVASPsInvoked       bool
	SearchVASPsInvoked     bool
	ListCertReqsInvoked    bool
	CreateCertReqInvoked   bool
	RetrieveCertReqInvoked bool
	UpdateCertReqInvoked   bool
	DeleteCertReqInvoked   bool
	ReindexInvoked         bool
	BackupInvoked          bool
}

func GetState() *MockState {
	return state
}

func ResetState() {
	state = &MockState{
		VASPs: make(map[string]pb.VASP),
		Keys:  []string{},
	}
}

// MockDB fulfills the store interface for testing.
type MockDB struct {
	OnClose           func() error
	OnCreateVASP      func(v *pb.VASP) (string, error)
	OnRetrieveVASP    func(id string) (*pb.VASP, error)
	OnUpdateVASP      func(v *pb.VASP) error
	OnDeleteVASP      func(id string) error
	OnListVASPs       func() iterator.DirectoryIterator
	OnSearchVASPs     func(query map[string]interface{}) ([]*pb.VASP, error)
	OnListCertReqs    func() iterator.CertificateIterator
	OnCreateCertReq   func(r *models.CertificateRequest) (string, error)
	OnRetrieveCertReq func(id string) (*models.CertificateRequest, error)
	OnUpdateCertReq   func(r *models.CertificateRequest) error
	OnDeleteCertReq   func(id string) error
	OnReindex         func() error
	OnBackup          func(string) error
}

func GetStore() store.Store {
	return db
}

func (m *MockDB) Close() error {
	state.CloseInvoked = true
	return m.OnClose()
}

func (m *MockDB) CreateVASP(v *pb.VASP) (string, error) {
	state.CreateVASPInvoked = true
	if v.Id == "" {
		return "", errors.New("VASP must contain an ID")
	}
	if _, ok := state.VASPs[v.Id]; ok {
		return "", fmt.Errorf("VASP with ID %s already exists", v.Id)
	}
	state.VASPs[v.Id] = *v
	state.Keys = append(state.Keys, v.Id)
	return v.Id, nil
}

func (m *MockDB) RetrieveVASP(id string) (*pb.VASP, error) {
	state.RetrieveVASPInvoked = true
	if id == "" {
		return nil, errors.New("missing VASP ID")
	}
	var v pb.VASP
	var ok bool
	if v, ok = state.VASPs[id]; !ok {
		return nil, fmt.Errorf("VASP with ID %s not found", id)
	}
	return &v, nil
}

func (m *MockDB) UpdateVASP(v *pb.VASP) error {
	state.UpdateVASPInvoked = true
	return m.OnUpdateVASP(v)
}

func (m *MockDB) DeleteVASP(id string) error {
	state.DeleteVASPInvoked = true
	return m.OnDeleteVASP(id)
}

func (m *MockDB) ListVASPs() iterator.DirectoryIterator {
	state.ListVASPsInvoked = true
	return m.OnListVASPs()
}

func (m *MockDB) SearchVASPs(query map[string]interface{}) ([]*pb.VASP, error) {
	state.SearchVASPsInvoked = true
	return m.OnSearchVASPs(query)
}

func (m *MockDB) ListCertReqs() iterator.CertificateIterator {
	state.ListCertReqsInvoked = true
	return m.OnListCertReqs()
}

func (m *MockDB) CreateCertReq(r *models.CertificateRequest) (string, error) {
	state.CreateCertReqInvoked = true
	return m.OnCreateCertReq(r)
}

func (m *MockDB) RetrieveCertReq(id string) (*models.CertificateRequest, error) {
	state.RetrieveCertReqInvoked = true
	return m.OnRetrieveCertReq(id)
}

func (m *MockDB) UpdateCertReq(r *models.CertificateRequest) error {
	state.UpdateCertReqInvoked = true
	return m.OnUpdateCertReq(r)
}

func (m *MockDB) DeleteCertReq(id string) error {
	state.DeleteCertReqInvoked = true
	return m.OnDeleteCertReq(id)
}

func (m *MockDB) Reindex() error {
	state.ReindexInvoked = true
	return m.OnReindex()
}

func (m *MockDB) Backup(path string) error {
	state.BackupInvoked = true
	return m.OnBackup(path)
}
