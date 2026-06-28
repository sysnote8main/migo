# migo

**migo** is a Go client library for the [Misskey](https://misskey-hub.net/) API.

> ⚠️ **Early development stage.** The API is not yet stable and may change without notice.

## Features

- 🎯 **Type-safe** — All request/response types are defined as Go structs
- 🧩 **Modular** — Services organized by domain (notes, users, drive, auth, etc.)
- 🧪 **Mock-friendly** — `Client` interface means easy testing
- 🪶 **Zero dependencies** — Only Go standard library
- 🚀 **Context-aware** — All methods accept `context.Context`

## Installation

```bash
go get github.com/sysnote8main/migo
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sysnote8main/migo"
    "github.com/sysnote8main/migo/notes"
    "github.com/sysnote8main/migo/timeline"
)

func ptr[T any](v T) *T { return &v }

func main() {
    ctx := context.Background()

    client := migo.NewClient(
        migo.WithBaseURL("https://your-instance.net/api"),
        migo.WithToken("YOUR_ACCESS_TOKEN"),
    )

    // Create a note
    ns := notes.NewService(client)
    note, err := ns.Create(ctx, &notes.CreateRequest{
        Text: ptr("Hello, Misskey from migo!"),
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created note: %s\n", note.ID)

    // Fetch home timeline
    ts := timeline.NewService(client)
    timelineNotes, _ := ts.Home(ctx, &timeline.Request{Limit: ptr(10)})
    fmt.Printf("Home timeline: %d notes\n", len(timelineNotes))
}
```

### Using the Services Bundle

```go
svc := migo.NewServices(client)

// All services available from one struct
me, _ := svc.Auth.I(ctx)
fmt.Printf("Logged in as @%s\n", me.Username)

driveInfo, _ := svc.Drive.Info(ctx)
fmt.Printf("Drive usage: %d / %d bytes\n", driveInfo.Usage, driveInfo.Capacity)
```

## Package Structure

```
migo/
├── migo.go              # Client interface, APIError, NewClient
├── client.go            # HTTP transport implementation
├── services.go          # Services convenience bundle
│
├── types/               # Shared domain types (Note, User, DriveFile, etc.)
│
├── notes/               # Note management
├── users/               # User queries
├── drive/               # File & folder management
├── auth/                # MiAuth, session, ping, /i
├── timeline/            # Home/Local/Hybrid/Global/Channel timelines
├── following/           # Follow/unfollow, requests
├── notification/        # Notification management
├── admin/               # Administration (emoji, announcements, roles, etc.)
└── chat/                # Direct messages & chat rooms
```

## Authentication

migo supports two authentication methods:

### 1. Bearer Token (recommended for CLI/bots)

```go
client := migo.NewClient(
    migo.WithBaseURL("https://instance.net/api"),
    migo.WithToken("YOUR_TOKEN"),
)
```

### 2. MiAuth (for end-user apps)

```go
as := auth.NewService(client)

// Step 1: Generate a session
session, _ := as.SessionGenerate(ctx, "your-app-secret")

// Step 2: Direct user to authorize
miauthURL := auth.MiAuthURL("https://instance.net", session.Token,
    []string{"write:notes", "read:account"})
fmt.Printf("Open this URL: %s\n", miauthURL)

// Step 3: After user authorizes, exchange for access token
key, _ := as.SessionUserKey(ctx, "your-app-secret", session.Token)
fmt.Printf("Access token: %s\n", key.AccessToken)
```

## Error Handling

All API errors return `*migo.APIError`:

```go
note, err := ns.Create(ctx, &notes.CreateRequest{...})
if err != nil {
    if apiErr, ok := err.(*migo.APIError); ok {
        fmt.Printf("API error [%d] %s: %s\n",
            apiErr.StatusCode, apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("Network error: %v\n", err)
    }
}
```

## Custom / Unimplemented Endpoints

Use `client.Do()` for endpoints not yet wrapped:

```go
var result map[string]any
client.Do(ctx, "/i/notifications", map[string]any{
    "limit": 10,
}, &result)
```

## Testing

The `Client` interface makes mocking easy:

```go
type mockClient struct {
    doFunc func(ctx context.Context, path string, req, resp any) error
}

func (m *mockClient) Do(ctx context.Context, path string, req, resp any) error {
    return m.doFunc(ctx, path, req, resp)
}

// In test:
client := &mockClient{
    doFunc: func(ctx context.Context, path string, req, resp any) error {
        // Return mock response
        return json.Unmarshal([]byte(`{"id":"test123"}`), resp)
    },
}
ns := notes.NewService(client)
note, _ := ns.Create(ctx, &notes.CreateRequest{Text: ptr("test")})
```

## License

MIT
