package migo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// client is the default Client implementation.
type client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	userAgent  string
}

// Do sends an authenticated JSON POST request to the Misskey API.
func (c *client) Do(ctx context.Context, path string, req, resp any) error {
	// Build URL.
	base := strings.TrimRight(c.baseURL, "/")
	endpoint := base + "/" + strings.TrimLeft(path, "/")

	// Serialize request body.
	var bodyReader io.Reader
	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("migo: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bodyReader)
	if err != nil {
		return fmt.Errorf("migo: create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("User-Agent", c.userAgent)
	if c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("migo: execute request: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("migo: read response body: %w", err)
	}

	if httpResp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: httpResp.StatusCode}
		// Try to parse the error body (best-effort).
		if err := json.Unmarshal(body, apiErr); err != nil {
			apiErr.Message = strings.TrimSpace(string(body))
		}
		return apiErr
	}

	if resp != nil {
		if err := json.Unmarshal(body, resp); err != nil {
			return fmt.Errorf("migo: unmarshal response: %w", err)
		}
	}

	return nil
}

// Ensure client implements Client at compile time.
var _ Client = (*client)(nil)
