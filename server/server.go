// File:		server.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/auth"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/game/nineball"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/game/snooker"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/notice"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/record"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/session"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitea.hoven.com/billiard/billiard-assistant-server/models"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/events"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/oss/minio"
	"gitea.hoven.com/core/auth-core/api/grpc/proto/authenticationpb"
	"gitea.hoven.com/core/auth-core/api/grpc/proto/userpb"
	"gitea.hoven.com/core/auth-core/api/grpc/proto/verifycodepb"
	"github.com/go-puzzles/puzzles/goredis"
	"gorm.io/gorm"

	gameDal "gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/game"
	noticeDal "gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/notice"
	recordDal "gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/record"
	roomDal "gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/room"
	userDal "gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/user"
	authSrv "gitea.hoven.com/billiard/billiard-assistant-server/server/auth"
	gameSrv "gitea.hoven.com/billiard/billiard-assistant-server/server/game"
	noticeSrv "gitea.hoven.com/billiard/billiard-assistant-server/server/notice"
	recordSrv "gitea.hoven.com/billiard/billiard-assistant-server/server/record"
	roomSrv "gitea.hoven.com/billiard/billiard-assistant-server/server/room"
	sessionSrv "gitea.hoven.com/billiard/billiard-assistant-server/server/session"
	userSrv "gitea.hoven.com/billiard/billiard-assistant-server/server/user"
)

type BilliardServer struct {
	redisClient *goredis.PuzzleRedisClient
	UserSrv     user.IUserService
	AuthSrv     auth.IAuthService
	GameSrv     game.IGameService
	RoomSrv     room.IRoomService
	NoticeSrv   notice.INoticeService
	SessionSrv  session.ISessionService
	RecordSrv   record.IRecordService
	EventBus    *events.EventBus
}

func NewBilliardServer(
	conf *models.Config,
	db *gorm.DB,
	redis *goredis.PuzzleRedisClient,
	minioClient *minio.MinioOss,
	authenticationClient authenticationpb.AuthCoreAuthenticationHandlerClient,
	userClient userpb.AuthCoreUserHandlerClient,
	verifycodeClient verifycodepb.AuthCoreVerifyCodeHandlerClient,
) *BilliardServer {
	userRepo := userDal.NewUserRepo(db)
	gameRepo := gameDal.NewGameRepo(db)
	roomRepo := roomDal.NewRoomRepo(db)
	noticeRepo := noticeDal.NewNoticeRepo(db)
	recordRepo := recordDal.NewRecordRepo(db)

	recordService := recordSrv.NewRecordService(recordRepo, roomRepo,
		recordSrv.WithGameStrategy(gameSrv.NewNineballService(redis)),
		recordSrv.WithGameStrategy(gameSrv.NewSnookerService(redis)),
		recordSrv.WithGameRecordTmp(shared.NineBall, &nineball.PlayerRecord{}),
		recordSrv.WithGameRecordTmp(shared.Snooker, &snooker.SnookerPlayerRecord{}),
	)

	s := &BilliardServer{
		redisClient: redis,
		EventBus:    events.NewEventBus(),
		UserSrv:     userSrv.NewUserService(userRepo, userClient),
		AuthSrv:     authSrv.NewAuthService(authenticationClient, verifycodeClient),
		GameSrv:     gameSrv.NewGameService(gameRepo, minioClient),
		RoomSrv:     roomSrv.NewRoomService(roomRepo, redis, conf.RoomConfig),
		SessionSrv:  sessionSrv.NewSessionService(),
		NoticeSrv:   noticeSrv.NewNoticeService(noticeRepo),
		RecordSrv:   recordService,
	}

	s.setupEventsSubscription()

	return s
}
