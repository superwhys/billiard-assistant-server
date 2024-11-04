package user

import (
	"context"
)

type IUserRepo interface {
	CreateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, userId int) error
	UpdateUser(ctx context.Context, u *User) error
	GetUserById(ctx context.Context, userId int) (*User, error)
	GetUserByName(ctx context.Context, username string) (*User, error)
	GetUserWithRoomById(ctx context.Context, userId int) (*User, error)
}
