package user

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/superwhys/snooker-assistant-server/pkg/wechat"
)

type IUserService interface {
	Login(ctx context.Context, username, password string) (*User, error)
	WechatLogin(ctx context.Context, wxSess *wechat.WechatSessionKeyResponse) (*User, error)
	Register(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userId int) error
	GetUserById(ctx context.Context, userId int) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	UploadAvatar(ctx context.Context, userId int, dest string, fh *multipart.FileHeader) (string, error)
	GetUserAvatar(ctx context.Context, avatarName string, dst io.Writer) error
}
