package migo

import (
	"github.com/sysnote8main/migo/account"
	"github.com/sysnote8main/migo/admin"
	"github.com/sysnote8main/migo/announcements"
	"github.com/sysnote8main/migo/antennas"
	"github.com/sysnote8main/migo/ap"
	"github.com/sysnote8main/migo/app"
	"github.com/sysnote8main/migo/auth"
	"github.com/sysnote8main/migo/blocking"
	"github.com/sysnote8main/migo/bubblegame"
	"github.com/sysnote8main/migo/channels"
	"github.com/sysnote8main/migo/charts"
	"github.com/sysnote8main/migo/chat"
	"github.com/sysnote8main/migo/clips"
	"github.com/sysnote8main/migo/drive"
	"github.com/sysnote8main/migo/emailaddress"
	"github.com/sysnote8main/migo/emoji"
	"github.com/sysnote8main/migo/emojis"
	"github.com/sysnote8main/migo/endpoint"
	"github.com/sysnote8main/migo/export"
	"github.com/sysnote8main/migo/federation"
	"github.com/sysnote8main/migo/fetch"
	"github.com/sysnote8main/migo/flash"
	"github.com/sysnote8main/migo/following"
	"github.com/sysnote8main/migo/gallery"
	"github.com/sysnote8main/migo/get"
	"github.com/sysnote8main/migo/hashtags"
	"github.com/sysnote8main/migo/invite"
	"github.com/sysnote8main/migo/meta"
	"github.com/sysnote8main/migo/misc"
	"github.com/sysnote8main/migo/mute"
	"github.com/sysnote8main/migo/notes"
	"github.com/sysnote8main/migo/notification"
	"github.com/sysnote8main/migo/pages"
	"github.com/sysnote8main/migo/promo"
	"github.com/sysnote8main/migo/renotemute"
	"github.com/sysnote8main/migo/reversi"
	"github.com/sysnote8main/migo/roles"
	"github.com/sysnote8main/migo/sw"
	"github.com/sysnote8main/migo/test"
	"github.com/sysnote8main/migo/timeline"
	"github.com/sysnote8main/migo/username"
	"github.com/sysnote8main/migo/users"
)

// Services bundles all service interfaces for convenient access.
type Services struct {
	Account       *account.Service
	Admin         *admin.Service
	Announcements *announcements.Service
	Antennas      *antennas.Service
	Ap            *ap.Service
	App           *app.Service
	Auth          *auth.Service
	Blocking      *blocking.Service
	BubbleGame    *bubblegame.Service
	Channels      *channels.Service
	Charts        *charts.Service
	Chat          *chat.Service
	Clips         *clips.Service
	Drive         *drive.Service
	EmailAddress  *emailaddress.Service
	Emoji         *emoji.Service
	Emojis        *emojis.Service
	Endpoint      *endpoint.Service
	Export        *export.Service
	Federation    *federation.Service
	Fetch         *fetch.Service
	Flash         *flash.Service
	Following     *following.Service
	Gallery       *gallery.Service
	Get           *get.Service
	Hashtags      *hashtags.Service
	Invite        *invite.Service
	Meta          *meta.Service
	Misc          *misc.Service
	Mute          *mute.Service
	Notes         *notes.Service
	Notification  *notification.Service
	Pages         *pages.Service
	Promo         *promo.Service
	RenoteMute    *renotemute.Service
	Reversi       *reversi.Service
	Roles         *roles.Service
	Sw            *sw.Service
	Test          *test.Service
	Timeline      *timeline.Service
	Users         *users.Service
	Username      *username.Service
}

// NewServices creates a bundle of all service instances using the given client.
func NewServices(cli Client) *Services {
	return &Services{
		Account:       account.NewService(cli),
		Admin:         admin.NewService(cli),
		Announcements: announcements.NewService(cli),
		Antennas:      antennas.NewService(cli),
		Ap:            ap.NewService(cli),
		App:           app.NewService(cli),
		Auth:          auth.NewService(cli),
		Blocking:      blocking.NewService(cli),
		BubbleGame:    bubblegame.NewService(cli),
		Channels:      channels.NewService(cli),
		Charts:        charts.NewService(cli),
		Chat:          chat.NewService(cli),
		Clips:         clips.NewService(cli),
		Drive:         drive.NewService(cli),
		EmailAddress:  emailaddress.NewService(cli),
		Emoji:         emoji.NewService(cli),
		Emojis:        emojis.NewService(cli),
		Endpoint:      endpoint.NewService(cli),
		Export:        export.NewService(cli),
		Federation:    federation.NewService(cli),
		Fetch:         fetch.NewService(cli),
		Flash:         flash.NewService(cli),
		Following:     following.NewService(cli),
		Gallery:       gallery.NewService(cli),
		Get:           get.NewService(cli),
		Hashtags:      hashtags.NewService(cli),
		Invite:        invite.NewService(cli),
		Meta:          meta.NewService(cli),
		Misc:          misc.NewService(cli),
		Mute:          mute.NewService(cli),
		Notes:         notes.NewService(cli),
		Notification:  notification.NewService(cli),
		Pages:         pages.NewService(cli),
		Promo:         promo.NewService(cli),
		RenoteMute:    renotemute.NewService(cli),
		Reversi:       reversi.NewService(cli),
		Roles:         roles.NewService(cli),
		Sw:            sw.NewService(cli),
		Test:          test.NewService(cli),
		Timeline:      timeline.NewService(cli),
		Users:         users.NewService(cli),
		Username:      username.NewService(cli),
	}
}
