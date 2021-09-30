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

	// Set extra on contact
	err = SetContactVerification(contact, "", true)
	require.NoError(t, err)

	// Fetch set extra
	token, verified, err = GetContactVerification(contact)
	require.NoError(t, err)
	require.True(t, verified)
	require.Equal(t, "", token)
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
