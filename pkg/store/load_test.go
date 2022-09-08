package store_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/store"
	"github.com/trisacrypto/directory/pkg/store/mockdb"
)

// Test that the Load functions fails to load improperly formatted files.
func TestLoadInvalid(t *testing.T) {
	db := mockdb.GetStore()

	// Invalid path
	require.Error(t, store.Load(db, filepath.Join("testdata", "invalid", "invalid.csv")))

	// Invalid CSV file
	require.Error(t, store.Load(db, filepath.Join("testdata", "bad_format.csv")))

	// Invalid country code
	require.Error(t, store.Load(db, filepath.Join("testdata", "invalid_country.csv")))

	// Invalid url
	require.Error(t, store.Load(db, filepath.Join("testdata", "invalid_url.csv")))

	// No records, should not error
	require.NoError(t, store.Load(db, filepath.Join("testdata", "empty.csv")))

	// No store calls
	require.False(t, mockdb.GetState().CreateVASPInvoked, "CreateVASP should not have been invoked")
	require.False(t, mockdb.GetState().RetrieveVASPInvoked, "RetrieveVASP should not have been invoked")
}

// Test that the Load function correctly loads VASPs from a valid CSV file.
func TestLoad(t *testing.T) {
	db := mockdb.GetStore()
	defer mockdb.ResetState()

	// Load a valid CSV file
	require.Nil(t, store.Load(db, filepath.Join("testdata", "vasps.csv")))
	state := mockdb.GetState()
	require.True(t, state.CreateVASPInvoked, "CreateVASP should have been invoked")
	require.True(t, state.RetrieveVASPInvoked, "RetrieveVASP should have been invoked")

	// Should be new VASPs in the store
	// Echo Funds should be skipped because there is no url provided
	require.Len(t, state.VASPs, 2)
	require.Len(t, state.Keys, 2)

	// Check that the VASPs were correctly loaded into the store
	charlieVASP, err := db.RetrieveVASP(state.Keys[0])
	require.NoError(t, err)
	require.NotEmpty(t, charlieVASP.Id)
	require.Equal(t, "CharlieBank", charlieVASP.Entity.Name.NameIdentifiers[0].LegalPersonName)
	require.Equal(t, "2140 Carson Mission Apt. 731", charlieVASP.Entity.GeographicAddresses[0].AddressLine[0])
	require.Equal(t, "https://trisa.charliebank.io", charlieVASP.Website)
	require.Equal(t, "CA", charlieVASP.Entity.CountryOfRegistration)
	require.Equal(t, "CA", charlieVASP.Entity.GeographicAddresses[0].Country)
	require.Equal(t, "trisa.charliebank.io", charlieVASP.CommonName)

	deltaVASP, err := db.RetrieveVASP(state.Keys[1])
	require.NoError(t, err)
	require.NotEmpty(t, deltaVASP.Id)
	require.NotEqual(t, charlieVASP.Id, deltaVASP.Id, "VASP IDs should be unique")
	require.Equal(t, "Delta Assets", deltaVASP.Entity.Name.NameIdentifiers[0].LegalPersonName)
	require.Equal(t, "0806 Neal Coves Suite 610", deltaVASP.Entity.GeographicAddresses[0].AddressLine[0])
	require.Equal(t, "http://trisa.delta.io", deltaVASP.Website)
	require.Equal(t, "XX", deltaVASP.Entity.CountryOfRegistration)
	require.Equal(t, "XX", deltaVASP.Entity.GeographicAddresses[0].Country)
	require.Equal(t, "trisa.delta.io", deltaVASP.CommonName)
}
