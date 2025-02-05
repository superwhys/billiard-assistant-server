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

	"gitea.hoven.com/billiard/billiard-assistant-server/api/handler"
	"gitea.hoven.com/billiard/billiard-assistant-server/api/middlewares"
	"gitea.hoven.com/billiard/billiard-assistant-server/models"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/oss/minio"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/token"
	"gitea.hoven.com/billiard/billiard-assistant-server/server"
	"github.com/go-puzzles/puzzles/goredis"
	"github.com/go-puzzles/puzzles/pgin"
)

type BilliardApi struct {
	handler http.Handler
}

func SetupRouter(
	srvConf *models.Config,
	redisClient *goredis.PuzzleRedisClient,
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
