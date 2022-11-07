package bff_test

import (
	"context"
	"net/http"

	"github.com/trisacrypto/directory/pkg/bff"
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
	require.NoError(s.client.Login(context.TODO(), &api.LoginParams{}))

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
	require.Len(org.Collaborators, 1, "organization should have one collaborator")
	collab := org.GetCollaborator(claims.Email)
	require.NotNil(collab, "collaborator should exist in organization")
	require.Equal(claims.Email, collab.Email, "collaborator email should match")
	require.NotEmpty(collab.UserId, "collaborator user id should not be empty")
	require.True(collab.Verified, "collaborator should be verified")

	// User should assume the leader role
	require.Equal([]string{bff.LeaderRole}, s.auth.GetUserRoles(), "user should have the leader role")

	// New TSP users should have the new organization added to the app metadata
	appdata, err = metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	s.auth.SetUserRoles([]string{bff.TSPRole})
	require.NoError(s.client.Login(context.TODO(), &api.LoginParams{}))

	// Appdata should contain a new organization and it should be added to the list
	userMeta = &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.NotEmpty(userMeta.OrgID, "app metadata should contain a new organization")
	require.Empty(userMeta.VASPs.TestNet, "app metadata should not contain a testnet VASP")
	require.Empty(userMeta.VASPs.MainNet, "app metadata should not contain a mainnet VASP")
	require.Len(userMeta.Organizations, 1, "app metadata organization list should contain one organization")
	require.Equal(map[string]struct{}{userMeta.OrgID: {}}, userMeta.Organizations, "app metadata organization list should contain the new organization")

	// User should exist as a collaborator in the new organization
	org, err = s.bff.OrganizationFromID(userMeta.OrgID)
	require.NoError(err, "could not get organization from ID")
	require.Len(org.Collaborators, 1, "organization should have one collaborator")
	collab = org.GetCollaborator(claims.Email)
	require.NotNil(collab, "collaborator should exist in organization")
	require.Equal(claims.Email, collab.Email, "collaborator email should match")
	require.NotEmpty(collab.UserId, "collaborator user id should not be empty")
	require.True(collab.Verified, "collaborator should be verified")

	// User should still have the TSP role
	require.Equal([]string{bff.TSPRole}, s.auth.GetUserRoles(), "user should have the TSP role")

	// TODO: The endpoint should return a 200 with the refresh_token flag set
}

func (s *bffTestSuite) TestReturningUserLogin() {
	require := s.Require()
	defer s.ResetDB()
	defer s.auth.ResetUserAppMetadata()

	// Test the existing user case - orgID in app metadata but not params
	claims := &authtest.Claims{
		Email: "leopold.wentzel@gmail.com",
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
	s.auth.SetUserRoles([]string{bff.LeaderRole})

	// Returns an error if the organization does not exist
	err = s.client.Login(context.TODO(), &api.LoginParams{})
	s.requireError(err, http.StatusNotFound, "organization not found")

	// Create the organization in the database without the collaborator
	org, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization")
	org.Id = metadata.OrgID
	org.Testnet = &models.DirectoryRecord{
		Id: metadata.VASPs.TestNet,
	}
	org.Mainnet = &models.DirectoryRecord{
		Id: metadata.VASPs.MainNet,
	}
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization")

	// User is not authorized to access the organization without being a collaborator
	err = s.client.Login(context.TODO(), &api.LoginParams{})
	s.requireError(err, http.StatusUnauthorized, "user is not authorized to access this organization")

	// Make the user a TSP
	newOrg, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization")
	metadata.OrgID = org.Id
	metadata.Organizations = map[string]struct{}{
		newOrg.Id: {},
	}
	appdata, err = metadata.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	s.auth.SetUserRoles([]string{bff.TSPRole})
	collab := &models.Collaborator{
		Email:    claims.Email,
		UserId:   "auth0|5f7b5f1b0b8b9b0069b0b1d5",
		Verified: true,
	}
	require.NoError(org.AddCollaborator(collab), "could not add collaborator to organization")
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization")

	// Valid login - appdata should be updated to include the added VASP
	require.NoError(s.client.Login(context.TODO(), &api.LoginParams{}))
	userMeta := &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.Equal(metadata.OrgID, userMeta.OrgID, "app metadata should contain the organization")
	require.Equal(metadata.VASPs.TestNet, userMeta.VASPs.TestNet, "app metadata should contain the testnet VASP")
	require.Equal(userMeta.VASPs.MainNet, userMeta.VASPs.MainNet, "app metadata should contain the mainnet VASP")
	require.Equal(map[string]struct{}{org.Id: {}, newOrg.Id: {}}, userMeta.Organizations, "app metadata should contain the organizations")

	// User should exist as a collaborator in the organization
	org, err = s.bff.OrganizationFromID(userMeta.OrgID)
	require.NoError(err, "could not get organization from ID")
	require.Len(org.Collaborators, 1, "organization should have one collaborator")
	collab = org.GetCollaborator(claims.Email)
	require.NotNil(collab, "collaborator should exist in organization")
	require.Equal(claims.Email, collab.Email, "collaborator email should match")
	require.NotEmpty(collab.UserId, "collaborator user id should not be empty")
	require.True(collab.Verified, "collaborator should be verified")

	// User role should not change
	require.Equal([]string{bff.TSPRole}, s.auth.GetUserRoles(), "user should have the leader role")

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
	s.auth.SetUserRoles([]string{bff.CollaboratorRole})

	// Return an error if the organization does not exist
	err := s.client.Login(context.TODO(), params)
	s.requireError(err, http.StatusNotFound, "organization not found")

	// Create the organization in the database
	org, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization")
	org.Id = params.OrgID
	org.Testnet = &models.DirectoryRecord{
		Id: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
	}
	org.Mainnet = &models.DirectoryRecord{
		Id: "87d92fd1-53cf-47d8-85b1-048e8a38ced9",
	}
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization")

	// Return an error if the user is not a collaborator in the organization
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
		Email:    claims.Email,
		UserId:   "auth0|5f7b5f1b0b8b9b0069b0b1d5",
		Verified: true,
	}
	require.NoError(org.AddCollaborator(collab), "could not add collaborator to organization")
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
	require.Equal([]string{bff.CollaboratorRole}, s.auth.GetUserRoles(), "user should have the collaborator role")

	// Create a new organization in the database
	newOrg, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization")
	newOrg.Testnet = org.Mainnet
	newOrg.Mainnet = org.Testnet

	// Add the collaborators to the new organization
	require.NoError(newOrg.AddCollaborator(collab), "could not add collaborator to organization")
	require.NoError(newOrg.AddCollaborator(leader), "could not add collaborator to organization")
	require.NoError(s.DB().UpdateOrganization(newOrg), "could not update organization")

	// Valid login - collab abandons the old organization and joins the new one
	params.OrgID = newOrg.Id
	s.auth.SetUserRoles([]string{bff.CollaboratorRole})
	require.NoError(s.client.Login(context.TODO(), params))
	userMeta = &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.Equal(newOrg.Id, userMeta.OrgID, "app metadata should contain the new organization")
	require.Equal(newOrg.Testnet.Id, userMeta.VASPs.TestNet, "app metadata should contain the testnet VASP")
	require.Equal(newOrg.Mainnet.Id, userMeta.VASPs.MainNet, "app metadata should contain the mainnet VASP")
	require.Empty(userMeta.Organizations, "app metadata organization list should be empty")

	// User should have the same role
	require.Equal([]string{bff.CollaboratorRole}, s.auth.GetUserRoles(), "user should have the collaborator role")

	// Valid login - leader abandons the old organization and joins the new one
	// TODO: Currently the Auth0 mock only supports one user at a time so these helpers
	// are a workaround to make sure we get the correct user data at test time
	s.auth.SetUserEmail(leader.Email)
	s.auth.SetUserRoles([]string{bff.LeaderRole})
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
	require.Equal([]string{bff.CollaboratorRole}, s.auth.GetUserRoles(), "user should have the collaborator role")

	// Previous organization should be deleted since it has no collaborators
	_, err = s.bff.OrganizationFromID(org.Id)
	require.Error(err, "organization should be deleted")

	// Add the user as a TSP collaborator in a few organizations
	org, err = s.DB().CreateOrganization()
	require.NoError(err, "could not create organization")
	org.Testnet = &models.DirectoryRecord{
		Id: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
	}
	org.Mainnet = &models.DirectoryRecord{
		Id: "87d92fd1-53cf-47d8-85b1-048e8a38ced9",
	}
	require.NoError(org.AddCollaborator(collab), "could not add collaborator to organization")
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization")
	userMeta.Organizations = map[string]struct{}{
		org.Id: {},
	}
	appdata, err = userMeta.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)
	s.auth.SetUserRoles([]string{bff.TSPRole})

	// Valid TSP user login - appdata should be updated with the organization list
	params.OrgID = newOrg.Id
	require.NoError(s.client.Login(context.TODO(), params))
	userMeta = &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.Equal(newOrg.Id, userMeta.OrgID, "app metadata should contain the organization")
	require.Equal(newOrg.Testnet.Id, userMeta.VASPs.TestNet, "app metadata should contain the testnet VASP")
	require.Equal(newOrg.Mainnet.Id, userMeta.VASPs.MainNet, "app metadata should contain the mainnet VASP")
	require.Equal(map[string]struct{}{org.Id: {}, newOrg.Id: {}}, userMeta.Organizations, "app metadata organization list should contain both organizations")

	// User should have the same role
	require.Equal([]string{bff.TSPRole}, s.auth.GetUserRoles(), "user should have the TSP role")
}
