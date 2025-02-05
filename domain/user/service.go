package user

import (
	"context"
	"io"
	"mime/multipart"
)

type IUserService interface {
	GetUserProfile(ctx context.Context, token string) (*User, error)
	UpdateUserName(ctx context.Context, token string, name string) error
	UpdateUserGender(ctx context.Context, token string, gender int) error
	UploadAvatar(ctx context.Context, token string, fh *multipart.FileHeader) (string, error)
	GetAvatar(ctx context.Context, avatarId string, writer io.Writer) error
}
