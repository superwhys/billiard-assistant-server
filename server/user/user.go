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

	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/auth"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/oss"
)

var _ user.IUserService = (*UserService)(nil)

type UserService struct {
	userRepo user.IUserRepo
	authRepo auth.IAuthRepo
	oss      oss.IOSS
}

func NewUserService(userRepo user.IUserRepo, authRepo auth.IAuthRepo, oss oss.IOSS) *UserService {
	return &UserService{userRepo: userRepo, authRepo: authRepo, oss: oss}
}

func (us *UserService) UserExists(ctx context.Context, userId int) (bool, error) {
	return us.userRepo.UserExists(ctx, userId)
}

func (us *UserService) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	existingUser, err := us.userRepo.GetUserByName(ctx, u.Name)
	if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return nil, errors.Wrap(err, "getUserByName")
	}
	if existingUser != nil {
		return nil, exception.ErrUserAlreadyExists
	}

	if u, err = us.userRepo.CreateUser(ctx, u); err != nil {
		return nil, errors.Wrap(err, "createUser")
	}

	return u, nil
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

func (us *UserService) GetUserByName(ctx context.Context, userName string) (*user.User, error) {
	return us.userRepo.GetUserByName(ctx, userName)
}

func (us *UserService) checkEqual(ou, nu *user.User) bool {
	return ou.UserId == nu.UserId &&
		ou.Name == nu.Name &&
		ou.UserInfo == nu.UserInfo &&
		ou.Status == nu.Status
}

func (us *UserService) UpdateUser(ctx context.Context, update *user.User) (*user.User, error) {
	oldUser, err := us.GetUserById(ctx, update.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "getUserById")
	}

	if !oldUser.HasUpdate(update) {
		return oldUser, nil
	}

	oldUser.Name = update.Name
	oldUser.UserInfo = update.UserInfo
	oldUser.Status = update.Status
	oldUser.Gender = update.Gender

	if err := us.userRepo.UpdateUser(ctx, oldUser); err != nil {
		return nil, errors.Wrap(err, "updateUser")
	}

	return oldUser, nil
}

func (us *UserService) UploadAvatar(ctx context.Context, userId int, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	avatarUrl, err := us.oss.UploadFile(ctx, file.Size, oss.SourceAvatar, file.Filename, src)
	if err != nil {
		return "", errors.Wrap(err, "uploadAvatar")
	}

	_, err = us.UpdateUser(ctx, &user.User{
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

func (us *UserService) UpdateUserStatus(ctx context.Context, userId int, status user.Status) error {
	u, err := us.GetUserById(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "getUserById")
	}
	u.Status = status
	return us.userRepo.UpdateUser(ctx, u)
}

func (us *UserService) UpdateUserRole(ctx context.Context, userId int, role user.Role) error {
	u, err := us.GetUserById(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "getUserById")
	}
	u.Role = role
	return us.userRepo.UpdateUser(ctx, u)
}
