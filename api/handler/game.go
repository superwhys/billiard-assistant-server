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

type GameHandlerApp interface {
	GetGameList(ctx context.Context) ([]*dto.Game, error)
	CreateGame(ctx context.Context, req *dto.CreateGameRequest) (*dto.Game, error)
	DeleteGame(ctx context.Context, gameId int) error
}

type GameHandler struct {
	gameApp    GameHandlerApp
	middleware *middlewares.SaMiddleware
}

func NewGameHandler(server *server.SaServer, middleware *middlewares.SaMiddleware) *GameHandler {
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
	gameAdmin.DELETE("/:gameId", pgin.RequestWithErrorHandler(g.deleteGameHandler))
}

func (g *GameHandler) getGamesList(ctx *gin.Context) (*dto.GetGameListResp, error) {
	games, err := g.gameApp.GetGameList(ctx)
	if exception.CheckSaException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetGameList
	}

	return &dto.GetGameListResp{Games: games}, nil
}

func (g *GameHandler) createGame(ctx *gin.Context, req *dto.CreateGameRequest) (*dto.Game, error) {
	game, err := g.gameApp.CreateGame(ctx, req)
	if exception.CheckSaException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrCreateGame
	}

	return game, nil
}

func (g *GameHandler) deleteGameHandler(ctx *gin.Context, req *dto.DeleteGameRequest) error {
	err := g.gameApp.DeleteGame(ctx, req.GameId)
	if exception.CheckSaException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrDeleteGame
	}

	return nil
}
