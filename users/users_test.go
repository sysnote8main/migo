package users

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

func TestService_Show(t *testing.T) {
	wantReq := &ShowRequest{UserID: idPtr("user1")}
	wantResp := &types.UserDetailed{}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/users/show",
		wantReq:  wantReq,
		resp:     wantResp,
	})

	user, err := svc.Show(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user is nil")
	}
}

func TestService_ShowBulk(t *testing.T) {
	wantReq := &ShowRequest{
		UserIDs: []types.ID{"user1", "user2"},
	}
	wantResp := []*types.UserDetailed{{}, {}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/users/show",
		wantReq:  wantReq,
		resp:     &wantResp,
	})

	users, err := svc.ShowBulk(context.Background(), []types.ID{"user1", "user2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Errorf("got %d users, want 2", len(users))
	}
}

func TestService_Search(t *testing.T) {
	wantReq := &SearchRequest{
		Query:  "alice",
		Limit:  intPtr(10),
		Origin: strPtr("local"),
	}
	wantResp := []*types.UserDetailed{{}, {}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/users/search",
		wantReq:  wantReq,
		resp:     &wantResp,
	})

	users, err := svc.Search(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Errorf("got %d users, want 2", len(users))
	}
}

func idPtr(id string) *types.ID {
	v := types.ID(id)
	return &v
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
