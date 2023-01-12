package bff_test

import (
	"context"
	"net/http"
	"time"

	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
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
	claims.Permissions = []string{auth.CreateOrganizations}
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

	// Invalid domains are rejected
	params.Domain = "alicevasp"
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusBadRequest, "invalid domain provided", "expected error when domain is invalid")

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
	require.Equal(authtest.Name, org.CreatedBy, "expected created by to be set")
	require.Equal(params.Name, org.Name, "organization name does not match")
	require.Equal(params.Domain, org.Domain, "organization domain does not match")

	// User should be added as a collaborator
	require.Len(org.Collaborators, 1, "expected one collaborator")
	collab := org.GetCollaborator(claims.Email)
	require.NotNil(collab, "expected user to exist as collaborator in new organization")
	require.Equal(authtest.Email, collab.Email, "expected collaborator email to match")
	require.Equal(authtest.UserID, collab.UserId, "expected collaborator user id to match")
	require.True(collab.Verified, "expected collaborator to be verified")

	// User app metadata should be updated with the organization id
	metadata := &auth.AppMetadata{}
	require.NoError(metadata.Load(s.auth.GetUserAppMetadata()))
	require.Len(metadata.Organizations, 1, "expected user to be a member of one organization")
	require.Equal(reply.ID, metadata.Organizations[0], "expected user to be a member of the organization")

	// Should not be able to create an organization with the same domain
	metadata.Organizations = []string{reply.ID}
	appdata, err := metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusConflict, "organization with domain already exists", "expected error when organization already exists")

	// Uniqueness check should be case insensitive
	params.Domain = "ALICEVASP.IO"
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusConflict, "organization with domain already exists", "expected error when organization already exists")

	// Uniqueness check should ignore leading and trailing whitespace
	params.Domain = " aliceVASP.io "
	_, err = s.client.CreateOrganization(context.TODO(), params)
	s.requireError(err, http.StatusConflict, "organization with domain already exists", "expected error when organization already exists")

	// Should not return an error if there is an organization on the app metadata that's not in the database
	metadata.Organizations = []string{"00000000-0000-0000-0000-000000000000"}
	params.Domain = " bobVASP.io "
	appdata, err = metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	require.NoError(s.DB().DeleteOrganization(org.UUID()), "could not delete organization from database")
	reply, err = s.client.CreateOrganization(context.TODO(), params)
	require.NoError(err, "create organization call failed")
	require.NotEmpty(reply.ID, "expected organization id to be set")
	require.Equal(params.Name, reply.Name, "expected name in reply to match")
	require.Equal("bobvasp.io", reply.Domain, "expected domain in reply to match")
	require.NotEmpty(reply.CreatedAt, "expected created at timestamp to be set")
	require.True(reply.RefreshToken, "refresh token should be set")
	org, err = s.bff.OrganizationFromID(reply.ID)
	require.NoError(err, "could not find organization in database")
	require.Equal(params.Name, org.Name, "organization name does not match")
	require.Equal("bobvasp.io", org.Domain, "organization domain does not match")

	// User should be added as a collaborator
	require.Len(org.Collaborators, 1, "expected one collaborator")
	collab = org.GetCollaborator(claims.Email)
	require.NotNil(collab, "expected user to exist as collaborator in new organization")
	require.Equal(authtest.Email, collab.Email, "expected collaborator email to match")
	require.Equal(authtest.UserID, collab.UserId, "expected collaborator user id to match")
	require.True(collab.Verified, "expected collaborator to be verified")

	// User app metadata should be updated with the organization id
	metadata = &auth.AppMetadata{}
	require.NoError(metadata.Load(s.auth.GetUserAppMetadata()))
	require.Len(metadata.Organizations, 2, "wrong number of organizations in app metadata")
	require.Equal(reply.ID, metadata.Organizations[1], "expected user to be a member of the organization")
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
	_, err := s.client.ListOrganizations(context.TODO(), &api.ListOrganizationsParams{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:collaborator permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.ListOrganizations(context.TODO(), &api.ListOrganizationsParams{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Should return empty response when user has no organizations
	claims.Permissions = []string{auth.ReadOrganizations}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")
	expected := &api.ListOrganizationsReply{
		Organizations: []*api.OrganizationReply{},
		Count:         0,
		Page:          1,
		PageSize:      8,
	}
	reply, err := s.client.ListOrganizations(context.TODO(), &api.ListOrganizationsParams{})
	require.NoError(err, "list organizations call failed")
	require.Equal(expected, reply, "expected default response")

	// Should enforce defaults for query parameters
	req := &api.ListOrganizationsParams{
		Page:     -1,
		PageSize: -1,
	}
	reply, err = s.client.ListOrganizations(context.TODO(), req)
	require.NoError(err, "list organizations call failed")
	require.Equal(expected, reply, "expected default response for invalid query parameters")

	// Should not return an error if there is an organization on the app metadata that's not in the database
	metadata := &auth.AppMetadata{}
	require.NoError(metadata.Load(s.auth.GetUserAppMetadata()))
	metadata.Organizations = []string{"00000000-0000-0000-0000-000000000000"}
	appdata, err := metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	reply, err = s.client.ListOrganizations(context.TODO(), &api.ListOrganizationsParams{})
	require.NoError(err, "list organizations call failed")
	require.Equal(expected, reply, "expected default response")

	// Create some organizations for the user
	alice := &models.Organization{
		Name:   "Alice VASP",
		Domain: "alicevasp.io",
	}
	aliceCollab := &models.Collaborator{
		Email: claims.Email,
	}
	require.NoError(alice.AddCollaborator(aliceCollab))
	_, err = s.DB().CreateOrganization(alice)
	require.NoError(err, "could not create organization")

	bob := &models.Organization{
		Name:   "Bob VASP",
		Domain: "bobvasp.io",
	}
	bobCollab := &models.Collaborator{
		Email:     claims.Email,
		LastLogin: time.Now().Format(time.RFC3339Nano),
	}
	require.NoError(bob.AddCollaborator(bobCollab))
	_, err = s.DB().CreateOrganization(bob)
	require.NoError(err, "could not create organization")

	charlie := &models.Organization{
		Name:   "Charlie VASP",
		Domain: "charlievasp.io",
	}
	charlieCollab := &models.Collaborator{
		Email:     claims.Email,
		LastLogin: time.Now().Format(time.RFC3339Nano),
	}
	require.NoError(charlie.AddCollaborator(charlieCollab))
	_, err = s.DB().CreateOrganization(charlie)
	require.NoError(err, "could not create organization")

	delta := &models.Organization{
		Name:   "Delta VASP",
		Domain: "deltavasp.io",
	}
	_, err = s.DB().CreateOrganization(delta)
	require.NoError(err, "could not create organization")

	// Update the app metadata to contain the organizations
	metadata.Organizations = []string{alice.Id, bob.Id, charlie.Id, delta.Id}
	appdata, err = metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)

	orgReplies := []*api.OrganizationReply{
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
			LastLogin: bobCollab.LastLogin,
		},
		{
			ID:        charlie.Id,
			Name:      charlie.Name,
			Domain:    charlie.Domain,
			CreatedAt: charlie.Created,
			LastLogin: charlieCollab.LastLogin,
		},
	}
	expected.Organizations = orgReplies
	expected.Count = 3

	// Should return all organizations the user is a collaborator on
	// If the user is not a collaborator, the endpoint should not return an error
	reply, err = s.client.ListOrganizations(context.TODO(), &api.ListOrganizationsParams{})
	require.NoError(err, "list organizations call failed")
	require.Equal(expected, reply, "expected returned organizations to match")

	// List organizations across multiple pages
	req = &api.ListOrganizationsParams{
		Page:     1,
		PageSize: 2,
	}
	expected.Organizations = orgReplies[0:2]
	expected.PageSize = 2
	reply, err = s.client.ListOrganizations(context.TODO(), req)
	require.NoError(err, "list organizations call failed")
	require.Equal(expected, reply, "wrong results for page 1")

	req.Page = 2
	expected.Organizations = orgReplies[2:3]
	expected.Page = 2
	reply, err = s.client.ListOrganizations(context.TODO(), req)
	require.NoError(err, "list organizations call failed")
	require.Equal(expected, reply, "wrong results for page 2")

	// Should return no organizations if the page is out of bounds
	req.Page = 3
	expected.Organizations = []*api.OrganizationReply{}
	expected.Page = 3
	reply, err = s.client.ListOrganizations(context.TODO(), req)
	require.NoError(err, "list organizations call failed")
	require.Equal(expected, reply, "wrong results for page out of bounds")
}
