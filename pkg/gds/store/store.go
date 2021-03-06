/*
Package store provides an interface to multiple types of embedded storage across
multiple objects. Unlike a SQL interface, the TRISA directory service relies on document
databases or key/value stores such as leveldb or trtl. It also manages multiple
namespaces (object types) - VASP records, CertificateRequests, Peers, etc. In general an
object store interface provides accesses to the objects, with one interface per
namespace as follows:

	type ObjectStore interface {
		List() *Iterator                               // Iterate over all objects
		Search(query map[string]interface{}) *Iterator // Create a query to list filtered objects
		Create(o *Object) (id string, err error)       // Make the object
		Retrieve(id string) (o *Object, err error)     // Fetch an object by ID or by key
		Update(o *Object) error                        // Make changes to an object
		Delete(id string) error                        // Delete an object
	}

Ideally there would be a store per namespace, but in order to generalize the store to
multiple embedded databases, the store interface affixes the object store methods with
the namespace. E.g. ListVASPs, CreateCertReq, etc.
*/
package store

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store/iterator"
	"github.com/trisacrypto/directory/pkg/gds/store/leveldb"
	"github.com/trisacrypto/directory/pkg/gds/store/trtl"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Open a directory storage provider with the specified URI. Database URLs should either
// specify protocol+transport://user:pass@host/dbname?opt1=a&opt2=b for servers or
// protocol:///relative/path/to/file for embedded databases (for absolute paths, specify
// protocol:////absolute/path/to/file).
func Open(conf config.DatabaseConfig) (s Store, err error) {
	var dsn *DSN
	if dsn, err = ParseDSN(conf.URL); err != nil {
		return nil, err
	}

	switch dsn.Scheme {
	case "leveldb":
		if s, err = leveldb.Open(dsn.Path); err != nil {
			return nil, err
		}
	case "trtl":
		if s, err = trtl.Open(conf); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unhandled database scheme %q", dsn.Scheme)
	}

	if conf.ReindexOnBoot {
		if indexer, ok := s.(Indexer); ok {
			if err = indexer.Reindex(); err != nil {
				// NOTE: the database is not closed here, so if reindexing fails,
				// something very bad might have occurred and the server should stop.
				return nil, err
			}
			log.Info().Str("scheme", dsn.Scheme).Msg("store reindexed")
		} else {
			log.Warn().Str("scheme", dsn.Scheme).Msg("store is not an indexer - skipping reindex")
		}
	}
	return s, nil
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
	CertificateRequestStore
}

// DirectoryStore describes how the service interacts with VASP identity records.
type DirectoryStore interface {
	ListVASPs() iterator.DirectoryIterator
	SearchVASPs(query map[string]interface{}) ([]*pb.VASP, error)
	CreateVASP(v *pb.VASP) (string, error)
	RetrieveVASP(id string) (*pb.VASP, error)
	UpdateVASP(v *pb.VASP) error
	DeleteVASP(id string) error
}

// CertificateRequestStore describes how the service interacts with Certificate requests.
type CertificateRequestStore interface {
	ListCertReqs() iterator.CertificateRequestIterator
	CreateCertReq(r *models.CertificateRequest) (string, error)
	RetrieveCertReq(id string) (*models.CertificateRequest, error)
	UpdateCertReq(r *models.CertificateRequest) error
	DeleteCertReq(id string) error
}

// CertificateStore describes how the service interacts with Certificate records.
type CertificateStore interface {
	ListCerts() iterator.CertificateIterator
	CreateCert(c *models.Certificate) (string, error)
	RetrieveCert(id string) (*models.Certificate, error)
	UpdateCert(c *models.Certificate) error
	DeleteCert(id string) error
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
