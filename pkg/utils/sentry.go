package utils

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// Middleware that tracks request performance with Sentry.
func SentryTrackPerformance() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := sentry.StartSpan(c.Request.Context(), c.Request.URL.Path)
		defer span.Finish()
		c.Next()
	}
}
