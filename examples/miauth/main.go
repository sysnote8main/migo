// Command miauth demonstrates the MiAuth authentication flow.
//
// MiAuth lets end users authorize third-party apps without sharing
// their password. The flow is:
//
//  1. Generate a session (obtain session token + MiAuth URL)
//  2. Direct the user to open the MiAuth URL in a browser
//  3. After the user approves, exchange the session token for an access token
//  4. Use the access token for regular API calls
//
// Usage:
//
//	export MIGO_URL="https://your-instance.net/api"
//	# Set via: export MIGO_APP_SECRET="<your-app-secret>"
//	go run ./examples/miauth
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sysnote8main/migo"
	"github.com/sysnote8main/migo/auth"
)

func main() {
	ctx := context.Background()

	baseURL := os.Getenv("MIGO_URL")
	if baseURL == "" {
		baseURL = "https://misskey.io/api"
	}

	appSecret := os.Getenv("MIGO_APP_SECRET")
	if appSecret == "" {
		log.Fatal("MIGO_APP_SECRET is required")
	}

	// -------------------------------------------------------------------
	// Step 1: Create an unauthenticated client and generate a session
	// -------------------------------------------------------------------
	// MiAuth does not require a token — authentication happens via
	// the session flow.
	client := migo.NewClient(migo.WithBaseURL(baseURL))
	as := auth.NewService(client)

	session, err := as.SessionGenerate(ctx, appSecret)
	if err != nil {
		log.Fatalf("failed to generate session: %v", err)
	}

	fmt.Printf("Session token: %s\n", session.Token)
	fmt.Printf("MiAuth URL:    %s\n\n", session.URL)

	// -------------------------------------------------------------------
	// Step 2: Build an authorization URL with requested permissions
	// -------------------------------------------------------------------
	permissions := []string{
		"write:notes",
		"read:account",
		"read:drive",
		"write:reactions",
	}

	authURL := auth.MiAuthURL(instanceBase(baseURL), session.Token, permissions)
	fmt.Println("Ask the user to open this URL in their browser:")
	fmt.Println("  " + authURL)

	// You can also use the URL returned by SessionGenerate directly:
	fmt.Println("\nOr use the session URL directly:")
	fmt.Println("  " + session.URL)

	// With a callback URL and app name:
	authURL2 := auth.MiAuthURLWithOptions(
		instanceBase(baseURL),
		session.Token,
		permissions,
		"My migo App",                  // app name shown on the consent page
		"https://example.com/callback", // callback URL after approval
	)
	fmt.Println("\nWith custom name and callback:")
	fmt.Println("  " + authURL2)

	fmt.Println("\n=== IMPORTANT ===")
	fmt.Println("The user must approve the request in their browser.")
	fmt.Println("After approval, call SessionUserKey to get the access token.")
	fmt.Println("(This example continues automatically — in a real app you")
	fmt.Println("would wait for the user to approve before proceeding.)")

	// -------------------------------------------------------------------
	// Step 3: Exchange the session token for an access token
	// -------------------------------------------------------------------
	// NOTE: This will fail with a "permission denied" error unless the
	// user has actually approved the session in their browser.
	key, err := as.SessionUserKey(ctx, appSecret, session.Token)
	if err != nil {
		log.Printf("SessionUserKey failed (expected unless approved): %v\n", err)
		fmt.Println("\nOnce approved, run this snippet to obtain the token:")
		fmt.Printf("  key, err := as.SessionUserKey(ctx, %q, %q)\n", appSecret, session.Token)
		fmt.Printf("  fmt.Println(\"Access token:\", key.AccessToken)\n")
		return
	}

	fmt.Printf("\nAccess token: %s\n", key.AccessToken)
	fmt.Printf("Authorized user: @%s (%s)\n", key.User.Username, key.User.ID)

	// -------------------------------------------------------------------
	// Step 4: Use the access token for authenticated API calls
	// -------------------------------------------------------------------
	authedClient := migo.NewClient(
		migo.WithBaseURL(baseURL),
		migo.WithToken(key.AccessToken),
	)
	authedAuth := auth.NewService(authedClient)

	me, err := authedAuth.I(ctx)
	if err != nil {
		log.Fatalf("failed to fetch profile with new token: %v", err)
	}
	fmt.Printf("\nAuthenticated as @%s — token works!\n", me.Username)
}

// instanceBase strips "/api" suffix from the API URL to get the instance base URL.
func instanceBase(apiURL string) string {
	if len(apiURL) >= 4 && apiURL[len(apiURL)-4:] == "/api" {
		return apiURL[:len(apiURL)-4]
	}
	return apiURL
}
