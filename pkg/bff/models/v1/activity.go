package models

import (
	"sort"
	"time"

	"github.com/trisacrypto/directory/pkg/utils/activity"
)

const (
	DateLayout = "2006-01-02"
)

func (a *ActivityCount) Add(network activity.Network, acv activity.Activity, count uint64) {
	switch network {
	case activity.MainNet:
		if a.Mainnet == nil {
			a.Mainnet = make(map[string]uint64)
		}
		a.Mainnet[acv.String()] += count
	case activity.TestNet:
		if a.Testnet == nil {
			a.Testnet = make(map[string]uint64)
		}
		a.Testnet[acv.String()] += count
	case activity.RVASP:
		if a.RVASP == nil {
			a.RVASP = make(map[string]uint64)
		}
		a.RVASP[acv.String()] += count
	}
}

// Create a new activity day from the date
func NewActivityDay(date string) *ActivityDay {
	return &ActivityDay{
		Date:         date,
		Activity:     &ActivityCount{},
		VaspActivity: make(map[string]*ActivityCount),
	}
}

// Add activity to the day.
func (d *ActivityDay) Add(a *activity.NetworkActivity) {
	for acv, count := range a.Activity {
		d.Activity.Add(a.Network, acv, count)
	}

	for id, counts := range a.VASPActivity {
		vasp := id.String()
		if _, ok := d.VaspActivity[vasp]; !ok {
			d.VaspActivity[vasp] = &ActivityCount{}
		}

		for acv, count := range counts {
			d.VaspActivity[vasp].Add(a.Network, acv, count)
		}
	}
}

// Add the activity to the month, this will ensure that the activity is added to
// the correct day and that the day is created if it does not exist.
// Note: This assumes that the activity window is less than 24 hours.
func (m *ActivityMonth) Add(a *activity.NetworkActivity) {
	// Get the date from the activity window
	date := a.WindowEnd().Format(DateLayout)

	// Find the day in the month
	i := sort.Search(len(m.Days), func(i int) bool {
		return m.Days[i].Date >= date
	})

	if i == len(m.Days) {
		// If there is no day yet, create it
		day := NewActivityDay(date)
		day.Add(a)
		m.Days = append(m.Days, day)
	} else if m.Days[i].Date == date {
		// If the day already exists, add the activity to it
		m.Days[i].Add(a)
	} else {
		// If the day is between two days, create the day and insert it between the two
		// days
		day := NewActivityDay(date)
		day.Add(a)
		m.Days = append(m.Days, nil)
		copy(m.Days[i+1:], m.Days[i:])
		m.Days[i] = day
	}
}

// Return the key used to lookup an ActivityMonth in the ActivityStore, which is the
// byte slice representation of the date in the form YYYY-MM. This method also
// validates that the date is correct.
func (m *ActivityMonth) Key() (_ []byte, err error) {
	if _, err = time.Parse(MonthLayout, m.Date); err != nil {
		return nil, err
	}
	return []byte(m.Date), nil
}
