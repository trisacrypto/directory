package utils

import (
	"context"
	"time"
)

// TODO: Timeout should be configurable.
const deadline time.Duration = time.Second * 30

func WithContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), deadline)
}
