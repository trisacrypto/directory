package bff

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"google.golang.org/protobuf/proto"
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
			if status, err := s.testnetGDS.Status(ctx, &gds.HealthCheck{}); err != nil {
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
			if status, err := s.mainnetGDS.Status(ctx, &gds.HealthCheck{}); err != nil {
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

// GetStatuses makes parallel calls to the directory service to get the status
// information for both testnet and mainnet.
func (s *Server) GetStatuses(ctx context.Context) (testnet, mainnet *gds.ServiceState, err error) {
	rpc := func(ctx context.Context, client GlobalDirectoryClient, network string) (rep proto.Message, err error) {
		return client.Status(ctx, &gds.HealthCheck{})
	}

	// Perform the parallel requests
	results, errs := s.ParallelGDSRequests(ctx, rpc, false)
	if len(errs) != 2 || len(results) != 2 {
		return nil, nil, fmt.Errorf("unexpected number of results from parallel requests: %d", len(results))
	}

	// Parse the results
	var ok bool
	if errs[0] != nil {
		log.Warn().Err(errs[0]).Msg("could not call testnet Status RPC")
		testnet = nil
	} else if results[0] == nil {
		return nil, nil, fmt.Errorf("nil testnet result returned from parallel requests")
	} else if testnet, ok = results[0].(*gds.ServiceState); !ok {
		return nil, nil, fmt.Errorf("unexpected testnet status result type returned from parallel requests: %T", results[0])
	}

	if errs[1] != nil {
		log.Warn().Err(errs[1]).Msg("could not call mainnet Status RPC")
		mainnet = nil
	} else if results[1] == nil {
		return nil, nil, fmt.Errorf("nil mainnet status result returned from parallel requests")
	} else if mainnet, ok = results[1].(*gds.ServiceState); !ok {
		return nil, nil, fmt.Errorf("unexpected mainnet status result type returned from parallel requests: %T", results[1])
	}

	return testnet, mainnet, nil
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
