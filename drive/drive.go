// Package drive provides the DriveService for file and folder management.
package drive

import (
	"context"

	"github.com/sysnote8main/migo/types"
)

// Service provides drive-related API operations.
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new DriveService.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// FilesRequest represents the request for drive/files.
type FilesRequest struct {
	Limit     *int      `json:"limit,omitempty"`
	SinceID   *types.ID `json:"sinceId,omitempty"`
	UntilID   *types.ID `json:"untilId,omitempty"`
	SinceDate *int64    `json:"sinceDate,omitempty"`
	UntilDate *int64    `json:"untilDate,omitempty"`
	FolderID  *types.ID `json:"folderId,omitempty"`
	Type      *string   `json:"type,omitempty"`
	Sort      *string   `json:"sort,omitempty"`
}

// FoldersRequest represents the request for drive/folders.
type FoldersRequest struct {
	Limit    *int      `json:"limit,omitempty"`
	SinceID  *types.ID `json:"sinceId,omitempty"`
	UntilID  *types.ID `json:"untilId,omitempty"`
	FolderID *types.ID `json:"folderId,omitempty"`
}

// FilesCreateRequest represents the metadata for drive/files/create.
// Note: for file upload (multipart/form-data), use UploadFile instead.
type FilesCreateRequest struct {
	FolderID    *types.ID `json:"folderId,omitempty"`
	Name        *string   `json:"name,omitempty"`
	IsSensitive *bool     `json:"isSensitive,omitempty"`
	Force       *bool     `json:"force,omitempty"`
}

// FolderCreateRequest represents the request for drive/folders/create.
type FolderCreateRequest struct {
	Name     string    `json:"name"`
	ParentID *types.ID `json:"parentId,omitempty"`
}

// FolderUpdateRequest represents the request for drive/folders/update.
type FolderUpdateRequest struct {
	FolderID types.ID  `json:"folderId"`
	Name     *string   `json:"name,omitempty"`
	ParentID *types.ID `json:"parentId,omitempty"`
}

// FolderDeleteRequest represents the request for drive/folders/delete.
type FolderDeleteRequest struct {
	FolderID types.ID `json:"folderId"`
}

// ---------------------------------------------------------------------------
// API methods
// ---------------------------------------------------------------------------

// Files lists drive files with optional filters.
func (s *Service) Files(ctx context.Context, req *FilesRequest) ([]*types.DriveFile, error) {
	var files []*types.DriveFile
	if err := s.cli.Do(ctx, "/drive/files", req, &files); err != nil {
		return nil, err
	}
	return files, nil
}

// FileInfo fetches a single drive file by its ID.
func (s *Service) FileInfo(ctx context.Context, fileID types.ID) (*types.DriveFile, error) {
	req := struct {
		FileID types.ID `json:"fileId"`
	}{FileID: fileID}
	var file types.DriveFile
	if err := s.cli.Do(ctx, "/drive/files/show", req, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

// Info returns the drive usage summary (capacity and usage).
func (s *Service) Info(ctx context.Context) (*types.Drive, error) {
	var drive types.Drive
	if err := s.cli.Do(ctx, "/drive", nil, &drive); err != nil {
		return nil, err
	}
	return &drive, nil
}

// Folders lists drive folders.
func (s *Service) Folders(ctx context.Context, req *FoldersRequest) ([]*types.DriveFolder, error) {
	var folders []*types.DriveFolder
	if err := s.cli.Do(ctx, "/drive/folders", req, &folders); err != nil {
		return nil, err
	}
	return folders, nil
}

// CreateFolder creates a new drive folder.
func (s *Service) CreateFolder(ctx context.Context, req *FolderCreateRequest) (*types.DriveFolder, error) {
	var folder types.DriveFolder
	if err := s.cli.Do(ctx, "/drive/folders/create", req, &folder); err != nil {
		return nil, err
	}
	return &folder, nil
}

// UpdateFolder updates a drive folder's name or parent.
func (s *Service) UpdateFolder(ctx context.Context, req *FolderUpdateRequest) (*types.DriveFolder, error) {
	var folder types.DriveFolder
	if err := s.cli.Do(ctx, "/drive/folders/update", req, &folder); err != nil {
		return nil, err
	}
	return &folder, nil
}

// DeleteFolder deletes a drive folder.
func (s *Service) DeleteFolder(ctx context.Context, folderID types.ID) error {
	req := &FolderDeleteRequest{FolderID: folderID}
	return s.cli.Do(ctx, "/drive/folders/delete", req, nil)
}
