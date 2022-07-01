package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	"google.golang.org/protobuf/proto"
)

const (
	NamespaceOrganizations = "organizations"
)

// The Organizations collection exposes methods for interacting with Organization
// records in the database. Organizations are given UUIDs when created, and those uuids
// are also used as keys in the database for quick Gets by ID.
//
// Organizations implements the Collection interface
type Organizations struct {
	db        *DB
	namespace string
}

// Ensure that Organizations implements the Collection interface.
var _ Collection = &Organizations{}

// Organizations constructs the collection type for db interactions with the namespace.
// This method is intended to be used with chaining, e.g. db.Organizations().Get(id). To
// reduce the number of allocations a singleton us used. Method calls to the collection
// are thread-safe.
func (db *DB) Organizations() *Organizations {
	db.makeOrganizations.Do(func() {
		db.organizations = &Organizations{
			db:        db,
			namespace: NamespaceOrganizations,
		}
	})
	return db.organizations
}

// Create an organization and get back the organization stub.
func (o *Organizations) Create(ctx context.Context) (org *models.Organization, err error) {
	// Create an empty organization
	ts := time.Now().Format(time.RFC3339Nano)
	uu := uuid.New()
	org = &models.Organization{
		Id:       uu.String(),
		Created:  ts,
		Modified: ts,
	}

	var value []byte
	if value, err = proto.Marshal(org); err != nil {
		return nil, err
	}

	if err = o.db.Put(ctx, uu[:], value, o.namespace); err != nil {
		return nil, err
	}

	return org, nil
}

// Retrieve an organization by it's ID, which can be either a []byte, string, or uuid.
func (o *Organizations) Retrieve(ctx context.Context, orgID interface{}) (org *models.Organization, err error) {
	var uu uuid.UUID
	if uu, err = models.ParseOrgID(orgID); err != nil {
		return nil, err
	}

	var data []byte
	if data, err = o.db.Get(ctx, uu[:], o.namespace); err != nil {
		return nil, err
	}

	org = &models.Organization{}
	if err = proto.Unmarshal(data, org); err != nil {
		return nil, err
	}

	return org, nil
}

// Update an organization with the record supplied.
func (o *Organizations) Update(ctx context.Context, org *models.Organization) (err error) {
	// Set modified timestamp and serialize
	org.Modified = time.Now().Format(time.RFC3339Nano)

	var data []byte
	if data, err = proto.Marshal(org); err != nil {
		return err
	}

	if err = o.db.Put(ctx, org.Key(), data, o.namespace); err != nil {
		return err
	}
	return nil
}

// Delete an organization by it's ID, which can be either a []byte, string, or uuid.
func (o *Organizations) Delete(ctx context.Context, orgID interface{}) (err error) {
	var uu uuid.UUID
	if uu, err = models.ParseOrgID(orgID); err != nil {
		return err
	}

	if err = o.db.Delete(ctx, uu[:], o.namespace); err != nil {
		return err
	}
	return nil
}

// Namespace implements the collection interface
func (o *Organizations) Namespace() string {
	return o.namespace
}
