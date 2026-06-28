package following

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

func TestService_Create(t *testing.T) {
	wantReq := &CreateRequest{UserID: "user1", WithReplies: boolPtr(true)}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/following/create",
		wantReq:  wantReq,
		resp:     &types.UserDetailed{},
	})
	user, err := svc.Create(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user is nil")
	}
}

func TestService_Delete(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/following/delete",
	})
	if err := svc.Delete(context.Background(), "user1"); err != nil {
		t.Fatal(err)
	}
}

func TestService_Invalidate(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/following/invalidate",
	})
	if err := svc.Invalidate(context.Background(), "user1"); err != nil {
		t.Fatal(err)
	}
}

func TestService_List(t *testing.T) {
	wantResp := []*Follower{{ID: "f1"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/following/list",
		resp:     &wantResp,
	})
	followers, err := svc.List(context.Background(), &ListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(followers) != 1 {
		t.Errorf("got %d followers", len(followers))
	}
}

func TestService_RequestsList(t *testing.T) {
	wantResp := []*Follower{{ID: "req1"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/following/requests/list",
		resp:     &wantResp,
	})
	reqs, err := svc.RequestsList(context.Background(), &RequestsListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(reqs) != 1 {
		t.Errorf("got %d requests", len(reqs))
	}
}

func TestService_RequestsAccept(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/following/requests/accept",
	})
	if err := svc.RequestsAccept(context.Background(), "user1"); err != nil {
		t.Fatal(err)
	}
}

func TestService_RequestsReject(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/following/requests/reject",
	})
	if err := svc.RequestsReject(context.Background(), "user1"); err != nil {
		t.Fatal(err)
	}
}

func TestService_RequestsCancel(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/following/requests/cancel",
	})
	if err := svc.RequestsCancel(context.Background(), "user1"); err != nil {
		t.Fatal(err)
	}
}

func boolPtr(b bool) *bool { return &b }
