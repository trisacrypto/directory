package contacts_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/contacts"
	"github.com/trisacrypto/directory/pkg/models/v1"
)

func TestContactHasEmail(t *testing.T) {
	require.False(t, contacts.ContactHasEmail(nil))
	require.False(t, contacts.ContactHasEmail(&models.Contact{Email: ""}))
	require.True(t, contacts.ContactHasEmail(&models.Contact{Email: "test@test.com"}))
}

func TestContactIsVerified(t *testing.T) {
	require.False(t, contacts.ContactIsVerified(nil))
	require.False(t, contacts.ContactIsVerified(&models.Contact{Verified: false}))
	require.True(t, contacts.ContactIsVerified(&models.Contact{Verified: true}))
}

func TestGetContactVerification(t *testing.T) {
	token, verified := contacts.GetContactVerification(nil)
	require.Empty(t, token)
	require.Empty(t, verified)

	token, verified = contacts.GetContactVerification(&models.Contact{
		Token:    "token",
		Verified: true,
	})
	require.Equal(t, token, "token")
	require.True(t, verified)
}

func TestSetContactVerification(t *testing.T) {
	err := contacts.SetContactVerification(nil, "", false)
	require.EqualError(t, err, "cannot set verification on nil contact")

	contact := &models.Contact{
		Token:    "",
		Verified: false,
	}
	err = contacts.SetContactVerification(contact, "new_token", true)
	require.NoError(t, err)
	require.Equal(t, contact.Token, "new_token")
	require.True(t, contact.Verified)
}

func TestGetEmailLog(t *testing.T) {
	require.Empty(t, contacts.GetEmailLog(nil))

	log := []*models.EmailLogEntry{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
			Recipient: "test@test.com",
		},
	}
	contact := models.Contact{EmailLog: log}
	retrievedLog := contacts.GetEmailLog(&contact)
	require.Equal(t, log, retrievedLog)
}

func TestAppendEmailLog(t *testing.T) {
	err := contacts.AppendEmailLog(nil, "", "")
	require.EqualError(t, err, "cannot append entry to nil contact")

	contact := &models.Contact{Email: "test@test.com"}
	require.NoError(t, contacts.AppendEmailLog(contact, "verify_contact", "verification"))
	require.Equal(t, contact.EmailLog[0].Reason, "verify_contact")
	require.Equal(t, contact.EmailLog[0].Subject, "verification")
	require.Equal(t, contact.EmailLog[0].Recipient, "test@test.com")

	contact.Email = "new_test@test.com"
	require.NoError(t, contacts.AppendEmailLog(contact, "new_verify_contact", "new_verification"))
	require.Equal(t, contact.EmailLog[1].Reason, "new_verify_contact")
	require.Equal(t, contact.EmailLog[1].Subject, "new_verification")
	require.Equal(t, contact.EmailLog[1].Recipient, "new_test@test.com")
}
