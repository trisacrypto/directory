package bff

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/go-ensign"
	pb "github.com/rotationalio/go-ensign/api/v1beta1"
	"github.com/rotationalio/go-ensign/mock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/models/v1"
	"github.com/trisacrypto/directory/pkg/store"
	storeerrors "github.com/trisacrypto/directory/pkg/store/errors"
	"github.com/trisacrypto/directory/pkg/utils/activity"
	"github.com/trisacrypto/directory/pkg/utils/sentry"
)

// NetworkActivity returns a time series representation of activity broken down by
// network. This endpoint is publicly accessible and requires no authentication.
//
// @Summary Get network activity for the last 30 days.
// @Description Returns a time series of activity in testnet and mainnet.
// @Produce json
// @Success 200 {object} api.NetworkActivityReply
// @Failure 500 {object} api.Reply
func (s *Server) NetworkActivity(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Get the month(s) for the last 30 days
	now := currentTime()
	currentMonth := now.Format(models.MonthLayout)
	startingMonth := now.AddDate(0, 0, -30).Format(models.MonthLayout)

	// Retrieve the months from the database
	months := make([]string, 0, 2)
	months = append(months, currentMonth)
	if startingMonth != currentMonth {
		months = append(months, startingMonth)
	}

	// Build the time series from months in the database
	out := &api.NetworkActivityReply{
		TestNet: make([]api.Activity, 0, 30),
		MainNet: make([]api.Activity, 0, 30),
	}

	// Walk backwards through the days and build the time series
	day := now
	remainingDays := 30
	for _, date := range months {
		var (
			month *models.ActivityMonth
			err   error
		)

		// Retrieve the activity month from the database
		if month, err = s.db.RetrieveActivityMonth(ctx, date); err != nil {
			switch {
			case errors.Is(err, storeerrors.ErrEntityNotFound):
				// If month is not found, assume no activity
				month = &models.ActivityMonth{
					Days: make([]*models.ActivityDay, 0),
				}
			default:
				sentry.Error(c).Err(err).Msg("failed to retrieve activity month")
				c.JSON(http.StatusInternalServerError, "could not retrieve network activity")
				return
			}
		}

		monthIndex := len(month.Days) - 1
	monthLoop:
		for remainingDays > 0 {
			date := day.Format(models.DateLayout)

			var activityDay *models.ActivityDay
			if monthIndex >= 0 && monthIndex < len(month.Days) {
				activityDay = month.Days[monthIndex]
			}
			testnet := api.Activity{
				Date: date,
			}
			mainnet := api.Activity{
				Date: date,
			}
			if activityDay != nil && activityDay.Date == date {
				// Aggregate the counts for this day
				for _, count := range activityDay.Activity.Testnet {
					testnet.Events += count
				}

				for _, count := range activityDay.Activity.Mainnet {
					mainnet.Events += count
				}

				// Move the month index back
				monthIndex--
			}

			// We are working backwards so prepend the activity
			out.TestNet = append([]api.Activity{testnet}, out.TestNet...)
			out.MainNet = append([]api.Activity{mainnet}, out.MainNet...)

			// Adjust the day and remaining days
			remainingDays--
			prevDay := day
			day = day.AddDate(0, 0, -1)
			if prevDay.Day() == 1 {
				// Reached the beginning of this month, go to the previous month
				break monthLoop
			}
		}
	}

	c.JSON(http.StatusOK, out)
}

// ActivitySubscriber is a struct with a go routine that subscribes to a network
// activity topic in Ensign and applies asynchronous updates to trtl.
type ActivitySubscriber struct {
	client *ensign.Client
	db     store.Store
	emock  *mock.Ensign
	topic  string
	stop   chan struct{}
}

func NewActivitySubscriber(conf activity.Config, db store.Store) (sub *ActivitySubscriber, err error) {
	if !conf.Enabled {
		return nil, errors.New("activity subscriber is disabled")
	}

	if err = conf.Validate(); err != nil {
		return nil, err
	}

	sub = &ActivitySubscriber{
		topic: conf.Topic,
		db:    db,
	}

	if conf.Testing {
		sub.emock = mock.New(nil)
		if sub.client, err = ensign.New(ensign.WithMock(sub.emock)); err != nil {
			return nil, err
		}
	} else {
		if sub.client, err = conf.Ensign.Client(); err != nil {
			return nil, err
		}
	}

	return sub, nil
}

// Run the subscriber under the waitgroup to allow the caller to wait for the
// subscriber to exit after calling Stop().
func (s *ActivitySubscriber) Run(wg *sync.WaitGroup) error {
	if s.stop != nil {
		return errors.New("activity subscriber is already running")
	}

	if wg == nil {
		return errors.New("waitgroup must be provided to run the activity subscriber")
	}

	s.stop = make(chan struct{})
	wg.Add(1)
	go func() {
		s.Subscribe()
		s.stop = nil
		wg.Done()
	}()
	return nil
}

// Stop the activity subscriber.
func (s *ActivitySubscriber) Stop() {
	if s.stop != nil {
		close(s.stop)
	}
}

func (s *ActivitySubscriber) Subscribe() {
	var (
		err error
		sub *ensign.Subscription
	)

	// Subscribe to the topic
	if sub, err = s.client.Subscribe(s.topic); err != nil {
		// Note: Using WithLevel with FatalLevel does not exit the program but we want
		// to know what the error was.
		log.WithLevel(zerolog.FatalLevel).Err(err).Msg("failed to subscribe to network activity topic")
		return
	}
	defer sub.Close()

	// Parse events and make the updates in trtl
	for {
		select {
		case <-s.stop:
			return
		case event, ok := <-sub.C:
			if !ok {
				return
			}

			// Parse the event into a network activity update
			var update *activity.NetworkActivity
			if update, err = activity.Parse(event); err != nil {
				log.Error().Err(err).Msg("failed to parse network activity event")
				event.Nack(pb.Nack_UNKNOWN_TYPE)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Fetch the network activity month from the database to update it.
			// Note: There is a potential race condition with this pattern if two
			// routines are updating the month concurrently, however currently this is
			// the only go routine that is writing to the activity store.
			date := update.WindowEnd().Format(models.MonthLayout)
			var month *models.ActivityMonth
			if month, err = s.db.RetrieveActivityMonth(ctx, date); err != nil {
				switch {
				case errors.Is(err, storeerrors.ErrEntityNotFound):
					// Create the activity month if it does not exist
					month = &models.ActivityMonth{
						Date: date,
					}
					if err = s.db.UpdateActivityMonth(ctx, month); err != nil {
						log.Error().Err(err).Str("month_date", date).Msg("failed to create activity month")
						event.Nack(pb.Nack_UNPROCESSED)
						continue
					}
				default:
					log.Error().Err(err).Str("month_date", date).Msg("failed to retrieve activity month")
					event.Nack(pb.Nack_UNPROCESSED)
					continue
				}
			}

			// Update the activity month
			month.Add(update)
			if err = s.db.UpdateActivityMonth(ctx, month); err != nil {
				log.Error().Err(err).Str("month_date", date).Msg("failed to update activity month")
				event.Nack(pb.Nack_UNPROCESSED)
				continue
			}

			// Acknowledge the event
			event.Ack()
		}
	}
}

// Expose the Ensign server mock to the tests
func (s *ActivitySubscriber) GetEnsignMock() *mock.Ensign {
	return s.emock
}
