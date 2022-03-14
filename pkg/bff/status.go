package bff

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
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

	// Check if the user has supplied a nogds query param
	params := &api.StatusParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Warn().Err(err).Msg("could not bind request with query params")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Get the status of the TestNet and MainNet GDS servers
	if !params.NoGDS {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			if status, err := s.testnet.Status(ctx, &gds.HealthCheck{}); err != nil {
				log.Warn().Err(err).Str("network", testnet).Msg("could not connect to GDS")
				out.TestNet = "unavailable"
			} else {
				log.Debug().Str("network", testnet).
					Str("status", status.Status.String()).
					Str("not_before", status.NotBefore).
					Str("not_after", status.NotAfter).
					Msg("GDS status")
				out.TestNet = strings.ToLower(status.Status.String())
			}
		}()

		go func() {
			defer wg.Done()
			if status, err := s.mainnet.Status(ctx, &gds.HealthCheck{}); err != nil {
				log.Warn().Err(err).Str("network", mainnet).Msg("could not connect to GDS")
				out.MainNet = "unavailable"
			} else {
				log.Debug().Str("network", mainnet).
					Str("status", status.Status.String()).
					Str("not_before", status.NotBefore).
					Str("not_after", status.NotAfter).
					Msg("GDS status")
				out.MainNet = strings.ToLower(status.Status.String())
			}

		}()

		wg.Wait()
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
