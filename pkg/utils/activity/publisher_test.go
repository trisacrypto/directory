package activity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	sdk "github.com/rotationalio/go-ensign"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	"github.com/rotationalio/go-ensign/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/ensign"
	"github.com/vmihailenco/msgpack/v5"
)

func TestPublisher(t *testing.T) {
	// Setup the Ensign mock and client
	emock := mock.New(nil)
	client, err := sdk.New(sdk.WithMock(emock))
	require.NoError(t, err, "expected no error creating mock ensign client")
	activity.SetClient(client)

	// Test publisher with bad configuration does not start
	conf := activity.Config{
		Enabled:           true,
		AggregationWindow: time.Minute * 5,
	}
	require.Error(t, activity.Start(conf), "expected error with bad configuration")

	// Test publisher in disabled mode starts
	activity.Reset()
	conf.Enabled = false
	require.NoError(t, activity.Start(conf), "expected no error for disabled publisher")

	// Activity methods are no ops but should not panic
	activity.Lookup().Add()
	activity.Search().VASP(uuid.New().String()).Add()
	activity.Register().Add()
	activity.KeyExchange().Add()
	activity.Stop()

	// Test publisher starts with valid configuration
	activity.Reset()
	conf.Enabled = true
	conf.Network = activity.TestNet
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
	vaspID := uuid.New().String()
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
		var lookup, search, register, keyExchange bool
		for lookup == false && search == false && register == false && keyExchange == false {
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

			if _, ok := acv.Activity[activity.RegisterActivity]; ok {
				_, ok = event.Metadata["has_activity"]
				assert.True(t, ok, "expected has_activity metadata to be present")
				register = true
			}

			if _, ok := acv.Activity[activity.KeyExchangeActivity]; ok {
				_, ok = event.Metadata["has_activity"]
				assert.True(t, ok, "expected has_activity metadata to be present")
				keyExchange = true
			}
		}
		published <- struct{}{}

		return nil
	}

	// Should be able to add activity entries
	activity.Lookup().Add()
	activity.Search().VASP(vaspID).Add()
	activity.Register().Add()
	activity.KeyExchange().Add()

	// Wait for the publisher to publish the event
	<-published
	activity.Stop()
}
