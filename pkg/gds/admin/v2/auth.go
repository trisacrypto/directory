package admin

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	bearer        = "Bearer "
	authorization = "Authorization"
)

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
