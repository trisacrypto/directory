package bff_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	records "github.com/trisacrypto/directory/pkg/bff/models/v1"
)

func (s *bffTestSuite) TestAnnouncements() {
	require := s.Require()

	// Keep track of what announcement months are being created to clean up at the end
	months := make(map[string]struct{})
	defer func() {
		for month := range months {
			err := s.db.DeleteAnnouncementMonth(month)
			require.NoError(err, "could not cleanup announcements")
		}
	}()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := s.client.Announcements(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	_, err = s.client.Announcements(context.TODO())
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Set valid credentials for the remainder of the tests
	claims.Permissions = []string{"read:vasp"}
	require.NoError(s.SetClientCredentials(claims), "could not create token from valid credentials")

	// Should be able to return empty results even when nothing is in the database
	posts, err := s.client.Announcements(context.TODO())
	require.NoError(err, "was unable to fetch announcements with valid claims")
	require.Len(posts.Announcements, 0, "expected no announcements returned")
	require.NotEmpty(posts.LastUpdated, "expected last updated to be set to now timestamp")

	// Add some announcements to the database for the past several months
	now := time.Now()
	for i := 0; i < 20; i++ {
		post := &records.Announcement{
			Title:  fmt.Sprintf("test post %d", i+1),
			Body:   fmt.Sprintf("this is a test post number %d", i+1),
			Author: "test@example.com",
		}

		// Create a random post date sometime in the past several months
		pd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, -1*rand.Intn(5), -1*rand.Intn(32))
		post.PostDate = pd.Format(records.PostDateLayout)
		months[pd.Format(records.MonthLayout)] = struct{}{}

		_, err = s.bff.PostAnnouncement(post)
		require.NoError(err, "could not post an announcement fixture")
	}

	// Create a post for yesterday to ensure there is at least one post returned
	months[time.Now().AddDate(0, 0, -1).Format(records.MonthLayout)] = struct{}{}
	_, err = s.bff.PostAnnouncement(&records.Announcement{
		Title:    "from the future",
		Body:     "this was posted yesterday",
		Author:   "future@example.com",
		PostDate: time.Now().AddDate(0, 0, -1).Format(records.PostDateLayout),
	})
	require.NoError(err, "could not post an announcement fixture")

	posts, err = s.client.Announcements(context.TODO())
	require.NoError(err, "could not fetch announcements with valid claims and data in the database")
	require.LessOrEqual(10, len(posts.Announcements), "expected 10 or less announcements")
	require.NotEmpty(posts.Announcements, "expected at least one result returned")
	require.NotEmpty(posts.LastUpdated, "expected last updated to be set")
}

func (s *bffTestSuite) TestMakeAnnouncement() {
	require := s.Require()

	// Keep track of what announcement months are being created to clean up at the end
	months := make([]string, 0, 1)
	months = append(months, time.Now().Format("2006-01"))

	defer func() {
		for _, month := range months {
			s.db.DeleteAnnouncementMonth(month)
		}
	}()

	// Create initial claims fixture
	claims := &authtest.Claims{
		Email:       "leopold.wentzel@gmail.com",
		Permissions: []string{"read:nothing"},
	}

	post := &records.Announcement{
		Title: "Hear ye, Hear ye",
		Body:  "We are conducting tests of the make announcements endpoint.",
	}

	// Endpoint requires CSRF protection
	err := s.client.MakeAnnouncement(context.TODO(), post)
	s.requireError(err, http.StatusForbidden, "csrf verification failed for request", "expected error when request is not CSRF protected")
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection on client")

	// Endpoint must be authenticated
	err = s.client.MakeAnnouncement(context.TODO(), post)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Endpoint requires the read:vasp permission
	require.NoError(s.SetClientCredentials(claims), "could not create token with incorrect permissions")
	err = s.client.MakeAnnouncement(context.TODO(), post)
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation", "expected error when user is not authorized")

	// Set valid credentials for the remainder of the tests
	claims.Permissions = []string{"create:announcements"}
	require.NoError(s.SetClientCredentials(claims), "could not create token from valid credentials")

	// Should be able to make an announcement
	err = s.client.MakeAnnouncement(context.TODO(), post)
	require.NoError(err, "was not able to make an announcement")

	// Check that the announcement exists in the database
	month, err := s.db.RetrieveAnnouncementMonth(months[0])
	require.NoError(err, "could not get announcements container")
	require.NotEmpty(month.Date, "expected month date to be set")
	require.Len(month.Announcements, 1, "expected announcements to contain 1 item")
	require.NotEmpty(month.Created, "expected created timestamp set")
	require.NotEmpty(month.Modified, "expected modified timestamp set")

	compat := month.Announcements[0]
	require.NotEmpty(compat.Id, "expected announcement ID to be set")
	require.Equal(post.Title, compat.Title, "expected announcement title to be same as orginal")
	require.Equal(post.Body, compat.Body, "expected announcement body to be same as orginal")
	require.Equal(time.Now().Format("2006-01-02"), compat.PostDate, "expected announcement post date to be set")
	require.Equal(claims.Email, compat.Author, "expected author to be set from claims")
	require.NotEmpty(compat.Created, "expected created timestamp set")
	require.NotEmpty(compat.Modified, "expected modified timestamp set")

	// Test Invalid Posts
	// Post should not have post_date set
	post.PostDate = "2022-07-04"
	err = s.client.MakeAnnouncement(context.TODO(), post)
	s.requireError(err, http.StatusBadRequest, "cannot set the post_date or author fields on the post", "expected post date required empty")

	// Post should not have author set
	post.PostDate = ""
	post.Author = "James Jillian"
	err = s.client.MakeAnnouncement(context.TODO(), post)
	s.requireError(err, http.StatusBadRequest, "cannot set the post_date or author fields on the post", "expected post date required empty")

	// Require email in claims to make announcement
	post.Author = ""
	claims.Email = ""
	require.NoError(s.SetClientCredentials(claims), "could not create token from valid credentials without email")
	err = s.client.MakeAnnouncement(context.TODO(), post)
	s.requireError(err, http.StatusBadRequest, "user claims are not correctly configured", "expected post date required empty")
}

func (s *bffTestSuite) TestAnnouncementsHelpers() {
	// Test creating and retrieving announcements using the helper methods
	var err error
	require := s.Require()

	// Keep track of what announcement months are being created to clean up at the end
	months := make(map[string]struct{})
	defer func() {
		for month := range months {
			err := s.db.DeleteAnnouncementMonth(month)
			require.NoError(err, "could not cleanup announcements")
		}
	}()

	// Should not be able to request unbounded time
	_, err = s.bff.RecentAnnouncements(10000, time.Time{}, time.Time{})
	require.ErrorIs(err, bff.ErrUnboundedRecent, "expected error when zero-valued time passed in as not before")
	nbf, _ := time.Parse("2006-01-02", "2020-12-01")
	stt, _ := time.Parse("2006-01-02", "2023-12-31")

	// There should be nothing in the database at the start of the test.
	recent, err := s.bff.RecentAnnouncements(10000, nbf, stt)
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
		pd, err := time.Parse("2006-01-02", post.PostDate)
		require.NoError(err, "could not parse post date from fixture")
		months[pd.Format(records.MonthLayout)] = struct{}{}

		id, err := s.bff.PostAnnouncement(post)
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
	recent, err = s.bff.RecentAnnouncements(10000, nbf, stt)
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
	recent, err = s.bff.RecentAnnouncements(5, nbf, stt)
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
	recent, err = s.bff.RecentAnnouncements(10000, nbf, stt)
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 5, "expected 5 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 11-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}

	// Should be able to set the not before timestamp AND max results
	recent, err = s.bff.RecentAnnouncements(2, nbf, stt)
	require.NoError(err, "could not fetch recent announcements")
	require.Len(recent.Announcements, 2, "expected 2 announcements returned")
	require.NotEmpty(recent.LastUpdated, "expected last updated to be set")

	// The Posts should be returned in the expected order
	for i, post := range recent.Announcements {
		expected := fmt.Sprintf("Post %d", 11-i)
		require.Equal(expected, post.Title, "posts seem to be out of order")
	}
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
