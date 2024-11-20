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
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/api/middlewares"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server/dto"
)

type GameHandlerApp interface {
	GetGameList(ctx context.Context) ([]*dto.Game, error)
	CreateGame(ctx context.Context, req *dto.CreateGameRequest) (*dto.Game, error)
	UpdateGame(ctx context.Context, req *dto.UpdateGameRequest) error
	DeleteGame(ctx context.Context, gameId int) error
	UploadGameIcon(ctx context.Context, fh *multipart.FileHeader) (string, error)
}

type GameHandler struct {
	gameApp    GameHandlerApp
	middleware *middlewares.BilliardMiddleware
}

func NewGameHandler(server *server.BilliardServer, middleware *middlewares.BilliardMiddleware) *GameHandler {
	return &GameHandler{
		gameApp:    server,
		middleware: middleware,
	}
}

func (g *GameHandler) Init(router gin.IRouter) {
	game := router.Group("game")
	game.GET("list", pgin.ResponseHandler(g.getGamesList))

	// game admin router
	gameAdmin := router.Group("game/admin", g.middleware.AdminRequired())
	gameAdmin.POST("create", pgin.RequestResponseHandler(g.createGame))
	gameAdmin.PUT("update/:gameId", pgin.RequestWithErrorHandler(g.updateGameHandler))
	gameAdmin.DELETE("/:gameId", pgin.RequestWithErrorHandler(g.deleteGameHandler))
	gameAdmin.POST("icon/upload", pgin.ResponseHandler(g.uploadGameIcon))
}

func (g *GameHandler) updateGameHandler(ctx *gin.Context, req *dto.UpdateGameRequest) error {
	err := g.gameApp.UpdateGame(ctx, req)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrUpdateGame
	}

	return nil
}

func (g *GameHandler) uploadGameIcon(ctx *gin.Context) (*dto.UploadGameIconResponse, error) {
	fh, err := ctx.FormFile("icon")
	if err != nil {
		plog.Errorc(ctx, "get fileHeader error: %v", err)
		return nil, exception.ErrUploadGameIcon
	}

	iconUrl, err := g.gameApp.UploadGameIcon(ctx, fh)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrUploadAvatar
	}

	return &dto.UploadGameIconResponse{IconUrl: iconUrl}, nil
}

func (g *GameHandler) getGamesList(ctx *gin.Context) (*dto.GetGameListResp, error) {
	games, err := g.gameApp.GetGameList(ctx)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetGameList
	}

	return &dto.GetGameListResp{Games: games}, nil
}

func (g *GameHandler) createGame(ctx *gin.Context, req *dto.CreateGameRequest) (*dto.Game, error) {
	game, err := g.gameApp.CreateGame(ctx, req)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrCreateGame
	}

	return game, nil
}

func (g *GameHandler) deleteGameHandler(ctx *gin.Context, req *dto.DeleteGameRequest) error {
	err := g.gameApp.DeleteGame(ctx, req.GameId)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrDeleteGame
	}

	return nil
}
