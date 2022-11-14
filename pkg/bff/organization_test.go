package bff_test

import (
	"context"
	"net/http"

	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
)

func (s *bffTestSuite) TestCreateOrganization() {
	require := s.Require()
	defer s.ResetDB()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"create:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
	_, err := s.client.CreateOrganization(context.TODO(), &api.OrganizationParams{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:collaborator permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.CreateOrganization(context.TODO(), &api.OrganizationParams{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Organization name is required
	claims.Permissions = []string{"create:organizations"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")
	params := &api.OrganizationParams{
		Domain: "alicevasp.io",
	}
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusBadRequest, "must provide name in request params", "expected error when name is not provided")

	// Organization domain is required
	params = &api.OrganizationParams{
		Name: "Alice VASP",
	}
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusBadRequest, "must provide domain in request params", "expected error when domain is not provided")

	// Valid request - organization should be created in the database
	params.Domain = "alicevasp.io"
	reply, err := s.client.CreateOrganization(context.TODO(), params)
	require.NoError(err, "create organization call failed")
	require.NotEmpty(reply.ID, "expected organization id to be set")
	require.Equal(params.Name, reply.Name, "expected name in reply to match")
	require.Equal(params.Domain, reply.Domain, "expected domain in reply to match")
	require.NotEmpty(reply.CreatedAt, "expected created at timestamp to be set")
	require.True(reply.RefreshToken, "refresh token should be set")
	org, err := s.bff.OrganizationFromID(reply.ID)
	require.NoError(err, "could not find organization in database")
	require.Equal(params.Name, org.Name, "organization name does not match")
	require.Equal(params.Domain, org.Domain, "organization domain does not match")

	// User app metadata should be updated to include the organization
	metadata := &auth.AppMetadata{}
	require.NoError(metadata.Load(s.auth.GetUserAppMetadata()))
	require.Equal(reply.ID, metadata.OrgID, "app metadata org id does not match")

	// Should not be able to create an organization with the same domain
	metadata.Organizations = []string{reply.ID}
	appdata, err := metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusConflict, "organization with domain name already exists", "expected error when organization already exists")

	// Uniqueness check should be case insensitive
	params.Domain = "ALICEVASP.IO"
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusConflict, "organization with domain name already exists", "expected error when organization already exists")

	// Uniqueness check should ignore leading and trailing whitespace
	params.Domain = " aliceVASP.io "
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusConflict, "organization with domain name already exists", "expected error when organization already exists")

	// Uniqueness check should ignore trailing periods
	params.Domain = "alicevasp.io."
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusConflict, "organization with domain name already exists", "expected error when organization already exists")

	// Should not return an error if there is an organization on the app metadata that's not in the database
	metadata.Organizations = []string{"00000000-0000-0000-0000-000000000000"}
	params.Domain = "bobvasp.io"
	appdata, err = metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	require.NoError(s.DB().DeleteOrganization(org.UUID()), "could not delete organization from database")
	reply, err = s.client.CreateOrganization(context.TODO(), params)
	require.NoError(err, "create organization call failed")
	require.NotEmpty(reply.ID, "expected organization id to be set")
	require.Equal(params.Name, reply.Name, "expected name in reply to match")
	require.Equal(params.Domain, reply.Domain, "expected domain in reply to match")
	require.NotEmpty(reply.CreatedAt, "expected created at timestamp to be set")
	require.True(reply.RefreshToken, "refresh token should be set")
	org, err = s.bff.OrganizationFromID(reply.ID)
	require.NoError(err, "could not find organization in database")
	require.Equal(params.Name, org.Name, "organization name does not match")
	require.Equal(params.Domain, org.Domain, "organization domain does not match")
}

func (s *bffTestSuite) TestListOrganizations() {
	require := s.Require()
	defer s.ResetDB()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := s.client.ListOrganizations(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:collaborator permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.ListOrganizations(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Should return empty response when user has no organizations
	claims.Permissions = []string{"read:organizations"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")
	reply, err := s.client.ListOrganizations(context.TODO())
	require.NoError(err, "list organizations call failed")
	require.Empty(reply, "expected empty response")

	// Should not return an error if there is an organization on the app metadata that's not in the database
	metadata := &auth.AppMetadata{}
	require.NoError(metadata.Load(s.auth.GetUserAppMetadata()))
	metadata.Organizations = []string{"00000000-0000-0000-0000-000000000000"}
	appdata, err := metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	reply, err = s.client.ListOrganizations(context.TODO())
	require.NoError(err, "list organizations call failed")
	require.Empty(reply, "expected empty response")

	// Create some organizations for the user
	alice, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization")
	alice.Name = "Alice VASP"
	alice.Domain = "alicevasp.io"
	require.NoError(s.DB().UpdateOrganization(alice), "could not update organization")

	bob, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization")
	bob.Name = "Bob VASP"
	bob.Domain = "bobvasp.io"
	require.NoError(s.DB().UpdateOrganization(bob), "could not update organization")

	charlie, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization")
	charlie.Name = "Charlie VASP"
	charlie.Domain = "charlievasp.io"
	require.NoError(s.DB().UpdateOrganization(charlie), "could not update organization")

	// Update the app metadata to contain the organizations
	metadata.Organizations = []string{alice.Id, bob.Id, charlie.Id}
	appdata, err = metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)

	expected := []*api.OrganizationReply{
		{
			ID:        alice.Id,
			Name:      alice.Name,
			Domain:    alice.Domain,
			CreatedAt: alice.Created,
		},
		{
			ID:        bob.Id,
			Name:      bob.Name,
			Domain:    bob.Domain,
			CreatedAt: bob.Created,
		},
		{
			ID:        charlie.Id,
			Name:      charlie.Name,
			Domain:    charlie.Domain,
			CreatedAt: charlie.Created,
		},
	}

	// Should return all organizations the user is a member of
	reply, err = s.client.ListOrganizations(context.TODO())
	require.NoError(err, "list organizations call failed")
	require.Equal(expected, reply, "expected returned organizations to match")
}
