package user

import (
	"context"
)

type IUserRepo interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	DeleteUser(ctx context.Context, userId int) error
	UpdateUser(ctx context.Context, u *User) error
}
