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
	"gitlab.hoven.com/billiard/billiard-assistant-server/api/handler"
	"gitlab.hoven.com/billiard/billiard-assistant-server/api/middlewares"
	"gitlab.hoven.com/billiard/billiard-assistant-server/models"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/oss/minio"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/token"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server"
)

type BilliardApi struct {
	handler http.Handler
}

func SetupRouter(
	srvConf *models.Config,
	redisClient *predis.RedisClient,
	minioClient *minio.MinioOss,
	server *server.BilliardServer,
) *BilliardApi {
	tokenManager := token.NewManager(
		redisClient,
		token.WithCacheTTL(srvConf.TokenTtl),
		token.WithCachePrefix(srvConf.TokenPrefix),
	)
	middleware := middlewares.NewBilliardMiddleware(tokenManager, server)

	router := pgin.NewServerHandlerWithOptions(
		pgin.WithMiddlewares(middleware.UserLoginStatMiddleware()),
		pgin.WithRouters("/v1",
			handler.NewUserHandler(server, middleware),
			handler.NewAuthHandler(server, middleware),
			handler.NewGameHandler(server, middleware),
			handler.NewRoomHandler(server, middleware),
			handler.NewNoticeHandler(server, middleware),
			handler.NewRecordHandler(server, middleware),
			minioClient,
		),
	)

	return &BilliardApi{
		handler: router,
	}
}

func (a *BilliardApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
