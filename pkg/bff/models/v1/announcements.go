package models

import (
	"sort"
	"time"
)

const (
	PostDateLayout = "2006-01-02"
	MonthLayout    = "2006-01"
)

// Month returns the postdate month in the form YYYY-MM to determine which
// AnnouncementsMonth the announcement should belong in.
func (a *Announcement) Month() (_ string, err error) {
	var postDate time.Time
	if postDate, err = a.ParsePostDate(); err != nil {
		return "", err
	}
	return postDate.Format(MonthLayout), nil
}

// Return the timestamp from the post date.
func (a *Announcement) ParsePostDate() (time.Time, error) {
	return time.Parse(PostDateLayout, a.PostDate)
}

// Add an announcement ensuring that they are stored sorted by post date.
// NOTE: can sort postdate strings in the YYYY-MM-DD format without parsing them,
// however the post date must be validated before adding it to the month.
func (m *AnnouncementMonth) Add(a *Announcement) {
	i := sort.Search(len(m.Announcements), func(i int) bool {
		return m.Announcements[i].PostDate < a.PostDate
	})

	if i == len(m.Announcements) {
		m.Announcements = append(m.Announcements, a)
		return
	}

	m.Announcements = append(m.Announcements[:i+1], m.Announcements[i:]...)
	m.Announcements[i] = a
}

// Return the key associated with the announcement month: the byte array of the string
// date in YYYY-MM form. This method also validates the Date is correct.
func (m *AnnouncementMonth) Key() (_ []byte, err error) {
	// Validate the announcement month is correct
	if _, err = time.Parse(MonthLayout, m.Date); err != nil {
		return nil, err
	}
	return []byte(m.Date), nil
}
