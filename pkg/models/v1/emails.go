package models

import (
	"net/mail"
	"strings"
	"time"
)

// Validate that the email record is complete and ensure the email and name are
// normalized correctly to ensure that the email record is handled uniformly.
func (e *Email) Validate() error {
	// Parse the name from the email if the name is not set
	if e.Name == "" {
		if addr, err := mail.ParseAddress(e.Email); err == nil {
			e.Name = addr.Name
		}
	}

	// Normalize the email address
	e.Email = NormalizeEmail(e.Email)

	// Ensure the email exists
	if e.Email == "" {
		return ErrNoEmailAddress
	}

	if e.Verified {
		if e.VerifiedOn == "" || e.Token != "" {
			return ErrVerifiedInvalid
		}
	} else {
		if e.VerifiedOn != "" || e.Token == "" {
			return ErrUnverifiedInvalid
		}
	}

	return nil
}

// Appends an email log entry to the send log of the email. Expects the log to be in
// time order. Note that the email record must be serialized and stored.
func (e *Email) Log(reason, subject string) {
	if e.SendLog == nil {
		e.SendLog = make([]*EmailLogEntry, 0, 1)
	}

	e.SendLog = append(e.SendLog, &EmailLogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Reason:    reason,
		Subject:   subject,
		Recipient: e.Email,
	})
}

// Normalize an email address to just the user@domain component.
func NormalizeEmail(email string) string {
	if addr, err := mail.ParseAddress(email); err == nil {
		return strings.ToLower(strings.TrimSpace(addr.Address))
	}

	// Otherwise just return the string lowercased without spaces
	return strings.ToLower(strings.TrimSpace(email))
}
