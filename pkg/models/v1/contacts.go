package models

import (
	"errors"
	"fmt"
	"time"

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

// Contacts wraps a VASPs contacts with their email records for easier access to
// contact records, including iteration over contact records. This type of record is
// often created from database records, and acts as a join record between the VASPs
// table and the emails table to prevent duplicate emails from being sent.
type Contacts struct {
	Contacts *pb.Contacts
	Emails   []*Email
}

// Contact wraps a VASP contact type record with it's email address record.
// TODO: rename to contact once we remove the contact protocol buffer
type ContactRecord struct {
	Kind    string      // e.g. technical, administrative, legal, billing, etc.
	Contact *pb.Contact // the contact record from the VASP object
	Email   *Email      // the email record associated with the contact
}

// Returns True if a Contact is nil or is empty.
func ContactIsZero(contact *pb.Contact) bool {
	return contact == nil || contact.IsZero()
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

// Get returns the contact for the specified kind if the contact on the VASP is nil then
// a nil contact is returned, otherwise the contact record is constructed and returned.
// This method panics if the specified kind is invalid.
func (c *Contacts) Get(kind string) *ContactRecord {
	record := &ContactRecord{Kind: kind}
	switch kind {
	case TechnicalContact:
		record.Contact = c.Contacts.Technical
	case AdministrativeContact:
		record.Contact = c.Contacts.Administrative
	case LegalContact:
		record.Contact = c.Contacts.Legal
	case BillingContact:
		record.Contact = c.Contacts.Billing
	default:
		panic(fmt.Errorf("invalid contact kind %q", kind))
	}

	// Return a nil record if the VASP contact is nil
	if record.Contact == nil {
		return nil
	}

	// Find the email for the given contact
	for _, email := range c.Emails {
		if email.Email == record.Contact.Email {
			record.Email = email
			break
		}
	}

	return record
}

// Index returns the contact at the specified index guaranteeing a specific ordering to
// contacts if you iterate from 0 to 3; the ordering is: technical, administrative,
// legal, then billing. You can also use the contacts iterator to loop over the contacts
// with specific options. If the contact does not exist then nil is returned; if the
// index is not in range [0,3] then this method panics.
func (c *Contacts) Index(i int) *ContactRecord {
	switch i {
	case 0:
		return c.Get(TechnicalContact)
	case 1:
		return c.Get(AdministrativeContact)
	case 2:
		return c.Get(LegalContact)
	case 3:
		return c.Get(BillingContact)
	default:
		panic(fmt.Errorf("index %d is not in range of contacts", i))
	}
}

type ContactsIterator struct {
	skipNoEmail    bool
	skipUnverified bool
	skipVerified   bool
	skipDuplicates bool
	index          int
	contacts       *Contacts
	emails         map[string]struct{}
}

// ContactIterOptions are used to configure iterator behavior
type ContactIterOption func(c *ContactsIterator)

// Skip contacts with no email address
func SkipNoEmail() ContactIterOption {
	return func(c *ContactsIterator) {
		c.skipNoEmail = true
	}
}

// Skip contacts that are not verified
func SkipUnverified() ContactIterOption {
	return func(c *ContactsIterator) {
		c.skipUnverified = true
	}
}

// Skip contacts that are verified
func SkipVerified() ContactIterOption {
	return func(c *ContactsIterator) {
		c.skipVerified = true
	}
}

// Skip contacts with duplicate email addresses, e.g. if the technical contact is the
// same as the administrative contact then only the technical contact will be returned
// in the iteration.
func SkipDuplicates() ContactIterOption {
	return func(c *ContactsIterator) {
		c.skipDuplicates = true
	}
}

// Returns a new ContactIterator which has Next() and Value() methods that can be used
// to iterate over contacts in a Contacts object. By default, the iterator will return
// each contact which is non-zero. SkipNoEmail(), SkipUnverified(), and
// SkipDuplicates() can be used to filter contacts returned by the iterator.
func (c *Contacts) NewIterator(opts ...ContactIterOption) (iter *ContactsIterator) {
	iter = &ContactsIterator{
		index:    -1,
		contacts: c,
		emails:   make(map[string]struct{}, 4),
	}

	for _, opt := range opts {
		opt(iter)
	}

	return iter
}

// Moves the ContactIterator to the next existing contact and returns true if there is
// one, false if the iteration is complete.
func (i *ContactsIterator) Next() bool {
	// Advance the index (note, index must be initialized at -1)
	i.index++

	// Stopping condition
	if i.index >= 4 {
		return false
	}

	// Filter checks - contact must exist
	current := i.contacts.Index(i.index)
	if current == nil || current.Contact.IsZero() {
		return i.Next()
	}

	// Filter if email required
	if i.skipNoEmail && !ContactHasEmail(current.Contact) {
		return i.Next()
	}

	// Filter if this is a duplicate
	_, ok := i.emails[current.Contact.Email]
	if i.skipDuplicates && ok {
		return i.Next()
	}

	// Filter if verified is required
	if i.skipUnverified && !current.Email.Verified {
		return i.Next()
	}

	// Filter if unverified is required
	if i.skipVerified && current.Email.Verified {
		return i.Next()
	}

	// Keep track of unique emails
	i.emails[current.Contact.Email] = struct{}{}
	return true
}

func (i *ContactsIterator) Contact() *ContactRecord {
	if i.index >= 0 && i.index < 4 {
		return i.contacts.Index(i.index)
	}
	return nil
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
func (c *Contacts) VerifiedContacts() (contacts map[string]string) {
	contacts = make(map[string]string)
	iter := c.NewIterator(SkipUnverified())
	for iter.Next() {
		contact := iter.Contact()
		contacts[contact.Kind] = contact.Email.Email
	}
	return contacts
}

// ContactVerifications returns a map of contact type to verified status, omitting any
// contacts that do not exist.
func (c *Contacts) ContactVerifications() (contacts map[string]bool) {
	contacts = make(map[string]bool)
	iter := c.NewIterator()
	for iter.Next() {
		contact := iter.Contact()
		contacts[contact.Kind] = contact.Email.Verified
	}
	return contacts
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
func AppendEmailLog(contact *pb.Contact, reason, subject string) (err error) {
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
		Recipient: contact.Email,
	}
	extra.EmailLog = append(extra.EmailLog, entry)

	// Serialize the extra data back to the VASP.
	if contact.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

// TODO: remove below since tests will no longer pass
type ContactIterator struct{}

func (c ContactIterator) Next() bool {
	return false
}

func (c ContactIterator) Value() (*pb.Contact, string) {
	return nil, ""
}

func (c ContactIterator) Error() error {
	return errors.New("no longer implemented")
}

func NewContactIterator(*pb.Contacts, ...ContactIterOption) *ContactIterator {
	return &ContactIterator{}
}
