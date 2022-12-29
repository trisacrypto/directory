package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestIterContacts(t *testing.T) {
	contacts := &pb.Contacts{
		Technical: &pb.Contact{
			Name: "technical",
		},
		Administrative: &pb.Contact{
			Email: "administrative@example.com",
		},
		Billing: &pb.Contact{
			Name: "billing",
		},
		Legal: &pb.Contact{
			Email: "legal@example.com",
		},
	}
	expectedContacts := []*pb.Contact{
		contacts.Technical,
		contacts.Administrative,
		contacts.Legal,
		contacts.Billing,
	}
	expectedKinds := []string{
		models.TechnicalContact,
		models.AdministrativeContact,
		models.LegalContact,
		models.BillingContact,
	}

	actualContacts := []*pb.Contact{}
	actualKinds := []string{}

	// Should iterate over all contacts.
	iter := models.NewContactIterator(contacts, false, false)
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)

	actualContacts = []*pb.Contact{}
	actualKinds = []string{}

	// Should skip contacts without an email address.
	expectedContacts = []*pb.Contact{
		contacts.Administrative,
		contacts.Legal,
	}
	expectedKinds = []string{
		models.AdministrativeContact,
		models.LegalContact,
	}
	iter = models.NewContactIterator(contacts, true, false)
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)

	actualContacts = []*pb.Contact{}
	actualKinds = []string{}

	// Should skip nil contacts.
	contacts.Technical = nil
	contacts.Billing = nil
	iter = models.NewContactIterator(contacts, false, false)
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)
}

func TestIterVerifiedContacts(t *testing.T) {
	contacts := &pb.Contacts{
		Technical: &pb.Contact{
			Email: "technical@example.com",
		},
		Administrative: &pb.Contact{
			Email: "administrative@example.com",
		},
		Billing: &pb.Contact{
			Email: "billing@example.com",
		},
		Legal: &pb.Contact{
			Email: "legal@example.com",
		},
	}

	actualContacts := []*pb.Contact{}
	actualKinds := []string{}

	// No contacts are verified.
	iter := models.NewContactIterator(contacts, false, true)
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, []*pb.Contact{}, actualContacts)
	require.Equal(t, []string{}, actualKinds)

	actualContacts = []*pb.Contact{}
	actualKinds = []string{}

	// Should only iterate through the verified contacts.
	require.NoError(t, models.SetContactVerification(contacts.Technical, "", true))
	require.NoError(t, models.SetContactVerification(contacts.Legal, "", true))
	expectedContacts := []*pb.Contact{
		contacts.Technical,
		contacts.Legal,
	}
	expectedKinds := []string{
		models.TechnicalContact,
		models.LegalContact,
	}
	iter = models.NewContactIterator(contacts, false, true)
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)
}

func TestGetSentEmailCount(t *testing.T) {
	contacts := &pb.Contacts{
		Technical: &pb.Contact{
			Email: "technical@example.com",
		},
		Administrative: &pb.Contact{
			Email: "administrative@example.com",
		},
		Billing: &pb.Contact{
			Email: "billing@example.com",
		},
		Legal: &pb.Contact{
			Email: "legal@example.com",
		},
	}

	// Log should initially be empty
	emailLog, err := models.GetEmailLog(contacts.Administrative)
	require.NoError(t, err)
	require.Len(t, emailLog, 0)

	// Error should be returned if the reason is empty
	_, err = models.CountSentEmails(emailLog, "", 30)
	require.EqualError(t, err, "cannot match on empty reason string")

	// Error should be returned if the time window is invalid
	_, err = models.CountSentEmails(emailLog, "test", -1)
	require.EqualError(t, err, "time window must be a positive number of days")

	// Append an entry to an empty log
	err = models.AppendEmailLog(contacts.Administrative, "verify_contact", "verification")
	require.NoError(t, err)

	// Append an entry to an empty log
	err = models.AppendEmailLog(contacts.Administrative, "verify_contact", "verification")
	require.NoError(t, err)

	// Get email log for contact
	emailLog, err = models.GetEmailLog(contacts.Administrative)
	require.NoError(t, err)
	require.Len(t, emailLog, 2)
	require.Equal(t, "verify_contact", emailLog[0].Reason)
	require.Equal(t, "verification", emailLog[0].Subject)

	// Should return 2 emails sent for contact
	sent, err := models.CountSentEmails(emailLog, "verify_contact", 30)
	require.NoError(t, err)
	require.Equal(t, 2, sent)

	// Get the technical contact's email log
	emailLog, err = models.GetEmailLog(contacts.Technical)
	require.NoError(t, err)

	// Should return 0 emails when the log is empty
	sent, err = models.CountSentEmails(emailLog, "verify_contact", 30)
	require.NoError(t, err)
	require.Equal(t, 0, sent)

	// Construct an email log with entries at different times
	log := []*models.EmailLogEntry{
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 31).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 29).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 28).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
	}
	require.NoError(t, SetEmailLog(contacts.Billing, log))

	// Get the billing contact's email log
	emailLog, err = models.GetEmailLog(contacts.Billing)
	require.NoError(t, err)

	// Should only return a count of emails within the time window
	sent, err = models.CountSentEmails(emailLog, "verify_contact", 32)
	require.NoError(t, err)
	require.Equal(t, 3, sent, "expected 3 emails sent within the last 32 days")

	sent, err = models.CountSentEmails(emailLog, "verify_contact", 30)
	require.NoError(t, err)
	require.Equal(t, 2, sent, "expected 2 emails sent within the last 30 days")

	sent, err = models.CountSentEmails(emailLog, "verify_contact", 27)
	require.NoError(t, err)
	require.Equal(t, 0, sent, "expected 0 emails sent within the last 27 days")

	// Construct an email log with entries of different reasons
	log = []*models.EmailLogEntry{
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 31).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 29).Format(time.RFC3339),
			Reason:    "rejection",
			Subject:   "rejected registration",
		},
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 28).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
	}
	require.NoError(t, SetEmailLog(contacts.Legal, log))

	// Get the legal contact's email log
	emailLog, err = models.GetEmailLog(contacts.Legal)
	require.NoError(t, err)

	// Should only return a count of emails that match the reason and are within the time window
	sent, err = models.CountSentEmails(emailLog, "verify_contact", 32)
	require.NoError(t, err)
	require.Equal(t, 2, sent, "expected 2 emails sent within the last 32 days")

	sent, err = models.CountSentEmails(emailLog, "rejection", 30)
	require.NoError(t, err)
	require.Equal(t, 1, sent, "expected 1 emails sent within the last 30 days")
}

// Helper function to serialize an email log onto a contact's extra data.
func SetEmailLog(contact *pb.Contact, log []*models.EmailLogEntry) (err error) {
	extra := &models.GDSContactExtraData{}
	extra.EmailLog = log
	if contact.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}
