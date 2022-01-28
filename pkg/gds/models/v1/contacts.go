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
	current  *contactType
	contacts []*contactType
}

// Returns a new ContactIterator which has Next() and Value() methods that can be used
// to iterate over contacts in a Contacts object.
func NewContactIterator(contacts *pb.Contacts, requireEmail bool, requireVerified bool) *ContactIterator {
	return &ContactIterator{
		email:    requireEmail,
		verified: requireVerified,
		contacts: contactOrder(contacts),
	}
}

// Moves the ContactIterator to the next existing contact and returns true if there is
// one, false if the iteration is complete.
func (i *ContactIterator) Next() bool {
	for ; i.index < len(i.contacts); i.index++ {
		i.current = i.contacts[i.index]
		contact := i.current.contact
		if contact != nil {
			if i.email && !ContactHasEmail(contact) {
				continue
			}
			if i.verified {
				var verified bool
				var err error
				if verified, err = ContactIsVerified(contact); err != nil {
					// If we can't retrieve the contact verification, let the caller
					// handle the error on Value().
					i.current.err = err
					i.index++
					return true
				}
				if !verified {
					continue
				}
			}
			i.index++
			return true
		}
	}
	i.current = nil
	return false
}

func (i *ContactIterator) Value() (*pb.Contact, string, error) {
	if i.current != nil {
		return i.current.contact, i.current.kind, i.current.err
	}
	return nil, "", errors.New("no more contacts")
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
func VerifiedContacts(vasp *pb.VASP) (contacts map[string]string) {
	contacts = make(map[string]string)
	iter := NewContactIterator(vasp.Contacts, false, true)
	for iter.Next() {
		if contact, kind, err := iter.Value(); err == nil {
			contacts[kind] = contact.Email
		}
	}
	return contacts
}

// ContactVerifications returns a map of contact type to verified status, omitting any
// contacts that do not exist.
func ContactVerifications(vasp *pb.VASP) (contacts map[string]bool, err error) {
	contacts = make(map[string]bool)
	iter := NewContactIterator(vasp.Contacts, false, false)
	for iter.Next() {
		if contact, kind, err := iter.Value(); err == nil {
			if verified, err := ContactIsVerified(contact); err == nil {
				contacts[kind] = verified
			}
		}
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
