package sqlite

import (
	"errors"

	"github.com/trisacrypto/directory/pkg/gds/global/v1"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/peers/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Open SQLite directory Store at the specified path.
func Open(path string) (*Store, error) {
	return &Store{}, errors.New("sqlite3 store is not implemented yet")
}

type Store struct{}

func (s *Store) Close() error {
	return nil
}

func (s *Store) CreateVASP(v *pb.VASP) (id string, err error) {
	return "", nil
}

func (s *Store) RetrieveVASP(id string) (v *pb.VASP, err error) {
	return nil, nil
}

func (s *Store) UpdateVASP(v *pb.VASP) (err error) {
	return nil
}

func (s *Store) DeleteVASP(id string) (err error) {
	return nil
}

func (s *Store) SearchVASPs(query map[string]interface{}) (vasps []*pb.VASP, err error) {
	return nil, nil
}

func (s *Store) ListCertReqs() (reqs []*models.CertificateRequest, err error) {
	return nil, nil
}

func (s *Store) CreateCertReq(r *models.CertificateRequest) (id string, err error) {
	return "", nil
}

func (s *Store) RetrieveCertReq(id string) (r *models.CertificateRequest, err error) {
	return nil, nil
}

func (s *Store) UpdateCertReq(r *models.CertificateRequest) (err error) {
	return nil
}

func (s *Store) DeleteCertReq(id string) (err error) {
	return nil
}

func (s *Store) ListPeers() (pl []*peers.Peer, err error) {
	return nil, nil
}

func (s *Store) CreatePeer(p *peers.Peer) (id string, err error) {
	return "", nil
}

func (s *Store) RetrievePeer(id string) (p *peers.Peer, err error) {
	return nil, nil
}

func (s *Store) DeletePeer(id string) error {
	return nil
}
func (s *Store) Backup(path string) (err error) {
	return nil
}

func (s *Store) WithVersionManager(vm *global.VersionManager) error {
	return nil
}
