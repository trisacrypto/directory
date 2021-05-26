package models

import (
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
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

// SetAdminVerificationToken on the extra data on the VASP record (completely replaces
// the old record, which may not be ideal).
func SetAdminVerificationToken(vasp *pb.VASP, token string) error {
	extra := &GDSExtraData{
		AdminVerificationToken: token,
	}
	return vasp.Extra.MarshalFrom(extra)
}

// GetContactVerification token and verified status from the extra data field on the Contact.
func GetContactVerification(contact *pb.Contact) (_ string, _ bool, err error) {
	// Return zero-valued defaults with no error if extra is nil.
	if contact.Extra == nil {
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
func SetContactVerification(contact *pb.Contact, token string, verified bool) error {
	extra := &GDSContactExtraData{
		Verified: verified,
		Token:    token,
	}
	return contact.Extra.MarshalFrom(extra)
}
