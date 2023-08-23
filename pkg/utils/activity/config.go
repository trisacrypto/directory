package activity

import "time"

type Config struct {
	Enabled           bool
	Topic             string
	AggregationWindow time.Duration
}
