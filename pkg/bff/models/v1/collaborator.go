package models

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"time"
)

// Returns the key which uniquely identifies this collaborator.
func (collab *Collaborator) Key() string {
	if collab.Email == "" {
		return ""
	}

	// Key must be url-safe since it will be used in the URL path on various endpoints
	hash := md5.Sum([]byte(collab.Email))
	return base64.RawURLEncoding.EncodeToString(hash[:])
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

// Helper to determine if a collaborator invite is valid based on the expiration date.
// If there is no expiration date, this method assumes that the collaborator invitation
// is still valid.
func (collab *Collaborator) ValidateInvitation() (err error) {
	if collab.ExpiresAt == "" {
		return nil
	}

	var expiration time.Time
	if expiration, err = time.Parse(time.RFC3339Nano, collab.ExpiresAt); err != nil {
		return err
	}

	if expiration.Before(time.Now()) {
		return errors.New("collaborator invite has expired")
	}
	return nil
}
