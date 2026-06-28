// Package types contains shared data types used across the migo packages.
//
// These types correspond to Misskey API response objects and are independent
// of any specific client implementation, making them importable by any
// sub-package without creating circular dependencies.
package types

import "time"

// ID is a Misskey object ID.
type ID string

// NoteVisibility represents the visibility of a note.
type NoteVisibility string

const (
	VisibilityPublic    NoteVisibility = "public"
	VisibilityHome      NoteVisibility = "home"
	VisibilityFollowers NoteVisibility = "followers"
	VisibilitySpecified NoteVisibility = "specified"
)

// ---------------------------------------------------------------------------
// User types
// ---------------------------------------------------------------------------

// UserLite is the minimal user representation returned by most endpoints.
type UserLite struct {
	ID                           ID                 `json:"id"`
	Name                         *string            `json:"name"`
	Username                     string             `json:"username"`
	Host                         *string            `json:"host"`
	AvatarURL                    string             `json:"avatarUrl"`
	AvatarBlurhash               *string            `json:"avatarBlurhash"`
	AvatarDecorations            []AvatarDecoration `json:"avatarDecorations"`
	IsBot                        bool               `json:"isBot,omitempty"`
	IsCat                        bool               `json:"isCat,omitempty"`
	RequireSigninToViewContents  bool               `json:"requireSigninToViewContents,omitempty"`
	MakeNotesFollowersOnlyBefore *int64             `json:"makeNotesFollowersOnlyBefore,omitempty"`
	MakeNotesHiddenBefore        *int64             `json:"makeNotesHiddenBefore,omitempty"`
	Instance                     *Instance          `json:"instance,omitempty"`
	Emojis                       map[string]string  `json:"emojis"`
	OnlineStatus                 string             `json:"onlineStatus"`
	BadgeRoles                   []BadgeRole        `json:"badgeRoles"`
}

// AvatarDecoration represents a user's avatar decoration.
type AvatarDecoration struct {
	ID      ID     `json:"id"`
	Angle   *int   `json:"angle,omitempty"`
	FlipH   bool   `json:"flipH,omitempty"`
	URL     string `json:"url"`
	OffsetX *int   `json:"offsetX,omitempty"`
	OffsetY *int   `json:"offsetY,omitempty"`
}

// Instance represents a remote instance info.
type Instance struct {
	Name            *string `json:"name,omitempty"`
	SoftwareName    *string `json:"softwareName,omitempty"`
	SoftwareVersion *string `json:"softwareVersion,omitempty"`
	IconURL         *string `json:"iconUrl,omitempty"`
	FaviconURL      *string `json:"faviconUrl,omitempty"`
	ThemeColor      *string `json:"themeColor,omitempty"`
}

// BadgeRole represents a badge role on user profile.
type BadgeRole struct {
	Name         string  `json:"name"`
	IconURL      *string `json:"iconUrl,omitempty"`
	DisplayOrder int     `json:"displayOrder"`
}

// UserField is a profile field entry.
type UserField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// UserDetailedNotMeOnly contains detailed user fields shown when viewing
// another user (not the authenticated user).
type UserDetailedNotMeOnly struct {
	URL                            *string     `json:"url"`
	URI                            *string     `json:"uri"`
	MovedTo                        *string     `json:"movedTo"`
	AlsoKnownAs                    []string    `json:"alsoKnownAs"`
	CreatedAt                      time.Time   `json:"createdAt"`
	UpdatedAt                      *time.Time  `json:"updatedAt"`
	LastFetchedAt                  *time.Time  `json:"lastFetchedAt"`
	BannerURL                      *string     `json:"bannerUrl"`
	BannerBlurhash                 *string     `json:"bannerBlurhash"`
	IsLocked                       bool        `json:"isLocked"`
	IsSilenced                     bool        `json:"isSilenced"`
	IsSuspended                    bool        `json:"isSuspended"`
	Description                    *string     `json:"description"`
	Location                       *string     `json:"location"`
	Birthday                       *string     `json:"birthday"`
	Lang                           *string     `json:"lang"`
	Fields                         []UserField `json:"fields"`
	VerifiedLinks                  []string    `json:"verifiedLinks"`
	FollowersCount                 int         `json:"followersCount"`
	FollowingCount                 int         `json:"followingCount"`
	NotesCount                     int         `json:"notesCount"`
	PinnedNoteIds                  []ID        `json:"pinnedNoteIds"`
	PinnedNotes                    []Note      `json:"pinnedNotes"`
	PinnedPageID                   *string     `json:"pinnedPageId"`
	PinnedPage                     *Page       `json:"pinnedPage"`
	PublicReactions                bool        `json:"publicReactions"`
	FollowingVisibility            string      `json:"followingVisibility"`
	FollowersVisibility            string      `json:"followersVisibility"`
	ChatScope                      string      `json:"chatScope"`
	CanChat                        bool        `json:"canChat"`
	Roles                          []UserRole  `json:"roles"`
	Memo                           *string     `json:"memo"`
	ModerationNote                 string      `json:"moderationNote"`
	IsFollowing                    bool        `json:"isFollowing"`
	IsFollowed                     bool        `json:"isFollowed"`
	HasPendingFollowRequestFromYou bool        `json:"hasPendingFollowRequestFromYou"`
	HasPendingFollowRequestToYou   bool        `json:"hasPendingFollowRequestToYou"`
	IsBlocking                     bool        `json:"isBlocking"`
	IsBlocked                      bool        `json:"isBlocked"`
	IsMuted                        bool        `json:"isMuted"`
	IsRenoteMuted                  bool        `json:"isRenoteMuted"`
	Notify                         string      `json:"notify"`
	WithReplies                    bool        `json:"withReplies"`
}

// MeDetailedOnly contains fields only available for the authenticated user.
type MeDetailedOnly struct {
	AvatarID                        *string        `json:"avatarId"`
	BannerID                        *string        `json:"bannerId"`
	FollowedMessage                 *string        `json:"followedMessage"`
	IsModerator                     bool           `json:"isModerator"`
	IsAdmin                         bool           `json:"isAdmin"`
	InjectFeaturedNote              bool           `json:"injectFeaturedNote"`
	ReceiveAnnouncementEmail        bool           `json:"receiveAnnouncementEmail"`
	AlwaysMarkNsfw                  bool           `json:"alwaysMarkNsfw"`
	AutoSensitive                   bool           `json:"autoSensitive"`
	CarefulBot                      bool           `json:"carefulBot"`
	AutoAcceptFollowed              bool           `json:"autoAcceptFollowed"`
	NoCrawle                        bool           `json:"noCrawle"`
	PreventAiLearning               bool           `json:"preventAiLearning"`
	IsExplorable                    bool           `json:"isExplorable"`
	IsDeleted                       bool           `json:"isDeleted"`
	TwoFactorBackupCodesStock       string         `json:"twoFactorBackupCodesStock"`
	HideOnlineStatus                bool           `json:"hideOnlineStatus"`
	HasUnreadSpecifiedNotes         bool           `json:"hasUnreadSpecifiedNotes"`
	HasUnreadMentions               bool           `json:"hasUnreadMentions"`
	HasUnreadAnnouncement           bool           `json:"hasUnreadAnnouncement"`
	UnreadAnnouncements             []Announcement `json:"unreadAnnouncements"`
	HasUnreadAntenna                bool           `json:"hasUnreadAntenna"`
	HasUnreadChannel                bool           `json:"hasUnreadChannel"`
	HasUnreadChatMessages           bool           `json:"hasUnreadChatMessages"`
	HasUnreadNotification           bool           `json:"hasUnreadNotification"`
	HasPendingReceivedFollowRequest bool           `json:"hasPendingReceivedFollowRequest"`
	UnreadNotificationsCount        int            `json:"unreadNotificationsCount"`
	MutedWords                      [][]string     `json:"mutedWords"`
	HardMutedWords                  [][]string     `json:"hardMutedWords"`
	MutedInstances                  []string       `json:"mutedInstances"`
	NotificationRecieveConfig       map[string]any `json:"notificationRecieveConfig"`
	EmailNotificationTypes          []string       `json:"emailNotificationTypes"`
	Achievements                    []Achievement  `json:"achievements"`
	LoggedInDays                    int            `json:"loggedInDays"`
	Policies                        RolePolicies   `json:"policies"`
	TwoFactorEnabled                bool           `json:"twoFactorEnabled"`
	UsePasswordLessLogin            bool           `json:"usePasswordLessLogin"`
	SecurityKeys                    bool           `json:"securityKeys"`
	Email                           *string        `json:"email"`
	EmailVerified                   *bool          `json:"emailVerified"`
	SecurityKeysList                []SecurityKey  `json:"securityKeysList"`
}

// UserDetailed is the full view of a non-self user (UserLite + extra fields).
type UserDetailed struct {
	UserLite
	UserDetailedNotMeOnly
}

// MeDetailed is the full view of the authenticated user.
type MeDetailed struct {
	UserLite
	UserDetailedNotMeOnly
	MeDetailedOnly
}

// UserRole represents a role assigned to a user.
type UserRole struct {
	ID           ID      `json:"id"`
	Name         string  `json:"name"`
	Color        *string `json:"color,omitempty"`
	IconURL      *string `json:"iconUrl,omitempty"`
	DisplayOrder int     `json:"displayOrder"`
}

// RolePolicies represents the effective policies for a user's roles.
type RolePolicies map[string]any

// SecurityKey represents a WebAuthn security key.
type SecurityKey struct {
	ID       ID         `json:"id"`
	Name     string     `json:"name"`
	LastUsed *time.Time `json:"lastUsed,omitempty"`
}

// Achievement represents a user achievement.
type Achievement struct {
	Name       string `json:"name"`
	UnlockedAt int64  `json:"unlockedAt"`
}

// ---------------------------------------------------------------------------
// Note types
// ---------------------------------------------------------------------------

// Note represents a Misskey note.
type Note struct {
	ID                       ID                `json:"id"`
	CreatedAt                time.Time         `json:"createdAt"`
	DeletedAt                *time.Time        `json:"deletedAt,omitempty"`
	Text                     *string           `json:"text"`
	Cw                       *string           `json:"cw,omitempty"`
	UserID                   ID                `json:"userId"`
	User                     UserLite          `json:"user"`
	ReplyID                  *ID               `json:"replyId,omitempty"`
	RenoteID                 *ID               `json:"renoteId,omitempty"`
	Reply                    *Note             `json:"reply,omitempty"`
	Renote                   *Note             `json:"renote,omitempty"`
	IsHidden                 bool              `json:"isHidden,omitempty"`
	Visibility               NoteVisibility    `json:"visibility"`
	Mentions                 []ID              `json:"mentions,omitempty"`
	VisibleUserIds           []ID              `json:"visibleUserIds,omitempty"`
	FileIDs                  []ID              `json:"fileIds,omitempty"`
	Files                    []DriveFile       `json:"files,omitempty"`
	Tags                     []string          `json:"tags,omitempty"`
	Poll                     *Poll             `json:"poll,omitempty"`
	Emojis                   map[string]string `json:"emojis,omitempty"`
	ChannelID                *ID               `json:"channelId,omitempty"`
	Channel                  *Channel          `json:"channel,omitempty"`
	LocalOnly                bool              `json:"localOnly,omitempty"`
	ReactionAcceptance       *string           `json:"reactionAcceptance,omitempty"`
	ReactionEmojis           map[string]string `json:"reactionEmojis,omitempty"`
	Reactions                map[string]int    `json:"reactions,omitempty"`
	ReactionCount            int               `json:"reactionCount"`
	RenoteCount              int               `json:"renoteCount"`
	RepliesCount             int               `json:"repliesCount"`
	URI                      *string           `json:"uri,omitempty"`
	URL                      *string           `json:"url,omitempty"`
	ReactionAndUserPairCache []string          `json:"reactionAndUserPairCache,omitempty"`
	ClippedCount             int               `json:"clippedCount,omitempty"`
	HasPoll                  bool              `json:"hasPoll,omitempty"`
	MyReaction               *string           `json:"myReaction,omitempty"`
}

// Poll represents a poll attached to a note.
type Poll struct {
	Choices   []PollChoice `json:"choices"`
	Multiple  bool         `json:"multiple"`
	ExpiresAt *time.Time   `json:"expiresAt,omitempty"`
}

// PollChoice is a single choice in a poll.
type PollChoice struct {
	Text    string `json:"text"`
	Votes   int    `json:"votes"`
	IsVoted bool   `json:"isVoted"`
}

// Channel represents a Misskey channel.
type Channel struct {
	ID             ID         `json:"id"`
	CreatedAt      time.Time  `json:"createdAt"`
	LastNotedAt    *time.Time `json:"lastNotedAt,omitempty"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	BannerURL      *string    `json:"bannerUrl,omitempty"`
	BannerBlurhash *string    `json:"bannerBlurhash,omitempty"`
	PinnedNoteIds  []ID       `json:"pinnedNoteIds,omitempty"`
	Color          *string    `json:"color,omitempty"`
	IsArchived     bool       `json:"isArchived,omitempty"`
	UsersCount     int        `json:"usersCount,omitempty"`
	NotesCount     int        `json:"notesCount,omitempty"`
	IsSensitive    bool       `json:"isSensitive,omitempty"`
	IsFollowing    bool       `json:"isFollowing,omitempty"`
	HasUnreadNote  bool       `json:"hasUnreadNote,omitempty"`
}

// Announcement represents a server announcement.
type Announcement struct {
	ID                     ID         `json:"id"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              *time.Time `json:"updatedAt,omitempty"`
	Text                   string     `json:"text"`
	Title                  string     `json:"title"`
	ImageURL               *string    `json:"imageUrl,omitempty"`
	Icon                   string     `json:"icon"`
	Display                string     `json:"display"`
	NeedConfirmationToRead bool       `json:"needConfirmationToRead"`
	Silence                bool       `json:"silence"`
	ForYou                 bool       `json:"forYou"`
}

// ---------------------------------------------------------------------------
// Drive types
// ---------------------------------------------------------------------------

// DriveFile represents a file stored in the user's drive.
type DriveFile struct {
	ID           ID                  `json:"id"`
	CreatedAt    time.Time           `json:"createdAt"`
	Name         string              `json:"name"`
	Type         string              `json:"type"`
	Md5          string              `json:"md5"`
	Size         int64               `json:"size"`
	IsSensitive  bool                `json:"isSensitive"`
	Blurhash     *string             `json:"blurhash"`
	Properties   DriveFileProperties `json:"properties"`
	URL          string              `json:"url"`
	ThumbnailURL *string             `json:"thumbnailUrl"`
	Comment      *string             `json:"comment"`
	FolderID     *ID                 `json:"folderId"`
	Folder       *DriveFolder        `json:"folder,omitempty"`
	UserID       *ID                 `json:"userId"`
	User         *UserLite           `json:"user,omitempty"`
}

// DriveFileProperties contains file metadata properties.
type DriveFileProperties struct {
	Width       *int `json:"width,omitempty"`
	Height      *int `json:"height,omitempty"`
	Orientation *int `json:"orientation,omitempty"`
	Duration    *int `json:"duration,omitempty"`
}

// DriveFolder represents a folder in the user's drive.
type DriveFolder struct {
	ID        ID           `json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	Name      string       `json:"name"`
	ParentID  *ID          `json:"parentId,omitempty"`
	Parent    *DriveFolder `json:"parent,omitempty"`
}

// Drive represents the drive usage summary.
type Drive struct {
	Capacity int64 `json:"capacity"`
	Usage    int64 `json:"usage"`
}

// ---------------------------------------------------------------------------
// Other types
// ---------------------------------------------------------------------------

// Page represents a Misskey page.
type Page struct {
	ID                  ID          `json:"id"`
	CreatedAt           time.Time   `json:"createdAt"`
	UpdatedAt           *time.Time  `json:"updatedAt,omitempty"`
	Title               string      `json:"title"`
	Name                string      `json:"name"`
	Summary             *string     `json:"summary,omitempty"`
	Content             []any       `json:"content,omitempty"`
	Variables           []any       `json:"variables,omitempty"`
	User                UserLite    `json:"user"`
	HideTitleWhenPinned bool        `json:"hideTitleWhenPinned,omitempty"`
	AlignCenter         bool        `json:"alignCenter,omitempty"`
	Font                string      `json:"font,omitempty"`
	Script              string      `json:"script,omitempty"`
	EyeCatchingImage    *DriveFile  `json:"eyeCatchingImage,omitempty"`
	AttachedFiles       []DriveFile `json:"attachedFiles,omitempty"`
	LikedCount          int         `json:"likedCount,omitempty"`
	IsLiked             bool        `json:"isLiked,omitempty"`
}

// ---------------------------------------------------------------------------
// Chat types
// ---------------------------------------------------------------------------

// ChatMessage represents a direct chat message.
type ChatMessage struct {
	ID         ID         `json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	FromUserID ID         `json:"fromUserId"`
	FromUser   UserLite   `json:"fromUser"`
	ToUserID   *ID        `json:"toUserId"`
	ToUser     *UserLite  `json:"toUser,omitempty"`
	ToRoomID   *ID        `json:"toRoomId"`
	ToRoom     *ChatRoom  `json:"toRoom,omitempty"`
	Text       *string    `json:"text"`
	FileID     *ID        `json:"fileId"`
	File       *DriveFile `json:"file,omitempty"`
	IsRead     bool       `json:"isRead"`
}

// ChatRoom represents a group chat room.
type ChatRoom struct {
	ID               ID        `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	OwnerID          ID        `json:"ownerId"`
	Owner            UserLite  `json:"owner"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	IsMuted          bool      `json:"isMuted"`
	InvitationExists bool      `json:"invitationExists"`
}

// ChatRoomInvitation represents a chat room invitation.
type ChatRoomInvitation struct {
	ID        ID        `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	RoomID    ID        `json:"roomId"`
	Room      ChatRoom  `json:"room"`
	UserID    ID        `json:"userId"`
	User      UserLite  `json:"user"`
}

// ---------------------------------------------------------------------------
// Notification types
// ---------------------------------------------------------------------------

type NotificationType string

const (
	NotificationTypeFollow                NotificationType = "follow"
	NotificationTypeMention               NotificationType = "mention"
	NotificationTypeReply                 NotificationType = "reply"
	NotificationTypeRenote                NotificationType = "renote"
	NotificationTypeQuote                 NotificationType = "quote"
	NotificationTypeReaction              NotificationType = "reaction"
	NotificationTypePollVote              NotificationType = "pollVote"
	NotificationTypePollEnded             NotificationType = "pollEnded"
	NotificationTypeReceiveFollowRequest  NotificationType = "receiveFollowRequest"
	NotificationTypeFollowRequestAccepted NotificationType = "followRequestAccepted"
	NotificationTypeGroupInvited          NotificationType = "groupInvited"
	NotificationTypeApp                   NotificationType = "app"
)

// Notification represents a Missky notification.
// The Type field determines which other fields are populated.
type Notification struct {
	ID         ID                  `json:"id"`
	CreatedAt  time.Time           `json:"createdAt"`
	Type       NotificationType    `json:"type"`
	User       *UserLite           `json:"user,omitempty"`
	UserID     *ID                 `json:"userId,omitempty"`
	Note       *Note               `json:"note,omitempty"`
	Reaction   *string             `json:"reaction,omitempty"`
	Choice     *int                `json:"choice,omitempty"`
	Invitation *ChatRoomInvitation `json:"invitation,omitempty"`
	Body       *string             `json:"body,omitempty"`
	Header     *string             `json:"header,omitempty"`
	Icon       *string             `json:"icon,omitempty"`
}

// ---------------------------------------------------------------------------
// Admin types
// ---------------------------------------------------------------------------

// EmojiDetailed represents a custom emoji with full details.
type EmojiDetailed struct {
	ID                                      ID       `json:"id"`
	Aliases                                 []string `json:"aliases"`
	Name                                    string   `json:"name"`
	Category                                *string  `json:"category"`
	Host                                    *string  `json:"host"`
	URL                                     string   `json:"url"`
	License                                 *string  `json:"license"`
	IsSensitive                             bool     `json:"isSensitive"`
	LocalOnly                               bool     `json:"localOnly"`
	RoleIDsThatCanBeUsedThisEmojiAsReaction []ID     `json:"roleIdsThatCanBeUsedThisEmojiAsReaction"`
}

// AbuseReport represents an abuse report.
type AbuseReport struct {
	ID         ID         `json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	TargetUser UserLite   `json:"targetUser"`
	Reporter   *UserLite  `json:"reporter,omitempty"`
	Assignee   *UserLite  `json:"assignee,omitempty"`
	Resolved   bool       `json:"resolved"`
	Forwarded  bool       `json:"forwarded"`
	Comment    string     `json:"comment"`
	ResolvedAt *time.Time `json:"resolvedAt,omitempty"`
}

// Meta represents the server meta information.
type Meta struct {
	MaintainerName  *string `json:"maintainerName,omitempty"`
	MaintainerEmail *string `json:"maintainerEmail,omitempty"`
	Version         string  `json:"version"`
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	BannerURL       *string `json:"bannerUrl,omitempty"`
	IconURL         *string `json:"iconUrl,omitempty"`
	// ... many more fields
}

// ---------------------------------------------------------------------------
// Pagination
// ---------------------------------------------------------------------------

// Pagination contains common pagination parameters.
type Pagination struct {
	SinceID   *ID    `json:"sinceId,omitempty"`
	UntilID   *ID    `json:"untilId,omitempty"`
	SinceDate *int64 `json:"sinceDate,omitempty"`
	UntilDate *int64 `json:"untilDate,omitempty"`
	Limit     *int   `json:"limit,omitempty"` // Max: 100
}
