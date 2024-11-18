// File:		server.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/predis"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/gorilla/websocket"
	"github.com/superwhys/snooker-assistant-server/domain/game"
	"github.com/superwhys/snooker-assistant-server/domain/notice"
	"github.com/superwhys/snooker-assistant-server/domain/room"
	"github.com/superwhys/snooker-assistant-server/domain/session"
	"github.com/superwhys/snooker-assistant-server/domain/user"
	"github.com/superwhys/snooker-assistant-server/models"
	"github.com/superwhys/snooker-assistant-server/pkg/events"
	"github.com/superwhys/snooker-assistant-server/pkg/exception"
	"github.com/superwhys/snooker-assistant-server/pkg/oss/minio"
	"github.com/superwhys/snooker-assistant-server/pkg/wechat"
	"github.com/superwhys/snooker-assistant-server/server/dto"
	"gorm.io/gorm"

	gameDal "github.com/superwhys/snooker-assistant-server/pkg/dal/game"
	noticeDal "github.com/superwhys/snooker-assistant-server/pkg/dal/notice"
	roomDal "github.com/superwhys/snooker-assistant-server/pkg/dal/room"
	userDal "github.com/superwhys/snooker-assistant-server/pkg/dal/user"
	gameSrv "github.com/superwhys/snooker-assistant-server/server/game"
	noticeSrv "github.com/superwhys/snooker-assistant-server/server/notice"
	roomSrv "github.com/superwhys/snooker-assistant-server/server/room"
	userSrv "github.com/superwhys/snooker-assistant-server/server/user"
)

type SaServer struct {
	avatarDir  string
	UserSrv    user.IUserService
	GameSrv    game.IGameService
	RoomSrv    room.IRoomService
	NoticeSrv  notice.INoticeService
	SessionSrv session.ISessionService
	EventBus   *events.EventBus
}

func NewSaServer(conf *models.SaConfig, db *gorm.DB, redis *predis.RedisClient, minioClient *minio.MinioOss) *SaServer {
	if !putils.FileExists(conf.AvatarDir) {
		err := os.MkdirAll(conf.AvatarDir, 0755)
		plog.PanicError(err, "createAvatarDir")
	}

	userRepo := userDal.NewUserRepo(db)
	gameRepo := gameDal.NewGameRepo(db)
	roomRepo := roomDal.NewRoomRepo(db)
	noticeRepo := noticeDal.NewNoticeRepo(db)

	s := &SaServer{
		avatarDir: conf.AvatarDir,
		EventBus:  events.NewEventBus(),
		UserSrv:   userSrv.NewUserService(userRepo, minioClient),
		GameSrv:   gameSrv.NewGameService(gameRepo, userRepo),
		RoomSrv:   roomSrv.NewRoomService(roomRepo, redis),
		NoticeSrv: noticeSrv.NewNoticeService(noticeRepo),
	}

	s.setupEventsSubscription()

	return s
}

func (s *SaServer) Login(ctx context.Context, username string, pwd string) (*dto.User, error) {
	u, err := s.UserSrv.Login(ctx, username, pwd)
	if err != nil {
		plog.Errorf("login error: %v", err)
		return nil, err
	}

	return dto.UserEntityToDto(u), nil
}

func (s *SaServer) WechatLogin(ctx context.Context, code string) (*dto.User, *wechat.WechatSessionKeyResponse, error) {
	wxSessionKey, err := wechat.GetSessionKey(ctx, code)
	if err != nil {
		plog.Errorf("get session key error: %v", err)
		return nil, nil, exception.ErrLoginFailed
	}

	u, err := s.UserSrv.WechatLogin(ctx, wxSessionKey)
	if err != nil {
		plog.Errorf("wechat login error: %v", err)
		return nil, nil, err
	}

	return dto.UserEntityToDto(u), wxSessionKey, nil
}

func (s *SaServer) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
	u := &user.User{
		Name:     req.Username,
		Password: req.Password,
		WechatId: req.WechatId,
		UserInfo: &user.BaseInfo{
			Phone: req.Phone,
			Email: req.Email,
		},
		Role: user.RoleUser,
	}
	err := s.UserSrv.Register(ctx, u)
	if err != nil {
		plog.Errorc(ctx, "create user error: %v", err)
		return nil, err
	}

	return dto.UserEntityToDto(u), nil
}

func (s *SaServer) UpdateUser(ctx context.Context, userId int, update *dto.UpdateUserRequest) error {
	u := &user.User{
		UserId:   userId,
		Name:     update.Username,
		Password: update.Password,
		UserInfo: &user.BaseInfo{
			Avatar: update.AvatarUrl,
			Phone:  update.Phone,
			Email:  update.Email,
		},
	}

	err := s.UserSrv.UpdateUser(ctx, u)
	if err != nil {
		plog.Errorc(ctx, "update user info error: %v", err)
		return err
	}

	return nil
}

func (s *SaServer) UploadAvatar(ctx context.Context, userId int, file *multipart.FileHeader) (string, error) {
	avatarUrl, err := s.UserSrv.UploadAvatar(ctx, userId, s.avatarDir, file)
	if err != nil {
		plog.Errorc(ctx, "upload avatar error: %v", err)
		return "", exception.ErrUploadAvatar
	}

	return avatarUrl, nil
}

func (s *SaServer) GetAvatar(ctx context.Context, avatarName string, dst io.Writer) error {
	err := s.UserSrv.GetUserAvatar(ctx, avatarName, dst)
	if err != nil {
		plog.Errorc(ctx, "get avatar error: %v", err)
		return exception.ErrGetAvatar
	}

	return nil
}

func (s *SaServer) GetNoticeList(ctx context.Context) ([]*dto.Notice, error) {
	notices, err := s.NoticeSrv.GetNoticeList(ctx)
	if err != nil {
		plog.Errorc(ctx, "get notice list error: %v", err)
		return nil, err
	}

	ret := make([]*dto.Notice, 0, len(notices))
	for _, n := range notices {
		ret = append(ret, dto.NoticeEntityToDto(n))
	}

	return ret, nil
}

func (s *SaServer) GetGameList(ctx context.Context) ([]*dto.Game, error) {
	gameList, err := s.GameSrv.GetGameList(ctx)
	if err != nil {
		plog.Errorc(ctx, "get game list error: %v", err)
		return nil, err
	}

	ret := make([]*dto.Game, 0, len(gameList))
	for _, g := range gameList {
		ret = append(ret, dto.GameEntityToDto(g))
	}

	return ret, nil
}

func (s *SaServer) GetUserGameRooms(ctx context.Context, userId int) ([]*dto.GameRoom, error) {
	rs, err := s.RoomSrv.GetUserGameRooms(ctx, userId)
	if err != nil {
		plog.Errorc(ctx, "get user game rooms error: %v", err)
		return nil, err
	}

	ret := make([]*dto.GameRoom, 0, len(rs))
	for _, r := range rs {
		ret = append(ret, dto.GameRoomEntityToDto(r))
	}

	return ret, nil
}

func (s *SaServer) CreateGame(ctx context.Context, req *dto.CreateGameRequest) (*dto.Game, error) {
	g := &game.Game{
		GameType: req.GameType,
		GameConfig: &game.Config{
			MaxPlayers: req.MaxPlayers,
			Desc:       req.Desc,
		},
	}
	err := s.GameSrv.CreateGame(ctx, g)
	if err != nil {
		plog.Errorc(ctx, "create game error: %v", err)
		return nil, err
	}

	return dto.GameEntityToDto(g), nil
}

func (s *SaServer) DeleteGame(ctx context.Context, gameId int) error {
	err := s.GameSrv.DeleteGame(ctx, gameId)
	if err != nil {
		plog.Errorc(ctx, "delete game error: %v", err)
		return err
	}

	return err
}

func (s *SaServer) CreateRoom(ctx context.Context, userId, gameId int) (*dto.GameRoom, error) {
	gr, err := s.RoomSrv.CreateGameRoom(ctx, userId, gameId)
	if err != nil {
		plog.Errorc(ctx, "create game room error: %v", err)
		return nil, err
	}

	return dto.GameRoomEntityToDto(gr), nil
}

func (s *SaServer) UpdateGameRoomStatus(ctx context.Context, userId int, req *dto.UpdateGameRoomRequest) error {
	gr := &room.Room{
		RoomId:        req.RoomId,
		OwnerId:       userId,
		GameStatus:    req.GameStatus,
		WinLoseStatus: req.WinLoseStatus,
	}

	err := s.RoomSrv.UpdateGameRoomStatus(ctx, gr)
	if err != nil {
		plog.Errorc(ctx, "update game room status error: %v", err)
		return err
	}

	return nil
}

func (s *SaServer) DeleteRoom(ctx context.Context, userId int, roomId int) error {
	err := s.RoomSrv.DeleteGameRoom(ctx, roomId, userId)
	if err != nil {
		plog.Errorc(ctx, "delete game room error: %v", err)
		return err
	}

	return nil
}

func (s *SaServer) EnterGameRoom(ctx context.Context, userId int, roomId int) error {
	enterUser, err := s.RoomSrv.EnterGameRoom(ctx, userId, roomId)
	if err != nil {
		plog.Errorc(ctx, "enter game room error: %v", err)
		return err
	}

	s.EventBus.Publish(room.NewEnterRoomEvent(roomId, enterUser))

	return nil
}

func (s *SaServer) LeaveGameRoom(ctx context.Context, userId int, roomId int) error {
	if _, err := s.RoomSrv.QuitGameRoom(ctx, userId, roomId); err != nil {
		plog.Errorc(ctx, "leave game room error: %v", err)
		return err
	}

	s.EventBus.Publish(room.NewLeaveRoomEvent(userId, roomId))

	return nil
}

func (s *SaServer) GetGameRoom(ctx context.Context, roomId int) (*dto.GameRoom, error) {
	r, err := s.RoomSrv.GetRoomById(ctx, roomId)
	if err != nil {
		plog.Errorc(ctx, "get game room error: %v", err)
		return nil, err
	}

	return dto.GameRoomEntityToDto(r), nil
}

func (s *SaServer) CreateRoomSession(ctx context.Context, userId, roomId int, conn *websocket.Conn) (*session.Session, error) {
	sess, err := s.SessionSrv.CreateSession(ctx, userId, roomId, conn)
	if err != nil {
		plog.Errorc(ctx, "register room session error: %v", err)
		return nil, err
	}

	s.SessionSrv.StartSession(sess)

	return sess, nil
}

func (s *SaServer) PrepareGame(ctx context.Context, userId, roomId int) error {
	err := s.RoomSrv.PrepareGame(ctx, userId, roomId)
	if err != nil {
		plog.Errorc(ctx, "prepare game error: %v", err)
		return err
	}

	s.EventBus.Publish(room.NewPrepareEvent(userId, roomId))
	return nil
}

func (s *SaServer) StartGame(ctx context.Context, userId, roomId int) error {
	currentGame, err := s.RoomSrv.StartGame(ctx, userId, roomId)
	if err != nil {
		plog.Errorc(ctx, "start game error: %v", err)
		return err
	}

	// game init
	gs, err := game.NewGameStrategy(currentGame.GetGameType())
	payload := gs.SetupGame(currentGame.GetGameConfig())

	s.EventBus.Publish(room.NewGameStartEvent(roomId, payload))
	return nil
}
