package activity

import "errors"

var (
	ErrMissingTopic   = errors.New("missing activity topic")
	ErrInvalidWindow  = errors.New("aggregation window must be greater than 0")
	ErrUnknownNetwork = errors.New("unknown network, expected testnet, mainnet, or rvasp")
)
