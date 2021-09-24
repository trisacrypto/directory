package models

import (
	"errors"
	"fmt"
	"strings"

	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/types/known/anypb"
)

// GetAdminVerificationToken from the extra data on the VASP record.
func GetAdminVerificationToken(vasp *pb.VASP) (_ string, err error) {
	// If the extra data is nil, return empty string with no error
	if vasp.Extra == nil {
		return "", nil
	}

	// Unmarshal the extra data field on the VASP
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return "", err
	}
	return extra.GetAdminVerificationToken(), nil
}

// SetAdminVerificationToken on the extra data on the VASP record.
func SetAdminVerificationToken(vasp *pb.VASP, token string) (err error) {
	// Must unmarshal previous extra to ensure that data besides the admin verification
	// token is not overwritten.
	extra := &GDSExtraData{}
	if vasp.Extra != nil {
		if err = vasp.Extra.UnmarshalTo(extra); err != nil {
			return fmt.Errorf("could not deserialize previous extra: %s", err)
		}
	}

	// Update the admin verification token
	extra.AdminVerificationToken = token

	// Serialize the extra back to the VASP.
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return err
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

// SetContactVerification token and verified status on the Contact record (completely
// replaces the old record, which may not be ideal).
func SetContactVerification(contact *pb.Contact, token string, verified bool) (err error) {
	if contact == nil {
		return errors.New("cannot set verification on nil contact")
	}

	extra := &GDSContactExtraData{
		Verified: verified,
		Token:    token,
	}
	if contact.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

// VerifiedContacts returns a map of contact type to email address for all verified
// contacts, omitting any contacts that are not verified or do not exist.
func VerifiedContacts(vasp *pb.VASP) (contacts map[string]string) {
	contacts = make(map[string]string)
	for key, verified := range ContactVerifications(vasp) {
		if verified {
			switch key {
			case "technical":
				contacts[key] = vasp.Contacts.Technical.Email
			case "administrative":
				contacts[key] = vasp.Contacts.Administrative.Email
			case "billing":
				contacts[key] = vasp.Contacts.Billing.Email
			case "legal":
				contacts[key] = vasp.Contacts.Legal.Email
			default:
				panic(fmt.Errorf("unknown contact type %q", key))
			}
		}
	}
	return contacts
}

// ContactVerifications returns a map of contact type to verified status, omitting any
// contacts that do not exist.
func ContactVerifications(vasp *pb.VASP) (contacts map[string]bool) {
	contacts = make(map[string]bool)
	pairs := []struct {
		key     string
		contact *pb.Contact
	}{
		{"technical", vasp.Contacts.Technical},
		{"administrative", vasp.Contacts.Administrative},
		{"billing", vasp.Contacts.Billing},
		{"legal", vasp.Contacts.Legal},
	}

	for _, pair := range pairs {
		if pair.contact != nil {
			_, contacts[pair.key], _ = GetContactVerification(pair.contact)
		}
	}
	return contacts
}

// IsTraveler returns true if the VASP common name ends in traveler.ciphertrace.com
func IsTraveler(vasp *pb.VASP) bool {
	return strings.HasSuffix(vasp.CommonName, "traveler.ciphertrace.com")
}
