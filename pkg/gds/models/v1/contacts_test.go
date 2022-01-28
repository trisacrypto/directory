package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
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
		contact, kind, err := iter.Value()
		require.NoError(t, err)
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
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
		contact, kind, err := iter.Value()
		require.NoError(t, err)
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)

	actualContacts = []*pb.Contact{}
	actualKinds = []string{}

	// Should skip nil contacts.
	contacts.Technical = nil
	contacts.Billing = nil
	iter = models.NewContactIterator(contacts, false, false)
	for iter.Next() {
		contact, kind, err := iter.Value()
		require.NoError(t, err)
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)
}

func TestIterVerfiiedContacts(t *testing.T) {
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
		contact, kind, err := iter.Value()
		require.NoError(t, err)
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
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
		contact, kind, err := iter.Value()
		require.NoError(t, err)
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)
}
