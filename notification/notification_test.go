package notification

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

func TestService_List(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/i/notifications",
		resp:     &[]*types.Notification{{ID: "n1"}},
	})
	notifs, err := svc.List(context.Background(), &ListRequest{Limit: intPtr(10)})
	if err != nil {
		t.Fatal(err)
	}
	if len(notifs) != 1 {
		t.Errorf("got %d notifications", len(notifs))
	}
}

func TestService_ListGrouped(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/i/notifications-grouped",
		resp:     &[]any{},
	})
	_, err := svc.ListGrouped(context.Background(), &ListRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_MarkAllAsRead(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notifications/mark-all-as-read",
	})
	if err := svc.MarkAllAsRead(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestService_Create(t *testing.T) {
	req := &CreateRequest{Body: "test notification"}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notifications/create",
		wantReq:  req,
	})
	if err := svc.Create(context.Background(), req); err != nil {
		t.Fatal(err)
	}
}

func TestService_Flush(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notifications/flush",
	})
	if err := svc.Flush(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestService_TestNotification(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notifications/test-notification",
	})
	if err := svc.TestNotification(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func intPtr(i int) *int { return &i }
