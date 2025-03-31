package emails_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"testing"

	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/bff/config"
	emails "github.com/trisacrypto/directory/pkg/bff/emails"
	emailutils "github.com/trisacrypto/directory/pkg/utils/emails"
	"github.com/trisacrypto/directory/pkg/utils/emails/mock"
	"github.com/trisacrypto/directory/pkg/utils/logger"
)

// If the eyeball flag is set, then the tests will write MIME emails to the testdata directory.
var eyeball = flag.Bool("eyeball", true, "Generate MIME emails for eyeball testing")

// Creates a directory for the MIME emails if the eyeball flag is set.
// If the eyeball flag is set, this will also purge the existing eyeball directory first.
func setupMIMEDir(t *testing.T) {
	if *eyeball {
		path := filepath.Join("testdata", fmt.Sprintf("eyeball%s", t.Name()))
		err := os.RemoveAll(path)
		require.NoError(t, err)
		err = os.MkdirAll(path, 0755)
		require.NoError(t, err)
	}
}

// generateMIME writes an SGMailV3 email to a MIME file for manual inspection if the eyeball flag is set.
func generateMIME(t *testing.T, msg *sgmail.SGMailV3, name string) {
	if *eyeball {
		err := emailutils.WriteMIME(msg, filepath.Join("testdata", fmt.Sprintf("eyeball%s", t.Name()), name))
		require.NoError(t, err, "failed to write MIME email")
	}
}

func TestEmailBuilders(t *testing.T) {
	var (
		sender         = "Sonic the Hedgehog"
		senderEmail    = "sonic@gottagofast.com"
		recipient      = "Tails the Fox"
		recipientEmail = "tails@gottagofast.com"
	)

	setupMIMEDir(t)

	// Test rendering with user names
	inviteData := emails.InviteUserData{
		UserName:     recipient,
		UserEmail:    recipientEmail,
		InviterName:  sender,
		InviterEmail: senderEmail,
		Organization: "Team Sonic",
		InviteURL:    "https://gottagofast.com/invite",
	}
	mail, err := emails.InviteUserEmail(sender, senderEmail, recipient, recipientEmail, inviteData)
	require.NoError(t, err, "failed to create user invite email")
	require.Equal(t, "Sonic the Hedgehog has invited you to collaborate on Team Sonic", mail.Subject, "user invite email subject is incorrect")
	generateMIME(t, mail, "invite_user_with_name.mime")

	// Test rendering with email addresses
	inviteData.UserName = ""
	inviteData.InviterName = ""
	mail, err = emails.InviteUserEmail(sender, senderEmail, recipient, recipientEmail, inviteData)
	require.NoError(t, err, "failed to create user invite email")
	require.Equal(t, "You have been invited to collaborate on Team Sonic", mail.Subject, "user invite email subject is incorrect")
	generateMIME(t, mail, "invite_user_with_email.mime")

	// Test rendering with no organization name
	inviteData.Organization = ""
	mail, err = emails.InviteUserEmail(sender, senderEmail, recipient, recipientEmail, inviteData)
	require.NoError(t, err, "failed to create user invite email")
	require.Equal(t, "You have been invited to collaborate on an organization", mail.Subject, "user invite email subject is incorrect")
	generateMIME(t, mail, "invite_user_no_org.mime")

	// Test rendering with inviter name but no organization name
	inviteData.InviterName = "Sonic the Hedgehog"
	mail, err = emails.InviteUserEmail(sender, senderEmail, recipient, recipientEmail, inviteData)
	require.NoError(t, err, "failed to create user invite email")
	require.Equal(t, "Sonic the Hedgehog has invited you to collaborate on their organization", mail.Subject, "user invite email subject is incorrect")
	generateMIME(t, mail, "invite_user_with_name_no_org.mime")
}

func (s *EmailTestSuite) TestUserInviteEmail() {
	require := s.Require()
	service, err := mail.ParseAddress(s.conf.ServiceEmail)
	require.NoError(err, "could not parse service email address")
	inviter, err := mail.ParseAddress("Sonic the Hedgehog <sonic@gottagofast.com>")
	require.NoError(err, "could not parse inviter email address")
	recipient, err := mail.ParseAddress("Tails the Fox <tails@gottagofast.com>")
	require.NoError(err, "could not parse email address")

	// Init the mocked SendGrid client
	email, err := emails.New(s.conf)
	require.NoError(err, "could not create email manager")

	data := emails.InviteUserData{
		UserName:     recipient.Name,
		UserEmail:    recipient.Address,
		InviterName:  inviter.Name,
		InviterEmail: inviter.Address,
		Organization: "Team Sonic",
		InviteURL:    "https://gottagofast.com/invite",
	}
	msg, err := emails.InviteUserEmail(service.Name, service.Address, recipient.Name, recipient.Address, data)
	require.NoError(err, "could not create user invite email")
	require.NoError(email.Send(msg), "could not send user invite email")
	require.Len(mock.Emails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err, "could not marshal invite email into bytes")
	require.Equal(expected, mock.Emails[0], "user email did not match expected")

	generateMIME(s.T(), msg, "invite_user.mime")
}

// This suite mocks the SendGrid email client to verify that email metadata is
// populated corectly and emails can be marshaled into bytes for transmission.
func TestEmailSuite(t *testing.T) {
	suite.Run(t, &EmailTestSuite{})
}

type EmailTestSuite struct {
	suite.Suite
	conf config.EmailConfig
}

func (s *EmailTestSuite) SetupSuite() {
	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overriden
	logger.Discard()

	s.conf = config.EmailConfig{
		Testing:      true,
		ServiceEmail: "TRISA Directory Service <admin@trisa.directory>",
		Storage:      "fixtures/emails",
	}
}

func (s *EmailTestSuite) BeforeTest(suiteName, testName string) {
	setupMIMEDir(s.T())
}

func (s *EmailTestSuite) AfterTest(suiteName, testName string) {
	mock.PurgeEmails()
}

func (s *EmailTestSuite) TearDownSuite() {
	logger.ResetLogger()
}
