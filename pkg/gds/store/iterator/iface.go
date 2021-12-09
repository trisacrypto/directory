package iterator

import (
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Iterators allow memory safe list operations from the Store.
type Iterator interface {
	Next() bool
	Error() error
	Release()
}

// DirectoryIterator allows access to DirectoryStore models
type DirectoryIterator interface {
	Iterator
	VASP() *pb.VASP
	All() ([]*pb.VASP, error)
}

// CertificateIterator allows access to CertificateStore models
type CertificateIterator interface {
	Iterator
	CertReq() *models.CertificateRequest
	All() ([]*models.CertificateRequest, error)
}
