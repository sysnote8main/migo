// Package admin provides the AdminService for administrative operations.
//
// With 99 admin endpoints, this package covers the most commonly used
// categories as typed methods. For less common admin endpoints, use
// the generic Do method.
package admin

import (
	"context"

	"github.com/sysnote8main/migo/types"
)

// Service provides administrative API operations (requires admin privileges).
type Service struct {
	cli doer
}

type doer interface {
	Do(ctx context.Context, path string, req, resp any) error
}

// NewService creates a new AdminService.
func NewService(cli doer) *Service {
	return &Service{cli: cli}
}

// Do sends a raw admin API request. Useful for admin endpoints
// that are not yet wrapped as typed methods.
func (s *Service) Do(ctx context.Context, path string, req, resp any) error {
	return s.cli.Do(ctx, path, req, resp)
}

// ---------------------------------------------------------------------------
// Emoji management
// ---------------------------------------------------------------------------

// EmojiAddRequest represents the request for admin/emoji/add.
type EmojiAddRequest struct {
	Name        string     `json:"name"`
	FileID      types.ID   `json:"fileId"`
	Category    *string    `json:"category,omitempty"`
	Aliases     []string   `json:"aliases,omitempty"`
	License     *string    `json:"license,omitempty"`
	IsSensitive *bool      `json:"isSensitive,omitempty"`
	LocalOnly   *bool      `json:"localOnly,omitempty"`
	RoleIDs     []types.ID `json:"roleIdsThatCanBeUsedThisEmojiAsReaction,omitempty"`
}

// EmojiUpdateRequest represents the request for admin/emoji/update.
type EmojiUpdateRequest struct {
	EmojiID     types.ID   `json:"id"`
	FileID      *types.ID  `json:"fileId,omitempty"`
	Category    *string    `json:"category,omitempty"`
	Aliases     []string   `json:"aliases,omitempty"`
	License     *string    `json:"license,omitempty"`
	IsSensitive *bool      `json:"isSensitive,omitempty"`
	LocalOnly   *bool      `json:"localOnly,omitempty"`
	RoleIDs     []types.ID `json:"roleIdsThatCanBeUsedThisEmojiAsReaction,omitempty"`
	Name        *string    `json:"name,omitempty"`
}

// EmojiListRequest represents the request for admin/emoji/list.
type EmojiListRequest struct {
	Query     *string   `json:"query,omitempty"`
	Limit     *int      `json:"limit,omitempty"`
	SinceID   *types.ID `json:"sinceId,omitempty"`
	UntilID   *types.ID `json:"untilId,omitempty"`
	SinceDate *int64    `json:"sinceDate,omitempty"`
	UntilDate *int64    `json:"untilDate,omitempty"`
}

// EmojiAdd creates a new custom emoji.
func (s *Service) EmojiAdd(ctx context.Context, req *EmojiAddRequest) (*types.EmojiDetailed, error) {
	var emoji types.EmojiDetailed
	if err := s.cli.Do(ctx, "/admin/emoji/add", req, &emoji); err != nil {
		return nil, err
	}
	return &emoji, nil
}

// EmojiList lists custom emojis.
func (s *Service) EmojiList(ctx context.Context, req *EmojiListRequest) ([]*types.EmojiDetailed, error) {
	var emojis []*types.EmojiDetailed
	if err := s.cli.Do(ctx, "/admin/emoji/list", req, &emojis); err != nil {
		return nil, err
	}
	return emojis, nil
}

// EmojiDelete deletes a custom emoji.
func (s *Service) EmojiDelete(ctx context.Context, emojiID types.ID) error {
	req := struct {
		ID types.ID `json:"id"`
	}{ID: emojiID}
	return s.cli.Do(ctx, "/admin/emoji/remove", req, nil)
}

// EmojiUpdate updates a custom emoji.
func (s *Service) EmojiUpdate(ctx context.Context, req *EmojiUpdateRequest) (*types.EmojiDetailed, error) {
	var emoji types.EmojiDetailed
	if err := s.cli.Do(ctx, "/admin/emoji/update", req, &emoji); err != nil {
		return nil, err
	}
	return &emoji, nil
}

// EmojiSetAliasesBulk sets aliases for multiple emojis at once.
func (s *Service) EmojiSetAliasesBulk(ctx context.Context, ids []types.ID, aliases []string) error {
	req := struct {
		IDs     []types.ID `json:"ids"`
		Aliases []string   `json:"aliases"`
	}{IDs: ids, Aliases: aliases}
	return s.cli.Do(ctx, "/admin/emoji/set-aliases-bulk", req, nil)
}

// EmojiCopy copies an emoji from a remote instance.
func (s *Service) EmojiCopy(ctx context.Context, emojiID types.ID) (*types.EmojiDetailed, error) {
	req := struct {
		EmojiID types.ID `json:"emojiId"`
	}{EmojiID: emojiID}
	var emoji types.EmojiDetailed
	if err := s.cli.Do(ctx, "/admin/emoji/copy", req, &emoji); err != nil {
		return nil, err
	}
	return &emoji, nil
}

// ---------------------------------------------------------------------------
// Announcements
// ---------------------------------------------------------------------------

// AnnouncementCreateRequest represents the request for admin/announcements/create.
type AnnouncementCreateRequest struct {
	Title    string  `json:"title"`
	Text     string  `json:"text"`
	ImageURL *string `json:"imageUrl,omitempty"`
	Icon     *string `json:"icon,omitempty"`
	Display  *string `json:"display,omitempty"`
	ForYou   *bool   `json:"forYou,omitempty"`
	Silence  *bool   `json:"silence,omitempty"`
}

// AnnouncementUpdateRequest represents the request for admin/announcements/update.
type AnnouncementUpdateRequest struct {
	ID       types.ID `json:"id"`
	Title    *string  `json:"title,omitempty"`
	Text     *string  `json:"text,omitempty"`
	ImageURL *string  `json:"imageUrl,omitempty"`
	Icon     *string  `json:"icon,omitempty"`
	Display  *string  `json:"display,omitempty"`
	ForYou   *bool    `json:"forYou,omitempty"`
	Silence  *bool    `json:"silence,omitempty"`
}

// AnnouncementCreate creates a new announcement.
func (s *Service) AnnouncementCreate(ctx context.Context, req *AnnouncementCreateRequest) (*types.Announcement, error) {
	var ann types.Announcement
	if err := s.cli.Do(ctx, "/admin/announcements/create", req, &ann); err != nil {
		return nil, err
	}
	return &ann, nil
}

// AnnouncementList lists all announcements.
func (s *Service) AnnouncementList(ctx context.Context) ([]*types.Announcement, error) {
	var anns []*types.Announcement
	if err := s.cli.Do(ctx, "/admin/announcements/list", nil, &anns); err != nil {
		return nil, err
	}
	return anns, nil
}

// AnnouncementUpdate updates an announcement.
func (s *Service) AnnouncementUpdate(ctx context.Context, req *AnnouncementUpdateRequest) (*types.Announcement, error) {
	var ann types.Announcement
	if err := s.cli.Do(ctx, "/admin/announcements/update", req, &ann); err != nil {
		return nil, err
	}
	return &ann, nil
}

// AnnouncementDelete deletes an announcement.
func (s *Service) AnnouncementDelete(ctx context.Context, id types.ID) error {
	req := struct {
		ID types.ID `json:"id"`
	}{ID: id}
	return s.cli.Do(ctx, "/admin/announcements/delete", req, nil)
}

// ---------------------------------------------------------------------------
// Accounts (user management)
// ---------------------------------------------------------------------------

// AccountCreateRequest represents the request for admin/accounts/create.
type AccountCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Bio      string `json:"bio,omitempty"`
}

// AccountCreate creates a new user account (no auth required).
func (s *Service) AccountCreate(ctx context.Context, req *AccountCreateRequest) (*types.UserDetailed, error) {
	var user types.UserDetailed
	if err := s.cli.Do(ctx, "/admin/accounts/create", req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// AccountDelete deletes a user account.
func (s *Service) AccountDelete(ctx context.Context, userID types.ID) error {
	req := struct {
		UserID types.ID `json:"userId"`
	}{UserID: userID}
	return s.cli.Do(ctx, "/admin/accounts/delete", req, nil)
}

// AccountFindByEmail finds a user by email address.
func (s *Service) AccountFindByEmail(ctx context.Context, email string) (*types.UserDetailed, error) {
	req := struct {
		Email string `json:"email"`
	}{Email: email}
	var user types.UserDetailed
	if err := s.cli.Do(ctx, "/admin/accounts/find-by-email", req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// ---------------------------------------------------------------------------
// Advertisements
// ---------------------------------------------------------------------------

// AdCreateRequest represents the request for admin/ad/create.
type AdCreateRequest struct {
	Memo      *string `json:"memo,omitempty"`
	Place     string  `json:"place"`    // "square" or "horizontal"
	Priority  string  `json:"priority"` // "high" or "middle" or "low"
	Ratio     *int    `json:"ratio,omitempty"`
	URL       string  `json:"url"`
	ImageURL  string  `json:"imageUrl"`
	ExpiresAt int64   `json:"expiresAt"`
	StartsAt  int64   `json:"startsAt"`
	DayOfWeek *int    `json:"dayOfWeek,omitempty"`
}

type AdListRequest struct {
	Limit   *int      `json:"limit,omitempty"`
	SinceID *types.ID `json:"sinceId,omitempty"`
	UntilID *types.ID `json:"untilId,omitempty"`
}

// AdCreate creates a new advertisement.
func (s *Service) AdCreate(ctx context.Context, req *AdCreateRequest) error {
	return s.cli.Do(ctx, "/admin/ad/create", req, nil)
}

// AdList lists advertisements.
func (s *Service) AdList(ctx context.Context, req *AdListRequest) ([]*Ad, error) {
	var ads []*Ad
	if err := s.cli.Do(ctx, "/admin/ad/list", req, &ads); err != nil {
		return nil, err
	}
	return ads, nil
}

// AdUpdate updates an advertisement.
func (s *Service) AdUpdate(ctx context.Context, req *AdUpdateRequest) error {
	return s.cli.Do(ctx, "/admin/ad/update", req, nil)
}

// AdDelete deletes an advertisement.
func (s *Service) AdDelete(ctx context.Context, adID types.ID) error {
	req := struct {
		ID types.ID `json:"id"`
	}{ID: adID}
	return s.cli.Do(ctx, "/admin/ad/delete", req, nil)
}

// AdUpdateRequest represents the request for admin/ad/update.
type AdUpdateRequest struct {
	ID        types.ID `json:"id"`
	Memo      *string  `json:"memo,omitempty"`
	Place     *string  `json:"place,omitempty"`
	Priority  *string  `json:"priority,omitempty"`
	Ratio     *int     `json:"ratio,omitempty"`
	URL       *string  `json:"url,omitempty"`
	ImageURL  *string  `json:"imageUrl,omitempty"`
	ExpiresAt *int64   `json:"expiresAt,omitempty"`
	StartsAt  *int64   `json:"startsAt,omitempty"`
	DayOfWeek *int     `json:"dayOfWeek,omitempty"`
}

// Ad type (response).
type Ad struct {
	ID        types.ID `json:"id"`
	ExpiresAt int64    `json:"expiresAt"`
	StartsAt  int64    `json:"startsAt"`
	Place     string   `json:"place"`
	Priority  string   `json:"priority"`
	Ratio     int      `json:"ratio"`
	URL       string   `json:"url"`
	ImageURL  string   `json:"imageUrl"`
	Memo      *string  `json:"memo,omitempty"`
	DayOfWeek *int     `json:"dayOfWeek,omitempty"`
}

// Ensure Ad is accessible from package
var _ = Ad{}

// ---------------------------------------------------------------------------
// Invites
// ---------------------------------------------------------------------------

// InviteCreateRequest represents the request for admin/invite/create.
type InviteCreateRequest struct {
	Count     *int   `json:"count,omitempty"`
	ExpiresAt *int64 `json:"expiresAt,omitempty"`
}

// InviteListRequest represents the request for admin/invite/list.
type InviteListRequest struct {
	Limit   *int      `json:"limit,omitempty"`
	SinceID *types.ID `json:"sinceId,omitempty"`
	UntilID *types.ID `json:"untilId,omitempty"`
	Type    *string   `json:"type,omitempty"` // "unused", "used", "expired"
}

// InviteCreate creates invitation codes.
func (s *Service) InviteCreate(ctx context.Context, req *InviteCreateRequest) error {
	return s.cli.Do(ctx, "/admin/invite/create", req, nil)
}

// InviteList lists invitation codes.
func (s *Service) InviteList(ctx context.Context, req *InviteListRequest) ([]any, error) {
	var invites []any
	if err := s.cli.Do(ctx, "/admin/invite/list", req, &invites); err != nil {
		return nil, err
	}
	return invites, nil
}

// ---------------------------------------------------------------------------
// Roles
// ---------------------------------------------------------------------------

// RoleCreateRequest represents the request for admin/roles/create.
type RoleCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Color       *string `json:"color,omitempty"`
	IconURL     *string `json:"iconUrl,omitempty"`
	Target      string  `json:"target"` // "manual" or "conditional"
	// ... many policy fields omitted for brevity
}

// RoleList lists all roles.
func (s *Service) RoleList(ctx context.Context) ([]*types.UserRole, error) {
	var roles []*types.UserRole
	if err := s.cli.Do(ctx, "/admin/roles/list", nil, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

// RoleCreate creates a new role.
func (s *Service) RoleCreate(ctx context.Context, req *RoleCreateRequest) error {
	return s.cli.Do(ctx, "/admin/roles/create", req, nil)
}

// RoleDelete deletes a role.
func (s *Service) RoleDelete(ctx context.Context, roleID types.ID) error {
	req := struct {
		RoleID types.ID `json:"roleId"`
	}{RoleID: roleID}
	return s.cli.Do(ctx, "/admin/roles/delete", req, nil)
}

// RoleAssign assigns a role to a user.
func (s *Service) RoleAssign(ctx context.Context, roleID, userID types.ID) error {
	req := struct {
		RoleID types.ID `json:"roleId"`
		UserID types.ID `json:"userId"`
	}{RoleID: roleID, UserID: userID}
	return s.cli.Do(ctx, "/admin/roles/assign", req, nil)
}

// RoleUnassign unassigns a role from a user.
func (s *Service) RoleUnassign(ctx context.Context, roleID, userID types.ID) error {
	req := struct {
		RoleID types.ID `json:"roleId"`
		UserID types.ID `json:"userId"`
	}{RoleID: roleID, UserID: userID}
	return s.cli.Do(ctx, "/admin/roles/unassign", req, nil)
}

// ---------------------------------------------------------------------------
// Abuse reports
// ---------------------------------------------------------------------------

// AbuseReports lists abuse reports.
func (s *Service) AbuseReports(ctx context.Context) ([]*types.AbuseReport, error) {
	var reports []*types.AbuseReport
	if err := s.cli.Do(ctx, "/admin/abuse-user-reports", nil, &reports); err != nil {
		return nil, err
	}
	return reports, nil
}

// AbuseReportResolve resolves an abuse report.
func (s *Service) AbuseReportResolve(ctx context.Context, reportID types.ID) error {
	req := struct {
		ReportID types.ID `json:"reportId"`
	}{ReportID: reportID}
	return s.cli.Do(ctx, "/admin/resolve-abuse-user-report", req, nil)
}

// ---------------------------------------------------------------------------
// Moderation
// ---------------------------------------------------------------------------

// SuspendUser suspends a user account.
func (s *Service) SuspendUser(ctx context.Context, userID types.ID) error {
	req := struct {
		UserID types.ID `json:"userId"`
	}{UserID: userID}
	return s.cli.Do(ctx, "/admin/suspend-user", req, nil)
}

// UnsuspendUser unsuspends a user account.
func (s *Service) UnsuspendUser(ctx context.Context, userID types.ID) error {
	req := struct {
		UserID types.ID `json:"userId"`
	}{UserID: userID}
	return s.cli.Do(ctx, "/admin/unsuspend-user", req, nil)
}

// ShowUser shows full user details (admin view).
func (s *Service) ShowUser(ctx context.Context, userID types.ID) (*types.UserDetailed, error) {
	req := struct {
		UserID types.ID `json:"userId"`
	}{UserID: userID}
	var user types.UserDetailed
	if err := s.cli.Do(ctx, "/admin/show-user", req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserNote updates the moderation note on a user.
func (s *Service) UpdateUserNote(ctx context.Context, userID types.ID, note string) error {
	req := struct {
		UserID types.ID `json:"userId"`
		Note   string   `json:"note"`
	}{UserID: userID, Note: note}
	return s.cli.Do(ctx, "/admin/update-user-note", req, nil)
}

// ---------------------------------------------------------------------------
// Server meta
// ---------------------------------------------------------------------------

// Meta returns the server meta information (admin view).
func (s *Service) Meta(ctx context.Context) (*types.Meta, error) {
	var meta types.Meta
	if err := s.cli.Do(ctx, "/admin/meta", nil, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// UpdateMeta updates server meta settings.
func (s *Service) UpdateMeta(ctx context.Context, req map[string]any) error {
	return s.cli.Do(ctx, "/admin/update-meta", req, nil)
}

// ---------------------------------------------------------------------------
// Drive management
// ---------------------------------------------------------------------------

// DriveFiles lists all drive files (admin view).
func (s *Service) DriveFiles(ctx context.Context, req map[string]any) ([]*types.DriveFile, error) {
	var files []*types.DriveFile
	if err := s.cli.Do(ctx, "/admin/drive/files", req, &files); err != nil {
		return nil, err
	}
	return files, nil
}

// DriveCleanup cleans up unused drive files.
func (s *Service) DriveCleanup(ctx context.Context) error {
	return s.cli.Do(ctx, "/admin/drive/cleanup", nil, nil)
}

// DriveCleanRemoteFiles cleans up remote drive files.
func (s *Service) DriveCleanRemoteFiles(ctx context.Context) error {
	return s.cli.Do(ctx, "/admin/drive/clean-remote-files", nil, nil)
}

// DriveShowFile shows a drive file details (admin view).
func (s *Service) DriveShowFile(ctx context.Context, fileID types.ID) (*types.DriveFile, error) {
	req := struct {
		FileID types.ID `json:"fileId"`
	}{FileID: fileID}
	var file types.DriveFile
	if err := s.cli.Do(ctx, "/admin/drive/show-file", req, &file); err != nil {
		return nil, err
	}
	return &file, nil
}
