// Package auth provides the AuthService for MiAuth session-based authentication.
package auth

import (
	"context"
	"net/url"

	"github.com/sysnote8main/migo/types"
)

// Service provides MiAuth-related API operations.
//
// MiAuth flow:
//  1. SessionGenerate → obtain session token + MiAuth URL
//  2. User visits the MiAuth URL in a browser and approves
//  3. SessionUserKey → exchange session token for access token
//  4. Use the access token with migo.WithToken() for regular API calls
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new AuthService.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

// ---------------------------------------------------------------------------
// Request / Response types
// ---------------------------------------------------------------------------

// SessionGenerateRequest represents the request for auth/session/generate.
type SessionGenerateRequest struct {
	AppSecret string `json:"appSecret"`
}

// SessionGenerateResponse is the response from auth/session/generate.
type SessionGenerateResponse struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

// SessionUserKeyRequest represents the request for auth/session/userkey.
type SessionUserKeyRequest struct {
	AppSecret string `json:"appSecret"`
	Token     string `json:"token"`
}

// SessionUserKeyResponse is the response from auth/session/userkey.
type SessionUserKeyResponse struct {
	AccessToken string         `json:"accessToken"`
	User        types.UserLite `json:"user"`
}

// ---------------------------------------------------------------------------
// API methods
// ---------------------------------------------------------------------------

// SessionGenerate generates a new MiAuth session and returns the
// session token and the MiAuth URL for the user to approve.
func (s *Service) SessionGenerate(ctx context.Context, appSecret string) (*SessionGenerateResponse, error) {
	req := &SessionGenerateRequest{AppSecret: appSecret}
	var resp SessionGenerateResponse
	if err := s.cli.Do(ctx, "/auth/session/generate", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SessionUserKey exchanges a session token for an access token.
// The user must have approved the session before calling this.
func (s *Service) SessionUserKey(ctx context.Context, appSecret, token string) (*SessionUserKeyResponse, error) {
	req := &SessionUserKeyRequest{AppSecret: appSecret, Token: token}
	var resp SessionUserKeyResponse
	if err := s.cli.Do(ctx, "/auth/session/userkey", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MiAuthURL builds a MiAuth permission-request URL for the given session token
// and list of permissions. The user visits this URL to approve the app.
//
// Example permissions: "write:notes", "read:account", "read:drive"
func MiAuthURL(instanceURL, sessionToken string, permissions []string) string {
	u, _ := url.Parse(string(instanceURL))
	u = u.JoinPath("miauth", sessionToken)

	q := u.Query()
	for _, p := range permissions {
		q.Add("permission", p)
	}
	// The callback URL is optional; set it via WithName/WithCallback params.
	u.RawQuery = q.Encode()

	return u.String()
}

// MiAuthURLWithOptions builds a MiAuth URL with additional parameters.
func MiAuthURLWithOptions(instanceURL, sessionToken string, permissions []string, name, callback string) string {
	u, _ := url.Parse(instanceURL)
	u = u.JoinPath("miauth", sessionToken)

	q := u.Query()
	for _, p := range permissions {
		q.Add("permission", p)
	}
	if name != "" {
		q.Add("name", name)
	}
	if callback != "" {
		q.Add("callback", callback)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// Ping sends a ping to check if the API is responsive.
// Returns the server's ping response time as a number.
func (s *Service) Ping(ctx context.Context) (int64, error) {
	var resp struct {
		Pong int64 `json:"pong"`
	}
	if err := s.cli.Do(ctx, "/ping", nil, &resp); err != nil {
		return 0, err
	}
	return resp.Pong, nil
}

// I fetches the authenticated user's detailed profile.
func (s *Service) I(ctx context.Context) (*types.MeDetailed, error) {
	var user types.MeDetailed
	if err := s.cli.Do(ctx, "/i", nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
