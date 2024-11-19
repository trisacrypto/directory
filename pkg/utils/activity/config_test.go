package activity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/ensign"
)

func TestActivityValidation(t *testing.T) {
	conf := activity.Config{
		Enabled: true,
		Ensign: ensign.Config{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Endpoint:     "api.ensign.world:443",
			AuthURL:      "https://auth.ensign.world",
		},
	}

	// Test error is returned for missing topic
	require.ErrorIs(t, conf.Validate(), activity.ErrMissingTopic)

	// Test error is returned for invalid ensign configuration
	conf.Topic = "gds-activity"
	conf.AggregationWindow = time.Duration(5 * time.Minute)
	conf.Ensign.ClientID = ""
	require.ErrorIs(t, conf.Validate(), ensign.ErrMissingClientID)

	// Test disabled configuration is valid
	conf.Enabled = false
	require.NoError(t, conf.Validate())

	// Test valid enabled configuration
	conf.Enabled = true
	conf.Network = activity.TestNet
	conf.Ensign.ClientID = "client-id"
	require.NoError(t, conf.Validate())
}
