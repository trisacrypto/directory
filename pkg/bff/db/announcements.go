package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	"google.golang.org/protobuf/proto"
)

const (
	NamespaceAnnouncements = "announcements"
)

// The Announcements type exposes methods for interacting with Network Announcements
// in the database. Announcements are stored as compact, compressed JSON to minimize the
// network requests and reduce storage requirements. This struct performs all necessary
// serialization on the announcements before storing and retrieving the model object. The
// announcement keys are ksuids - timestamp ordered unique IDs so that it is easy to scan
// the trtl database to find the most recent announcements.
//
// Announcements implements the Collection interface
type Announcements struct {
	db        *DB
	namespace string
}

// Ensure that Announcements implements the Collection interface.
var _ Collection = &Announcements{}

// Announcements constructs the collection type for interactions with network
// announcement objects. This method is intended to be used with chaining, e.g. as
// db.Announcements().Recent(), so to reduce the number of allocations a singleton
// intermediate struct is used. Method calls to the collection are thread-safe.
func (db *DB) Announcements() *Announcements {
	db.makeAnnouncements.Do(func() {
		db.announcements = &Announcements{
			db:        db,
			namespace: NamespaceAnnouncements,
		}
	})
	return db.announcements
}

// Recent returns the set of results whose post date is after the not before timestamp,
// limited to the maximum number of results. Last updated returns the timestamp that
// any announcement was added or changed.
func (a *Announcements) Recent(ctx context.Context, maxResults int, notBefore, start time.Time) (out *api.AnnouncementsReply, err error) {
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
		if crate, err = a.GetMonth(ctx, month.Format(models.MonthLayout)); err != nil {
			if errors.Is(err, ErrNotFound) {
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

// Post an announcement, putting it to the trtl database. This method does no
// verification of duplicate announcements or any content verification except for a
// check that an empty announcement is not being put to the database. Announcements are
// stored in announcement months, so the month for the announcement is extracted and the
// announcement is inserted into the correct month, creating it if necessary.
func (a *Announcements) Post(ctx context.Context, in *models.Announcement) (_ string, err error) {
	// Make sure we don't post empty announcements
	if in.Title == "" && in.Body == "" && in.PostDate == "" && in.Author == "" {
		return "", ErrEmptyAnnouncement
	}

	// Set the ID and timestamp metadata on the Post
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
	if crate, err = a.GetOrCreateMonth(ctx, month); err != nil {
		return "", err
	}

	// Add the announcement in sorted order and update its modified timestamp
	crate.Add(in)
	crate.Modified = time.Now().Format(time.RFC3339Nano)

	if err = a.SaveMonth(ctx, crate); err != nil {
		return "", err
	}
	return in.Id, nil
}

// GetOrCreateMonth from a month timestamp in the form YYYY-MM.
func (a *Announcements) GetOrCreateMonth(ctx context.Context, months string) (month *models.AnnouncementMonth, err error) {
	if month, err = a.GetMonth(ctx, months); err != nil {
		if errors.Is(err, ErrNotFound) {
			return a.CreateMonth(ctx, months)
		}
		return nil, err
	}
	return month, nil
}

// GetMonth from a month timestamp in the form YYYY-MM.
func (a *Announcements) GetMonth(ctx context.Context, months string) (month *models.AnnouncementMonth, err error) {
	// Get the key by creating an intermediate announcement month to ensure that
	// validation and key creation always happens the same way.
	var key, value []byte
	month = &models.AnnouncementMonth{Date: months}
	if key, err = month.Key(); err != nil {
		return nil, err
	}

	if value, err = a.db.Get(ctx, key, a.namespace); err != nil {
		return nil, err
	}

	if err = proto.Unmarshal(value, month); err != nil {
		return nil, err
	}
	return month, nil
}

// CreateMonth from a month timestamp in the form YYYY-MM.
func (a *Announcements) CreateMonth(ctx context.Context, months string) (month *models.AnnouncementMonth, err error) {
	month = &models.AnnouncementMonth{
		Date:          months,
		Announcements: make([]*models.Announcement, 0, 1),
	}

	if err = a.SaveMonth(ctx, month); err != nil {
		return nil, err
	}
	return month, nil
}

// Save a month, storing it in the database
func (a *Announcements) SaveMonth(ctx context.Context, month *models.AnnouncementMonth) (err error) {
	// Compute the key for the month
	var key []byte
	if key, err = month.Key(); err != nil {
		return err
	}

	// Update the modified timestamp and serialize
	month.Modified = time.Now().Format(time.RFC3339Nano)
	if month.Created == "" {
		month.Created = month.Modified
	}

	var value []byte
	if value, err = proto.Marshal(month); err != nil {
		return err
	}

	if err = a.db.Put(ctx, key, value, a.namespace); err != nil {
		return err
	}
	return nil
}

// Namespace implements the collection interface
func (a *Announcements) Namespace() string {
	return a.namespace
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
