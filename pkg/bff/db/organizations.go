package db

import (
	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

// Create an organization and get back the organization stub with an assigned ID.
func (store *DB) CreateOrganization() (*models.Organization, error) {
	return store.db.CreateOrganization()
}

// Retrieve an organization by its ID, which can be either a []byte, string, or uuid.
func (store *DB) RetrieveOrganization(orgID interface{}) (_ *models.Organization, err error) {
	var uu uuid.UUID
	if uu, err = models.ParseOrgID(orgID); err != nil {
		return nil, err
	}

	return store.db.RetrieveOrganization(uu)
}

// Update an organization with the record supplied.
func (store *DB) UpdateOrganization(org *models.Organization) error {
	return store.db.UpdateOrganization(org)
}

// Delete an organization by its ID, which can be either a []byte, string, or uuid.
func (store *DB) DeleteOrganization(orgID interface{}) (err error) {
	var uu uuid.UUID
	if uu, err = models.ParseOrgID(orgID); err != nil {
		return err
	}

	return store.db.DeleteOrganization(uu)
}
