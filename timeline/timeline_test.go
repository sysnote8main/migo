package timeline

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

func TestService_Home(t *testing.T) {
	wantResp := []*types.Note{{ID: "n1"}, {ID: "n2"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/timeline",
		resp:     &wantResp,
	})
	notes, err := svc.Home(context.Background(), &Request{Limit: intPtr(10)})
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 2 {
		t.Errorf("got %d notes", len(notes))
	}
}

func TestService_Local(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/local-timeline",
		resp:     &[]*types.Note{},
	})
	notes, err := svc.Local(context.Background(), &Request{})
	if err != nil {
		t.Fatal(err)
	}
	if notes == nil {
		t.Fatal("notes is nil")
	}
}

func TestService_Hybrid(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/hybrid-timeline",
		resp:     &[]*types.Note{},
	})
	_, err := svc.Hybrid(context.Background(), &Request{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_Global(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/global-timeline",
		resp:     &[]*types.Note{},
	})
	_, err := svc.Global(context.Background(), &Request{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_Channel(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/channels/timeline",
		resp:     &[]*types.Note{},
	})
	_, err := svc.Channel(context.Background(), &ChannelRequest{ChannelID: "ch1"})
	if err != nil {
		t.Fatal(err)
	}
}

func intPtr(i int) *int { return &i }
