// Package migo provides a Misskey API client.
//
// The root package defines the Client interface, transport implementation,
// and API error handling. Domain types live in the sub-package migo/types.
//
// Basic usage:
//
//	client := migo.NewClient(
//	    migo.WithBaseURL("https://your-instance.net/api"),
//	    migo.WithToken("YOUR_TOKEN"),
//	)
//	ns := notes.NewService(client)
//	note, err := ns.Create(ctx, &notes.CreateRequest{
//	    Text: proto.String("Hello, Misskey!"),
//	})
package migo

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Client is the core interface for making Misskey API requests.
// Implementations must handle authentication, serialization, and
// error response parsing.
//
// To mock in tests, implement this interface.
type Client interface {
	// Do sends a POST request to the given API endpoint path.
	//   - path: e.g. "/notes/create" (the /api prefix is added automatically)
	//   - req:  request body, serialized to JSON (may be nil)
	//   - resp: response body destination, JSON-deserialized (may be nil)
	// Returns an *APIError for non-2xx responses.
	Do(ctx context.Context, path string, req, resp any) error
}

// requestDoer is an internal alias used for documentation.
type requestDoer = Client

// APIError represents a Misskey API error response.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	ID      string `json:"id"`

	// HTTP status code from the response.
	StatusCode int `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("misskey API error [%d] %s: %s (id=%s)", e.StatusCode, e.Code, e.Message, e.ID)
}

// ---------------------------------------------------------------------------
// Client options & constructor
// ---------------------------------------------------------------------------

// ClientOption configures a Client.
type ClientOption func(*clientOptions)

type clientOptions struct {
	baseURL    string
	token      string
	httpClient *http.Client
	userAgent  string
}

// WithBaseURL sets the Misskey instance API base URL.
// Default: "https://misskey.io/api"
func WithBaseURL(url string) ClientOption {
	return func(o *clientOptions) {
		o.baseURL = url
	}
}

// WithToken sets the Bearer token for authenticated requests.
func WithToken(token string) ClientOption {
	return func(o *clientOptions) {
		o.token = token
	}
}

// WithHTTPClient sets a custom HTTP client (useful for mocking or proxies).
func WithHTTPClient(c *http.Client) ClientOption {
	return func(o *clientOptions) {
		o.httpClient = c
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(ua string) ClientOption {
	return func(o *clientOptions) {
		o.userAgent = ua
	}
}

// NewClient creates a new Misskey API client.
func NewClient(opts ...ClientOption) Client {
	o := &clientOptions{
		baseURL: "https://misskey.io/api",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "migo (Go Misskey client)",
	}
	for _, opt := range opts {
		opt(o)
	}

	return &client{
		baseURL:    o.baseURL,
		token:      o.token,
		httpClient: o.httpClient,
		userAgent:  o.userAgent,
	}
}
