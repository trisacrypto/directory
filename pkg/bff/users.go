package bff

import (
	"fmt"
	"net/http"

	"github.com/auth0/go-auth0/management"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/auth"
)

const (
	// TODO: do not hard code this value but make it a configuration
	DefaultRole = "Organization Collaborator"
)

// Login performs post-authentication checks and ensures that the user has the proper
// permissions and roles after they sign in with Auth0. The front-end should call the
// BFF login endpoint after the user signs in, providing the access_token in the
// request. If there is no access token a 401 is returned. This endpoint verifies that
// the user has a role and organization assigned to it and that the organization is up
// to date with the auth0 app_data. If the user does not have an organization, it is
// assumed that this is the first time the user has logged in and an organization is
// created for the user and they are assigned the organization leader role. If they have
// an organization but no role, they are assigned the organization collaborator role.
func (s *Server) Login(c *gin.Context) {
	var (
		err   error
		user  *management.User
		roles *management.RoleList
	)

	// Fetch the user from the context
	if user, err = auth.GetUserInfo(c); err != nil {
		log.Error().Err(err).Msg("login handler requires user info; expected middleware to return 401")
		c.JSON(http.StatusInternalServerError, "could not identify user to login")
		return
	}

	// Fetch user roles.
	if roles, err = s.auth0.User.Roles(*user.ID); err != nil {
		log.Error().Err(err).Msg("could not fetch roles associated with the user")
		c.JSON(http.StatusInternalServerError, "could not complete user login")
		return
	}

	if len(roles.Roles) == 0 {
		// Assign the user the organization collaborator role
		var collaborator *management.Role
		if collaborator, err = s.FindRoleByName(DefaultRole); err != nil {
			log.Error().Err(err).Msg("could not identify the default role to assign the user")
			c.JSON(http.StatusInternalServerError, "could not complete user login")
			return
		}

		if err = s.auth0.Role.AssignUsers(*collaborator.ID, []*management.User{user}); err != nil {
			log.Error().Err(err).Msg("could not assign the default role to the user")
			c.JSON(http.StatusInternalServerError, "could not complete user login")
			return
		}
	}

	// TODO: deal with user resources (will happen in a different story).

	// Once work has been performed reply with success no content
	c.Status(http.StatusNoContent)
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
