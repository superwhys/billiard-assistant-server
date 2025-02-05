// File:		game.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package handler

import (
	"context"
	"io"
	"mime/multipart"

	"gitea.hoven.com/billiard/billiard-assistant-server/api/middlewares"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitea.hoven.com/billiard/billiard-assistant-server/server"
	"gitea.hoven.com/billiard/billiard-assistant-server/server/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/pkg/errors"
)

type UserHandlerApp interface {
	UpdateUserName(ctx context.Context, token string, userName string) error
	UpdateUserGender(ctx context.Context, token string, gender int) error
	UploadAvatar(ctx context.Context, token string, fh *multipart.FileHeader) (string, error)
	GetAvatar(ctx context.Context, avatarId string, writer io.Writer) error
}

type UserHandler struct {
	userApp    UserHandlerApp
	middleware *middlewares.BilliardMiddleware
}

func NewUserHandler(server *server.BilliardServer, middleware *middlewares.BilliardMiddleware) *UserHandler {
	return &UserHandler{
		userApp:    server,
		middleware: middleware,
	}
}

func (u *UserHandler) Init(router gin.IRouter) {
	userNeedLogin := router.Group("user", u.middleware.UserLoginRequired())
	userNeedLogin.GET("info", pgin.ResponseHandler(u.getUserInfoHandler))
	userNeedLogin.PUT("nickname/update", pgin.RequestWithErrorHandler(u.updateUserNameHandler))
	userNeedLogin.PUT("gender/update", pgin.RequestWithErrorHandler(u.updateUserGenderHander))
	userNeedLogin.POST("avatar/upload", pgin.ResponseHandler(u.uploadAvatarHandler))
	userNeedLogin.GET("avatar/:avatar_id", pgin.RequestWithErrorHandler(u.getAvatarHandler))
}

func (u *UserHandler) getUserInfoHandler(ctx *gin.Context) (*dto.User, error) {
	user, err := u.middleware.CurrentUser(ctx)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetUserInfo
	}

	return dto.UserEntityToDto(user), nil
}

func (u *UserHandler) updateUserNameHandler(ctx *gin.Context, req *dto.UpdateUserNameRequest) error {
	token := u.middleware.GetLoginToken(ctx).GetKey()

	err := u.userApp.UpdateUserName(ctx, token, req.UserName)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrUpdateUserName
	}

	return nil
}

func (u *UserHandler) updateUserGenderHander(ctx *gin.Context, req *dto.UpdateUserGenderRequest) error {
	token := u.middleware.GetLoginToken(ctx).GetKey()

	err := u.userApp.UpdateUserGender(ctx, token, req.Gender)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrUpdateUserGender
	}

	return nil
}

func (u *UserHandler) uploadAvatarHandler(ctx *gin.Context) (*dto.UploadAvatarResponse, error) {
	token := u.middleware.GetLoginToken(ctx).GetKey()

	fh, err := ctx.FormFile("avatar")
	if err != nil {
		return nil, exception.ErrUploadAvatar
	}

	avatarUrl, err := u.userApp.UploadAvatar(ctx, token, fh)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrUploadAvatar
	}

	return &dto.UploadAvatarResponse{AvatarUrl: avatarUrl}, nil
}

func (u *UserHandler) getAvatarHandler(ctx *gin.Context, req *dto.GetUserAvatarRequest) error {
	return u.userApp.GetAvatar(ctx.Request.Context(), req.AvatarId, ctx.Writer)
}
