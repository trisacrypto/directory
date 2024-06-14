package iterator

import (
	bff "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Iterators allow memory safe list operations from the Store.
type Iterator interface {
	Next() bool
	Prev() bool
	Error() error
	Release()
}

// DirectoryIterator allows access to DirectoryStore models
type DirectoryIterator interface {
	Iterator
	Id() string
	VASP() (*pb.VASP, error)
	All() ([]*pb.VASP, error)
	SeekId(vaspID string) bool
}

// CertificateRequestIterator allows access to CertificateRequestStore models
type CertificateRequestIterator interface {
	Iterator
	CertReq() (*models.CertificateRequest, error)
	All() ([]*models.CertificateRequest, error)
}

// CertificateIterator allows access to CertificateStore models
type CertificateIterator interface {
	Iterator
	Cert() (*models.Certificate, error)
	All() ([]*models.Certificate, error)
}

// OrganizationIterator allows access to OrganizationStore models
type OrganizationIterator interface {
	Iterator
	ID() string
	Organization() (*bff.Organization, error)
}

// EmailIterator allows access to Email models
type EmailIterator interface {
	Iterator
	Email() (*models.Email, error)
}
