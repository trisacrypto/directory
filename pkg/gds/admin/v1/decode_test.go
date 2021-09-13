package admin_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/gds/admin/v1"
)

func TestResendActionSerialization(t *testing.T) {
	// Test valid enums
	cases := []ResendAction{ResendVerifyContact, ResendReview, ResendDeliverCerts, ResendRejection}
	for _, tc := range cases {
		data, err := json.Marshal(tc)
		require.NoError(t, err)

		var action ResendAction
		err = json.Unmarshal(data, &action)
		require.NoError(t, err)

		require.Equal(t, tc, action)
	}

	// Test white-space- and case-insensitivity
	for _, tc := range []string{"REVIEW", "Review", "   review", "  ReView ", "review   "} {
		data, err := json.Marshal(tc)
		require.NoError(t, err)

		var action ResendAction
		err = json.Unmarshal(data, &action)
		require.NoError(t, err)

		require.Equal(t, ResendReview, action)
	}

	// Test invalid enums
	for _, tc := range []string{"", "foo"} {
		data, err := json.Marshal(tc)
		require.NoError(t, err)

		var action ResendAction
		err = json.Unmarshal(data, &action)
		require.ErrorIs(t, err, ErrInvalidResendAction)
	}
}
