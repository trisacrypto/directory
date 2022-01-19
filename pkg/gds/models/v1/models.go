package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	ErrorAlreadyExists = errors.New("already exists")
	ErrorNotFound      = errors.New("not found")
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

// UpdateVerificationStatus changes the verification state of a VASP and appends a new
// entry to the audit log on the extra.
func UpdateVerificationStatus(vasp *pb.VASP, state pb.VerificationState, description string, source string) error {
	// Append a new entry to the audit log.
	entry := &AuditLogEntry{
		Timestamp:     time.Now().Format(time.RFC3339),
		PreviousState: vasp.VerificationStatus,
		CurrentState:  state,
		Description:   description,
		Source:        source,
	}
	if err := AppendAuditLog(vasp, entry); err != nil {
		return err
	}

	// Set the new state on the VASP.
	vasp.VerificationStatus = state
	return nil
}

// GetAuditLog from the extra data on the VASP record.
func GetAuditLog(vasp *pb.VASP) (_ []*AuditLogEntry, err error) {
	// If the extra data is nil, return nil (no audit log).
	if vasp.Extra == nil {
		return nil, nil
	}

	// Unmarshal the extra data field on the VASP.
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return nil, err
	}
	return extra.GetAuditLog(), nil
}

// Append an AuditLogEntry to the extra data on the VASP record.
func AppendAuditLog(vasp *pb.VASP, entry *AuditLogEntry) (err error) {
	// Entry must be non-nil.
	if entry == nil {
		return errors.New("cannot append nil entry to AuditLog")
	}

	// Validate current state.
	if entry.CurrentState < 0 || entry.CurrentState > pb.VerificationState_ERRORED {
		return fmt.Errorf("cannot set verification state to unsupported value %d", entry.CurrentState)
	}

	// Unmarshal previous extra data.
	extra := &GDSExtraData{}
	if vasp.Extra != nil {
		if err = vasp.Extra.UnmarshalTo(extra); err != nil {
			return fmt.Errorf("could not deserialize previous extra: %s", err)
		}
	} else {
		extra.AuditLog = make([]*AuditLogEntry, 0, 1)
	}

	// Set the previous state for the new entry.
	if entry.PreviousState == 0 && len(extra.AuditLog) > 0 {
		entry.PreviousState = extra.AuditLog[len(extra.AuditLog)-1].CurrentState
	}

	// Append entry to the previous log.
	extra.AuditLog = append(extra.AuditLog, entry)

	// Serialize the extra data back to the VASP.
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

// GetCertReqIDs returns the list of associated CertificateRequest IDs for the VASP record.
func GetCertReqIDs(vasp *pb.VASP) (_ []string, err error) {
	// If the extra data is nil, return nil (no certificate requests).
	if vasp.Extra == nil {
		return nil, nil
	}

	// Unmarshal the extra data field on the VASP.
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return nil, err
	}
	return extra.GetCertificateRequests(), nil
}

// AppendCertReqID adds the certificate request ID to the VASP if its not already added.
func AppendCertReqID(vasp *pb.VASP, certreqID string) (err error) {
	// Entry must be non-nil.
	if certreqID == "" {
		return errors.New("cannot append empty certificate request ID to extra")
	}

	// Unmarshal previous extra data.
	extra := &GDSExtraData{}
	if vasp.Extra != nil {
		if err = vasp.Extra.UnmarshalTo(extra); err != nil {
			return fmt.Errorf("could not deserialize previous extra: %s", err)
		}
	} else {
		extra.CertificateRequests = make([]string, 0, 1)
	}

	// Do not allow duplicate certificate requests to be appended
	for _, containsID := range extra.CertificateRequests {
		if certreqID == containsID {
			// Do not return an error
			return nil
		}
	}

	// Append certificate request ID to the array.
	extra.CertificateRequests = append(extra.CertificateRequests, certreqID)

	// Serialize the extra data back to the VASP.
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

// UpdateCertificateRequestStatus changes the status of a CertificateRequest and appends
// an entry to the audit log.
func UpdateCertificateRequestStatus(request *CertificateRequest, state CertificateRequestState, description string, source string) (err error) {
	// CertificateRequest must be non-nil
	if request == nil {
		return fmt.Errorf("cannot set certificate request status on a nil CertificateRequest")
	}

	// Validate the new state.
	if state < 0 || state > CertificateRequestState_CR_ERRORED {
		return fmt.Errorf("cannot set certificate request status to unsupported value %d", state)
	}

	// Append a new entry to the audit log.
	entry := &CertificateRequestLogEntry{
		Timestamp:     time.Now().Format(time.RFC3339),
		PreviousState: request.Status,
		CurrentState:  state,
		Description:   description,
		Source:        source,
	}
	request.AuditLog = append(request.AuditLog, entry)

	// Set the new state on the CertificateRequest.
	request.Status = state
	return nil
}

// GetReviewNotes returns all of the review notes for a VASP as a map.
func GetReviewNotes(vasp *pb.VASP) (_ map[string]*ReviewNote, err error) {
	// If the extra data is nil, return an empty map (no review notes).
	if vasp.Extra == nil {
		return map[string]*ReviewNote{}, nil
	}

	// Unmarshal the extra data field on the VASP.
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return nil, err
	}
	return extra.ReviewNotes, nil
}

// CreateReviewNote creates a note on the VASP given the specified note id.
func CreateReviewNote(vasp *pb.VASP, id string, author string, text string) (note *ReviewNote, err error) {
	// Validate note id.
	if id == "" {
		return nil, errors.New("must specify a valid note id")
	}

	// Unmarshal previous extra data.
	extra := &GDSExtraData{}
	if vasp.Extra != nil {
		if err = vasp.Extra.UnmarshalTo(extra); err != nil {
			return nil, fmt.Errorf("could not deserialize previous extra: %s", err)
		}
	}

	if extra.ReviewNotes == nil {
		extra.ReviewNotes = make(map[string]*ReviewNote)
	}

	if _, exists := extra.ReviewNotes[id]; exists {
		return nil, ErrorAlreadyExists
	}

	// Create the new note.
	note = &ReviewNote{
		Id:      id,
		Created: time.Now().Format(time.RFC3339),
		Author:  author,
		Text:    text,
	}
	extra.ReviewNotes[id] = note

	// Serialize the extra data back to the VASP.
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return nil, err
	}

	return note, nil
}

// UpdateReviewNote updates a specified note on the VASP.
func UpdateReviewNote(vasp *pb.VASP, id string, editor string, text string) (note *ReviewNote, err error) {
	// Validate note id.
	if id == "" {
		return nil, errors.New("must specify a valid note id")
	}

	// Update is invalid if the extra data doesn't exist.
	if vasp.Extra == nil {
		return nil, errors.New("extra does not exist")
	}

	// Unmarshal previous extra data.
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return nil, fmt.Errorf("could not deserialize previous extra: %s", err)
	}

	// Get the specified note.
	var exists bool
	if note, exists = extra.ReviewNotes[id]; !exists {
		return nil, ErrorNotFound
	}

	// Update the note.
	note.Modified = time.Now().Format(time.RFC3339)
	note.Editor = editor
	note.Text = text

	// Serialize the extra data back to the VASP.
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return nil, err
	}

	return note, nil
}

// DeleteReviewNote deletes a specified note on the VASP.
func DeleteReviewNote(vasp *pb.VASP, id string) (err error) {
	// Validate note id.
	if id == "" {
		return errors.New("must specify a valid note id")
	}

	// Delete is invalid if the extra data doesn't exist.
	if vasp.Extra == nil {
		return errors.New("extra does not exist")
	}

	// Unmarshal previous extra data.
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return fmt.Errorf("could not deserialize previous extra: %s", err)
	}

	if _, exists := extra.ReviewNotes[id]; !exists {
		return ErrorNotFound
	}

	// Delete the note
	delete(extra.ReviewNotes, id)

	// Serialize the extra data back to the VASP.
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return err
	}

	return nil
}

// VerifiedContacts returns a map of contact type to email address for all verified
// contacts, omitting any contacts that are not verified or do not exist.
func VerifiedContacts(vasp *pb.VASP) (contacts map[string]string, err error) {
	contacts = make(map[string]string)
	next := IterContacts(vasp.Contacts, true)
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
	contact, kind, err := next()
	for ; err == nil && contact != nil; contact, kind, err = next() {
		var verified bool
		if verified, err = ContactVerified(contact); err != nil {
			return nil, err
		}
		contacts[kind] = verified
	}
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

// IsTraveler returns true if the VASP common name ends in traveler.ciphertrace.com
func IsTraveler(vasp *pb.VASP) bool {
	return strings.HasSuffix(vasp.CommonName, "traveler.ciphertrace.com")
}
