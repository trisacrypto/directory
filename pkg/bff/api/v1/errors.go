package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff/config"
)

var (
	unsuccessful = Reply{Success: false}
	notFound     = Reply{Success: false, Error: "resource not found"}
	notAllowed   = Reply{Success: false, Error: "method not allowed"}
)

var (
	ErrNetworkRequired    = fmt.Errorf("request requires a valid network (%s or %s)", config.TestNet, config.MainNet)
	ErrInvalidCredentials = errors.New("auth0 credentials are missing or invalid")
	ErrExpiredCredentials = errors.New("auth0 credentials have expired")
	ErrPathRequired       = errors.New("local credentials requires a path to the stored json credential")
	ErrIDRequired         = errors.New("request requires a valid ID to determine the endpoint path")
	ErrMissingMemberID    = errors.New("member ID is required to query member details")
)

// ErrorResponse constructs an new response from the error or returns a success: false.
func ErrorResponse(err interface{}) Reply {
	if err == nil {
		return unsuccessful
	}

	rep := Reply{Success: false}
	switch err := err.(type) {
	case error:
		rep.Error = err.Error()
	case string:
		rep.Error = err
	case fmt.Stringer:
		rep.Error = err.String()
	case json.Marshaler:
		data, e := err.MarshalJSON()
		if e != nil {
			panic(err)
		}
		rep.Error = string(data)
	default:
		rep.Error = "unhandled error response"
	}

	return rep
}

// NotFound returns a JSON 404 response for the API.
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, notFound)
}

// NotAllowed returns a JSON 405 response for the API.
func NotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, notAllowed)
}

// MustRefreshToken returns a JSON 401 response with the refresh_token flag set to true.
func MustRefreshToken(c *gin.Context, err interface{}) {
	rep := ErrorResponse(err)
	rep.RefreshToken = true
	c.JSON(http.StatusUnauthorized, rep)
}
