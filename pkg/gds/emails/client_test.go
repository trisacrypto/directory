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
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func TestClientSendEmails(t *testing.T) {
	// NOTE: if you place a .env file in this directory alongside the test file, it
	// will be read, making it simpler to run tests and set environment variables.
	godotenv.Load()

	if os.Getenv("GDS_TEST_SENDING_CLIENT_EMAILS") == "" {
		t.Skip("skip client send emails test")
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

	receipient, err := mail.ParseAddress(conf.AdminEmail)
	require.NoError(t, err)

	vasp := &pb.VASP{
		Id:            uuid.NewString(),
		CommonName:    "test.example.com",
		TrisaEndpoint: "test.example.com:443",
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
	err = models.SetContactVerification(vasp.Contacts.Technical, "", true)
	require.NoError(t, err)
	err = models.SetContactVerification(vasp.Contacts.Administrative, "", true)
	require.NoError(t, err)
	err = models.SetContactVerification(vasp.Contacts.Legal, "12345token1234", false)
	require.NoError(t, err)

	sent, err := email.SendVerifyContacts(vasp)
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
	require.True(t, verified)
	require.Equal(t, "", token)

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
	require.Equal(t, 2, sent)

	sent, err = email.SendDeliverCertificates(vasp, "testdata/foo.zip")
	require.NoError(t, err)
	require.Equal(t, 1, sent)

	reissueDate := time.Date(2022, time.July, 25, 12, 0, 0, 0, time.Local)
	sent, err = email.SendExpiresAdminNotification(vasp, 0, reissueDate)
	require.NoError(t, err)
	require.Equal(t, 1, sent)
	sent, err = email.SendExpiresAdminNotification(vasp, 1, reissueDate)
	require.NoError(t, err)
	require.Equal(t, 0, sent, "should not have sent duplicate expiration email to the admin")

	sent, err = email.SendReissuanceReminder(vasp, reissueDate)
	require.NoError(t, err)
	require.Equal(t, 2, sent)

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

	// Technical is verified and first so should get Rejection and DeliverCerts emails
	// It should also receive the reissuance started email after the reminder.
	emailLog, err := models.GetEmailLog(vasp.Contacts.Technical)
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

	// Administrative is verified so should get Rejection and Reissue Reminder emails
	emailLog, err = models.GetEmailLog(vasp.Contacts.Administrative)
	require.NoError(t, err)
	require.Len(t, emailLog, 2)
	require.Equal(t, string(admin.ResendRejection), emailLog[0].Reason)
	require.Equal(t, emails.RejectRegistrationRE, emailLog[0].Subject)
	require.Equal(t, string(admin.ReissuanceReminder), emailLog[1].Reason)
	require.Equal(t, emails.ReissuanceReminderRE, emailLog[1].Subject)

	// Legal is not verified so should get VerifyContact email
	emailLog, err = models.GetEmailLog(vasp.Contacts.Legal)
	require.NoError(t, err)
	require.Len(t, emailLog, 1)
	require.Equal(t, string(admin.ResendVerifyContact), emailLog[0].Reason)
	require.Equal(t, emails.VerifyContactRE, emailLog[0].Subject)

	// Billing doesn't have an associated email so shouldn't get anything
	emailLog, err = models.GetEmailLog(vasp.Contacts.Billing)
	require.NoError(t, err)
	require.Len(t, emailLog, 0)
}
