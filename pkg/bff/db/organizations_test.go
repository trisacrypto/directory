package db_test

func (s *dbTestSuite) TestOrganizations() {
	// Test basic interactions with the organizations collection
	require := s.Require()

	// Create an organization
	org, err := s.db.CreateOrganization()
	require.NoError(err, "could not create an empty organization in the database")
	require.NotEmpty(org.Id, "no uuid was created for the organization")
	require.NotEmpty(org.Created, "no created timestamp was added on the organization")
	require.NotEmpty(org.Modified, "no modified timestamp was added on the organization")
	require.Equal(org.Created, org.Modified, "created and modified timestamps should be identical on create")

	// Update an organization
	org.Name = "BestCoin SuperFun"
	err = s.db.UpdateOrganization(org)
	require.NoError(err, "could not update organization")

	// Retrieve an organization
	ret, err := s.db.RetrieveOrganization(org.Id)
	require.NoError(err, "could not retrieve organization")
	require.NotNil(ret, "no organization returned from retrieve")
	require.NotSame(org, ret, "a new organization should have been returned")

	require.Equal(org.Id, ret.Id, "original and retrieved Id should be the same")
	require.Equal(org.Name, ret.Name, "original and retrieved name should be the same")
	require.Equal(org.Created, ret.Created, "original and retrieved created should be the same")
	require.NotEqual(org.Created, ret.Modified, "original created and retrieved modified should not be the same")

	// Delete an organization
	err = s.db.DeleteOrganization(ret.Id)
	require.NoError(err, "could not delete organization")

	_, err = s.db.RetrieveOrganization(org.Id)
	require.Error(err, "should error trying to retrieve deleted organization")
}
