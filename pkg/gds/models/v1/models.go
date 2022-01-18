package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/trisa/pkg/ivms101"
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

// NewCertificateRequest creates and returns a new certificate request to be associated with a VASP.
func NewCertificateRequest(vasp *pb.VASP) (certRequest *CertificateRequest, err error) {
	var (
		organizationName    string
		localityName        string
		stateOrProvinceName string
		countryName         string
	)

	if vasp == nil {
		return nil, errors.New("must supply a VASP object for certificate request creation")
	}

	certRequest = &CertificateRequest{
		Id:         uuid.New().String(),
		Vasp:       vasp.Id,
		CommonName: vasp.CommonName,
		Params:     make(map[string]string),
	}

	// Populate the organization name, if available.
	if vasp.Entity.Name != nil {
		for _, name := range vasp.Entity.Name.NameIdentifiers {
			if name.LegalPersonNameIdentifierType == ivms101.LegalPersonLegal {
				organizationName = name.LegalPersonName
				break
			}
		}
	}
	if organizationName != "" {
		certRequest.Params["organizationName"] = organizationName
	} else {
		log.Info().
			Str("vasp_id", vasp.Id).
			Str("certreq_id", certRequest.Id).
			Msg("organization name not found, populating new certificate request with default value")
		certRequest.Params["organizationName"] = "TRISA Production"
	}

	// Populate the location information, if available.
	if len(vasp.Entity.GeographicAddresses) > 0 {
		address := vasp.Entity.GeographicAddresses[0]
		localityName = address.TownLocationName
		stateOrProvinceName = address.CountrySubDivision
		countryName = address.Country
	}
	if localityName != "" && stateOrProvinceName != "" && countryName != "" {
		certRequest.Params["localityName"] = localityName
		certRequest.Params["stateOrProvinceName"] = stateOrProvinceName
		certRequest.Params["countryName"] = countryName
	} else {
		log.Info().
			Str("vasp_id", vasp.Id).
			Str("certreq_id", certRequest.Id).
			Msg("localtion information not found or incomplete, populating new certificate request with default values")
		certRequest.Params["localityName"] = "Menlo Park"
		certRequest.Params["stateOrProvinceName"] = "California"
		certRequest.Params["countryName"] = "US"
	}

	return certRequest, nil
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
