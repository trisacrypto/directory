package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/proto"
)

func TestVASPExtra(t *testing.T) {
	// Ensures that the Get/Set methods on the VASP do not overwrite values other than
	// the values the method is intended to interact with.
	vasp := &pb.VASP{}

	// Attempt to get an admin verification token on a nil extra
	token, err := GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "", token)

	// Attempt to set admin verification token on a nil extra
	err = SetAdminVerificationToken(vasp, "pontoonboatz")
	require.NoError(t, err)

	// Should be able to fetch the token
	token, err = GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "pontoonboatz", token)

	// Attempt to append to audit log
	entry := &AuditLogEntry{
		Timestamp:    time.Now().Format(time.RFC3339),
		CurrentState: pb.VerificationState_VERIFIED,
		Description:  "description",
		Source:       "pontoon@boatz.com",
	}
	err = AppendAuditLog(vasp, entry)
	require.NoError(t, err)

	// Should be able to fetch the new audit log
	auditLog, err := GetAuditLog(vasp)
	require.NoError(t, err)
	require.Len(t, auditLog, 1)
	require.True(t, proto.Equal(entry, auditLog[0]))

	// Attempt to create a review note
	note, err := CreateReviewNote(vasp, "boats", "pontoon@boatz.com", "boats are cool")
	require.NoError(t, err)
	require.Equal(t, "boats", note.Id)
	require.Equal(t, "", note.Modified)
	require.Equal(t, "pontoon@boatz.com", note.Author)
	require.Equal(t, "", note.Editor)
	require.Equal(t, "boats are cool", note.Text)

	// Should be able to fetch the note
	notes, err := GetReviewNotes(vasp)
	require.NoError(t, err)
	require.Len(t, notes, 1)
	require.Equal(t, "boats", notes["boats"].Id)
	require.Equal(t, "", notes["boats"].Modified)
	require.Equal(t, "pontoon@boatz.com", notes["boats"].Author)
	require.Equal(t, "", notes["boats"].Editor)
	require.Equal(t, "boats are cool", notes["boats"].Text)

	// Verification token should be unchanged
	token, err = GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "pontoonboatz", token)

	// Attempt to set a new verification token
	err = SetAdminVerificationToken(vasp, "jetskis")
	require.NoError(t, err)

	// Verify that the new token was set
	token, err = GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "jetskis", token)

	// Audit log should be unchanged
	auditLog, err = GetAuditLog(vasp)
	require.NoError(t, err)
	require.Len(t, auditLog, 1)
	require.True(t, proto.Equal(entry, auditLog[0]))

	// Review notes should be unchanged
	notes, err = GetReviewNotes(vasp)
	require.NoError(t, err)
	require.Len(t, notes, 1)
	require.Equal(t, "boats", notes["boats"].Id)
	require.Equal(t, "", notes["boats"].Modified)
	require.Equal(t, "pontoon@boatz.com", notes["boats"].Author)
	require.Equal(t, "", notes["boats"].Editor)
	require.Equal(t, "boats are cool", notes["boats"].Text)

	// Attempt to append certificate request IDs
	for _, cfid := range []string{"b5841869-105f-411c-8722-4045aad72717", "230d5e77-9983-4f1f-80ea-d379d56519af"} {
		err = AppendCertReqID(vasp, cfid)
		require.NoError(t, err)
	}

	// Should be able to fetch the certificate request IDs
	ids, err := GetCertReqIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 2)

	// Verification token should be unchanged
	token, err = GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "jetskis", token)
}

func TestAuditLog(t *testing.T) {
	// Test that the audit log functions are working as expected
	vasp := &pb.VASP{}

	// Audit log should initially be empty
	auditLog, err := GetAuditLog(vasp)
	require.NoError(t, err)
	require.Len(t, auditLog, 0)

	// Should not be able to append a nil entry to the audit log
	err = AppendAuditLog(vasp, nil)
	require.Error(t, err)

	// Should not be able to append an invalid state to the audit log
	entry := &AuditLogEntry{
		Timestamp:    time.Now().Format(time.RFC3339),
		CurrentState: -1,
		Description:  "invalid state",
		Source:       "pontoon@boatz.com",
	}
	err = AppendAuditLog(vasp, entry)
	require.Error(t, err)
	entry.CurrentState = pb.VerificationState_ERRORED + 1
	err = AppendAuditLog(vasp, entry)
	require.Error(t, err)

	// Append an entry to an empty log
	entry = &AuditLogEntry{
		Timestamp:    time.Now().Format(time.RFC3339),
		CurrentState: pb.VerificationState_SUBMITTED,
		Description:  "description",
		Source:       "automated",
	}
	err = AppendAuditLog(vasp, entry)
	require.NoError(t, err)
	auditLog, err = GetAuditLog(vasp)
	require.NoError(t, err)
	require.Len(t, auditLog, 1)
	require.True(t, proto.Equal(entry, auditLog[0]))

	// Append an entry without specifying PreviousState
	entry2 := &AuditLogEntry{
		Timestamp:    time.Now().Format(time.RFC3339),
		CurrentState: pb.VerificationState_EMAIL_VERIFIED,
		Description:  "sent verification emails",
		Source:       "pontoon@boatz.com",
	}
	expected2 := entry2
	expected2.PreviousState = entry.CurrentState
	err = AppendAuditLog(vasp, entry2)
	require.NoError(t, err)
	auditLog, err = GetAuditLog(vasp)
	require.NoError(t, err)
	require.Len(t, auditLog, 2)
	require.True(t, proto.Equal(entry, auditLog[0]))
	require.True(t, proto.Equal(expected2, auditLog[1]))

	// Append an entry, specifying the PreviousState
	entry3 := &AuditLogEntry{
		Timestamp:     time.Now().Format(time.RFC3339),
		PreviousState: pb.VerificationState_APPEALED,
		CurrentState:  pb.VerificationState_REJECTED,
		Description:   "appeal rejected",
		Source:        "admin@example.com",
	}
	err = AppendAuditLog(vasp, entry3)
	require.NoError(t, err)
	auditLog, err = GetAuditLog(vasp)
	require.NoError(t, err)
	require.Len(t, auditLog, 3)
	require.True(t, proto.Equal(entry, auditLog[0]))
	require.True(t, proto.Equal(expected2, auditLog[1]))
	require.True(t, proto.Equal(entry3, auditLog[2]))
}

func TestReviewNotes(t *testing.T) {
	// Test review note operations (retrieve, create, update, delete)
	vasp := &pb.VASP{}

	// Should initially be no review notes
	notes, err := GetReviewNotes(vasp)
	require.NoError(t, err)
	require.Len(t, notes, 0)

	// Attempt to update a note from an empty map
	_, err = UpdateReviewNote(vasp, "boats", "pontoon@boatz.com", "boats are cool")
	require.Error(t, err)

	// Attempt to update a note from an empty map
	err = DeleteReviewNote(vasp, "boats")
	require.Error(t, err)

	// Create a new note
	note, err := CreateReviewNote(vasp, "boats", "pontoon@boatz.com", "boats are cool")
	require.NoError(t, err)
	require.Equal(t, "boats", note.Id)
	require.Equal(t, "", note.Modified)
	require.Equal(t, "pontoon@boatz.com", note.Author)
	require.Equal(t, "", note.Editor)
	require.Equal(t, "boats are cool", note.Text)

	notes, err = GetReviewNotes(vasp)
	require.NoError(t, err)
	require.Len(t, notes, 1)
	require.Equal(t, "boats", notes["boats"].Id)
	require.Equal(t, "", notes["boats"].Modified)
	require.Equal(t, "pontoon@boatz.com", notes["boats"].Author)
	require.Equal(t, "", notes["boats"].Editor)
	require.Equal(t, "boats are cool", notes["boats"].Text)

	// Attempt to update a note that doesn't exist
	_, err = UpdateReviewNote(vasp, "jetskis", "admin@example.com", "jetskis are fun")
	require.Error(t, err)

	// Attempt to delete a note that doesn't exist
	err = DeleteReviewNote(vasp, "jetskis")
	require.Error(t, err)

	// Create a new note
	note, err = CreateReviewNote(vasp, "jetskis", "admin@example.com", "jetskis are fun")
	require.NoError(t, err)
	require.Equal(t, "jetskis are fun", note.Text)

	notes, err = GetReviewNotes(vasp)
	require.NoError(t, err)
	require.Len(t, notes, 2)
	require.Equal(t, "boats are cool", notes["boats"].Text)
	require.Equal(t, "jetskis are fun", notes["jetskis"].Text)

	// Update an existing note
	expectedTime := time.Now()
	note, err = UpdateReviewNote(vasp, "jetskis", "pontoon@boatz.com", "jetskis are loud")
	require.NoError(t, err)
	require.Equal(t, "jetskis", note.Id)
	require.Equal(t, "jetskis are loud", note.Text)
	require.Equal(t, "admin@example.com", note.Author)
	require.Equal(t, "pontoon@boatz.com", note.Editor)
	modifiedTime, err := time.Parse(time.RFC3339, note.Modified)
	require.NoError(t, err)
	require.LessOrEqual(t, modifiedTime.Sub(expectedTime), time.Minute)

	// Editor and modified should be updated
	notes, err = GetReviewNotes(vasp)
	require.NoError(t, err)
	require.Len(t, notes, 2)
	require.Equal(t, "boats are cool", notes["boats"].Text)
	require.Equal(t, "jetskis are loud", notes["jetskis"].Text)
	require.Equal(t, "pontoon@boatz.com", notes["jetskis"].Editor)
	require.LessOrEqual(t, modifiedTime.Sub(expectedTime), time.Minute)

	// Delete an existing note
	err = DeleteReviewNote(vasp, "boats")
	require.NoError(t, err)

	notes, err = GetReviewNotes(vasp)
	require.NoError(t, err)
	require.Len(t, notes, 1)
	require.Equal(t, "jetskis are loud", notes["jetskis"].Text)
}

func TestCertReqIDs(t *testing.T) {
	vasp := &pb.VASP{}

	// No extra, Get should return nil
	ids, err := GetCertReqIDs(vasp)
	require.NoError(t, err)
	require.Empty(t, ids)

	// Cannot append an empty ID
	err = AppendCertReqID(vasp, "")
	require.EqualError(t, err, "cannot append empty certificate request ID to extra")

	// Append an ID from empty
	err = AppendCertReqID(vasp, "1df61840-7033-40fb-8ce9-538c87e242f5")
	require.NoError(t, err)
	ids, err = GetCertReqIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 1)

	// Append an ID with data already inside
	err = AppendCertReqID(vasp, "9676bf6a-ffdb-4185-8fa5-87cdae6f6eef")
	require.NoError(t, err)
	ids, err = GetCertReqIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 2)

	// Do not allow duplicate IDs
	err = AppendCertReqID(vasp, "9676bf6a-ffdb-4185-8fa5-87cdae6f6eef")
	require.NoError(t, err)
	ids, err = GetCertReqIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 2)
}

func TestContactExtra(t *testing.T) {
	// Test contact is nil
	token, verified, err := GetContactVerification(nil)
	require.NoError(t, err, "nil contact returns error")
	require.False(t, verified)
	require.Empty(t, token)

	// Test extra is nil
	contact := &pb.Contact{
		Email: "pontoon@boatz.com",
		Name:  "Sailor Moon",
		Phone: "555-5555",
	}
	token, verified, err = GetContactVerification(contact)
	require.NoError(t, err, "nil contact extra returns error")
	require.False(t, verified)
	require.Empty(t, token)

	// Contact cannot be nil to set extra
	err = SetContactVerification(nil, "12345", false)
	require.Error(t, err)

	// Set extra on contact
	err = SetContactVerification(contact, "12345", false)
	require.NoError(t, err)

	// Fetch set extra
	token, verified, err = GetContactVerification(contact)
	require.NoError(t, err)
	require.False(t, verified)
	require.Equal(t, "12345", token)

	// Append to email log
	err = AppendEmailLog(contact, "verify_contact", "verification")
	require.NoError(t, err)
	require.False(t, verified)
	require.Equal(t, "12345", token)

	// Fetch email log
	emailLog, err := GetEmailLog(contact)
	require.NoError(t, err)
	require.Len(t, emailLog, 1)
	require.Equal(t, "verify_contact", emailLog[0].Reason)
	require.Equal(t, "verification", emailLog[0].Subject)

	// Should not overwrite contact verification
	token, verified, err = GetContactVerification(contact)
	require.NoError(t, err)
	require.False(t, verified)
	require.Equal(t, "12345", token)

	// Set extra on contact
	err = SetContactVerification(contact, "", true)
	require.NoError(t, err)

	// Fetch set extra
	token, verified, err = GetContactVerification(contact)
	require.NoError(t, err)
	require.True(t, verified)
	require.Equal(t, "", token)

	// Should not overwrite email log
	emailLog, err = GetEmailLog(contact)
	require.NoError(t, err)
	require.Len(t, emailLog, 1)
	require.Equal(t, "verify_contact", emailLog[0].Reason)
	require.Equal(t, "verification", emailLog[0].Subject)
}

func TestEmailLog(t *testing.T) {
	// Test that the email log functions are working as expected
	contact := &pb.Contact{}

	// Audit log should initially be empty
	emailLog, err := GetEmailLog(contact)
	require.NoError(t, err)
	require.Len(t, emailLog, 0)

	// Should not be able to append on a nil contact
	err = AppendEmailLog(nil, "verify_contact", "verification")
	require.Error(t, err)

	// Append an entry to an empty log
	err = AppendEmailLog(contact, "verify_contact", "verification")
	require.NoError(t, err)
	emailLog, err = GetEmailLog(contact)
	require.NoError(t, err)
	require.Len(t, emailLog, 1)
	require.Equal(t, "verify_contact", emailLog[0].Reason)
	require.Equal(t, "verification", emailLog[0].Subject)

	// Append another entry to the email log
	err = AppendEmailLog(contact, "review", "review resend")
	require.NoError(t, err)
	emailLog, err = GetEmailLog(contact)
	require.NoError(t, err)
	require.Len(t, emailLog, 2)
	require.Equal(t, "verify_contact", emailLog[0].Reason)
	require.Equal(t, "verification", emailLog[0].Subject)
	require.Equal(t, "review", emailLog[1].Reason)
	require.Equal(t, "review resend", emailLog[1].Subject)
}

func TestVeriedContacts(t *testing.T) {
	vasp := &pb.VASP{
		Contacts: &pb.Contacts{
			Administrative: &pb.Contact{
				Name:  "Admin Person",
				Email: "admin@example.com",
			},
			Technical: &pb.Contact{
				Name:  "Technical Person",
				Email: "tech@example.com",
			},
			Legal: &pb.Contact{
				Name:  "Legal Person",
				Email: "legal@example.com",
			},
		},
	}

	contacts := VerifiedContacts(vasp)
	require.Len(t, contacts, 0)

	err := SetContactVerification(vasp.Contacts.Administrative, "", true)
	require.NoError(t, err)

	err = SetContactVerification(vasp.Contacts.Technical, "12345", false)
	require.NoError(t, err)

	contacts = VerifiedContacts(vasp)
	require.Len(t, contacts, 1)

	err = SetContactVerification(vasp.Contacts.Technical, "", true)
	require.NoError(t, err)

	err = SetContactVerification(vasp.Contacts.Legal, "12345", false)
	require.NoError(t, err)

	contacts = VerifiedContacts(vasp)
	require.Len(t, contacts, 2)
}

func TestUpdateCertificateRequestStatus(t *testing.T) {
	// Attempt to set request status on a nil object
	err := UpdateCertificateRequestStatus(nil, CertificateRequestState_READY_TO_SUBMIT, "ready to submit", "automated")
	require.Error(t, err)

	// Attempt to set request status to an invalid state
	request := &CertificateRequest{}
	err = UpdateCertificateRequestStatus(request, -1, "invalid", "automated")
	require.Error(t, err)
	err = UpdateCertificateRequestStatus(request, CertificateRequestState_CR_ERRORED+1, "invalid", "automated")
	require.Error(t, err)

	// Update a brand new request with no previous state
	expectedTime := time.Now()
	err = UpdateCertificateRequestStatus(request, CertificateRequestState_READY_TO_SUBMIT, "ready to submit", "automated")
	require.NoError(t, err)

	// Check request status and audit log
	require.Equal(t, request.Status, CertificateRequestState_READY_TO_SUBMIT)
	require.Len(t, request.AuditLog, 1)

	// Timestamps should be close
	actualTime, err := time.Parse(time.RFC3339, request.AuditLog[0].Timestamp)
	require.NoError(t, err)
	require.LessOrEqual(t, expectedTime.Sub(actualTime), time.Duration(time.Minute))

	// Verify audit log entry
	require.Equal(t, CertificateRequestState_INITIALIZED, request.AuditLog[0].PreviousState)
	require.Equal(t, CertificateRequestState_READY_TO_SUBMIT, request.AuditLog[0].CurrentState)
	require.Equal(t, "ready to submit", request.AuditLog[0].Description)
	require.Equal(t, "automated", request.AuditLog[0].Source)

	// Change the status of the request again
	expectedTime = time.Now()
	err = UpdateCertificateRequestStatus(request, CertificateRequestState_PROCESSING, "processing", "automated")
	require.NoError(t, err)

	// Check request status and audit log
	require.Equal(t, request.Status, CertificateRequestState_PROCESSING)
	require.Len(t, request.AuditLog, 2)

	// Timestamps should be close
	actualTime, err = time.Parse(time.RFC3339, request.AuditLog[1].Timestamp)
	require.NoError(t, err)
	require.LessOrEqual(t, expectedTime.Sub(actualTime), time.Duration(time.Minute))

	// Verify audit log entries
	require.Equal(t, CertificateRequestState_INITIALIZED, request.AuditLog[0].PreviousState)
	require.Equal(t, CertificateRequestState_READY_TO_SUBMIT, request.AuditLog[0].CurrentState)
	require.Equal(t, "ready to submit", request.AuditLog[0].Description)
	require.Equal(t, "automated", request.AuditLog[0].Source)
	require.Equal(t, CertificateRequestState_READY_TO_SUBMIT, request.AuditLog[1].PreviousState)
	require.Equal(t, CertificateRequestState_PROCESSING, request.AuditLog[1].CurrentState)
	require.Equal(t, "processing", request.AuditLog[1].Description)
	require.Equal(t, "automated", request.AuditLog[1].Source)
}

func TestIsTraveler(t *testing.T) {
	vasp := &pb.VASP{CommonName: "trisa.example.com"}
	require.False(t, IsTraveler(vasp))

	vasp.CommonName = "trisa-1234.traveler.ciphertrace.com"
	require.True(t, IsTraveler(vasp))
}

func TestUpdateVerificationStatus(t *testing.T) {
	vasp := &pb.VASP{}

	// Update a brand new VASP with no previous state
	expectedTime := time.Now()
	err := UpdateVerificationStatus(vasp, pb.VerificationState_SUBMITTED, "submitted", "automated")
	require.NoError(t, err)

	// Check verification status and audit log
	require.Equal(t, vasp.VerificationStatus, pb.VerificationState_SUBMITTED)
	auditLog, err := GetAuditLog(vasp)
	require.NoError(t, err)
	require.Len(t, auditLog, 1)

	// Timestamps should be close
	actualTime, err := time.Parse(time.RFC3339, auditLog[0].Timestamp)
	require.NoError(t, err)
	require.LessOrEqual(t, expectedTime.Sub(actualTime), time.Duration(time.Minute))

	// Verify audit log entry
	require.Equal(t, pb.VerificationState_NO_VERIFICATION, auditLog[0].PreviousState)
	require.Equal(t, pb.VerificationState_SUBMITTED, auditLog[0].CurrentState)
	require.Equal(t, "submitted", auditLog[0].Description)
	require.Equal(t, "automated", auditLog[0].Source)

	// Change the state of the VASP again
	expectedTime = time.Now()
	err = UpdateVerificationStatus(vasp, pb.VerificationState_REVIEWED, "review completed", "pontoon@boatz.com")
	require.NoError(t, err)

	// Check verification status and audit log
	require.Equal(t, vasp.VerificationStatus, pb.VerificationState_REVIEWED)
	auditLog, err = GetAuditLog(vasp)
	require.NoError(t, err)
	require.Len(t, auditLog, 2)

	// Timestamps should be close
	actualTime, err = time.Parse(time.RFC3339, auditLog[1].Timestamp)
	require.NoError(t, err)
	require.LessOrEqual(t, expectedTime.Sub(actualTime), time.Duration(time.Minute))

	// Verify audit log entry
	require.Equal(t, pb.VerificationState_SUBMITTED, auditLog[1].PreviousState)
	require.Equal(t, pb.VerificationState_REVIEWED, auditLog[1].CurrentState)
	require.Equal(t, "review completed", auditLog[1].Description)
	require.Equal(t, "pontoon@boatz.com", auditLog[1].Source)
}
