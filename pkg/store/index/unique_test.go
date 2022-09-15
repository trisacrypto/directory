package index_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/store/index"
)

func TestUniqueIndex(t *testing.T) {
	names := make(index.Unique)
	require.Empty(t, names, 0)

	// Add a new single entry
	id1 := uuid.New().String()
	require.True(t, names.Add("    Foo   ", id1, index.Normalize))
	require.Len(t, names, 1)
	require.Contains(t, names, "foo")
	require.Equal(t, names["foo"], id1)

	// Add a duplicate entry, no overwrite
	id2 := uuid.New().String()
	require.False(t, names.Add("FOO", id2, index.Normalize))
	require.Len(t, names, 1)
	require.NotContains(t, names, "FOO")
	require.Equal(t, names["foo"], id1)

	// Overwrite an entry
	require.True(t, names.Overwrite("Foo", id2, index.Normalize))
	require.Len(t, names, 1)
	require.Contains(t, names, "foo")
	require.Equal(t, names["foo"], id2)

	// Add a few more entries
	id3 := uuid.New().String()
	require.True(t, names.Add("Bar", id1, index.Normalize))
	require.True(t, names.Add("bAZ  ", id3, index.Normalize))
	require.Len(t, names, 3)

	// Remove an entry
	require.True(t, names.Remove("Baz", index.Normalize))
	require.Len(t, names, 2)
	require.NotContains(t, names, "baz")
	require.NotContains(t, names, "Baz")

	// Remove entry twice
	require.False(t, names.Remove(" baZ", index.Normalize))
	require.Len(t, names, 2)

	// Find some entries
	val, ok := names.Find("BAR", index.Normalize)
	require.True(t, ok)
	require.Equal(t, val, id1)

	_, ok = names.Find("BaZ", index.Normalize)
	require.False(t, ok)

	// Create some duplicate values
	require.True(t, names.Add("zing", id1, index.Normalize))
	vals, ok := names.Reverse(id1, index.Normalize)
	require.True(t, ok)
	require.Len(t, vals, 2)

	vals, ok = names.Reverse("foo", index.Normalize)
	require.False(t, ok)
	require.Empty(t, vals)

	// Test serialization and deserialization
	data, err := names.Dump()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	other := make(index.Unique)
	require.Empty(t, other)
	err = other.Load(data)
	require.NoError(t, err)
	require.NotEmpty(t, other)

	for key, value := range names {
		require.Contains(t, other, key)
		require.Equal(t, other[key], value)
	}
}

func TestUniqueNonNormalize(t *testing.T) {
	idx := make(index.Unique)
	require.True(t, idx.Add("foo", "bar", nil))
	require.True(t, idx.Add("Foo", "bar", nil))
	require.Len(t, idx, 2)
	require.Contains(t, idx, "foo")
	require.Contains(t, idx, "Foo")
	require.True(t, idx.Remove("Foo", nil))
	require.Len(t, idx, 1)
	require.Contains(t, idx, "foo")
	require.NotContains(t, idx, "Foo")
}
