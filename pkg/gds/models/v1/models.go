package models

import (
	"fmt"
	"time"

	"github.com/trisacrypto/directory/pkg/gds/global/v1"
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

// GetMetadata returns the object metadata and deleted on timestamp from the VASP record.
func GetMetadata(vasp *pb.VASP) (_ *global.Object, deletedOn time.Time, err error) {
	if vasp.Extra == nil {
		return nil, deletedOn, nil
	}

	// Unmarshal the extra data field on the VASP
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return nil, deletedOn, err
	}

	// Parse the deleted on timestamp
	if extra.DeletedOn != "" {
		if deletedOn, err = time.Parse(time.RFC3339, extra.DeletedOn); err != nil {
			return nil, deletedOn, fmt.Errorf("could not parse deleted on timestamp: %s", err)
		}
	}

	// Return the metadata
	return extra.Metadata, deletedOn, nil
}

// SetMetadata updates the VASP record with the new object metadata and deleted on
// timestamp. If the record is nil or the timestamp is zero, then it will be set to nil
// or zero on the extra record, overwriting any previous value.
func SetMetadata(vasp *pb.VASP, metadata *global.Object, deletedOn time.Time) (err error) {
	// Must unmarshal previous extra to ensure that data besides the object metadata is
	// not overwritten (such as the admin verification token).
	extra := &GDSExtraData{}
	if vasp.Extra != nil {
		if err = vasp.Extra.UnmarshalTo(extra); err != nil {
			return fmt.Errorf("could not deserialize previous extra: %s", err)
		}
	}

	// Update the extra record with new data
	extra.Metadata = metadata

	// If we don't do the iszero check, then epoch time will be written to string.
	if deletedOn.IsZero() {
		extra.DeletedOn = ""
	} else {
		extra.DeletedOn = deletedOn.Format(time.RFC3339)
	}

	// Serialize the extra back to the VASP.
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
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
func SetContactVerification(contact *pb.Contact, token string, verified bool) (err error) {
	extra := &GDSContactExtraData{
		Verified: verified,
		Token:    token,
	}
	if contact.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}
