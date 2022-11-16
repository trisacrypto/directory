package bff

import (
	"fmt"
	"net/http"
	"time"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
)

const (
	// TODO: Need to make sure these roles are in sync with the roles in Auth0
	LeaderRole         = "Organization Leader"
	CollaboratorRole   = "Organization Collaborator"
	TSPRole            = "TRISA Service Provider"
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
// should be assigned to. This parameter is used to faciliate organization switching
// from the frontend as well as completing the invite workflow for new collaborators
// joining an organization. If the orgID is not provided as part of the request or does
// not exist in the user's app metadata, a new organization is automatically created for
// them and they are assigned the organization leader role. If the auth0 app data was
// changed, this returns a response with the refresh_token field set to true,
// indicating that the frontend should refresh the access token to ensure that the user
// claims are up to date.
func (s *Server) Login(c *gin.Context) {
	var (
		err   error
		user  *management.User
		roles *management.RoleList
	)

	// Parse optional params
	params := &api.LoginParams{}
	if err = c.ShouldBind(params); err != nil {
		log.Error().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Fetch the user from the context
	if user, err = auth.GetUserInfo(c); err != nil {
		log.Error().Err(err).Msg("login handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not identify user to login"))
		return
	}

	// Ensure the user resources are correctly populated.
	// If the user is not associated with an organization, create it.
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		log.Error().Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse user app metadata"))
		return
	}

	// Fetch the current user role.
	if roles, err = s.auth0.User.Roles(*user.ID); err != nil {
		log.Error().Err(err).Msg("could not fetch roles associated with the user")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
		return
	}

	// Users should only have one role
	var prevRole, userRole string
	switch len(roles.Roles) {
	case 0:
		// Default users to the organization collaborator role
		userRole = CollaboratorRole
	case 1:
		prevRole = *roles.Roles[0].Name
		userRole = prevRole
	default:
		// TODO: Resolve the conflict rather than returning an error
		log.Error().Err(err).Msg("user has multiple roles")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
		return
	}

	var org *models.Organization
	if params.OrgID == "" && appdata.OrgID == "" {
		// This is a new user so create a new organization for them
		org = &models.Organization{}
		if _, err = s.db.CreateOrganization(org); err != nil {
			log.Error().Err(err).Msg("could not create organization for new user")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
			return
		}

		// Add the user to the organization in the database
		collaborator := &models.Collaborator{
			Email:    *user.Email,
			UserId:   *user.ID,
			Verified: *user.EmailVerified,
		}
		if err = org.AddCollaborator(collaborator); err != nil {
			log.Error().Err(err).Str("user_id", collaborator.UserId).Msg("could not add collaborator to organization")
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
			log.Error().Err(err).Str("org_id", orgID).Msg("could not fetch organization for invited user")
			c.JSON(http.StatusNotFound, api.ErrorResponse("organization not found"))
			return
		}

		// Retrieve the user's collaborator record from the organization
		// Note: This is a critical security check because it ensures that the user was
		// really invited by an organization leader via the AddCollaborator endpoint
		// which started the invite workflow. Without this check, any user could log
		// into any organization simply by providing the orgID in the request.
		var collaborator *models.Collaborator
		if collaborator = org.GetCollaborator(*user.Email); collaborator == nil {
			log.Debug().Str("email", *user.Email).Str("org_id", org.Id).Msg("could not find user in organization")
			c.JSON(http.StatusUnauthorized, api.ErrorResponse("user is not authorized to access this organization"))
			return
		}

		// Other endpoints expect the user's verification status to be up to date
		collaborator.Verified = *user.EmailVerified
	}

	if userRole == TSPRole {
		// TSP users can be added to multiple organizations
		appdata.AddOrganization(org.Id)
	} else {
		// Organizations with one collaborator need a leader to add other collaborators
		if userRole == CollaboratorRole && len(org.Collaborators) == 1 {
			userRole = LeaderRole
		}

		// Non-TSP users can only exist in one organization
		if appdata.OrgID != "" && org.Id != appdata.OrgID {
			// When switching to a new organization, make sure leader roles are not
			// unintentionally preserved
			if userRole == LeaderRole && len(org.Collaborators) > 1 {
				userRole = CollaboratorRole
			}

			// Remove the collaborator record from the previous organization
			// TODO: This might require a user confirmation prompt
			var prevOrg *models.Organization
			if prevOrg, err = s.OrganizationFromID(appdata.OrgID); err != nil {
				log.Error().Err(err).Str("org_id", appdata.OrgID).Msg("could not fetch organization for user migration")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
				return
			}
			prevOrg.DeleteCollaborator(*user.Email)

			// If the previous organization has no collaborators, delete it
			if len(prevOrg.Collaborators) == 0 {
				if err = s.db.DeleteOrganization(prevOrg.UUID()); err != nil {
					log.Error().Err(err).Str("org_id", prevOrg.Id).Msg("could not delete organization")
					c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
					return
				}
			} else if err = s.db.UpdateOrganization(prevOrg); err != nil {
				log.Error().Err(err).Str("org_id", prevOrg.Id).Msg("could not update organization")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
				return
			}
		}
	}

	// Update the user app metadata to reflect the user's currently selected
	// organization and make sure the metadata is up to date in Auth0.
	appdata.UpdateOrganization(org)
	if err = s.SaveAuth0AppMetadata(*user.ID, *appdata); err != nil {
		log.Error().Err(err).Str("user_id", *user.ID).Msg("could not save user app_metadata")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
		return
	}

	// Update the organization record in the database
	if err = s.db.UpdateOrganization(org); err != nil {
		log.Error().Err(err).Str("org_id", org.Id).Msg("could not update organization")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
		return
	}

	// Assign a new user role if necessary
	if userRole != prevRole {
		if err = s.AssignRoles(*user.ID, []string{userRole}); err != nil {
			log.Error().Err(err).Str("user_id", *user.ID).Msg("could not assign user role")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not complete user login"))
			return
		}
	}

	// Protect the front-end by setting double cookie tokens for CSRF protection.
	// TODO: should we set expires at to the expiration of the access token? What happens on refresh?
	expiresAt := time.Now().Add(DoubleCookieMaxAge)
	if err := auth.SetDoubleCookieToken(c, s.conf.CookieDomain, expiresAt); err != nil {
		log.Error().Err(err).Msg("could not set double cookie csrf protection")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not set csrf protection"))
		return
	}

	// Get the old user app metadata for comparison
	oldAppdata := &auth.AppMetadata{}
	if err = oldAppdata.Load(user.AppMetadata); err != nil {
		log.Error().Err(err).Msg("could not parse user app metadata")
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

// AssignRoles assigns a set of roles to a user by ID, removing the existing roles and
// replacing them with the new set.
func (s *Server) AssignRoles(userID string, roles []string) (err error) {
	// TODO: There might be a more atomic way to do this.

	// Validate the specified roles in Auth0
	var newRoles []*management.Role
	for _, name := range roles {
		var role *management.Role
		if role, err = s.FindRoleByName(name); err != nil {
			log.Error().Err(err).Str("role", name).Msg("could not find role in Auth0")
			return ErrInvalidUserRole
		}
		newRoles = append(newRoles, role)
	}

	// Get the existing roles for the user
	var userRoles *management.RoleList
	if userRoles, err = s.auth0.User.Roles(userID); err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("could not fetch user roles from Auth0")
		return err
	}

	// Remove the existing roles from the user
	// The management endpoint requires a non-empty list, otherwise it returns a 400
	if len(userRoles.Roles) > 0 {
		if err = s.auth0.User.RemoveRoles(userID, userRoles.Roles); err != nil {
			log.Error().Err(err).Str("user_id", userID).Msg("could not remove existing roles from user")
			return err
		}
	}

	// Assign the new roles to the user
	// The management endpoint requires a non-empty list, otherwise it returns a 400
	if len(newRoles) > 0 {
		if err = s.auth0.User.AssignRoles(userID, newRoles); err != nil {
			log.Error().Err(err).Str("user_id", userID).Msg("could not add new roles to user")
			return err
		}
	}

	return nil
}

// ListUserRoles returns the list of assignable user roles.
func (s *Server) ListUserRoles(c *gin.Context) {
	// TODO: This is currently a static list which must be maintained to be in sync
	// with the roles defined in Auth0.
	c.JSON(http.StatusOK, []string{CollaboratorRole, LeaderRole})
}

func (s *Server) FindRoleByName(name string) (*management.Role, error) {
	roles, err := s.auth0.Role.List()
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

func (s *Server) SaveAuth0AppMetadata(uid string, appdata auth.AppMetadata) (err error) {
	// Create a blank user with no data but the appdata
	user := &management.User{}

	// Send the updated user app_metadata back to auth0
	if user.AppMetadata, err = appdata.Dump(); err != nil {
		return err
	}

	// Patch the user with the specified user ID
	if err = s.auth0.User.Update(uid, user); err != nil {
		return err
	}

	return nil
}
