package mockdb

import (
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/gds/store/iterator"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

var _ store.Store = &MockDB{}

// MockDB fulfills the store interface for testing.
type MockDB struct {
	OnClose      func() error
	CloseInvoked bool

	OnCreateVASP      func(v *pb.VASP) (string, error)
	CreateVASPInvoked bool

	OnRetrieveVASP      func(id string) (*pb.VASP, error)
	RetrieveVASPInvoked bool

	OnUpdateVASP      func(v *pb.VASP) error
	UpdateVASPInvoked bool

	OnDeleteVASP      func(id string) error
	DeleteVASPInvoked bool

	OnListVASPs      func() iterator.DirectoryIterator
	ListVASPsInvoked bool

	OnSearchVASPs      func(query map[string]interface{}) ([]*pb.VASP, error)
	SearchVASPsInvoked bool

	OnListCertReqs      func() iterator.CertificateIterator
	ListCertReqsInvoked bool

	OnCreateCertReq      func(r *models.CertificateRequest) (string, error)
	CreateCertReqInvoked bool

	OnRetrieveCertReq      func(id string) (*models.CertificateRequest, error)
	RetrieveCertReqInvoked bool

	OnUpdateCertReq      func(r *models.CertificateRequest) error
	UpdateCertReqInvoked bool

	OnDeleteCertReq      func(id string) error
	DeleteCertReqInvoked bool

	OnListPeers      func() iterator.ReplicaIterator
	ListPeersInvoked bool

	OnCreatePeer      func(p *peers.Peer) (string, error)
	CreatePeerInvoked bool

	OnRetrievePeer      func(id string) (*peers.Peer, error)
	RetrievePeerInvoked bool

	OnDeletePeer      func(id string) error
	DeletePeerInvoked bool

	OnReindex      func() error
	ReindexInvoked bool

	OnBackup      func(string) error
	BackupInvoked bool
}

func (m *MockDB) Close() error {
	m.CloseInvoked = true
	return m.OnClose()
}

func (m *MockDB) CreateVASP(v *pb.VASP) (string, error) {
	m.CreateVASPInvoked = true
	return m.OnCreateVASP(v)
}

func (m *MockDB) RetrieveVASP(id string) (*pb.VASP, error) {
	m.RetrieveVASPInvoked = true
	return m.OnRetrieveVASP(id)
}

func (m *MockDB) UpdateVASP(v *pb.VASP) error {
	m.UpdateVASPInvoked = true
	return m.OnUpdateVASP(v)
}

func (m *MockDB) DeleteVASP(id string) error {
	m.DeleteVASPInvoked = true
	return m.OnDeleteVASP(id)
}

func (m *MockDB) ListVASPs() iterator.DirectoryIterator {
	m.ListVASPsInvoked = true
	return m.OnListVASPs()
}

func (m *MockDB) SearchVASPs(query map[string]interface{}) ([]*pb.VASP, error) {
	m.SearchVASPsInvoked = true
	return m.OnSearchVASPs(query)
}

func (m *MockDB) ListCertReqs() iterator.CertificateIterator {
	m.ListCertReqsInvoked = true
	return m.OnListCertReqs()
}

func (m *MockDB) CreateCertReq(r *models.CertificateRequest) (string, error) {
	m.CreateCertReqInvoked = true
	return m.OnCreateCertReq(r)
}

func (m *MockDB) RetrieveCertReq(id string) (*models.CertificateRequest, error) {
	m.RetrieveCertReqInvoked = true
	return m.OnRetrieveCertReq(id)
}

func (m *MockDB) UpdateCertReq(r *models.CertificateRequest) error {
	m.UpdateCertReqInvoked = true
	return m.OnUpdateCertReq(r)
}

func (m *MockDB) DeleteCertReq(id string) error {
	m.DeleteCertReqInvoked = true
	return m.OnDeleteCertReq(id)
}

func (m *MockDB) ListPeers() iterator.ReplicaIterator {
	m.ListPeersInvoked = true
	return m.OnListPeers()
}

func (m *MockDB) CreatePeer(p *peers.Peer) (string, error) {
	m.CreatePeerInvoked = true
	return m.OnCreatePeer(p)
}

func (m *MockDB) RetrievePeer(id string) (*peers.Peer, error) {
	m.RetrievePeerInvoked = true
	return m.OnRetrievePeer(id)
}

func (m *MockDB) DeletePeer(id string) error {
	m.DeletePeerInvoked = true
	return m.OnDeletePeer(id)
}

func (m *MockDB) Reindex() error {
	m.ReindexInvoked = true
	return m.OnReindex()
}

func (m *MockDB) Backup(path string) error {
	m.BackupInvoked = true
	return m.OnBackup(path)
}
