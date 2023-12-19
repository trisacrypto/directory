package models_test

import (
	"bytes"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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

	// Deleting certificate request IDs from an empty slice should not error
	require.NoError(t, DeleteCertReqID(vasp, "b5841869-105f-411c-8722-4045aad72717"))

	// Attempt to fetch the latest certificate request ID, should return empty string
	id, err := GetLatestCertReqID(vasp)
	require.NoError(t, err)
	require.Empty(t, id)

	// Attempt to append certificate request IDs
	certReqs := []string{
		"b5841869-105f-411c-8722-4045aad72717",
		"230d5e77-9983-4f1f-80ea-d379d56519af",
	}
	for _, cfid := range certReqs {
		err = AppendCertReqID(vasp, cfid)
		require.NoError(t, err)
		latest, err := GetLatestCertReqID(vasp)
		require.NoError(t, err)
		require.Equal(t, cfid, latest)
	}

	// Should be able to fetch the certificate request IDs
	ids, err := GetCertReqIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 2)
	require.Equal(t, certReqs[0], ids[0])
	require.Equal(t, certReqs[1], ids[1])

	// Should be able to delete a certificate request ID
	require.NoError(t, DeleteCertReqID(vasp, certReqs[0]))
	ids, err = GetCertReqIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 1)
	require.Equal(t, certReqs[1], ids[0])

	// Latest certificate request ID should be the remaining one
	latest, err := GetLatestCertReqID(vasp)
	require.NoError(t, err)
	require.Equal(t, certReqs[1], latest)

	// Deleting the certificate request ID again should not error
	require.NoError(t, DeleteCertReqID(vasp, certReqs[0]))

	// Deleting a non-existent certificate request ID should not error
	require.NoError(t, DeleteCertReqID(vasp, "does-not-exist"))

	// Certificate request IDs should be unchanged
	ids, err = GetCertReqIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 1)
	require.Equal(t, certReqs[1], ids[0])

	// Attempt to append certificate IDs
	certs := []string{
		"c9838f8f-f8f8-4f8f-8f8f-f8f8f8f8f8f8",
		"f8f2349d-f8f8-4f8f-8f8f-f8f8f8f8f8f8",
	}
	for _, certID := range certs {
		err = AppendCertID(vasp, certID)
		require.NoError(t, err)
	}

	// Should be able to fetch the certificate IDs
	ids, err = GetCertIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 2)
	require.Equal(t, certs[0], ids[0])
	require.Equal(t, certs[1], ids[1])

	// Verification token should be unchanged
	token, err = GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "jetskis", token)

	// Initial email log should be nil
	emailLog, err := GetAdminEmailLog(vasp)
	require.NoError(t, err)
	require.Len(t, emailLog, 0)

	// Should not be able to append on a nil VASP
	err = AppendAdminEmailLog(nil, "certificate_reissued", "certificate has been reissued to Alice VASP")
	require.Error(t, err)

	// Append an entry to an empty log
	err = AppendAdminEmailLog(vasp, "certificate_reissued", "certificate has been reissued to Alice VASP")
	require.NoError(t, err)
	emailLog, err = GetAdminEmailLog(vasp)
	require.NoError(t, err)
	require.Len(t, emailLog, 1)
	require.Equal(t, "certificate_reissued", emailLog[0].Reason)
	require.Equal(t, "certificate has been reissued to Alice VASP", emailLog[0].Subject)
	require.NotEmpty(t, emailLog[0].Timestamp)

	// Append another entry to the email log
	err = AppendAdminEmailLog(vasp, "certificate_reissued_again", "certificate has been reissued to Alice VASP (again)")
	require.NoError(t, err)
	emailLog, err = GetAdminEmailLog(vasp)
	require.NoError(t, err)
	require.Len(t, emailLog, 2)
	require.Equal(t, "certificate_reissued", emailLog[0].Reason)
	require.Equal(t, "certificate has been reissued to Alice VASP", emailLog[0].Subject)
	require.NotEmpty(t, emailLog[0].Timestamp)
	require.Equal(t, "certificate_reissued_again", emailLog[1].Reason)
	require.Equal(t, "certificate has been reissued to Alice VASP (again)", emailLog[1].Subject)
	require.NotEmpty(t, emailLog[1].Timestamp)
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

func TestCertIDs(t *testing.T) {
	vasp := &pb.VASP{}

	// No extra, Get should return nil
	ids, err := GetCertIDs(vasp)
	require.NoError(t, err)
	require.Empty(t, ids)

	// Cannot append an empty ID
	err = AppendCertID(vasp, "")
	require.EqualError(t, err, "cannot append empty certificate ID to extra")

	// Append an ID from empty
	err = AppendCertID(vasp, "1df61840-7033-40fb-8ce9-538c87e242f5")
	require.NoError(t, err)
	ids, err = GetCertIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 1)

	// Append an ID with data already inside
	err = AppendCertID(vasp, "9676bf6a-ffdb-4185-8fa5-87cdae6f6eef")
	require.NoError(t, err)
	ids, err = GetCertIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 2)

	// Do not allow duplicate IDs
	err = AppendCertID(vasp, "9676bf6a-ffdb-4185-8fa5-87cdae6f6eef")
	require.NoError(t, err)
	ids, err = GetCertIDs(vasp)
	require.NoError(t, err)
	require.Len(t, ids, 2)
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

func TestNewCertificate(t *testing.T) {
	vasp := &pb.VASP{
		Id: "b5841869-105f-411c-8722-4045aad72717",
	}

	certReq := &CertificateRequest{
		Id: "c8f8f8f8-f8f8-f8f8-f8f8-f8f8f8f8f8f8",
	}

	// The serial number must be a capital hex-encoded value to mirror the sectigo format
	serial := "ABC83132333435363738"
	serialBytes, err := hex.DecodeString(serial)
	require.NoError(t, err)
	pub := &pb.Certificate{
		SerialNumber: serialBytes,
	}

	// Should not be able to create a certificate with a nil vasp
	_, err = NewCertificate(nil, certReq, pub)
	require.Error(t, err)

	// Should not be able to create a certificate with a nil request
	_, err = NewCertificate(vasp, nil, pub)
	require.Error(t, err)

	// Should not be able to create a certificate with nil certificate data
	_, err = NewCertificate(vasp, certReq, nil)
	require.Error(t, err)

	// Certificate correctly created
	cert, err := NewCertificate(vasp, certReq, pub)
	require.NoError(t, err)
	require.Equal(t, serial, cert.Id)
	require.Equal(t, certReq.Id, cert.Request)
	require.Equal(t, vasp.Id, cert.Vasp)
	require.Equal(t, CertificateState_ISSUED, cert.Status)
	require.True(t, proto.Equal(pub, cert.Details))
}

func TestNewCertificateRequest(t *testing.T) {
	// Should not be able to create a certificate request with a nil vasp
	_, err := NewCertificateRequest(nil)
	require.Error(t, err)

	// If the name does not exist, it should not be populated
	vasp := &pb.VASP{
		Id:         "b5841869-105f-411c-8722-4045aad72717",
		CommonName: "charlieVASP",
		Entity:     &ivms101.LegalPerson{},
	}
	cr, err := NewCertificateRequest(vasp)
	require.NoError(t, err)
	require.Equal(t, vasp.Id, cr.Vasp)
	require.Equal(t, vasp.CommonName, cr.CommonName)
	require.NotContains(t, cr.Params, sectigo.ParamOrganizationName)

	vasp.Entity.Name = &ivms101.LegalPersonName{
		NameIdentifiers: []*ivms101.LegalPersonNameId{
			{
				LegalPersonName:               "Charlie Inc.",
				LegalPersonNameIdentifierType: ivms101.LegalPersonShort,
			},
		},
	}
	cr, err = NewCertificateRequest(vasp)
	require.NoError(t, err)
	require.Contains(t, cr.Params, sectigo.ParamOrganizationName)
	require.Equal(t, "Charlie Inc.", cr.Params[sectigo.ParamOrganizationName])

	// Valid organization name, unspecified location info
	vasp.Entity.Name.NameIdentifiers[0].LegalPersonNameIdentifierType = ivms101.LegalPersonLegal
	cr, err = NewCertificateRequest(vasp)
	require.NoError(t, err)
	require.Contains(t, cr.Params, sectigo.ParamOrganizationName)
	require.Equal(t, "Charlie Inc.", cr.Params[sectigo.ParamOrganizationName])
	require.NotContains(t, cr.Params, sectigo.ParamLocalityName)
	require.NotContains(t, cr.Params, sectigo.ParamStateOrProvinceName)
	require.NotContains(t, cr.Params, sectigo.ParamCountryName)

	// Valid organization name, partial location info is not accepted
	vasp.Entity.GeographicAddresses = []*ivms101.Address{
		{
			Country: "CA",
		},
	}
	cr, err = NewCertificateRequest(vasp)
	require.NoError(t, err)
	require.Contains(t, cr.Params, sectigo.ParamOrganizationName)
	require.Equal(t, "Charlie Inc.", cr.Params[sectigo.ParamOrganizationName])
	require.NotContains(t, cr.Params, sectigo.ParamLocalityName)
	require.NotContains(t, cr.Params, sectigo.ParamStateOrProvinceName)
	require.NotContains(t, cr.Params, sectigo.ParamCountryName)

	// Complete organization name and location info
	vasp.Entity.GeographicAddresses[0].TownLocationName = "Toronto"
	vasp.Entity.GeographicAddresses[0].CountrySubDivision = "Ontario"
	vasp.Entity.GeographicAddresses[0].Country = "CA"
	cr, err = NewCertificateRequest(vasp)
	require.NoError(t, err)
	require.Contains(t, cr.Params, sectigo.ParamOrganizationName)
	require.Equal(t, "Charlie Inc.", cr.Params[sectigo.ParamOrganizationName])
	require.Contains(t, cr.Params, sectigo.ParamLocalityName)
	require.Equal(t, "Toronto", cr.Params[sectigo.ParamLocalityName])
	require.Contains(t, cr.Params, sectigo.ParamStateOrProvinceName)
	require.Equal(t, "Ontario", cr.Params[sectigo.ParamStateOrProvinceName])
	require.Contains(t, cr.Params, sectigo.ParamCountryName)
	require.Equal(t, "CA", cr.Params[sectigo.ParamCountryName])
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

func TestValidateVASP(t *testing.T) {
	vasp := &pb.VASP{
		CommonName:         "trisa.charliebank.com",
		TrisaEndpoint:      "trisa.charliebank.io:123",
		VerificationStatus: pb.VerificationState_SUBMITTED,
		Contacts: &pb.Contacts{
			Administrative: &pb.Contact{
				Email: "glenn@charliebank.com",
				Name:  "Glenn Davis",
			},
		},
		Entity: &ivms101.LegalPerson{
			CountryOfRegistration: "CA",
			GeographicAddresses: []*ivms101.Address{
				{
					AddressLine: []string{"123 Main St"},
					AddressType: ivms101.AddressTypeCode_ADDRESS_TYPE_CODE_BIZZ,
					Country:     "CA",
				},
			},
			Name: &ivms101.LegalPersonName{
				NameIdentifiers: []*ivms101.LegalPersonNameId{
					{
						LegalPersonName:               "Charlie Inc.",
						LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_LEGL,
					},
				},
			},
			NationalIdentification: &ivms101.NationalIdentification{
				NationalIdentifier:     "123456789",
				NationalIdentifierType: ivms101.NationalIdentifierLEIX,
			},
		},
	}
	require.Error(t, ValidateVASP(vasp, false), "expected failed validation if fields are missing and partial is false")
	require.NoError(t, ValidateVASP(vasp, true), "expected successful validation if fields are missing but partial is true")

	vasp.Id = "b5841869-105f-411c-8722-4045aad72717"
	vasp.RegisteredDirectory = "trisatest.net"
	vasp.FirstListed = time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	vasp.LastUpdated = time.Now().Format(time.RFC3339)
	vasp.Signature = []byte("abc123")
	require.NoError(t, ValidateVASP(vasp, false), "expected successful validation when partial is false")

	// Verify that the validation helper ignores C9 constraint errors
	vasp.Entity.NationalIdentification.CountryOfIssue = "CA"
	require.NoError(t, ValidateVASP(vasp, false), "expected successful validation even when C9 constraint is violated")
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

func TestVASPSignature(t *testing.T) {
	// Signature comparison assertion (can require equal or require not equal using cmp)
	compare := func(a, b *pb.VASP, cmp require.BoolAssertionFunc, msg ...interface{}) {
		siga, err := VASPSignature(a)
		require.NoError(t, err, "could not get signature from vaspa")

		sigb, err := VASPSignature(b)
		require.NoError(t, err, "could not get signature from vaspb")

		cmp(t, bytes.Equal(siga, sigb), msg...)
	}

	// empty VASPs should have the same signature
	vaspa := &pb.VASP{}
	vaspb := &pb.VASP{}
	compare(vaspa, vaspb, require.True)

	// load a VASP fixture from disk
	vaspa, err := loadFixture("testdata/vasp.json")
	require.NoError(t, err, "could not load vasp fixture from disk")

	// full vasp should not be the same as empty vasp
	compare(vaspa, vaspb, require.False)

	// load a second VASP fixture from disk (different pointer)
	vaspb, err = loadFixture("testdata/vasp.json")
	require.NoError(t, err, "could not load vasp fixture from disk again")

	compare(vaspa, vaspb, require.True)

	// make changes to vaspa to ensure its signature changes
	siga, err := VASPSignature(vaspa)
	require.NoError(t, err, "could not compute vaspa signature")

	vaspa.Extra, _ = anypb.New(&GDSExtraData{})

	_, err = CreateReviewNote(vaspa, "123", "admin@example.com", "this is a test note")
	require.NoError(t, err, "could not add review note")

	siga2, err := VASPSignature(vaspa)
	require.NoError(t, err, "could not compute vaspa signature")
	require.NotEqual(t, siga, siga2, "changing vasp a did not change singature")

	err = AppendAdminEmailLog(vaspa, "manual entry", "this is a test email log")
	require.NoError(t, err, "could not append admin email log")

	siga3, err := VASPSignature(vaspa)
	require.NoError(t, err, "could not compute vaspa signature")
	require.NotEqual(t, siga2, siga3, "changing vasp a did not change singature")
}

func TestGetVASPEmailLog(t *testing.T) {
	// Create a function that will create a contacts fixture for testing GetVASPEmailLog
	makeContacts := func(contacts *pb.Contacts) *Contacts {
		if contacts == nil {
			contacts = &pb.Contacts{
				Administrative: &pb.Contact{
					Name:  "Ashley Quickstar",
					Email: "admin@example.com",
				},
				Technical: &pb.Contact{
					Name:  "Billy Rester",
					Email: "tech@example.com",
				},
				Legal: &pb.Contact{
					Name:  "Cathleen Studeville",
					Email: "legal@example.com",
				},
				Billing: &pb.Contact{
					Name:  "David Teeter",
					Email: "billing@example.com",
				},
			}
		}

		fixture := &Contacts{
			VASP:     "b2fc0f56-3121-492f-8cd5-540f15456f6f",
			Contacts: contacts,
			Emails:   make([]*Email, 0),
		}

		seen := make(map[string]struct{})
		emails := []*pb.Contact{contacts.Administrative, contacts.Technical, contacts.Billing, contacts.Legal}

		for _, email := range emails {
			if email != nil && email.Email != "" {
				if _, ok := seen[email.Email]; !ok {
					record := &Email{
						Name:    email.Name,
						Email:   email.Email,
						Vasps:   []string{fixture.VASP},
						SendLog: make([]*EmailLogEntry, 0),
					}

					fixture.Emails = append(fixture.Emails, record)
					seen[email.Email] = struct{}{}
				}
			}
		}

		return fixture
	}

	// Email fixtures for the tests below
	now := time.Now()
	messages := []*EmailLogEntry{
		{
			Reason:    "verify_contact",
			Subject:   "verify_admin",
			Timestamp: now.Format(time.RFC3339),
			Recipient: "admin@example.com",
		},
		{
			Reason:    "reissuance",
			Subject:   "reissuance_admin",
			Timestamp: now.AddDate(0, 0, 1).Format(time.RFC3339),
			Recipient: "admin@example.com",
		},
		{
			Reason:    "verify_contact",
			Subject:   "verify_tech",
			Timestamp: now.Add(time.Hour).Format(time.RFC3339),
			Recipient: "tech@example.com",
		},
		{
			Reason:    "reissuance",
			Subject:   "reissuance_tech",
			Timestamp: now.Add(time.Hour * 2).Format(time.RFC3339),
			Recipient: "tech@example.com",
		},
		{
			Reason:    "verify_contact",
			Subject:   "verify_billing",
			Timestamp: now.Add(-time.Hour).Format(time.RFC3339),
			Recipient: "billing@example.com",
		},
		{
			Reason:    "resend",
			Subject:   "resend_billing",
			Timestamp: now.Add(time.Hour * 3).Format(time.RFC3339),
			Recipient: "billing@example.com",
		},
		{
			Reason:    "reissuance",
			Subject:   "reissuance_billing",
			Timestamp: now.AddDate(0, 0, 2).Format(time.RFC3339),
			Recipient: "billing@example.com",
		},
	}

	t.Run("Nil", func(t *testing.T) {
		emails, err := GetVASPEmailLog(nil)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emails, 0)
	})

	t.Run("Empty", func(t *testing.T) {
		// Create a contacts data structure that is completely empty
		contacts := makeContacts(&pb.Contacts{})

		// Should return an empty slice if there are no contacts
		emails, err := GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emails, 0)
	})

	t.Run("Single", func(t *testing.T) {
		// Create a single contact with some email log entries
		contacts := makeContacts(&pb.Contacts{Administrative: &pb.Contact{Name: "Ashley Quickstar", Email: "admin@example.com"}})
		admin := contacts.Get(AdministrativeContact)

		admin.Email.SendLog = append(admin.Email.SendLog, messages[0], messages[1])

		// Should preserve ordering of the email log entries
		emails, err := GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emails, 2)
		require.Equal(t, messages[0].Subject, emails[0].Subject)
		require.Equal(t, messages[1].Subject, emails[1].Subject)
	})

	t.Run("SingleLog", func(t *testing.T) {
		// Only a single contact has log entries, the others should be empty logs
		contacts := makeContacts(nil)
		admin := contacts.Get(AdministrativeContact)

		// Should ignore contacts with no log entries
		admin.Email.SendLog = append(admin.Email.SendLog, messages[0], messages[1])
		emails, err := GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emails, 2)
		require.Equal(t, messages[0].Subject, emails[0].Subject)
		require.Equal(t, messages[1].Subject, emails[1].Subject)
	})

	t.Run("Duplicates", func(t *testing.T) {
		// Create four contacts all with the same email address and therefore the same
		// email log, but the logs should be deduplicated.
		contacts := makeContacts(&pb.Contacts{
			Administrative: &pb.Contact{Name: "Ashley Quickstar", Email: "admin@example.com"},
			Technical:      &pb.Contact{Name: "Ashley Quickstar", Email: "admin@example.com"},
			Billing:        &pb.Contact{Name: "Ashley Quickstar", Email: "admin@example.com"},
			Legal:          &pb.Contact{Name: "Ashley Quickstar", Email: "admin@example.com"},
		})
		admin := contacts.Get(AdministrativeContact)

		// Should ignore contacts with no log entries
		admin.Email.SendLog = append(admin.Email.SendLog, messages[0], messages[1])
		emails, err := GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emails, 2)
		require.Equal(t, messages[0].Subject, emails[0].Subject)
		require.Equal(t, messages[1].Subject, emails[1].Subject)
	})

	t.Run("Double", func(t *testing.T) {
		// Create contacts with two records, both with log entries
		contacts := makeContacts(&pb.Contacts{
			Administrative: &pb.Contact{Name: "Ashley Quickstar", Email: "admin@example.com"},
			Technical:      &pb.Contact{Name: "Ashley Quickstar", Email: "tech@example.com"},
		})

		admin := contacts.Get(AdministrativeContact)
		tech := contacts.Get(TechnicalContact)

		admin.Email.SendLog = append(admin.Email.SendLog, messages[0], messages[1])
		tech.Email.SendLog = append(tech.Email.SendLog, messages[2], messages[3])

		// Should properly merge the two email logs
		emails, err := GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emails, 4)
		require.Equal(t, messages[0].Subject, emails[0].Subject)
		require.Equal(t, messages[2].Subject, emails[1].Subject)
		require.Equal(t, messages[3].Subject, emails[2].Subject)
		require.Equal(t, messages[1].Subject, emails[3].Subject)
	})

	t.Run("Triple", func(t *testing.T) {
		// Create contacts where three contacts have log entries
		contacts := makeContacts(nil)

		admin := contacts.Get(AdministrativeContact)
		tech := contacts.Get(TechnicalContact)
		billing := contacts.Get(BillingContact)

		admin.Email.SendLog = append(admin.Email.SendLog, messages[0], messages[1])
		tech.Email.SendLog = append(tech.Email.SendLog, messages[2], messages[3])
		billing.Email.SendLog = append(billing.Email.SendLog, messages[4], messages[5], messages[6])

		// Should properly merge the three email logs
		emails, err := GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emails, 7)
		require.Equal(t, messages[4].Subject, emails[0].Subject)
		require.Equal(t, messages[0].Subject, emails[1].Subject)
		require.Equal(t, messages[2].Subject, emails[2].Subject)
		require.Equal(t, messages[3].Subject, emails[3].Subject)
		require.Equal(t, messages[5].Subject, emails[4].Subject)
		require.Equal(t, messages[1].Subject, emails[5].Subject)
		require.Equal(t, messages[6].Subject, emails[6].Subject)
	})

	t.Run("NoEmail", func(t *testing.T) {
		// Create contacts that do not have an email address.
		contacts := makeContacts(&pb.Contacts{
			Administrative: &pb.Contact{Name: "Ashley Quickstar", Email: "admin@example.com"},
			Technical:      &pb.Contact{Name: "Sydney Sneaks", Email: ""},
		})
		admin := contacts.Get(AdministrativeContact)

		// Should ignore contacts with no log entries
		admin.Email.SendLog = append(admin.Email.SendLog, messages[0], messages[1])
		emails, err := GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emails, 2)
		require.Equal(t, messages[0].Subject, emails[0].Subject)
		require.Equal(t, messages[1].Subject, emails[1].Subject)
	})
}

func loadFixture(path string) (vasp *pb.VASP, err error) {
	pbjson := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: false,
	}

	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return nil, err
	}

	vasp = &pb.VASP{}
	if err = pbjson.Unmarshal(data, vasp); err != nil {
		return nil, err
	}

	return vasp, nil
}
