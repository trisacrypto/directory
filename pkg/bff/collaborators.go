package bff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

// AddCollaborator creates a new collaborator with the email address in the request.
// The endpoint adds the collaborator to the organization record associated with the
// user and sends a verification email to the provided email address, so the user must
// have the update:collaborators permission.
func (s *Server) AddCollaborator(c *gin.Context) {
	var (
		err          error
		collaborator *models.Collaborator
		org          *models.Organization
	)

	// Fetch the organization from the claims
	// NOTE: This method handles the error logging and response
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Unmarshal the collaborator from the POST request
	collaborator = &models.Collaborator{}
	if err = c.ShouldBind(collaborator); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Add the collaborator to the organization
	if err = org.AddCollaborator(collaborator); err != nil {
		log.Error().Err(err).Msg("could not add new collaborator to organization")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// TODO: Send invite/verification email to the collaborator

	// Save the updated organization
	if err = s.db.UpdateOrganization(org); err != nil {
		log.Error().Err(err).Msg("could not save organization with new collaborator")
		c.JSON(http.StatusInternalServerError, "could not add collaborator")
		return
	}

	c.JSON(http.StatusOK, collaborator)
}

// ReplaceCollaborator completely replaces a collaborator on the user's organization
// with the collaborator in the request. The collaborator object in the request must be
// valid and the user must have the update:collaborators permission.
func (s *Server) ReplaceCollaborator(c *gin.Context) {
	var (
		err          error
		collaborator *models.Collaborator
		org          *models.Organization
	)

	// Fetch the organization from the claims
	// NOTE: This method handles the error logging and response
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Unmarshal the collaborator from the PUT request
	collaborator = &models.Collaborator{}
	if err = c.ShouldBind(collaborator); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Replace the collaborator on the organization
	if err = org.ReplaceCollaborator(collaborator); err != nil {
		log.Error().Err(err).Msg("could not replace collaborator on organization")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Save the updated organization
	if err = s.db.UpdateOrganization(org); err != nil {
		log.Error().Err(err).Msg("could not save organization with new collaborator")
		c.JSON(http.StatusInternalServerError, "could not replace collaborator")
		return
	}

	c.JSON(http.StatusOK, collaborator)
}
