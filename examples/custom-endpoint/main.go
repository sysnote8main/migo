// Command custom-endpoint demonstrates calling Misskey API endpoints
// that are not yet wrapped by migo's typed services.
//
// The Client.Do() method gives you full access to every Misskey API
// endpoint by sending raw JSON requests and decoding the response into
// any Go type, including map[string]any for ad-hoc access.
//
// Usage:
//
//	# Set via: export MIGO_TOKEN=\"<your-token>\"
//	export MIGO_URL="https://your-instance.net/api"
//	# Set via: export MIGO_TOKEN="<your-token>"
//	go run ./examples/custom-endpoint
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sysnote8main/migo"
	"github.com/sysnote8main/migo/auth"
)

func main() {
	ctx := context.Background()

	baseURL := envOrDefault("MIGO_URL", "https://misskey.io/api")
	token := os.Getenv("MIGO_TOKEN")
	if token == "" {
		log.Fatal("MIGO_TOKEN is required")
	}

	client := migo.NewClient(
		migo.WithBaseURL(baseURL),
		migo.WithToken(token),
	)

	// -------------------------------------------------------------------
	// 1. Fetch raw response with map[string]any
	// -------------------------------------------------------------------
	fmt.Println("=== 1. Raw map response ===")
	var raw map[string]any
	if err := client.Do(ctx, "/i", nil, &raw); err != nil {
		log.Fatalf("failed: %v", err)
	}
	prettyJSON(raw)

	// -------------------------------------------------------------------
	// 2. Custom request body, custom response struct
	// -------------------------------------------------------------------
	fmt.Println("\n=== 2. Custom struct response ===")

	// Define response structs for endpoints that don't have wrappers yet.
	var stats struct {
		NotesCount         int   `json:"notesCount"`
		UsersCount         int   `json:"usersCount"`
		Instances          int   `json:"instances"`
		OriginalNotesCount int   `json:"originalNotesCount"`
		OriginalUsersCount int   `json:"originalUsersCount"`
		DriveUsage         int64 `json:"driveUsage"`
		DriveCapacity      int64 `json:"driveCapacity"`
		OnlineUserCount    int   `json:"onlineUserCount"`
	}
	if err := client.Do(ctx, "/stats", nil, &stats); err != nil {
		log.Printf("stats failed: %v\n", err)
	} else {
		fmt.Printf("Instance stats:\n")
		fmt.Printf("  Notes:   %d\n", stats.NotesCount)
		fmt.Printf("  Users:   %d\n", stats.UsersCount)
		fmt.Printf("  Online:  %d\n", stats.OnlineUserCount)
		fmt.Printf("  Drive:   %d / %d bytes\n", stats.DriveUsage, stats.DriveCapacity)
	}

	// -------------------------------------------------------------------
	// 3. Pagination with typed parameters
	// -------------------------------------------------------------------
	fmt.Println("\n=== 3. Paginated request ===")

	type NotesRequest struct {
		Limit *int `json:"limit,omitempty"`
	}

	var latestNotes []struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		Text      *string   `json:"text"`
		User      struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	}
	if err := client.Do(ctx, "/notes/local-timeline", &NotesRequest{Limit: ptr(3)}, &latestNotes); err != nil {
		log.Printf("local timeline failed: %v\n", err)
	} else {
		fmt.Printf("Latest local notes (%d):\n", len(latestNotes))
		for _, n := range latestNotes {
			text := ""
			if n.Text != nil {
				text = *n.Text
			}
			if len(text) > 60 {
				text = text[:60] + "..."
			}
			fmt.Printf("  [%s] @%s: %s\n", n.ID[:8], n.User.Username, text)
		}
	}

	// -------------------------------------------------------------------
	// 4. Error handling with raw map
	// -------------------------------------------------------------------
	fmt.Println("\n=== 4. Error handling ===")
	var badResp map[string]any
	if err := client.Do(ctx, "/notes/show", map[string]any{
		"noteId": "nonexistent",
	}, &badResp); err != nil {
		if apiErr, ok := err.(*migo.APIError); ok {
			fmt.Printf("API error caught:\n")
			fmt.Printf("  Status:  %d\n", apiErr.StatusCode)
			fmt.Printf("  Code:    %s\n", apiErr.Code)
			fmt.Printf("  Message: %s\n", apiErr.Message)
			fmt.Printf("  ID:      %s\n", apiErr.ID)
		} else {
			fmt.Printf("Non-API error: %v\n", err)
		}
	}

	// -------------------------------------------------------------------
	// 5. Using services alongside raw Do
	// -------------------------------------------------------------------
	fmt.Println("\n=== 5. Mixing typed service with raw Do ===")

	// Use the typed service for the profile...
	as := auth.NewService(client)
	me, err := as.I(ctx)
	if err != nil {
		log.Fatalf("failed: %v", err)
	}
	fmt.Printf("Logged in as @%s\n", me.Username)

	// ...but fetch additional data via raw Do for fields not in the types.
	var extra struct {
		TwoFactorEnabled bool `json:"twoFactorEnabled"`
	}
	if err := client.Do(ctx, "/i", nil, &extra); err == nil {
		fmt.Printf("2FA enabled: %v\n", extra.TwoFactorEnabled)
	}

	// -------------------------------------------------------------------
	// 6. Server ping
	// -------------------------------------------------------------------
	fmt.Println("\n=== 6. Ping ===")
	var pong struct {
		Pong int64 `json:"pong"`
	}
	if err := client.Do(ctx, "/ping", nil, &pong); err != nil {
		log.Printf("ping failed: %v\n", err)
	} else {
		fmt.Printf("Pong: %dms\n", pong.Pong)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func ptr[T any](v T) *T { return &v }

func prettyJSON(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return
	}
	// Print a subset of interesting keys.
	var m map[string]any
	if err := json.Unmarshal(b, &m); err == nil {
		for _, key := range []string{"id", "username", "name", "host", "avatarUrl"} {
			if val, ok := m[key]; ok {
				fmt.Printf("  %s: %v\n", key, val)
			}
		}
	}
}
