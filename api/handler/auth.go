// File:		auth.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package handler

import (
	"context"

	"gitea.hoven.com/billiard/billiard-assistant-server/api/middlewares"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/auth"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitea.hoven.com/billiard/billiard-assistant-server/server"
	"gitea.hoven.com/billiard/billiard-assistant-server/server/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/pkg/errors"
)

type AuthHandlerApp interface {
	AccountRegister(ctx context.Context, username, password string) error
	AccountLogin(ctx context.Context, device, username, password string) (*auth.Token, error)
	WechatLogin(ctx context.Context, device, code string) (*auth.Token, error)
	Logout(ctx context.Context, token string) error
	BindEmail(ctx context.Context, token string, email, code string) error
	SendEmailCode(ctx context.Context, email string) error
}

type AuthHandler struct {
	authApp    AuthHandlerApp
	middleware *middlewares.BilliardMiddleware
}

func NewAuthHandler(server *server.BilliardServer, middleware *middlewares.BilliardMiddleware) *AuthHandler {
	return &AuthHandler{
		authApp:    server,
		middleware: middleware,
	}
}

func (a *AuthHandler) Init(router gin.IRouter) {
	auth := router.Group("auth")
	auth.POST("login/wx", pgin.RequestResponseHandler(a.wechatLoginHandler))
	auth.POST("login/account", pgin.RequestResponseHandler(a.loginAccountHandler))
	auth.POST("register/account", pgin.RequestWithErrorHandler(a.registerAccountHandler))

	authNeedLogin := router.Group("auth", a.middleware.UserLoginRequired())
	authNeedLogin.GET("logout", pgin.ErrorReturnHandler(a.logoutHandler))
	authNeedLogin.POST("bind/email", pgin.RequestWithErrorHandler(a.bindEmailHandler))
	authNeedLogin.POST("send/email_code", pgin.RequestWithErrorHandler(a.sendEmailCodeHandler))
}

func (a *AuthHandler) loginAccountHandler(ctx *gin.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	resp, err := a.authApp.AccountLogin(ctx, req.DeviceId, req.Username, req.Password)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrLoginFailed
	}

	token := middlewares.NewUserLoginToken(resp.UserId, resp.AccessToken, resp.RefreshToken)
	a.middleware.SaveToken(token, ctx)

	return &dto.LoginResponse{
		UserId: resp.UserId,
		Token:  token.GetKey(),
	}, nil
}

func (a *AuthHandler) registerAccountHandler(ctx *gin.Context, req *dto.RegisterRequest) error {
	err := a.authApp.AccountRegister(ctx, req.Username, req.Password)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrRegisterUser
	}

	return nil
}

func (a *AuthHandler) wechatLoginHandler(ctx *gin.Context, req *dto.WechatLoginRequest) (*dto.LoginResponse, error) {
	resp, err := a.authApp.WechatLogin(ctx, req.DeviceId, req.Code)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrLoginFailed
	}

	token := middlewares.NewUserLoginToken(resp.UserId, resp.AccessToken, resp.RefreshToken)
	a.middleware.SaveToken(token, ctx)

	return &dto.LoginResponse{
		UserId: resp.UserId,
		Token:  token.GetKey(),
	}, nil
}

func (a *AuthHandler) logoutHandler(ctx *gin.Context) error {
	token := a.middleware.GetLoginToken(ctx).GetKey()

	err := a.authApp.Logout(ctx, token)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrLogoutFailed
	}

	return nil
}

func (a *AuthHandler) bindEmailHandler(ctx *gin.Context, req *dto.BindEmailRequest) error {
	token := a.middleware.GetLoginToken(ctx).GetKey()

	err := a.authApp.BindEmail(ctx, token, req.Email, req.Code)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrBindEmail
	}

	return nil
}

func (a *AuthHandler) sendEmailCodeHandler(ctx *gin.Context, req *dto.SendEmailCodeRequest) error {
	err := a.authApp.SendEmailCode(ctx, req.Email)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrSendEmailCode
	}
	return nil
}
