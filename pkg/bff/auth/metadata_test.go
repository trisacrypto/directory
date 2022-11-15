package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff/auth"
)

func TestAppMetadata(t *testing.T) {
	// Test loading and dumping app_metadata to/from the auth0 response.
	testCases := []struct {
		appdata  map[string]interface{}
		expected *AppMetadata
	}{
		{
			nil, &AppMetadata{},
		},
		{
			map[string]interface{}{}, &AppMetadata{},
		},
		{
			map[string]interface{}{
				"orgid": "67428be4-3fa4-4bf2-9e15-edbf043f8670",
			},
			&AppMetadata{
				OrgID: "67428be4-3fa4-4bf2-9e15-edbf043f8670",
			},
		},
		{
			map[string]interface{}{
				"vasps": map[string]string{
					"testnet": "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
				},
			},
			&AppMetadata{
				VASPs: VASPs{
					TestNet: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
				},
			},
		},
		{
			map[string]interface{}{
				"vasps": map[string]string{
					"mainnet": "2ac8d50a-ff4c-479e-8eec-a35d96d90911",
				},
			},
			&AppMetadata{
				VASPs: VASPs{
					MainNet: "2ac8d50a-ff4c-479e-8eec-a35d96d90911",
				},
			},
		},
		{
			map[string]interface{}{
				"vasps": map[string]string{
					"testnet": "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
					"mainnet": "2ac8d50a-ff4c-479e-8eec-a35d96d90911",
				},
			},
			&AppMetadata{
				VASPs: VASPs{
					TestNet: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
					MainNet: "2ac8d50a-ff4c-479e-8eec-a35d96d90911",
				},
			},
		},
		{
			map[string]interface{}{
				"orgid": "67428be4-3fa4-4bf2-9e15-edbf043f8670",
				"vasps": map[string]string{
					"testnet": "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
					"mainnet": "2ac8d50a-ff4c-479e-8eec-a35d96d90911",
				},
			},
			&AppMetadata{
				OrgID: "67428be4-3fa4-4bf2-9e15-edbf043f8670",
				VASPs: VASPs{
					TestNet: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
					MainNet: "2ac8d50a-ff4c-479e-8eec-a35d96d90911",
				},
			},
		},
	}

	for _, tc := range testCases {
		actual := &AppMetadata{}

		err := actual.Load(tc.appdata)
		require.NoError(t, err, "could not load appdata")
		require.Equal(t, tc.expected, actual, "app_metadata did not load correctly")

		appdata, err := actual.Dump()
		require.NoError(t, err, "could not dump app_metdata")

		require.Contains(t, appdata, "orgid")
		require.Equal(t, actual.OrgID, appdata["orgid"])

		require.Contains(t, appdata, "vasps")
		vasps, ok := appdata["vasps"].(map[string]interface{})
		require.True(t, ok, "appdata vasps is wrong type")

		require.Contains(t, appdata["vasps"], "testnet")
		require.Equal(t, actual.VASPs.TestNet, vasps["testnet"])
		require.Contains(t, appdata["vasps"], "mainnet")
		require.Equal(t, actual.VASPs.MainNet, vasps["mainnet"])
	}

}

func TestEquals(t *testing.T) {
	a := &AppMetadata{
		OrgID: "67428be4-3fa4-4bf2-9e15-edbf043f8670",
		VASPs: VASPs{
			TestNet: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
			MainNet: "2ac8d50a-ff4c-479e-8eec-a35d96d90911",
		},
		Organizations: []string{"67428be4-3fa4-4bf2-9e15-edbf043f8670"},
	}

	b := &AppMetadata{
		OrgID: "b4f2e15e-dbf0-3f86-428b-4e3fa4f2e15e",
		VASPs: VASPs{
			TestNet: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
			MainNet: "2ac8d50a-ff4c-479e-8eec-a35d96d90911",
		},
		Organizations: []string{"67428be4-3fa4-4bf2-9e15-edbf043f8670"},
	}

	// If the orgid is different, they are not equal
	require.False(t, a.Equals(b), "orgids are different, but Equals returned true")

	// If the VASPs are different, they are not equal
	b.OrgID = a.OrgID
	b.VASPs.MainNet = "00000000-0000-0000-0000-000000000000"
	require.False(t, a.Equals(b), "VASPs are different, but Equals returned true")

	// If the Organizations are different, they are not equal
	b.VASPs.MainNet = a.VASPs.MainNet
	b.Organizations = []string{
		"00000000-0000-0000-0000-000000000000",
		"67428be4-3fa4-4bf2-9e15-edbf043f8670",
	}
	require.False(t, a.Equals(b), "Organizations are different, but Equals returned true")
	require.False(t, b.Equals(a), "Organizations are different, but Equals returned true")
	b.Organizations = []string{"00000000-0000-0000-0000-000000000000"}
	require.False(t, a.Equals(b), "Organizations are different, but Equals returned true")

	// Only equal if all fields are equal
	b.Organizations = []string{"67428be4-3fa4-4bf2-9e15-edbf043f8670"}
	require.True(t, a.Equals(b), "all fields are equal, but Equals returned false")
}
