package user

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/superwhys/snooker-assistant-server/pkg/wechat"
)

type IUserService interface {
	// Authentication related
	Login(ctx context.Context, username, password string) (*User, error)
	WechatLogin(ctx context.Context, wxSess *wechat.WechatSessionKeyResponse) (*User, error)
	Register(ctx context.Context, username, password string) (*User, error)

	// Basic user operations
	DeleteUser(ctx context.Context, userId int) error
	GetUserById(ctx context.Context, userId int) (*User, error)
	GetUserWithRoom(ctx context.Context, userId int) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)

	// Auth management
	CreateUserAuth(ctx context.Context, userId int, auth *UserAuth) error
	UpdateUserAuth(ctx context.Context, auth *UserAuth) error
	DeleteUserAuth(ctx context.Context, authId int) error
	GetUserAuths(ctx context.Context, userId int) ([]*UserAuth, error)

	// User status and role management
	UpdateUserStatus(ctx context.Context, userId int, status Status) error
	UpdateUserRole(ctx context.Context, userId int, role Role) error

	// Avatar management
	UploadAvatar(ctx context.Context, userId int, dest string, fh *multipart.FileHeader) (string, error)
	GetUserAvatar(ctx context.Context, avatarName string, dst io.Writer) error
}
