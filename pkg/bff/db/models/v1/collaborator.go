package models

import "errors"

// Returns the key which uniquely identifies this collaborator.
func (collab *Collaborator) Key() string {
	return collab.Email
}

// Validate a collaborator record, ensuring that all the required fields exist for
// storage.
func (collab *Collaborator) Validate() error {
	// Collaborators must be indexable by email address
	if collab.Email == "" {
		return errors.New("collaborator is missing email address")
	}

	// TODO: More comprehensive validation may be required
	return nil
}
