package userDal

import (
	"context"
	"errors"

	"github.com/superwhys/snooker-assistant-server/domain/user"
	"github.com/superwhys/snooker-assistant-server/pkg/dal/base"
	"github.com/superwhys/snooker-assistant-server/pkg/dal/model"
	"github.com/superwhys/snooker-assistant-server/pkg/exception"
	"gorm.io/gorm"
)

var _ user.IUserRepo = (*UserRepoImpl)(nil)

type UserRepoImpl struct {
	db *base.Query
}

func NewUserRepo(db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{base.Use(db)}
}

func (u *UserRepoImpl) CreateUser(ctx context.Context, user *user.User) error {
	up := new(model.UserPo)
	up.FromEntity(user)

	userDb := u.db.UserPo
	return userDb.WithContext(ctx).Create(up)
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
