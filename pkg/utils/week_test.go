package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils"
)

func TestDateFromISOWeek(t *testing.T) {
	var (
		actual   time.Time
		expected time.Time
	)

	// Year starts on a Monday
	expected = time.Date(2019, time.March, 18, 0, 0, 0, 0, time.UTC)
	actual = utils.DateFromISOWeek(2019, 12)
	require.Equal(t, expected, actual)

	// Year doesn't start on a Monday
	expected = time.Date(2020, time.June, 15, 0, 0, 0, 0, time.UTC)
	actual = utils.DateFromISOWeek(2020, 25)
	require.Equal(t, expected, actual)

	// First week of the year
	expected = time.Date(2021, time.January, 4, 0, 0, 0, 0, time.UTC)
	actual = utils.DateFromISOWeek(2021, 1)
	require.Equal(t, expected, actual)

	// Last week of the year
	expected = time.Date(2021, time.December, 27, 0, 0, 0, 0, time.UTC)
	actual = utils.DateFromISOWeek(2021, 52)
	require.Equal(t, expected, actual)
}

func TestWeekSub(t *testing.T) {
	var (
		a *utils.Week
		b *utils.Week
	)

	// Subtracting identical Weeks
	a = utils.NewWeek(time.Date(2019, time.March, 18, 0, 0, 0, 0, time.UTC))
	require.Equal(t, 0, a.Sub(a))

	// Subtracting Weeks in the same year
	a = utils.NewWeek(time.Date(2021, time.January, 4, 0, 0, 0, 0, time.UTC))
	b = utils.NewWeek(time.Date(2021, time.January, 28, 0, 0, 0, 0, time.UTC))
	require.Equal(t, -3, a.Sub(b))
	require.Equal(t, 3, b.Sub(a))

	// Subtracting Weeks in different years
	a = utils.NewWeek(time.Date(2019, time.December, 11, 0, 0, 0, 0, time.UTC))
	b = utils.NewWeek(time.Date(2020, time.January, 10, 0, 0, 0, 0, time.UTC))
	require.Equal(t, -4, a.Sub(b))
	require.Equal(t, 4, b.Sub(a))

	// Subtracting Weeks in different years
	a = utils.NewWeek(time.Date(2019, time.January, 7, 0, 0, 0, 0, time.UTC))
	b = utils.NewWeek(time.Date(2020, time.January, 10, 0, 0, 0, 0, time.UTC))
	require.Equal(t, -52, a.Sub(b))
	require.Equal(t, 52, b.Sub(a))
}

func TestWeekIter(t *testing.T) {
	var (
		err   error
		ok    bool
		week  *utils.Week
		prev  *utils.Week
		iter  *utils.WeekIterator
		start time.Time
		end   time.Time
	)

	// Should error if start time is after end time
	start = time.Date(2021, time.January, 28, 0, 0, 0, 0, time.UTC)
	end = time.Date(2021, time.January, 4, 0, 0, 0, 0, time.UTC)
	iter, err = utils.GetWeekIterator(start, end)
	require.Error(t, err)
	require.Nil(t, iter)

	// Time range of one week
	start = time.Date(2021, time.January, 4, 0, 0, 0, 0, time.UTC)
	end = time.Date(2021, time.January, 8, 0, 0, 0, 0, time.UTC)
	iter, err = utils.GetWeekIterator(start, end)
	require.NoError(t, err)
	week, ok = iter.Next()
	require.True(t, ok)
	require.Equal(t, utils.NewWeek(start), week)
	week, ok = iter.Next()
	require.False(t, ok)
	require.Nil(t, week)

	// Should be able to iterate over multiple weeks
	start = time.Date(2020, time.December, 25, 0, 0, 0, 0, time.UTC)
	end = time.Date(2021, time.January, 27, 0, 0, 0, 0, time.UTC)
	iter, err = utils.GetWeekIterator(start, end)
	require.NoError(t, err)
	require.NotNil(t, iter)
	for i := 0; ; i++ {
		prev = week
		if week, ok = iter.Next(); !ok {
			break
		}
		require.NotNil(t, week)
		require.LessOrEqual(t, i, 5)
	}
	require.Equal(t, utils.NewWeek(time.Date(2021, time.January, 25, 0, 0, 0, 0, time.UTC)), prev)
}
