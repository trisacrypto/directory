package sentry

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Initialize the Sentry SDK with the given configuration. This should be called before
// any servers are started.
func Init(conf Config) (err error) {
	if err = sentry.Init(sentry.ClientOptions{
		Dsn:              conf.DSN,
		Environment:      conf.Environment,
		Release:          conf.GetRelease(),
		AttachStacktrace: true,
		Debug:            conf.Debug,
		TracesSampleRate: conf.SampleRate,
	}); err != nil {
		return fmt.Errorf("could not initialize sentry: %w", err)
	}

	log.Info().Bool("track_performance", conf.TrackPerformance).Float64("sample_rate", conf.SampleRate).Msg("sentry tracing is enabled")

	return nil
}

// Middleware that tracks request performance with Sentry.
func TrackPerformance() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := sentry.StartSpan(c.Request.Context(), c.Request.URL.Path)
		defer span.Finish()
		c.Next()
	}
}

// Flush the Sentry log, this is usually called before shutting down the servers.
func Flush(timeout time.Duration) bool {
	return sentry.Flush(timeout)
}
