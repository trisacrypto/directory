package models_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/models/v1"
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
