package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	TechnicalContact        = "technical"
	AdministrativeContact   = "administrative"
	LegalContact            = "legal"
	BillingContact          = "billing"
	VerificationTokenLength = 48
)

type contactType struct {
	contact *pb.Contact
	kind    string
	err     error
}

// Returns True if a Contact is not nil and has an email address.
func ContactHasEmail(contact *pb.Contact) bool {
	return contact != nil && contact.Email != ""
}

// Returns True if a Contact is verified.
func ContactIsVerified(contact *pb.Contact) (verified bool, err error) {
	if _, verified, err = GetContactVerification(contact); err != nil {
		return false, err
	}
	return verified, nil
}

// Returns True if the contact kind is one of the recognized strings.
func ContactKindIsValid(kind string) bool {
	kinds := map[string]struct{}{
		AdministrativeContact: {},
		BillingContact:        {},
		LegalContact:          {},
		TechnicalContact:      {},
	}
	_, ok := kinds[kind]
	return ok
}

// Returns the corresponding contact object for the given contact type.
func ContactFromType(contacts *pb.Contacts, kind string) *pb.Contact {
	switch kind {
	case AdministrativeContact:
		return contacts.Administrative
	case BillingContact:
		return contacts.Billing
	case LegalContact:
		return contacts.Legal
	case TechnicalContact:
		return contacts.Technical
	}
	return nil
}

// Adds a contact on the VASP object.
func AddContact(vasp *pb.VASP, kind string, contact *pb.Contact) error {
	switch kind {
	case AdministrativeContact:
		vasp.Contacts.Administrative = contact
	case BillingContact:
		vasp.Contacts.Billing = contact
	case LegalContact:
		vasp.Contacts.Legal = contact
	case TechnicalContact:
		vasp.Contacts.Technical = contact
	default:
		return fmt.Errorf("invalid contact type: %s", kind)
	}
	return nil
}

// Deletes a contact on the VASP object by setting it to nil.
func DeleteContact(vasp *pb.VASP, kind string) error {
	switch kind {
	case AdministrativeContact:
		vasp.Contacts.Administrative = nil
	case BillingContact:
		vasp.Contacts.Billing = nil
	case LegalContact:
		vasp.Contacts.Legal = nil
	case TechnicalContact:
		vasp.Contacts.Technical = nil
	default:
		return fmt.Errorf("invalid contact type: %s", kind)
	}
	return nil
}

// Returns a standardized order of iterating through contacts.
func contactOrder(contacts *pb.Contacts) []*contactType {
	return []*contactType{
		{contact: contacts.Technical, kind: TechnicalContact},
		{contact: contacts.Administrative, kind: AdministrativeContact},
		{contact: contacts.Legal, kind: LegalContact},
		{contact: contacts.Billing, kind: BillingContact},
	}
}

type ContactIterator struct {
	email    bool
	verified bool
	index    int
	contacts []*contactType
}

// Returns a new ContactIterator which has Next() and Value() methods that can be used
// to iterate over contacts in a Contacts object.
func NewContactIterator(contacts *pb.Contacts, requireEmail bool, requireVerified bool) *ContactIterator {
	return &ContactIterator{
		email:    requireEmail,
		verified: requireVerified,
		index:    -1,
		contacts: contactOrder(contacts),
	}
}

// Moves the ContactIterator to the next existing contact and returns true if there is
// one, false if the iteration is complete.
func (i *ContactIterator) Next() bool {
	// Advance the index (note, index must be initialized at -1)
	i.index++

	// Stopping condition
	if i.index >= len(i.contacts) {
		return false
	}

	// Filter checks - contact must exist
	current := i.contacts[i.index]
	if current.contact == nil || current.contact.IsZero() {
		return i.Next()
	}

	// Filter if email required
	if i.email && !ContactHasEmail(current.contact) {
		return i.Next()
	}

	// Filter if verified is required
	if i.verified {
		if verified, err := ContactIsVerified(current.contact); err != nil || !verified {
			// Even in an error we're skipping the contact, errors have to be
			// fetched as a multi-error after the iteration is complete.
			if err != nil {
				current.err = fmt.Errorf("error retrieving verification status for %s contact: %s", current.kind, err.Error())
			}
			return i.Next()
		}
	}

	return true
}

func (i *ContactIterator) Value() (*pb.Contact, string) {
	if i.index < len(i.contacts) {
		// Note that no checking of the contact occurs here to allow
		// us to allow iteration with different filtering mechanisms.
		current := i.contacts[i.index]
		return current.contact, current.kind
	}
	return nil, ""
}

// Error returns a multi-error that describes any errors that occurred during
// iteration for any contact. This will allow us to log all errors once.
func (i *ContactIterator) Error() (err error) {
	var errs *multierror.Error
	for _, contact := range i.contacts {
		errs = multierror.Append(errs, contact.err)
	}
	return errs.ErrorOrNil()
}

// GetContactVerification token and verified status from the extra data field on the Contact.
func GetContactVerification(contact *pb.Contact) (_ string, _ bool, err error) {
	// Return zero-valued defaults with no error if extra is nil.
	if contact == nil || contact.Extra == nil {
		return "", false, nil
	}

	// Unmarshal the extra data field on the Contact
	extra := &GDSContactExtraData{}
	if err = contact.Extra.UnmarshalTo(extra); err != nil {
		return "", false, err
	}
	return extra.GetToken(), extra.GetVerified(), nil
}

// SetContactVerification token and verified status on the Contact record.
func SetContactVerification(contact *pb.Contact, token string, verified bool) (err error) {
	if contact == nil || contact.IsZero() {
		return errors.New("cannot set verification on nil contact")
	}

	// Unmarshal previous extra data.
	extra := &GDSContactExtraData{}
	if contact.Extra != nil {
		if err = contact.Extra.UnmarshalTo(extra); err != nil {
			return fmt.Errorf("could not deserialize previous extra: %s", err)
		}
	}

	// Set contact verification.
	extra.Verified = verified
	extra.Token = token
	if contact.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

// VerifiedContacts returns a map of contact type to email address for all verified
// contacts, omitting any contacts that are not verified or do not exist.
func VerifiedContacts(vasp *pb.VASP) (contacts map[string]string) {
	contacts = make(map[string]string)
	iter := NewContactIterator(vasp.Contacts, false, true)
	for iter.Next() {
		contact, kind := iter.Value()
		contacts[kind] = contact.Email
	}
	return contacts
}

// ContactVerifications returns a map of contact type to verified status, omitting any
// contacts that do not exist.
func ContactVerifications(vasp *pb.VASP) (contacts map[string]bool, errs *multierror.Error) {
	contacts = make(map[string]bool)
	iter := NewContactIterator(vasp.Contacts, false, false)
	for iter.Next() {
		contact, kind := iter.Value()
		if verified, err := ContactIsVerified(contact); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("error retrieving verification status for %s contact: %s", kind, err))
		} else {
			contacts[kind] = verified
		}
	}
	return contacts, errs
}

// GetEmailLog from the extra data on the Contact record.
func GetEmailLog(contact *pb.Contact) (_ []*EmailLogEntry, err error) {
	// If the extra data is nil, return nil (no email log).
	if contact == nil || contact.Extra == nil {
		return nil, nil
	}

	// Unmarshal the extra data field on the VASP.
	extra := &GDSContactExtraData{}
	if err = contact.Extra.UnmarshalTo(extra); err != nil {
		return nil, err
	}
	return extra.GetEmailLog(), nil
}

// Create and add a new entry to the EmailLog on the extra data on the Contact record.
func AppendEmailLog(contact *pb.Contact, reason string, subject string) (err error) {
	// Contact must be non-nil.
	if contact == nil || contact.IsZero() {
		return errors.New("cannot append entry to nil contact")
	}

	// Unmarshal previous extra data.
	extra := &GDSContactExtraData{}
	if contact.Extra != nil {
		if err = contact.Extra.UnmarshalTo(extra); err != nil {
			return fmt.Errorf("could not deserialize previous extra: %s", err)
		}
	} else {
		extra.EmailLog = make([]*EmailLogEntry, 0, 1)
	}

	// Append entry to the previous log.
	entry := &EmailLogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Reason:    reason,
		Subject:   subject,
	}
	extra.EmailLog = append(extra.EmailLog, entry)

	// Serialize the extra data back to the VASP.
	if contact.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

// GetSentEmailCount returns the number of emails sent to the contact for the given reason.
func GetSentEmailCount(contact *pb.Contact, reason string, timeWindowDays int) (sent int, err error) {
	if reason == "" {
		return 0, errors.New("cannot match on empty reason string")
	}

	if timeWindowDays < 0 {
		return 0, errors.New("time window must be a positive number of days")
	}

	emailLog, err := GetEmailLog(contact)
	if err != nil {
		return 0, err
	}

	for _, value := range emailLog {
		strTimestamp := value.Timestamp
		timestamp, err := time.Parse(time.RFC3339, strTimestamp)
		if err != nil {
			return 0, fmt.Errorf("error parsing timestamp: %v", err)
		}

		matchedReason := reason == value.Reason
		withinTimeWindow := timestamp.After(time.Now().AddDate(0, 0, -timeWindowDays))

		if matchedReason && withinTimeWindow {
			sent++
		}
	}
	return sent, nil
}
