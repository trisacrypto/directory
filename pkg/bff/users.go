package bff

import "github.com/gin-gonic/gin"

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
	// Fetch the user from the context
}
