// File:		api.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package api

import (
	"net/http"
	
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/go-puzzles/puzzles/predis"
	"github.com/superwhys/billiard-assistant-server/api/handler"
	"github.com/superwhys/billiard-assistant-server/api/middlewares"
	"github.com/superwhys/billiard-assistant-server/pkg/token"
	"github.com/superwhys/billiard-assistant-server/server"
)

type BilliardApi struct {
	handler http.Handler
}

func SetupRouter(redisClient *predis.RedisClient, server *server.BilliardServer) *BilliardApi {
	tokenManager := token.NewManager(redisClient, token.WithCachePrefix("billiard"))
	middleware := middlewares.NewBilliardMiddleware(tokenManager, server)
	
	router := pgin.NewServerHandlerWithOptions(
		pgin.WithMiddlewares(middleware.UserLoginStatMiddleware()),
		pgin.WithRouters("/v1",
			handler.NewUserHandler(server, middleware),
			handler.NewGameHandler(server, middleware),
			handler.NewRoomHandler(server, middleware),
			handler.NewNoticeHandler(server, middleware),
		),
	)
	
	return &BilliardApi{
		handler: router,
	}
}

func (a *BilliardApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
