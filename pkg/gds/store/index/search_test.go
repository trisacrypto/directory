package index_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/store/index"
)

func TestNormalize(t *testing.T) {
	tt := []struct {
		in       string
		expected string
	}{
		{"foo", "foo"},
		{"FOO", "foo"},
		{"Foo", "foo"},
		{"A Red Line", "a red line"},
		{"   TeST", "test"},
		{"   TeST   ", "test"},
		{"TeST   ", "test"},
		{"United States", "united states"},
	}

	for i, tc := range tt {
		require.Equal(t, tc.expected, index.Normalize(tc.in), "test case %d failed", i)
	}
}

func TestNormalizeCountry(t *testing.T) {
	tt := []struct {
		in       string
		expected string
	}{
		{"foo", "foo"},
		{"United States", "US"},
		{"united states", "US"},
		{"   JAPAN  ", "JP"},
		{" Falkland Islands", "FK"},
		{"Falkland", "FK"},
		{"United", "united"},
	}

	for i, tc := range tt {
		require.Equal(t, tc.expected, index.NormalizeCountry(tc.in), "test case %d failed", i)
	}
}

func TestNormalizeURL(t *testing.T) {
	tt := []struct {
		in       string
		expected string
	}{
		{"foo", ""},
		{"https://example.com", "example.com"},
		{"HTTP://EXAMPLE.COM", "example.com"},
		{"example.com", ""}, // TODO: make this test work
		{"ftp://example.com/path/to/nowhere.html", "example.com"},
	}

	for i, tc := range tt {
		require.Equal(t, tc.expected, index.NormalizeURL(tc.in), "test case %d failed", i)
	}
}

func TestParseQuery(t *testing.T) {
	query1 := map[string]interface{}{
		"name":     "BobVASP",
		"country":  "Guyana",
		"website":  "http://example.com/",
		"category": "DEX",
	}

	query2 := map[string]interface{}{
		"name":     []string{"BobVASP"},
		"country":  []string{"Guyana"},
		"website":  []string{"http://example.com/"},
		"category": []string{"DEX"},
	}

	queryTypes := map[string]map[string]interface{}{
		"string": query1,
		"slice":  query2,
	}

	for kind, query := range queryTypes {
		// Test index not in query
		vals, ok := index.ParseQuery("foo", query, index.Normalize)
		require.False(t, ok, "did not return not found for query type %q", kind)
		require.Nil(t, vals, "did not return not found for query type %q", kind)

		// Test names with normalize
		vals, ok = index.ParseQuery("name", query, index.Normalize)
		require.True(t, ok, "test failed for query type %q", kind)
		require.Equal(t, []string{"bobvasp"}, vals, "test failed for query type %q", kind)

		// Test country with normalize country
		vals, ok = index.ParseQuery("country", query, index.NormalizeCountry)
		require.True(t, ok, "test failed for query type %q", kind)
		require.Equal(t, []string{"GY"}, vals, "test failed for query type %q", kind)

		// Test website with normalize url
		vals, ok = index.ParseQuery("website", query, index.NormalizeURL)
		require.True(t, ok, "test failed for query type %q", kind)
		require.Equal(t, []string{"example.com"}, vals, "test failed for query type %q", kind)

		// Test category with no normalization
		vals, ok = index.ParseQuery("category", query, nil)
		require.True(t, ok, "test failed for query type %q", kind)
		require.Equal(t, []string{"DEX"}, vals, "test failed for query type %q", kind)
	}

	// Test Bad Parse Query
	query3 := map[string]interface{}{
		"name":     1,
		"country":  1,
		"website":  1,
		"category": 1,
	}
	vals, ok := index.ParseQuery("country", query3, index.NormalizeCountry)
	require.False(t, ok, "bad query type returned ok")
	require.Nil(t, vals, "bad query type returned values")
}

func TestNameSearch(t *testing.T) {
	// Search an empty index
	idx := index.NewNamesIndex()
	results := idx.Search(queryFixture)
	require.Empty(t, results, "failed search test with empty index")

	aliceID := uuid.NewString()
	bobID := uuid.NewString()

	// Add a term that matches the search query
	idx.Add("Alice VASP, LLC", aliceID)
	results = idx.Search(queryFixture)
	require.Len(t, results, 1, "failed search test with one entry")
	require.Contains(t, results, aliceID, "search doesn't contain alice")

	// Add a term that does not match the search query
	idx.Add("bob", bobID)
	results = idx.Search(queryFixture)
	require.Len(t, results, 1, "failed search test with two entries")
	require.NotContains(t, results, bobID, "search shouldn't contain bob")

	// Add the same term that does match the search query
	idx.Add("BOBVASP", bobID)
	results = idx.Search(queryFixture)
	require.Len(t, results, 2, "failed search test with three entries")
	require.Contains(t, results, aliceID, "search doesn't contain alice")
	require.Contains(t, results, bobID, "search doesn't contain bob")

	// Add multiple terms that will match the search query
	idx.Add("BOBVASP Industries", bobID)
	idx.Add("Alice", aliceID)
	idx.Add("Just Some Junk", uuid.NewString())
	idx.Add("Alice Again", aliceID)
	results = idx.Search(queryFixture)
	require.Len(t, results, 2, "failed search test with multiple entries")
	require.Contains(t, results, aliceID, "search doesn't contain alice")
	require.Contains(t, results, bobID, "search doesn't contain bob")
}

func TestWebsiteSearch(t *testing.T) {
	// Search an empty index
	idx := index.NewWebsiteIndex()
	results := idx.Search(queryFixture)
	require.Empty(t, results, "failed search test with empty index")

	aliceID := uuid.NewString()
	bobID := uuid.NewString()

	// Add a term that matches the search query
	idx.Add("https://alice.vaspbot.net", aliceID)
	results = idx.Search(queryFixture)
	require.Len(t, results, 1, "failed search test with one entry")
	require.Contains(t, results, aliceID, "search doesn't contain alice")

	// Add a term that does not match the search query
	idx.Add("http://bob.vaspbot.net", bobID)
	results = idx.Search(queryFixture)
	require.Len(t, results, 1, "failed search test with two entries")
	require.NotContains(t, results, bobID, "search shouldn't contain bob")

	newquery := map[string]interface{}{
		"website": []string{"https://alice.vaspbot.net", "http://alice.us/info.html", "https://bob.vaspbot.net/index.html"},
	}

	// Add the same term that does match the search query
	results = idx.Search(newquery)
	require.Len(t, results, 2, "failed search test with three entries")
	require.Contains(t, results, aliceID, "search doesn't contain alice")
	require.Contains(t, results, bobID, "search doesn't contain bob")

	// Add multiple terms that will match the search query
	idx.Add("http://bob.co.uk", bobID)
	idx.Add("https://alice.us/index.html", aliceID)
	idx.Add("https://example.com", uuid.NewString())
	idx.Add("http://alice.us", aliceID)
	results = idx.Search(newquery)
	require.Len(t, results, 2, "failed search test with multiple entries")
	require.Contains(t, results, aliceID, "search doesn't contain alice")
	require.Contains(t, results, bobID, "search doesn't contain bob")
}
