package bff

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
)

// Certificates returns the list of certificates for the authenticated user.
func (s *Server) Certificates(c *gin.Context) {
	// TODO: Call the GDS admin API to get the VASP certificates
	c.JSON(http.StatusNotImplemented, api.ErrorResponse(fmt.Errorf("not implemented")))
}
