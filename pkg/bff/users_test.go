package bff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff"
)

func TestAppDataHelpers(t *testing.T) {
	// Test AppData helpers to fetch user data stored in Auth0
	// Create a valid appdata data structure
	appdata := map[string]interface{}{
		"orgid": "a2de6ede-1706-4746-8e85-d85a5fbc203e",
		"vasps": map[string]string{
			"testnet": "c87df2af-93aa-4281-996e-1c13bcd731d0",
			"mainnet": "9c84beb3-dcad-4e85-be55-2378fd226123",
		},
		"color": "red",
	}

	// Test GetOrgID
	orgId, ok := GetOrgID(appdata)
	require.True(t, ok, "should be able to fetch org id from valid appdata")
	require.Equal(t, "a2de6ede-1706-4746-8e85-d85a5fbc203e", orgId, "should be able to fetch correct org id from valid appdata")

	// Test GetVASPs
	vasps, ok := GetVASPs(appdata)
	require.True(t, ok, "should be able to fetch vasps from valid appdata")
	require.Contains(t, vasps, "testnet", "vasps should contain testnet")
	require.Equal(t, "c87df2af-93aa-4281-996e-1c13bcd731d0", vasps["testnet"])
	require.Contains(t, vasps, "mainnet", "vasps should contain mainnet")
	require.Equal(t, "9c84beb3-dcad-4e85-be55-2378fd226123", vasps["mainnet"])

	// Test AppData with no data in it
	empty := make(map[string]interface{})

	orgId, ok = GetOrgID(empty)
	require.False(t, ok, "shouldn't be able to fetch org id from empty appdata")
	require.Empty(t, orgId, "should not return an orgId when appadata is empty")

	vasps, ok = GetVASPs(empty)
	require.False(t, ok, "shouldn't be able to fetch vasps from empty appdata")
	require.Nil(t, vasps, "should not return vasps when appdata is empty")

	// Test AppData with wrong types in it
	invalid := map[string]interface{}{
		"orgid": 42,
		"vasp":  21,
		"color": "red",
	}

	vasps, ok = GetVASPs(invalid)
	require.False(t, ok, "shouldn't be able to fetch vasps from invalid appdata")
	require.Nil(t, vasps, "should not return vasps when appdata is invalid")

	orgId, ok = GetOrgID(invalid)
	require.False(t, ok, "shouldn't be able to fetch org id from invalid appdata")
	require.Empty(t, orgId, "should not return an orgId when appdata is invalid")
}

func TestMapEqual(t *testing.T) {
	testCases := []map[string]string{
		{"a": "b", "b": "c", "d": "e"},
		{"a": "b", "b": "c"},
		{"a": "c", "b": "f", "d": "12"},
		{"c": "d", "e": "f", "j": "k"},
		{"c": "d"},
		{"q": "r", "u": "t"},
	}

	for i, m1 := range testCases {
		for j, m2 := range testCases {
			if i == j {
				require.True(t, MapEqual(m1, m2), "expected maps %d and %d to be equal but they weren't", i, j)
			} else {
				require.False(t, MapEqual(m1, m2), "expected maps %d and %d to not be equal", i, j)
			}
		}

		// Compare with nil and empty map
		require.False(t, MapEqual(m1, nil), "expected map %d to not be equal with nil", i)
		require.False(t, MapEqual(m1, make(map[string]string)), "expected map %d to not be equal with empty map", i)
	}

}
