package admin

import "net/http"

// ClientOption allows us to configure the APIv2 client when it is created.
type ClientOption func(c *APIv2) error

func WithClient(client *http.Client) ClientOption {
	return func(c *APIv2) error {
		c.client = client
		return nil
	}
}
