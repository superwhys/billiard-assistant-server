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
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/password"
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

func (us *UserService) Login(ctx context.Context, username, pwd string) (*user.User, error) {
	u, err := us.userRepo.GetUserByName(ctx, username)
	if err != nil {
		return nil, errors.Wrap(err, "getUserByName")
	}

	userAuth, err := us.authRepo.GetUserAuthByType(ctx, u.UserId, auth.AuthTypePassword)
	if err != nil {
		return nil, errors.Wrap(err, "getUserAuth")
	}

	if !password.CheckPasswordHash(pwd, userAuth.Credential) {
		return nil, exception.ErrPasswordNotCorrect
	}

	if u, err = us.UpdateUser(ctx, u); err != nil {
		return nil, exception.ErrUpdateUserInfo
	}

	return u, nil
}

func (us *UserService) Register(ctx context.Context, username string) (*user.User, error) {
	existingUser, err := us.userRepo.GetUserByName(ctx, username)
	if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return nil, errors.Wrap(err, "getUserByName")
	}
	if existingUser != nil {
		return nil, exception.ErrUserAlreadyExists
	}

	newUser := &user.User{
		Name:   username,
		Status: user.StatusActive,
	}

	if newUser, err = us.userRepo.CreateUser(ctx, newUser); err != nil {
		return nil, errors.Wrap(err, "createUser")
	}

	return newUser, nil
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

func (us *UserService) GetUserWithRoom(ctx context.Context, userId int) (*user.User, error) {
	return us.userRepo.GetUserWithRoomById(ctx, userId)
}

func (us *UserService) UpdateUser(ctx context.Context, update *user.User) (*user.User, error) {
	oldUser, err := us.GetUserById(ctx, update.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "getUserById")
	}

	oldUser.Name = update.Name
	oldUser.UserInfo = update.UserInfo
	oldUser.Status = update.Status

	if err := us.userRepo.UpdateUser(ctx, oldUser); err != nil {
		return nil, errors.Wrap(err, "updateUser")
	}

	return oldUser, nil
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
