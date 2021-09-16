package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func TestVASPExtra(t *testing.T) {
	// Ensures that the Get/Set methods on the VASP do not overwrite values other than
	// the values the method is intended to interact with.
	vasp := &pb.VASP{}

	// Attempt to get an admin verification token on a nil extra
	token, err := GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "", token)

	// Attempt to set admin verification token on a nil extra
	err = SetAdminVerificationToken(vasp, "pontoonboatz")
	require.NoError(t, err)

	// Should be able to fetch the token
	token, err = GetAdminVerificationToken(vasp)
	require.NoError(t, err)
	require.Equal(t, "pontoonboatz", token)
}

func TestContactExtra(t *testing.T) {
	// Test contact is nil
	token, verified, err := GetContactVerification(nil)
	require.NoError(t, err, "nil contact returns error")
	require.False(t, verified)
	require.Empty(t, token)

	// Test extra is nil
	contact := &pb.Contact{
		Email: "pontoon@boatz.com",
		Name:  "Sailor Moon",
		Phone: "555-5555",
	}
	token, verified, err = GetContactVerification(contact)
	require.NoError(t, err, "nil contact extra returns error")
	require.False(t, verified)
	require.Empty(t, token)

	// Contact cannot be nil to set extra
	err = SetContactVerification(nil, "12345", false)
	require.Error(t, err)

	// Set extra on contact
	err = SetContactVerification(contact, "12345", false)
	require.NoError(t, err)

	// Fetch set extra
	token, verified, err = GetContactVerification(contact)
	require.NoError(t, err)
	require.False(t, verified)
	require.Equal(t, "12345", token)

	// Set extra on contact
	err = SetContactVerification(contact, "", true)
	require.NoError(t, err)

	// Fetch set extra
	token, verified, err = GetContactVerification(contact)
	require.NoError(t, err)
	require.True(t, verified)
	require.Equal(t, "", token)
}

func TestVeriedContacts(t *testing.T) {
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

	contacts := VerifiedContacts(vasp)
	require.Len(t, contacts, 0)

	err := SetContactVerification(vasp.Contacts.Administrative, "", true)
	require.NoError(t, err)

	err = SetContactVerification(vasp.Contacts.Technical, "12345", false)
	require.NoError(t, err)

	contacts = VerifiedContacts(vasp)
	require.Len(t, contacts, 1)

	err = SetContactVerification(vasp.Contacts.Technical, "", true)
	require.NoError(t, err)

	err = SetContactVerification(vasp.Contacts.Legal, "12345", false)
	require.NoError(t, err)

	contacts = VerifiedContacts(vasp)
	require.Len(t, contacts, 2)
}

func TestIsTraveler(t *testing.T) {
	vasp := &pb.VASP{CommonName: "trisa.example.com"}
	require.False(t, IsTraveler(vasp))

	vasp.CommonName = "trisa-1234.traveler.ciphertrace.com"
	require.True(t, IsTraveler(vasp))
}
