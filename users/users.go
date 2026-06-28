// Package users provides the UserService for querying Misskey users.
package users

import (
	"context"

	"github.com/sysnote8main/migo/types"
)

// Service provides user-related API operations.
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new UserService.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// ShowRequest represents the request body for users/show.
// Either UserID, UserIDs, or Username must be set.
type ShowRequest struct {
	UserID   *types.ID  `json:"userId,omitempty"`
	UserIDs  []types.ID `json:"userIds,omitempty"`
	Username *string    `json:"username,omitempty"`
	Host     *string    `json:"host,omitempty"`
}

// SearchRequest represents the request body for users/search.
type SearchRequest struct {
	Query  string  `json:"query"`
	Offset *int    `json:"offset,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	Origin *string `json:"origin,omitempty"`
	Detail *bool   `json:"detail,omitempty"`
}

// ---------------------------------------------------------------------------
// API methods
// ---------------------------------------------------------------------------

// Show fetches user details. The response is either a single UserDetailed
// (when UserID or Username is specified) or []UserDetailed (when UserIDs is set).
func (s *Service) Show(ctx context.Context, req *ShowRequest) (*types.UserDetailed, error) {
	var user types.UserDetailed
	if err := s.cli.Do(ctx, "/users/show", req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// ShowBulk fetches multiple users by their IDs.
func (s *Service) ShowBulk(ctx context.Context, ids []types.ID) ([]*types.UserDetailed, error) {
	req := &ShowRequest{UserIDs: ids}
	var users []*types.UserDetailed
	if err := s.cli.Do(ctx, "/users/show", req, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// Search searches for users by query.
func (s *Service) Search(ctx context.Context, req *SearchRequest) ([]*types.UserDetailed, error) {
	var users []*types.UserDetailed
	if err := s.cli.Do(ctx, "/users/search", req, &users); err != nil {
		return nil, err
	}
	return users, nil
}
