package bff

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

const (
	// TODO: Need to make sure these roles are in sync with the roles in Auth0
	DoubleCookieMaxAge = 24 * time.Hour
	OrgIDKey           = "orgid"
	VASPsKey           = "vasps"
)

// Login performs post-authentication checks and ensures that the user has the proper
// permissions and roles after they sign in with Auth0. The front-end should call the
// BFF login endpoint after the user signs in, providing the access_token in the
// request. If there is no access token a 401 is returned. This endpoint verifies that
// the user has a role and organization assigned to it and that the organization is up
// to date with the auth0 app_data.
//
// By default, this endpoint attempts to log the user into their last used
// organization, using the orgID in the user app metadata. The endpoint also accepts an
// orgID parameter as part of the request which determines the organization the user
// should be assigned to. This parameter is used to facilitate organization switching
// from the frontend as well as completing the invite workflow for new collaborators
// joining an organization. If the orgID is not provided as part of the request or does
// not exist in the user's app metadata, a new organization is automatically created for
// them and they are assigned the organization leader role. If the auth0 app data was
// changed, this returns a response with the refresh_token field set to true,
// indicating that the frontend should refresh the access token to ensure that the user
// claims are up to date.
//
// @Summary Login a user to the BFF
// @Description Completes the user login process by assigning the user to an organization and verifying that the user has the proper roles.
// @Tags users
// @Accept json
// @Produce json
// @Param params body api.LoginParams true "Login parameters"
// @Success 200 {object} api.Reply "Login successful, token refresh required"
// @Success 204 "Login successful"
// @Failure 400 {object} api.Reply
// @Failure 401 {object} api.Reply
// @Failure 403 {object} api.Reply "User invitation has expired"
// @Failure 404 {object} api.Reply "Organization not found"
// @Failure 500 {object} api.Reply
// @Router /users/login [post]
func (s *Server) Login(c *gin.Context) {
	var (
		err   error
		user  *management.User
		roles *management.RoleList
	)

	// Parse optional params
	params := &api.LoginParams{}
	if err = c.ShouldBind(params); err != nil {
		sentry.Error(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Fetch the user from the context
	if user, err = auth.GetUserInfo(c); err != nil {
		sentry.Error(c).Err(err).Msg("login handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify user to login"))
		return
	}

	// Fetch the user's claims
	var claims *auth.Claims
	if claims, err = auth.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not fetch user claims")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
		return
	}

	// Ensure the user resources are correctly populated.
	// If the user is not associated with an organization, create it.
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		sentry.Error(c).Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
		return
	}

	// Fetch the current user role.
	if roles, err = s.auth0.User.Roles(c.Request.Context(), *user.ID); err != nil {
		sentry.Error(c).Err(err).Msg("could not fetch roles associated with the user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
		return
	}

	// Users should only have one role
	var prevRole, userRole string
	switch len(roles.Roles) {
	case 0:
		// Default users to the organization collaborator role
		userRole = auth.CollaboratorRole
	case 1:
		prevRole = *roles.Roles[0].Name
		userRole = prevRole
	default:
		// TODO: Resolve the conflict rather than returning an error
		sentry.Error(c).Err(err).Msg("user has multiple roles")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
		return
	}

	var (
		org          *models.Organization
		collaborator *models.Collaborator
	)
	if params.OrgID == "" && appdata.OrgID == "" {
		// This is a new user so create a new organization for them
		var userName string
		if userName, err = auth.UserDisplayName(user); err != nil {
			sentry.Error(c).Err(err).Str("user_id", *user.ID).Msg("could not get user display name")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
			return
		}
		org = &models.Organization{
			CreatedBy: userName,
		}
		ctx, cancel := utils.WithDeadline(context.Background())
		defer cancel()
		if _, err = s.db.CreateOrganization(ctx, org); err != nil {
			sentry.Error(c).Err(err).Str("user_id", *user.ID).Msg("could not create organization")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
			return
		}

		// Add the user to the organization in the database
		collaborator = &models.Collaborator{
			Email:    *user.Email,
			UserId:   *user.ID,
			Verified: *user.EmailVerified,
		}
		if err = org.AddCollaborator(collaborator); err != nil {
			sentry.Error(c).Err(err).Str("user_id", collaborator.UserId).Msg("could not add collaborator to organization")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
			return
		}
	} else {
		var orgID string
		if params.OrgID != "" {
			// Try to login the user to the requested organization if provided
			orgID = params.OrgID
		} else {
			// Default to the last used organization
			orgID = appdata.OrgID
		}

		// Fetch the organization from the database
		if org, err = s.OrganizationFromID(orgID); err != nil {
			sentry.Error(c).Err(err).Str("org_id", orgID).Msg("could not fetch organization for invited user")
			c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
			return
		}

		// Retrieve the user's collaborator record from the organization
		// Note: This is a critical security check because it ensures that the user was
		// really invited by an organization leader via the AddCollaborator endpoint
		// which started the invite workflow. Without this check, any user could log
		// into any organization simply by providing the orgID in the request.
		if collaborator = org.GetCollaborator(*user.Email); collaborator == nil {
			sentry.Warn(c).Str("email", *user.Email).Str("org_id", org.Id).Msg("could not find user in organization")
			c.JSON(http.StatusUnauthorized, api.ErrorResponse("user is not authorized to access this organization"))
			return
		}

		// Verify that pending invitations have not expired
		if err = collaborator.ValidateInvitation(); err != nil {
			sentry.Warn(c).Err(err).Str("email", collaborator.Email).Str("org_id", org.Id).Msg("invalid user invitation")
			c.JSON(http.StatusForbidden, api.ErrorResponse("user invitation has expired"))
			return
		}

		// Other endpoints expect the user's verification status to be up to date
		collaborator.Verified = *user.EmailVerified
		collaborator.UserId = *user.ID
		collaborator.ExpiresAt = ""
	}

	// Update collaborator metadata timestamps when the user logs in
	collaborator.LastLogin = time.Now().Format(time.RFC3339Nano)
	if collaborator.JoinedAt == "" {
		collaborator.JoinedAt = collaborator.LastLogin
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	if claims.HasPermission(auth.SwitchOrganizations) {
		// If the user can be in multiple organizations, add the selected organization
		// to the user's list.
		appdata.AddOrganization(org.Id)
	} else {
		// Make sure users are able to add collaborators if they are the only member of
		// their organization.
		if !claims.HasPermission(auth.UpdateCollaborators) && len(org.Collaborators) == 1 {
			userRole = auth.LeaderRole
		}

		// Non-TSP users can only exist in one organization
		if appdata.OrgID != "" && org.Id != appdata.OrgID {
			// Users should not be able to add collaborators if they are migrated to an
			// existing organization.
			if claims.HasPermission(auth.UpdateCollaborators) && len(org.Collaborators) > 1 {
				userRole = auth.CollaboratorRole
			}

			// Remove the collaborator record from the previous organization
			// TODO: This might require a user confirmation prompt
			var prevOrg *models.Organization
			if prevOrg, err = s.OrganizationFromID(appdata.OrgID); err != nil {
				sentry.Error(c).Err(err).Str("org_id", appdata.OrgID).Msg("could not fetch organization for user migration")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
				return
			}
			prevOrg.DeleteCollaborator(*user.Email)

			// If the previous organization has no collaborators, delete it
			if len(prevOrg.Collaborators) == 0 {
				if err = s.db.DeleteOrganization(ctx, prevOrg.UUID()); err != nil {
					sentry.Error(c).Err(err).Str("org_id", prevOrg.Id).Msg("could not delete organization")
					c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
					return
				}
			} else if err = s.db.UpdateOrganization(ctx, prevOrg); err != nil {
				sentry.Error(c).Err(err).Str("org_id", prevOrg.Id).Msg("could not update organization")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
				return
			}
		}
	}

	// Update the user app metadata to reflect the user's currently selected
	// organization and make sure the metadata is up to date in Auth0.
	appdata.UpdateOrganization(org)
	if err = s.SaveAuth0AppMetadata(c.Request.Context(), *user.ID, *appdata); err != nil {
		sentry.Error(c).Err(err).Str("user_id", *user.ID).Msg("could not save user app_metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
		return
	}

	// Update the organization record in the database
	if err = s.db.UpdateOrganization(ctx, org); err != nil {
		sentry.Error(c).Err(err).Str("org_id", org.Id).Msg("could not update organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
		return
	}

	// Assign a new user role if necessary
	if userRole != prevRole {
		if err = s.AssignRoles(c.Request.Context(), *user.ID, []string{userRole}); err != nil {
			sentry.Error(c).Err(err).Str("user_id", *user.ID).Msg("could not assign user role")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
			return
		}
	}

	// Protect the front-end by setting double cookie tokens for CSRF protection.
	// TODO: should we set expires at to the expiration of the access token? What happens on refresh?
	expiresAt := time.Now().Add(DoubleCookieMaxAge)
	if err := auth.SetDoubleCookieToken(c, s.conf.CookieDomain, expiresAt); err != nil {
		sentry.Error(c).Err(err).Msg("could not set double cookie csrf protection")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not set csrf protection"))
		return
	}

	// Get the old user app metadata for comparison
	oldAppdata := &auth.AppMetadata{}
	if err = oldAppdata.Load(user.AppMetadata); err != nil {
		sentry.Error(c).Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
		return
	}

	// If the user app metadata has changed, set the refresh flag in the response
	if !appdata.Equals(oldAppdata) {
		c.JSON(http.StatusOK, api.Reply{Success: true, RefreshToken: true})
	} else {
		c.Status(http.StatusNoContent)
	}
}

// UpdateUser updates the user's profile information in Auth0.
//
// @Summary Update the user's profile
// @Description Update the user's profile information in Auth0.
// @Tags users
// @Accept json
// @Success 204 {object} api.Reply
// @Failure 400 {object} api.Reply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /users [patch]
func (s *Server) UpdateUser(c *gin.Context) {
	var (
		user *management.User
		err  error
	)

	// Parse the params from the request body
	params := &api.UpdateUserParams{}
	if err = c.ShouldBind(params); err != nil {
		sentry.Error(c).Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Fetch the user from the context
	if user, err = auth.GetUserInfo(c); err != nil {
		sentry.Error(c).Err(err).Msg("login handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify user to login"))
		return
	}

	// At least one field must be provided
	if *params == (api.UpdateUserParams{}) {
		sentry.Warn(c).Msg("no fields were provided to update user")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("no fields were provided"))
		return
	}

	// Update the user's name if provided
	patch := &management.User{}
	if params.Name != "" {
		patch.Name = &params.Name
	}

	// Commit the update to Auth0
	if err = s.auth0.User.Update(c.Request.Context(), *user.ID, patch); err != nil {
		sentry.Error(c).Err(err).Str("user_id", *user.ID).Msg("could not update user name")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update user name"))
		return
	}

	// Invalidate the user's cache entry so that the updated profile is returned from
	// the backend.
	s.users.Remove(*user.ID)

	c.Status(http.StatusNoContent)
}

// UserOrganization returns the current organization that the user is logged into. The
// user must have the read:organizations permission to perform this action.
//
// @Summary Get the user's current organization [read:organizations]
// @Description Get high level info about the user's current organization
// @Tags users
// @Produce json
// @Success 200 {object} api.OrganizationReply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /users/organization [get]
func (s *Server) UserOrganization(c *gin.Context) {
	var (
		err error
		org *models.Organization
	)

	// Fetch the organization from the claims
	// Note: This method handles the error logging and response
	if org, err = s.OrganizationFromClaims(c); err != nil {
		return
	}

	// Build the response
	reply := &api.OrganizationReply{
		ID:        org.Id,
		Name:      org.ResolveName(),
		Domain:    org.Domain,
		CreatedAt: org.Created,
	}

	c.JSON(http.StatusOK, reply)
}

// AssignRoles assigns a set of roles to a user by ID, removing the existing roles and
// replacing them with the new set.
func (s *Server) AssignRoles(ctx context.Context, userID string, roles []string) (err error) {
	// TODO: There might be a more atomic way to do this.

	// Validate the specified roles in Auth0
	var newRoles []*management.Role
	for _, name := range roles {
		var role *management.Role
		if role, err = s.FindRoleByName(ctx, name); err != nil {
			sentry.Error(nil).Err(err).Str("role", name).Msg("could not find role in Auth0")
			return ErrInvalidUserRole
		}
		newRoles = append(newRoles, role)
	}

	// Get the existing roles for the user
	var userRoles *management.RoleList
	if userRoles, err = s.auth0.User.Roles(ctx, userID); err != nil {
		sentry.Error(nil).Err(err).Str("user_id", userID).Msg("could not fetch user roles from Auth0")
		return err
	}

	// Remove the existing roles from the user
	// The management endpoint requires a non-empty list, otherwise it returns a 400
	if len(userRoles.Roles) > 0 {
		if err = s.auth0.User.RemoveRoles(ctx, userID, userRoles.Roles); err != nil {
			sentry.Error(nil).Err(err).Str("user_id", userID).Msg("could not remove existing roles from user")
			return err
		}
	}

	// Assign the new roles to the user
	// The management endpoint requires a non-empty list, otherwise it returns a 400
	if len(newRoles) > 0 {
		if err = s.auth0.User.AssignRoles(ctx, userID, newRoles); err != nil {
			sentry.Error(nil).Err(err).Str("user_id", userID).Msg("could not add new roles to user")
			return err
		}
	}

	// Invalidate the user's cache entry so updated roles are returned from the backend
	s.users.Remove(userID)

	return nil
}

// ListUserRoles returns the list of assignable user roles.
//
// @Summary Get the list of assignable user roles
// @Description Get the list of assignable user roles
// @Tags users
// @Produce json
// @Success 200 {list} string
// @Router /users/roles [get]
func (s *Server) ListUserRoles(c *gin.Context) {
	// TODO: This is currently a static list which must be maintained to be in sync
	// with the roles defined in Auth0.
	c.JSON(http.StatusOK, []string{auth.CollaboratorRole, auth.LeaderRole})
}

// FindUserByEmail returns the Auth0 user record by email address. This method returns
// an ErrUserEmailNotFound error if the user does not exist and returns the first user
// if there are multiple users with the same email address.
func (s *Server) FindUserByEmail(ctx context.Context, email string) (user *management.User, err error) {
	var users []*management.User
	if users, err = s.auth0.User.ListByEmail(ctx, strings.ToLower(email)); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrUserEmailNotFound
	}

	if len(users) > 1 {
		// TODO: This can happen if the user has authenticated with Auth0 using
		// multiple identities (e.g. email and Google). We might be able to handle this
		// by linking the identities.
		log.Warn().Str("email", email).Int("count", len(users)).Msg("multiple users found with same email address")
	}

	return users[0], nil
}

func (s *Server) FindRoleByName(ctx context.Context, name string) (*management.Role, error) {
	roles, err := s.auth0.Role.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, role := range roles.Roles {
		if *role.Name == name {
			return role, nil
		}
	}
	return nil, fmt.Errorf("could not find role %q in %d available roles", name, len(roles.Roles))
}

// Switch the user to an available organization by updating their app metadata on Auth0.
// This always clears the current organization info from the app metadata but only
// replaces it if another organization is found.
// TODO: This switches the user to the first valid organization in the list. Should we
// switch the user to their last used organization instead?
func (s *Server) SwitchUserOrganization(ctx context.Context, user *management.User, appdata *auth.AppMetadata) (err error) {
	// Clear out the old organization info
	appdata.ClearOrganization()

	// Find the first organization that the user is a collaborator in
	var org *models.Organization
	for _, id := range appdata.GetOrganizations() {
		if org, err = s.OrganizationFromID(id); err == nil && user.Email != nil && org.GetCollaborator(*user.Email) != nil {
			break
		}
	}

	// Update the app metadata with the organization ID if one was found
	if org != nil {
		appdata.UpdateOrganization(org)
	}

	// Save the updated app metadata to Auth0
	if err = s.SaveAuth0AppMetadata(ctx, *user.ID, *appdata); err != nil {
		return err
	}

	return nil
}

func (s *Server) SaveAuth0AppMetadata(ctx context.Context, uid string, appdata auth.AppMetadata) (err error) {
	// Create a blank user with no data but the appdata
	user := &management.User{}

	// Send the updated user app_metadata back to auth0
	if user.AppMetadata, err = appdata.Dump(); err != nil {
		return err
	}

	// Patch the user with the specified user ID
	if err = s.auth0.User.Update(ctx, uid, user); err != nil {
		return err
	}

	return nil
}
