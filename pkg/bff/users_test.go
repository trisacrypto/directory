package bff_test

import (
	"context"

	"github.com/trisacrypto/directory/pkg/bff"
)

func (s *bffTestSuite) TestListUserRoles() {
	require := s.Require()

	// Test listing the assignable roles
	expected := []string{
		bff.CollaboratorRole,
		bff.LeaderRole,
	}
	roles, err := s.client.ListUserRoles(context.TODO())
	require.NoError(err, "could not list assignable roles")
	require.Equal(expected, roles, "roles do not match")
}
