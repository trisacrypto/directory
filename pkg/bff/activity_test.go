package bff_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	pb "github.com/rotationalio/go-ensign/api/v1beta1"
	mimetype "github.com/rotationalio/go-ensign/mimetype/v1beta1"
	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/ensign"
	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *bffTestSuite) TestNetworkActivity() {
	require := s.Require()
	defer bff.ResetTime()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// If there is no activity in the database, should return zeros
	rep, err := s.client.NetworkActivity(ctx)
	require.NoError(err, "could not get network activity")
	require.Len(rep.MainNet, 30, "should have 30 days of activity for mainnet")
	require.Len(rep.TestNet, 30, "should have 30 days of activity for testnet")
	assertNoCounts(require, rep.MainNet)
	assertNoCounts(require, rep.TestNet)

	// Mock the current time for consistent testing
	now, err := time.Parse(time.RFC3339, "2023-08-24T12:00:00Z")
	require.NoError(err, "could not create activity time")
	bff.MockTime(now)

	// Add an empty month to the database but with no data
	month := &models.ActivityMonth{
		Date: now.Format(models.MonthLayout),
	}
	require.NoError(s.DB().UpdateActivityMonth(ctx, month), "could not create activity month")
	rep, err = s.client.NetworkActivity(ctx)
	require.NoError(err, "could not get network activity")
	assertNoCounts(require, rep.MainNet)
	assertNoCounts(require, rep.TestNet)

	// Add some activity to the database
	month.Add(&activity.NetworkActivity{
		Network: activity.MainNet,
		Activity: activity.ActivityCount{
			activity.LookupActivity: 1,
		},
		Timestamp: now,
		Window:    time.Minute * 5,
	})
	require.NoError(s.DB().UpdateActivityMonth(ctx, month), "could not update activity month")
	rep, err = s.client.NetworkActivity(ctx)
	require.NoError(err, "could not get network activity")
	require.Len(rep.MainNet, 30, "should have 30 days of activity for mainnet")
	require.Len(rep.TestNet, 30, "should have 30 days of activity for testnet")
	require.Equal(uint64(1), rep.MainNet[29].Events, "should have one mainnet event")
	require.Zero(rep.TestNet[28].Events, "should have zero testnet events")

	// Create some more activity
	acv := &activity.NetworkActivity{
		Network: activity.TestNet,
		Activity: activity.ActivityCount{
			activity.LookupActivity: 3,
		},
		Timestamp: now.AddDate(0, 0, -29),
		Window:    time.Minute * 5,
	}

	// Test the case of being less than 30 days into the month
	prevMonth := &models.ActivityMonth{
		Date: acv.Timestamp.Format(models.MonthLayout),
	}
	prevMonth.Add(acv)
	require.NoError(s.DB().UpdateActivityMonth(ctx, prevMonth), "could not update activity month")
	rep, err = s.client.NetworkActivity(ctx)
	require.NoError(err, "could not get network activity")
	require.Len(rep.MainNet, 30, "should have 30 days of activity for mainnet")
	require.Len(rep.TestNet, 30, "should have 30 days of activity for testnet")
	require.Equal(uint64(1), rep.MainNet[29].Events, "should have one mainnet event on the last day")
	require.Equal(uint64(3), rep.TestNet[0].Events, "should have three testnet events on the first day")

	// Set the time to the end of the month to test the case of being more than 30 days into the month
	now, err = time.Parse(time.RFC3339, "2023-08-31T12:00:00Z")
	require.NoError(err, "could not create activity time")
	bff.MockTime(now)
	rep, err = s.client.NetworkActivity(ctx)
	require.NoError(err, "could not get network activity")
	require.Len(rep.MainNet, 30, "should have 30 days of activity for mainnet")
	require.Len(rep.TestNet, 30, "should have 30 days of activity for testnet")
	require.Equal(uint64(1), rep.MainNet[22].Events, "should have one mainnet event on August 24th")
	assertNoCounts(require, rep.TestNet)
}

// Convenience method that asserts an activity reply has no counts
func assertNoCounts(require *require.Assertions, counts []api.Activity) {
	dates := map[string]struct{}{}
	var prevDate string
	for _, acv := range counts {
		require.Zero(acv.Events, "expected zero events")
		_, ok := dates[acv.Date]
		require.False(ok, "should not have duplicate dates")
		dates[acv.Date] = struct{}{}
		if prevDate != "" {
			require.True(acv.Date > prevDate, "dates should be in chronological order")
		}
		prevDate = acv.Date
	}
}

func (s *bffTestSuite) TestActivitySubscriber() {
	require := s.Require()

	// If not enabled, should return an error
	conf := activity.Config{}
	_, err := bff.NewActivitySubscriber(conf, s.DB())
	require.Error(err, "should return an error if not enabled")

	// If config is not valid, should return an error
	conf.Enabled = true
	_, err = bff.NewActivitySubscriber(conf, s.DB())
	require.Error(err, "should return an error if config is invalid")

	// Test running the subscriber under a waitgroup
	conf.Topic = "network-activity"
	conf.Testing = true
	conf.Network = activity.TestNet
	conf.Ensign = ensign.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Endpoint:     "mock ensign endpoint",
		AuthURL:      "mock ensign auth url",
	}
	sub, err := bff.NewActivitySubscriber(conf, s.DB())
	require.NoError(err, "could not create subscriber")

	// Setup the network activity fixtures
	events := make([]*pb.EventWrapper, 0)
	aliceVASP := uuid.New().String()
	bobVASP := uuid.New().String()

	// First activity
	acv := &activity.NetworkActivity{
		Network: activity.MainNet,
		Activity: map[activity.Activity]uint64{
			activity.LookupActivity: 1,
		},
		VASPActivity: map[string]activity.ActivityCount{
			aliceVASP: {
				activity.LookupActivity: 1,
			},
		},
		Window: time.Minute * 5,
	}
	acv.Timestamp, err = time.Parse(time.RFC3339, "2023-08-24T12:00:00Z")
	require.NoError(err, "could not create first activity time")
	data, err := msgpack.Marshal(acv)
	require.NoError(err, "could not marshal first activity")
	event := &pb.Event{
		Data:     data,
		Mimetype: mimetype.MIME_APPLICATION_MSGPACK,
		Type:     &activity.NetworkActivityEventType,
		Created:  timestamppb.Now(),
	}
	wrapper := &pb.EventWrapper{
		Id:        []byte("eventID"),
		TopicId:   []byte("topicID"),
		Committed: timestamppb.Now(),
	}
	require.NoError(wrapper.Wrap(event), "could not wrap first event")
	events = append(events, wrapper)

	// Second activity - on the same day
	acv = &activity.NetworkActivity{
		Network: activity.MainNet,
		Activity: map[activity.Activity]uint64{
			activity.LookupActivity: 1,
			activity.SearchActivity: 3,
		},
		VASPActivity: map[string]activity.ActivityCount{
			aliceVASP: {
				activity.LookupActivity: 1,
			},
			bobVASP: {
				activity.SearchActivity: 3,
			},
		},
		Window: time.Minute * 5,
	}
	acv.Timestamp, err = time.Parse(time.RFC3339, "2023-08-24T12:05:00Z")
	require.NoError(err, "could not create second activity time")
	data, err = msgpack.Marshal(acv)
	require.NoError(err, "could not marshal second activity")
	event = &pb.Event{
		Data:     data,
		Mimetype: mimetype.MIME_APPLICATION_MSGPACK,
		Type:     &activity.NetworkActivityEventType,
		Created:  timestamppb.Now(),
	}
	wrapper = &pb.EventWrapper{
		Id:        []byte("eventID"),
		TopicId:   []byte("topicID"),
		Committed: timestamppb.Now(),
	}
	require.NoError(wrapper.Wrap(event), "could not wrap second event")
	events = append(events, wrapper)

	// Third activity - on the next day
	acv = &activity.NetworkActivity{
		Network: activity.MainNet,
		Activity: map[activity.Activity]uint64{
			activity.SearchActivity: 1,
		},
		VASPActivity: map[string]activity.ActivityCount{
			bobVASP: {
				activity.SearchActivity: 1,
			},
		},
		Window: time.Minute * 5,
	}
	acv.Timestamp, err = time.Parse(time.RFC3339, "2023-08-25T12:00:00Z")
	require.NoError(err, "could not create third activity time")
	data, err = msgpack.Marshal(acv)
	require.NoError(err, "could not marshal third activity")
	event = &pb.Event{
		Data:     data,
		Mimetype: mimetype.MIME_APPLICATION_MSGPACK,
		Type:     &activity.NetworkActivityEventType,
		Created:  timestamppb.Now(),
	}
	wrapper = &pb.EventWrapper{
		Id:        []byte("eventID"),
		TopicId:   []byte("topicID"),
		Committed: timestamppb.Now(),
	}
	require.NoError(wrapper.Wrap(event), "could not wrap third event")
	events = append(events, wrapper)

	// Fourth activity - in the next month
	acv = &activity.NetworkActivity{
		Network: activity.TestNet,
		Activity: map[activity.Activity]uint64{
			activity.LookupActivity: 2,
		},
		VASPActivity: map[string]activity.ActivityCount{},
		Window:       time.Minute * 5,
	}
	acv.Timestamp, err = time.Parse(time.RFC3339, "2023-09-01T12:00:00Z")
	require.NoError(err, "could not create fourth activity time")
	data, err = msgpack.Marshal(acv)
	require.NoError(err, "could not marshal fourth activity")
	event = &pb.Event{
		Data:     data,
		Mimetype: mimetype.MIME_APPLICATION_MSGPACK,
		Type:     &activity.NetworkActivityEventType,
		Created:  timestamppb.Now(),
	}
	wrapper = &pb.EventWrapper{
		Id:        []byte("eventID"),
		TopicId:   []byte("topicID"),
		Committed: timestamppb.Now(),
	}
	require.NoError(wrapper.Wrap(event), "could not wrap fourth event")
	events = append(events, wrapper)

	// Configure the Ensign mock to stream back the network activity
	emock := sub.GetEnsignMock()
	server := &sync.WaitGroup{}
	server.Add(1)
	emock.OnSubscribe = func(stream pb.Ensign_SubscribeServer) (err error) {
		defer server.Done()

		// Read the initial subscription request
		_, err = stream.Recv()
		if err != nil {
			return err
		}

		// Send the ready response back to the client
		if err = stream.Send(&pb.SubscribeReply{
			Embed: &pb.SubscribeReply_Ready{
				Ready: &pb.StreamReady{
					ClientId: "client-id",
					ServerId: "server-id",
					Topics: map[string][]byte{
						"network-activity": ulid.Make().Bytes(),
					},
				},
			},
		}); err != nil {
			return err
		}

		// Send the activity updates
		for _, event := range events {
			if err = stream.Send(&pb.SubscribeReply{
				Embed: &pb.SubscribeReply_Event{
					Event: event,
				},
			}); err != nil {
				return err
			}
		}

		// Should receive all the acks from the subscruber
		for i := 0; i < len(events); i++ {
			var req *pb.SubscribeRequest
			if req, err = stream.Recv(); err != nil {
				return err
			}

			if req.GetAck() == nil {
				return status.Errorf(codes.InvalidArgument, "expected ack")
			}
		}

		return nil
	}

	// Run the subscriber
	wg := &sync.WaitGroup{}
	require.NoError(sub.Run(wg), "could not run subscriber")

	// Wait for the subscriber to ack all the events
	server.Wait()
	sub.Stop()
	wg.Wait()

	// Should be two activity months in the store
	ctx := context.Background()
	expected := &models.ActivityMonth{
		Date: "2023-08",
		Days: []*models.ActivityDay{
			{
				Date: "2023-08-24",
				Activity: &models.ActivityCount{
					Mainnet: map[string]uint64{
						"lookup": 2,
						"search": 3,
					},
				},
				VaspActivity: map[string]*models.ActivityCount{
					aliceVASP: {
						Mainnet: map[string]uint64{
							"lookup": 2,
						},
					},
					bobVASP: {
						Mainnet: map[string]uint64{
							"search": 3,
						},
					},
				},
			},
			{
				Date: "2023-08-25",
				Activity: &models.ActivityCount{
					Mainnet: map[string]uint64{
						"search": 1,
					},
				},
				VaspActivity: map[string]*models.ActivityCount{
					bobVASP: {
						Mainnet: map[string]uint64{
							"search": 1,
						},
					},
				},
			},
		},
	}
	month, err := s.DB().RetrieveActivityMonth(ctx, "2023-08")
	require.NoError(err, "could not retrieve activity month")
	require.Equal(expected.Date, month.Date, "wrong activity month retrieved")
	require.Len(month.Days, len(expected.Days), "wrong number of activity days in month")
	for i, day := range month.Days {
		require.Equal(expected.Days[i].Date, day.Date, fmt.Sprintf("wrong date for day %d", i))
		require.Equal(expected.Days[i].Activity, day.Activity, fmt.Sprintf("wrong activity for day %d", i))
		require.Equal(expected.Days[i].VaspActivity, day.VaspActivity, fmt.Sprintf("wrong vasp activity for day %d", i))
	}

	expected = &models.ActivityMonth{
		Date: "2023-09",
		Days: []*models.ActivityDay{
			{
				Date: "2023-09-01",
				Activity: &models.ActivityCount{
					Testnet: map[string]uint64{
						"lookup": 2,
					},
				},
			},
		},
	}
	month, err = s.DB().RetrieveActivityMonth(ctx, "2023-09")
	require.NoError(err, "could not retrieve activity month")
	require.Equal("2023-09", month.Date, "wrong activity month retrieved")
	require.Len(month.Days, 1, "wrong number of activity days in month")
	require.Equal(expected.Days[0].Date, month.Days[0].Date, "wrong date for day")
	require.Equal(expected.Days[0].Activity, month.Days[0].Activity, "wrong activity for day")
	require.Nil(month.Days[0].VaspActivity, "expected no VASP activity for day")
}
