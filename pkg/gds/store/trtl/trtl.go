package trtl

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/client"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store/iterator"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
)

// Open a connection to the Trtl database.
func Open(profile *client.TrtlProfile) (store *Store, err error) {
	store = &Store{}
	if store.conn, err = profile.Connect(); err != nil {
		return nil, err
	}
	store.client = pb.NewTrtlClient(store.conn)
	return &Store{}, nil
}

// Store implements the store.Store interface for the Trtl replicated database.
type Store struct {
	conn   *grpc.ClientConn
	client pb.TrtlClient
}

//===========================================================================
// Store Implementation
//===========================================================================

// Close the connection to the database.
func (s *Store) Close() error {
	return s.conn.Close()
}

//===========================================================================
// DirectoryStore Implementation
//===========================================================================

func (s *Store) ListVASPs() iterator.DirectoryIterator {
	log.Debug().Msg("not implemented")
	return nil
}

func (s *Store) SearchVASPs(query map[string]interface{}) ([]*gds.VASP, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) CreateVASP(vasp *gds.VASP) error {
	return errors.New("not implemented")
}

func (s *Store) RetrieveVASP(id string) (*gds.VASP, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) UpdateVASP(v *gds.VASP) error {
	return errors.New("not implemented")
}

func (s *Store) DeleteVASP(id string) error {
	return errors.New("not implemented")
}

//===========================================================================
// CertificateStore Implementation
//===========================================================================

func (s *Store) ListCertReqs() iterator.CertificateIterator {
	log.Debug().Msg("not implemented")
	return nil
}

func (s *Store) CreateCertReq(r *models.CertificateRequest) (string, error) {
	return "", errors.New("not implemented")
}

func (s *Store) RetrieveCertReq(id string) (*models.CertificateRequest, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) UpdateCertReq(r *models.CertificateRequest) error {
	return errors.New("not implemented")
}

func (s *Store) DeleteCertReq(id string) error {
	return errors.New("not implemented")
}
