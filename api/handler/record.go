// File:		record.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package handler

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"gitlab.hoven.com/billiard/billiard-assistant-server/api/middlewares"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server/dto"
)

type RecordApp interface {
	HandleRoomAction(ctx context.Context, roomId, userId int, action json.RawMessage) error
	HandleRoomRecord(ctx context.Context, roomId int, record json.RawMessage) error
	GetRoomActions(ctx context.Context, roomId int) (*dto.Action, error)
	GetRoomRecoed(ctx context.Context, roomId int) (*dto.Record, error)
}

type RecordHandler struct {
	recordApp  RecordApp
	middleware *middlewares.BilliardMiddleware
}

func NewRecordHandler(server *server.BilliardServer, middleware *middlewares.BilliardMiddleware) *RecordHandler {
	return &RecordHandler{
		recordApp:  server,
		middleware: middleware,
	}
}

func (r *RecordHandler) Init(router gin.IRouter) {
	recordRouter := router.Group("record", r.middleware.UserLoginRequired())
	recordRouter.POST("action/:roomId", pgin.RequestWithErrorHandler(r.roomActionHandler))
	recordRouter.POST("record/:roomId", pgin.RequestWithErrorHandler(r.roomRecordHandler))
	recordRouter.GET("record/:roomId", pgin.RequestResponseHandler(r.getRoomRecordHandler))
	recordRouter.GET("action/:roomId", pgin.RequestResponseHandler(r.getRoomActionHandler))
}

func (r *RecordHandler) roomActionHandler(ctx *gin.Context, req *dto.RoomActionRequest) error {
	userId, err := r.middleware.CurrentUserId(ctx)
	if err != nil {
		return err
	}

	err = r.recordApp.HandleRoomAction(ctx, req.RoomId, userId, req.Action)
	if exception.CheckException(err) {
		return err
	} else if err != nil {
		return exception.ErrHandleAction
	}
	return nil
}

func (r *RecordHandler) roomRecordHandler(ctx *gin.Context, req *dto.RoomRecordRequest) error {
	plog.Debugc(ctx, "%v", string(req.Records))
	err := r.recordApp.HandleRoomRecord(ctx, req.RoomId, req.Records)
	if exception.CheckException(err) {
		return err
	} else if err != nil {
		return exception.ErrHandleRecord
	}
	return nil
}

func (r *RecordHandler) getRoomRecordHandler(ctx *gin.Context, req *dto.RoomUriRequest) (*dto.Record, error) {
	record, err := r.recordApp.GetRoomRecoed(ctx, req.RoomId)
	if exception.CheckException(err) {
		return nil, err
	} else if err != nil {
		return nil, exception.ErrGetRoomRecord
	}
	return record, nil
}

func (r *RecordHandler) getRoomActionHandler(ctx *gin.Context, req *dto.RoomUriRequest) (*dto.Action, error) {
	action, err := r.recordApp.GetRoomActions(ctx, req.RoomId)
	if exception.CheckException(err) {
		return nil, err
	} else if err != nil {
		return nil, exception.ErrGetRoomAction
	}
	return action, nil
}
