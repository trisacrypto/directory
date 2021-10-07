package utils

import (
	"errors"
	"time"
)

type Week struct {
	Date    time.Time
	Year    int
	WeekNum int
}

type WeekIterator struct {
	Start   *Week
	Current *Week
	End     *Week
}

// Calculate the actual date which represents the start of the given ISOWeek
func DateFromISOWeek(isoYear int, isoWeek int) (date time.Time) {
	// Start at the beginning of the target year, iterate backwards to Monday
	date = time.Date(isoYear, 0, 0, 0, 0, 0, 0, time.UTC)
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}

	// Iterate forwards to the first week of the target year, this needs to be done
	// if the target year doesn't start on a Monday
	year, week := date.ISOWeek()
	for year < isoYear {
		date = date.AddDate(0, 0, 7)
		year, week = date.ISOWeek()
	}

	// Iterate forwards to the correct start week
	for week < isoWeek {
		date = date.AddDate(0, 0, 7)
		_, week = date.ISOWeek()
	}

	return date
}

// Returns a new Week object given a time.Time
func NewWeek(date time.Time) (week *Week) {
	week = &Week{}

	// Align to the nearest ISOWeek
	week.Year, week.WeekNum = date.ISOWeek()

	// Get the start date of the ISOWeek
	week.Date = DateFromISOWeek(week.Year, week.WeekNum)

	return week
}

// Returns the difference w-u between two Weeks, in terms of number of weeks
func (w *Week) Sub(u *Week) int {
	if w.Year >= u.Year {
		return ((w.Year - u.Year) * 52) + w.WeekNum - u.WeekNum
	} else {
		return ((u.Year - w.Year) * -52) - (u.WeekNum - w.WeekNum)
	}
}

// Takes in a time range and returns a WeekIterator
func GetWeekIterator(start time.Time, end time.Time) (iter *WeekIterator, err error) {
	if start.After(end) {
		return nil, errors.New("start time cannot be after end time")
	}
	return &WeekIterator{Start: NewWeek(start), End: NewWeek(end)}, nil
}

// Returns the next Week from the WeekIterator
func (i *WeekIterator) Next() (week *Week, ok bool) {
	if i.Current == nil {
		week = i.Start
		i.Current = i.Start
		return week, true
	}

	i.Current.Date = i.Current.Date.AddDate(0, 0, 7)
	if i.Current.Date.After(i.End.Date) {
		return nil, false
	}

	week = &Week{Date: i.Current.Date}
	week.Year, week.WeekNum = week.Date.ISOWeek()
	return week, true
}
