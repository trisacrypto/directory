package db_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff/db"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

func (s *dbTestSuite) TestAnnouncements() {
	// Test creating and retrieving announcements
	var err error
	require := s.Require()

	// Should not be able to request unbounded time
	_, err = s.db.RecentAnnouncements(10000, time.Time{}, time.Time{})
	require.ErrorIs(err, ErrUnboundedRecent, "expected error when zero-valued time passed in as not before")
	nbf, _ := time.Parse("2006-01-02", "2020-12-01")
	stt, _ := time.Parse("2006-01-02", "2023-12-31")

	// There should be nothing in the database at the start of the test.
	recent, err := s.db.RecentAnnouncements(10000, nbf, stt)
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 0, "expected no announcements returned")
	require.Empty(recent.LastUpdated, "expected last updated to be zero-valued")

	// Load announcements fixtures
	fixture, err := loadAnnouncements()
	require.NoError(err, "could not load announcement fixtures")
	require.Len(fixture, 11, "fixtures have changed without test update")

	// Post each fixture one at a time and ensure that we can fetch them all
	ids := make([]string, 0, 10)
	for i, post := range fixture {
		id, err := s.db.PostAnnouncement(post)
		require.NoError(err, "could not post announcement %d", i)
		require.NotEmpty(id, "expected an ordered kid to be assigned")
		ids = append(ids, id)
	}

	// Ensure all IDs assigned are unique
	for i, id := range ids {
		for j, jd := range ids {
			if i == j {
				continue
			}
			require.NotEqual(id, jd, "expected assigned IDs to be unique in database")
		}
	}

	// Should be able to retrieve all recent announcements
	recent, err = s.db.RecentAnnouncements(10000, nbf, stt)
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 11, "expected 11 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	// NOTE: Post titles must be ordered by post date
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 11-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}

	// Should be able to limit the number of results returned
	recent, err = s.db.RecentAnnouncements(5, nbf, stt)
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 5, "expected 5 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 11-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}

	// Should be able to set the not before timestamp
	nbf, _ = time.Parse("2006-01-02", "2022-03-18")
	recent, err = s.db.RecentAnnouncements(10000, nbf, stt)
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 5, "expected 5 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 11-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}

	// Should be able to set the not before timestamp AND max results
	recent, err = s.db.RecentAnnouncements(2, nbf, stt)
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 2, "expected 2 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 11-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}
}

func TestLatest(t *testing.T) {
	alpha := "2022-04-07T20:04:21.000Z"
	bravo := "2022-04-07T08:36:07.000Z"

	require.Equal(t, alpha, Latest(alpha, ""))
	require.Equal(t, bravo, Latest("", bravo))
	require.Empty(t, Latest("", ""))
	require.Equal(t, alpha, Latest(alpha, bravo))
}

func loadAnnouncements() (fixture []*models.Announcement, err error) {
	var f *os.File
	if f, err = os.Open("testdata/announcements.json"); err != nil {
		return nil, err
	}
	defer f.Close()

	fixture = make([]*models.Announcement, 0, 10)
	if err = json.NewDecoder(f).Decode(&fixture); err != nil {
		return nil, err
	}
	return fixture, nil
}
