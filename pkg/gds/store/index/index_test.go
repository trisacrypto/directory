package index_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/store/index"
)

var queryFixture = map[string]interface{}{
	"name":     []string{"Alice", "BobVASP"},
	"website":  "https://alice.vaspbot.net",
	"country":  []string{"KY", "Cayman Islands", "Guyana"},
	"category": []string{"ATM", "PRIVATE_ORGANIZATION"},
}

func TestNameIndex(t *testing.T) {
	// Load the index from disk
	index := index.NewNamesIndex()
	require.True(t, index.Empty(), "new index is not empty")

	data, err := ioutil.ReadFile("testdata/names.json.gz")
	require.NoError(t, err, "could not read testdata/names.json.gz")
	require.NoError(t, index.Load(data), "could not load index from disk")
	require.False(t, index.Empty(), "no fixtures were loaded")
	require.Equal(t, 11, index.Len(), "fixtures length has changed")

	// Test dumping the data matches the data loaded
	dump, err := index.Dump()
	require.NoError(t, err, "could not dump index")
	require.True(t, bytes.Equal(data, dump), "dump does not match loaded data")

	// Check an existing ID in the fixture
	bobID := "69324932-286f-4708-abaa-2bb3a5df9557"
	val, ok := index.Find("bobvasp")
	require.True(t, ok, "could not find bob")
	require.Equal(t, bobID, val, "incorrect ID returned, has fixture changed?")

	vals, ok := index.Reverse(bobID)
	require.True(t, ok, "could not reverse bob")
	require.Len(t, vals, 4, "could not load bob correctly")

	// Add, overwrite and remove an item
	ok = index.Add("Foo", "7a4b54aa-17e7-4175-ba74-b34a0d9e97c8")
	require.True(t, ok, "could not add an item")

	val, ok = index.Find("FOO")
	require.True(t, ok, "could not find added item")
	require.Equal(t, "7a4b54aa-17e7-4175-ba74-b34a0d9e97c8", val, "incorrect item returned")

	ok = index.Overwrite("    FOO   ", "bar")
	require.True(t, ok, "could not overwrite item")

	val, ok = index.Find("foo")
	require.True(t, ok, "could not find overwritten item")
	require.Equal(t, "bar", val, "incorrect item returned")

	ok = index.Remove("   Foo  ")
	require.True(t, ok, "could not remove item")

	val, ok = index.Find("foo")
	require.False(t, ok, "found removed item")
	require.Equal(t, "", val, "empty string not returned")

	// Test search functionality
	results := index.Search(queryFixture)
	require.NotEmpty(t, results, "search returned no results from fixture")
	require.Equal(t, []string{"69324932-286f-4708-abaa-2bb3a5df9557", "7510e5cd-e63e-4218-b867-6a1cf08cb691"}, results, "unexpected results returned")
}

func TestWebsiteIndex(t *testing.T) {
	// Load the index from disk
	index := index.NewWebsiteIndex()
	require.True(t, index.Empty(), "new index is not empty")

	data, err := ioutil.ReadFile("testdata/websites.json.gz")
	require.NoError(t, err, "could not read testdata/websites.json.gz")
	require.NoError(t, index.Load(data), "could not load index from disk")
	require.False(t, index.Empty(), "no fixtures were loaded")
	require.Equal(t, 3, index.Len(), "fixtures length has changed")

	// Test dumping the data matches the data loaded
	dump, err := index.Dump()
	require.NoError(t, err, "could not dump index")
	require.True(t, bytes.Equal(data, dump), "dump does not match loaded data")

	// Check an existing ID in the fixture
	bobID := "69324932-286f-4708-abaa-2bb3a5df9557"
	val, ok := index.Find("https://bob.vaspbot.net")
	require.True(t, ok, "could not find bob.vaspbot.net")
	require.Equal(t, bobID, val, "incorrect ID returned, has fixture changed?")

	vals, ok := index.Reverse(bobID)
	require.True(t, ok, "could not reverse bob")
	require.Len(t, vals, 1, "could not load bob correctly")

	// Add, overwrite and remove an item
	ok = index.Add("https://example.com", "7a4b54aa-17e7-4175-ba74-b34a0d9e97c8")
	require.True(t, ok, "could not add an item")

	val, ok = index.Find("http://example.com")
	require.True(t, ok, "could not find added item")
	require.Equal(t, "7a4b54aa-17e7-4175-ba74-b34a0d9e97c8", val, "incorrect item returned")

	ok = index.Overwrite("FTP://EXAMPLE.COM", "bar")
	require.True(t, ok, "could not overwrite item")

	val, ok = index.Find("https://example.com/help")
	require.True(t, ok, "could not find overwritten item")
	require.Equal(t, "bar", val, "incorrect item returned")

	ok = index.Remove("http://example.com:443/file.html#zoo")
	require.True(t, ok, "could not remove item")

	val, ok = index.Find("http://example.com")
	require.False(t, ok, "found removed item")
	require.Equal(t, "", val, "empty string not returned")

	// Test search functionality
	results := index.Search(queryFixture)
	require.NotEmpty(t, results, "search returned no results from fixture")
	require.Equal(t, []string{"7510e5cd-e63e-4218-b867-6a1cf08cb691"}, results, "unexpected results returned")
}

func TestCountryIndex(t *testing.T) {
	// Load the index from disk
	index := index.NewCountryIndex()
	require.True(t, index.Empty(), "new index is not empty")

	data, err := ioutil.ReadFile("testdata/countries.json.gz")
	require.NoError(t, err, "could not read testdata/countries.json.gz")
	require.NoError(t, index.Load(data), "could not load index from disk")
	require.False(t, index.Empty(), "no fixtures were loaded")
	require.Equal(t, 10, index.Len(), "fixtures length has changed")

	// Test dumping the data matches the data loaded
	dump, err := index.Dump()
	require.NoError(t, err, "could not dump index")
	require.True(t, bytes.Equal(data, dump), "dump does not match loaded data")

	// Check an existing ID in the fixture
	catID := "bd86def4-51ea-471d-ac32-7df0e3a9725b"
	vals, ok := index.Find("Guyana")
	require.True(t, ok, "could not find country")
	require.Len(t, vals, 9, "incorrect index length, has fixture changed?")
	require.Contains(t, vals, catID, "index doesn't contain value")

	vals, ok = index.Reverse(catID)
	require.True(t, ok, "could not find value in index")
	require.Len(t, vals, 1, "could not load index value correctly")

	// Add multiple items to the index
	require.True(t, index.Add("Svalbard and Jan Mayen", "foo"))
	require.True(t, index.Add("SJ", "bar"))
	require.True(t, index.Add("svalbard and jan mayen", "baz"))
	require.True(t, index.Contains("svalbard", "bar"))

	vals, ok = index.Find("Svalbard and Jan Mayen")
	require.True(t, ok)
	require.Equal(t, []string{"bar", "baz", "foo"}, vals)

	// Remove items from index
	require.True(t, index.Remove("SJ", "foo"))
	require.True(t, index.Remove("svalbard", "bar"))
	require.True(t, index.Remove("Svalbard and Jan Mayen", "baz"))

	vals, _ = index.Find("svalbard and jan mayen")
	require.Empty(t, vals, "index contains values even after removal")

	// Test search functionality
	results := index.Search(queryFixture)
	require.NotEmpty(t, results, "search returned no results from fixture")

	expected := []string{"154de709-a8d3-4520-9aeb-627ffbabc36d", "1f1e11d6-4bfc-48a2-9c1e-e91ff9d17dc3", "20e999eb-5d42-4d21-a4bd-fe62ea6feac3", "3da2b71c-eef0-4491-bf82-69f26597cb80", "6eb58428-3a28-4a1d-a43b-fafceabddc8a", "72e5a1da-c1fb-4f75-bf94-26b9edc429d2", "8cb7e45e-e512-4e53-8c4f-b2f2f9487cc2", "9ed2686f-3cb9-4b33-af6b-cdb1c15f9ca7", "a4945722-d814-4f78-9aaf-dc6bfa8e7c5f", "a7d5971e-b474-4c39-92fb-234cd7323aac", "bd86def4-51ea-471d-ac32-7df0e3a9725b", "c0bc2f50-96c9-48cb-bf60-49d5769b997d", "e6f88483-4234-4428-9380-af3db8d907be", "f0ad2704-df38-4f24-8a53-2d4fad7ca30a", "ff98e7a8-2299-47a1-a855-024895f92857"}
	require.Equal(t, expected, results, "unexpected results returned")
}

func TestCategoryIndex(t *testing.T) {
	// Load the index from disk
	index := index.NewCategoryIndex()
	require.True(t, index.Empty(), "new index is not empty")

	data, err := ioutil.ReadFile("testdata/categories.json.gz")
	require.NoError(t, err, "could not read testdata/categories.json.gz")
	require.NoError(t, index.Load(data), "could not load index from disk")
	require.False(t, index.Empty(), "no fixtures were loaded")
	require.Equal(t, 12, index.Len(), "fixtures length has changed")

	// Test dumping the data matches the data loaded
	dump, err := index.Dump()
	require.NoError(t, err, "could not dump index")
	require.True(t, bytes.Equal(data, dump), "dump does not match loaded data")

	// Check an existing ID in the fixture
	catID := "9213b7de-585b-40a3-b2d4-081c81a1b94b"
	vals, ok := index.Find("private_organization")
	require.True(t, ok, "could not find category")
	require.Len(t, vals, 8, "incorrect index length, has fixture changed?")
	require.Contains(t, vals, catID, "index doesn't contain value")

	vals, ok = index.Reverse(catID)
	require.True(t, ok, "could not find category ID")
	require.Len(t, vals, 2, "could not load category ID correctly")

	// Add multiple items to the index
	require.True(t, index.Add("red", "foo"))
	require.True(t, index.Add("RED", "bar"))
	require.True(t, index.Add("   Red ", "baz"))
	require.True(t, index.Contains("red", "bar"))

	vals, ok = index.Find("red")
	require.True(t, ok)
	require.Equal(t, []string{"bar", "baz", "foo"}, vals)

	// Remove items from index
	require.True(t, index.Remove(" red  ", "foo"))
	require.True(t, index.Remove("rEd  ", "bar"))
	require.True(t, index.Remove("RED", "baz"))

	vals, _ = index.Find("red")
	require.Empty(t, vals, "index contains values even after removal")

	// Test search functionality
	results := index.Search(queryFixture)
	require.NotEmpty(t, results, "search returned no results from fixture")
	expected := []string{"076ab2a0-8122-42ca-9e0d-ee4f364f313e", "2d250cb2-d628-4ad8-827a-8480cd307624", "471f3ac2-970c-4b1f-8e86-0e0960999ce3", "530b2aa6-93ce-4f95-b1a5-6f18fd048893", "59c76571-3cca-4b5f-a06f-0c30d318f6be", "63604c34-2760-484a-9e43-ac929c7142ea", "75063bb9-02dc-4890-bfad-35ceeeba411f", "9213b7de-585b-40a3-b2d4-081c81a1b94b", "b0dfabac-863f-4343-a1e7-7d702007db46", "b5ad0cee-1099-40f8-b8e1-de96a29ee17f", "b6ff7902-6614-41e5-ba2e-abd41f197ed6", "d3333a46-a4ab-4a41-9a79-9d1708e4ffa3", "ee9cd4ba-d81a-4ba3-b374-d04845c32912"}
	require.Equal(t, expected, results, "unexpected results returned")
}
