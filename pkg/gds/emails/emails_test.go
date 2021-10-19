package emails_test

import (
	"encoding/json"
	"net/mail"
	"net/url"
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

	vcdata := emails.VerifyContactData{Name: recipient, Token: "abcdef1234567890", VID: "42", BaseURL: "http://localhost:8080/verify-contact"}
	mail, err := emails.VerifyContactEmail(sender, senderEmail, recipient, recipientEmail, vcdata)
	require.NoError(t, err)
	require.Equal(t, emails.VerifyContactRE, mail.Subject)

	rrdata := emails.ReviewRequestData{Request: "foo", Token: "abcdef1234567890", VID: "42", BaseURL: "http://localhost:8081/vasps/"}
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
	require.Empty(t, data.VerifyContactURL(), "if no base url is provided, VerifyContactURL() should return empty string")

	data = emails.VerifyContactData{
		Name:    "Darlene Ulmsted",
		Token:   "1234defg4321",
		VID:     "42",
		BaseURL: "http://localhost:8080/verify-contact",
	}
	link, err := url.Parse(data.VerifyContactURL())
	require.NoError(t, err)
	require.Equal(t, "http", link.Scheme)
	require.Equal(t, "localhost:8080", link.Host)
	require.Equal(t, "/verify-contact", link.Path)
	params := link.Query()
	require.Equal(t, data.Token, params.Get("token"))
	require.Equal(t, data.VID, params.Get("vaspID"))
}

func TestAdminReviewURL(t *testing.T) {
	data := emails.ReviewRequestData{
		VID:   "42",
		Token: "1234defg4321",
	}
	require.Empty(t, data.AdminReviewURL(), "if no base url provided, AdminReviewURL() should return empty string")

	data = emails.ReviewRequestData{
		VID:     "42",
		Token:   "1234defg4321",
		BaseURL: "http://localhost:8088/vasps/",
	}
	link, err := url.Parse(data.AdminReviewURL())
	require.NoError(t, err)
	require.Equal(t, "http", link.Scheme)
	require.Equal(t, "localhost:8088", link.Host)
	require.Equal(t, "/vasps/42", link.Path)
}

func TestSendEmails(t *testing.T) {
	// NOTE: if you place a .env file in this directory alongside the test file, it
	// will be read, making it simpler to run tests and set environment variables.
	godotenv.Load()

	// This test uses the environment to send rendered emails with context specific
	// emails - this is to test the rendering of the emails with data only; it does not
	// go through any of the server workflow for generating tokens, etc.
	var conf config.EmailConfig
	err := envconfig.Process("gds", &conf)
	require.NoError(t, err)

	if !conf.EmailTesting {
		t.Skip("skip generate and send emails test")
	}

	// This test mocks the SendGrid method for sending emails and stores them to an
	// in-memory data structure instead.
	email, err := emails.New(conf)
	require.NoError(t, err)

	sender, err := mail.ParseAddress(conf.ServiceEmail)
	require.NoError(t, err)

	receipient, err := mail.ParseAddress(conf.AdminEmail)
	require.NoError(t, err)

	var expected []byte
	vcdata := emails.VerifyContactData{Name: receipient.Name, Token: "Hk79ZIhCSrYJtSaaMECZZKI1BtsCY9zDLPq9c1amyK2zJY6T", VID: "9e069e01-8515-4d57-b9a5-e249f7ab4fca", BaseURL: "http://localhost:3000/verify-contact"}
	msg, err := emails.VerifyContactEmail(sender.Name, sender.Address, receipient.Name, receipient.Address, vcdata)
	require.NoError(t, err)
	require.NoError(t, email.Send(msg))
	require.Len(t, emails.MockEmails, 1)
	expected, err = json.Marshal(msg)
	require.NoError(t, err)
	require.Equal(t, expected, emails.MockEmails[0])

	rrdata := emails.ReviewRequestData{Request: "foo", Token: "abcdef1234567890", VID: "42", Attachment: []byte(`{"hello": "world"}`)}
	msg, err = emails.ReviewRequestEmail(sender.Name, sender.Address, receipient.Name, receipient.Address, rrdata)
	require.NoError(t, err)
	require.NoError(t, email.Send(msg))
	require.Len(t, emails.MockEmails, 2)
	expected, err = json.Marshal(msg)
	require.NoError(t, err)
	require.Equal(t, expected, emails.MockEmails[1])

	rjdata := emails.RejectRegistrationData{Name: receipient.Name, Reason: "not a good time", VID: "42"}
	msg, err = emails.RejectRegistrationEmail(sender.Name, sender.Address, receipient.Name, receipient.Address, rjdata)
	require.NoError(t, err)
	require.NoError(t, email.Send(msg))
	require.Len(t, emails.MockEmails, 3)
	expected, err = json.Marshal(msg)
	require.NoError(t, err)
	require.Equal(t, expected, emails.MockEmails[2])

	dcdata := emails.DeliverCertsData{Name: receipient.Name, VID: "42", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443"}
	msg, err = emails.DeliverCertsEmail(sender.Name, sender.Address, receipient.Name, receipient.Address, "testdata/foo.zip", dcdata)
	require.NoError(t, err)
	require.NoError(t, email.Send(msg))
	require.Len(t, emails.MockEmails, 4)
	expected, err = json.Marshal(msg)
	require.NoError(t, err)
	require.Equal(t, expected, emails.MockEmails[3])

	t.Cleanup(emails.PurgeMockEmails)
}
