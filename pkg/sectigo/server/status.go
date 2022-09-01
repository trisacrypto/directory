package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg"
)

// Available is middleware that checks the healthy boolean and returns service
// unavailable if the server is shutting down.
func (s *Server) Available() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.RLock()
		if !s.healthy {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "stopping",
				"version": pkg.Version(),
				"uptime":  time.Since(s.started).String(),
			})

			c.Abort()
			s.RUnlock()
			return
		}

		s.RUnlock()
		c.Next()
	}
}

// Status implements the heartbeat endpoint of the API server.
func (s *Server) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": pkg.Version(),
		"uptime":  time.Since(s.started).String(),
	})
}

// NotFound returns a JSON 404 response for the API
func (s *Server) NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   "resource not found",
	})
}

// NotAllowed returns a JSON 405 response for the API.
func (s *Server) NotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"success": false,
		"error":   "method not allowed",
	})
}
