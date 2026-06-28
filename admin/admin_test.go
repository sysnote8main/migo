package admin

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

func TestService_EmojiAdd(t *testing.T) {
	wantReq := &EmojiAddRequest{Name: "test", FileID: "file1"}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/emoji/add",
		wantReq:  wantReq,
		resp:     &types.EmojiDetailed{ID: "emoji1", Name: "test"},
	})
	emoji, err := svc.EmojiAdd(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if emoji.Name != "test" {
		t.Errorf("Name = %q", emoji.Name)
	}
}

func TestService_EmojiList(t *testing.T) {
	resp := []*types.EmojiDetailed{{ID: "e1"}, {ID: "e2"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/emoji/list",
		resp:     &resp,
	})
	emojis, err := svc.EmojiList(context.Background(), &EmojiListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(emojis) != 2 {
		t.Errorf("got %d emojis", len(emojis))
	}
}

func TestService_EmojiDelete(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/emoji/remove",
	})
	if err := svc.EmojiDelete(context.Background(), "e1"); err != nil {
		t.Fatal(err)
	}
}

func TestService_AnnouncementCreate(t *testing.T) {
	wantReq := &AnnouncementCreateRequest{Title: "Notice", Text: "Server maintenance"}
	resp := &types.Announcement{ID: "ann1", Title: "Notice"}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/announcements/create",
		wantReq:  wantReq,
		resp:     resp,
	})
	ann, err := svc.AnnouncementCreate(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if ann.Title != "Notice" {
		t.Errorf("Title = %q", ann.Title)
	}
}

func TestService_AnnouncementList(t *testing.T) {
	resp := []*types.Announcement{{ID: "ann1"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/announcements/list",
		resp:     &resp,
	})
	anns, err := svc.AnnouncementList(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(anns) != 1 {
		t.Errorf("got %d announcements", len(anns))
	}
}

func TestService_SuspendUser(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/suspend-user",
	})
	if err := svc.SuspendUser(context.Background(), "user1"); err != nil {
		t.Fatal(err)
	}
}

func TestService_UnsuspendUser(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/unsuspend-user",
	})
	if err := svc.UnsuspendUser(context.Background(), "user1"); err != nil {
		t.Fatal(err)
	}
}

func TestService_AbuseReports(t *testing.T) {
	resp := []*types.AbuseReport{{ID: "r1"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/abuse-user-reports",
		resp:     &resp,
	})
	reports, err := svc.AbuseReports(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 {
		t.Errorf("got %d reports", len(reports))
	}
}

func TestService_Do(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/some/custom/endpoint",
	})
	if err := svc.Do(context.Background(), "/admin/some/custom/endpoint", nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestService_RoleList(t *testing.T) {
	resp := []*types.UserRole{{ID: "role1", Name: "moderator"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/admin/roles/list",
		resp:     &resp,
	})
	roles, err := svc.RoleList(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 1 || roles[0].Name != "moderator" {
		t.Errorf("got roles: %+v", roles)
	}
}
