package bff_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

func (s *bffTestSuite) TestAddCollaborator() {
	require := s.Require()

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
	org, err := s.db.CreateOrganization()
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.DeleteOrganization(org.Id)
	}()

	// Create valid credentials with the organization ID
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")

	// Should return an error if the collaborator email is missing from the request
	_, err = s.client.AddCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusBadRequest, "email address is required to add an organization collaborator", "expected error when collaborator email is missing")

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
	org, err = s.db.RetrieveOrganization(org.Id)
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 1, "expected one collaborator in the organization")
	collab, ok := org.Collaborators[request.Email]
	require.True(ok, "expected collaborator to be in the organization")
	require.Equal(request.Email, collab.Email, "expected collaborator email to match")
	require.NotEmpty(collab.CreatedAt, "expected collaborator to have a created at timestamp")
	require.Empty(collab.VerifiedAt, "expected collaborator to not have a verified at timestamp")

	// Should return an error if the collaborator already exists
	_, err = s.client.AddCollaborator(context.TODO(), request)
	s.requireError(err, http.StatusBadRequest, fmt.Sprintf("collaborator with email address %s already exists", request.Email), "expected error when collaborator already exists")
}
