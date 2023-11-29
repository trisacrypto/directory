package models

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
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

// GetAdminEmailLog from the extra data on the VASP record.
func GetAdminEmailLog(vasp *pb.VASP) (_ []*EmailLogEntry, err error) {
	// If the extra data is nil, return nil (no email log).
	if vasp.Extra == nil {
		return nil, nil
	}

	// Unmarshal the extra data field on the VASP.
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return nil, err
	}
	return extra.GetEmailLog(), nil
}

// Create and add a new entry to the EmailLog on the extra data on the VASP record.
func AppendAdminEmailLog(vasp *pb.VASP, reason string, subject string) (err error) {
	// VASP must be non-nil.
	if vasp == nil {
		return errors.New("cannot append to nil VASP")
	}

	// Unmarshal previous extra data.
	extra := &GDSExtraData{}
	if vasp.Extra != nil {
		if err = vasp.Extra.UnmarshalTo(extra); err != nil {
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
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

// GetVASPEmailLog computes a unified email log which is a time ordered representation
// of all the emails sent to a VASP's contacts.
func GetVASPEmailLog(vasp *pb.VASP) (emails []*EmailLogEntry, err error) {
	emails = make([]*EmailLogEntry, 0)

	// Iterate over all the contacts on the VASP, this skips any nil contacts.
	iter := NewContactIterator(vasp.Contacts)
	for iter.Next() {
		contact, _ := iter.Value()
		var contactLog []*EmailLogEntry
		if contactLog, err = GetEmailLog(contact); err != nil {
			return nil, err
		}

		// Merge the contact log into the unified log while maintaining timestamp order
		i, j := 0, 0
		for i < len(emails) && j < len(contactLog) {
			var logTime, contactTime time.Time
			if logTime, err = time.Parse(time.RFC3339, emails[i].Timestamp); err != nil {
				return nil, err
			}

			if contactTime, err = time.Parse(time.RFC3339, contactLog[j].Timestamp); err != nil {
				return nil, err
			}

			if contactTime.Before(logTime) {
				emails = append(emails, nil)
				copy(emails[i+1:], emails[i:])
				emails[i] = contactLog[j]
				j++
			} else {
				i++
			}
		}

		// Append any remaining entries
		for j < len(contactLog) {
			emails = append(emails, contactLog[j])
			j++
		}
	}

	return emails, nil
}

// Create and add a new entry to the EmailLog on the extra data on the Contact record.
func (c *Contact) AppendEmailLog(reason, subject string) {
	// Contact must be non-nil.
	if c == nil {
		return
	}

	// Create the EmailLog if it is nil.
	if c.EmailLog == nil {
		c.EmailLog = make([]*EmailLogEntry, 0, 1)
	}

	// Append entry to the previous log.
	entry := &EmailLogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Reason:    reason,
		Subject:   subject,
		Recipient: c.Email,
	}
	c.EmailLog = append(c.EmailLog, entry)
}

// Normalize the email and convert to bytes
func NormalizeEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	normalized := strings.ToLower(trimmed)
	return normalized
}

// Counts emails within the given EmailLogEntry slice for the given reason within the given time frame.
func CountSentEmails(emailLog []*EmailLogEntry, reason string, timeWindowDays int) (sent int, err error) {
	if reason == "" {
		return 0, errors.New("cannot match on empty reason string")
	}
	if timeWindowDays < 0 {
		return 0, errors.New("time window must be a positive number of days")
	}

	for _, value := range emailLog {
		var timestamp time.Time
		if timestamp, err = time.Parse(time.RFC3339, value.Timestamp); err != nil {
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

// GetLatestCertReqID returns the latest CertificateRequest ID for the VASP record.
// TODO: This relies on the assumption that IDs can only be appended to the list.
func GetLatestCertReqID(vasp *pb.VASP) (_ string, err error) {
	// If the extra data is nil, return nil (no certificate requests).
	if vasp.Extra == nil {
		return "", nil
	}

	// Unmarshal the extra data field on the VASP.
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return "", err
	}

	// Return nil if there are no certificate requests.
	if len(extra.CertificateRequests) == 0 {
		return "", nil
	}

	// Get the latest certificate request ID
	return extra.CertificateRequests[len(extra.CertificateRequests)-1], nil
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

// DeleteCertReqID removes the certificate request ID from the VASP if it exists.
func DeleteCertReqID(vasp *pb.VASP, certreqID string) (err error) {
	// ID is required
	if certreqID == "" {
		return errors.New("cannot delete empty certificate request ID from extra")
	}

	// Unmarshal previous extra data.
	extra := &GDSExtraData{}
	if vasp.Extra != nil {
		if err = vasp.Extra.UnmarshalTo(extra); err != nil {
			return fmt.Errorf("could not deserialize previous extra: %s", err)
		}
	} else {
		extra.CertificateRequests = make([]string, 0)
	}

	// Search the slice for the certificate request ID
	for i, containsID := range extra.CertificateRequests {
		if certreqID == containsID {
			// Remove the certificate request ID from the array
			extra.CertificateRequests = append(extra.CertificateRequests[:i], extra.CertificateRequests[i+1:]...)
			break
		}
	}

	// Serialize the extra data back to the VASP.
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

// GetCertIDs returns the list of associated Certificate IDs for the VASP record.
func GetCertIDs(vasp *pb.VASP) (_ []string, err error) {
	// If the extra data is nil, return nil (no certificates).
	if vasp.Extra == nil {
		return nil, nil
	}

	// Unmarshal the extra data field on the VASP.
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return nil, err
	}
	return extra.GetCertificates(), nil
}

// AppendCertID adds the certificate ID to the VASP if it's not already added.
func AppendCertID(vasp *pb.VASP, certID string) (err error) {
	// Entry must be non-nil.
	if certID == "" {
		return errors.New("cannot append empty certificate ID to extra")
	}

	// Unmarshal previous extra data.
	extra := &GDSExtraData{}
	if vasp.Extra != nil {
		if err = vasp.Extra.UnmarshalTo(extra); err != nil {
			return fmt.Errorf("could not deserialize previous extra: %s", err)
		}
	} else {
		extra.Certificates = make([]string, 0, 1)
	}

	// Do not allow duplicate certificates to be appended.
	for _, containsID := range extra.Certificates {
		if certID == containsID {
			// Do not return an error.
			return nil
		}
	}

	// Append certificate ID to the slice.
	extra.Certificates = append(extra.Certificates, certID)

	// Serialize the extra data back to the VASP.
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

// NewCertificate creates and returns a certificate associated with a VASP.
func NewCertificate(vasp *pb.VASP, certRequest *CertificateRequest, data *pb.Certificate) (cert *Certificate, err error) {
	// VASP must be not nil.
	if vasp == nil {
		return nil, errors.New("must supply a VASP object for certificate creation")
	}

	// Certificate request must be not nil.
	if certRequest == nil {
		return nil, errors.New("must supply a certificate request for certificate creation")
	}

	// Certificate data must be not nil.
	if data == nil {
		return nil, errors.New("must supply certificate data for certificate creation")
	}

	cert = &Certificate{
		Id:      GetCertID(data), // capital hex encoded serial number to match sectigo
		Request: certRequest.Id,
		Vasp:    vasp.Id,
		Status:  CertificateState_ISSUED,
		Details: data,
	}

	return cert, nil
}

// String encode the SerialNumber of a certificate
func GetCertID(identityCert *pb.Certificate) string {
	return fmt.Sprintf("%X", identityCert.SerialNumber)
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
		Id:              uuid.New().String(),
		Vasp:            vasp.Id,
		CommonName:      vasp.CommonName,
		Params:          make(map[string]string),
		Webhook:         vasp.CertificateWebhook,
		NoEmailDelivery: vasp.NoEmailDelivery,
	}

	// Populate the organization name, if available.
	if organizationName, err = vasp.Name(); err == nil {
		UpdateCertificateRequestParams(certRequest, sectigo.ParamOrganizationName, organizationName)
	} else {
		sentry.Warn(nil).Err(err).Str("vasp_id", vasp.Id).Str("certreq_id", certRequest.Id).Msg("organization name not found")
	}

	// Populate the location information, if available.
	if len(vasp.Entity.GeographicAddresses) > 0 {
		address := vasp.Entity.GeographicAddresses[0]
		localityName = address.TownLocationName
		stateOrProvinceName = address.CountrySubDivision
		countryName = address.Country
	}
	if localityName != "" && stateOrProvinceName != "" && countryName != "" {
		UpdateCertificateRequestParams(certRequest, sectigo.ParamLocalityName, localityName)
		UpdateCertificateRequestParams(certRequest, sectigo.ParamStateOrProvinceName, stateOrProvinceName)
		UpdateCertificateRequestParams(certRequest, sectigo.ParamCountryName, countryName)
	} else {
		log.Debug().Str("vasp_id", vasp.Id).Str("certreq_id", certRequest.Id).Msg("location information not found or incomplete")
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

// IsTraveler returns true if the VASP common name ends in traveler.ciphertrace.com
func IsTraveler(vasp *pb.VASP) bool {
	return strings.HasSuffix(vasp.CommonName, "traveler.ciphertrace.com")
}

// ValidateVASP checks if a VASP model is valid using the method from the TRISA
// codebase and returns an error if it is not. If partial is true, "partial" validation
// is performed which allows for some fields such as timestamps and IDs to be empty.
func ValidateVASP(vasp *pb.VASP, partial bool) (err error) {
	err = vasp.Validate(partial)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ivms101.ErrCompleteNationalIdentifierCountry):
		// TODO: ErrCompleteNationalIdentifierLegalPerson must be ignored to support older
		// VASP records, see issue #34
		sentry.Warn(nil).Str("vasp_id", vasp.Id).Err(err).Msg("ignoring ErrCompleteNationalIdentifierLegalPerson validation error")
		return nil
	default:
		return err
	}
}

// VASPSignature returns a SHA512 hash of the protocol buffer bytes that represents the
// current data in the VASP. This signature can be used to compare two VASP records for
// equality or to detect if a VASP record has been changed.
func VASPSignature(vasp *pb.VASP) (_ []byte, err error) {
	var data []byte
	if data, err = proto.Marshal(vasp); err != nil {
		return nil, err
	}

	sum := sha512.Sum512(data)
	return sum[:], nil
}
