package emails_test

import (
	"net/mail"
	"net/url"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
)

func TestEmailBuilders(t *testing.T) {
	var (
		sender         = "Lewis Hudson"
		senderEmail    = "lewis@example.com"
		recipient      = "Rachel Lendt"
		recipientEmail = "rachel@example.com"
	)

	vcdata := emails.VerifyContactData{Name: recipient, Token: "abcdef1234567890", VID: "42"}
	mail, err := emails.VerifyContactEmail(sender, senderEmail, recipient, recipientEmail, vcdata)
	require.NoError(t, err)
	require.Equal(t, emails.VerifyContactRE, mail.Subject)

	rrdata := emails.ReviewRequestData{Request: "foo", Token: "abcdef1234567890", VID: "42"}
	mail, err = emails.ReviewRequestEmail(sender, senderEmail, recipient, recipientEmail, rrdata)
	require.NoError(t, err)
	require.Equal(t, emails.ReviewRequestRE, mail.Subject)

	rjdata := emails.RejectRegistrationData{Name: recipient, Reason: "not a good time", VID: "42"}
	mail, err = emails.RejectRegistrationEmail(sender, senderEmail, recipient, recipientEmail, rjdata)
	require.NoError(t, err)
	require.Equal(t, emails.RejectRegistrationRE, mail.Subject)

	dcdata := emails.DeliverCertsData{Name: recipient, VID: "42", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443"}
	mail, err = emails.DeliverCertsEmail(sender, senderEmail, recipient, recipientEmail, "testdata/foo.zip", dcdata)
	require.NoError(t, err)
	require.Equal(t, emails.DeliverCertsRE, mail.Subject)
}

func TestVerifyContactURL(t *testing.T) {
	data := emails.VerifyContactData{
		Name:  "Darlene Ulmsted",
		Token: "1234defg4321",
		VID:   "42",
	}
	link, err := url.Parse(data.VerifyContactURL())
	require.NoError(t, err)
	require.Equal(t, "https", link.Scheme)
	require.Equal(t, "vaspdirectory.net", link.Host)
	require.Equal(t, "/verify-contact", link.Path)
	params := link.Query()
	require.Equal(t, data.Token, params.Get("token"))
	require.Equal(t, data.VID, params.Get("vaspID"))

	data = emails.VerifyContactData{
		Name:  "Darlene Ulmsted",
		Token: "1234defg4321",
		VID:   "42",
		URL:   "http://localhost:8080/verify-contact",
	}
	link, err = url.Parse(data.VerifyContactURL())
	require.NoError(t, err)
	require.Equal(t, "http", link.Scheme)
	require.Equal(t, "localhost:8080", link.Host)
	require.Equal(t, "/verify-contact", link.Path)
	params = link.Query()
	require.Equal(t, data.Token, params.Get("token"))
	require.Equal(t, data.VID, params.Get("vaspID"))
}

func TestSendEmails(t *testing.T) {
	// NOTE: if you place a .env file in this directory alongside the test file, it
	// will be read, making it simpler to run tests and set environment variables.
	godotenv.Load()

	if os.Getenv("GDS_TEST_SENDING_EMAILS") == "" {
		t.Skip("skip generate and send emails test")
	}

	// This test uses the environment to send rendered emails with context specific
	// emails - this is to test the rendering of the emails with data only; it does not
	// go through any of the server workflow for generating tokens, etc.
	var conf config.EmailConfig
	err := envconfig.Process("gds", &conf)
	require.NoError(t, err)

	// This test sends emails from the serviceEmail using SendGrid to the adminsEmail
	email, err := emails.New(conf)
	require.NoError(t, err)

	sender, err := mail.ParseAddress(conf.ServiceEmail)
	require.NoError(t, err)

	receipient, err := mail.ParseAddress(conf.AdminEmail)
	require.NoError(t, err)

	vcdata := emails.VerifyContactData{Name: receipient.Name, Token: "Hk79ZIhCSrYJtSaaMECZZKI1BtsCY9zDLPq9c1amyK2zJY6T", VID: "9e069e01-8515-4d57-b9a5-e249f7ab4fca", URL: "http://localhost:3000/verify-contact"}
	msg, err := emails.VerifyContactEmail(sender.Name, sender.Address, receipient.Name, receipient.Address, vcdata)
	require.NoError(t, err)
	require.NoError(t, email.Send(msg))

	rrdata := emails.ReviewRequestData{Request: "foo", Token: "abcdef1234567890", VID: "42"}
	msg, err = emails.ReviewRequestEmail(sender.Name, sender.Address, receipient.Name, receipient.Address, rrdata)
	require.NoError(t, err)
	require.NoError(t, email.Send(msg))

	rjdata := emails.RejectRegistrationData{Name: receipient.Name, Reason: "not a good time", VID: "42"}
	msg, err = emails.RejectRegistrationEmail(sender.Name, sender.Address, receipient.Name, receipient.Address, rjdata)
	require.NoError(t, err)
	require.NoError(t, email.Send(msg))

	dcdata := emails.DeliverCertsData{Name: receipient.Name, VID: "42", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443"}
	msg, err = emails.DeliverCertsEmail(sender.Name, sender.Address, receipient.Name, receipient.Address, "testdata/foo.zip", dcdata)
	require.NoError(t, err)
	require.NoError(t, email.Send(msg))
}
