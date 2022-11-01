package bff_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"google.golang.org/protobuf/proto"
)

func (s *bffTestSuite) TestAddCollaborator() {
	require := s.Require()
	defer s.ResetDB()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
	_, err := s.client.AddCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:collaborator permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.AddCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{"update:collaborators"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	_, err = s.client.AddCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic but should
	// return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	_, err = s.client.AddCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when user claims are valid but the organization is not in the database")

	// Create an organization in the database without any collaborators
	org, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization in the database")

	// Create valid credentials with the organization ID
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")

	// Should return an error if the collaborator email is missing from the request
	_, err = s.client.AddCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusBadRequest, "collaborator is missing email address", "expected error when collaborator email is missing")

	// Successfully adding a collaborator to the organization
	request := &models.Collaborator{
		Email: "alice@example.com",
	}
	collab, err := s.client.AddCollaborator(context.TODO(), request)
	require.NoError(err, "could not add collaborator to organization")
	require.Equal(request.Email, collab.Email, "expected collaborator email to match request email")
	require.NotEmpty(collab.CreatedAt, "expected collaborator to have a created at timestamp")
	require.Empty(collab.VerifiedAt, "expected collaborator to not have a verified at timestamp")

	// Collaborator should be in the database
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 1, "expected one collaborator in the organization")
	collab, ok := org.Collaborators[request.Key()]
	require.True(ok, "expected collaborator to be in the organization")
	require.Equal(request.Email, collab.Email, "expected collaborator email to match")
	require.NotEmpty(collab.CreatedAt, "expected collaborator to have a created at timestamp")
	require.Empty(collab.VerifiedAt, "expected collaborator to not have a verified at timestamp")

	// Should return an error if the collaborator already exists
	_, err = s.client.AddCollaborator(context.TODO(), request)
	s.requireError(err, http.StatusConflict, "collaborator already exists", "expected error when collaborator already exists")
}

func (s *bffTestSuite) TestListCollaborators() {
	require := s.Require()
	defer s.ResetDB()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
	_, err := s.client.ListCollaborators(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:collaborators permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.ListCollaborators(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{"read:collaborators"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	_, err = s.client.ListCollaborators(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic but should
	// return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	_, err = s.client.ListCollaborators(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when user claims are valid but the organization is not in the database")

	// Create an organization in the database without any collaborators
	org, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization in the database")

	// Create valid credentials with the organization ID
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")

	// Should return an empty list if there are no collaborators in the organization
	collabs, err := s.client.ListCollaborators(context.TODO())
	require.NoError(err, "could not list collaborators")
	require.Len(collabs.Collaborators, 0, "expected empty collaborators list to be returned")

	// Add a new collaborator to the organization
	leopold := &models.Collaborator{
		Email: "leopold.wentzel@gmail.com",
		UserId: authtest.UserID,
		VerifiedAt: time.Now().Format(time.RFC3339Nano),
	}
	leopoldRoles := []string{authtest.UserRole}
	org.Collaborators = make(map[string]*models.Collaborator)
	org.Collaborators[leopold.Key()] = leopold
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization in the database")

	// Should return a list with one collaborator in it
	collabs, err = s.client.ListCollaborators(context.TODO())
	require.NoError(err, "could not list collaborators")
	require.Len(collabs.Collaborators, 1, "expected one collaborator in the list")
	require.Equal(leopold.Email, collabs.Collaborators[0].Email, "expected collaborator email to match")

	// Add some more collaborators to the organization
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	bob := &models.Collaborator{
		Email: "bob@example.com",
	}
	org.Collaborators[bob.Key()] = bob

	charlie := &models.Collaborator{
		Email: "charlie@example.com",
	}
	org.Collaborators[charlie.Key()] = charlie

	yogg := &models.Collaborator{
		Email: "yogg-sothoth@hpl.org",
	}
	org.Collaborators[yogg.Key()] = yogg
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization in the database")

	// Should return a list with collaborators ordered by email address
	collabs, err = s.client.ListCollaborators(context.TODO())
	require.NoError(err, "could not list collaborators")
	require.True(proto.Equal(bob, collabs.Collaborators[0]), "expected bob to be first in the list")
	require.True(proto.Equal(charlie, collabs.Collaborators[1]), "expected charlie to be second in the list")
	require.Equal(leopold.Email, collabs.Collaborators[2].Email, "expected leopold to be third in the list")
	require.Equal(authtest.Name, collabs.Collaborators[2].Name, "expected leopold name to be set")
	require.Equal(leopoldRoles, collabs.Collaborators[2].Roles, "expected leopold roles to be set")
	require.True(proto.Equal(yogg, collabs.Collaborators[3]), "expected yogg to be fourth in the list")
}

func (s *bffTestSuite) TestUpdateCollaboratorRoles() {
	require := s.Require()
	defer s.ResetDB()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
	_, err := s.client.UpdateCollaboratorRoles(context.TODO(), "invalid", &api.UpdateRolesParams{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:collaborator permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.UpdateCollaboratorRoles(context.TODO(), "invalid", &api.UpdateRolesParams{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{"update:collaborators"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	_, err = s.client.UpdateCollaboratorRoles(context.TODO(), "invalid", &api.UpdateRolesParams{})
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic but should
	// return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	_, err = s.client.UpdateCollaboratorRoles(context.TODO(), "invalid", &api.UpdateRolesParams{})
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when user claims are valid but the organization is not in the database")

	// Create an organization in the database without any collaborators
	org, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization in the database")

	// Create valid credentials with the organization ID
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")

	// Should return an error if the collaborator does not exist
	_, err = s.client.UpdateCollaboratorRoles(context.TODO(), "invalid", &api.UpdateRolesParams{})
	s.requireError(err, http.StatusNotFound, "collaborator does not exist", "expected error when collaborator does not exist")

	// Add a new collaborator to the organization
	request := &models.Collaborator{
		Email: "alice@example.com",
	}
	collab, err := s.client.AddCollaborator(context.TODO(), request)
	require.NoError(err, "could not add collaborator to organization")

	// Should return an error if the collaborator is not verified
	params := &api.UpdateRolesParams{
		Roles: []string{"notarole"},
	}
	_, err = s.client.UpdateCollaboratorRoles(context.TODO(), collab.Id, params)
	s.requireError(err, http.StatusBadRequest, "cannot update roles for unverified collaborator", "expected error when collaborator is not verified")

	// Create a verified collaborator in the database
	collab.VerifiedAt = time.Now().Format(time.RFC3339Nano)
	collab.UserId = authtest.UserID
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	org.Collaborators[collab.Key()] = collab
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization in the database")

	// Should return an error if there is an invalid role
	_, err = s.client.UpdateCollaboratorRoles(context.TODO(), collab.Id, params)
	s.requireError(err, http.StatusBadRequest, fmt.Sprintf("could not find role %q in 3 available roles", params.Roles[0]), "expected error when role is invalid")

	// Successfully updating the roles of a collaborator
	params.Roles = []string{"Organization Collaborator", "Organization Leader"}
	modified, err := s.client.UpdateCollaboratorRoles(context.TODO(), collab.Id, params)
	require.NoError(err, "could not update collaborator roles")
	require.Equal(collab.Id, modified.Id, "expected collaborator ID to match")

	// Updated collaborator should be in the database
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 1, "expected one collaborator in the organization")
	_, ok := org.Collaborators[collab.Key()]
	require.True(ok, "expected collaborator to be in the organization")

	// TODO: It's difficult to verify correctness without more deeply mocking the Auth0
	// role management.
}

func (s *bffTestSuite) TestDeleteCollaborator() {
	require := s.Require()
	defer s.ResetDB()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
	err := s.client.DeleteCollaborator(context.TODO(), "invalid")
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:collaborators permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	err = s.client.DeleteCollaborator(context.TODO(), "invalid")
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{"update:collaborators"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	err = s.client.DeleteCollaborator(context.TODO(), "invalid")
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic but should
	// return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	err = s.client.DeleteCollaborator(context.TODO(), "invalid")
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when user claims are valid but the organization is not in the database")

	// Create an organization in the database without any collaborators
	org, err := s.DB().CreateOrganization()
	require.NoError(err, "could not create organization in the database")

	// Create valid credentials with the organization ID
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")

	// Should return an error if the collaborator does not exist
	err = s.client.DeleteCollaborator(context.TODO(), "c160f8cc69a4f0bf2b0362752353d060")
	s.requireError(err, http.StatusNotFound, "collaborator not found", "expected an error when collaborator does not exist")

	// Add a new collaborator to the organization
	collab := &models.Collaborator{
		Email: "alice@example.com",
	}
	collab, err = s.client.AddCollaborator(context.TODO(), collab)
	require.NoError(err, "could not add collaborator to organization")

	// If the collaborator is not verified, then the record is just deleted from the
	// organization
	require.NoError(s.client.DeleteCollaborator(context.TODO(), collab.Id))
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 0, "expected no collaborators in the organization")

	// Add the collaborator again
	collab, err = s.client.AddCollaborator(context.TODO(), collab)
	require.NoError(err, "could not add collaborator to organization")

	// Configure a verified collaborator in the database
	collab.VerifiedAt = time.Now().Format(time.RFC3339Nano)
	collab.UserId = authtest.UserID
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	org.Collaborators[collab.Key()] = collab
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization in the database")

	// If a verified collaborator is deleted, then the record should still be deleted
	// from the organization
	require.NoError(s.client.DeleteCollaborator(context.TODO(), collab.Id))
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 0, "expected no collaborators in the organization")

	// The user app metadata should also be updated
	appdata := &auth.AppMetadata{}
	require.NoError(appdata.Load(s.auth.GetUserAppMetadata()))
	require.Empty(appdata.OrgID, "expected orgid in app metadata to be empty")
}
