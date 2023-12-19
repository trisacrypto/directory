package emails_test

import (
	"net/mail"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/emails"
	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/emails/mock"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func TestClientSend(t *testing.T) {
	// NOTE: if you place a .env file in this directory alongside the test file, it
	// will be read, making it simpler to run tests and set environment variables.
	godotenv.Load()

	// This test uses the environment to send rendered emails with context specific
	// emails - this is to test the rendering of the emails with data only; it does not
	// go through any of the server workflow for generating tokens, etc. If the
	// GDS_TEST_SENDING_CLIENT_EMAILS environment variable is not specified, then a
	// mock email client is used instead which will not send any real emails but will
	// still exercise the logic in the client methods.
	var conf config.EmailConfig
	if os.Getenv("GDS_TEST_SENDING_CLIENT_EMAILS") != "" {
		require.NoError(t, envconfig.Process("gds", &conf), "failed to load environment variables for email client tests")
	} else {
		conf = config.EmailConfig{
			ServiceEmail:         "GDS <service@gds.dev>",
			AdminEmail:           "GDS Admin <admin@gds.dev>",
			SendGridAPIKey:       "notarealsendgridapikey",
			DirectoryID:          "gds.dev",
			VerifyContactBaseURL: "https://gds.dev/verify",
			AdminReviewBaseURL:   "https://admin.gds.dev/vasps/",
			Testing:              true,
		}

		defer mock.PurgeEmails()
	}

	// This test sends emails from the serviceEmail using SendGrid to the adminsEmail
	email, err := emails.New(conf)
	require.NoError(t, err)

	recipient, err := mail.ParseAddress(conf.AdminEmail)
	require.NoError(t, err)

	vasp, contacts := makeClientFixtures(t, recipient)

	// Helper Functions
	resetLogs := func(t *testing.T) {
		for _, contact := range contacts.Emails {
			contact.SendLog = make([]*models.EmailLogEntry, 0)
		}
	}

	verifyContact := func(kind string) {
		contact := contacts.Get(kind)
		contact.Email.Verified = true
		contact.Email.VerifiedOn = "2023-12-19T09:08:42Z"
	}

	resetContacts := func(t *testing.T) {
		iter := contacts.NewIterator()
		for iter.Next() {
			contact := iter.Contact()
			contact.Email.Verified = false
			contact.Email.VerifiedOn = ""
		}
	}

	t.Run("VerifyContacts", func(t *testing.T) {
		defer resetLogs(t)

		sent, err := email.SendVerifyContacts(contacts)
		require.NoError(t, err)
		require.Equal(t, 1, sent)

		// Make sure that the contact pointer was not modified
		require.False(t, contacts.Emails[0].Verified)
		require.NotEmpty(t, contacts.Emails[0].Token)

		// The contacts email log should contain one item
		emailLog, err := models.GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emailLog, 1)
		require.Equal(t, string(admin.ResendVerifyContact), emailLog[0].Reason)
		require.Equal(t, emails.VerifyContactRE, emailLog[0].Subject)
	})

	t.Run("ReviewRequest", func(t *testing.T) {
		defer resetLogs(t)

		sent, err := email.SendReviewRequest(vasp)
		require.NoError(t, err)
		require.Equal(t, 1, sent)

		// Make sure that the VASP pointer was not modified
		token, err := models.GetAdminVerificationToken(vasp)
		require.NoError(t, err)
		require.Equal(t, "12345token1234", token)

		// Ensure that there are no email logs since this should have been sent to an admin
		emailLog, err := models.GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emailLog, 0)
	})

	t.Run("RejectRegistration", func(t *testing.T) {
		defer resetLogs(t)
		defer resetContacts(t)

		// Cannot send an email if there are no verified contacts
		sent, err := email.SendRejectRegistration(vasp, contacts, "this is a test rejection from the test runner")
		require.Error(t, err, "expected no email sent if no verified contacts")
		require.Equal(t, 0, sent)

		// Verify the technical contact
		verifyContact(models.TechnicalContact)
		sent, err = email.SendRejectRegistration(vasp, contacts, "this is a test rejection from the test runner")
		require.NoError(t, err)
		require.Equal(t, 1, sent)

		// Ensure that at least one email was sent
		emailLog, err := models.GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emailLog, 1)
	})

	t.Run("DeliverCertificates", func(t *testing.T) {
		defer resetLogs(t)
		defer resetContacts(t)

		// Cannot send an email if there are no verified contacts
		sent, err := email.SendDeliverCertificates(vasp, contacts, "testdata/foo.zip")
		require.Error(t, err, "expected no emails sent if no verified contacts")
		require.Equal(t, 0, sent)

		// Verify the administrative contact
		verifyContact(models.AdministrativeContact)
		sent, err = email.SendDeliverCertificates(vasp, contacts, "testdata/foo.zip")
		require.NoError(t, err)
		require.Equal(t, 1, sent)

		// Ensure that at least one email was sent
		emailLog, err := models.GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emailLog, 1)
	})

	t.Run("ExpiresAdminNotification", func(t *testing.T) {
		defer resetLogs(t)

		// TODO: For reissuance related emails, test that emails are not sent twice
		reissueDate := time.Date(2022, time.July, 25, 12, 0, 0, 0, time.Local)
		sent, err := email.SendExpiresAdminNotification(vasp, 0, reissueDate)
		require.NoError(t, err)
		require.Equal(t, 1, sent)
		sent, err = email.SendExpiresAdminNotification(vasp, 1, reissueDate)
		require.NoError(t, err)
		require.Equal(t, 0, sent, "should not have sent duplicate expiration email to the admin")

		// Ensure that there are no email logs since this should have been sent to an admin
		emailLog, err := models.GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emailLog, 0)
	})

	t.Run("ReissuanceReminder", func(t *testing.T) {
		defer resetLogs(t)
		defer resetContacts(t)

		// Cannot send reminder if there are no verified contacts
		reissueDate := time.Date(2022, time.July, 25, 12, 0, 0, 0, time.Local)
		sent, err := email.SendReissuanceReminder(vasp, contacts, reissueDate)
		require.Error(t, err, "expected no emails sent if no verified contacts")
		require.Equal(t, 0, sent)

		// Verify the technical contact
		verifyContact(models.TechnicalContact)
		sent, err = email.SendReissuanceReminder(vasp, contacts, reissueDate)
		require.NoError(t, err)
		require.Equal(t, 1, sent)

		// Ensure that at least one email was sent
		emailLog, err := models.GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emailLog, 1)
	})

	t.Run("ReissuanceStarted", func(t *testing.T) {
		defer resetLogs(t)
		defer resetContacts(t)

		// Cannot send reissuance started if there are no verified contacts
		sent, err := email.SendReissuanceStarted(vasp, contacts, "https://whisper.dev/supersecret")
		require.Error(t, err, "expected no emails sent if no verified contacts")
		require.Equal(t, 0, sent)

		// Verify legal contact
		verifyContact(models.LegalContact)
		sent, err = email.SendReissuanceStarted(vasp, contacts, "https://whisper.dev/supersecret")
		require.NoError(t, err)
		require.Equal(t, 1, sent)

		// Ensure that at least one email was sent
		emailLog, err := models.GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emailLog, 1)
	})

	t.Run("ReissuanceAdminNotification", func(t *testing.T) {
		defer resetLogs(t)

		reissuedDate := time.Date(2022, time.July, 25, 12, 0, 0, 0, time.Local)
		sent, err := email.SendReissuanceAdminNotification(vasp, 0, reissuedDate)
		require.NoError(t, err)
		require.Equal(t, 1, sent)
		sent, err = email.SendReissuanceAdminNotification(vasp, 1, reissuedDate)
		require.NoError(t, err)
		require.Equal(t, 0, sent, "should not have sent duplicate reissuance email to the admin")

		// Ensure that there are no email logs since this should have been sent to an admin
		emailLog, err := models.GetVASPEmailLog(contacts)
		require.NoError(t, err, "could not get email log")
		require.Len(t, emailLog, 0)
	})
}

func makeClientFixtures(t *testing.T, recipient *mail.Address) (*pb.VASP, *models.Contacts) {
	vasp := &pb.VASP{
		Id:            uuid.NewString(),
		CommonName:    "test.example.com",
		TrisaEndpoint: "test.example.com:443",
		Entity: &ivms101.LegalPerson{
			Name: &ivms101.LegalPersonName{
				NameIdentifiers: []*ivms101.LegalPersonNameId{
					{
						LegalPersonName:               "Acme, Inc.",
						LegalPersonNameIdentifierType: ivms101.LegalPersonLegal,
					},
				},
			},
		},
		Contacts: &pb.Contacts{
			Technical: &pb.Contact{
				Name:  recipient.Name,
				Email: recipient.Address,
			},
			Administrative: &pb.Contact{
				Name:  recipient.Name,
				Email: recipient.Address,
			},
			Legal: &pb.Contact{
				Name:  recipient.Name,
				Email: recipient.Address,
			},
		},
		IdentityCertificate: &pb.Certificate{
			SerialNumber: []byte("notarealcertificate"),
			NotAfter:     "2022-07-18T17:18:55Z",
		},
	}

	contacts := &models.Contacts{
		VASP:     vasp.Id,
		Contacts: vasp.Contacts,
		Emails: []*models.Email{
			{
				Name:       recipient.Name,
				Email:      recipient.Address,
				Token:      "12345token1234",
				Verified:   false,
				VerifiedOn: "",
				SendLog:    make([]*models.EmailLogEntry, 0),
				Created:    "2023-09-01T08:46:16-05:00",
				Modified:   "2023-09-01T08:46:16-05:00",
			},
		},
	}

	err := models.SetAdminVerificationToken(vasp, "12345token1234")
	require.NoError(t, err)

	return vasp, contacts
}
