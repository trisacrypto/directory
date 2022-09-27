package bff_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
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
		s.db.DeleteOrganization(org.UUID())
	}()

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
	org, err = s.db.RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 1, "expected one collaborator in the organization")
	collab, ok := org.Collaborators[request.Email]
	require.True(ok, "expected collaborator to be in the organization")
	require.Equal(request.Email, collab.Email, "expected collaborator email to match")
	require.NotEmpty(collab.CreatedAt, "expected collaborator to have a created at timestamp")
	require.Empty(collab.VerifiedAt, "expected collaborator to not have a verified at timestamp")

	// Should return an error if the collaborator already exists
	_, err = s.client.AddCollaborator(context.TODO(), request)
	s.requireError(err, http.StatusBadRequest, fmt.Sprintf("collaborator %q already exists", request.Email), "expected error when collaborator already exists")
}

func (s *bffTestSuite) TestReplaceCollaborator() {
	require := s.Require()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
	_, err := s.client.ReplaceCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:collaborator permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.ReplaceCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{"update:collaborators"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	_, err = s.client.ReplaceCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic but should
	// return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	_, err = s.client.ReplaceCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when user claims are valid but the organization is not in the database")

	// Create an organization in the database without any collaborators
	org, err := s.db.CreateOrganization()
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.DeleteOrganization(org.UUID())
	}()

	// Create valid credentials with the organization ID
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")

	// Should return an error if the collaborator email is missing from the request
	_, err = s.client.ReplaceCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusBadRequest, "collaborator is missing email address", "expected error when collaborator email is missing")

	// Should return an error if the collaborator does not exist
	request := &models.Collaborator{
		Email: "alice@example.com",
	}
	_, err = s.client.ReplaceCollaborator(context.TODO(), request)
	s.requireError(err, http.StatusBadRequest, fmt.Sprintf("collaborator %q does not exist", request.Email), "expected error when collaborator does not exist")

	// Add a new collaborator to the organization
	collab, err := s.client.AddCollaborator(context.TODO(), request)
	require.NoError(err, "could not add collaborator to organization")

	// Edit collaborator data
	collab.VerifiedAt = time.Now().Format(time.RFC3339Nano)
	modified, err := s.client.ReplaceCollaborator(context.TODO(), collab)
	require.NoError(err, "could not replace collaborator in organization")
	require.Equal(collab.Email, modified.Email, "expected email to match original email")
	require.Equal(collab.CreatedAt, modified.CreatedAt, "expected created at timestamp to match original timestamp")
	require.Equal(collab.VerifiedAt, modified.VerifiedAt, "expected verified at timestamp to match original timestamp")

	// Updated collaborator should be in the database
	org, err = s.db.RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 1, "expected one collaborator in the organization")
	retrieved, ok := org.Collaborators[collab.Email]
	require.True(ok, "expected collaborator to be in the organization")
	require.Equal(collab.Email, retrieved.Email, "expected email to match original email")
	require.Equal(collab.CreatedAt, retrieved.CreatedAt, "expected created at timestamp to match original timestamp")
	require.Equal(collab.VerifiedAt, retrieved.VerifiedAt, "expected verified at timestamp to match original timestamp")
}

func (s *bffTestSuite) TestDeleteCollaborator() {
	require := s.Require()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")
	err := s.client.DeleteCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the update:collaborators permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	err = s.client.DeleteCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Claims must have an organization ID
	claims.Permissions = []string{"update:collaborators"}
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	err = s.client.DeleteCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "missing claims info, try logging out and logging back in", "expected error when user claims does not have an orgid")

	// Create valid claims but no record in the database - should not panic but should
	// return an error
	claims.OrgID = "2295c698-afdc-4aaf-9443-85a4515217e3"
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")
	err = s.client.DeleteCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusUnauthorized, "no organization found, try logging out and logging back in", "expected error when user claims are valid but the organization is not in the database")

	// Create an organization in the database without any collaborators
	org, err := s.db.CreateOrganization()
	require.NoError(err, "could not create organization in the database")
	defer func() {
		// Ensure organization is deleted at the end of the tests
		s.db.DeleteOrganization(org.UUID())
	}()

	// Create valid credentials with the organization ID
	claims.OrgID = org.Id
	require.NoError(s.SetClientCredentials(claims), "could not create token with valid credentials")

	// Should return an error if the collaborator email is missing from the request
	err = s.client.DeleteCollaborator(context.TODO(), &models.Collaborator{})
	s.requireError(err, http.StatusBadRequest, "invalid collaborator in request", "expected error when collaborator email is missing")

	// Should not return an error if the collaborator does not exist
	request := &models.Collaborator{
		Email: "alice@example.com",
	}
	err = s.client.DeleteCollaborator(context.TODO(), request)
	require.NoError(err, "expected no error when collaborator does not exist")

	// Add a new collaborator to the organization
	collab, err := s.client.AddCollaborator(context.TODO(), request)
	require.NoError(err, "could not add collaborator to organization")

	// Delete collaborator from the organization
	err = s.client.DeleteCollaborator(context.TODO(), collab)
	require.NoError(err, "could not delete collaborator from organization")

	// Deleted collaborator should not be in the database
	org, err = s.db.RetrieveOrganization(org.UUID())
	require.NoError(err, "could not retrieve organization from the database")
	require.Len(org.Collaborators, 0, "expected no collaborators in the organization")
}
