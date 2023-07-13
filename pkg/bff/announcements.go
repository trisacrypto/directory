package bff

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/auth"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

const (
	maxAnnouncements = 10
	subMonths        = -2
)

// @Summary Get announcements [read:announcements]
// @Description Get the most recent network announcements
// @Tags announcements
// @Produce json
// @Success 200 {object} api.AnnouncementsReply
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /announcements [get]
func (s *Server) Announcements(c *gin.Context) {
	// Only fetch the previous 10 announcements from the last two months
	nbf := time.Now().AddDate(0, subMonths, 0)
	nbf = time.Date(nbf.Year(), nbf.Month(), 1, 0, 0, 0, 0, time.UTC)

	out, err := s.RecentAnnouncements(maxAnnouncements, nbf, time.Now())
	if err != nil {
		sentry.Error(c).Err(err).Msg("could not fetch recent announcements")
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

// @Summary Post an announcement [create:announcements]
// @Description Post a new announcement to the network
// @Tags announcements
// @Accept json
// @Produce json
// @Param announcement body models.Announcement true "Announcement to post"
// @Success 204
// @Failure 400 {object} api.Reply "Post date and author are required"
// @Failure 401 {object} api.Reply
// @Failure 500 {object} api.Reply
// @Router /announcements [post]
func (s *Server) MakeAnnouncement(c *gin.Context) {
	var (
		id     string
		err    error
		claims *auth.Claims
		post   *models.Announcement
	)

	if err = c.BindJSON(&post); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse announcement post data")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse announcement JSON data"))
		return
	}

	if post.PostDate != "" || post.Author != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("cannot set the post_date or author fields on the post"))
		return
	}

	if claims, err = auth.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not fetch claims from request")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not save announcement"))
		return
	}

	if claims.Email == "" {
		sentry.Warn(c).Msg("missing email on claims, cannot set author of network announcement")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("user claims are not correctly configured"))
		return
	}

	// Set the post date and the author
	post.PostDate = time.Now().Format("2006-01-02")
	post.Author = claims.Email

	if id, err = s.PostAnnouncement(post); err != nil {
		sentry.Error(c).Err(err).Msg("could not put announcement to trtl database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not save announcement"))
		return
	}

	// Return a 204 No Content to indicate the post happened successfully
	log.Info().Str("id", id).Str("title", post.Title).Str("author", post.Author).Msg("network announcement added")
	c.JSON(http.StatusNoContent, nil)
}

// RecentAnnouncements returns the set of results whose post date is after the not
// before timestamp, limited to the maximum number of results. Last updated returns the
// timestamp that any announcement was added or changed.
func (s *Server) RecentAnnouncements(maxResults int, notBefore, start time.Time) (out *api.AnnouncementsReply, err error) {
	// Do not allow unbounded requests in recent
	if notBefore.IsZero() {
		return nil, ErrUnboundedRecent
	}

	out = &api.AnnouncementsReply{
		Announcements: make([]*models.Announcement, 0, maxResults),
	}

	// Get the last day of this month or the start-after to begin querying announcements
	if start.IsZero() {
		start = time.Now()
	}
	month := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0).Add(-1 * time.Second)

	for !month.Before(notBefore) {
		ctx, cancel := utils.WithDeadline(context.Background())
		defer cancel()

		var crate *models.AnnouncementMonth
		if crate, err = s.db.RetrieveAnnouncementMonth(ctx, month.Format(models.MonthLayout)); err != nil {
			if errors.Is(err, storeerrors.ErrEntityNotFound) {
				// Decrement month and continue; see notes below about month decrement
				month = month.AddDate(0, -1, -5)
				month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0).Add(-1 * time.Second)
				continue
			}
			return nil, err
		}

		// Loop through the announcements adding them to the reply
		// NOTE: this expects announcements are in PostDate order
		for _, post := range crate.Announcements {
			// Stop if we've reached the maximum number of announcements
			if len(out.Announcements) >= maxResults {
				break
			}

			// Stop if the post is before the notBefore limit
			var pd time.Time
			if pd, err = post.ParsePostDate(); err != nil {
				return nil, fmt.Errorf("could not parse post date for announcement %s: %s", post.Id, err)
			}

			if pd.Before(notBefore) {
				break
			}

			out.Announcements = append(out.Announcements, post)
			out.LastUpdated = utils.Latest(out.LastUpdated, post.Modified)
		}

		// Decrement the month to check for more announcements
		// We need to find the last day of the previous month, to do this, we subtract
		// 1 month and 5 days (to ensure we go past the 28th of February), then compute
		// the last day of that month using the AddDate function.
		month = month.AddDate(0, -1, -5)
		month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0).Add(-1 * time.Second)
	}

	return out, nil
}

// Post an announcement, putting it to the database. This method does no
// verification of duplicate announcements or any content verification except for a
// check that an empty announcement is not being put to the database. Announcements are
// stored in announcement months, so the month for the announcement is extracted and the
// announcement is inserted into the correct month, creating it if necessary.
func (s *Server) PostAnnouncement(in *models.Announcement) (_ string, err error) {
	// Make sure we don't post empty announcements
	if in.Title == "" && in.Body == "" && in.PostDate == "" && in.Author == "" {
		return "", ErrEmptyAnnouncement
	}

	// Set the ID and timestamp metadata on the Post
	// Announcement keys are ksuids - timestamp ordered unique IDs so that it is easy
	// to scan the trtl database to find the most recent announcements.
	in.Id = ksuid.New().String()
	in.Created = time.Now().Format(time.RFC3339Nano)
	in.Modified = in.Created

	// Get the month to store the announcement in
	var month string
	if month, err = in.Month(); err != nil {
		return "", fmt.Errorf("could not identify month from post date: %s", err)
	}

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	// Get or Create the announcement month "crate"
	var crate *models.AnnouncementMonth
	if crate, err = s.db.RetrieveAnnouncementMonth(ctx, month); err != nil {
		switch {
		case errors.Is(err, storeerrors.ErrEntityNotFound):
			crate = &models.AnnouncementMonth{
				Date:          month,
				Announcements: make([]*models.Announcement, 0, 1),
			}

			// Update creates the announcement month if it does not exist in the database
			if err = s.db.UpdateAnnouncementMonth(ctx, crate); err != nil {
				return "", fmt.Errorf("could not create new announcement month: %s", err)
			}
		default:
			return "", fmt.Errorf("could not retrieve announcement month: %s", err)
		}
	}

	// Add the announcement in sorted order and update its modified timestamp
	crate.Add(in)
	crate.Modified = time.Now().Format(time.RFC3339Nano)

	if err = s.db.UpdateAnnouncementMonth(ctx, crate); err != nil {
		return "", err
	}
	return in.Id, nil
}
