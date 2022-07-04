package db_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/directory/pkg/bff/db"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	"google.golang.org/protobuf/proto"
)

func (s *dbTestSuite) TestAnnouncements() {
	// Test creating and retrieving announcements
	var err error
	require := s.Require()

	// Announcements should implement the Collection interface
	require.Equal(NamespaceAnnouncements, s.db.Announcements().Namespace())

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// There should be nothing in the database at the start of the test.
	recent, err := s.db.Announcements().Recent(ctx, 10000, time.Time{})
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 0, "expected no announcements returned")
	require.Empty(recent.LastUpdated, "expected last updated to be zero-valued")

	// Load announcements fixtures
	fixture, err := loadAnnouncements()
	require.NoError(err, "could not load announcement fixtures")
	require.Len(fixture, 10, "fixtures have changed without test update")

	// Post each fixture one at a time and ensure that we can fetch them all
	ids := make([]string, 0, 10)
	for i, post := range fixture {
		id, err := s.db.Announcements().Post(ctx, post)
		require.NoError(err, "could not post announcement %d", i)
		require.NotEmpty(id, "expected an ordered kid to be assigned")
		ids = append(ids, id)

		// Ensure there is at least 1 second between each post or the ksuid will
		// identify the posts as concurrent and produce an incorrect ordering.
		// Normally posts won't be published within the same second, but this is a bad
		// situation for tests because it means that this test will take a min of 10
		// seconds to run, which adds up in CI and beyond ...
		time.Sleep(1 * time.Second)
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
	recent, err = s.db.Announcements().Recent(ctx, 10000, time.Time{})
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 10, "expected 10 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	// NOTE: we will get a random order of our posts if we don't sleep 1 second between Post
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 10-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}

	// Should be able to limit the number of results returned
	recent, err = s.db.Announcements().Recent(ctx, 5, time.Time{})
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 5, "expected 5 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 10-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}

	// Should be able to set the not before timestamp
	nbf, _ := time.Parse("2006-01-02", "2022-03-31")
	recent, err = s.db.Announcements().Recent(ctx, 10000, nbf)
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 3, "expected 3 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 10-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}

	// Should be able to set the not before timestamp AND max results
	recent, err = s.db.Announcements().Recent(ctx, 2, nbf)
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 2, "expected 2 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 10-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}
}

func TestAnnouncementsSerialization(t *testing.T) {
	// Should be able to encode and decode an announcement to and from bytes.
	announcement := &models.Announcement{
		Title:    "It is the Solstice!",
		Body:     "Today is the longest day of the year, make sure it's a grilling day!",
		PostDate: "2022-06-21",
		Author:   "summer@solstice.ninja",
	}

	// Technically an empty collection should be able to encode and decode the model.
	// The only reason these methods are attached to the collection is for namespacing.
	collection := &Announcements{}
	data, err := collection.Encode(announcement)
	require.NoError(t, err, "could not encode announcement into bytes")
	require.NotEmpty(t, data, "expected some data returned from encoding")

	// Decode the announcement
	decoded, err := collection.Decode(data)
	require.NoError(t, err, "could not decode announcement from bytes")
	require.True(t, proto.Equal(announcement, decoded), "the decoded announcement did not match the original")
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
