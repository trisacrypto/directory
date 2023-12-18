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

// Returns true if the contact kind is one of the recognized strings.
func ContactKindIsValid(kind string) bool {
	return kind == TechnicalContact || kind == AdministrativeContact || kind == LegalContact || kind == BillingContact
}

//=====================================================================================
// Contacts and ContactRecords
//=====================================================================================

// Contacts wraps a VASPs contacts with their email records for easier access to
// contact records, including iteration over contact records. This type of record is
// often created from database records, and acts as a join record between the VASPs
// table and the emails table to prevent duplicate emails from being sent.
type Contacts struct {
	VASP     string
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

// Returns true if the contact record has an email and it has been verified.
func (c *ContactRecord) IsVerified() bool {
	return c.Email != nil && c.Email.Verified
}

// Return true if the contact is nil or is completely zero-valued.
func (c *ContactRecord) IsZero() bool {
	return c.Contact == nil || c.Contact.IsZero()
}

// Returns true if a Contact is not nil and has an email address.
func (c *ContactRecord) HasEmail() bool {
	return c.Contact != nil && c.Contact.Email != ""
}

// Returns the sent email logs if they're on the contact email address otherwise nil.
func (c *ContactRecord) Logs() []*EmailLogEntry {
	if c.Email != nil {
		return c.Email.SendLog
	}
	return nil
}

// Update a contact record with the other contact data and return a new email address
// if the email address has changed along with a bool indicating if it has to be saved.
func (c *ContactRecord) Update(contact *pb.Contact) (*Email, bool) {
	contactEmail := NormalizeEmail(contact.Email)
	c.Contact = contact
	if c.Email == nil {
		return &Email{Name: contact.Name, Email: contactEmail}, true
	}

	if c.Email.Email != contactEmail {
		if contact.Name != "" {
			c.Email.Name = contact.Name
		}

		c.Email.Email = contactEmail
		return c.Email, true
	}

	return c.Email, false
}

//=====================================================================================
// Contacts Methods
//=====================================================================================

// Has returns true if the conact for the specified kind is not nil or zero.
func (c *Contacts) Has(kind string) bool {
	switch kind {
	case TechnicalContact:
		return c.Contacts.Technical != nil && !c.Contacts.Technical.IsZero()
	case AdministrativeContact:
		return c.Contacts.Administrative != nil && !c.Contacts.Administrative.IsZero()
	case LegalContact:
		return c.Contacts.Legal != nil && !c.Contacts.Legal.IsZero()
	case BillingContact:
		return c.Contacts.Billing != nil && !c.Contacts.Billing.IsZero()
	default:
		panic(fmt.Errorf("invalid contact kind %q", kind))
	}
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

	// Return a nil record if the VASP contact is nil but returns a zero-valued contact.
	if record.Contact == nil {
		return nil
	}

	// Find the email for the given contact
	contactEmail := NormalizeEmail(record.Contact.Email)
	for _, email := range c.Emails {
		if email.Email == contactEmail {
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

// Adds a contact on the VASP object. Returns the email address record associated with
// the contact and returns true if the email address needs to be created or updated.
func (c *Contacts) Add(kind string, contact *pb.Contact) (*Email, bool) {
	switch kind {
	case TechnicalContact:
		c.Contacts.Technical = contact
	case AdministrativeContact:
		c.Contacts.Administrative = contact
	case LegalContact:
		c.Contacts.Legal = contact
	case BillingContact:
		c.Contacts.Billing = contact
	default:
		panic(fmt.Errorf("invalid contact kind %q", kind))
	}

	// Determine if another contact already has this email; if so we don't have to
	// create or update the email since it was already stored for another contact.
	contactEmail := NormalizeEmail(contact.Email)
	for _, email := range c.Emails {
		if email.Email == contactEmail {
			return email, false
		}
	}

	// If we didn't find the email on the contacts, then it needs to be created and a
	// verification email sent or updated with the vaspID or the specified record.
	return &Email{Name: contact.Name, Email: contactEmail}, true
}

// Deletes a contact on the VASP object by setting it to nil. Returns the email address
// record associated with the contact and true if the email address needs to be saved.
func (c *Contacts) Delete(kind string) (*Email, bool) {
	contact := c.Get(kind)
	if contact == nil {
		return nil, false
	}

	switch kind {
	case TechnicalContact:
		c.Contacts.Technical = nil
	case AdministrativeContact:
		c.Contacts.Administrative = nil
	case LegalContact:
		c.Contacts.Legal = nil
	case BillingContact:
		c.Contacts.Billing = nil
	default:
		panic(fmt.Errorf("invalid contact kind %q", kind))
	}

	// Check if any of the other contacts have the same email
	found := false
	iter := c.NewIterator(SkipNoEmail())
	for iter.Next() {
		other := iter.Contact()
		if other.Kind != kind && other.Email.Email == contact.Contact.Email {
			found = true
		}
	}

	// If we did not find another contact with the same email, we should remove the
	// vaspID from the email address and return true that it needs to be saved.
	if !found {
		contact.Email.RmVASP(c.VASP)
		return contact.Email, true
	}

	return nil, false
}

// IsVerified returns true if the contact kind is not nil or zero and has a verified
// email address record in the database.
func (c *Contacts) IsVerified(kind string) bool {
	contact := c.Get(kind)
	if contact == nil {
		return false
	}
	return contact.Email != nil && contact.Email.Verified
}

// Logs returns the email logs for the specified contact kind. If the specified contact
// is not found or the contact does not have an email then an empty log is returned.
func (c *Contacts) Logs(kind string) []*EmailLogEntry {
	var contact *ContactRecord
	if contact = c.Get(kind); contact == nil {
		return nil
	}
	return contact.Logs()
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

// Length returns the number of contact records that would be returned from the iterator
// given the specified iteration constraints. This method can be used to count the
// number of verified or unverified contacts, the number of unduplicated contacts, etc.
func (c *Contacts) Length(opts ...ContactIterOption) int {
	n := 0
	iter := c.NewIterator(opts...)
	for iter.Next() {
		n++
	}
	return n
}

//=====================================================================================
// ContactsIterator
//=====================================================================================

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
	if i.skipNoEmail && !current.HasEmail() {
		return i.Next()
	}

	// Filter if this is a duplicate
	_, ok := i.emails[current.Contact.Email]
	if i.skipDuplicates && ok {
		return i.Next()
	}

	// Filter if verified is required
	if i.skipUnverified && !current.IsVerified() {
		return i.Next()
	}

	// Filter if unverified is required
	if i.skipVerified && current.IsVerified() {
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

//=====================================================================================
// Deprecrated
//=====================================================================================

// GetContactVerification token and verified status from the extra data field on the Contact.
//
// Deprecated: Use the emails model to manage email verification state.
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
//
// Deprecated: Use the emails model to manage email verification state.
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

// Create and add a new entry to the EmailLog on the extra data on the Contact record.
//
// Deprecated: Use the emails model to manage email logs.
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
