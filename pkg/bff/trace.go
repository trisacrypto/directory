package bff

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Middleware that tracks request performance with Sentry.
func (s *Server) TrackPerformance() gin.HandlerFunc {
	log.Info().Float64("sample_rate", s.conf.Sentry.SampleRate).Msg("sentry performance tracking enabled")
	return func(c *gin.Context) {
		span := sentry.StartSpan(c.Request.Context(), c.Request.URL.Path)
		defer span.Finish()
		c.Next()
	}
}
