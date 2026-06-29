// Command basic demonstrates common migo operations: authentication check,
// note creation, timeline fetching, reactions, and search.
//
// Usage:
//
//	export MIGO_URL="https://your-instance.net/api"
//	# Set via: export MIGO_TOKEN="<your-token>"
//	go run ./examples/basic
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sysnote8main/migo"
	"github.com/sysnote8main/migo/auth"
	"github.com/sysnote8main/migo/notes"
	"github.com/sysnote8main/migo/timeline"
	"github.com/sysnote8main/migo/types"
)

func ptr[T any](v T) *T { return &v }

func main() {
	ctx := context.Background()

	// -----------------------------------------------------------------------
	// 1. Create the client
	// -----------------------------------------------------------------------
	baseURL := envOrDefault("MIGO_URL", "https://misskey.io/api")
	token := os.Getenv("MIGO_TOKEN")
	if token == "" {
		log.Fatal("MIGO_TOKEN is required")
	}

	client := migo.NewClient(
		migo.WithBaseURL(baseURL),
		migo.WithToken(token),
	)

	// -----------------------------------------------------------------------
	// 2. Verify authentication — fetch the authenticated user
	// -----------------------------------------------------------------------
	as := auth.NewService(client)
	me, err := as.I(ctx)
	if err != nil {
		log.Fatalf("failed to fetch profile: %v", err)
	}
	fmt.Printf("Logged in as @%s (%s)\n", me.Username, me.ID)
	fmt.Printf("  Notes: %d  |  Following: %d  |  Followers: %d\n",
		me.NotesCount, me.FollowingCount, me.FollowersCount)

	// -----------------------------------------------------------------------
	// 3. Create a note
	// -----------------------------------------------------------------------
	ns := notes.NewService(client)

	note, err := ns.Create(ctx, &notes.CreateRequest{
		Text: ptr("Hello, Misskey! This note was posted via migo."),
	})
	if err != nil {
		log.Fatalf("failed to create note: %v", err)
	}
	fmt.Printf("\nCreated note: %s\n", note.ID)

	// -----------------------------------------------------------------------
	// 4. Fetch the note we just created (notes/show)
	// -----------------------------------------------------------------------
	fetched, err := ns.Show(ctx, note.ID)
	if err != nil {
		log.Fatalf("failed to fetch note: %v", err)
	}
	fmt.Printf("Fetched note: %s\n", ellipsis(ptrStr(fetched.Text), 60))

	// -----------------------------------------------------------------------
	// 5. Add a reaction, then remove it
	// -----------------------------------------------------------------------
	if err := ns.CreateReaction(ctx, note.ID, "👍"); err != nil {
		log.Printf("failed to add reaction: %v", err)
	} else {
		fmt.Println("Added reaction: 👍")
	}

	if err := ns.DeleteReaction(ctx, note.ID); err != nil {
		log.Printf("failed to remove reaction: %v", err)
	} else {
		fmt.Println("Removed reaction")
	}

	// -----------------------------------------------------------------------
	// 6. Delete the note
	// -----------------------------------------------------------------------
	if err := ns.Delete(ctx, note.ID); err != nil {
		log.Printf("failed to delete note: %v", err)
	} else {
		fmt.Println("Deleted note")
	}

	// -----------------------------------------------------------------------
	// 7. Fetch home timeline (first 5 notes)
	// -----------------------------------------------------------------------
	ts := timeline.NewService(client)
	timelineNotes, err := ts.Home(ctx, &timeline.Request{
		Limit: ptr(5),
	})
	if err != nil {
		log.Fatalf("failed to fetch timeline: %v", err)
	}
	fmt.Printf("\nHome timeline (latest %d):\n", len(timelineNotes))
	for _, tlNote := range timelineNotes {
		text := ptrStr(tlNote.Text)
		if text == "" {
			text = "(no text — may be a renote)"
		}
		fmt.Printf("  • @%s: %s\n", tlNote.User.Username, ellipsis(text, 80))
	}

	// -----------------------------------------------------------------------
	// 8. Search notes containing "misskey"
	// -----------------------------------------------------------------------
	results, err := ns.Search(ctx, &notes.SearchRequest{
		Query: "misskey",
		Limit: ptr(3),
	})
	if err != nil {
		log.Printf("search failed: %v", err)
	} else {
		fmt.Printf("\nSearch results for \"misskey\" (%d):\n", len(results))
		for _, r := range results {
			fmt.Printf("  • @%s: %s\n", r.User.Username, ellipsis(ptrStr(r.Text), 80))
		}
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

// Ensure types is used (referenced in the structs)
var _ = types.VisibilityPublic
