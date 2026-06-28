// Package notes provides the NoteService for managing Misskey notes.
package notes

import (
	"context"

	"github.com/sysnote8main/migo/types"
)

// Service provides note-related API operations.
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new NoteService.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// CreateRequest represents the request body for notes/create.
type CreateRequest struct {
	Visibility         types.NoteVisibility `json:"visibility,omitempty"`
	VisibleUserIds     []types.ID           `json:"visibleUserIds,omitempty"`
	Cw                 *string              `json:"cw,omitempty"`
	LocalOnly          bool                 `json:"localOnly,omitempty"`
	ReactionAcceptance *string              `json:"reactionAcceptance,omitempty"`
	NoExtractMentions  bool                 `json:"noExtractMentions,omitempty"`
	NoExtractHashtags  bool                 `json:"noExtractHashtags,omitempty"`
	NoExtractEmojis    bool                 `json:"noExtractEmojis,omitempty"`
	ReplyID            *types.ID            `json:"replyId,omitempty"`
	RenoteID           *types.ID            `json:"renoteId,omitempty"`
	ChannelID          *types.ID            `json:"channelId,omitempty"`
	Text               *string              `json:"text,omitempty"`
	FileIDs            []types.ID           `json:"fileIds,omitempty"`
	MediaIDs           []types.ID           `json:"mediaIds,omitempty"`
	Poll               *CreatePollRequest   `json:"poll,omitempty"`
}

// CreatePollRequest represents a poll definition in a note create request.
type CreatePollRequest struct {
	Choices   []string `json:"choices"`
	Multiple  bool     `json:"multiple,omitempty"`
	ExpiresAt *int64   `json:"expiresAt,omitempty"`
}

// CreateResponse is the response from notes/create.
type CreateResponse struct {
	CreatedNote types.Note `json:"createdNote"`
}

// ShowRequest represents the request body for notes/show.
type ShowRequest struct {
	NoteID types.ID `json:"noteId"`
}

// DeleteRequest represents the request body for notes/delete.
type DeleteRequest struct {
	NoteID types.ID `json:"noteId"`
}

// SearchRequest represents the request body for notes/search.
type SearchRequest struct {
	Query     string    `json:"query"`
	SinceID   *types.ID `json:"sinceId,omitempty"`
	UntilID   *types.ID `json:"untilId,omitempty"`
	SinceDate *int64    `json:"sinceDate,omitempty"`
	UntilDate *int64    `json:"untilDate,omitempty"`
	Limit     *int      `json:"limit,omitempty"`
	Offset    *int      `json:"offset,omitempty"`
	Host      *string   `json:"host,omitempty"`
	UserID    *types.ID `json:"userId,omitempty"`
	ChannelID *types.ID `json:"channelId,omitempty"`
}

// ReactionCreateRequest represents the request for notes/reactions/create.
type ReactionCreateRequest struct {
	NoteID   types.ID `json:"noteId"`
	Reaction string   `json:"reaction"`
}

// ReactionDeleteRequest represents the request for notes/reactions/delete.
type ReactionDeleteRequest struct {
	NoteID types.ID `json:"noteId"`
}

// ---------------------------------------------------------------------------
// API methods
// ---------------------------------------------------------------------------

// Create creates a new note.
func (s *Service) Create(ctx context.Context, req *CreateRequest) (*types.Note, error) {
	var resp CreateResponse
	if err := s.cli.Do(ctx, "/notes/create", req, &resp); err != nil {
		return nil, err
	}
	return &resp.CreatedNote, nil
}

// Show fetches a single note by its ID.
func (s *Service) Show(ctx context.Context, noteID types.ID) (*types.Note, error) {
	req := &ShowRequest{NoteID: noteID}
	var note types.Note
	if err := s.cli.Do(ctx, "/notes/show", req, &note); err != nil {
		return nil, err
	}
	return &note, nil
}

// Delete deletes a note by its ID.
func (s *Service) Delete(ctx context.Context, noteID types.ID) error {
	req := &DeleteRequest{NoteID: noteID}
	return s.cli.Do(ctx, "/notes/delete", req, nil)
}

// Search searches notes by query.
func (s *Service) Search(ctx context.Context, req *SearchRequest) ([]*types.Note, error) {
	var notes []*types.Note
	if err := s.cli.Do(ctx, "/notes/search", req, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// CreateReaction adds a reaction to a note.
func (s *Service) CreateReaction(ctx context.Context, noteID types.ID, reaction string) error {
	req := &ReactionCreateRequest{NoteID: noteID, Reaction: reaction}
	return s.cli.Do(ctx, "/notes/reactions/create", req, nil)
}

// DeleteReaction removes a reaction from a note.
func (s *Service) DeleteReaction(ctx context.Context, noteID types.ID) error {
	req := &ReactionDeleteRequest{NoteID: noteID}
	return s.cli.Do(ctx, "/notes/reactions/delete", req, nil)
}
