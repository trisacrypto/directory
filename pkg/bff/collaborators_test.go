package bff_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/emails/mock"
	"google.golang.org/protobuf/proto"
)

func (s *bffTestSuite) TestAddCollaborator() {
	require := s.Require()
	defer s.ResetDB()
	defer mock.PurgeEmails()

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
	claims.Permissions = []string{auth.UpdateCollaborators}
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
	org := &models.Organization{}
	_, err = s.DB().CreateOrganization(org)
	require.NoError(err, "could not create organization in the database")

	// Create valid credentials with the organization ID
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")

	// Should return an error if the collaborator email is missing from the request
	_, err = s.client.AddCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusBadRequest, "collaborator record is invalid", "expected error when collaborator email is missing")

	// Successfully adding a collaborator to the organization
	request := &models.Collaborator{
		Email: "alice@example.com",
	}
	collab, err := s.client.AddCollaborator(context.TODO(), request)
	require.NoError(err, "could not add collaborator to organization")
	require.Equal(request.Email, collab.Email, "expected collaborator email to match request email")
	require.NotEmpty(collab.CreatedAt, "expected collaborator to have a created at timestamp")
	require.False(collab.Verified, "expected collaborator to not be verified")
	require.NotEmpty(collab.ExpiresAt, "expected collaborator to have an expires at timestamp")

	// Collaborator should be in the database
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 1, "expected one collaborator in the organization")
	collab, ok := org.Collaborators[request.Key()]
	require.True(ok, "expected collaborator to be in the organization")
	require.Equal(request.Email, collab.Email, "expected collaborator email to match")
	require.NotEmpty(collab.CreatedAt, "expected collaborator to have a created at timestamp")
	require.False(collab.Verified, "expected collaborator to not be verified")

	// Email should be sent to the collaborator
	require.Len(mock.Emails, 1, "expected one email to be sent")

	// Should return an error if the collaborator already exists
	_, err = s.client.AddCollaborator(context.TODO(), request)
	s.requireError(err, http.StatusConflict, "collaborator already exists in organization", "expected error when collaborator already exists")

	// Test that number of collaborators on an organization is limited
	err = s.client.DeleteCollaborator(context.TODO(), collab.Id)
	require.NoError(err, "could not delete collaborator from organization")
	for i := 0; i < models.MaxCollaborators; i++ {
		request.Email = fmt.Sprintf("alice%d@example.com", i)
		_, err = s.client.AddCollaborator(context.TODO(), request)
		require.NoError(err, "could not add collaborator to organization")
	}
	_, err = s.client.AddCollaborator(context.TODO(), request)
	s.requireError(err, http.StatusForbidden, "maximum number of collaborators reached", "expected error when maximum number of collaborators is reached")
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
	claims.Permissions = []string{auth.ReadCollaborators}
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
	org := &models.Organization{}
	_, err = s.DB().CreateOrganization(org)
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
		Email:    "leopold.wentzel@gmail.com",
		UserId:   authtest.UserID,
		Verified: true,
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
	require.Equal(leopoldRoles, collabs.Collaborators[0].Roles, "expected collaborator roles to match")

	// Collaborator should be updated in the database
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 1, "expected one collaborator in the organization")
	collab, ok := org.Collaborators[leopold.Key()]
	require.True(ok, "expected collaborator to be in the organization")
	require.Equal(leopold.Email, collab.Email, "expected collaborator email in database to match")
	require.Equal(leopoldRoles, collab.Roles, "expected collaborator roles in database to match")

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
	claims.Permissions = []string{auth.UpdateCollaborators}
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
	org := &models.Organization{}
	_, err = s.DB().CreateOrganization(org)
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
	collab.Verified = true
	collab.UserId = authtest.UserID
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	org.Collaborators[collab.Key()] = collab
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization in the database")

	// Should return an error if there is an invalid role
	_, err = s.client.UpdateCollaboratorRoles(context.TODO(), collab.Id, params)
	s.requireError(err, http.StatusBadRequest, bff.ErrInvalidUserRole.Error(), "expected error when role is invalid")

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
	claims.Permissions = []string{auth.UpdateCollaborators}
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
	org := &models.Organization{}
	_, err = s.DB().CreateOrganization(org)
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
	collab.Verified = true
	collab.UserId = authtest.UserID
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	org.Collaborators[collab.Key()] = collab
	require.NoError(s.DB().UpdateOrganization(org), "could not update organization in the database")

	// Make sure the user has some app metadata
	userMeta := &auth.AppMetadata{
		OrgID: org.Id,
		VASPs: auth.VASPs{
			MainNet: "1bcacaf5-4b43-4e14-b70c-a47107d3a56c",
			TestNet: "87d92fd1-53cf-47d8-85b1-048e8a38ced9",
		},
	}
	appdata, err := userMeta.Dump()
	require.NoError(err, "could not dump app metadata")
	s.auth.SetUserAppMetadata(appdata)

	// If a verified collaborator is deleted, then the record should still be deleted
	// from the organization
	require.NoError(s.client.DeleteCollaborator(context.TODO(), collab.Id))
	org, err = s.DB().RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 0, "expected no collaborators in the organization")

	// The org in the user's app metadata should be cleared
	userMeta = &auth.AppMetadata{}
	require.NoError(userMeta.Load(s.auth.GetUserAppMetadata()))
	require.Empty(userMeta.OrgID, "expected orgid in app metadata to be empty")
	require.Empty(userMeta.VASPs.MainNet, "expected mainnet VASP in app metadata to be empty")
	require.Empty(userMeta.VASPs.TestNet, "expected testnet VASP in app metadata to be empty")
}

func (s *bffTestSuite) TestInsortCollaborator() {
	require := s.Require()

	// Should handle nil values without panicking
	require.Nil(bff.InsortCollaborator(nil, nil, nil))
	require.Nil(bff.InsortCollaborator([]*models.Collaborator{}, nil, nil))
	require.Nil(bff.InsortCollaborator([]*models.Collaborator{}, &models.Collaborator{}, nil))

	// Create some collaborators
	alice := &models.Collaborator{
		Email:     "alice@example.com",
		CreatedAt: time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC).Format(time.RFC3339Nano),
	}
	bob := &models.Collaborator{
		Email:     "bob@example.com",
		CreatedAt: time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC).Format(time.RFC3339Nano),
	}
	charlie := &models.Collaborator{
		Email:     "charlie@example.com",
		CreatedAt: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339Nano),
	}

	// Test ordering by email
	f := func(a, b *models.Collaborator) bool {
		return a.Email < b.Email
	}

	// Insort some collaborators into a slice
	collabs := bff.InsortCollaborator([]*models.Collaborator{}, charlie, f)
	require.Len(collabs, 1, "expected one collaborator in the slice")
	require.Equal(charlie.Email, collabs[0].Email, "wrong collaborator in the slice")

	collabs = bff.InsortCollaborator(collabs, alice, f)
	require.Len(collabs, 2, "expected two collaborators in the slice")
	require.Equal(alice.Email, collabs[0].Email, "collaborator not in the correct position")
	require.Equal(charlie.Email, collabs[1].Email, "collaborator not in the correct position")

	collabs = bff.InsortCollaborator(collabs, bob, f)
	require.Len(collabs, 3, "expected three collaborators in the slice")
	require.Equal(alice.Email, collabs[0].Email, "collaborator not in the correct position")
	require.Equal(bob.Email, collabs[1].Email, "collaborator not in the correct position")
	require.Equal(charlie.Email, collabs[2].Email, "collaborator not in the correct position")

	// Test ordering by timestamp
	f = func(a, b *models.Collaborator) bool {
		return a.CreatedAt < b.CreatedAt
	}

	// Insort some collaborators into a slice
	collabs = bff.InsortCollaborator([]*models.Collaborator{}, charlie, f)
	require.Len(collabs, 1, "expected one collaborator in the slice")
	require.Equal(charlie.CreatedAt, collabs[0].CreatedAt, "wrong collaborator in the slice")

	collabs = bff.InsortCollaborator(collabs, alice, f)
	require.Len(collabs, 2, "expected two collaborators in the slice")
	require.Equal(charlie.CreatedAt, collabs[0].CreatedAt, "collaborator not in the correct position")
	require.Equal(alice.CreatedAt, collabs[1].CreatedAt, "collaborator not in the correct position")

	collabs = bff.InsortCollaborator(collabs, bob, f)
	require.Len(collabs, 3, "expected three collaborators in the slice")
	require.Equal(charlie.CreatedAt, collabs[0].CreatedAt, "collaborator not in the correct position")
	require.Equal(bob.CreatedAt, collabs[1].CreatedAt, "collaborator not in the correct position")
	require.Equal(alice.CreatedAt, collabs[2].CreatedAt, "collaborator not in the correct position")
}
