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
	"time"

	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	emailutils "github.com/trisacrypto/directory/pkg/utils/emails"
	"github.com/trisacrypto/directory/pkg/utils/emails/mock"
	"github.com/trisacrypto/directory/pkg/utils/logger"
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
		err := emailutils.WriteMIME(msg, filepath.Join("testdata", fmt.Sprintf("eyeball%s", t.Name()), name))
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

	vcdata := emails.VerifyContactData{Name: recipient, Token: "abcdef1234567890", VID: "42", BaseURL: "http://localhost:8080/verify", DirectoryID: "testnet.io"}
	mail, err := emails.VerifyContactEmail(sender, senderEmail, recipient, recipientEmail, vcdata)
	require.NoError(t, err)
	require.Equal(t, emails.VerifyContactRE, mail.Subject)
	generateMIME(t, mail, "verify-contact.mim")

	rrdata := emails.ReviewRequestData{Request: "foo", Token: "abcdef1234567890", VID: "42", BaseURL: "http://localhost:8081/vasps/"}
	mail, err = emails.ReviewRequestEmail(sender, senderEmail, recipient, recipientEmail, rrdata)
	require.NoError(t, err)
	require.Equal(t, emails.ReviewRequestRE, mail.Subject)
	generateMIME(t, mail, "review-request.mim")

	rjdata := emails.RejectRegistrationData{Name: recipient, Reason: "not a good time", VID: "42", CommonName: "example.com", Organization: "Acme, Inc.", RegisteredDirectory: "trisatest.dev"}
	mail, err = emails.RejectRegistrationEmail(sender, senderEmail, recipient, recipientEmail, rjdata)
	require.NoError(t, err)
	require.Equal(t, emails.RejectRegistrationRE, mail.Subject)
	generateMIME(t, mail, "reject-registration.mim")

	dcdata := emails.DeliverCertsData{Name: recipient, VID: "42", Organization: "Acme, Inc", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443"}
	mail, err = emails.DeliverCertsEmail(sender, senderEmail, recipient, recipientEmail, "testdata/foo.zip", dcdata)
	require.NoError(t, err)
	require.Equal(t, emails.DeliverCertsRE, mail.Subject)
	generateMIME(t, mail, "deliver-certs.mim")

	expires := time.Date(2022, time.July, 18, 12, 11, 35, 0, time.UTC)
	reissuance := time.Date(2022, time.July, 25, 14, 0, 0, 0, time.UTC)

	eandata := emails.ExpiresAdminNotificationData{VID: "42", Organization: "Acme, Inc", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443", RegisteredDirectory: "trisatest.net", Expiration: expires, Reissuance: reissuance, BaseURL: "http://localhost:8081/vasps/"}
	mail, err = emails.ExpiresAdminNotificationEmail(sender, senderEmail, recipient, recipientEmail, eandata)
	require.NoError(t, err)
	require.Equal(t, emails.ExpiresAdminNotificationRE, mail.Subject, "incorrect subject")
	generateMIME(t, mail, "expires-admin-notification.mim")

	rmdata := emails.ReissuanceReminderData{Name: recipient, VID: "42", Organization: "Acme, Inc", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443", RegisteredDirectory: "trisatest.net", Expiration: expires, Reissuance: reissuance}
	mail, err = emails.ReissuanceReminderEmail(sender, senderEmail, recipient, recipientEmail, rmdata)
	require.NoError(t, err)
	require.Equal(t, emails.ReissuanceReminderRE, mail.Subject, "incorrect subject")
	generateMIME(t, mail, "reissuance-reminder.mim")

	rsdata := emails.ReissuanceStartedData{Name: recipient, VID: "42", Organization: "Acme, Inc", CommonName: "example.com", Endpoint: "trisa.example.com:443", RegisteredDirectory: "trisatest.net", WhisperURL: "http://localhost/secret"}
	mail, err = emails.ReissuanceStartedEmail(sender, senderEmail, recipient, recipientEmail, rsdata)
	require.NoError(t, err)
	require.Equal(t, emails.ReissuanceStartedRE, mail.Subject, "incorrect subject")
	generateMIME(t, mail, "reissuance-started.mim")

	randata := emails.ReissuanceAdminNotificationData{VID: "42", Organization: "Acme, Inc", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443", RegisteredDirectory: "trisatest.net", Expiration: expires, Reissuance: reissuance, BaseURL: "http://localhost:8081/vasps/"}
	mail, err = emails.ReissuanceAdminNotificationEmail(sender, senderEmail, recipient, recipientEmail, randata)
	require.NoError(t, err)
	require.Equal(t, emails.ReissuanceAdminNotificationRE, mail.Subject, "incorrect subject")
	generateMIME(t, mail, "reissuance-admin-notification.mim")
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
		BaseURL: "http://localhost:8080/verify",
	}
	link := data.VerifyContactURL()
	require.Equal(t, "http", link.Scheme)
	require.Equal(t, "localhost:8080", link.Host)
	require.Equal(t, "/verify", link.Path)
	params := link.Query()
	require.Equal(t, data.Token, params.Get("token"))
	require.Equal(t, data.VID, params.Get("vaspID"))
	require.Equal(t, data.DirectoryID, params.Get("registered_directory"))
}

func TestAdminReviewURL(t *testing.T) {
	emptyCases := []emails.AdminReview{
		emails.ReviewRequestData{
			VID:   "42",
			Token: "1234defg4321",
		},
		emails.ExpiresAdminNotificationData{
			VID:        "42",
			CommonName: "test.example.com",
		},
	}

	for _, tc := range emptyCases {
		require.Empty(t, tc.AdminReviewURL(), "if no base url provided, AdminReviewURL() should return empty string")
	}

	testCases := []emails.AdminReview{
		emails.ReviewRequestData{
			VID:     "42",
			Token:   "1234defg4321",
			BaseURL: "http://localhost:8088/vasps/",
		},
		emails.ExpiresAdminNotificationData{
			VID:        "42",
			CommonName: "test.example.com",
			BaseURL:    "http://localhost:8088/vasps/",
		},
	}

	for _, tc := range testCases {
		link, err := url.Parse(tc.AdminReviewURL())
		require.NoError(t, err)
		require.Equal(t, "http", link.Scheme)
		require.Equal(t, "localhost:8088", link.Host)
		require.Equal(t, "/vasps/42", link.Path)
	}
}

func TestDateStrings(t *testing.T) {
	emptyCases := []emails.FutureReissuer{
		emails.ExpiresAdminNotificationData{
			VID:        "42",
			CommonName: "test.example.com",
		},
		emails.ReissuanceReminderData{
			VID:        "42",
			CommonName: "test.example.com",
		},
	}

	for _, tc := range emptyCases {
		require.Equal(t, emails.UnknownDate, tc.ExpirationDate(), "expected expiration date to be unknown when no expiration timestamp")
		require.Equal(t, emails.UnknownDate, tc.ReissueDate(), "expected reissue date to be unknown when no reissuance timestamp")
	}

	// Create some timestamps
	expires, err := time.Parse(time.RFC3339, "2022-07-18T16:28:51-05:00")
	require.NoError(t, err, "could not parse timestamp fixture")

	reissue, err := time.Parse(time.RFC3339, "2022-07-25T12:30:00-05:00")
	require.NoError(t, err, "could not parse timestamp fixture")

	testCases := []emails.FutureReissuer{
		emails.ExpiresAdminNotificationData{
			VID:        "42",
			CommonName: "test.example.com",
			Expiration: expires,
			Reissuance: reissue,
		},
		emails.ReissuanceReminderData{
			VID:        "42",
			CommonName: "test.example.com",
			Expiration: expires,
			Reissuance: reissue,
		},
	}

	for _, tc := range testCases {
		require.Equal(t, "Monday, July 18, 2022", tc.ExpirationDate(), "expected expiration date to be formated correctly with timestamp")
		require.Equal(t, "Monday, July 25, 2022", tc.ReissueDate(), "expected reissuance date to be formated correctly with timestamp")
	}
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
	// Discard logging from the application to focus on test logs
	// NOTE: ConsoleLog MUST be false otherwise this will be overriden
	logger.Discard()

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
	mock.PurgeEmails()
}

func (suite *EmailTestSuite) TearDownSuite() {
	logger.ResetLogger()
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

	data := emails.VerifyContactData{Name: recipient.Name, Token: "Hk79ZIhCSrYJtSaaMECZZKI1BtsCY9zDLPq9c1amyK2zJY6T", VID: "9e069e01-8515-4d57-b9a5-e249f7ab4fca", BaseURL: "http://localhost:3000/verify", DirectoryID: "testnet.io"}
	msg, err := emails.VerifyContactEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, data)
	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(mock.Emails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, mock.Emails[0])

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
	require.Len(mock.Emails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, mock.Emails[0])

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

	data := emails.RejectRegistrationData{Name: recipient.Name, Reason: "not a good time", VID: "42", Organization: "Acme, Inc.", CommonName: "example.com", RegisteredDirectory: "trisatest.dev"}
	msg, err := emails.RejectRegistrationEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, data)
	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(mock.Emails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, mock.Emails[0])

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

	data := emails.DeliverCertsData{Name: recipient.Name, VID: "42", Organization: "Acme, Inc.", CommonName: "example.com", SerialNumber: "1234abcdef56789", Endpoint: "trisa.example.com:443"}
	msg, err := emails.DeliverCertsEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, "testdata/foo.zip", data)
	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(mock.Emails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, mock.Emails[0])

	generateMIME(suite.T(), msg, "deliver-certs.mim")
}

func (suite *EmailTestSuite) TestSendExpiresAdminNotificationEmail() {
	// Load the test suite config
	require := suite.Require()
	sender, err := mail.ParseAddress(suite.conf.ServiceEmail)
	require.NoError(err)
	recipient, err := mail.ParseAddress(suite.conf.AdminEmail)
	require.NoError(err)

	// Init the mocked SendGrid client
	email, err := emails.New(suite.conf)
	require.NoError(err)

	// Create the expires admin notification email.
	data := emails.ExpiresAdminNotificationData{
		VID:                 "42",
		Organization:        "Acme, Inc.",
		CommonName:          "test.example.com",
		SerialNumber:        "1234abcdef56789",
		Endpoint:            "test.example.com:443",
		RegisteredDirectory: "trisatest.net",
		Expiration:          time.Date(2022, time.July, 18, 16, 38, 38, 0, time.Local),
		Reissuance:          time.Date(2022, time.July, 25, 12, 30, 0, 0, time.Local),
		BaseURL:             "http://localhost:8080/vasps",
	}
	msg, err := emails.ExpiresAdminNotificationEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, data)

	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(mock.Emails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, mock.Emails[0])

	generateMIME(suite.T(), msg, "expires-admin-notification.mim")
}

func (suite *EmailTestSuite) TestSendReissuanceReminderEmail() {
	// Load the test suite config
	require := suite.Require()
	sender, err := mail.ParseAddress(suite.conf.ServiceEmail)
	require.NoError(err)
	recipient, err := mail.ParseAddress(suite.conf.AdminEmail)
	require.NoError(err)

	// Init the mocked SendGrid client
	email, err := emails.New(suite.conf)
	require.NoError(err)

	// Create the reissuance reminder email.
	data := emails.ReissuanceReminderData{
		Name:                recipient.Name,
		VID:                 "42",
		Organization:        "Acme, Inc.",
		CommonName:          "test.example.com",
		SerialNumber:        "1234abcdef56789",
		Endpoint:            "test.example.com:443",
		RegisteredDirectory: "trisatest.net",
		Expiration:          time.Date(2022, time.July, 18, 16, 38, 38, 0, time.Local),
		Reissuance:          time.Date(2022, time.July, 25, 12, 30, 0, 0, time.Local),
	}
	msg, err := emails.ReissuanceReminderEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, data)

	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(mock.Emails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, mock.Emails[0])

	generateMIME(suite.T(), msg, "reissuance-reminder.mim")
}

func (suite *EmailTestSuite) TestSendReissuanceStartedEmail() {
	// Load the test suite config
	require := suite.Require()
	sender, err := mail.ParseAddress(suite.conf.ServiceEmail)
	require.NoError(err)
	recipient, err := mail.ParseAddress(suite.conf.AdminEmail)
	require.NoError(err)

	// Init the mocked SendGrid client
	email, err := emails.New(suite.conf)
	require.NoError(err)

	// Create the reissuance started email.
	data := emails.ReissuanceStartedData{
		Name:                recipient.Name,
		VID:                 "42",
		Organization:        "Acme, Inc.",
		CommonName:          "test.example.com",
		Endpoint:            "test.example.com:443",
		RegisteredDirectory: "trisatest.net",
		WhisperURL:          "http://whisper.rotational.dev/secret/foo",
	}
	msg, err := emails.ReissuanceStartedEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, data)

	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(mock.Emails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, mock.Emails[0])

	generateMIME(suite.T(), msg, "reissuance-started.mim")
}

func (suite *EmailTestSuite) TestSendReissuanceAdminNotificationEmail() {
	// Load the test suite config
	require := suite.Require()
	sender, err := mail.ParseAddress(suite.conf.ServiceEmail)
	require.NoError(err)
	recipient, err := mail.ParseAddress(suite.conf.AdminEmail)
	require.NoError(err)

	// Init the mocked SendGrid client
	email, err := emails.New(suite.conf)
	require.NoError(err)

	// Create the reissuance admin notification email.
	data := emails.ReissuanceAdminNotificationData{
		VID:                 "42",
		Organization:        "Acme, Inc.",
		CommonName:          "test.example.com",
		SerialNumber:        "1234abcdef56789",
		Endpoint:            "test.example.com:443",
		RegisteredDirectory: "trisatest.net",
		Expiration:          time.Date(2022, time.July, 18, 16, 38, 38, 0, time.Local),
		Reissuance:          time.Date(2022, time.July, 25, 12, 30, 0, 0, time.Local),
		BaseURL:             "http://localhost:8080/vasps",
	}
	msg, err := emails.ReissuanceAdminNotificationEmail(sender.Name, sender.Address, recipient.Name, recipient.Address, data)

	require.NoError(err)
	require.NoError(email.Send(msg))
	require.Len(mock.Emails, 1)
	expected, err := json.Marshal(msg)
	require.NoError(err)
	require.Equal(expected, mock.Emails[0])

	generateMIME(suite.T(), msg, "reissuance-admin-notification.mim")
}
