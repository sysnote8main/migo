package streaming

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sysnote8main/migo/types"
)

// Test WebSocket server that simulates a Misskey streaming server.
type testWSServer struct {
	*httptest.Server
	url      string
	connCh   chan *websocket.Conn
	upgrader websocket.Upgrader
}

func newTestWSServer(t *testing.T) *testWSServer {
	ts := &testWSServer{
		connCh: make(chan *websocket.Conn, 1),
	}
	ts.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := ts.upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("upgrade error: %v", err)
			return
		}
		ts.connCh <- conn
	}))
	ts.url = strings.Replace(ts.Server.URL, "http://", "ws://", 1)
	return ts
}

func (ts *testWSServer) waitConn(t *testing.T) *websocket.Conn {
	t.Helper()
	select {
	case conn := <-ts.connCh:
		return conn
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for WebSocket connection")
		return nil
	}
}

func TestNewStream(t *testing.T) {
	s := NewStream("wss://misskey.io", "test-token")
	if s == nil {
		t.Fatal("NewStream returned nil")
	}
	if !strings.Contains(s.url, "test-token") {
		t.Errorf("url does not contain token: %s", s.url)
	}
}

func TestNewStream_HTTPS(t *testing.T) {
	s := NewStream("https://misskey.io", "tok")
	if !strings.Contains(s.url, "wss://") {
		t.Errorf("https should become wss: %s", s.url)
	}
}

func TestStream_ConnectAndClose(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := NewStream(ts.url, "test-token", WithReconnect(false))

	// Wait for the server to accept before connecting (or handle connect in background)
	connected := make(chan struct{})
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Connect(ctx)
		close(connected)
	}()

	// Wait for server connection
	serverConn := ts.waitConn(t)

	// Wait a tiny bit for client to fully enter read loop
	time.Sleep(50 * time.Millisecond)

	// Close the stream from the client side
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	// Server connection should also close
	serverConn.Close()

	// Connect should return
	select {
	case err := <-errCh:
		t.Logf("Connect returned: %v", err)
	case <-connected:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for Connect to return")
	}
}

func TestStream_Subscribe(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := NewStream(ts.url, "tok", WithReconnect(false))

	go func() {
		s.Connect(ctx)
	}()

	serverConn := ts.waitConn(t)

	// Subscribe
	id, err := s.Subscribe("homeTimeline", nil)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("subscription ID is empty")
	}

	// Read connect message from server
	_, msg, err := serverConn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	var cm struct {
		Type string `json:"type"`
		Body struct {
			Channel string `json:"channel"`
			ID      string `json:"id"`
		} `json:"body"`
	}
	if err := json.Unmarshal(msg, &cm); err != nil {
		t.Fatal(err)
	}
	if cm.Type != "connect" {
		t.Errorf("type = %q, want connect", cm.Type)
	}
	if cm.Body.Channel != "homeTimeline" {
		t.Errorf("channel = %q", cm.Body.Channel)
	}
	if cm.Body.ID != id {
		t.Errorf("id = %q, want %q", cm.Body.ID, id)
	}

	s.Close()
}

func TestStream_OnNote(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := NewStream(ts.url, "tok", WithReconnect(false))
	defer s.Close()

	noteCh := make(chan *types.Note, 1)
	_, err := s.OnNote(func(note *types.Note) {
		noteCh <- note
	})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		s.Connect(ctx)
	}()

	serverConn := ts.waitConn(t)

	// Read the connect message from client
	_, _, _ = serverConn.ReadMessage()

	// Simulate a note event from server
	note := types.Note{ID: "note123", Text: strPtr("hello")}
	bodyData, _ := json.Marshal(note)

	event := map[string]any{
		"type": "channel",
		"body": map[string]any{
			"id":   "sub_0",
			"type": "note",
			"body": json.RawMessage(bodyData),
		},
	}
	eventData, _ := json.Marshal(event)
	if err := serverConn.WriteMessage(websocket.TextMessage, eventData); err != nil {
		t.Fatal(err)
	}

	select {
	case received := <-noteCh:
		if received.ID != "note123" {
			t.Errorf("note ID = %q", received.ID)
		}
		if *received.Text != "hello" {
			t.Errorf("text = %q", *received.Text)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for note callback")
	}
}

func TestStream_OnNotification(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := NewStream(ts.url, "tok", WithReconnect(false))
	defer s.Close()

	notifCh := make(chan *types.Notification, 1)
	_, err := s.OnNotification(func(notif *types.Notification) {
		notifCh <- notif
	})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		s.Connect(ctx)
	}()

	serverConn := ts.waitConn(t)
	_, _, _ = serverConn.ReadMessage()

	// Simulate a notification
	bodyData, _ := json.Marshal(types.Notification{
		ID:   "notif1",
		Type: types.NotificationTypeFollow,
	})

	event := map[string]any{
		"type": "channel",
		"body": map[string]any{
			"id":   "sub_0",
			"type": "notification",
			"body": json.RawMessage(bodyData),
		},
	}
	eventData, _ := json.Marshal(event)
	serverConn.WriteMessage(websocket.TextMessage, eventData)

	select {
	case received := <-notifCh:
		if received.ID != "notif1" {
			t.Errorf("notification ID = %q", received.ID)
		}
		if received.Type != types.NotificationTypeFollow {
			t.Errorf("type = %q", received.Type)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for notification callback")
	}
}

func TestStream_CaptureNote(t *testing.T) {
	ts := newTestWSServer(t)
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := NewStream(ts.url, "tok", WithReconnect(false))
	defer s.Close()

	noteCh := make(chan *types.Note, 1)
	err := s.CaptureNote("note123", func(note *types.Note) {
		noteCh <- note
	})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		s.Connect(ctx)
	}()

	serverConn := ts.waitConn(t)
	// Read the note capture subscribe message from client
	_, _, _ = serverConn.ReadMessage()

	// Simulate a note capture event
	note := types.Note{ID: "note123", Text: strPtr("updated")}
	noteData, _ := json.Marshal(note)

	event := map[string]any{
		"type": "noteCaptured",
		"body": json.RawMessage(noteData),
	}
	eventData, _ := json.Marshal(event)
	serverConn.WriteMessage(websocket.TextMessage, eventData)

	select {
	case received := <-noteCh:
		if received.ID != "note123" {
			t.Errorf("note ID = %q", received.ID)
		}
		if *received.Text != "updated" {
			t.Errorf("text = %q", *received.Text)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for capture callback")
	}
}

func TestStream_DoubleClose(t *testing.T) {
	s := NewStream("wss://localhost:1", "tok")
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestStream_SubscribeAfterClose(t *testing.T) {
	s := NewStream("wss://localhost:1", "tok")
	s.Close()
	_, err := s.Subscribe("homeTimeline", nil)
	if err == nil {
		t.Fatal("expected error subscribing to closed stream")
	}
}

func strPtr(s string) *string { return &s }
