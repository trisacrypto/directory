package sentry

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
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

// Gin middleware that tracks HTTP request performance with Sentry.
func TrackPerformance(tags map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		span := sentry.StartSpan(c.Request.Context(), "http handler", sentry.TransactionName(request))
		for k, v := range tags {
			span.SetTag(k, v)
		}
		defer span.Finish()
		c.Next()
	}
}

// Gin middleware that adds request-level tags to the current Sentry scope. This
// also accepts a map of service-level tags to uniquely identify properties of the http
// service for monolithic server setups.
func UseTags(tags map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			// Service-level tags
			for k, v := range tags {
				hub.Scope().SetTag(k, v)
			}

			// Request-level tags
			hub.Scope().SetTag("path", c.Request.URL.Path)
			hub.Scope().SetTag("method", c.Request.Method)
		}
		c.Next()
	}
}

// Flush the Sentry log, this is usually called before shutting down the servers.
func Flush(timeout time.Duration) bool {
	return sentry.Flush(timeout)
}
