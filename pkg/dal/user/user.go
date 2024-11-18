package userDal

import (
	"context"
	
	"github.com/pkg/errors"
	"github.com/superwhys/billiard-assistant-server/domain/user"
	"github.com/superwhys/billiard-assistant-server/pkg/dal/base"
	"github.com/superwhys/billiard-assistant-server/pkg/dal/model"
	"github.com/superwhys/billiard-assistant-server/pkg/exception"
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
func (u *UserRepoImpl) CreateUser(ctx context.Context, user *user.User, userAuth *user.UserAuth) error {
	up := new(model.UserPo)
	up.FromEntity(user)
	
	userDb := u.db.UserPo
	if err := userDb.WithContext(ctx).Create(up); err != nil {
		return errors.Wrap(err, "createUser")
	}
	user.UserId = up.ID
	userAuth.UserId = up.ID
	
	return u.CreateUserAuth(ctx, up.ID, userAuth)
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

// User queries
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

// User authentication operations
func (u *UserRepoImpl) CreateUserAuth(ctx context.Context, userId int, auth *user.UserAuth) error {
	userDb := u.db.UserPo
	userAuthDb := u.db.UserAuthPo
	
	userAuthPo := new(model.UserAuthPo)
	userAuthPo.FromEntity(auth)
	userAuthPo.UserPoID = userId
	if err := userAuthDb.WithContext(ctx).Create(userAuthPo); err != nil {
		return err
	}
	
	return userDb.UserAuthPos.Model(&model.UserPo{ID: userId}).Append(userAuthPo)
}

func (u *UserRepoImpl) UpdateUserAuth(ctx context.Context, auth *user.UserAuth) error {
	userAuthPo := new(model.UserAuthPo)
	userAuthPo.FromEntity(auth)
	
	userAuthDb := u.db.UserAuthPo
	_, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.ID.Eq(auth.Id)).
		Updates(userAuthPo)
	return err
}

func (u *UserRepoImpl) DeleteUserAuth(ctx context.Context, authId int) error {
	userAuthDb := u.db.UserAuthPo
	_, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.ID.Eq(authId)).
		Delete()
	return err
}

func (u *UserRepoImpl) GetUserByAuth(ctx context.Context, authType user.AuthType, identifier string) (*user.User, error) {
	userDb := u.db.UserPo
	userAuthDb := u.db.UserAuthPo
	
	auth, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.AuthType.Eq(int(authType))).
		Where(userAuthDb.Identifier.Eq(identifier)).
		First()
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	
	usr, err := userDb.WithContext(ctx).
		Where(userDb.ID.Eq(auth.UserPoID)).
		First()
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	
	return usr.ToEntity(), nil
}

// User authentication queries
func (u *UserRepoImpl) GetUserAuths(ctx context.Context, userId int) ([]*user.UserAuth, error) {
	userAuthDb := u.db.UserAuthPo
	auths, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.UserPoID.Eq(userId)).
		Find()
	if err != nil {
		return nil, err
	}
	
	result := make([]*user.UserAuth, 0, len(auths))
	for _, auth := range auths {
		result = append(result, auth.ToEntity())
	}
	return result, nil
}

func (u *UserRepoImpl) GetUserAuthByIdentifier(ctx context.Context, authType user.AuthType, identifier string) (*user.UserAuth, error) {
	userAuthDb := u.db.UserAuthPo
	auth, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.AuthType.Eq(int(authType))).
		Where(userAuthDb.Identifier.Eq(identifier)).
		First()
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserAuthNotFound
	} else if err != nil {
		return nil, err
	}
	
	return auth.ToEntity(), nil
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

// Extended query methods
func (u *UserRepoImpl) GetUserByWechatId(ctx context.Context, wechatId string) (*user.User, error) {
	return u.GetUserByAuth(ctx, user.AuthTypeWechat, wechatId)
}

func (u *UserRepoImpl) GetUserAuthByType(ctx context.Context, userId int, authType user.AuthType) (*user.UserAuth, error) {
	userAuthDb := u.db.UserAuthPo
	auth, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.UserPoID.Eq(userId)).
		Where(userAuthDb.AuthType.Eq(int(authType))).
		First()
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserAuthNotFound
	} else if err != nil {
		return nil, err
	}
	
	return auth.ToEntity(), nil
}

func (u *UserRepoImpl) GetUserWithRoom(ctx context.Context, userId int) (*user.User, error) {
	return u.GetUserWithRoomById(ctx, userId)
}
