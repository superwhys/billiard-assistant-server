// File:		game.go
// Created by:	Hoven
// Created on:	2024-10-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package userSrv

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/pkg/errors"
	"github.com/superwhys/snooker-assistant-server/domain/user"
	"github.com/superwhys/snooker-assistant-server/pkg/exception"
	"github.com/superwhys/snooker-assistant-server/pkg/oss"
	"github.com/superwhys/snooker-assistant-server/pkg/password"
	"github.com/superwhys/snooker-assistant-server/pkg/wechat"
)

var _ user.IUserService = (*UserService)(nil)

type UserService struct {
	userRepo user.IUserRepo
	oss      oss.IOSS
}

func NewUserService(userRepo user.IUserRepo, oss oss.IOSS) *UserService {
	return &UserService{userRepo: userRepo, oss: oss}
}

func (us *UserService) Login(ctx context.Context, username, pwd string) (*user.User, error) {
	u, err := us.userRepo.GetUserByName(ctx, username)
	if err != nil {
		return nil, errors.Wrap(err, "getUserByName")
	}

	if !password.CheckPasswordHash(pwd, u.Password) {
		return nil, exception.ErrPasswordNotCorrect
	}
	u.Password = ""
	u.LastLoginAt = time.Now()
	if err := us.UpdateUser(ctx, u); err != nil {
		return nil, exception.ErrUpdateUserInfo
	}

	return u, nil
}

func (us *UserService) WechatLogin(ctx context.Context, wxSess *wechat.WechatSessionKeyResponse) (*user.User, error) {
	u, err := us.userRepo.GetUserByWechatId(ctx, wxSess.OpenID)
	if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return nil, errors.Wrap(err, "GetUserByWechatOpenId")
	}

	if u == nil {
		u = &user.User{
			Name:     wxSess.OpenID,
			WechatId: wxSess.OpenID,
			Status:   user.StatusActive,
		}
		if err := us.userRepo.CreateUser(ctx, u); err != nil {
			return nil, errors.Wrap(err, "CreateUser")
		}
	}

	return u, nil
}

func (us *UserService) Register(ctx context.Context, u *user.User) error {
	ret, err := us.userRepo.GetUserByName(ctx, u.Name)
	if ret != nil {
		return exception.ErrUserAlreadyExists
	}

	hashPwd, err := password.HashPassword(u.Password)
	if err != nil {
		return errors.Wrap(err, "hashPassword")
	}
	u.Password = hashPwd
	u.Status = user.StatusActive

	return us.userRepo.CreateUser(ctx, u)
}

func (us *UserService) DeleteUser(ctx context.Context, userId int) error {
	_, err := us.GetUserById(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "getUserById")
	}

	return us.userRepo.DeleteUser(ctx, userId)
}

func (us *UserService) GetUserById(ctx context.Context, userId int) (*user.User, error) {
	return us.userRepo.GetUserById(ctx, userId)
}

func (us *UserService) UpdateUser(ctx context.Context, user *user.User) error {
	oldUser, err := us.GetUserById(ctx, user.UserId)
	if err != nil {
		return errors.Wrap(err, "getUserById")
	}

	if user.Password != "" {
		hashPwd, err := password.HashPassword(user.Password)
		if err != nil {
			return errors.Wrap(err, "hashPassword")
		}
		user.Password = hashPwd
	}

	// Only this following fields can be updated
	oldUser.Name = user.Name
	oldUser.Password = user.Password
	oldUser.UserInfo = user.UserInfo
	oldUser.Status = user.Status
	oldUser.LastLoginAt = user.LastLoginAt

	return us.userRepo.UpdateUser(ctx, oldUser)
}

func (us *UserService) UploadAvatar(ctx context.Context, userId int, dest string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	objName := fmt.Sprintf("%s/%s", dest, file.Filename)
	avatarUrl, err := us.oss.UploadFile(ctx, file.Size, objName, src)
	if err != nil {
		return "", errors.Wrap(err, "uploadAvatar")
	}

	err = us.UpdateUser(ctx, &user.User{
		UserId: userId,
		UserInfo: &user.BaseInfo{
			Avatar: avatarUrl,
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "updateUserAvatarUrl")
	}

	return avatarUrl, nil
}

func (us *UserService) GetUserAvatar(ctx context.Context, avatarName string, dst io.Writer) error {
	objName := fmt.Sprintf("avatar/%s", avatarName)
	return us.oss.GetFile(ctx, objName, dst)
}
