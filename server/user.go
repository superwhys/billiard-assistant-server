// File:		user.go
// Created by:	Hoven
// Created on:	2025-01-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"context"
	"io"
	"mime/multipart"

	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
)

// UpdateUserName
// FIX: maybe need to use locker to avoid concurrency problem caused by update a same name
func (s *BilliardServer) UpdateUserName(ctx context.Context, token string, userName string) error {
	err := s.UserSrv.UpdateUserName(ctx, token, userName)
	if errors.Is(err, exception.ErrUserNotFound) {
		return exception.ErrUserNotFound
	} else if err != nil {
		plog.Errorc(ctx, "update user name error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) UpdateUserGender(ctx context.Context, token string, gender int) error {
	err := s.UserSrv.UpdateUserGender(ctx, token, gender)
	if errors.Is(err, exception.ErrUserNotFound) {
		return exception.ErrUserNotFound
	} else if err != nil {
		plog.Errorc(ctx, "update user gender error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) UploadAvatar(ctx context.Context, token string, file *multipart.FileHeader) (string, error) {
	avatarUrl, err := s.UserSrv.UploadAvatar(ctx, token, file)
	if err != nil {
		plog.Errorc(ctx, "upload avatar error: %v", err)
		return "", exception.ErrUploadAvatar
	}

	return avatarUrl, nil
}

func (s *BilliardServer) GetAvatar(ctx context.Context, avatarId string, writer io.Writer) error {
	err := s.UserSrv.GetAvatar(ctx, avatarId, writer)
	if errors.Is(err, exception.ErrGetAvatar) {
		return exception.ErrGetAvatar
	} else if err != nil {
		plog.Errorc(ctx, "get avatar error: %v", err)
		return err
	}

	return nil
}
