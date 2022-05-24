package bff_test

func (s *bffTestSuite) TestOverview() {
	// TODO: need to mock authentication before these tests will work.
	s.T().Skip("not implemented yet")

	// Test 401 with no access token
	// Test 401 authenticated user without read:vasp permission
	// Test 200 response with authenticated user with read:vasp permission
}
