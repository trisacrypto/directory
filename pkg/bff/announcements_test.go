package bff_test

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/trisacrypto/directory/pkg/bff/auth/authtest"
	"github.com/trisacrypto/directory/pkg/bff/db"
	records "github.com/trisacrypto/directory/pkg/bff/db/models/v1"
	"google.golang.org/protobuf/proto"
)

func (s *bffTestSuite) TestAnnouncements() {
	require := s.Require()

	// Keep track of what announcement months are being created to clean up at the end
	months := make(map[string]struct{})
	defer func() {
		for month := range months {
			err := s.db.Delete(context.TODO(), []byte(month), db.NamespaceAnnouncements)
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

		_, err = s.db.Announcements().Post(context.TODO(), post)
		require.NoError(err, "could not post an announcement fixture")
	}

	// Create a post for yesterday to ensure there is at least one post returned
	months[time.Now().AddDate(0, 0, -1).Format(records.MonthLayout)] = struct{}{}
	_, err = s.db.Announcements().Post(context.TODO(), &records.Announcement{
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
			s.db.Delete(context.TODO(), []byte(month), db.NamespaceAnnouncements)
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
	monthData, err := s.db.Get(context.TODO(), []byte(months[0]), db.NamespaceAnnouncements)
	require.NoError(err, "could not get announcements container")
	require.NotEmpty(monthData, "expected month date to be populated")

	month := &records.AnnouncementMonth{}
	require.NoError(proto.Unmarshal(monthData, month), "could not unmarshal announcement month")

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
