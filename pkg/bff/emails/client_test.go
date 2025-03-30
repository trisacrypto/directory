package emails_test

import (
	"net/mail"
	"net/url"
	"os"
	"testing"

	"github.com/auth0/go-auth0/management"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/config"
	emails "github.com/trisacrypto/directory/pkg/bff/emails"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/emails/mock"
)

func TestLiveEmailTests(t *testing.T) {
	// NOTE: if you place a .env file in this directory alongside the test file, it
	// will be read, making it simpler to run tests and set environment variables.
	godotenv.Load()

	if os.Getenv("GDS_BFF_TEST_SENDING_CLIENT_EMAILS") == "" {
		t.Skip("skip client send emails test")
	}

	recipientEmail := os.Getenv("GDS_BFF_TEST_RECIPIENT_EMAIL")
	require.NotEmpty(t, recipientEmail, "GDS_BFF_TEST_RECIPIENT_EMAIL must be set for live email tests")

	recipient, err := mail.ParseAddress(recipientEmail)
	require.NoError(t, err, "could not parse recipient email address from environment")

	// This test uses the environment to actually send rendered emails with specific
	// context data.
	var conf config.EmailConfig
	require.NoError(t, envconfig.Process("gds_bff", &conf), "could not process email config")

	// This test sends emails from the serviceEmail using SendGrid to the email address
	// specified in the environment.
	email, err := emails.New(conf)
	require.NoError(t, err, "could not create email manager")

	// Send a user invite email
	inviteName := "Alice Ables"
	inviteEmail := "alice@example.com"
	inviteURL, err := url.Parse("https://testnet.directory/invite/1234")
	require.NoError(t, err, "could not parse invite URL")

	user := &management.User{
		Email: &recipient.Address,
		Name:  &recipient.Name,
	}
	inviter := &management.User{
		Email: &inviteEmail,
		Name:  &inviteName,
	}
	org := &models.Organization{
		Name: "Alice VASP",
	}
	err = email.SendUserInvite(user, inviter, org, inviteURL)
	require.NoError(t, err, "could not send user invite email")
}

func (s *EmailTestSuite) TestSendUserInvite() {
	require := s.Require()

	// Init the mocked SendGrid client
	email, err := emails.New(s.conf)
	require.NoError(err, "could not create email manager")

	var userName, userEmail, inviterName, inviterEmail string

	user := &management.User{
		Name:  &userName,
		Email: &userEmail,
	}

	inviter := &management.User{
		Name:  &inviterName,
		Email: &inviterEmail,
	}

	org := &models.Organization{
		Name: "Team Sonic",
	}

	inviteURL, err := url.Parse("https://gottagofast.com/invite/1234")
	require.NoError(err, "could not parse invite URL")

	// Should return an error if user has no email address
	userName = "Tails the Fox"
	userEmail = ""
	inviterName = "Sonic the Hedgehog"
	inviterEmail = "sonic@gottagofast.com"
	err = email.SendUserInvite(user, inviter, org, inviteURL)
	require.EqualError(err, "user has no email address", "should return an error if user has no email address")

	// Should return an error if inviter has no name or email address
	userEmail = "tails@gottagofast.com"
	inviterName = ""
	inviterEmail = ""
	err = email.SendUserInvite(user, inviter, org, inviteURL)
	require.EqualError(err, "inviter has no email address", "should return an error if inviter has no name or email address")

	// If no user name is provided, the email address is used instead
	userName = ""
	userEmail = "tails@gottagofast.com"
	inviterEmail = "sonic@gottagofast.com"
	err = email.SendUserInvite(user, inviter, org, inviteURL)
	require.NoError(err, "could not send user invite email")
	require.Len(mock.Emails, 1)

	// Does not error if no organization name is provided
	org.Name = ""
	err = email.SendUserInvite(user, inviter, org, inviteURL)
	require.NoError(err, "could not send user invite email")
	require.Len(mock.Emails, 2)

	// Does not error if all fields are provided
	userName = "Tails the Fox"
	userEmail = "tails@gottagofast.com"
	inviterName = "Sonic the Hedgehog"
	inviterEmail = "sonic@gottagofast.com"
	org.Name = "Team Sonic"
	err = email.SendUserInvite(user, inviter, org, inviteURL)
	require.NoError(err, "could not send user invite email")
	require.Len(mock.Emails, 3)
}
