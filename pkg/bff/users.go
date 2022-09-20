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
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

const (
	// TODO: do not hard code this value but make it a configuration
	DefaultRole        = "Organization Collaborator"
	DoubleCookieMaxAge = 24 * time.Hour
	OrgIDKey           = "orgid"
	VASPsKey           = "vasps"
)

// Login performs post-authentication checks and ensures that the user has the proper
// permissions and roles after they sign in with Auth0. The front-end should call the
// BFF login endpoint after the user signs in, providing the access_token in the
// request. If there is no access token a 401 is returned. This endpoint verifies that
// the user has a role and organization assigned to it and that the organization is up
// to date with the auth0 app_data. If the user does not have an organization, it is
// assumed that this is the first time the user has logged in and an organization is
// created for the user and they are assigned the organization leader role. If they have
// an organization but no role, they are assigned the organization collaborator role. If
// the auth0 app data was changed, this returns a response with the refresh_token field
// set to true, indicating that the frontend should refresh the access token to ensure
// that the user claims are up to date.
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

	// Ensure the user resources are correctly populated.
	// If the user is not associated with an organization, create it.
	appdata := &auth.AppMetadata{}
	if err = appdata.Load(user.AppMetadata); err != nil {
		log.Error().Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, "could not parse user app metadata")
		return
	}

	// Retrieve the user's organization from the database
	var org *models.Organization
	if appdata.OrgID == "" {
		// Create the organization
		org, err = s.db.Organizations().Create(c.Request.Context())
		if err != nil {
			log.Error().Err(err).Msg("could not create organization for new user")
			c.JSON(http.StatusInternalServerError, "could not complete user login")
			return
		}

		// Set the organization ID in the user app metadata
		appdata.OrgID = org.Id
	} else {
		// Get the organization for the specified user
		org, err = s.db.Organizations().Retrieve(c.Request.Context(), appdata.OrgID)
		if err != nil {
			log.Error().Err(err).Str("orgid", appdata.OrgID).Msg("could not retrieve organization for user VASP verification")
			c.JSON(http.StatusInternalServerError, "could not complete user login")
			return
		}

		// Ensure the VASP record is correct for the user
		if org.Testnet != nil && org.Testnet.Id != "" {
			appdata.VASPs.TestNet = org.Testnet.Id
		}
		if org.Mainnet != nil && org.Mainnet.Id != "" {
			appdata.VASPs.MainNet = org.Mainnet.Id
		}
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

		// Update the collaborators on the organization
		if err = AddOrganizationCollaborator(org, *user.Email); err != nil {
			log.Error().Err(err).Msg("could not add the user to the organization collaborators")
			c.JSON(http.StatusInternalServerError, "could not complete user login")
			return
		}

		// TODO: this will require the user to login again
		if err = s.auth0.Role.AssignUsers(*collaborator.ID, []*management.User{user}); err != nil {
			log.Error().Err(err).Msg("could not assign the default role to the user")
			c.JSON(http.StatusInternalServerError, "could not complete user login")
			return
		}
	}

	if err = s.SaveAuth0AppMetadata(*user.ID, *appdata); err != nil {
		log.Error().Err(err).Str("user_id", *user.ID).Msg("could not save user app_metadata")
		c.JSON(http.StatusInternalServerError, "could not complete user login")
		return
	}

	// Protect the front-end by setting double cookie tokens for CSRF protection.
	// TODO: should we set expires at to the expiration of the access token? What happens on refresh?
	expiresAt := time.Now().Add(DoubleCookieMaxAge)
	if err := auth.SetDoubleCookieToken(c, s.conf.CookieDomain, expiresAt); err != nil {
		log.Error().Err(err).Msg("could not set double cookie csrf protection")
		c.JSON(http.StatusInternalServerError, "could not set csrf protection")
		return
	}

	// Get the old user app metadata for comparison
	oldAppdata := &auth.AppMetadata{}
	if err = oldAppdata.Load(user.AppMetadata); err != nil {
		log.Error().Err(err).Msg("could not parse user app metadata")
		c.JSON(http.StatusInternalServerError, "could not parse user app metadata")
		return
	}

	// If the user app metadata has changed, set the refresh flag in the response
	if *appdata != *oldAppdata {
		c.JSON(http.StatusOK, api.Reply{Success: true, RefreshToken: true})
	} else {
		c.Status(http.StatusNoContent)
	}
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
