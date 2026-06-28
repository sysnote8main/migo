// Package following provides the FollowingService for follow management.
package following

import (
	"context"

	"github.com/sysnote8main/migo/types"
)

// Service provides following-related API operations.
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new FollowingService.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// CreateRequest represents the request for following/create.
type CreateRequest struct {
	UserID      types.ID `json:"userId"`
	WithReplies *bool    `json:"withReplies,omitempty"`
}

// DeleteRequest represents the request for following/delete.
type DeleteRequest struct {
	UserID types.ID `json:"userId"`
}

// InvalidateRequest represents the request for following/invalidate.
type InvalidateRequest struct {
	UserID types.ID `json:"userId"`
}

// ListRequest represents the request for following/list.
type ListRequest struct {
	Pagination
	Notification *bool `json:"notify,omitempty"`
}

// RequestsListRequest represents the request for following/requests/list.
type RequestsListRequest struct {
	Pagination
}

// RequestsActionRequest represents the request for accept/reject/cancel.
type RequestsActionRequest struct {
	UserID types.ID `json:"userId"`
}

// ---------------------------------------------------------------------------
// Pagination embedding helper
// ---------------------------------------------------------------------------

// Pagination is embedded into list requests.
type Pagination struct {
	SinceID   *types.ID `json:"sinceId,omitempty"`
	UntilID   *types.ID `json:"untilId,omitempty"`
	SinceDate *int64    `json:"sinceDate,omitempty"`
	UntilDate *int64    `json:"untilDate,omitempty"`
	Limit     *int      `json:"limit,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// Follower represents a follower/following entry.
type Follower struct {
	ID          types.ID        `json:"id"`
	CreatedAt   string          `json:"createdAt"`
	Followee    *types.UserLite `json:"followee,omitempty"`
	Follower    *types.UserLite `json:"follower,omitempty"`
	WithReplies bool            `json:"withReplies"`
}

// ---------------------------------------------------------------------------
// API methods
// ---------------------------------------------------------------------------

// Create follows a user.
func (s *Service) Create(ctx context.Context, req *CreateRequest) (*types.UserDetailed, error) {
	var user types.UserDetailed
	if err := s.cli.Do(ctx, "/following/create", req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete unfollows a user.
func (s *Service) Delete(ctx context.Context, userID types.ID) error {
	req := &DeleteRequest{UserID: userID}
	return s.cli.Do(ctx, "/following/delete", req, nil)
}

// Invalidate removes a follower (force-unfollow them).
func (s *Service) Invalidate(ctx context.Context, userID types.ID) error {
	req := &InvalidateRequest{UserID: userID}
	return s.cli.Do(ctx, "/following/invalidate", req, nil)
}

// List returns the list of users the authenticated user is following.
func (s *Service) List(ctx context.Context, req *ListRequest) ([]*Follower, error) {
	var followers []*Follower
	if err := s.cli.Do(ctx, "/following/list", req, &followers); err != nil {
		return nil, err
	}
	return followers, nil
}

// RequestsList returns pending follow requests.
func (s *Service) RequestsList(ctx context.Context, req *RequestsListRequest) ([]*Follower, error) {
	var followers []*Follower
	if err := s.cli.Do(ctx, "/following/requests/list", req, &followers); err != nil {
		return nil, err
	}
	return followers, nil
}

// RequestsAccept accepts a follow request.
func (s *Service) RequestsAccept(ctx context.Context, userID types.ID) error {
	req := &RequestsActionRequest{UserID: userID}
	return s.cli.Do(ctx, "/following/requests/accept", req, nil)
}

// RequestsReject rejects a follow request.
func (s *Service) RequestsReject(ctx context.Context, userID types.ID) error {
	req := &RequestsActionRequest{UserID: userID}
	return s.cli.Do(ctx, "/following/requests/reject", req, nil)
}

// RequestsCancel cancels a pending follow request.
func (s *Service) RequestsCancel(ctx context.Context, userID types.ID) error {
	req := &RequestsActionRequest{UserID: userID}
	return s.cli.Do(ctx, "/following/requests/cancel", req, nil)
}
