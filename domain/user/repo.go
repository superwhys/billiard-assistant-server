package user

import (
	"context"
)

type IUserRepo interface {
	CreateUser(ctx context.Context, u *User, ua *UserAuth) error
	DeleteUser(ctx context.Context, userId int) error
	UpdateUser(ctx context.Context, u *User) error

	GetUserById(ctx context.Context, userId int) (*User, error)
	GetUserByName(ctx context.Context, username string) (*User, error)
	GetUserWithRoomById(ctx context.Context, userId int) (*User, error)

	CreateUserAuth(ctx context.Context, userId int, auth *UserAuth) error
	UpdateUserAuth(ctx context.Context, auth *UserAuth) error
	DeleteUserAuth(ctx context.Context, authId int) error
	GetUserByAuth(ctx context.Context, authType AuthType, identifier string) (*User, error)

	GetUserAuths(ctx context.Context, userId int) ([]*UserAuth, error)
	GetUserAuthByIdentifier(ctx context.Context, authType AuthType, identifier string) (*UserAuth, error)

	UpdateUserStatus(ctx context.Context, userId int, status Status) error
	UpdateUserRole(ctx context.Context, userId int, role Role) error

	GetUserByWechatId(ctx context.Context, wechatId string) (*User, error)
	GetUserAuthByType(ctx context.Context, userId int, authType AuthType) (*UserAuth, error)
	GetUserWithRoom(ctx context.Context, userId int) (*User, error)
}
