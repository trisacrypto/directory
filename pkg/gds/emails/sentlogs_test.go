package emails_test

import (
	"testing"

	"github.com/trisacrypto/directory/pkg/gds/emails"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func TestVerifyEmailSentVASPContacts(t *testing.T) {
	vasp := &pb.VASP{
		Contacts: &pb.Contacts{
			Administrative: &pb.Contact{
				Name:  "Admin Person",
				Email: "admin@example.com",
			},
			Technical: &pb.Contact{
				Name:  "Technical Person",
				Email: "tech@example.com",
			},
			Legal: &pb.Contact{
				Name:  "Legal Person",
				Email: "legal@example.com",
			},
		},
	}

	//  log should initially be empty
	emailLog, err := models.GetEmailLog(vasp.Contacts.Administrative)
	require.NoError(t, err)
	require.Len(t, emailLog, 0)

	// Append an entry to an empty log
	err = models.AppendEmailLog(vasp.Contacts.Administrative, "verify_contact", "verification")
	require.NoError(t, err)

	// get email log for contact
	emailLog, err = models.GetEmailLog(vasp.Contacts.Administrative)
	require.NoError(t, err)
	require.Len(t, emailLog, 1)
	require.Equal(t, "verify_contact", emailLog[0].Reason)
	require.Equal(t, "verification", emailLog[0].Subject)

	status, err := emails.VerifyEmailSentVASPContact(vasp, "administrative", 30)
	require.NoError(t, err)
	require.Equal(t, true, status)
}
