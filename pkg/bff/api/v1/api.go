package api

import (
	"context"
)

//===========================================================================
// Service Interface
//===========================================================================

type BFFClient interface {
	Status(ctx context.Context, in *StatusParams) (out *StatusReply, err error)
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are used for generic API responses and errors
type Reply struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
}

// StatusParams is parsed from the query parameters of the GET request
type StatusParams struct {
	NoGDS bool `url:"nogds,omitempty" form:"nogds" default:"false"`
}

// StatusReply is returned on status requests. Note that no request is needed.
type StatusReply struct {
	Status  string `json:"status"`
	Uptime  string `json:"uptime,omitempty"`
	Version string `json:"version,omitempty"`
	TestNet string `json:"testnet,omitempty"`
	MainNet string `json:"mainnet,omitempty"`
}
