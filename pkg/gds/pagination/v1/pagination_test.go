package pagination_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/gds/pagination/v1"
)

// This is a very similar test to trtl.internal but decoupled so Trtl and GDS aren't
// dependent on each other.
func TestPageCursor(t *testing.T) {
	cursor := &PageCursor{
		PageSize: int32(16),
		NextVasp: "c1bbd6ac-1d56-42ad-be9e-04512b0a2066",
	}

	// Testing dumping a cursor
	pageToken, err := cursor.Dump()
	require.NoError(t, err, "could not dump pageToken")
	require.Equal(t, "CBASJGMxYmJkNmFjLTFkNTYtNDJhZC1iZTllLTA0NTEyYjBhMjA2Ng", pageToken)

	// Test loading a cursor
	other := &PageCursor{}
	require.NoError(t, other.Load(pageToken), "could not load pageToken")
	require.Equal(t, cursor.PageSize, other.PageSize)
	require.Equal(t, cursor.NextVasp, other.NextVasp)
}
