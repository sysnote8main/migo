package chat

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sysnote8main/migo/types"
)

type mockDoer struct {
	t        *testing.T
	wantPath string
	wantReq  any
	resp     any
}

func (m *mockDoer) Do(ctx context.Context, path string, req, resp any) error {
	if m.wantPath != "" && path != m.wantPath {
		m.t.Errorf("path = %q, want %q", path, m.wantPath)
	}
	if m.wantReq != nil {
		gotJSON, _ := json.Marshal(req)
		wantJSON, _ := json.Marshal(m.wantReq)
		if string(gotJSON) != string(wantJSON) {
			m.t.Errorf("req = %s, want %s", gotJSON, wantJSON)
		}
	}
	if m.resp != nil && resp != nil {
		data, _ := json.Marshal(m.resp)
		json.Unmarshal(data, resp)
	}
	return nil
}

func TestService_History(t *testing.T) {
	resp := []*types.ChatMessage{{ID: "m1"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/history",
		resp:     &resp,
	})
	msgs, err := svc.History(context.Background(), &HistoryRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 {
		t.Errorf("got %d messages", len(msgs))
	}
}

func TestService_SendToUser(t *testing.T) {
	wantReq := &MessageCreateToUserRequest{
		Text:     strPtr("hello"),
		ToUserID: "user1",
	}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/messages/create-to-user",
		wantReq:  wantReq,
		resp:     &types.ChatMessage{ID: "m1", Text: strPtr("hello")},
	})
	msg, err := svc.SendToUser(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if msg.ID != "m1" {
		t.Errorf("ID = %q", msg.ID)
	}
}

func TestService_SendToRoom(t *testing.T) {
	wantReq := &MessageCreateToRoomRequest{
		Text:     strPtr("hello room"),
		ToRoomID: "room1",
	}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/messages/create-to-room",
		wantReq:  wantReq,
		resp:     &types.ChatMessage{ID: "m2"},
	})
	msg, err := svc.SendToRoom(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if msg.ID != "m2" {
		t.Errorf("ID = %q", msg.ID)
	}
}

func TestService_DeleteMessage(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/messages/delete",
	})
	if err := svc.DeleteMessage(context.Background(), "m1"); err != nil {
		t.Fatal(err)
	}
}

func TestService_ShowMessage(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/messages/show",
		resp:     &types.ChatMessage{ID: "m1"},
	})
	msg, err := svc.ShowMessage(context.Background(), "m1")
	if err != nil {
		t.Fatal(err)
	}
	if msg.ID != "m1" {
		t.Errorf("ID = %q", msg.ID)
	}
}

func TestService_ReactToMessage(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/messages/react",
	})
	if err := svc.ReactToMessage(context.Background(), "m1", "👍"); err != nil {
		t.Fatal(err)
	}
}

func TestService_CreateRoom(t *testing.T) {
	wantReq := &RoomCreateRequest{Name: "general"}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/rooms/create",
		wantReq:  wantReq,
		resp:     &types.ChatRoom{ID: "room1", Name: "general"},
	})
	room, err := svc.CreateRoom(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if room.Name != "general" {
		t.Errorf("Name = %q", room.Name)
	}
}

func TestService_ShowRoom(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/rooms/show",
		resp:     &types.ChatRoom{ID: "room1"},
	})
	room, err := svc.ShowRoom(context.Background(), "room1")
	if err != nil {
		t.Fatal(err)
	}
	if room.ID != "room1" {
		t.Errorf("ID = %q", room.ID)
	}
}

func TestService_JoinLeaveRoom(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/rooms/join",
	})
	if err := svc.JoinRoom(context.Background(), "room1"); err != nil {
		t.Fatal(err)
	}
}

func TestService_ReadAll(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/chat/read-all",
	})
	if err := svc.ReadAll(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func strPtr(s string) *string { return &s }
