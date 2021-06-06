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
		},
		IdentityCertificate: &pb.Certificate{
			SerialNumber: []byte("notarealcertificate"),
		},
	}

	err = models.SetContactVerification(vasp.Contacts.Technical, "", true)
	require.NoError(t, err)

	require.NoError(t, email.SendVerifyContacts(vasp))
	require.NoError(t, email.SendReviewRequest(vasp))
	require.NoError(t, email.SendRejectRegistration(vasp, "this is a test rejection from the test runner"))
	require.NoError(t, email.SendDeliverCertificates(vasp, "testdata/foo.zip"))
}
