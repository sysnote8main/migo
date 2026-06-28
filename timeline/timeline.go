// Package timeline provides access to the Misskey timeline endpoints.
package timeline

import (
	"context"

	"github.com/sysnote8main/migo/types"
)

// Service provides timeline-related API operations.
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new TimelineService.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// Request contains common timeline pagination parameters.
type Request struct {
	Limit                 *int      `json:"limit,omitempty"`
	SinceID               *types.ID `json:"sinceId,omitempty"`
	UntilID               *types.ID `json:"untilId,omitempty"`
	SinceDate             *int64    `json:"sinceDate,omitempty"`
	UntilDate             *int64    `json:"untilDate,omitempty"`
	AllowPartial          *bool     `json:"allowPartial,omitempty"`
	IncludeMyRenotes      *bool     `json:"includeMyRenotes,omitempty"`
	IncludeRenotedMyNotes *bool     `json:"includeRenotedMyNotes,omitempty"`
	IncludeLocalRenotes   *bool     `json:"includeLocalRenotes,omitempty"`
	WithFiles             *bool     `json:"withFiles,omitempty"`
	WithRenotes           *bool     `json:"withRenotes,omitempty"`
}

// ChannelRequest extends Request with a channel ID filter.
type ChannelRequest struct {
	Request
	ChannelID types.ID `json:"channelId"`
}

// ---------------------------------------------------------------------------
// API methods
// ---------------------------------------------------------------------------

// Home fetches the home timeline (notes from users you follow).
// Requires authentication (read:account permission).
func (s *Service) Home(ctx context.Context, req *Request) ([]*types.Note, error) {
	var notes []*types.Note
	if err := s.cli.Do(ctx, "/notes/timeline", req, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// Local fetches the local timeline (notes from the local instance).
// Does not require authentication.
func (s *Service) Local(ctx context.Context, req *Request) ([]*types.Note, error) {
	var notes []*types.Note
	if err := s.cli.Do(ctx, "/notes/local-timeline", req, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// Hybrid fetches the hybrid timeline (local + remote follow notes).
// Requires authentication.
func (s *Service) Hybrid(ctx context.Context, req *Request) ([]*types.Note, error) {
	var notes []*types.Note
	if err := s.cli.Do(ctx, "/notes/hybrid-timeline", req, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// Global fetches the global timeline (all known notes).
// Does not require authentication.
func (s *Service) Global(ctx context.Context, req *Request) ([]*types.Note, error) {
	var notes []*types.Note
	if err := s.cli.Do(ctx, "/notes/global-timeline", req, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// Channel fetches a channel's timeline.
func (s *Service) Channel(ctx context.Context, req *ChannelRequest) ([]*types.Note, error) {
	var notes []*types.Note
	if err := s.cli.Do(ctx, "/channels/timeline", req, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}
