package notes

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sysnote8main/migo/types"
)

// mockDoer implements the doer interface for testing.
type mockDoer struct {
	t        *testing.T
	wantPath string
	wantReq  any
	resp     any
	errResp  error
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
	if m.errResp != nil {
		return m.errResp
	}
	if m.resp != nil && resp != nil {
		data, _ := json.Marshal(m.resp)
		json.Unmarshal(data, resp)
	}
	return nil
}

func TestService_Create(t *testing.T) {
	wantReq := &CreateRequest{
		Text: strPtr("hello"),
		Poll: nil,
	}
	wantResp := &CreateResponse{
		CreatedNote: types.Note{ID: "note123", Text: strPtr("hello")},
	}

	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/create",
		wantReq:  wantReq,
		resp:     wantResp,
	})

	note, err := svc.Create(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if note.ID != "note123" {
		t.Errorf("note.ID = %q, want %q", note.ID, "note123")
	}
}

func TestService_Show(t *testing.T) {
	wantReq := &ShowRequest{NoteID: "note123"}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/show",
		wantReq:  wantReq,
		resp:     &types.Note{ID: "note123", Text: strPtr("test")},
	})

	note, err := svc.Show(context.Background(), "note123")
	if err != nil {
		t.Fatal(err)
	}
	if note.ID != "note123" {
		t.Errorf("note.ID = %q", note.ID)
	}
}

func TestService_Delete(t *testing.T) {
	wantReq := &DeleteRequest{NoteID: "note123"}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/delete",
		wantReq:  wantReq,
	})
	if err := svc.Delete(context.Background(), "note123"); err != nil {
		t.Fatal(err)
	}
}

func TestService_Search(t *testing.T) {
	wantReq := &SearchRequest{
		Query: "misskey",
		Limit: intPtr(10),
	}
	wantResp := []*types.Note{
		{ID: "n1", Text: strPtr("misskey is great")},
	}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/search",
		wantReq:  wantReq,
		resp:     &wantResp,
	})
	notes, err := svc.Search(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 1 || notes[0].ID != "n1" {
		t.Errorf("got %d notes, first ID = %v", len(notes), notes[0].ID)
	}
}

func TestService_CreateReaction(t *testing.T) {
	wantReq := &ReactionCreateRequest{NoteID: "n1", Reaction: "👍"}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/reactions/create",
		wantReq:  wantReq,
	})
	if err := svc.CreateReaction(context.Background(), "n1", "👍"); err != nil {
		t.Fatal(err)
	}
}

func TestService_DeleteReaction(t *testing.T) {
	wantReq := &ReactionDeleteRequest{NoteID: "n1"}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/notes/reactions/delete",
		wantReq:  wantReq,
	})
	if err := svc.DeleteReaction(context.Background(), "n1"); err != nil {
		t.Fatal(err)
	}
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
