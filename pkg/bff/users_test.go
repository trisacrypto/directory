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

func (s *bffTestSuite) TestNewUserLogin() {
	require := s.Require()
	defer s.ResetDB()
	defer s.auth.ResetUserAppMetadata()

	// Test the new user case - no orgID in params or app metadata
	claims := &authtest.Claims{
		Email: "leopold.wentzel@gmail.com",
	}
	metadata := &auth.AppMetadata{}
	require.NoError(s.SetClientCredentials(claims))
	appdata, err := metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	s.auth.SetUserRoles([]string{})
	require.NoError(s.client.Login(context.TODO(), nil))

	// Appdata should contain a new organization
	userMeta := &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.NotEmpty(userMeta.OrgID, "app metadata should contain a new organization")
	require.Empty(userMeta.VASPs.TestNet, "app metadata should not contain a testnet VASP")
	require.Empty(userMeta.VASPs.MainNet, "app metadata should not contain a mainnet VASP")
	require.Len(userMeta.Organizations, 0, "app metadata organization list should be empty")

	// User should exist as a collaborator in the new organization
	org, err := s.bff.OrganizationFromID(userMeta.OrgID)
	require.NoError(err, "could not get organization from ID")
	require.Equal(authtest.Name, org.CreatedBy, "organization created by should be set")
	require.Len(org.Collaborators, 1, "organization should have one collaborator")
	collab := org.GetCollaborator(claims.Email)
	require.NotNil(collab, "collaborator should exist in organization")
	require.Equal(claims.Email, collab.Email, "collaborator email should match")
	require.NotEmpty(collab.UserId, "collaborator user id should not be empty")
	require.True(collab.Verified, "collaborator should be verified")
	require.NotEmpty(collab.JoinedAt, "collaborator should have a joined at timestamp")
	require.Equal(collab.JoinedAt, collab.LastLogin, "collaborator should have a last login timestamp")

	// User should assume the leader role
	require.Equal([]string{auth.LeaderRole}, s.auth.GetUserRoles(), "user should have the leader role")

	// New TSP users should have the new organization added to the app metadata
	appdata, err = metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	s.auth.SetUserRoles([]string{"TRISA Service Provider"})
	claims.Permissions = []string{auth.SwitchOrganizations, auth.UpdateCollaborators}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	require.NoError(s.client.Login(context.TODO(), nil))

	// Appdata should contain a new organization and it should be added to the list
	userMeta = &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.NotEmpty(userMeta.OrgID, "app metadata should contain a new organization")
	require.Empty(userMeta.VASPs.TestNet, "app metadata should not contain a testnet VASP")
	require.Empty(userMeta.VASPs.MainNet, "app metadata should not contain a mainnet VASP")
	require.Len(userMeta.Organizations, 1, "app metadata organization list should contain one organization")
	require.Equal([]string{userMeta.OrgID}, userMeta.Organizations, "app metadata organization list should contain the new organization")

	// User should exist as a collaborator in the new organization
	org, err = s.bff.OrganizationFromID(userMeta.OrgID)
	require.NoError(err, "could not get organization from ID")
	require.Len(org.Collaborators, 1, "organization should have one collaborator")
	collab = org.GetCollaborator(claims.Email)
	require.NotNil(collab, "collaborator should exist in organization")
	require.Equal(claims.Email, collab.Email, "collaborator email should match")
	require.NotEmpty(collab.UserId, "collaborator user id should not be empty")
	require.True(collab.Verified, "collaborator should be verified")
	require.NotEmpty(collab.JoinedAt, "collaborator should have a joined at timestamp")
	require.Equal(collab.JoinedAt, collab.LastLogin, "collaborator should have a last login timestamp")

	// User should still have the TSP role
	require.Equal([]string{"TRISA Service Provider"}, s.auth.GetUserRoles(), "user should have the TSP role")

	// TODO: The endpoint should return a 200 with the refresh_token flag set
}

func (s *bffTestSuite) TestReturningUserLogin() {
	require := s.Require()
	defer s.ResetDB()
	defer s.auth.ResetUserAppMetadata()

	// Test the existing user case - orgID in app metadata but not params
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{auth.UpdateCollaborators},
	}
	metadata := &auth.AppMetadata{
		OrgID: "67428be4-3fa4-4bf2-9e15-edbf043f8670",
		VASPs: auth.VASPs{
			TestNet: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
		},
	}
	require.NoError(s.SetClientCredentials(claims))
	appdata, err := metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	s.auth.SetUserRoles([]string{auth.LeaderRole})

	// Returns an error if the organization does not exist
	err = s.client.Login(context.TODO(), nil)
	s.requireError(err, http.StatusNotFound, "organization not found")

	// Create the organization in the database without the collaborator
	org := &models.Organization{
		Id: metadata.OrgID,
		Testnet: &models.DirectoryRecord{
			Id: metadata.VASPs.TestNet,
		},
		Mainnet: &models.DirectoryRecord{
			Id: metadata.VASPs.MainNet,
		},
	}
	_, err = s.DB().CreateOrganization(org)
	require.NoError(err, "could not create organization")

	// User is not authorized to access the organization without being a collaborator
	err = s.client.Login(context.TODO(), nil)
	s.requireError(err, http.StatusUnauthorized, "user is not authorized to access this organization")

	// Make the user a TSP
	newOrg := &models.Organization{}
	_, err = s.DB().CreateOrganization(newOrg)
	require.NoError(err, "could not create organization")
	metadata.OrgID = org.Id
	metadata.Organizations = []string{newOrg.Id}
	appdata, err = metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	s.auth.SetUserRoles([]string{"TRISA Service Provider"})
	claims.Permissions = []string{auth.SwitchOrganizations, auth.UpdateCollaborators}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	now := time.Now().Format(time.RFC3339Nano)
	collab := &models.Collaborator{
		Email:     claims.Email,
		UserId:    "auth0|5f7b5f1b0b8b9b0069b0b1d5",
		Verified:  true,
		JoinedAt:  now,
		LastLogin: now,
	}
	require.NoError(org.AddCollaborator(collab), "could not add collaborator to organization")
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization")

	// Valid login - appdata should be updated to include the added VASP
	require.NoError(s.client.Login(context.TODO(), nil))
	userMeta := &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.Equal(metadata.OrgID, userMeta.OrgID, "app metadata should contain the organization")
	require.Equal(metadata.VASPs.TestNet, userMeta.VASPs.TestNet, "app metadata should contain the testnet VASP")
	require.Equal(userMeta.VASPs.MainNet, userMeta.VASPs.MainNet, "app metadata should contain the mainnet VASP")
	require.ElementsMatch([]string{org.Id, newOrg.Id}, userMeta.Organizations, "app metadata should contain the organizations")

	// User should exist as a collaborator in the organization
	org, err = s.bff.OrganizationFromID(userMeta.OrgID)
	require.NoError(err, "could not get organization from ID")
	require.Len(org.Collaborators, 1, "organization should have one collaborator")
	collab = org.GetCollaborator(claims.Email)
	require.NotNil(collab, "collaborator should exist in organization")
	require.Equal(claims.Email, collab.Email, "collaborator email should match")
	require.NotEmpty(collab.UserId, "collaborator user id should not be empty")
	require.True(collab.Verified, "collaborator should be verified")

	// Last login timestamp should be later than joined timestamp
	lastLogin, err := time.Parse(time.RFC3339Nano, collab.LastLogin)
	require.NoError(err, "could not parse last login timestamp")
	joinedAt, err := time.Parse(time.RFC3339Nano, collab.JoinedAt)
	require.NoError(err, "could not parse joined at timestamp")
	require.True(lastLogin.After(joinedAt), "last login timestamp should be after joined at timestamp")

	// User role should not change
	require.Equal([]string{"TRISA Service Provider"}, s.auth.GetUserRoles(), "user should have the leader role")

	// User should be able to login to the same organization if requested
	params := &api.LoginParams{
		OrgID: org.Id,
	}
	require.NoError(s.client.Login(context.TODO(), params))

	// User app metadata should not change
	nextMeta := &auth.AppMetadata{}
	require.NoError(nextMeta.Load(s.auth.GetUserAppMetadata()))
	require.Equal(userMeta, nextMeta, "app metadata should not change")

	// TODO: The endpoint should return a 204 without the refresh_token flag
}

func (s *bffTestSuite) TestUserInviteLogin() {
	require := s.Require()
	defer s.ResetDB()
	defer s.auth.ResetUserAppMetadata()

	// Test the invited user case - orgID in the params but not in the app metadata
	claims := &authtest.Claims{
		Email: "leopold.wentzel@gmail.com",
	}
	require.NoError(s.SetClientCredentials(claims))
	params := &api.LoginParams{
		OrgID: "67428be4-3fa4-4bf2-9e15-edbf043f8670",
	}
	s.auth.SetUserRoles([]string{auth.CollaboratorRole})

	// Return an error if the organization does not exist
	err := s.client.Login(context.TODO(), params)
	s.requireError(err, http.StatusNotFound, "organization not found")

	// Create the organization in the database
	org := &models.Organization{
		Id: params.OrgID,
		Testnet: &models.DirectoryRecord{
			Id: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
		},
		Mainnet: &models.DirectoryRecord{
			Id: "87d92fd1-53cf-47d8-85b1-048e8a38ced9",
		},
	}
	_, err = s.DB().CreateOrganization(org)
	require.NoError(err, "could not create organization")

	// Return an error if the user is not a collaborator in the organization
	// Note: This is a critical test case because it ensures that a user cannot login
	// to any organization by simply passing the orgID in the params. Instead, they
	// must already be listed as a collaborator in the organization in the database.
	err = s.client.Login(context.TODO(), params)
	s.requireError(err, http.StatusUnauthorized, "user is not authorized to access this organization")

	// Add a pre-existing collaborator who sent the invite
	leader := &models.Collaborator{
		Email:    "orgleader@example.com",
		UserId:   "auth0|6f7b5f1b123b9b0069b0fab3",
		Verified: true,
	}
	require.NoError(org.AddCollaborator(leader), "could not add collaborator to organization")

	// Add the user as a collaborator in the organization
	collab := &models.Collaborator{
		Email:     claims.Email,
		UserId:    "auth0|5f7b5f1b0b8b9b0069b0b1d5",
		Verified:  true,
		ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339Nano),
	}
	require.NoError(org.AddCollaborator(collab), "could not add collaborator to organization")
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization")

	// User should not be able to access the organization if the invitation has expired
	err = s.client.Login(context.TODO(), params)
	s.requireError(err, http.StatusForbidden, "user invitation has expired")

	// Configure a valid invitation
	collab.ExpiresAt = time.Now().Add(time.Hour).Format(time.RFC3339Nano)
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization")

	// Valid login - appdata should be updated with the organization
	require.NoError(s.client.Login(context.TODO(), params))
	userMeta := &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.Equal(org.Id, userMeta.OrgID, "app metadata should contain the organization")
	require.Equal(org.Testnet.Id, userMeta.VASPs.TestNet, "app metadata should contain the testnet VASP")
	require.Equal(org.Mainnet.Id, userMeta.VASPs.MainNet, "app metadata should contain the mainnet VASP")
	require.Empty(userMeta.Organizations, "app metadata organization list should be empty")

	// User should have the same role
	require.Equal([]string{auth.CollaboratorRole}, s.auth.GetUserRoles(), "user should have the collaborator role")

	// Collaborator should contain updated timestamps
	org, err = s.bff.OrganizationFromID(org.Id)
	require.NoError(err, "could not get organization from ID")
	collab = org.GetCollaborator(claims.Email)
	require.NotNil(collab, "collaborator should exist in organization")
	require.NotEmpty(collab.JoinedAt, "collaborator last login timestamp should not be empty")
	require.Equal(collab.JoinedAt, collab.LastLogin, "collaborator joined at timestamp should not be empty")

	// Create a new organization in the database
	newOrg := &models.Organization{
		Testnet: org.Mainnet,
		Mainnet: org.Testnet,
	}
	_, err = s.DB().CreateOrganization(newOrg)
	require.NoError(err, "could not create organization")

	// Add the collaborators to the new organization
	collab.JoinedAt = ""
	collab.LastLogin = ""
	require.NoError(newOrg.AddCollaborator(collab), "could not add collaborator to organization")
	require.NoError(newOrg.AddCollaborator(leader), "could not add collaborator to organization")
	require.NoError(s.DB().UpdateOrganization(newOrg), "could not update organization")

	// Valid login - collab abandons the old organization and joins the new one
	params.OrgID = newOrg.Id
	s.auth.SetUserRoles([]string{auth.CollaboratorRole})
	require.NoError(s.client.Login(context.TODO(), params))
	userMeta = &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.Equal(newOrg.Id, userMeta.OrgID, "app metadata should contain the new organization")
	require.Equal(newOrg.Testnet.Id, userMeta.VASPs.TestNet, "app metadata should contain the testnet VASP")
	require.Equal(newOrg.Mainnet.Id, userMeta.VASPs.MainNet, "app metadata should contain the mainnet VASP")
	require.Empty(userMeta.Organizations, "app metadata organization list should be empty")

	// User should have the same role
	require.Equal([]string{auth.CollaboratorRole}, s.auth.GetUserRoles(), "user should have the collaborator role")

	// Collaborator should contain updated timestamps
	newOrg, err = s.bff.OrganizationFromID(newOrg.Id)
	require.NoError(err, "could not get organization from ID")
	collab = newOrg.GetCollaborator(claims.Email)
	require.NotNil(collab, "collaborator should exist in organization")
	require.NotEmpty(collab.JoinedAt, "collaborator last login timestamp should not be empty")
	require.Equal(collab.JoinedAt, collab.LastLogin, "collaborator joined at timestamp should not be empty")

	// Valid login - leader abandons the old organization and joins the new one
	// TODO: Currently the Auth0 mock only supports one user at a time so these helpers
	// are a workaround to make sure we get the correct user data at test time
	s.auth.SetUserEmail(leader.Email)
	s.auth.SetUserRoles([]string{auth.LeaderRole})
	claims.Permissions = []string{auth.UpdateCollaborators}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	userMeta.OrgID = org.Id
	appdata, err := userMeta.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	require.NoError(s.client.Login(context.TODO(), params))
	userMeta = &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.Equal(newOrg.Id, userMeta.OrgID, "app metadata should contain the new organization")
	require.Equal(newOrg.Testnet.Id, userMeta.VASPs.TestNet, "app metadata should contain the testnet VASP")
	require.Equal(newOrg.Mainnet.Id, userMeta.VASPs.MainNet, "app metadata should contain the mainnet VASP")
	require.Empty(userMeta.Organizations, "app metadata organization list should be empty")

	// Leader user should be demoted to a collaborator in the new organization
	require.Equal([]string{auth.CollaboratorRole}, s.auth.GetUserRoles(), "user should have the collaborator role")

	// Previous organization should be deleted since it has no collaborators
	_, err = s.bff.OrganizationFromID(org.Id)
	require.Error(err, "organization should be deleted")

	// Leader's collab record should contain updated timestamps
	newOrg, err = s.bff.OrganizationFromID(newOrg.Id)
	require.NoError(err, "could not get organization from ID")
	leader = newOrg.GetCollaborator(leader.Email)
	require.NotNil(leader, "leader should exist in organization")
	require.NotEmpty(leader.JoinedAt, "leader last login timestamp should not be empty")
	require.Equal(leader.JoinedAt, leader.LastLogin, "leader joined at timestamp should not be empty")

	// Add the user as a TSP collaborator in a few organizations
	org = &models.Organization{
		Testnet: &models.DirectoryRecord{
			Id: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
		},
		Mainnet: &models.DirectoryRecord{
			Id: "87d92fd1-53cf-47d8-85b1-048e8a38ced9",
		},
	}
	_, err = s.DB().CreateOrganization(org)
	require.NoError(err, "could not create organization")
	require.NoError(org.AddCollaborator(collab), "could not add collaborator to organization")
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization")
	userMeta.Organizations = []string{org.Id}
	appdata, err = userMeta.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	s.auth.SetUserRoles([]string{"TRISA Service Provider"})
	claims.Permissions = []string{auth.SwitchOrganizations, auth.UpdateCollaborators}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")

	// Valid TSP user login - appdata should be updated with the organization list
	params.OrgID = newOrg.Id
	require.NoError(s.client.Login(context.TODO(), params))
	userMeta = &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.Equal(newOrg.Id, userMeta.OrgID, "app metadata should contain the organization")
	require.Equal(newOrg.Testnet.Id, userMeta.VASPs.TestNet, "app metadata should contain the testnet VASP")
	require.Equal(newOrg.Mainnet.Id, userMeta.VASPs.MainNet, "app metadata should contain the mainnet VASP")
	require.ElementsMatch([]string{org.Id, newOrg.Id}, userMeta.Organizations, "app metadata organization list should contain both organizations")

	// User should have the same role
	require.Equal([]string{"TRISA Service Provider"}, s.auth.GetUserRoles(), "user should have the TSP role")

	// Collaborator record should contain updated timestamps
	newOrg, err = s.bff.OrganizationFromID(newOrg.Id)
	require.NoError(err, "could not get organization from ID")
	collab = newOrg.GetCollaborator(claims.Email)
	require.NotNil(collab, "collaborator should exist in organization")
	require.NotEmpty(collab.JoinedAt, "collaborator last login timestamp should not be empty")
	require.Equal(collab.JoinedAt, collab.LastLogin, "collaborator joined at timestamp should not be empty")
}

func (s *bffTestSuite) TestListUserRoles() {
	require := s.Require()

	// Test listing the assignable roles
	expected := []string{
		auth.CollaboratorRole,
		auth.LeaderRole,
	}
	roles, err := s.client.ListUserRoles(context.TODO())
	require.NoError(err, "could not list assignable roles")
	require.Equal(expected, roles, "roles do not match")
}

func (s *bffTestSuite) TestUserOrganization() {
	require := s.Require()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := s.client.UserOrganization(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:organizations permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.UserOrganization(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{"read:organizations"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")
	_, err = s.client.UserOrganization(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Valid claims but no record in the database - should not panic but should return
	// an error
	claims.OrgID = "67428be4-3fa4-4bf2-9e15-edbf043f8670"
	require.NoError(s.SetClientCredentials(claims), "could not create token with correct permissions")
	_, err = s.client.UserOrganization(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when user claims are valid but the organization is not in the database")

	// Successful request returns organization info
	org := &models.Organization{
		Id:     claims.OrgID,
		Name:   "Alice VASP",
		Domain: "alice.io",
	}
	_, err = s.DB().CreateOrganization(org)
	require.NoError(err, "could not create organization in the database")
	expected := &api.OrganizationReply{
		ID:        org.Id,
		Name:      org.Name,
		Domain:    org.Domain,
		CreatedAt: org.Created,
	}
	reply, err := s.client.UserOrganization(context.TODO())
	require.NoError(err, "could not get user organization info")
	require.Equal(expected, reply, "organization info does not match")
}
