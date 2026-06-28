// Package notification provides the NotificationService for managing
// Misskey notifications.
package notification

import (
	"context"

	"github.com/sysnote8main/migo/types"
)

// Service provides notification-related API operations.
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new NotificationService.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// ListRequest represents the request for i/notifications.
type ListRequest struct {
	Limit      *int      `json:"limit,omitempty"`
	SinceID    *types.ID `json:"sinceId,omitempty"`
	UntilID    *types.ID `json:"untilId,omitempty"`
	SinceDate  *int64    `json:"sinceDate,omitempty"`
	UntilDate  *int64    `json:"untilDate,omitempty"`
	MarkAsRead *bool     `json:"markAsRead,omitempty"`
}

// CreateRequest represents the request for notifications/create
// (admin push notification).
type CreateRequest struct {
	Body   string  `json:"body"`
	Header *string `json:"header,omitempty"`
	Icon   *string `json:"icon,omitempty"`
}

// ---------------------------------------------------------------------------
// API methods
// ---------------------------------------------------------------------------

// List fetches the authenticated user's notifications.
func (s *Service) List(ctx context.Context, req *ListRequest) ([]*types.Notification, error) {
	var notifications []*types.Notification
	if err := s.cli.Do(ctx, "/i/notifications", req, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

// ListGrouped fetches the authenticated user's notifications grouped.
func (s *Service) ListGrouped(ctx context.Context, req *ListRequest) ([]*types.Notification, error) {
	var notifications []*types.Notification
	if err := s.cli.Do(ctx, "/i/notifications-grouped", req, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

// MarkAllAsRead marks all notifications as read.
func (s *Service) MarkAllAsRead(ctx context.Context) error {
	return s.cli.Do(ctx, "/notifications/mark-all-as-read", nil, nil)
}

// Create creates a notification (admin push to all users).
func (s *Service) Create(ctx context.Context, req *CreateRequest) error {
	return s.cli.Do(ctx, "/notifications/create", req, nil)
}

// Flush processes pending notification emails.
func (s *Service) Flush(ctx context.Context) error {
	return s.cli.Do(ctx, "/notifications/flush", nil, nil)
}

// TestNotification sends a test notification to the authenticated user.
func (s *Service) TestNotification(ctx context.Context) error {
	return s.cli.Do(ctx, "/notifications/test-notification", nil, nil)
}
