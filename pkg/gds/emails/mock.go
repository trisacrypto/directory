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
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

type EmailMeta struct {
	Contact   *pb.Contact
	To        string
	From      string
	Subject   string
	Reason    string
	Timestamp time.Time
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
			log, err := models.GetEmailLog(msg.Contact)
			require.NoError(t, err)
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
