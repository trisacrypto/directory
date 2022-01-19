package models

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	TechnicalContact      = "technical"
	AdministrativeContact = "administrative"
	LegalContact          = "legal"
	BillingContact        = "billing"
)

type contactType struct {
	contact *pb.Contact
	kind    string
}

// Returns True if a Contact is not nil and has an email address.
func ContactHasEmail(contact *pb.Contact) bool {
	return contact != nil && contact.Email != ""
}

// Returns True if a Contact is verified.
func ContactVerified(contact *pb.Contact) (verified bool, err error) {
	if _, verified, err = GetContactVerification(contact); err != nil {
		return false, err
	}
	return verified, nil
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

// Returns a function which iterates over the verified contacts in a Contacts object.
func IterContacts(contacts *pb.Contacts, requireEmail bool) func() (*pb.Contact, string) {
	all := contactOrder(contacts)
	i := 0
	return func() (contact *pb.Contact, kind string) {
		// Return the next contact that exists or nil to indicate end of iteration.
		for i < len(all) {
			contact = all[i].contact
			kind = all[i].kind
			i++
			if contact != nil && (!requireEmail || ContactHasEmail(contact)) {
				return contact, kind
			}
		}
		return nil, ""
	}
}

// Returns a function which iterates over the verified contacts in a Contacts object.
func IterVerifiedContacts(contacts *pb.Contacts) func() (*pb.Contact, string, error) {
	all := contactOrder(contacts)
	i := 0
	return func() (contact *pb.Contact, kind string, err error) {
		// Return the next contact that exists or nil to indicate end of iteration.
		for i < len(all) {
			contact = all[i].contact
			kind = all[i].kind
			i++
			if ContactHasEmail(contact) {
				var verified bool
				if verified, err = ContactVerified(contact); err != nil {
					return nil, "", err
				}
				if verified {
					return contact, kind, nil
				}
			}
		}
		return nil, "", nil
	}
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
	if contact == nil {
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
func VerifiedContacts(vasp *pb.VASP) (contacts map[string]string, err error) {
	contacts = make(map[string]string)
	next := IterVerifiedContacts(vasp.Contacts)
	contact, kind, err := next()
	for ; err == nil && contact != nil; contact, kind, err = next() {
		contacts[kind] = contact.Email
	}
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

// ContactVerifications returns a map of contact type to verified status, omitting any
// contacts that do not exist.
func ContactVerifications(vasp *pb.VASP) (contacts map[string]bool, err error) {
	contacts = make(map[string]bool)
	next := IterContacts(vasp.Contacts, false)
	for contact, kind := next(); contact != nil; contact, kind = next() {
		var verified bool
		if verified, err = ContactVerified(contact); err != nil {
			return nil, err
		}
		contacts[kind] = verified
	}
	return contacts, nil
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
	if contact == nil {
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
