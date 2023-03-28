package trtl_test

import (
	"context"

	store "github.com/trisacrypto/directory/pkg/store/trtl"
)

func (s *trtlStoreTestSuite) TestIndexSync() {
	require := s.Require()
	require.NoError(s.grpc.Connect(context.Background()), "could not connect to grpc bufconn")
	defer s.grpc.Close()

	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err, "could not create mock trtl store")

	// Ensure that indices are all empty
	require.True(db.GetNamesIndex().Empty(), "name index not empty, have fixtures changed?")
	require.True(db.GetWebsitesIndex().Empty(), "website index not empty, have fixtures changed?")
	require.True(db.GetCountriesIndex().Empty(), "country index not empty, have fixtures changed?")
	require.True(db.GetCategoriesIndex().Empty(), "category index not empty, have fixtures changed?")

	// Create a bunch of records
	err = createVASPs(db, 100, 1)
	require.NoError(err, "could not create 100 vasps for index tests")
	require.Equal(200, db.GetNamesIndex().Len(), "names index has an unexpected length")
	require.Equal(100, db.GetWebsitesIndex().Len(), "website index has an unexpected length")
	require.Equal(7, db.GetCountriesIndex().Len(), "countries index has an unexpected length")
	require.Equal(3, db.GetCategoriesIndex().Len(), "categories index has an unexpected length")

	// Sync the indices to disk
	// NOTE: this should also test any conflicts with reserved namespaces in trtl
	require.NoError(db.Sync("all"), "could not sync indices to disk")

	// TODO: check that we can load the indices from disk

	// TODO: check updating records

	// Test deleting the records clears the indices
	require.NoError(deleteVASPs(db), "could not delete vasps during index test")
	require.True(db.GetNamesIndex().Empty(), "name index not empty after delete")
	require.True(db.GetWebsitesIndex().Empty(), "website index not empty after delete")

	// TODO: multi-index has empty arrays but still contains country/categories
	// require.True(db.GetCountriesIndex().Empty(), "country index not empty after delete")
	// require.True(db.GetCategoriesIndex().Empty(), "category index not empty after delete")
}

func (s *trtlStoreTestSuite) TestSearch() {
	require := s.Require()
	require.NoError(s.grpc.Connect(context.Background()), "could not connect to grpc bufconn")
	defer s.grpc.Close()

	db, err := store.NewMock(s.grpc.Conn)
	require.NoError(err, "could not create mock trtl store")

	// Create a bunch of records removing any records that were there before
	err = createVASPs(db, 100, 1)
	require.NoError(err, "could not create 100 vasps for search test")

	// Test a simple search
	query := map[string]interface{}{
		"name":     []string{"Test VASP 00A1", "Test VASP F32A", "Test VASP 0014"},
		"website":  []string{"https://test0003.net/", "https://test00FA.net"},
		"country":  "CC",
		"category": "PRIVATE_ORGANIZATION",
	}

	vasps, err := db.SearchVASPs(context.Background(), query)
	require.NoError(err, "could not search vasps with query")
	require.Len(vasps, 1, "no vasps returned from search")
	require.Equal("trisa0003.test.net", vasps[0].CommonName)
	require.NoError(deleteVASPs(db), "could not delete vasps after search test")
}
