package contacts

import (
	"errors"
	"time"

	"github.com/trisacrypto/directory/pkg/models/v1"
)

// Returns True if a Contact is not nil and has an email address.
func ContactHasEmail(contact *models.Contact) bool {
	return contact != nil && contact.Email != ""
}

// Returns True if a Contact is verified.
func ContactIsVerified(contact *models.Contact) (verified bool) {
	// Return zero-valued defaults with no error if extra is nil.
	if contact == nil {
		return false
	}
	return contact.Verified
}

// GetContactVerification token and verified status from the extra data field on the Contact.
func GetContactVerification(contact *models.Contact) (string, bool) {
	// Return zero-valued defaults with no error if extra is nil.
	if contact == nil {
		return "", false
	}
	return contact.Token, contact.Verified
}

// SetContactVerification token and verified status on the Contact record.
func SetContactVerification(contact *models.Contact, token string, verified bool) (err error) {
	if contact == nil {
		return errors.New("cannot set verification on nil contact")
	}

	// Set contact verification.
	contact.Verified = verified
	contact.Token = token
	return nil
}

// GetEmailLog from the extra data on the Contact record.
func GetEmailLog(contact *models.Contact) (_ []*models.EmailLogEntry) {
	// If the extra data is nil, return nil (no email log).
	if contact == nil {
		return nil
	}
	return contact.EmailLog
}

// Create and add a new entry to the EmailLog on the extra data on the Contact record.
func AppendEmailLog(contact *models.Contact, reason, subject string) (err error) {
	// Contact must be non-nil.
	if contact == nil {
		return errors.New("cannot append entry to nil contact")
	}

	// Create the EmailLog if it is nil.
	if contact.EmailLog == nil {
		contact.EmailLog = make([]*models.EmailLogEntry, 0, 1)
	}

	// Append entry to the previous log.
	entry := &models.EmailLogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Reason:    reason,
		Subject:   subject,
		Recipient: contact.Email,
	}
	contact.EmailLog = append(contact.EmailLog, entry)
	return nil
}
