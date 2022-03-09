package bff

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
)

const (
	serverStatusOK          = "ok"
	serverStatusStopping    = "stopping"
	serverStatusMaintenance = "maintenance"
)

func (s *Server) Status(c *gin.Context) {
	// The available middleware handles stopping and maintenance mode. If the request
	// has come this far, the status is necessarily ok.
	out := api.StatusReply{
		Status:  serverStatusOK,
		Version: pkg.Version(),
		Uptime:  time.Since(s.started).String(),
	}

	// Render the response
	c.JSON(http.StatusOK, out)
}

// Available is middleware that uses the healthy boolean to return a service unavailable
// http status code if the server is shutting down. This middleware must be first in the
// chain to ensure that complex handling to slow the shutdown of the server.
func (s *Server) Available() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check health status
		s.RLock()
		if !s.healthy {
			var status string
			if s.conf.Maintenance {
				status = serverStatusMaintenance
			} else {
				status = serverStatusStopping
			}

			c.JSON(http.StatusServiceUnavailable, api.StatusReply{
				Status:  status,
				Uptime:  time.Since(s.started).String(),
				Version: pkg.Version(),
			})

			c.Abort()
			s.RUnlock()
			return
		}
		s.RUnlock()
		c.Next()
	}
}
