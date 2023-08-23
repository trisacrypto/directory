package activity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/ensign"
)

func TestPublisher(t *testing.T) {
	// Test publisher with bad configuration does not start
	conf := activity.Config{
		Enabled:           true,
		AggregationWindow: time.Minute * 5,
	}
	require.ErrorIs(t, activity.ErrMissingTopic, activity.Start(conf))

	// Test publisher in disabled mode starts
	activity.Reset()
	conf.Enabled = false
	require.NoError(t, activity.Start(conf))

	// Activity methods are no ops but should not panic
	activity.Lookup().Add()
	activity.Search().VASP(uuid.New()).Add()
	activity.Stop()

	// Test publisher with valid configuration starts
	activity.Reset()
	conf.Enabled = true
	conf.Topic = "network-activity"
	conf.Ensign = ensign.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Endpoint:     "mock ensign endpoint",
		AuthURL:      "mock ensign auth url",
	}
	conf.Testing = true
	require.NoError(t, activity.Start(conf))

	// Should be able to add activity entries
	activity.Lookup().Add()
	activity.Search().VASP(uuid.New()).Add()
	activity.Stop()
}
