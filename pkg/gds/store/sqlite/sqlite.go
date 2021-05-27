package sqlite

import (
	"errors"

	"github.com/trisacrypto/directory/pkg/gds/models/v1"
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

func (s *Store) Create(v *pb.VASP) (id string, err error) {
	return "", nil
}

func (s *Store) Retrieve(id string) (v *pb.VASP, err error) {
	return nil, nil
}

func (s *Store) Update(v *pb.VASP) (err error) {
	return nil
}

func (s *Store) Destroy(id string) (err error) {
	return nil
}

func (s *Store) Search(query map[string]interface{}) (vasps []*pb.VASP, err error) {
	return nil, nil
}

func (s *Store) ListCertRequests() (reqs []*models.CertificateRequest, err error) {
	return nil, nil
}

func (s *Store) GetCertRequest(id string) (r *models.CertificateRequest, err error) {
	return nil, nil
}

func (s *Store) SaveCertRequest(r *models.CertificateRequest) (err error) {
	return nil
}

func (s *Store) DeleteCertRequest(id string) (err error) {
	return nil
}

func (s *Store) Backup(path string) (err error) {
	return nil
}
