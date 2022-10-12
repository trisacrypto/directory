package bff

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
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

	// Make sure the collaborator is valid for storage
	if err = collaborator.Validate(); err != nil {
		log.Warn().Err(err).Msg("invalid collaborator in request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Don't overwrite an existing collaborator
	id := collaborator.Key()
	if _, ok := org.Collaborators[id]; ok {
		log.Warn().Str("collabID", id).Msg("collaborator already exists")
		c.JSON(http.StatusConflict, api.ErrorResponse("collaborator already exists"))
		return
	}

	// Make sure the record has a created timestamp
	if collaborator.CreatedAt == "" {
		collaborator.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}

	// Add the collaborator to the organization
	if org.Collaborators == nil {
		org.Collaborators = make(map[string]*models.Collaborator)
	}
	org.Collaborators[id] = collaborator

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

	// Get the collabID from the URL
	collabID := c.Param("collabID")
	if collabID == "" {
		log.Warn().Msg("missing ID in replace collaborator request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("ID is required in order to replace a collaborator"))
		return
	}

	// Unmarshal the collaborator from the PUT request
	collaborator = &models.Collaborator{}
	if err = c.ShouldBind(collaborator); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Make sure a 404 is returned if the collaborator does not exist
	if _, ok := org.Collaborators[collabID]; !ok {
		log.Warn().Str("collabID", collabID).Msg("collaborator does not exist")
		c.JSON(http.StatusNotFound, api.ErrorResponse("collaborator does not exist"))
		return
	}

	// Make sure the collaborator is valid for storage
	// Note: The user will not be able to update the email address or ID
	if err = collaborator.Validate(); err != nil {
		log.Warn().Err(err).Msg("invalid collaborator")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Make sure the collaborator record has updated timestamps
	collaborator.ModifiedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if collaborator.CreatedAt == "" {
		collaborator.CreatedAt = collaborator.ModifiedAt
	}

	// Replace the collaborator on the organization
	if org.Collaborators == nil {
		org.Collaborators = make(map[string]*models.Collaborator)
	}
	org.Collaborators[collabID] = collaborator

	// Save the updated organization
	if err = s.db.UpdateOrganization(org); err != nil {
		log.Error().Err(err).Msg("could not save organization with new collaborator")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not replace collaborator"))
		return
	}

	c.JSON(http.StatusOK, collaborator)
}

// DeleteCollaborator deletes the collaborator in the request from the user's
// organization. The user must have the update:collaborators permission.
// Note: This does not return an error if the collaborator does not exist on the
// organization and instead returns a 200 OK response.
func (s *Server) DeleteCollaborator(c *gin.Context) {
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

	// Get the collabID from the URL
	collabID := c.Param("collabID")
	if collabID == "" {
		log.Warn().Msg("missing ID in delete collaborator request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("ID is required in order to delete a collaborator"))
		return
	}

	// Make sure a 404 is returned if the collaborator does not exist
	if _, ok := org.Collaborators[collabID]; !ok {
		log.Warn().Str("collabID", collabID).Msg("collaborator does not exist")
		c.JSON(http.StatusNotFound, api.ErrorResponse("collaborator not found"))
		return
	}

	// Delete the collaborator from the organization
	delete(org.Collaborators, collabID)

	// Save the updated organization
	if err = s.db.UpdateOrganization(org); err != nil {
		log.Error().Err(err).Msg("could not save organization without collaborator")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete collaborator"))
		return
	}

	c.JSON(http.StatusOK, collaborator)
}
