// Command admin demonstrates administrative operations using the AdminService.
//
// Admin operations require elevated privileges on the Misskey instance.
// The token must belong to an admin or moderator account.
//
// Usage:
//
//	# Set via: export MIGO_TOKEN=\"<admin-token>\"
//	export MIGO_URL="https://your-instance.net/api"
//	# Set via: export MIGO_TOKEN="<admin-token>"
//	go run ./examples/admin
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sysnote8main/migo"
	"github.com/sysnote8main/migo/admin"
	"github.com/sysnote8main/migo/types"
)

func ptr[T any](v T) *T { return &v }

func main() {
	ctx := context.Background()

	baseURL := envOrDefault("MIGO_URL", "https://misskey.io/api")
	token := os.Getenv("MIGO_TOKEN")
	if token == "" {
		log.Fatal("MIGO_TOKEN is required (admin/moderator token)")
	}

	client := migo.NewClient(
		migo.WithBaseURL(baseURL),
		migo.WithToken(token),
	)

	as := admin.NewService(client)

	// -------------------------------------------------------------------
	// 1. List announcements
	// -------------------------------------------------------------------
	announcements, err := as.AnnouncementList(ctx)
	if err != nil {
		log.Printf("failed to list announcements: %v\n", err)
	} else {
		fmt.Printf("Announcements (%d):\n", len(announcements))
		for _, a := range announcements {
			fmt.Printf("  • [%s] %s\n", a.ID, a.Title)
		}
		fmt.Println()
	}

	// -------------------------------------------------------------------
	// 2. Create an announcement
	// -------------------------------------------------------------------
	icon := "info"
	display := "normal"
	silence := false

	newAnn, err := as.AnnouncementCreate(ctx, &admin.AnnouncementCreateRequest{
		Title:   "Maintenance Notice",
		Text:    "The server will undergo maintenance at 02:00 UTC.",
		Icon:    &icon,
		Display: &display,
		Silence: &silence,
	})
	if err != nil {
		log.Printf("failed to create announcement: %v\n", err)
	} else {
		fmt.Printf("Created announcement: %s — %s\n", newAnn.ID, newAnn.Title)

		// Clean up — delete the announcement we just created.
		if err := as.AnnouncementDelete(ctx, newAnn.ID); err != nil {
			log.Printf("failed to delete announcement: %v\n", err)
		} else {
			fmt.Println("Deleted announcement (cleanup)")
		}
		fmt.Println()
	}

	// -------------------------------------------------------------------
	// 3. List custom emojis
	// -------------------------------------------------------------------
	emojis, err := as.EmojiList(ctx, &admin.EmojiListRequest{Limit: ptr(10)})
	if err != nil {
		log.Printf("failed to list emojis: %v\n", err)
	} else {
		fmt.Printf("Custom emojis (first %d):\n", len(emojis))
		for _, e := range emojis {
			fmt.Printf("  • :%s: (category: %s)\n", e.Name, ptrStr(e.Category))
		}
		fmt.Println()
	}

	// -------------------------------------------------------------------
	// 4. List roles
	// -------------------------------------------------------------------
	roles, err := as.RoleList(ctx)
	if err != nil {
		log.Printf("failed to list roles: %v\n", err)
	} else {
		fmt.Printf("Roles (%d):\n", len(roles))
		for _, r := range roles {
			fmt.Printf("  • %s (%s)\n", r.Name, r.ID)
		}
		fmt.Println()
	}

	// -------------------------------------------------------------------
	// 5. Fetch server meta (admin view)
	// -------------------------------------------------------------------
	meta, err := as.Meta(ctx)
	if err != nil {
		log.Printf("failed to fetch meta: %v\n", err)
	} else {
		fmt.Printf("Server meta:\n")
		fmt.Printf("  Name:    %s\n", ptrStr(meta.Name))
		fmt.Printf("  Version: %s\n", meta.Version)
		fmt.Printf("  Desc:    %s\n", ptrStr(meta.Description))
		fmt.Println()
	}

	// -------------------------------------------------------------------
	// 6. List abuse reports
	// -------------------------------------------------------------------
	reports, err := as.AbuseReports(ctx)
	if err != nil {
		log.Printf("failed to list abuse reports: %v\n", err)
	} else {
		fmt.Printf("Abuse reports (%d):\n", len(reports))
		for _, r := range reports {
			status := "open"
			if r.Resolved {
				status = "resolved"
			}
			fmt.Printf("  • [%s] @%s — %s\n", status, r.TargetUser.Username, r.Comment)
		}
		fmt.Println()
	}

	// -------------------------------------------------------------------
	// 7. Use the generic Do method for a custom admin endpoint
	// -------------------------------------------------------------------
	fmt.Println("Server info via generic Do:")
	result := make(map[string]any)
	if err := as.Do(ctx, "/admin/meta", nil, &result); err != nil {
		log.Printf("generic meta call failed: %v\n", err)
	} else {
		for k, v := range result {
			if k == "version" || k == "name" || k == "description" {
				fmt.Printf("  %s: %v\n", k, v)
			}
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

// Ensure types is referenced.
var _ = types.VisibilityPublic
