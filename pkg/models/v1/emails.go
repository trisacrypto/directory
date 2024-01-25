package models

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/exp/slices"
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

// Adds a vaspID to the list of vasps associated with the email address; if the vasp is
// already in the list it is not duplicated and the list is maintained in sorted order.
func (e *Email) AddVASP(vaspID string) {
	pos, found := slices.BinarySearch(e.Vasps, vaspID)
	if found {
		return
	}

	e.Vasps = append(e.Vasps, "")
	copy(e.Vasps[pos+1:], e.Vasps[pos:])
	e.Vasps[pos] = vaspID
}

// Removes a vaspID from the list of vasps associated with the email address.
func (e *Email) RmVASP(vaspID string) {
	pos, found := slices.BinarySearch(e.Vasps, vaspID)
	if !found {
		return
	}

	e.Vasps = append(e.Vasps[:pos], e.Vasps[pos+1:]...)
}

// Counts emails within the given EmailLogEntry slice for the given reason within the given time frame.
func CountSentEmails(emailLog []*EmailLogEntry, reason string, timeWindowDays int) (sent int, err error) {
	if reason == "" {
		return 0, ErrNoLogReason
	}

	if timeWindowDays < 0 {
		return 0, ErrInvalidWindow
	}

	for _, value := range emailLog {
		if reason != value.Reason {
			continue
		}

		var timestamp time.Time
		if timestamp, err = time.Parse(time.RFC3339, value.Timestamp); err != nil {
			return 0, fmt.Errorf("invalid timestamp: %w", err)
		}

		if timestamp.After(time.Now().AddDate(0, 0, -timeWindowDays)) {
			sent++
		}
	}
	return sent, nil
}

// Normalize an email address to just the user@domain component.
func NormalizeEmail(email string) string {
	if addr, err := mail.ParseAddress(email); err == nil {
		return strings.ToLower(strings.TrimSpace(addr.Address))
	}

	// Otherwise just return the string lowercased without spaces
	return strings.ToLower(strings.TrimSpace(email))
}
