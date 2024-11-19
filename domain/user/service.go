package user

import (
	"context"
	"io"
	"mime/multipart"
)

type IUserService interface {
	// Login(ctx context.Context, username, password string) (*User, error)
	CreateUser(ctx context.Context, u *User) (*User, error)
	DeleteUser(ctx context.Context, userId int) error
	GetUserById(ctx context.Context, userId int) (*User, error)
	GetUserWithRoom(ctx context.Context, userId int) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)

	// User status and role management
	UpdateUserStatus(ctx context.Context, userId int, status Status) error
	UpdateUserRole(ctx context.Context, userId int, role Role) error

	// Avatar management
	UploadAvatar(ctx context.Context, userId int, dest string, fh *multipart.FileHeader) (string, error)
	GetUserAvatar(ctx context.Context, avatarName string, dst io.Writer) error
}
