package models_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"golang.org/x/exp/slices"
)

func TestEmailValidation(t *testing.T) {
	testCases := []struct {
		email         *models.Email
		err           error
		expectedName  string
		expectedEmail string
	}{
		{&models.Email{}, models.ErrNoEmailAddress, "", ""},
		{&models.Email{Email: "\t\n\t\n\t\t\t"}, models.ErrNoEmailAddress, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: true, VerifiedOn: "", Token: ""}, models.ErrVerifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: true, VerifiedOn: "", Token: "foo"}, models.ErrVerifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: "foo"}, models.ErrVerifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, nil, "", "ted@example.com"},
		{&models.Email{Email: "ted@example.com", Verified: false, VerifiedOn: "", Token: ""}, models.ErrUnverifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: false, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, models.ErrUnverifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: false, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: "foo"}, models.ErrUnverifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: false, VerifiedOn: "", Token: "foo"}, nil, "", "ted@example.com"},
		{&models.Email{Email: "TED@example.com", Verified: false, VerifiedOn: "", Token: "foo"}, nil, "", "ted@example.com"},
		{&models.Email{Email: "Ted Tonks <TED@example.com>", Verified: false, VerifiedOn: "", Token: "foo"}, nil, "Ted Tonks", "ted@example.com"},
		{&models.Email{Name: "James Surry", Email: "Ted Tonks <TED@example.com>", Verified: false, VerifiedOn: "", Token: "foo"}, nil, "James Surry", "ted@example.com"},
		{&models.Email{Email: "TED@example.com", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, nil, "", "ted@example.com"},
		{&models.Email{Email: "Ted Tonks <TED@example.com>", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, nil, "Ted Tonks", "ted@example.com"},
		{&models.Email{Name: "James Surry", Email: "Ted Tonks <TED@example.com>", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, nil, "James Surry", "ted@example.com"},
	}

	for i, tc := range testCases {
		err := tc.email.Validate()
		if tc.err == nil {
			require.NoError(t, err, "test case %d failed with error", i)
			require.Equal(t, tc.expectedName, tc.email.Name, "test case %d failed with name mismatch", i)
			require.Equal(t, tc.expectedEmail, tc.email.Email, "test case %d failed with email mismatch", i)
		} else {
			require.ErrorIs(t, err, tc.err, "test case %d failed with incorrect error", i)
		}
	}

}

func TestEmailLog(t *testing.T) {
	email := &models.Email{Email: "James Surry", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00"}
	require.NoError(t, email.Validate())

	for i := 1; i < 11; i++ {
		email.Log("testing", fmt.Sprintf("test email %d", i))
	}

	require.Len(t, email.SendLog, 10, "expected 10 log entries in the database")
	for _, entry := range email.SendLog {
		require.NotEmpty(t, entry.Timestamp, "entry did not have a timestamp")
		require.Equal(t, entry.Reason, "testing")
		require.True(t, strings.HasPrefix(entry.Subject, "test email"))
		require.Equal(t, entry.Recipient, email.Email)
	}
}

func TestEmailVASPs(t *testing.T) {

	vaspIDs := []string{
		"56f0251a-6670-4c8b-ac5a-328ec1f89a4b",
		"93d21f48-f71a-4d13-a4cb-512f09a0755a",
		"56f0251a-6670-4c8b-ac5a-328ec1f89a4b",
		"a054e7db-b6d5-49f0-9ad7-c406fdfc01df",
		"4358418b-d1d8-46b6-be44-4462b3f2afa8",
		"5f09a45d-c9d2-4f04-99ca-9c9caa22ce9d",
		"5641f91a-6102-4202-8a49-7c9295d439d2",
		"a054e7db-b6d5-49f0-9ad7-c406fdfc01df",
		"589baa68-f870-4d7f-a556-0c3150aa2589",
		"306f540c-baf1-46be-84a4-94c0cb9ec293",
		"306f540c-baf1-46be-84a4-94c0cb9ec293",
	}

	email := &models.Email{Name: "Ralph Stein", Email: "ralph@example.com", Verified: true, VerifiedOn: "2023-09-11T07:41:16-05:00"}

	// Add vaspIDs to the email
	for _, vaspID := range vaspIDs {
		email.AddVASP(vaspID)
	}

	require.Len(t, email.Vasps, 8, "without duplicates expected 8 vaspIDs")
	require.True(t, slices.IsSorted(email.Vasps), "expected the vasps slice to be sorted")

	// Remove vaspIDs from the email
	for _, vaspID := range vaspIDs {
		email.RmVASP(vaspID)
	}

	require.Empty(t, email.Vasps, "expected vasp list to be empty after removal")
}

func TestCountSentEmails(t *testing.T) {
	data, err := os.ReadFile("testdata/log.pb.json")
	require.NoError(t, err, "could not read testdata/log.pb.json")

	var emailLog []*models.EmailLogEntry
	err = json.Unmarshal(data, &emailLog)
	require.NoError(t, err, "could not unmarshal email log")
	require.Len(t, emailLog, 10, "has the email log fixture changed?")

	// Timestamps in the email log must be relative to now
	now := time.Now()
	emailLog[0].Timestamp = now.Add(-30 * 24 * time.Hour).Format(time.RFC3339)
	emailLog[1].Timestamp = now.Add(-20 * 24 * time.Hour).Format(time.RFC3339)
	emailLog[2].Timestamp = now.Add(-10 * 24 * time.Hour).Format(time.RFC3339)
	emailLog[3].Timestamp = now.Add(-7 * 24 * time.Hour).Format(time.RFC3339)
	emailLog[4].Timestamp = now.Add(-4*24*time.Hour - 2*time.Hour).Format(time.RFC3339)
	emailLog[5].Timestamp = now.Add(-4*25*time.Hour - 1*time.Hour).Format(time.RFC3339)
	emailLog[6].Timestamp = now.Add(-4 * 24 * time.Hour).Format(time.RFC3339)
	emailLog[7].Timestamp = now.Add(-8 * time.Hour).Format(time.RFC3339)
	emailLog[8].Timestamp = now.Add(-7 * time.Hour).Format(time.RFC3339)
	emailLog[9].Timestamp = now.Add(-6 * time.Hour).Format(time.RFC3339)

	t.Run("EmptyLog", func(t *testing.T) {
		sent, err := models.CountSentEmails(nil, "test", 30)
		require.NoError(t, err, "expected no error on nil email log")
		require.Zero(t, sent, "expected sent to be zero")

		sent, err = models.CountSentEmails([]*models.EmailLogEntry{}, "test", 30)
		require.NoError(t, err, "expected no error on nil email log")
		require.Zero(t, sent, "expected sent to be zero")
	})

	t.Run("Invariants", func(t *testing.T) {
		// Error should be returned if the reason is empty
		_, err := models.CountSentEmails(emailLog, "", 30)
		require.ErrorIs(t, err, models.ErrNoLogReason)

		// Error should be returned if the time window is invalid
		_, err = models.CountSentEmails(emailLog, "test", -1)
		require.ErrorIs(t, err, models.ErrInvalidWindow)
	})

	t.Run("Happy", func(t *testing.T) {
		testCases := []struct {
			reason   string
			window   int
			expected int
		}{
			{"verify_contact", 5, 0},
			{"verify_contact", 8, 1},
			{"verify_contact", 15, 2},
			{"verify_contact", 31, 4},
			{"reissuance_started", 1, 1},
			{"reissuance_started", 3, 1},
			{"reissuance_started", 7, 2},
			{"deliver_certs", 1, 1},
			{"deliver_certs", 3, 1},
			{"deliver_certs", 7, 2},
			{"reissuance_reminder", 1, 1},
			{"reissuance_reminder", 3, 1},
			{"reissuance_reminder", 7, 2},
		}

		for i, tc := range testCases {
			actual, err := models.CountSentEmails(emailLog, tc.reason, tc.window)
			require.NoError(t, err, "test case %d failed with error", i)
			require.Equal(t, tc.expected, actual, "test case %d failed with sent mismatch", i)
		}
	})

	t.Run("TimestampParseError", func(t *testing.T) {
		badLog := []*models.EmailLogEntry{
			{
				Timestamp: "foo",
				Reason:    "deadline",
				Subject:   "this is a bad eamil",
				Recipient: "jill@example.com",
			},
		}

		// Error returned when timestamp is unparsable
		parseError := &time.ParseError{}
		_, err := models.CountSentEmails(badLog, "deadline", 30)
		require.ErrorAs(t, err, &parseError)

		// No error returned if reason doesn't match
		_, err = models.CountSentEmails(badLog, "verify_email", 30)
		require.NoError(t, err)
	})

}

func TestNormalizeEmail(t *testing.T) {
	testCases := []struct {
		email    string
		expected string
	}{
		{"support@trisa.io", "support@trisa.io"},
		{"Gary.Verdun@example.com", "gary.verdun@example.com"},
		{"   jessica@blankspace.net       ", "jessica@blankspace.net"},
		{"\t\t\nweird@foo.co.uk\t\n", "weird@foo.co.uk"},
		{"ALLCAPSCREAM@WILD.FR", "allcapscream@wild.fr"},
		{"Gary Verdun <gary@example.com>", "gary@example.com"},
	}

	for i, tc := range testCases {
		require.Equal(t, tc.expected, models.NormalizeEmail(tc.email), "test case %d failed", i)
	}
}
