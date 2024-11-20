package user

import (
	"context"
	"io"
	"mime/multipart"
)

type IUserService interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	DeleteUser(ctx context.Context, userId int) error
	GetUserById(ctx context.Context, userId int) (*User, error)
	GetUserWithRoom(ctx context.Context, userId int) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	UpdateUserStatus(ctx context.Context, userId int, status Status) error
	UpdateUserRole(ctx context.Context, userId int, role Role) error

	// Avatar management
	UploadAvatar(ctx context.Context, userId int, fh *multipart.FileHeader) (string, error)
	GetUserAvatar(ctx context.Context, avatarName string, dst io.Writer) error
}
