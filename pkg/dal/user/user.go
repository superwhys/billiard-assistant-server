package userDal

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/base"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gorm.io/gorm"
)

var _ user.IUserRepo = (*UserRepoImpl)(nil)

type UserRepoImpl struct {
	db *base.Query
}

func NewUserRepo(db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{base.Use(db)}
}

// Basic user operations
func (u *UserRepoImpl) CreateUser(ctx context.Context, user *user.User) (*user.User, error) {
	up := new(model.UserPo)
	up.FromEntity(user)

	userDb := u.db.UserPo
	if err := userDb.WithContext(ctx).Create(up); err != nil {
		return nil, errors.Wrap(err, "createUser")
	}
	user.UserId = up.ID
	return user, nil
}

func (u *UserRepoImpl) DeleteUser(ctx context.Context, userId int) error {
	userDb := u.db.UserPo
	_, err := userDb.WithContext(ctx).Where(userDb.ID.Eq(userId)).Delete()
	return err
}

func (u *UserRepoImpl) UpdateUser(ctx context.Context, user *user.User) error {
	up := new(model.UserPo)
	up.FromEntity(user)

	userDb := u.db.UserPo
	_, err := userDb.WithContext(ctx).Where(userDb.ID.Eq(up.ID)).Updates(up)
	return err
}

func (u *UserRepoImpl) GetUserById(ctx context.Context, userId int) (*user.User, error) {
	userDb := u.db.UserPo
	usr, err := userDb.WithContext(ctx).
		Where(userDb.ID.Eq(userId)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return usr.ToEntity(), nil
}

func (u *UserRepoImpl) GetUserByName(ctx context.Context, username string) (*user.User, error) {
	userDb := u.db.UserPo
	usr, err := userDb.WithContext(ctx).
		Where(userDb.Name.Eq(username)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return usr.ToEntity(), nil
}

func (u *UserRepoImpl) GetUserWithRoomById(ctx context.Context, userId int) (*user.User, error) {
	userDb := u.db.UserPo
	usr, err := userDb.WithContext(ctx).
		Preload(userDb.Rooms).
		Preload(userDb.Rooms.Game).
		Where(userDb.ID.Eq(userId)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return usr.ToEntity(), nil
}

// User status management
func (u *UserRepoImpl) UpdateUserStatus(ctx context.Context, userId int, status user.Status) error {
	userDb := u.db.UserPo
	_, err := userDb.WithContext(ctx).
		Where(userDb.ID.Eq(userId)).
		UpdateSimple(userDb.Status.Value(int(status)))
	return err
}

func (u *UserRepoImpl) UpdateUserRole(ctx context.Context, userId int, role user.Role) error {
	userDb := u.db.UserPo
	_, err := userDb.WithContext(ctx).
		Where(userDb.ID.Eq(userId)).
		UpdateSimple(userDb.Role.Value(int(role)))
	return err
}
