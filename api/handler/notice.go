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
	
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/api/middlewares"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server/dto"
)

type NoticeHandlerApp interface {
	GetNoticeList(ctx context.Context) ([]*dto.Notice, error)
}

type NoticeHandler struct {
	noticeApp  NoticeHandlerApp
	middleware *middlewares.BilliardMiddleware
}

func NewNoticeHandler(server *server.BilliardServer, middleware *middlewares.BilliardMiddleware) *NoticeHandler {
	return &NoticeHandler{
		noticeApp:  server,
		middleware: middleware,
	}
}

func (g *NoticeHandler) Init(router gin.IRouter) {
	notic := router.Group("notice")
	notic.GET("list", pgin.ResponseHandler(g.getNoticeList))
}

func (g *NoticeHandler) getNoticeList(ctx *gin.Context) (*dto.GetNoticeListResp, error) {
	notices, err := g.noticeApp.GetNoticeList(ctx)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetNoticeList
	}
	
	return &dto.GetNoticeListResp{
		Notices: notices,
	}, nil
}
