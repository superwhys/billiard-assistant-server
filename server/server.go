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
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"time"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/predis"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/gorilla/websocket"
	"github.com/superwhys/billiard-assistant-server/domain/game"
	"github.com/superwhys/billiard-assistant-server/domain/notice"
	"github.com/superwhys/billiard-assistant-server/domain/room"
	"github.com/superwhys/billiard-assistant-server/domain/session"
	"github.com/superwhys/billiard-assistant-server/domain/user"
	"github.com/superwhys/billiard-assistant-server/models"
	"github.com/superwhys/billiard-assistant-server/pkg/events"
	"github.com/superwhys/billiard-assistant-server/pkg/exception"
	"github.com/superwhys/billiard-assistant-server/pkg/oss/minio"
	"github.com/superwhys/billiard-assistant-server/pkg/wechat"
	"github.com/superwhys/billiard-assistant-server/server/dto"
	"gorm.io/gorm"
	
	gameDal "github.com/superwhys/billiard-assistant-server/pkg/dal/game"
	noticeDal "github.com/superwhys/billiard-assistant-server/pkg/dal/notice"
	roomDal "github.com/superwhys/billiard-assistant-server/pkg/dal/room"
	userDal "github.com/superwhys/billiard-assistant-server/pkg/dal/user"
	gameSrv "github.com/superwhys/billiard-assistant-server/server/game"
	noticeSrv "github.com/superwhys/billiard-assistant-server/server/notice"
	roomSrv "github.com/superwhys/billiard-assistant-server/server/room"
	userSrv "github.com/superwhys/billiard-assistant-server/server/user"
)

type BilliardServer struct {
	avatarDir   string
	redisClient *predis.RedisClient
	UserSrv     user.IUserService
	GameSrv     game.IGameService
	RoomSrv     room.IRoomService
	NoticeSrv   notice.INoticeService
	SessionSrv  session.ISessionService
	EventBus    *events.EventBus
}

func NewBilliardServer(conf *models.Config, db *gorm.DB, redis *predis.RedisClient, minioClient *minio.MinioOss) *BilliardServer {
	if !putils.FileExists(conf.AvatarDir) {
		err := os.MkdirAll(conf.AvatarDir, 0755)
		plog.PanicError(err, "createAvatarDir")
	}
	
	userRepo := userDal.NewUserRepo(db)
	gameRepo := gameDal.NewGameRepo(db)
	roomRepo := roomDal.NewRoomRepo(db)
	noticeRepo := noticeDal.NewNoticeRepo(db)
	
	s := &BilliardServer{
		avatarDir:   conf.AvatarDir,
		redisClient: redis,
		EventBus:    events.NewEventBus(),
		UserSrv:     userSrv.NewUserService(userRepo, minioClient),
		GameSrv:     gameSrv.NewGameService(gameRepo, userRepo),
		RoomSrv:     roomSrv.NewRoomService(roomRepo, redis),
		NoticeSrv:   noticeSrv.NewNoticeService(noticeRepo),
	}
	
	s.setupEventsSubscription()
	
	return s
}

func (s *BilliardServer) Login(ctx context.Context, username string, pwd string) (*dto.User, error) {
	u, err := s.UserSrv.Login(ctx, username, pwd)
	if err != nil {
		plog.Errorc(ctx, "login error: %v", err)
		return nil, err
	}
	
	return dto.UserEntityToDto(u), nil
}

func (s *BilliardServer) WechatLogin(ctx context.Context, code string) (*dto.User, *wechat.WechatSessionKeyResponse, error) {
	wxSessionKey, err := wechat.GetSessionKey(ctx, code)
	if err != nil {
		plog.Errorc(ctx, "get session key error: %v", err)
		return nil, nil, exception.ErrLoginFailed
	}
	
	u, err := s.UserSrv.WechatLogin(ctx, wxSessionKey)
	if err != nil {
		plog.Errorc(ctx, "wechat login error: %v", err)
		return nil, nil, err
	}
	
	return dto.UserEntityToDto(u), wxSessionKey, nil
}

func (s *BilliardServer) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
	u, err := s.UserSrv.Register(ctx, req.Username, req.Password)
	if err != nil {
		plog.Errorc(ctx, "create user error: %v", err)
		return nil, err
	}
	
	return dto.UserEntityToDto(u), nil
}

func (s *BilliardServer) UpdateUser(ctx context.Context, userId int, update *dto.UpdateUserRequest) (*dto.User, error) {
	u := &user.User{
		UserId: userId,
		Name:   update.Username,
		UserInfo: &user.BaseInfo{
			Avatar: update.AvatarUrl,
		},
	}
	
	user, err := s.UserSrv.UpdateUser(ctx, u)
	if err != nil {
		plog.Errorc(ctx, "update user info error: %v", err)
		return nil, err
	}
	
	return dto.UserEntityToDto(user), nil
}

func (s *BilliardServer) UploadAvatar(ctx context.Context, userId int, file *multipart.FileHeader) (string, error) {
	avatarUrl, err := s.UserSrv.UploadAvatar(ctx, userId, s.avatarDir, file)
	if err != nil {
		plog.Errorc(ctx, "upload avatar error: %v", err)
		return "", exception.ErrUploadAvatar
	}
	
	return avatarUrl, nil
}

func (s *BilliardServer) GetAvatar(ctx context.Context, avatarName string, dst io.Writer) error {
	err := s.UserSrv.GetUserAvatar(ctx, avatarName, dst)
	if err != nil {
		plog.Errorc(ctx, "get avatar error: %v", err)
		return exception.ErrGetAvatar
	}
	
	return nil
}

func (s *BilliardServer) GetNoticeList(ctx context.Context) ([]*dto.Notice, error) {
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

func (s *BilliardServer) GetGameList(ctx context.Context) ([]*dto.Game, error) {
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

func (s *BilliardServer) GetUserGameRooms(ctx context.Context, userId int) ([]*dto.GameRoom, error) {
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

func (s *BilliardServer) CreateGame(ctx context.Context, req *dto.CreateGameRequest) (*dto.Game, error) {
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

func (s *BilliardServer) DeleteGame(ctx context.Context, gameId int) error {
	err := s.GameSrv.DeleteGame(ctx, gameId)
	if err != nil {
		plog.Errorc(ctx, "delete game error: %v", err)
		return err
	}
	
	return err
}

func (s *BilliardServer) CreateRoom(ctx context.Context, userId, gameId int) (*dto.GameRoom, error) {
	gr, err := s.RoomSrv.CreateGameRoom(ctx, userId, gameId)
	if err != nil {
		plog.Errorc(ctx, "create game room error: %v", err)
		return nil, err
	}
	
	return dto.GameRoomEntityToDto(gr), nil
}

func (s *BilliardServer) UpdateGameRoomStatus(ctx context.Context, userId int, req *dto.UpdateGameRoomRequest) error {
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

func (s *BilliardServer) DeleteRoom(ctx context.Context, userId int, roomId int) error {
	err := s.RoomSrv.DeleteGameRoom(ctx, roomId, userId)
	if err != nil {
		plog.Errorc(ctx, "delete game room error: %v", err)
		return err
	}
	
	return nil
}

func (s *BilliardServer) EnterGameRoom(ctx context.Context, userId int, roomId int) error {
	enterUser, err := s.RoomSrv.EnterGameRoom(ctx, userId, roomId)
	if err != nil {
		plog.Errorc(ctx, "enter game room error: %v", err)
		return err
	}
	
	s.EventBus.Publish(room.NewEnterRoomEvent(roomId, enterUser))
	
	return nil
}

func (s *BilliardServer) LeaveGameRoom(ctx context.Context, userId int, roomId int) error {
	if _, err := s.RoomSrv.QuitGameRoom(ctx, userId, roomId); err != nil {
		plog.Errorc(ctx, "leave game room error: %v", err)
		return err
	}
	
	s.EventBus.Publish(room.NewLeaveRoomEvent(userId, roomId))
	
	return nil
}

func (s *BilliardServer) GetGameRoom(ctx context.Context, roomId int) (*dto.GameRoom, error) {
	r, err := s.RoomSrv.GetRoomById(ctx, roomId)
	if err != nil {
		plog.Errorc(ctx, "get game room error: %v", err)
		return nil, err
	}
	
	return dto.GameRoomEntityToDto(r), nil
}

func (s *BilliardServer) CreateRoomSession(ctx context.Context, userId, roomId int, conn *websocket.Conn) (*session.Session, error) {
	sess, err := s.SessionSrv.CreateSession(ctx, userId, roomId, conn)
	if err != nil {
		plog.Errorc(ctx, "register room session error: %v", err)
		return nil, err
	}
	
	s.SessionSrv.StartSession(sess)
	
	return sess, nil
}

func (s *BilliardServer) PrepareGame(ctx context.Context, userId, roomId int) error {
	err := s.RoomSrv.PrepareGame(ctx, userId, roomId)
	if err != nil {
		plog.Errorc(ctx, "prepare game error: %v", err)
		return err
	}
	
	s.EventBus.Publish(room.NewPrepareEvent(userId, roomId))
	return nil
}

func (s *BilliardServer) StartGame(ctx context.Context, userId, roomId int) error {
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

func (s *BilliardServer) BindPhone(ctx context.Context, userId int, phone, code string) error {
	if err := s.verifyPhoneCode(ctx, phone, code); err != nil {
		return err
	}
	
	u := &user.User{
		UserId: userId,
		UserInfo: &user.BaseInfo{
			Phone: phone,
		},
	}
	if _, err := s.UserSrv.UpdateUser(ctx, u); err != nil {
		return err
	}
	
	auth := &user.UserAuth{
		UserId:     userId,
		AuthType:   user.AuthTypePhone,
		Identifier: phone,
	}
	if err := s.UserSrv.CreateUserAuth(ctx, userId, auth); err != nil {
		plog.Errorc(ctx, "create phone auth error: %v", err)
		return err
	}
	
	return nil
}

func (s *BilliardServer) BindEmail(ctx context.Context, userId int, email, code string) error {
	if err := s.verifyEmailCode(ctx, email, code); err != nil {
		return err
	}
	
	u := &user.User{
		UserId: userId,
		UserInfo: &user.BaseInfo{
			Email: email,
		},
	}
	if _, err := s.UserSrv.UpdateUser(ctx, u); err != nil {
		return err
	}
	
	auth := &user.UserAuth{
		UserId:     userId,
		AuthType:   user.AuthTypeEmail,
		Identifier: email,
	}
	if err := s.UserSrv.CreateUserAuth(ctx, userId, auth); err != nil {
		plog.Errorc(ctx, "create email auth error: %v", err)
		return err
	}
	
	return nil
}

func (s *BilliardServer) SendPhoneCode(ctx context.Context, phone string) error {
	code := s.generateVerificationCode()
	expireAt := time.Now().Add(5 * time.Minute)
	
	key := fmt.Sprintf("phone_code:%s", phone)
	if err := s.redisClient.SetWithTTL(key, code, 5*time.Minute); err != nil {
		plog.Errorc(ctx, "save phone code error: %v", err)
		return exception.ErrSendPhoneCode
	}
	
	s.EventBus.Publish(user.NewSendPhoneCodeEvent(phone, code, expireAt))
	return nil
}

func (s *BilliardServer) SendEmailCode(ctx context.Context, email string) error {
	code := s.generateVerificationCode()
	expireAt := time.Now().Add(5 * time.Minute)
	
	key := fmt.Sprintf("email_code:%s", email)
	if err := s.redisClient.SetWithTTL(key, code, 5*time.Minute); err != nil {
		plog.Errorc(ctx, "save email code error: %v", err)
		return exception.ErrSendEmailCode
	}
	
	s.EventBus.Publish(user.NewSendEmailCodeEvent(email, code, expireAt))
	return nil
}

func (s *BilliardServer) verifyPhoneCode(ctx context.Context, phone, code string) error {
	// TODO: verify phone code
	panic("verify phone code")
}

func (s *BilliardServer) verifyEmailCode(ctx context.Context, email, code string) error {
	// TODO: verify email code
	panic("verify email code")
}

func (s *BilliardServer) generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
