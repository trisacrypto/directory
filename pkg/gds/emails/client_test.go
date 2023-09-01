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

func TestClientSendEmails(t *testing.T) {
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

	receipient, err := mail.ParseAddress(conf.AdminEmail)
	require.NoError(t, err)

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
				Name:  receipient.Name,
				Email: receipient.Address,
			},
			Administrative: &pb.Contact{
				Name:  receipient.Name,
				Email: receipient.Address,
			},
			Legal: &pb.Contact{
				Name:  receipient.Name,
				Email: receipient.Address,
			},
		},
		IdentityCertificate: &pb.Certificate{
			SerialNumber: []byte("notarealcertificate"),
			NotAfter:     "2022-07-18T17:18:55Z",
		},
	}

	err = models.SetAdminVerificationToken(vasp, "12345token1234")
	require.NoError(t, err)
	err = models.SetContactVerification(vasp.Contacts.Technical, "12345token1234", false)
	require.NoError(t, err)
	err = models.SetContactVerification(vasp.Contacts.Administrative, "", true)
	require.NoError(t, err)
	err = models.SetContactVerification(vasp.Contacts.Legal, "12345token1234", false)
	require.NoError(t, err)

	sent, err := email.SendVerifyContacts(vasp, nil)
	require.NoError(t, err)
	require.Equal(t, 1, sent)

	sent, err = email.SendReviewRequest(vasp)
	require.NoError(t, err)
	require.Equal(t, 1, sent)

	// Make sure that the VASP pointer was not modified
	token, err := models.GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "12345token1234", token)

	token, verified, err := models.GetContactVerification(vasp.Contacts.Technical)
	require.NoError(t, err)
	require.False(t, verified)
	require.Equal(t, "12345token1234", token)

	token, verified, err = models.GetContactVerification(vasp.Contacts.Administrative)
	require.NoError(t, err)
	require.True(t, verified)
	require.Equal(t, "", token)

	token, verified, err = models.GetContactVerification(vasp.Contacts.Legal)
	require.NoError(t, err)
	require.False(t, verified)
	require.Equal(t, "12345token1234", token)

	token, verified, err = models.GetContactVerification(vasp.Contacts.Billing)
	require.NoError(t, err)
	require.False(t, verified)
	require.Equal(t, "", token)

	sent, err = email.SendRejectRegistration(vasp, "this is a test rejection from the test runner")
	require.NoError(t, err)
	require.Equal(t, 1, sent)

	sent, err = email.SendDeliverCertificates(vasp, "testdata/foo.zip")
	require.NoError(t, err)
	require.Equal(t, 1, sent)

	// TODO: For reissuance related emails, test that emails are not sent twice
	reissueDate := time.Date(2022, time.July, 25, 12, 0, 0, 0, time.Local)
	sent, err = email.SendExpiresAdminNotification(vasp, 0, reissueDate)
	require.NoError(t, err)
	require.Equal(t, 1, sent)
	sent, err = email.SendExpiresAdminNotification(vasp, 1, reissueDate)
	require.NoError(t, err)
	require.Equal(t, 0, sent, "should not have sent duplicate expiration email to the admin")

	sent, err = email.SendReissuanceReminder(vasp, reissueDate)
	require.NoError(t, err)
	require.Equal(t, 1, sent)

	sent, err = email.SendReissuanceStarted(vasp, "https://whisper.dev/supersecret")
	require.NoError(t, err)
	require.Equal(t, 1, sent)

	reissuedDate := time.Date(2022, time.July, 25, 12, 0, 0, 0, time.Local)
	sent, err = email.SendReissuanceAdminNotification(vasp, 0, reissuedDate)
	require.NoError(t, err)
	require.Equal(t, 1, sent)
	sent, err = email.SendReissuanceAdminNotification(vasp, 1, reissuedDate)
	require.NoError(t, err)
	require.Equal(t, 0, sent, "should not have sent duplicate reissuance email to the admin")

	// TRISA Admin should get an expiration notification email and a reissuance started email
	log, err := models.GetAdminEmailLog(vasp)
	require.NoError(t, err)
	require.Len(t, log, 2)
	require.Equal(t, string(admin.ReissuanceReminder), log[0].Reason)
	require.Equal(t, emails.ExpiresAdminNotificationRE, log[0].Subject)
	require.Equal(t, string(admin.ReissuanceStarted), log[1].Reason)
	require.Equal(t, emails.ReissuanceAdminNotificationRE, log[1].Subject)

	// Technical is not verified so it should get the VerifyContact email
	emailLog, err := models.GetEmailLog(vasp.Contacts.Technical)
	require.NoError(t, err)
	require.Len(t, emailLog, 1)
	require.Equal(t, string(admin.ResendVerifyContact), emailLog[0].Reason)
	require.Equal(t, emails.VerifyContactRE, emailLog[0].Subject)

	// Administrative is the first verified contact so it should get Rejection, Deliver
	// Certs, Reissue Reminder, and Reissuance Started after the reminder
	emailLog, err = models.GetEmailLog(vasp.Contacts.Administrative)
	require.NoError(t, err)
	require.Len(t, emailLog, 4)
	require.Equal(t, string(admin.ResendRejection), emailLog[0].Reason)
	require.Equal(t, emails.RejectRegistrationRE, emailLog[0].Subject)
	require.Equal(t, string(admin.ResendDeliverCerts), emailLog[1].Reason)
	require.Equal(t, emails.DeliverCertsRE, emailLog[1].Subject)
	require.Equal(t, string(admin.ReissuanceReminder), emailLog[2].Reason)
	require.Equal(t, emails.ReissuanceReminderRE, emailLog[2].Subject)
	require.Equal(t, string(admin.ReissuanceStarted), emailLog[3].Reason)
	require.Equal(t, emails.ReissuanceStartedRE, emailLog[3].Subject)

	// Legal is not verified and it has the same email as Administrative so it should
	// not get any emails
	emailLog, err = models.GetEmailLog(vasp.Contacts.Legal)
	require.NoError(t, err)
	require.Len(t, emailLog, 0)

	// Billing doesn't have an associated email so shouldn't get anything
	emailLog, err = models.GetEmailLog(vasp.Contacts.Billing)
	require.NoError(t, err)
	require.Len(t, emailLog, 0)
}
