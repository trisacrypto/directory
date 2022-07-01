package models_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

func TestOrganizationKey(t *testing.T) {
	org := &models.Organization{}
	require.Equal(t, uuid.Nil[:], org.Key(), "expected nil uuid when organization id is empty string")

	uu := uuid.New()
	org.Id = uu.String()
	require.Equal(t, uu[:], org.Key(), "expected key to be uuid bytes")

	require.Panics(t, func() {
		org.Id = "notauuid"
		org.Key()
	}, "if the organization id is not a uuid string, expect a panic")
}

func TestParseOrgID(t *testing.T) {
	example := uuid.New()

	testCases := []struct {
		expected uuid.UUID
		input    interface{}
		err      error
	}{
		{example, example.String(), nil},       // parse string
		{example, example[:], nil},             // parse bytes
		{example, example, nil},                // parse uuid
		{uuid.Nil, 14, models.ErrInvalidOrgID}, // unknown type
	}

	for i, tc := range testCases {
		uu, err := models.ParseOrgID(tc.input)
		if tc.err != nil {
			require.Equal(t, uuid.Nil, uu, "expected nil uuid in test case %d", i)
			require.ErrorIs(t, err, tc.err, "expected error in test case %d", i)
		} else {
			require.NoError(t, err, "expected no error to occur in test case %d", i)
			require.Equal(t, tc.expected, uu, "unexpected org id returned when parsed in test case %d", i)
		}
	}
}
