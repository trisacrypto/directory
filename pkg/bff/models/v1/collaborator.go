package models

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
)

// Returns the key which uniquely identifies this collaborator.
func (collab *Collaborator) Key() string {
	if collab.Email == "" {
		return ""
	}

	hash := md5.Sum([]byte(collab.Email))
	return base64.RawStdEncoding.EncodeToString(hash[:])
}

// Validate a collaborator record, ensuring that all the required fields exist for
// storage and generating missing fields such as the ID.
func (collab *Collaborator) Validate() error {
	// Collaborator must be indexable by email address
	if collab.Email == "" {
		return errors.New("collaborator is missing email address")
	}

	// Collaborator must have an ID
	id := collab.Key()
	if collab.Id == "" {
		collab.Id = id
	}

	// Make sure the ID matches the email address
	if collab.Id != id {
		return errors.New("collaborator has invalid id")
	}

	// TODO: More comprehensive validation may be required
	return nil
}
