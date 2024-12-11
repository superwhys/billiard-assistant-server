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

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/pkg/errors"
	"gitea.hoven.com/billiard/billiard-assistant-server/api/middlewares"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitea.hoven.com/billiard/billiard-assistant-server/server"
	"gitea.hoven.com/billiard/billiard-assistant-server/server/dto"
)

type AuthHandlerApp interface {
	BindPhone(ctx context.Context, userId int, phone, code string) error
	BindEmail(ctx context.Context, userId int, email, code string) error
	SendPhoneCode(ctx context.Context, phone string) error
	SendEmailCode(ctx context.Context, email string) error
	CheckPhoneBind(ctx context.Context, phone string) (bool, error)
	CheckEmailBind(ctx context.Context, email string) (bool, error)
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
	auth := router.Group("auth", a.middleware.UserLoginRequired())
	auth.POST("bind/phone", pgin.RequestWithErrorHandler(a.bindPhoneHandler))
	auth.POST("bind/email", pgin.RequestWithErrorHandler(a.bindEmailHandler))
	auth.POST("check/phone", pgin.RequestResponseHandler(a.checkPhoneBindHandler))
	auth.POST("check/email", pgin.RequestResponseHandler(a.checkEmailBindHandler))
	auth.POST("send/phone_code", pgin.RequestWithErrorHandler(a.sendPhoneCodeHandler))
	auth.POST("send/email_code", pgin.RequestWithErrorHandler(a.sendEmailCodeHandler))
}

func (a *AuthHandler) bindPhoneHandler(ctx *gin.Context, req *dto.BindPhoneRequest) error {
	userId, err := a.middleware.CurrentUserId(ctx)
	if err != nil {
		return err
	}

	err = a.authApp.BindPhone(ctx, userId, req.Phone, req.Code)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrBindPhone
	}

	return nil
}

func (a *AuthHandler) bindEmailHandler(ctx *gin.Context, req *dto.BindEmailRequest) error {
	userId, err := a.middleware.CurrentUserId(ctx)
	if err != nil {
		return err
	}

	err = a.authApp.BindEmail(ctx, userId, req.Email, req.Code)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrBindEmail
	}

	return nil
}

func (a *AuthHandler) sendPhoneCodeHandler(ctx *gin.Context, req *dto.SendPhoneCodeRequest) error {
	err := a.authApp.SendPhoneCode(ctx, req.Phone)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrSendPhoneCode
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

func (a *AuthHandler) checkPhoneBindHandler(ctx *gin.Context, req *dto.CheckPhoneBindRequest) (*dto.CheckPhoneBindResponse, error) {
	isBound, err := a.authApp.CheckPhoneBind(ctx, req.Phone)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetUserInfo
	}

	return &dto.CheckPhoneBindResponse{IsBound: isBound}, nil
}

func (a *AuthHandler) checkEmailBindHandler(ctx *gin.Context, req *dto.CheckPhoneBindRequest) (*dto.CheckPhoneBindResponse, error) {
	isBound, err := a.authApp.CheckEmailBind(ctx, req.Phone)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetUserInfo
	}

	return &dto.CheckPhoneBindResponse{IsBound: isBound}, nil
}
