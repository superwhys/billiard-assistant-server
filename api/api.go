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
	
	"github.com/go-puzzles/pgin"
	"github.com/go-puzzles/predis"
	"github.com/superwhys/snooker-assistant-server/api/handler"
	"github.com/superwhys/snooker-assistant-server/api/middlewares"
	"github.com/superwhys/snooker-assistant-server/pkg/token"
	"github.com/superwhys/snooker-assistant-server/server"
)

type SaApi struct {
	handler    http.Handler
	server     *server.SaServer
	middleware *middlewares.SaMiddleware
}

func SetupRouter(redisClient *predis.RedisClient, saServer *server.SaServer) *SaApi {
	tokenManager := token.NewManager(redisClient, token.WithCachePrefix("Sa"))
	saMiddleware := middlewares.NewSaMiddleware(tokenManager, saServer)
	
	router := pgin.NewServerHandlerWithOptions(
		pgin.WithMiddlewares(saMiddleware.UserLoginStatMiddleware()),
		pgin.WithRouters("/v1",
			handler.NewUserHandler(saServer, saMiddleware),
			handler.NewGameHandler(saServer, saMiddleware),
			handler.NewRoomHandler(saServer, saMiddleware),
			handler.NewNoticeHandler(saServer, saMiddleware),
		),
	)
	
	return &SaApi{
		handler:    router,
		server:     saServer,
		middleware: saMiddleware,
	}
}

func (a *SaApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
