package bff

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Overview endpoint is an authenticated endpoint that requires the read:vasp permission.
// TODO: implement; this is just a mock endpoint to facilitate front-end development.
func (s *Server) Overview(c *gin.Context) {
	data := map[string]interface{}{
		"testnet": map[string]interface{}{
			"status":              "healthy",
			"vasps_count":         51,
			"certificates_issued": 49,
			"new_members":         2,
		},
		"mainnet": map[string]interface{}{
			"status":              "healthy",
			"vasps_count":         14,
			"certificates_issued": 10,
			"new_members":         0,
		},
		"organization": map[string]interface{}{
			"vasp_id":             "0ffe693f-9752-45cb-830f-7dbae12d6baf",
			"verification_status": "VERIFIED",
			"country":             "US",
			"certificate": map[string]interface{}{
				"common_name":     "trisa.rotational.io",
				"alternate_names": []string{},
				"issued_on":       "2021-06-30T17:16:14Z",
				"expires":         "2022-07-01T17:16:14Z",
			},
		},
	}

	c.JSON(http.StatusOK, data)
}
