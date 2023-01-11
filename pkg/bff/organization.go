package bff

import (
	"errors"
	"net/http"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
)

// CreateOrganization creates a new organization in the database. This endpoint returns
// an error if the organization already exists and the user is assigned to it. The user
// must have the create:organizations permission to perform this action.
//
// @Summary Create a new organization [create:organizations]
// @Description Create a new organization with the specified name and domain for the user.
// @Tags organizations
// @Accept json
// @Produce json
// @Param params body api.OrganizationParams true "Name and domain"
// @Success 200 {object} api.OrganizationReply
// @Failure 400 {object} api.Reply "Must provide name and domain"
// @Failure 401 {object} api.Reply
// @Failure 409 {object} api.Reply "Domain already exists"
// @Failure 500 {object} api.Reply
// @Router /organizations [post]
func (s *Server) CreateOrganization(c *gin.Context) {
	var (
		err  error
		user *management.User
	)

	// Fetch the user from the context
	if user, err = auth.GetUserInfo(c); err != nil {
		log.Error().Err(err).Msg("create organization handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify user to create organization"))
		return
	}

	// Load the user app metadata to check their organization assignments
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		log.Error().Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
		return
	}

	// Unmarshal the params from the POST request
	params := &api.OrganizationParams{}
	if err := c.ShouldBind(params); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Name is a required parameter
	if params.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("must provide name in request params"))
		return
	}

	// Domain is a required parameter
	if params.Domain == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("must provide domain in request params"))
		return
	}

	// Don't allow the user to create duplicate organizations
	// TODO: Should we do a universal check against the database using an index?
	var domain string
	if domain, err = s.ValidateOrganizationDomain(params.Domain, appdata); err != nil {
		log.Error().Err(err).Str("domain", params.Domain).Msg("could not validate organization domain")
		if errors.Is(err, ErrDomainAlreadyExists) {
			c.JSON(http.StatusConflict, api.ErrorResponse("organization with domain already exists"))
			return
		}
		c.JSON(http.StatusBadRequest, api.ErrorResponse("invalid domain provided"))
		return
	}

	// Create a new organization in the database with the provided name and domain
	org := &models.Organization{
		Name:   params.Name,
		Domain: domain,
	}

	// CreatedBy is used to render the organization name for the frontend if no other
	// organization name is available
	if org.CreatedBy, err = auth.UserDisplayName(user); err != nil {
		log.Error().Err(err).Msg("could not get user display name")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not resolve name for user"))
		return
	}

	// Add the user to their new organization
	// Note: The UserInfo middleware ensures that these fields are present in the Auth0
	// user record
	collaborator := &models.Collaborator{
		Email:    *user.Email,
		UserId:   *user.ID,
		Verified: *user.EmailVerified,
	}
	if err = org.AddCollaborator(collaborator); err != nil {
		log.Error().Err(err).Str("user_id", collaborator.UserId).Msg("could not add user as collaborator in organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create organization"))
		return
	}

	if _, err = s.db.CreateOrganization(org); err != nil {
		log.Error().Err(err).Msg("could not create organization in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create organization"))
		return
	}

	// Assign the user to the organization
	appdata.AddOrganization(org.Id)
	if err = s.SaveAuth0AppMetadata(*user.ID, *appdata); err != nil {
		log.Error().Err(err).Msg("could not update user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user app metadata"))
		return
	}

	// Build the response
	out := &api.OrganizationReply{
		ID:           org.Id,
		Name:         org.Name,
		Domain:       org.Domain,
		CreatedAt:    org.Created,
		RefreshToken: true,
	}
	c.JSON(http.StatusOK, out)
}

// ListOrganizations returns a list of organizations that the user is a member of. The
// user must have the read:organizations permission to perform this action.
//
// @Summary List organizations [read:organizations]
// @Description Return the list of organizations that the user is assigned to.
// @Tags organizations
// @Produce json
// @Success 200 {list} api.OrganizationReply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /organizations [get]
func (s *Server) ListOrganizations(c *gin.Context) {
	var (
		err  error
		user *management.User
	)

	// Fetch the user from the context
	if user, err = auth.GetUserInfo(c); err != nil {
		log.Error().Err(err).Msg("list organizations handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify user to list organizations"))
		return
	}

	// Load the user app metadata to check their organization assignments
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		log.Error().Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
		return
	}

	// Build the response
	out := make([]*api.OrganizationReply, 0)
	for _, id := range appdata.GetOrganizations() {
		if org, err := s.OrganizationFromID(id); err != nil {
			log.Error().Err(err).Str("org_id", id).Msg("could not retrieve organization from database")
		} else {
			// User data is on the collaborator record.
			// Note: The UserInfo middleware ensures that the user email address is
			// present in the claims so we can safely dereference it here.
			var collaborator *models.Collaborator
			if collaborator = org.GetCollaborator(*user.Email); collaborator == nil {
				log.Error().Str("org_id", id).Str("email", *user.Email).Msg("could not find user in organization collaborators")
				continue
			}

			out = append(out, &api.OrganizationReply{
				ID:        org.Id,
				Name:      org.ResolveName(),
				Domain:    org.Domain,
				CreatedAt: org.Created,
				LastLogin: collaborator.LastLogin,
			})
		}
	}

	c.JSON(http.StatusOK, out)
}

// DeleteOrganization deletes an organization from the database. The user must have the
// delete:organizations permission and also be a collaborator in the organization to
// perform this action.
//
// @Summary Delete an organization [delete:organizations]
// @Description Completely delete an organization, including the registration and collaborators.
// @Tags organizations
// @Success 200 {object} api.Reply
// @Failure 401 {object} api.Reply
// @Failure 403 {object} api.Reply "User is not a collaborator in the organization"
// @Failure 404 {object} api.Reply "Organization not found"
// @Failure 500 {object} api.Reply
// @Router /organizations/{id} [delete]
func (s *Server) DeleteOrganization(c *gin.Context) {
	var (
		err  error
		user *management.User
		org  *models.Organization
	)

	// Fetch the user from the context
	if user, err = auth.GetUserInfo(c); err != nil {
		log.Error().Err(err).Msg("delete organization handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify user to delete organization"))
		return
	}

	// Parse the organization ID from the URL
	orgID := c.Param("orgID")

	// Fetch the organization to be deleted
	if org, err = s.OrganizationFromID(orgID); err != nil {
		log.Error().Err(err).Str("org_id", orgID).Msg("could not retrieve organization from database")
		c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
		return
	}

	// The user must be a collaborator in the organization to delete it
	if org.GetCollaborator(*user.Email) == nil {
		log.Error().Err(err).Str("org_id", orgID).Str("email", *user.Email).Msg("could not find user in organization collaborators")
		c.JSON(http.StatusForbidden, api.ErrorResponse("user is not authorized to access this organization"))
		return
	}

	// Delete the organization from the database
	if err = s.db.DeleteOrganization(org.UUID()); err != nil {
		log.Error().Err(err).Str("org_id", orgID).Msg("could not delete organization from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete organization from database"))
		return
	}

	// Remove the organization from all collaborators so that they don't get an error
	// when they try to log in. If a user no longer has an organization, they will be
	// assigned one automatically the next time they log in.
	for _, collab := range org.GetCollaborators() {
		// If the collaborator doesn't have an Auth0 ID then they haven't logged in yet
		if collab.UserId == "" {
			log.Debug().Str("org_id", orgID).Str("email", collab.Email).Msg("ignoring unverified collaborator during organization deletion")
			continue
		}

		// Retrieve the Auth0 user record from the user ID
		var collabUser *management.User
		if collabUser, err = s.auth0.User.Read(collab.UserId); err != nil {
			log.Error().Err(err).Str("user_id", collab.UserId).Msg("could not retrieve user from Auth0")
			continue
		}

		// Remove the organization from the user's organization list
		appdata := &auth.AppMetadata{}
		if err = appdata.Load(collabUser.AppMetadata); err != nil {
			log.Error().Err(err).Str("user_id", collab.UserId).Msg("could not parse user app metadata")
			continue
		}
		appdata.RemoveOrganization(org.Id)

		// If the user is currently assigned to the deleted organization, then switch
		// them to a different one to allow them to login.
		if appdata.OrgID == org.Id {
			// This method both modifies the app metadata and pushes the updates to
			// Auth0.
			if err = s.SwitchUserOrganization(collabUser, appdata); err != nil {
				log.Error().Err(err).Str("user_id", collab.UserId).Msg("could not switch user to new organization")
				continue
			}
		}
	}

	c.JSON(http.StatusOK, api.Reply{Success: true, RefreshToken: true})
}

// ValidateOrganizationDomain performs any necessary normalization and validation of an
// organization domain name, ensuring that the domain is not already in use by another
// organization on the specified app metadata and returning the normalized domain name
// for storage.
func (s *Server) ValidateOrganizationDomain(domain string, appdata *auth.AppMetadata) (string, error) {
	var err error

	// Normalize the domain
	if domain, err = NormalizeDomain(domain); err != nil {
		return "", err
	}

	// Check that the domain is valid
	if err = ValidateDomain(domain); err != nil {
		return "", err
	}

	// Check for duplicate domains
	for _, id := range appdata.GetOrganizations() {
		if org, err := s.OrganizationFromID(id); err != nil {
			log.Error().Err(err).Str("org_id", id).Msg("could not retrieve organization from database")
		} else if org.Domain == domain {
			return "", ErrDomainAlreadyExists
		}
	}

	return domain, nil
}

// OrganizationFromClaims is a helper method to retrieve the organization for a
// particular request by fetching the orgID from the claims and querying the database.
// If there is an error fetching the organization, the appropriate error response is
// made on the gin writer and logged. The caller should check for error and return.
func (s *Server) OrganizationFromClaims(c *gin.Context) (org *models.Organization, err error) {
	// Retrieve the organization ID from the claims
	var claims *auth.Claims
	if claims, err = auth.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not retrieve claims to fetch orgID")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify organization"))
		return nil, err
	}

	// If there is no organization ID, something went wrong
	if claims.OrgID == "" {
		log.Warn().Msg("claims do not contain an orgID")
		api.MustRefreshToken(c, "missing claims info, try logging out and logging back in")
		return nil, errors.New("missing organization ID in claims")
	}

	// Fetch the organization from the database
	if org, err = s.OrganizationFromID(claims.OrgID); err != nil {
		if errors.Is(err, storeerrors.ErrEntityNotFound) {
			log.Warn().Err(err).Msg("could not find organization in database from orgID in claims")
			api.MustRefreshToken(c, "no organization found, try logging out and logging back in")
			return nil, err
		}

		log.Error().Err(err).Msg("could not retrieve organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify organization"))
		return nil, err
	}

	return org, nil
}

func (s *Server) OrganizationFromID(id string) (org *models.Organization, err error) {
	// Organizations are stored by UUID in the database
	var uuid uuid.UUID
	if uuid, err = models.ParseOrgID(id); err != nil {
		return nil, err
	}

	// Fetch the record from the database
	if org, err = s.db.RetrieveOrganization(uuid); err != nil {
		return nil, err
	}

	return org, nil
}
