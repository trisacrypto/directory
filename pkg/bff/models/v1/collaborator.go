package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

// Returns the key which uniquely identifies this collaborator.
func (collab *Collaborator) Key() string {
	if collab.Email == "" {
		return ""
	}

	hash := md5.Sum([]byte(collab.Email))
	return hex.EncodeToString(hash[:])
}

// Validate a collaborator record, ensuring that all the required fields exist for
// storage and generating missing fields such as the ID.
func (collab *Collaborator) Validate() error {
	switch {
	case collab.Email == "":
		return errors.New("collaborator is missing email address")
	case collab.Id == "":
		collab.Id = collab.Key()
	case collab.Id != collab.Key():
		return errors.New("collaborator has invalid id")
	}

	// TODO: More comprehensive validation may be required
	return nil
}
