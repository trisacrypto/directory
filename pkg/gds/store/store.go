/*
Package store provides an interface to database storage for the TRISA directory service.
*/
package store

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/trisacrypto/directory/pkg/gds/store/leveldb"
	"github.com/trisacrypto/directory/pkg/gds/store/sqlite"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Open a directory storage provider with the specified URI. Database URLs should either
// specify protocol+transport://user:pass@host/dbname?opt1=a&opt2=b for servers or
// protocol:///relative/path/to/file for embedded databases (for absolute paths, specify
// protocol:////absolute/path/to/file).
func Open(uri string) (_ Store, err error) {
	var dsn *DSN
	if dsn, err = ParseDSN(uri); err != nil {
		return nil, err
	}

	switch dsn.Scheme {
	case "leveldb":
		return leveldb.Open(dsn.Path)
	case "sqlite", "sqlite3":
		return sqlite.Open(dsn.Path)
	default:
		return nil, fmt.Errorf("unhandled database scheme %q", dsn.Scheme)
	}
}

// Store provides an interface for directory storage services to abstract the underlying
// database provider. The storage methods correspond to directory service requests,
// which are currently implemented with a simple CRUD and search interface for VASP
// records and certificate requests. The underlying database can be a simple embedded
// store or a distributed SQL server, so long as it can interact with identity records.
type Store interface {
	Close() error
	DirectoryStore
	CertificateStore
}

// DirectoryStore describes how the service interacts with VASP identity records.
type DirectoryStore interface {
	Create(v *pb.VASP) (string, error)
	Retrieve(id string) (*pb.VASP, error)
	Update(v *pb.VASP) error
	Destroy(id string) error
	Search(query map[string]interface{}) ([]*pb.VASP, error)
}

// CertificateStore describes how the service interacts with Certificate requests.
type CertificateStore interface {
	ListCertRequests() ([]*pb.CertificateRequest, error)
	GetCertRequest(id string) (*pb.CertificateRequest, error)
	SaveCertRequest(r *pb.CertificateRequest) error
	DeleteCertRequest(id string) error
}

// Indexer allows external methods to access the index function of the store if it has
// them. E.g. a leveldb embedded database or other store that uses an in-memory index
// needs to be an Indexer but not a SQL database.
type Indexer interface {
	Reindex() error
}

// Backup means that the Store can be backed up to a compressed location on disk,
// optionally with encryption if its required.
type Backup interface {
	Backup(string) error
}

// DSN represents the parsed components of an embedded database service.
type DSN struct {
	Scheme string
	Path   string
}

// DSN Parsing and Handling
func ParseDSN(uri string) (_ *DSN, err error) {
	dsn, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("could not parse dsn: %s", err)
	}

	if dsn.Scheme == "" || dsn.Path == "" {
		return nil, errors.New("could not parse dsn, specify scheme:///relative/path/to/db")
	}

	return &DSN{
		Scheme: dsn.Scheme,
		Path:   strings.TrimPrefix(dsn.Path, "/"),
	}, nil
}
