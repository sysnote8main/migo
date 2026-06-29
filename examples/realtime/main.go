// Command streaming demonstrates real-time event streaming via WebSocket.
//
// The streaming API receives new notes and notifications as they happen,
// without polling. This example subscribes to the home timeline and
// notification events, runs for a set duration, then exits.
//
// Usage:
//
//	export MIGO_URL="https://your-instance.net"
//	# Set MIGO_TOKEN via environment variable
//	go run ./examples/realtime
//
// The MIGO_URL should be the instance root (not the /api path).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sysnote8main/migo/streaming"
	"github.com/sysnote8main/migo/types"
)

func main() {
	baseURL := os.Getenv("MIGO_URL")
	if baseURL == "" {
		baseURL = "https://misskey.io"
	}

	token := os.Getenv("MIGO_TOKEN")
	if token == "" {
		log.Fatal("MIGO_TOKEN is required")
	}

	// -------------------------------------------------------------------
	// 1. Create a stream and configure event handlers
	// -------------------------------------------------------------------
	s := streaming.NewStream(
		baseURL,
		token,
		streaming.WithReconnect(true), // auto-reconnect on disconnect
	)

	// Register a callback for home timeline notes.
	subNote, err := s.OnNote(func(note *types.Note) {
		fmt.Printf("[NOTE   ] @%s: %s\n",
			note.User.Username, ellipsis(ptrStr(note.Text), 100))
	})
	if err != nil {
		log.Fatalf("failed to subscribe to notes: %v", err)
	}
	fmt.Printf("Subscribed to home timeline (id=%s)\n", subNote)

	// Register a callback for notifications.
	subNotif, err := s.OnNotification(func(notif *types.Notification) {
		fmt.Printf("[NOTIF  ] type=%s", notif.Type)
		if notif.User != nil {
			fmt.Printf(" from=@%s", notif.User.Username)
		}
		if notif.Reaction != nil {
			fmt.Printf(" reaction=%s", *notif.Reaction)
		}
		fmt.Println()
	})
	if err != nil {
		log.Fatalf("failed to subscribe to notifications: %v", err)
	}
	fmt.Printf("Subscribed to notifications (id=%s)\n", subNotif)

	// Register a callback for mentions.
	subMention, err := s.OnMention(func(note *types.Note) {
		fmt.Printf("[MENTION] @%s mentioned you: %s\n",
			note.User.Username, ellipsis(ptrStr(note.Text), 100))
	})
	if err != nil {
		log.Fatalf("failed to subscribe to mentions: %v", err)
	}
	_ = subMention

	// -------------------------------------------------------------------
	// 2. Connect and listen for events
	// -------------------------------------------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Handle SIGINT for graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	fmt.Println("\nListening for events (30s timeout). Press Ctrl+C to stop early.")
	fmt.Println("--------------------------------------------------")

	// Connect blocks until the context is cancelled or a fatal error occurs.
	if err := s.Connect(ctx); err != nil {
		// Ignore context.Cancelled — it's our graceful shutdown signal.
		if ctx.Err() == nil {
			log.Printf("stream ended: %v", err)
		}
	}

	// -------------------------------------------------------------------
	// 3. Clean up subscriptions
	// -------------------------------------------------------------------
	if err := s.Unsubscribe(subNote); err != nil {
		log.Printf("failed to unsubscribe note: %v", err)
	}
	if err := s.Unsubscribe(subNotif); err != nil {
		log.Printf("failed to unsubscribe notification: %v", err)
	}

	if err := s.Close(); err != nil {
		log.Printf("close error: %v", err)
	}
	fmt.Println("\nDone.")
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ellipsis(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}

// Ensure types is referenced.
var _ = types.VisibilityPublic
