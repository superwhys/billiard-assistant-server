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
	"github.com/superwhys/snooker-assistant-server/api/middlewares"
	"github.com/superwhys/snooker-assistant-server/pkg/exception"
	"github.com/superwhys/snooker-assistant-server/server"
	"github.com/superwhys/snooker-assistant-server/server/dto"
)

type NoticeHandlerApp interface {
	GetNoticeList(ctx context.Context) ([]*dto.Notice, error)
}

type NoticeHandler struct {
	noticeApp  NoticeHandlerApp
	middleware *middlewares.SaMiddleware
}

func NewNoticeHandler(server *server.SaServer, middleware *middlewares.SaMiddleware) *NoticeHandler {
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
	if exception.CheckSaException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetNoticeList
	}

	return &dto.GetNoticeListResp{
		Notices: notices,
	}, nil
}
