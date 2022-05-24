package admin

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/tokens"
)

const (
	bearer        = "Bearer "
	authorization = "Authorization"
	UserClaims    = "user_claims"
)

// Authorization middleware ensures that the request has a valid Bearer JWT in the
// Authorization header of the request otherwise it returns a 401 unauthorized error.
func Authorization(tm *tokens.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err         error
			accessToken string
			claims      *tokens.Claims
		)

		// Get access token from the request and ensure it is supplied
		if accessToken, err = GetAccessToken(c); err != nil {
			log.Debug().Err(err).Msg("no access token requested with secure endpoint")
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse("a valid authorization is required to access this endpoint"))
			return
		}

		// Verify that the access_token is valid and signed with server keys
		if claims, err = tm.Verify(accessToken); err != nil {
			log.Debug().Err(err).Msg("access token is invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse("a valid authorization is required to access this endpoint"))
			return
		}

		// Add claims to context for use in downstream processing
		c.Set(UserClaims, claims)

		// Continue with middleware handling
		c.Next()
	}
}

// GetAccessToken retrieves the bearer token from the authorization header and parses
// it to return only the access token. If the header is missing or the token is not
// available an error is returned.
func GetAccessToken(c *gin.Context) (tks string, err error) {
	header := c.GetHeader(authorization)
	if header != "" {
		parts := strings.Split(header, bearer)
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1]), nil
		}
		return "", errors.New("could not parser Bearer token from Authorization header")
	}
	return "", errors.New("no access token found in request")
}
