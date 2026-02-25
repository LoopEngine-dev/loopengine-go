package loopengine

import "net/http"

// Option configures a Client. Pass to New as optional arguments.
type Option func(*Client)

// WithHTTPClient sets the HTTP client. Default is http.DefaultClient.
// Use to set timeouts or a custom transport.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}
