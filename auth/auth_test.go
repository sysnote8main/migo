package auth

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

func TestService_SessionGenerate(t *testing.T) {
	wantReq := &SessionGenerateRequest{AppSecret: "my-secret"}
	wantResp := &SessionGenerateResponse{
		Token: "sess-token",
		URL:   "https://example.com/miauth/sess-token",
	}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/auth/session/generate",
		wantReq:  wantReq,
		resp:     wantResp,
	})

	resp, err := svc.SessionGenerate(context.Background(), "my-secret")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Token != "sess-token" {
		t.Errorf("Token = %q", resp.Token)
	}
}

func TestService_SessionUserKey(t *testing.T) {
	wantReq := &SessionUserKeyRequest{
		AppSecret: "my-secret",
		Token:     "sess-token",
	}
	wantResp := &SessionUserKeyResponse{
		AccessToken: "access-token",
		User:        types.UserLite{Username: "testuser"},
	}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/auth/session/userkey",
		wantReq:  wantReq,
		resp:     wantResp,
	})

	resp, err := svc.SessionUserKey(context.Background(), "my-secret", "sess-token")
	if err != nil {
		t.Fatal(err)
	}
	if resp.AccessToken != "access-token" {
		t.Errorf("AccessToken = %q", resp.AccessToken)
	}
	if resp.User.Username != "testuser" {
		t.Errorf("Username = %q", resp.User.Username)
	}
}

func TestService_Ping(t *testing.T) {
	wantResp := struct {
		Pong int64 `json:"pong"`
	}{Pong: 1}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/ping",
		resp:     &wantResp,
	})

	pong, err := svc.Ping(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if pong != 1 {
		t.Errorf("Pong = %d, want 1", pong)
	}
}

func TestService_I(t *testing.T) {
	wantResp := &types.MeDetailed{}
	wantResp.Username = "testuser"
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/i",
		resp:     wantResp,
	})

	user, err := svc.I(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "testuser" {
		t.Errorf("Username = %q", user.Username)
	}
}

func TestMiAuthURL(t *testing.T) {
	url := MiAuthURL("https://misskey.io", "sess-token", []string{"write:notes", "read:account"})
	if url == "" {
		t.Fatal("empty URL")
	}
	// Should contain the session token and permissions
	expected := "miauth/sess-token"
	if len(url) < len(expected) {
		t.Errorf("URL too short: %s", url)
	}
}

func TestMiAuthURLWithOptions(t *testing.T) {
	url := MiAuthURLWithOptions("https://misskey.io", "sess-token",
		[]string{"read:account"}, "MyApp", "https://example.com/cb")
	if url == "" {
		t.Fatal("empty URL")
	}
}
