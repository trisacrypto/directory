package models_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/activity"
)

func TestActivityMonth(t *testing.T) {
	// Add activity to a month
	month := &models.ActivityMonth{
		Date: "2023-08",
	}
	aliceID := uuid.New()
	acv := &activity.NetworkActivity{
		Network: activity.MainNet,
		Activity: activity.ActivityCount{
			activity.LookupActivity: 1,
		},
		VASPActivity: map[uuid.UUID]activity.ActivityCount{
			aliceID: {
				activity.LookupActivity: 1,
			},
		},
		Window: time.Minute * 5,
	}
	var err error
	acv.Timestamp, err = time.Parse(time.RFC3339, "2023-08-24T00:00:00Z")
	require.NoError(t, err, "could not create activity timestamp")
	month.Add(acv)
	require.Len(t, month.Days, 1, "should have one day in the month")
	require.Equal(t, "2023-08-24", month.Days[0].Date, "day has the wrong date")

	// Add an activity to the same day
	acv.Timestamp, err = time.Parse(time.RFC3339, "2023-08-24T00:01:00Z")
	require.NoError(t, err, "could not create activity timestamp")
	month.Add(acv)
	require.Len(t, month.Days, 1, "should have one day in the month")

	// Add an activity to a different day
	acv.Timestamp, err = time.Parse(time.RFC3339, "2023-08-26T00:01:00Z")
	require.NoError(t, err, "could not create activity timestamp")
	month.Add(acv)
	require.Len(t, month.Days, 2, "should have two days in the month")
	require.Equal(t, "2023-08-26", month.Days[1].Date, "day has the wrong date")

	// Add an activity which is between the two days
	acv.Timestamp, err = time.Parse(time.RFC3339, "2023-08-25T00:01:00Z")
	require.NoError(t, err, "could not create activity timestamp")
	month.Add(acv)
	require.Len(t, month.Days, 3, "should have three days in the month")
	require.Equal(t, "2023-08-24", month.Days[0].Date, "day has the wrong date")
	require.Equal(t, "2023-08-25", month.Days[1].Date, "day has the wrong date")
	require.Equal(t, "2023-08-26", month.Days[2].Date, "day has the wrong date")
}
