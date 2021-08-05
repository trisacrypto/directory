package emails_test

import (
	"net/mail"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"
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
}
