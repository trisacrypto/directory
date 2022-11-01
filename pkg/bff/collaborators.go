package bff

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
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

	// TODO: We can search the email address in Auth0 to see if the user already exists
	// TODO: Send invite/verification email to the collaborator

	// Save the updated organization
	if err = s.db.UpdateOrganization(org); err != nil {
		log.Error().Err(err).Msg("could not save organization with new collaborator")
		c.JSON(http.StatusInternalServerError, "could not add collaborator")
		return
	}

	c.JSON(http.StatusOK, collaborator)
}

// ListCollaborators lists all the collaborators on the user's organization. The user
// must have the read:collaborators permission to make this request.
func (s *Server) ListCollaborators(c *gin.Context) {
	var (
		err          error
		org 		*models.Organization
	)

	// Fetch the organization from the claims
	// NOTE: This method handles the error logging and response
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Build the response from the internal map
	out := &api.ListCollaboratorsReply{
		Collaborators: make([]*models.Collaborator, 0),
	}

	for _, collab := range org.Collaborators {
		if err = s.LoadCollaboratorDetails(collab); err != nil {
			log.Error().Err(err).Str("collabID", collab.Key()).Msg("could not load collaborator details from Auth0")
		}
		out.Collaborators = append(out.Collaborators, collab)
	}

	// Sort by email address so we return consistent responses
	sort.Slice(out.Collaborators, func(i, j int) bool {
		return out.Collaborators[i].Email < out.Collaborators[j].Email
	})

	c.JSON(http.StatusOK, out)
}

// UpdateCollaboratorRoles updates the roles of the collaborator ID in the request,
// ensuring that the roles are updated both on the organization record and in Auth0.
// The user must have the update:collaborators permission to make this request.
func (s *Server) UpdateCollaboratorRoles(c *gin.Context) {
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

	// Unmarshal the roles from the POST request
	params := &api.UpdateRolesParams{}
	if err = c.ShouldBind(params); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Make sure a 404 is returned if the collaborator does not exist
	var ok bool
	if collaborator, ok = org.Collaborators[collabID]; !ok {
		log.Warn().Str("collabID", collabID).Msg("collaborator does not exist")
		c.JSON(http.StatusNotFound, api.ErrorResponse("collaborator does not exist"))
		return
	}

	// Collaborator needs to be a verified user in Auth0
	if collaborator.VerifiedAt == "" || collaborator.UserId == "" {
		log.Warn().Str("collabID", collabID).Msg("cannot update roles for unverified collaborator")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("cannot update roles for unverified collaborator"))
		return
	}

	// Validate the specified roles in Auth0
	var newRoles []*management.Role
	for _, name := range params.Roles {
		var role *management.Role
		if role, err = s.FindRoleByName(name); err != nil {
			log.Warn().Err(err).Str("role", name).Msg("could not find role in Auth0")
			c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
			return
		}
		newRoles = append(newRoles, role)
	}

	// Remove the user's existing roles so we can add the new ones
	// TODO: There might be a more atomic way to do this in the management API
	var userRoles *management.RoleList
	if userRoles, err = s.auth0.User.Roles(collaborator.UserId); err != nil {
		log.Error().Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not fetch user roles from Auth0")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch user roles from Auth0"))
		return
	}

	if err = s.auth0.User.RemoveRoles(collaborator.UserId, userRoles.Roles); err != nil {
		log.Error().Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not update user roles in Auth0")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user roles in Auth0"))
		fmt.Println(err)
		return
	}

	// Add the new roles to the user
	if err = s.auth0.User.AssignRoles(collaborator.UserId, newRoles); err != nil {
		log.Error().Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not update user roles in Auth0")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user roles in Auth0"))
		return
	}

	// Update the collaborator record with the info in Auth0
	if err = s.LoadCollaboratorDetails(collaborator); err != nil {
		log.Error().Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not update collaborator record")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update collaborator record"))
		return
	}

	// Make sure the collaborator record has updated timestamps
	collaborator.ModifiedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if collaborator.CreatedAt == "" {
		collaborator.CreatedAt = collaborator.ModifiedAt
	}

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

	// Make sure a 404 is returned if the collaborator does not exist
	var ok bool
	if collaborator, ok = org.Collaborators[collabID]; !ok {
		log.Warn().Str("collabID", collabID).Msg("collaborator does not exist")
		c.JSON(http.StatusNotFound, api.ErrorResponse("collaborator not found"))
		return
	}

	// If the collabroator is already verified in Auth0, then remove them from the
	// organization
	if collaborator.VerifiedAt != "" && collaborator.UserId != "" {
		// Fetch the user from Auth0
		var user *management.User
		if user, err = s.auth0.User.Read(collaborator.UserId); err != nil {
			log.Error().Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not fetch user from Auth0")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch user from Auth0"))
			return
		}

		// Fetch the user app metadata
		appdata := &auth.AppMetadata{}
		if err = appdata.Load(user.AppMetadata); err != nil {
			log.Error().Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not parse user app metadata")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
			return
		}

		// Update the app metadata with the removed organization
		appdata.OrgID = ""
		if err = s.SaveAuth0AppMetadata(*user.ID, *appdata); err != nil {
			log.Error().Err(err).Str("collabID", collabID).Str("auth0_id", collaborator.UserId).Msg("could not save user app metadata")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not save user app metadata"))
			return
		}
	}

	// Delete the collaborator from the organization record
	delete(org.Collaborators, collabID)

	// Save the updated organization
	if err = s.db.UpdateOrganization(org); err != nil {
		log.Error().Err(err).Msg("could not save organization without collaborator")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete collaborator"))
		return
	}

	c.Status(http.StatusOK)
}

// LoadCollaboratorDetails updates a collaborator record with the user details in
// Auth0. The collaborator must have an user ID on it and the data in Auth0 will
// overwrite the data on the collaborator record.
func (s *Server) LoadCollaboratorDetails(collab *models.Collaborator) (err error) {
	// If the user is not verified in Auth0 then we can't retrieve the details
	if collab.VerifiedAt == "" {
		return nil
	}

	// Record must have a user ID to query Auth0
	if collab.UserId == "" {
		return errors.New("collaborator does not have a user ID")
	}

	// Load the user details
	// Note: This assumes that the user does not change or have multiple email addresses
	var user *management.User
	if user, err = s.auth0.User.Read(collab.UserId); err != nil {
		return fmt.Errorf("could not load user %q from Auth0: %w", collab.UserId, err)
	}

	if user.Name != nil {
		collab.Name = *user.Name
	}

	// Load the user roles
	var roles *management.RoleList
	if roles, err = s.auth0.User.Roles(collab.UserId); err != nil {
		return fmt.Errorf("could not load roles for user %q from Auth0: %w", collab.UserId, err)
	}
	collab.Roles = make([]string, len(roles.Roles))
	for i, role := range roles.Roles {
		collab.Roles[i] = *role.Name
	}

	return nil
}
