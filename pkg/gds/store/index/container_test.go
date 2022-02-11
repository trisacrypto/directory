package index_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/store/index"
)

func TestContainerIndex(t *testing.T) {
	countries := make(index.Container)
	require.Empty(t, countries)

	// Add a single entry
	id1 := uuid.New().String()
	require.True(t, countries.Add("United States", id1, index.NormalizeCountry))
	require.Len(t, countries, 1)
	require.Contains(t, countries, "US")
	require.Contains(t, countries["US"], id1)

	// Add a second entry
	id2 := uuid.New().String()
	require.True(t, countries.Add("USA", id2, index.NormalizeCountry))
	require.Len(t, countries, 1)
	require.Contains(t, countries, "US")
	require.Contains(t, countries["US"], id2)
	require.Len(t, countries["US"], 2)

	// Add a third entry
	id3 := uuid.New().String()
	require.True(t, countries.Add("united states", id3, index.NormalizeCountry))
	require.Len(t, countries, 1)
	require.Contains(t, countries, "US")
	require.Contains(t, countries["US"], id3)
	require.Len(t, countries["US"], 3)

	// Add a duplicate entry
	require.False(t, countries.Add("US", id2, index.NormalizeCountry))
	require.Len(t, countries["US"], 3)

	// Add another country
	id4 := uuid.New().String()
	require.True(t, countries.Add("spain", id4, index.NormalizeCountry))
	require.Len(t, countries, 2)
	require.Contains(t, countries, "ES")
	require.Contains(t, countries["ES"], id4)
	require.Len(t, countries["ES"], 1)

	// Add a multiple country entry
	require.True(t, countries.Add("ES", id2, index.NormalizeCountry))
	require.Len(t, countries["US"], 3)
	require.Contains(t, countries["ES"], id2)
	require.Len(t, countries["ES"], 2)

	// Add another country
	id5 := uuid.New().String()
	require.True(t, countries.Add("New Zealand", id5, index.NormalizeCountry))
	require.Len(t, countries, 3)
	require.Contains(t, countries, "NZ")
	require.Contains(t, countries["NZ"], id5)
	require.Len(t, countries["NZ"], 1)

	// Remove an entry
	require.True(t, countries.Remove("USA", id3, index.NormalizeCountry))
	require.Len(t, countries, 3)
	require.Contains(t, countries, "US")
	require.NotContains(t, countries["US"], id3)
	require.Len(t, countries["US"], 2)

	// Remove an entry that doesn't exist
	require.False(t, countries.Remove("USA", id3, index.NormalizeCountry))
	require.Len(t, countries, 3)
	require.Contains(t, countries, "US")
	require.NotContains(t, countries["US"], id3)
	require.Len(t, countries["US"], 2)

	// Empty a country
	require.True(t, countries.Remove("New Zealand", id5, index.NormalizeCountry))
	require.Len(t, countries, 3)
	require.Contains(t, countries, "NZ")
	require.Empty(t, countries["NZ"])

	// Test Find and reverse
	vals, ok := countries.Find("Spain", index.NormalizeCountry)
	require.True(t, ok)
	require.Len(t, vals, 2)

	vals, ok = countries.Find("Russia", index.NormalizeCountry)
	require.False(t, ok)
	require.Empty(t, vals)

	vals, ok = countries.Reverse(id2, nil)
	require.True(t, ok)
	require.Len(t, vals, 2)

	vals, ok = countries.Reverse("foo", nil)
	require.False(t, ok)
	require.Empty(t, vals)

	// Test serialization and deserialization
	data, err := countries.Dump()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	other := make(index.Container)
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

func TestContainerNonNormalize(t *testing.T) {
	idx := make(index.Container)
	require.True(t, idx.Add("foo", "bar", nil))
	require.True(t, idx.Add("Foo", "bar", nil))
	require.Len(t, idx, 2)
	require.Contains(t, idx, "foo")
	require.Contains(t, idx, "Foo")
	require.True(t, idx.Remove("Foo", "bar", nil))
	require.NotContains(t, idx["Foo"], "bar")
	require.Contains(t, idx["foo"], "bar")
}
