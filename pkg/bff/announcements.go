package bff

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

const (
	maxAnnouncements = 10
	subMonths        = -2
)

func (s *Server) Announcements(c *gin.Context) {
	// Only fetch the previous 10 announcements from the last two months
	nbf := time.Now().AddDate(0, subMonths, 0)
	nbf = time.Date(nbf.Year(), nbf.Month(), 1, 0, 0, 0, 0, time.UTC)

	out, err := s.db.Announcements().Recent(c.Request.Context(), maxAnnouncements, nbf, time.Now())
	if err != nil {
		log.Error().Err(err).Msg("could not fetch recent announcements")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("unable to fetch recent announcements"))
		return
	}

	// If the database is empty, set the last updated timestamp to now.
	if len(out.Announcements) == 0 && out.LastUpdated == "" {
		out.LastUpdated = time.Now().Format(time.RFC3339)
	}

	// Return the results
	c.JSON(http.StatusOK, out)
}

func (s *Server) MakeAnnouncement(c *gin.Context) {
	var (
		id     string
		err    error
		claims *auth.Claims
		post   *models.Announcement
	)

	if err = c.BindJSON(&post); err != nil {
		log.Warn().Err(err).Msg("could not parse announcement post data")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse announcement JSON data"))
		return
	}

	if post.PostDate != "" || post.Author != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("cannot set the post_date or author fields on the post"))
		return
	}

	if claims, err = auth.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not fetch claims from request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not save announcement"))
		return
	}

	if claims.Email == "" {
		log.Warn().Msg("missing email on claims, cannot set author of network announcement")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("user claims are not correctly configured"))
		return
	}

	// Set the post date and the author
	post.PostDate = time.Now().Format("2006-01-02")
	post.Author = claims.Email

	if id, err = s.db.Announcements().Post(c.Request.Context(), post); err != nil {
		log.Error().Err(err).Msg("could not put announcement to trtl database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not save announcement"))
		return
	}

	// Return a 204 No Content to indicate the post happened successfully
	log.Info().Str("id", id).Str("title", post.Title).Str("author", post.Author).Msg("network announcement added")
	c.JSON(http.StatusNoContent, nil)
}
