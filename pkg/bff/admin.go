package bff

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
)

func (s *Server) ListCertificates(c *gin.Context) {
	// TODO: Call the GDS admin API to get the VASP certificates
	c.JSON(http.StatusNotImplemented, api.ErrorResponse(fmt.Errorf("not implemented")))
}
