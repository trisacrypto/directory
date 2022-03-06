package api

import (
	"context"
	"time"
)

//===========================================================================
// Service Interface
//===========================================================================

type BFFClient interface {
	Status(ctx context.Context) (out *StatusReply, err error)
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are used for generic API responses and errors
type Reply struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
}

// StatusReply is returned on status requests. Note that no request is needed.
type StatusReply struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Version   string    `json:"version,omitempty"`
}
