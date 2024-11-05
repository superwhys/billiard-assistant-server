package user

import (
	"context"
	"mime/multipart"
)

type IUserService interface {
	Login(ctx context.Context, username, password string) (*User, error)
	Register(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userId int) error
	GetUserById(ctx context.Context, userId int) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	UploadAvatar(ctx context.Context, userId int, dest string, fh *multipart.FileHeader) (string, error)
}
