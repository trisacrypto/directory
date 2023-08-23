package activity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/ensign"
	"github.com/vmihailenco/msgpack/v5"
)

func TestPublisher(t *testing.T) {
	// Test publisher with bad configuration does not start
	conf := activity.Config{
		Enabled:           true,
		AggregationWindow: time.Minute * 5,
	}
	require.ErrorIs(t, activity.ErrMissingTopic, activity.Start(conf), "expected missing topic error")

	// Test publisher in disabled mode starts
	activity.Reset()
	conf.Enabled = false
	require.NoError(t, activity.Start(conf), "expected no error for disabled publisher")

	// Activity methods are no ops but should not panic
	activity.Lookup().Add()
	activity.Search().VASP(uuid.New()).Add()
	activity.Stop()

	// Test publisher starts with valid configuration
	activity.Reset()
	conf.Enabled = true
	conf.Network = "testnet"
	conf.Topic = "network-activity"
	conf.AggregationWindow = time.Millisecond
	conf.Ensign = ensign.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Endpoint:     "mock ensign endpoint",
		AuthURL:      "mock ensign auth url",
	}
	conf.Testing = true
	require.NoError(t, activity.Start(conf), "expected no error for valid configuration")

	// Configure the Ensign mock to assert that events are being published
	vaspID := uuid.New()
	emock := activity.GetEnsignMock()
	published := make(chan struct{})
	emock.OnPublish = func(stream api.Ensign_PublishServer) (err error) {
		// Receive the initial request from the client
		_, err = stream.Recv()
		assert.NoError(t, err, "expected no error receiving initial request from mock server")

		// Send back the topic map to the client
		err = stream.Send(&api.PublisherReply{
			Embed: &api.PublisherReply_Ready{
				Ready: &api.StreamReady{
					ClientId: "client-id",
					ServerId: "server-id",
					Topics: map[string][]byte{
						"network-activity": ulid.Make().Bytes(),
					},
				},
			},
		})
		assert.NoError(t, err, "expected no error sending ready message from mock server")

		// Receive events from the client. Publish happens on intervals, so wait until
		// both activity entries are published. If not received, the test should timeout.
		var lookup, search bool
		for lookup == false && search == false {
			msg, err := stream.Recv()
			assert.NoError(t, err, "expected no error receiving event from activity publisher")

			// Should be able to unmarshal and parse the event
			event, err := msg.GetEvent().Unwrap()
			assert.NoError(t, err, "expected no error unwrapping event")
			assert.Equal(t, activity.NetworkActivityMimeType, event.Mimetype, "expected network activity mime type for event")
			assert.True(t, event.Type.Equals(&activity.NetworkActivityEventType), "expected network activity event type")
			network, ok := event.Metadata["network"]
			assert.True(t, ok, "expected network metadata to be present")
			assert.Equal(t, "testnet", network, "expected network in metadata to be testnet")
			acv := &activity.NetworkActivity{}
			assert.NoError(t, msgpack.Unmarshal(event.Data, acv), "expected no error unmarshaling event")

			assert.Equal(t, activity.TestNet, acv.Network, "expected network to be testnet")
			assert.Equal(t, time.Millisecond, acv.Window, "expected window to be 1ms")

			if _, ok := acv.Activity[activity.LookupActivity]; ok {
				_, ok = event.Metadata["has_activity"]
				assert.True(t, ok, "expected has_activity metadata to be present")
				lookup = true
			}

			if _, ok := acv.Activity[activity.SearchActivity]; ok {
				_, ok = event.Metadata["has_activity"]
				assert.True(t, ok, "expected has_activity metadata to be present")
				assert.Equal(t, uint64(1), acv.VASPActivity[vaspID][activity.SearchActivity], "expected search activity to be 1 for VASP")
				search = true
			}

			published <- struct{}{}
			return nil
		}

		return nil
	}

	// Should be able to add activity entries
	activity.Lookup().Add()
	activity.Search().VASP(vaspID).Add()

	// Wait for the publisher to publish the event
	<-published
	activity.Stop()
}
