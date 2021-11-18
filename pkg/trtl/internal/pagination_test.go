package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/trtl/internal"
)

func TestPageCursor(t *testing.T) {
	cursor := &internal.PageCursor{
		PageSize: int32(42),
		NextKey:  []byte("default::e23c012c-c5a7-4cae-b944-7a65c8ece4e0"),
	}

	// Testing dumping a cursor
	pageToken, err := cursor.Dump()
	require.NoError(t, err, "could not dump pageToken")
	require.Equal(t, "CCoSLWRlZmF1bHQ6OmUyM2MwMTJjLWM1YTctNGNhZS1iOTQ0LTdhNjVjOGVjZTRlMA", pageToken)

	// Test loading a cursor
	other := &internal.PageCursor{}
	require.NoError(t, other.Load(pageToken), "could not load pageToken")
	require.Equal(t, cursor.PageSize, other.PageSize)
	require.Equal(t, cursor.NextKey, other.NextKey)
}
