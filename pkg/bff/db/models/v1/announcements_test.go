package models_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/db/models/v1"
)

func TestAnnouncements(t *testing.T) {
	announcement := &models.Announcement{
		Title:    "An announcement",
		PostDate: "2019-07-21",
	}

	month, err := announcement.Month()
	require.NoError(t, err, "expected the post date to be parseable")
	require.Equal(t, "2019-07", month, "expected the month to be YYYY-MM extracted from announcement")

	announcement.PostDate = "07/21/2019"
	month, err = announcement.Month()
	require.Error(t, err, "expected invalid post date to return an error")
	require.Empty(t, month, "expected month to be empty when post date is invalid")
}

func TestAnnouncementsMonth(t *testing.T) {
	month := &models.AnnouncementMonth{
		Date: "2019-07",
		Announcements: []*models.Announcement{
			{
				Title:    "An announcement",
				PostDate: "2019-07-21",
			},
		},
	}

	key, err := month.Key()
	require.NoError(t, err, "could not fetch key from valid month date")
	require.Equal(t, []byte("2019-07"), key, "unexpected key returned")

	month.Date = "07-2019"
	key, err = month.Key()
	require.Error(t, err, "expected error when date is invalid")
	require.Empty(t, key, "expected nil key when date is invalid")

	// Test adding announcements with random dates
	month.Date = "2019-07"
	for i := 0; i < 100; i++ {
		day := rand.Intn(30) + 1
		month.Add(&models.Announcement{
			Title:    fmt.Sprintf("The number is %04X", day),
			PostDate: fmt.Sprintf("%s-%02d", month.Date, day),
		})
	}

	// All post dates should be in sorted order
	day, err := time.Parse(models.PostDateLayout, "2019-08-01")
	require.NoError(t, err, "could not parse starting day with post date layout")

	for i, a := range month.Announcements {
		pd, err := a.ParsePostDate()
		require.NoError(t, err, "could not parse post date %d", i)

		if !pd.Equal(day) {
			require.True(t, pd.Before(day), "post date %d not equal to or before previous day", i)
		} else {
			require.True(t, day.Equal(pd), "post date %d not equal to or before previous day", i)
		}

		day = pd
	}
}
