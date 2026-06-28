package drive

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

func TestService_Files(t *testing.T) {
	wantReq := &FilesRequest{Limit: intPtr(10)}
	wantResp := []*types.DriveFile{{ID: "f1"}, {ID: "f2"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/drive/files",
		wantReq:  wantReq,
		resp:     &wantResp,
	})
	files, err := svc.Files(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("got %d files", len(files))
	}
}

func TestService_FileInfo(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/drive/files/show",
		resp:     &types.DriveFile{ID: "f1", Name: "test.png"},
	})
	file, err := svc.FileInfo(context.Background(), "f1")
	if err != nil {
		t.Fatal(err)
	}
	if file.Name != "test.png" {
		t.Errorf("Name = %q", file.Name)
	}
}

func TestService_Info(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/drive",
		resp:     &types.Drive{Capacity: 1_000_000, Usage: 500_000},
	})
	drive, err := svc.Info(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if drive.Capacity != 1_000_000 {
		t.Errorf("Capacity = %d", drive.Capacity)
	}
}

func TestService_Folders(t *testing.T) {
	wantResp := []*types.DriveFolder{{ID: "fold1"}}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/drive/folders",
		resp:     &wantResp,
	})
	folders, err := svc.Folders(context.Background(), &FoldersRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(folders) != 1 {
		t.Errorf("got %d folders", len(folders))
	}
}

func TestService_CreateFolder(t *testing.T) {
	wantReq := &FolderCreateRequest{Name: "new-folder"}
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/drive/folders/create",
		wantReq:  wantReq,
		resp:     &types.DriveFolder{ID: "fold1", Name: "new-folder"},
	})
	folder, err := svc.CreateFolder(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}
	if folder.Name != "new-folder" {
		t.Errorf("Name = %q", folder.Name)
	}
}

func TestService_DeleteFolder(t *testing.T) {
	svc := NewService(&mockDoer{
		t:        t,
		wantPath: "/drive/folders/delete",
	})
	if err := svc.DeleteFolder(context.Background(), "fold1"); err != nil {
		t.Fatal(err)
	}
}

func intPtr(i int) *int { return &i }
