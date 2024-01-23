package emails

import (
	"encoding/json"
	"net/mail"
	"testing"
	"time"

	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/models/v1"
	emailutils "github.com/trisacrypto/directory/pkg/utils/emails"
	mock "github.com/trisacrypto/directory/pkg/utils/emails/mock"
)

type EmailMeta struct {
	Contact   *models.ContactRecord
	To        string
	From      string
	Subject   string
	Reason    string
	Timestamp time.Time
}

// Quick helper function to create expected emails in tests
func Expected(contact *models.ContactRecord, subject string, timestamp time.Time) *EmailMeta {
	return &EmailMeta{
		Contact: contact,
		To:      contact.Email.Email,
		Subject: subject,
		Reason:  Subject2Reason[subject],
	}
}

// Factory function creates this type of function that prepares email logs for tests.
type ExpectedEmails func(entries ...*EmailMeta) []*EmailMeta

// ExpectedEmailsFactory returns a function that makes it easier to create email logs
// to assert expected email behavior against. This factory function will usually be
// used to populate a test suite.
func ExpectedEmailsFactory(serviceEmail string) ExpectedEmails {
	return func(entries ...*EmailMeta) []*EmailMeta {
		for _, entry := range entries {
			// Set From email address to the service email from the configuration
			if entry.From == "" {
				entry.From = serviceEmail
			}

			// Set the To email address from the contact
			if entry.To == "" {
				entry.To = entry.Contact.Email.Email
			}

			// Map subject to reason if reason is not populated
			if entry.Reason == "" && entry.Subject != "" {
				entry.Reason = Subject2Reason[entry.Subject]
			}

			// Map reason to subject if subject is not populated
			if entry.Subject == "" && entry.Reason != "" {
				entry.Subject = Reason2Subject[entry.Reason]
			}
		}
		return entries
	}
}

// CheckEmails verifies that the provided email messages exist in both the email mock
// and the audit log on the contact, if the email was sent to a contact. This method is
// meant to be run from a test context.
// TODO: refactor to expect multiple emails per contact/recipient
func CheckEmails(t *testing.T, messages []*EmailMeta) {
	var sentEmails []*sgmail.SGMailV3

	// Check total number of emails sent
	require.Len(t, mock.Emails, len(messages), "incorrect number of emails sent")

	// Get emails from the mock
	for _, data := range mock.Emails {
		msg := &sgmail.SGMailV3{}
		require.NoError(t, json.Unmarshal(data, msg))
		sentEmails = append(sentEmails, msg)
	}

	for i, msg := range messages {
		// If the email was sent to a contact, check the sent email log
		if msg.Contact != nil {
			log := msg.Contact.Logs()
			require.GreaterOrEqual(t, len(log), 1, "contact %s is expected to have at least one email log", msg.Contact.Email)
			require.Equal(t, msg.Reason, log[0].Reason)

			ts, err := time.Parse(time.RFC3339, log[0].Timestamp)
			require.NoError(t, err)
			require.True(t, ts.Sub(msg.Timestamp) < time.Minute, "timestamp in email log is too old")
		}

		expectedRecipient, err := mail.ParseAddress(msg.To)
		require.NoError(t, err)

		// Search for the sent email in the mock and check the metadata
		found := false
		for _, sent := range sentEmails {
			recipient, err := emailutils.GetRecipient(sent)
			require.NoError(t, err)
			if recipient == expectedRecipient.Address {
				found = true
				sender, err := mail.ParseAddress(msg.From)
				require.NoError(t, err)
				require.Equal(t, sender.Address, sentEmails[i].From.Address)
				require.Equal(t, msg.Subject, sentEmails[i].Subject)
				break
			}
		}
		require.True(t, found, "email not sent for recipient %s", msg.To)
	}
}
