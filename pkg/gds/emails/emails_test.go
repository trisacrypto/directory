package emails_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
)

// If the eyeball flag is set, then the tests will write MIME emails to the testdata directory.
var eyeball = flag.Bool("eyeball", false, "Generate MIME emails for eyeball testing")

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
		err := emails.WriteMIME(msg, filepath.Join("testdata", fmt.Sprintf("eyeball%s", t.Name()), name))
		require.NoError(t, err)
	}
}

func TestEmailBuilders(t *testing.T) {
	var (
		sender         = "Lewis Hudson"
		senderEmail    = "lewis@example.com"
		recipient      = "Rachel Lendt"
		recipientEmail = "rachel@example.com"
	)

	setupMIMEDir(t)

	vcdata := emails.VerifyContactData{Name: recipient, Token: "abcdef1234567890", VID: "42", BaseURL: "http://localhost:8080/verify-contact"}
	mail, err := emails.VerifyContactEmail(sender, senderEmail, recipient, recipientEmail, vcdata)
	require.NoError(t, err)
	require.Equal(t, emails.VerifyContactRE, mail.Subject)
	generateMIME(t, mail, "verify-contact.mim")

	rrdata := emails.ReviewRequestData{Request: "foo", Token: "abcdef1234567890", VID: "42", BaseURL: "http://localhost:8081/vasps/"}
	mail, err = emails.ReviewRequestEmail(sender, senderEmail, recipient, recipientEmail, rrdata)
	require.NoError(t, err)
	require.Equal(t, emails.ReviewRequestRE, mail.Subject)
	generateMIME(t, mail, "review-request.mim")

	rjdata := emails.RejectRegistrationData{Name: recipient, Reason: "not a good time", VID: "42"}
	mail, err = emails.RejectRegistrationEmail(sender, senderEmail, recipient, recipientEmail, rjdata)
	require.NoError(t, err)
	require.Equal(t, emails.RejectRegistrationRE, mail.Subject)
	generateMIME(t, mail, "reject-registration.mim")

	dcdata := emails.DeliverCertsData{Name: recipient, VID: "42", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443"}
	mail, err = emails.DeliverCertsEmail(sender, senderEmail, recipient, recipientEmail, "testdata/foo.zip", dcdata)
	require.NoError(t, err)
	require.Equal(t, emails.DeliverCertsRE, mail.Subject)
	generateMIME(t, mail, "deliver-certs.mim")
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

// This suite mocks the SendGrid email client to verify that email metadata is
// populated corectly and emails can be marshaled into bytes for transmission.
func TestEmailSuite(t *testing.T) {
	suite.Run(t, &EmailTestSuite{})
}

type EmailTestSuite struct {
	suite.Suite
	conf config.EmailConfig
}

func (suite *EmailTestSuite) SetupSuite() {
	suite.conf = config.EmailConfig{
		Testing:      true,
		ServiceEmail: "service@example.com",
		AdminEmail:   "admin@example.com",
		Storage:      "fixtures/emails",
	}
}

func (suite *EmailTestSuite) BeforeTest(suiteName, testName string) {
	setupMIMEDir(suite.T())
}

func (suite *EmailTestSuite) AfterTest(suiteName, testName string) {
	emails.PurgeMockEmails()
}

func (suite *EmailTestSuite) TestSendVerifyContactEmail() {
	// Load the test suite config
	require := suite.Require()
	sender, err := mail.ParseAddress(suite.conf.ServiceEmail)
	require.NoError(err)
	recipient, err := mail.ParseAddress(suite.conf.AdminEmail)
	require.NoError(err)

	// Init the mocked SendGrid client
	email, err := emails.New(suite.conf)
	require.NoError(err)

	data := emails.VerifyContactData{Name: recipient.Name, Token: "Hk79ZIhCSrYJtSaaMECZZKI1BtsCY9zDLPq9c1amyK2zJY6T", VID: "9e069e01-8515-4d57-b9a5-e249f7ab4fca", BaseURL: "http://localhost:3000/verify-contact"}
	msg, err := emails.VerifyContactEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, data)
	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(emails.MockEmails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, emails.MockEmails[0])

	generateMIME(suite.T(), msg, "verify-contact.mim")
}

func (suite *EmailTestSuite) TestSendReviewRequestEmail() {
	// Load the test suite config
	require := suite.Require()
	sender, err := mail.ParseAddress(suite.conf.ServiceEmail)
	require.NoError(err)
	recipient, err := mail.ParseAddress(suite.conf.AdminEmail)
	require.NoError(err)

	// Init the mocked SendGrid client
	email, err := emails.New(suite.conf)
	require.NoError(err)

	data := emails.ReviewRequestData{Request: "foo", Token: "abcdef1234567890", VID: "42", Attachment: []byte(`{"hello": "world"}`)}
	msg, err := emails.ReviewRequestEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, data)
	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(emails.MockEmails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, emails.MockEmails[0])

	generateMIME(suite.T(), msg, "review-request.mim")
}

func (suite *EmailTestSuite) TestSendRejectRegistrationEmail() {
	// Load the test suite config
	require := suite.Require()
	sender, err := mail.ParseAddress(suite.conf.ServiceEmail)
	require.NoError(err)
	recipient, err := mail.ParseAddress(suite.conf.AdminEmail)
	require.NoError(err)

	// Init the mocked SendGrid client
	email, err := emails.New(suite.conf)
	require.NoError(err)

	data := emails.RejectRegistrationData{Name: recipient.Name, Reason: "not a good time", VID: "42"}
	msg, err := emails.RejectRegistrationEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, data)
	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(emails.MockEmails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, emails.MockEmails[0])

	generateMIME(suite.T(), msg, "reject-registration.mim")
}

func (suite *EmailTestSuite) TestSendDeliverCertsEmail() {
	// Load the test suite config
	require := suite.Require()
	sender, err := mail.ParseAddress(suite.conf.ServiceEmail)
	require.NoError(err)
	recipient, err := mail.ParseAddress(suite.conf.AdminEmail)
	require.NoError(err)

	// Init the mocked SendGrid client
	email, err := emails.New(suite.conf)
	require.NoError(err)

	data := emails.DeliverCertsData{Name: recipient.Name, VID: "42", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443"}
	msg, err := emails.DeliverCertsEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, "testdata/foo.zip", data)
	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(emails.MockEmails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, emails.MockEmails[0])

	generateMIME(suite.T(), msg, "deliver-certs.mim")
}
