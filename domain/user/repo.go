package user

import (
	"context"
)

type IUserRepo interface {
	GetUser(ctx context.Context, userId int) (*User, error)
	CreateUser(ctx context.Context, u *User) (*User, error)
	DeleteUser(ctx context.Context, userId int) error
	UpdateUser(ctx context.Context, u *User) error
}
