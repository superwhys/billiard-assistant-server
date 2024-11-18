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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/snooker-assistant-server/api/middlewares"
	"github.com/superwhys/snooker-assistant-server/pkg/exception"
	"github.com/superwhys/snooker-assistant-server/pkg/wechat"
	"github.com/superwhys/snooker-assistant-server/server"
	"github.com/superwhys/snooker-assistant-server/server/dto"
)

type UserHandlerApp interface {
	Login(ctx context.Context, username, password string) (*dto.User, error)
	WechatLogin(ctx context.Context, code string) (*dto.User, *wechat.WechatSessionKeyResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error)
	UpdateUser(ctx context.Context, userId int, update *dto.UpdateUserRequest) error
	UploadAvatar(ctx context.Context, userId int, fh *multipart.FileHeader) (string, error)
	GetAvatar(ctx context.Context, avatarName string, dst io.Writer) error
}

type UserHandler struct {
	userApp    UserHandlerApp
	middleware *middlewares.SaMiddleware
}

func NewUserHandler(server *server.SaServer, middleware *middlewares.SaMiddleware) *UserHandler {
	return &UserHandler{
		userApp:    server,
		middleware: middleware,
	}
}

func (u *UserHandler) Init(router gin.IRouter) {
	user := router.Group("user")
	user.POST("login", pgin.RequestResponseHandler(u.loginHandler))
	user.POST("login/wx", pgin.RequestResponseHandler(u.wechatLoginHandler))
	user.POST("register", pgin.RequestResponseHandler(u.registerHandler))
	user.GET("avatar/:avatar_name", pgin.RequestHandler(u.getUserAvatarHandler))

	userAuth := router.Group("user", u.middleware.UserLoginRequired())
	userAuth.GET("info", pgin.ResponseHandler(u.getUserInfoHandler))
	userAuth.PUT("info/update", pgin.RequestWithErrorHandler(u.updateUserHandler))
	userAuth.POST("avatar/upload", pgin.ResponseHandler(u.uploadAvatarHandler))
}

func (u *UserHandler) getUserAvatarHandler(ctx *gin.Context, req *dto.GetUserAvatarRequest) {
	err := u.userApp.GetAvatar(ctx, req.AvatarName, ctx.Writer)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
}

func (u *UserHandler) wechatLoginHandler(ctx *gin.Context, req *dto.WechatLoginRequest) (*dto.WechatLoginResponse, error) {
	plog.Debugc(ctx, "wechat login code: %s", req.Code)

	user, wxSessionKey, err := u.userApp.WechatLogin(ctx, req.Code)
	if exception.CheckSaException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrLoginFailed
	}

	token := &middlewares.UserToken{
		Uid:              user.UserId,
		WechatId:         wxSessionKey.OpenID,
		WechatUnionId:    wxSessionKey.UnionID,
		WechatSessionKey: wxSessionKey.SessionKey,
	}
	u.middleware.SaveToken(token, ctx)
	return &dto.WechatLoginResponse{Token: token.GetKey(), User: user}, nil
}

func (u *UserHandler) loginHandler(ctx *gin.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.userApp.Login(ctx, req.Username, req.Password)
	if exception.CheckSaException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrLoginFailed
	}

	token := middlewares.NewUserToken(user.UserId, user.WechatId, user.Name)
	u.middleware.SaveToken(token, ctx)
	return &dto.LoginResponse{Token: token.GetKey(), User: user}, nil
}

func (u *UserHandler) registerHandler(ctx *gin.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	user, err := u.userApp.Register(ctx, req)
	if exception.CheckSaException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrRegisterUser
	}

	return &dto.RegisterResponse{UserId: user.UserId, Username: user.Name}, nil
}

func (u *UserHandler) getUserInfoHandler(ctx *gin.Context) (*dto.User, error) {
	user, err := u.middleware.CurrentUser(ctx)
	if exception.CheckSaException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetUserInfo
	}

	return dto.UserEntityToDto(user), nil
}

func (u *UserHandler) updateUserHandler(ctx *gin.Context, req *dto.UpdateUserRequest) error {
	userDomain, err := u.middleware.CurrentUser(ctx)
	if exception.CheckSaException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrGetUserInfo
	}

	err = u.userApp.UpdateUser(ctx, userDomain.UserId, req)
	if exception.CheckSaException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrUpdateUserInfo
	}

	return nil
}

func (u *UserHandler) uploadAvatarHandler(ctx *gin.Context) (*dto.UploadAvatarResponse, error) {
	userId, err := u.middleware.CurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	fh, err := ctx.FormFile("avatar")
	if err != nil {
		return nil, exception.ErrUploadAvatar
	}

	avatarUrl, err := u.userApp.UploadAvatar(ctx, userId, fh)
	if exception.CheckSaException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrUploadAvatar
	}

	return &dto.UploadAvatarResponse{AvatarUrl: avatarUrl}, nil
}
