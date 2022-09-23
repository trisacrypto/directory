package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
)

// RecentAnnouncements returns the set of results whose post date is after the not
// before timestamp, limited to the maximum number of results. Last updated returns the
// timestamp that any announcement was added or changed.
func (store *DB) RecentAnnouncements(maxResults int, notBefore, start time.Time) (out *api.AnnouncementsReply, err error) {
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
		var crate *models.AnnouncementMonth
		if crate, err = store.RetrieveAnnouncementMonth(month.Format(models.MonthLayout)); err != nil {
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
			out.LastUpdated = Latest(out.LastUpdated, post.Modified)
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
func (store *DB) PostAnnouncement(in *models.Announcement) (_ string, err error) {
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

	// Get or Create the announcement month "crate"
	var crate *models.AnnouncementMonth
	if crate, err = store.GetOrCreateMonth(month); err != nil {
		return "", err
	}

	// Add the announcement in sorted order and update its modified timestamp
	crate.Add(in)
	crate.Modified = time.Now().Format(time.RFC3339Nano)

	if err = store.db.UpdateAnnouncementMonth(crate); err != nil {
		return "", err
	}
	return in.Id, nil
}

// GetOrCreateMonth from a month timestamp in the form YYYY-MM.
func (store *DB) GetOrCreateMonth(date string) (month *models.AnnouncementMonth, err error) {
	if month, err = store.db.RetrieveAnnouncementMonth(date); err != nil {
		if errors.Is(err, storeerrors.ErrEntityNotFound) {
			return store.CreateMonth(date)
		}
		return nil, err
	}
	return month, nil
}

// CreateMonth from a month timestamp in the form YYYY-MM.
func (store *DB) CreateMonth(date string) (month *models.AnnouncementMonth, err error) {
	month = &models.AnnouncementMonth{
		Date:          date,
		Announcements: make([]*models.Announcement, 0, 1),
	}

	if err = store.db.UpdateAnnouncementMonth(month); err != nil {
		return nil, err
	}
	return month, nil
}

// Retrieve an announcement month from the database by month timestamp in the form YYYY-MM.
func (store *DB) RetrieveAnnouncementMonth(date string) (month *models.AnnouncementMonth, err error) {
	return store.db.RetrieveAnnouncementMonth(date)
}

// Delete an announcement month from the database by month timestamp in the form YYYY-MM.
func (store *DB) DeleteAnnouncementMonth(date string) (err error) {
	return store.db.DeleteAnnouncementMonth(date)
}

// Helper method to return the latest string timestamp from the two RFC3339 timestamps
func Latest(a, b string) string {
	// Parse without checking errors - will use zero-valued ts for checks
	ats, _ := time.Parse(time.RFC3339Nano, a)
	bts, _ := time.Parse(time.RFC3339Nano, b)

	switch {
	case ats.IsZero() && bts.IsZero():
		return ""
	case !ats.IsZero() && ats.After(bts):
		return a
	case !bts.IsZero() && bts.After(ats):
		return b
	case !ats.IsZero() && !bts.IsZero() && ats.Equal(bts):
		return a
	}
	return ""
}
