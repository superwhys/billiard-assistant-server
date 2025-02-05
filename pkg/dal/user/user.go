package userDal

import (
	"context"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/base"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"github.com/pkg/errors"
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
