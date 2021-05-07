package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUniqueIndex(t *testing.T) {
	names := make(uniqueIndex)
	require.Empty(t, names, 0)

	// Add a new single entry
	id1 := uuid.New().String()
	require.True(t, names.add("    Foo   ", id1, normalize))
	require.Len(t, names, 1)
	require.Contains(t, names, "foo")
	require.Equal(t, names["foo"], id1)

	// Add a duplicate entry, no overwrite
	id2 := uuid.New().String()
	require.False(t, names.add("FOO", id2, normalize))
	require.Len(t, names, 1)
	require.NotContains(t, names, "FOO")
	require.Equal(t, names["foo"], id1)

	// Overwrite an entry
	require.True(t, names.overwrite("Foo", id2, normalize))
	require.Len(t, names, 1)
	require.Contains(t, names, "foo")
	require.Equal(t, names["foo"], id2)

	// Add a few more entries
	id3 := uuid.New().String()
	require.True(t, names.add("Bar", id1, normalize))
	require.True(t, names.add("bAZ  ", id3, normalize))
	require.Len(t, names, 3)

	// Remove an entry
	require.True(t, names.rm("Baz", normalize))
	require.Len(t, names, 2)
	require.NotContains(t, names, "baz")
	require.NotContains(t, names, "Baz")

	// Remove entry twice
	require.False(t, names.rm(" baZ", normalize))
	require.Len(t, names, 2)

	// Find some entries
	val, ok := names.find("BAR", normalize)
	require.True(t, ok)
	require.Equal(t, val, id1)

	val, ok = names.find("BaZ", normalize)
	require.False(t, ok)

	// Create some duplicate values
	require.True(t, names.add("zing", id1, normalize))
	vals, ok := names.reverse(id1, normalize)
	require.True(t, ok)
	require.Len(t, vals, 2)

	vals, ok = names.reverse("foo", normalize)
	require.False(t, ok)
	require.Empty(t, vals)

	// Test serialization and deserialization
	data, err := names.Dump()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	other := make(uniqueIndex)
	require.Empty(t, other)
	err = other.Load(data)
	require.NoError(t, err)
	require.NotEmpty(t, other)

	for key, value := range names {
		require.Contains(t, other, key)
		require.Equal(t, other[key], value)
	}
}

func TestContainerIndex(t *testing.T) {
	countries := make(containerIndex)
	require.Empty(t, countries)

	// Add a single entry
	id1 := uuid.New().String()
	require.True(t, countries.add("United States", id1, normalizeCountry))
	require.Len(t, countries, 1)
	require.Contains(t, countries, "US")
	require.Contains(t, countries["US"], id1)

	// Add a second entry
	id2 := uuid.New().String()
	require.True(t, countries.add("USA", id2, normalizeCountry))
	require.Len(t, countries, 1)
	require.Contains(t, countries, "US")
	require.Contains(t, countries["US"], id2)
	require.Len(t, countries["US"], 2)

	// Add a third entry
	id3 := uuid.New().String()
	require.True(t, countries.add("united states", id3, normalizeCountry))
	require.Len(t, countries, 1)
	require.Contains(t, countries, "US")
	require.Contains(t, countries["US"], id3)
	require.Len(t, countries["US"], 3)

	// Add a duplicate entry
	require.False(t, countries.add("US", id2, normalizeCountry))
	require.Len(t, countries["US"], 3)

	// Add another country
	id4 := uuid.New().String()
	require.True(t, countries.add("spain", id4, normalizeCountry))
	require.Len(t, countries, 2)
	require.Contains(t, countries, "ES")
	require.Contains(t, countries["ES"], id4)
	require.Len(t, countries["ES"], 1)

	// Add a multiple country entry
	require.True(t, countries.add("ES", id2, normalizeCountry))
	require.Len(t, countries["US"], 3)
	require.Contains(t, countries["ES"], id2)
	require.Len(t, countries["ES"], 2)

	// Add another country
	id5 := uuid.New().String()
	require.True(t, countries.add("New Zealand", id5, normalizeCountry))
	require.Len(t, countries, 3)
	require.Contains(t, countries, "NZ")
	require.Contains(t, countries["NZ"], id5)
	require.Len(t, countries["NZ"], 1)

	// Remove an entry
	require.True(t, countries.rm("USA", id3, normalizeCountry))
	require.Len(t, countries, 3)
	require.Contains(t, countries, "US")
	require.NotContains(t, countries["US"], id3)
	require.Len(t, countries["US"], 2)

	// Remove an entry that doesn't exist
	require.False(t, countries.rm("USA", id3, normalizeCountry))
	require.Len(t, countries, 3)
	require.Contains(t, countries, "US")
	require.NotContains(t, countries["US"], id3)
	require.Len(t, countries["US"], 2)

	// Empty a country
	require.True(t, countries.rm("New Zealand", id5, normalizeCountry))
	require.Len(t, countries, 3)
	require.Contains(t, countries, "NZ")
	require.Empty(t, countries["NZ"])

	// Test Find and reverse
	vals, ok := countries.find("Spain", normalizeCountry)
	require.True(t, ok)
	require.Len(t, vals, 2)

	vals, ok = countries.find("Russia", normalizeCountry)
	require.False(t, ok)
	require.Empty(t, vals)

	vals, ok = countries.reverse(id2, nil)
	require.True(t, ok)
	require.Len(t, vals, 2)

	vals, ok = countries.reverse("foo", nil)
	require.False(t, ok)
	require.Empty(t, vals)

	// Test serialization and deserialization
	data, err := countries.Dump()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	other := make(containerIndex)
	require.Empty(t, other)
	err = other.Load(data)
	require.NoError(t, err)
	require.NotEmpty(t, other)

	for key, values := range countries {
		require.Contains(t, other, key)
		for i, val := range values {
			require.Equal(t, other[key][i], val)
		}
	}
}

func TestSequence(t *testing.T) {
	var pk sequence
	require.Equal(t, uint64(pk), uint64(0))

	// Perform some operations on the pk
	for i := 0; i < 3; i++ {
		pk++
	}

	require.Equal(t, uint64(pk), uint64(3))

	var id uint64
	if pk > 0 {
		id = uint64(pk)
	}

	for i := 0; i < 5; i++ {
		pk++
	}

	require.Equal(t, uint64(3), id)
	require.Equal(t, uint64(8), uint64(pk))
	require.Equal(t, sequence(8), pk)

	data, err := pk.Dump()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	for i := 0; i < 4; i++ {
		pk++
	}

	other, err := pk.Load(data)
	require.NoError(t, err)
	require.Equal(t, sequence(12), pk)
	require.NotEqual(t, other, pk)
	require.Equal(t, sequence(8), other)
}

func TestContainerNonNormalize(t *testing.T) {
	idx := make(containerIndex)
	require.True(t, idx.add("foo", "bar", nil))
	require.True(t, idx.add("Foo", "bar", nil))
	require.Len(t, idx, 2)
	require.Contains(t, idx, "foo")
	require.Contains(t, idx, "Foo")
	require.True(t, idx.rm("Foo", "bar", nil))
	require.NotContains(t, idx["Foo"], "bar")
	require.Contains(t, idx["foo"], "bar")
}

func TestUniqueNonNormalize(t *testing.T) {
	idx := make(uniqueIndex)
	require.True(t, idx.add("foo", "bar", nil))
	require.True(t, idx.add("Foo", "bar", nil))
	require.Len(t, idx, 2)
	require.Contains(t, idx, "foo")
	require.Contains(t, idx, "Foo")
	require.True(t, idx.rm("Foo", nil))
	require.Len(t, idx, 1)
	require.Contains(t, idx, "foo")
	require.NotContains(t, idx, "Foo")
}
