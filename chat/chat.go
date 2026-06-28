// Package chat provides the ChatService for direct messaging and chat rooms.
package chat

import (
	"context"

	"github.com/sysnote8main/migo/types"
)

// Service provides chat-related API operations.
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new ChatService.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// HistoryRequest represents the request for chat/history.
type HistoryRequest struct {
	Limit *int  `json:"limit,omitempty"`
	Room  *bool `json:"room,omitempty"`
}

// MessageCreateToUserRequest represents the request for chat/messages/create-to-user.
type MessageCreateToUserRequest struct {
	Text     *string   `json:"text,omitempty"`
	FileID   *types.ID `json:"fileId,omitempty"`
	ToUserID types.ID  `json:"toUserId"`
}

// MessageCreateToRoomRequest represents the request for chat/messages/create-to-room.
type MessageCreateToRoomRequest struct {
	Text     *string   `json:"text,omitempty"`
	FileID   *types.ID `json:"fileId,omitempty"`
	ToRoomID types.ID  `json:"toRoomId"`
}

// MessageDeleteRequest represents the request for chat/messages/delete.
type MessageDeleteRequest struct {
	MessageID types.ID `json:"messageId"`
}

// MessageReactRequest represents the request for chat/messages/react.
type MessageReactRequest struct {
	MessageID types.ID `json:"messageId"`
	Reaction  string   `json:"reaction"`
}

// MessageTimelineRequest represents common timeline pagination.
type MessageTimelineRequest struct {
	Limit     *int      `json:"limit,omitempty"`
	SinceID   *types.ID `json:"sinceId,omitempty"`
	UntilID   *types.ID `json:"untilId,omitempty"`
	SinceDate *int64    `json:"sinceDate,omitempty"`
	UntilDate *int64    `json:"untilDate,omitempty"`
}

// UserTimelineRequest extends timeline with a user filter.
type UserTimelineRequest struct {
	MessageTimelineRequest
	UserID types.ID `json:"userId"`
}

// RoomTimelineRequest extends timeline with a room filter.
type RoomTimelineRequest struct {
	MessageTimelineRequest
	RoomID types.ID `json:"roomId"`
}

// MessageSearchRequest represents the request for chat/messages/search.
type MessageSearchRequest struct {
	Query  string    `json:"query"`
	Limit  *int      `json:"limit,omitempty"`
	UserID *types.ID `json:"userId,omitempty"`
	RoomID *types.ID `json:"roomId,omitempty"`
}

// RoomCreateRequest represents the request for chat/rooms/create.
type RoomCreateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// RoomUpdateRequest represents the request for chat/rooms/update.
type RoomUpdateRequest struct {
	RoomID      types.ID `json:"roomId"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
}

// RoomInviteCreateRequest represents the request for invitations/create.
type RoomInviteCreateRequest struct {
	RoomID types.ID `json:"roomId"`
	UserID types.ID `json:"userId"`
}

// RoomInviteActionRequest represents accept/ignore requests.
type RoomInviteActionRequest struct {
	RoomID types.ID `json:"roomId"`
}

// RoomMembersRequest represents the request for rooms/members.
type RoomMembersRequest struct {
	MessageTimelineRequest
	RoomID types.ID `json:"roomId"`
}

// RoomMuteRequest represents the request for rooms/mute.
type RoomMuteRequest struct {
	RoomID types.ID `json:"roomId"`
	Mute   bool     `json:"mute"`
}

// ---------------------------------------------------------------------------
// API methods — Messages
// ---------------------------------------------------------------------------

// History returns the list of recent chat conversations.
func (s *Service) History(ctx context.Context, req *HistoryRequest) ([]*types.ChatMessage, error) {
	var msgs []*types.ChatMessage
	if err := s.cli.Do(ctx, "/chat/history", req, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// SendToUser sends a direct message to a user.
func (s *Service) SendToUser(ctx context.Context, req *MessageCreateToUserRequest) (*types.ChatMessage, error) {
	var msg types.ChatMessage
	if err := s.cli.Do(ctx, "/chat/messages/create-to-user", req, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// SendToRoom sends a message to a chat room.
func (s *Service) SendToRoom(ctx context.Context, req *MessageCreateToRoomRequest) (*types.ChatMessage, error) {
	var msg types.ChatMessage
	if err := s.cli.Do(ctx, "/chat/messages/create-to-room", req, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// DeleteMessage deletes a chat message.
func (s *Service) DeleteMessage(ctx context.Context, messageID types.ID) error {
	req := &MessageDeleteRequest{MessageID: messageID}
	return s.cli.Do(ctx, "/chat/messages/delete", req, nil)
}

// ShowMessage shows a single chat message.
func (s *Service) ShowMessage(ctx context.Context, messageID types.ID) (*types.ChatMessage, error) {
	req := struct {
		MessageID types.ID `json:"messageId"`
	}{MessageID: messageID}
	var msg types.ChatMessage
	if err := s.cli.Do(ctx, "/chat/messages/show", req, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ReactToMessage adds a reaction to a chat message.
func (s *Service) ReactToMessage(ctx context.Context, messageID types.ID, reaction string) error {
	req := &MessageReactRequest{MessageID: messageID, Reaction: reaction}
	return s.cli.Do(ctx, "/chat/messages/react", req, nil)
}

// UnreactMessage removes a reaction from a chat message.
func (s *Service) UnreactMessage(ctx context.Context, messageID types.ID, reaction string) error {
	req := &MessageReactRequest{MessageID: messageID, Reaction: reaction}
	return s.cli.Do(ctx, "/chat/messages/unreact", req, nil)
}

// UserTimeline fetches the message timeline with a specific user.
func (s *Service) UserTimeline(ctx context.Context, req *UserTimelineRequest) ([]*types.ChatMessage, error) {
	var msgs []*types.ChatMessage
	if err := s.cli.Do(ctx, "/chat/messages/user-timeline", req, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// RoomTimeline fetches the message timeline in a specific room.
func (s *Service) RoomTimeline(ctx context.Context, req *RoomTimelineRequest) ([]*types.ChatMessage, error) {
	var msgs []*types.ChatMessage
	if err := s.cli.Do(ctx, "/chat/messages/room-timeline", req, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// SearchMessages searches chat messages.
func (s *Service) SearchMessages(ctx context.Context, req *MessageSearchRequest) ([]*types.ChatMessage, error) {
	var msgs []*types.ChatMessage
	if err := s.cli.Do(ctx, "/chat/messages/search", req, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// ReadAll marks all chat messages as read.
func (s *Service) ReadAll(ctx context.Context) error {
	return s.cli.Do(ctx, "/chat/read-all", nil, nil)
}

// ---------------------------------------------------------------------------
// API methods — Rooms
// ---------------------------------------------------------------------------

// CreateRoom creates a new chat room.
func (s *Service) CreateRoom(ctx context.Context, req *RoomCreateRequest) (*types.ChatRoom, error) {
	var room types.ChatRoom
	if err := s.cli.Do(ctx, "/chat/rooms/create", req, &room); err != nil {
		return nil, err
	}
	return &room, nil
}

// UpdateRoom updates a chat room.
func (s *Service) UpdateRoom(ctx context.Context, req *RoomUpdateRequest) (*types.ChatRoom, error) {
	var room types.ChatRoom
	if err := s.cli.Do(ctx, "/chat/rooms/update", req, &room); err != nil {
		return nil, err
	}
	return &room, nil
}

// DeleteRoom deletes a chat room.
func (s *Service) DeleteRoom(ctx context.Context, roomID types.ID) error {
	req := struct {
		RoomID types.ID `json:"roomId"`
	}{RoomID: roomID}
	return s.cli.Do(ctx, "/chat/rooms/delete", req, nil)
}

// ShowRoom shows a chat room details.
func (s *Service) ShowRoom(ctx context.Context, roomID types.ID) (*types.ChatRoom, error) {
	req := struct {
		RoomID types.ID `json:"roomId"`
	}{RoomID: roomID}
	var room types.ChatRoom
	if err := s.cli.Do(ctx, "/chat/rooms/show", req, &room); err != nil {
		return nil, err
	}
	return &room, nil
}

// JoinRoom joins a chat room.
func (s *Service) JoinRoom(ctx context.Context, roomID types.ID) error {
	req := struct {
		RoomID types.ID `json:"roomId"`
	}{RoomID: roomID}
	return s.cli.Do(ctx, "/chat/rooms/join", req, nil)
}

// LeaveRoom leaves a chat room.
func (s *Service) LeaveRoom(ctx context.Context, roomID types.ID) error {
	req := struct {
		RoomID types.ID `json:"roomId"`
	}{RoomID: roomID}
	return s.cli.Do(ctx, "/chat/rooms/leave", req, nil)
}

// MuteRoom toggles muting a chat room.
func (s *Service) MuteRoom(ctx context.Context, roomID types.ID, mute bool) error {
	req := &RoomMuteRequest{RoomID: roomID, Mute: mute}
	return s.cli.Do(ctx, "/chat/rooms/mute", req, nil)
}

// JoinedRooms lists rooms the authenticated user has joined.
func (s *Service) JoinedRooms(ctx context.Context, req *MessageTimelineRequest) ([]*types.ChatRoom, error) {
	var rooms []*types.ChatRoom
	if err := s.cli.Do(ctx, "/chat/rooms/joining", req, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

// OwnedRooms lists rooms the authenticated user owns.
func (s *Service) OwnedRooms(ctx context.Context, req *MessageTimelineRequest) ([]*types.ChatRoom, error) {
	var rooms []*types.ChatRoom
	if err := s.cli.Do(ctx, "/chat/rooms/owned", req, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

// RoomMembers lists the members of a chat room.
func (s *Service) RoomMembers(ctx context.Context, req *RoomMembersRequest) ([]*types.UserLite, error) {
	var users []*types.UserLite
	if err := s.cli.Do(ctx, "/chat/rooms/members", req, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// ---------------------------------------------------------------------------
// API methods — Room invitations
// ---------------------------------------------------------------------------

// CreateInvitation invites a user to a chat room.
func (s *Service) CreateInvitation(ctx context.Context, roomID, userID types.ID) error {
	req := &RoomInviteCreateRequest{RoomID: roomID, UserID: userID}
	return s.cli.Do(ctx, "/chat/rooms/invitations/create", req, nil)
}

// IgnoreInvitation ignores a room invitation.
func (s *Service) IgnoreInvitation(ctx context.Context, roomID types.ID) error {
	req := &RoomInviteActionRequest{RoomID: roomID}
	return s.cli.Do(ctx, "/chat/rooms/invitations/ignore", req, nil)
}

// InvitationsInbox lists pending room invitations for the user.
func (s *Service) InvitationsInbox(ctx context.Context, req *MessageTimelineRequest) ([]*types.ChatRoomInvitation, error) {
	var invs []*types.ChatRoomInvitation
	if err := s.cli.Do(ctx, "/chat/rooms/invitations/inbox", req, &invs); err != nil {
		return nil, err
	}
	return invs, nil
}

// InvitationsOutbox lists sent room invitations from a room.
func (s *Service) InvitationsOutbox(ctx context.Context, roomID types.ID, req *MessageTimelineRequest) ([]*types.ChatRoom, error) {
	// Combine roomID with pagination
	fullReq := struct {
		RoomID    types.ID  `json:"roomId"`
		Limit     *int      `json:"limit,omitempty"`
		SinceID   *types.ID `json:"sinceId,omitempty"`
		UntilID   *types.ID `json:"untilId,omitempty"`
		SinceDate *int64    `json:"sinceDate,omitempty"`
		UntilDate *int64    `json:"untilDate,omitempty"`
	}{
		RoomID:    roomID,
		Limit:     req.Limit,
		SinceID:   req.SinceID,
		UntilID:   req.UntilID,
		SinceDate: req.SinceDate,
		UntilDate: req.UntilDate,
	}
	var rooms []*types.ChatRoom
	if err := s.cli.Do(ctx, "/chat/rooms/invitations/outbox", fullReq, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
