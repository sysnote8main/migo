// Package streaming provides real-time access to Misskey's Streaming API
// via WebSocket. It supports channel subscription, note capture, and
// automatic reconnection.
//
// Basic usage:
//
//	s := streaming.NewStream("wss://instance.net", "YOUR_TOKEN")
//	defer s.Close()
//
//	s.OnNote(func(note *types.Note) {
//	    fmt.Printf("@%s: %s\n", note.User.Username, *note.Text)
//	})
//
//	if err := s.Connect(ctx); err != nil {
//	    log.Fatal(err)
//	}
package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sysnote8main/migo/types"
)

// Defaults
const (
	defaultPingInterval = 30 * time.Second
	defaultReadLimit    = 65536 // 64KB
)

// Stream manages a WebSocket connection to the Misskey streaming API.
type Stream struct {
	url    string
	token  string
	dialer *websocket.Dialer

	mu     sync.RWMutex
	conn   *websocket.Conn
	subs   map[string]*subscription // id -> subscription
	nextID int
	closed bool

	// Auto-reconnect
	reconnect bool

	// Callbacks
	onNote         func(*types.Note)
	onNotification func(*types.Notification)
	onFollow       func(*types.UserLite)
	onUnfollow     func(*types.UserLite)
	onMention      func(*types.Note)
	onReply        func(*types.Note)
	onRenote       func(*types.Note)
	onReacted      func(*types.Note)

	// Note capture callbacks
	onNoteCapture map[types.ID]func(*types.Note)

	// Raw event callback for custom handling
	onEvent func(*Event)
}

// Option configures the Stream.
type Option func(*Stream)

// WithReconnect enables automatic reconnection on connection loss.
func WithReconnect(enabled bool) Option {
	return func(s *Stream) {
		s.reconnect = enabled
	}
}

// WithDialer sets a custom WebSocket dialer.
func WithDialer(d *websocket.Dialer) Option {
	return func(s *Stream) {
		s.dialer = d
	}
}

// NewStream creates a new Stream.
// The baseURL should be the instance URL (e.g., "wss://misskey.io" or "https://misskey.io").
func NewStream(baseURL, token string, opts ...Option) *Stream {
	// Convert http/https to ws/wss if needed
	wsURL := baseURL
	parsed, err := url.Parse(baseURL)
	if err == nil {
		switch parsed.Scheme {
		case "https":
			parsed.Scheme = "wss"
		case "http":
			parsed.Scheme = "ws"
		}
		wsURL = parsed.String()
	}

	s := &Stream{
		url:       fmt.Sprintf("%s/streaming?i=%s", wsURL, url.QueryEscape(token)),
		token:     token,
		dialer:    websocket.DefaultDialer,
		subs:      make(map[string]*subscription),
		reconnect: true,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// subscription represents an active channel subscription.
type subscription struct {
	id      string
	channel string
	params  map[string]any
}

// Event represents a raw event received from a channel.
type Event struct {
	ChannelID string
	Type      string
	Body      json.RawMessage
}

// Connect establishes the WebSocket connection and starts reading events.
// It blocks until the context is cancelled or a fatal error occurs.
func (s *Stream) Connect(ctx context.Context) error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return fmt.Errorf("streaming: stream is closed")
	}
	s.mu.Unlock()

	for {
		if err := s.connectOnce(ctx); err != nil {
			s.mu.RLock()
			rc := s.reconnect
			closed := s.closed
			s.mu.RUnlock()

			if closed || !rc {
				return err
			}
			log.Printf("streaming: reconnecting in 5s: %v", err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
			}
			continue
		}
		return nil
	}
}

func (s *Stream) connectOnce(ctx context.Context) error {
	conn, _, err := s.dialer.DialContext(ctx, s.url, nil)
	if err != nil {
		return fmt.Errorf("streaming: dial: %w", err)
	}

	s.mu.Lock()
	// If stream was closed while dialing, close and bail out
	if s.closed {
		s.mu.Unlock()
		conn.Close()
		return fmt.Errorf("streaming: stream closed during dial")
	}
	s.conn = conn
	// Resubscribe to active subscriptions (we hold the lock)
	for _, sub := range s.subs {
		writeJSONTo(conn, clientMessage{
			Type: "connect",
			Body: map[string]any{
				"channel": sub.channel,
				"id":      sub.id,
				"params":  sub.params,
			},
		})
	}
	s.mu.Unlock()

	conn.SetReadLimit(defaultReadLimit)
	conn.SetPingHandler(func(data string) error {
		return conn.WriteMessage(websocket.PongMessage, []byte(data))
	})

	// Read loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.mu.Lock()
			if s.conn == conn {
				s.conn = nil
			}
			s.mu.Unlock()
			return fmt.Errorf("streaming: read: %w", err)
		}
		s.handleMessage(message)
	}
}

// Close closes the WebSocket connection and marks the stream as closed.
func (s *Stream) Close() error {
	s.mu.Lock()
	s.closed = true
	s.reconnect = false
	conn := s.conn
	s.conn = nil
	s.mu.Unlock()

	if conn != nil {
		// Close the underlying connection directly to unblock the read loop.
		// We avoid WriteMessage(CloseMessage) here because the write may block
		// if the peer is not reading (deadlock with the read loop goroutine).
		conn.Close()
	}
	return nil
}

// Subscribe subscribes to a channel with the given parameters.
// Returns a subscription ID that can be used with Unsubscribe.
func (s *Stream) Subscribe(channel string, params map[string]any) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return "", fmt.Errorf("streaming: stream is closed")
	}

	id := fmt.Sprintf("sub_%d", s.nextID)
	s.nextID++

	sub := &subscription{
		id:      id,
		channel: channel,
		params:  params,
	}
	s.subs[id] = sub

	if s.conn != nil {
		// Write directly without acquiring lock (we already hold it)
		if err := s.conn.WriteJSON(clientMessage{
			Type: "connect",
			Body: map[string]any{
				"channel": sub.channel,
				"id":      sub.id,
				"params":  sub.params,
			},
		}); err != nil {
			log.Printf("streaming: write connect error: %v", err)
		}
	}

	return id, nil
}

// Unsubscribe unsubscribes from a channel by its subscription ID.
func (s *Stream) Unsubscribe(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subs, id)
	if s.conn != nil {
		// Write directly without acquiring lock (we already hold it)
		if err := s.conn.WriteJSON(clientMessage{
			Type: "disconnect",
			Body: map[string]any{
				"id": id,
			},
		}); err != nil {
			log.Printf("streaming: write disconnect error: %v", err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Convenience subscriptions
// ---------------------------------------------------------------------------

// OnNote registers a callback for home timeline notes.
// It subscribes to the homeTimeline channel automatically.
// Returns the subscription ID for unsubscription.
func (s *Stream) OnNote(cb func(*types.Note)) (string, error) {
	s.mu.Lock()
	s.onNote = cb
	s.mu.Unlock()
	return s.Subscribe("homeTimeline", nil)
}

// OnNotification registers a callback for notifications.
// It subscribes to the main channel automatically.
// Returns the subscription ID for unsubscription.
func (s *Stream) OnNotification(cb func(*types.Notification)) (string, error) {
	s.mu.Lock()
	s.onNotification = cb
	s.mu.Unlock()
	return s.Subscribe("main", nil)
}

// OnFollow registers a callback for follow events.
func (s *Stream) OnFollow(cb func(*types.UserLite)) (string, error) {
	s.mu.Lock()
	s.onFollow = cb
	s.mu.Unlock()
	return s.Subscribe("main", nil)
}

// OnMention registers a callback for mention events.
func (s *Stream) OnMention(cb func(*types.Note)) (string, error) {
	s.mu.Lock()
	s.onMention = cb
	s.mu.Unlock()
	return s.Subscribe("main", nil)
}

// OnReply registers a callback for reply events.
func (s *Stream) OnReply(cb func(*types.Note)) (string, error) {
	s.mu.Lock()
	s.onReply = cb
	s.mu.Unlock()
	return s.Subscribe("main", nil)
}

// SubscribeHomeTimeline subscribes to the home timeline.
func (s *Stream) SubscribeHomeTimeline(cb func(*types.Note)) (string, error) {
	return s.OnNote(cb)
}

// SubscribeLocalTimeline subscribes to the local timeline.
func (s *Stream) SubscribeLocalTimeline(cb func(*types.Note)) (string, error) {
	id, err := s.Subscribe("localTimeline", nil)
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.onNote = cb
	s.mu.Unlock()
	return id, nil
}

// SubscribeGlobalTimeline subscribes to the global timeline.
func (s *Stream) SubscribeGlobalTimeline(cb func(*types.Note)) (string, error) {
	id, err := s.Subscribe("globalTimeline", nil)
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.onNote = cb
	s.mu.Unlock()
	return id, nil
}

// SubscribeHybridTimeline subscribes to the hybrid timeline.
func (s *Stream) SubscribeHybridTimeline(cb func(*types.Note)) (string, error) {
	id, err := s.Subscribe("hybridTimeline", nil)
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.onNote = cb
	s.mu.Unlock()
	return id, nil
}

// OnEvent registers a raw event callback for custom handling.
func (s *Stream) OnEvent(cb func(*Event)) {
	s.mu.Lock()
	s.onEvent = cb
	s.mu.Unlock()
}

// ---------------------------------------------------------------------------
// Note capture
// ---------------------------------------------------------------------------

// CaptureNote subscribes to real-time updates for a specific note.
// The callback is called when the note is updated (reactions, etc.).
func (s *Stream) CaptureNote(noteID types.ID, cb func(*types.Note)) error {
	_, err := s.Subscribe("note", map[string]any{
		"noteId": noteID,
	})
	if err != nil {
		return err
	}
	// Store callback per noteID
	s.mu.Lock()
	if s.onNoteCapture == nil {
		s.onNoteCapture = make(map[types.ID]func(*types.Note))
	}
	s.onNoteCapture[noteID] = cb
	s.mu.Unlock()
	return nil
}

// UncaptureNote stops capturing updates for a specific note.
func (s *Stream) UncaptureNote(noteID types.ID) error {
	s.mu.Lock()
	delete(s.onNoteCapture, noteID)
	s.mu.Unlock()
	return nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// clientMessage is the JSON structure sent to the server.
type clientMessage struct {
	Type string `json:"type"`
	Body any    `json:"body"`
}

// serverMessage is the JSON structure received from the server.
type serverMessage struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

// channelBody is the body of a channel event.
type channelBody struct {
	ID   string          `json:"id"`
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

func writeJSONTo(conn *websocket.Conn, v any) {
	if conn == nil {
		return
	}
	if err := conn.WriteJSON(v); err != nil {
		log.Printf("streaming: write error: %v", err)
	}
}

func (s *Stream) writeJSON(v any) {
	s.mu.RLock()
	conn := s.conn
	s.mu.RUnlock()
	if conn != nil {
		writeJSONTo(conn, v)
	}
}

func (s *Stream) handleMessage(data []byte) {
	var msg serverMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("streaming: parse error: %v", err)
		return
	}

	switch msg.Type {
	case "channel":
		var cb channelBody
		if err := json.Unmarshal(msg.Body, &cb); err != nil {
			log.Printf("streaming: parse channel body: %v", err)
			return
		}
		s.handleChannelEvent(&cb)

	case "noteCaptured":
		// Note capture event
		var note types.Note
		if err := json.Unmarshal(msg.Body, &note); err != nil {
			log.Printf("streaming: parse captured note: %v", err)
			return
		}
		s.mu.RLock()
		cb := s.onNoteCapture[note.ID]
		s.mu.RUnlock()
		if cb != nil {
			cb(&note)
		}

	default:
		log.Printf("streaming: unknown message type: %s", msg.Type)
	}
}

func (s *Stream) handleChannelEvent(cb *channelBody) {
	// Dispatch to raw event callback
	s.mu.RLock()
	onEvent := s.onEvent
	s.mu.RUnlock()
	if onEvent != nil {
		onEvent(&Event{
			ChannelID: cb.ID,
			Type:      cb.Type,
			Body:      cb.Body,
		})
	}

	switch cb.Type {
	case "note":
		var note types.Note
		if err := json.Unmarshal(cb.Body, &note); err != nil {
			log.Printf("streaming: parse note: %v", err)
			return
		}
		s.mu.RLock()
		fn := s.onNote
		s.mu.RUnlock()
		if fn != nil {
			fn(&note)
		}

	case "notification":
		var notif types.Notification
		if err := json.Unmarshal(cb.Body, &notif); err != nil {
			log.Printf("streaming: parse notification: %v", err)
			return
		}
		s.mu.RLock()
		fn := s.onNotification
		s.mu.RUnlock()
		if fn != nil {
			fn(&notif)
		}

	case "follow":
		var user types.UserLite
		if err := json.Unmarshal(cb.Body, &user); err != nil {
			log.Printf("streaming: parse follow: %v", err)
			return
		}
		s.mu.RLock()
		fn := s.onFollow
		s.mu.RUnlock()
		if fn != nil {
			fn(&user)
		}

	case "unfollow":
		var user types.UserLite
		if err := json.Unmarshal(cb.Body, &user); err != nil {
			log.Printf("streaming: parse unfollow: %v", err)
			return
		}
		s.mu.RLock()
		fn := s.onUnfollow
		s.mu.RUnlock()
		if fn != nil {
			fn(&user)
		}

	case "mention":
		var note types.Note
		if err := json.Unmarshal(cb.Body, &note); err != nil {
			log.Printf("streaming: parse mention: %v", err)
			return
		}
		s.mu.RLock()
		fn := s.onMention
		s.mu.RUnlock()
		if fn != nil {
			fn(&note)
		}

	case "reply":
		var note types.Note
		if err := json.Unmarshal(cb.Body, &note); err != nil {
			log.Printf("streaming: parse reply: %v", err)
			return
		}
		s.mu.RLock()
		fn := s.onReply
		s.mu.RUnlock()
		if fn != nil {
			fn(&note)
		}

	case "renote":
		var note types.Note
		if err := json.Unmarshal(cb.Body, &note); err != nil {
			log.Printf("streaming: parse renote: %v", err)
			return
		}
		s.mu.RLock()
		fn := s.onRenote
		s.mu.RUnlock()
		if fn != nil {
			fn(&note)
		}

	case "reacted":
		var note types.Note
		if err := json.Unmarshal(cb.Body, &note); err != nil {
			log.Printf("streaming: parse reacted: %v", err)
			return
		}
		s.mu.RLock()
		fn := s.onReacted
		s.mu.RUnlock()
		if fn != nil {
			fn(&note)
		}
	}
}
