package bff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	records "github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

// AddCollaborator creates a new collaborator with the email address in the request.
// The endpoint adds the collaborator to the organization record associated with the
// user and sends a verification email to the provided email address, so the user must
// have the update:collaborators permission.
func (s *Server) AddCollaborator(c *gin.Context) {
	var (
		err          error
		request      *records.Collaborator
		collaborator *records.Collaborator
		org          *models.Organization
	)

	// Fetch the organization from the claims
	if org, err = s.OrganizationFromClaims(c); err != nil {
		log.Error().Err(err).Msg("could not fetch organization from claims")
		c.JSON(http.StatusInternalServerError, "could not add collaborator")
		return
	}

	// Unmarshal the collaborator from the POST request
	request = &records.Collaborator{}
	if err = c.ShouldBind(request); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Add the collaborator to the organization
	if err = AddOrganizationCollaborator(org, request.Email); err != nil {
		log.Error().Err(err).Msg("could not add new collaborator to organization")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// TODO: Send invite/verification email to the collaborator

	// Save the updated organization
	if err = s.db.Organizations().Update(c.Request.Context(), org); err != nil {
		log.Error().Err(err).Msg("could not save organization with new collaborator")
		c.JSON(http.StatusInternalServerError, "could not add collaborator")
		return
	}

	// Return the updated collaborator
	if collaborator, err = GetOrganizationCollaborator(org, request.Email); err != nil {
		log.Error().Err(err).Msg("could not retrieve collaborator from organization")
		c.JSON(http.StatusInternalServerError, "could not add collaborator")
		return
	}

	c.JSON(http.StatusOK, collaborator)
}
