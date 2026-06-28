package migo

import (
	"github.com/sysnote8main/migo/admin"
	"github.com/sysnote8main/migo/auth"
	"github.com/sysnote8main/migo/chat"
	"github.com/sysnote8main/migo/drive"
	"github.com/sysnote8main/migo/following"
	"github.com/sysnote8main/migo/notes"
	"github.com/sysnote8main/migo/notification"
	"github.com/sysnote8main/migo/timeline"
	"github.com/sysnote8main/migo/users"
)

// Services bundles all service interfaces for convenient access.
//
//	svc := migo.NewServices(client)
//	note, _ := svc.Notes.Create(ctx, &notes.CreateRequest{Text: proto.String("hi")})
type Services struct {
	Notes        *notes.Service
	Users        *users.Service
	Drive        *drive.Service
	Auth         *auth.Service
	Timeline     *timeline.Service
	Following    *following.Service
	Notification *notification.Service
	Admin        *admin.Service
	Chat         *chat.Service
}

// NewServices creates a bundle of all service instances using the given client.
func NewServices(cli Client) *Services {
	return &Services{
		Notes:        notes.NewService(cli),
		Users:        users.NewService(cli),
		Drive:        drive.NewService(cli),
		Auth:         auth.NewService(cli),
		Timeline:     timeline.NewService(cli),
		Following:    following.NewService(cli),
		Notification: notification.NewService(cli),
		Admin:        admin.NewService(cli),
		Chat:         chat.NewService(cli),
	}
}
