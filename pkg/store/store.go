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
		Count() (uint64, error) 					   // Count the number of objects
	}

Ideally there would be a store per namespace, but in order to generalize the store to
multiple embedded databases, the store interface affixes the object store methods with
the namespace. E.g. ListVASPs, CreateCertReq, etc.
*/
package store

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	bff "github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/store/config"
	"github.com/trisacrypto/directory/pkg/store/iterator"
	"github.com/trisacrypto/directory/pkg/store/leveldb"
	"github.com/trisacrypto/directory/pkg/store/mock"
	"github.com/trisacrypto/directory/pkg/store/trtl"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Open a directory storage provider with the specified URI. Database URLs should either
// specify protocol+transport://user:pass@host/dbname?opt1=a&opt2=b for servers or
// protocol:///relative/path/to/file for embedded databases (for absolute paths, specify
// protocol:////absolute/path/to/file).
//
// To open a mock store, use the DSN mock:///
func Open(conf config.StoreConfig) (s Store, err error) {
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
	case "mock":
		if s, err = mock.Open(); err != nil {
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
	AnnouncementStore
	ActivityStore
	OrganizationStore
	ContactStore
	EmailStore
	DirectoryContactStore
}

// leveldb.Store and trtl.Store must implement the Store interface.
var _ Store = &leveldb.Store{}
var _ Store = &trtl.Store{}
var _ Store = &mock.Store{}

// DirectoryStore describes how services interact with VASP identity records.
type DirectoryStore interface {
	ListVASPs(ctx context.Context) iterator.DirectoryIterator
	SearchVASPs(ctx context.Context, query map[string]interface{}) ([]*pb.VASP, error)
	CreateVASP(ctx context.Context, v *pb.VASP) (string, error)
	RetrieveVASP(ctx context.Context, id string) (*pb.VASP, error)
	UpdateVASP(ctx context.Context, v *pb.VASP) error
	DeleteVASP(ctx context.Context, id string) error
	CountVASPs(ctx context.Context) (uint64, error)
}

// CertificateRequestStore describes how services interact with Certificate requests.
type CertificateRequestStore interface {
	ListCertReqs(ctx context.Context) iterator.CertificateRequestIterator
	CreateCertReq(ctx context.Context, r *models.CertificateRequest) (string, error)
	RetrieveCertReq(ctx context.Context, id string) (*models.CertificateRequest, error)
	UpdateCertReq(ctx context.Context, r *models.CertificateRequest) error
	DeleteCertReq(ctx context.Context, id string) error
	CountCertReqs(context.Context) (uint64, error)
}

// CertificateStore describes how services interact with Certificate records.
type CertificateStore interface {
	ListCerts(ctx context.Context) iterator.CertificateIterator
	CreateCert(ctx context.Context, c *models.Certificate) (string, error)
	RetrieveCert(ctx context.Context, id string) (*models.Certificate, error)
	UpdateCert(ctx context.Context, c *models.Certificate) error
	DeleteCert(ctx context.Context, id string) error
	CountCerts(context.Context) (uint64, error)
}

// AnnouncementStore describes how services interact with the Announcement records.
type AnnouncementStore interface {
	RetrieveAnnouncementMonth(ctx context.Context, date string) (*bff.AnnouncementMonth, error)
	UpdateAnnouncementMonth(ctx context.Context, m *bff.AnnouncementMonth) error
	DeleteAnnouncementMonth(ctx context.Context, date string) error
	CountAnnouncementMonths(context.Context) (uint64, error)
}

// ActivityStore describes how services interact with the Activity records.
type ActivityStore interface {
	RetrieveActivityMonth(ctx context.Context, date string) (*bff.ActivityMonth, error)
	UpdateActivityMonth(ctx context.Context, m *bff.ActivityMonth) error
	DeleteActivityMonth(ctx context.Context, date string) error
	CountActivityMonth(context.Context) (uint64, error)
}

// OrganizationStore describes how services interact with the Organization records.
type OrganizationStore interface {
	ListOrganizations(ctx context.Context) iterator.OrganizationIterator
	CreateOrganization(ctx context.Context, o *bff.Organization) (string, error)
	RetrieveOrganization(ctx context.Context, id uuid.UUID) (*bff.Organization, error)
	UpdateOrganization(ctx context.Context, o *bff.Organization) error
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	CountOrganizations(context.Context) (uint64, error)
}

// ContactStore describes how services interact with the Contact records.
type ContactStore interface {
	ListContacts(ctx context.Context) []*models.Contact
	CreateContact(ctx context.Context, c *models.Contact) (string, error)
	RetrieveContact(ctx context.Context, email string) (*models.Contact, error)
	UpdateContact(ctx context.Context, c *models.Contact) error
	DeleteContact(ctx context.Context, email string) error
	CountContacts(context.Context) (uint64, error)
}

// EmailStore describes how services interact with Email records.
type EmailStore interface {
	ListEmails(ctx context.Context) iterator.EmailIterator
	CreateEmail(ctx context.Context, c *models.Email) (string, error)
	RetrieveEmail(ctx context.Context, email string) (*models.Email, error)
	UpdateEmail(ctx context.Context, c *models.Email) error
	DeleteEmail(ctx context.Context, email string) error
	CountEmails(context.Context) (uint64, error)
}

// DirectoryContactStore joins VASP contact records with Email records.
type DirectoryContactStore interface {
	VASPContacts(ctx context.Context, vasp *pb.VASP) (*models.Contacts, error)
	RetrieveVASPContacts(ctx context.Context, vaspID string) (*models.Contacts, error)
	UpdateVASPContacts(ctx context.Context, vaspID string, contacts *models.Contacts) error
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
