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
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/predis"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/auth"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/notice"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/session"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitlab.hoven.com/billiard/billiard-assistant-server/models"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/email"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/events"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/oss/minio"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/password"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/wechat"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server/dto"
	"gorm.io/gorm"

	authDal "gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/auth"
	gameDal "gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/game"
	noticeDal "gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/notice"
	roomDal "gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/room"
	userDal "gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/user"
	authSrv "gitlab.hoven.com/billiard/billiard-assistant-server/server/auth"
	gameSrv "gitlab.hoven.com/billiard/billiard-assistant-server/server/game"
	noticeSrv "gitlab.hoven.com/billiard/billiard-assistant-server/server/notice"
	roomSrv "gitlab.hoven.com/billiard/billiard-assistant-server/server/room"
	userSrv "gitlab.hoven.com/billiard/billiard-assistant-server/server/user"
)

type BilliardServer struct {
	redisClient *predis.RedisClient
	UserSrv     user.IUserService
	AuthSrv     auth.IAuthService
	GameSrv     game.IGameService
	RoomSrv     room.IRoomService
	NoticeSrv   notice.INoticeService
	SessionSrv  session.ISessionService
	EventBus    *events.EventBus
	emailSender email.EmailSender
}

func NewBilliardServer(
	conf *models.Config,
	db *gorm.DB,
	redis *predis.RedisClient,
	minioClient *minio.MinioOss,
	emailSender email.EmailSender,
) *BilliardServer {
	userRepo := userDal.NewUserRepo(db)
	authRepo := authDal.NewAuthRepo(db)
	gameRepo := gameDal.NewGameRepo(db)
	roomRepo := roomDal.NewRoomRepo(db)
	noticeRepo := noticeDal.NewNoticeRepo(db)

	s := &BilliardServer{
		redisClient: redis,
		EventBus:    events.NewEventBus(),
		UserSrv:     userSrv.NewUserService(userRepo, authRepo, minioClient),
		AuthSrv:     authSrv.NewAuthService(authRepo),
		GameSrv:     gameSrv.NewGameService(gameRepo, minioClient),
		RoomSrv:     roomSrv.NewRoomService(roomRepo, redis, conf.RoomConfig),
		NoticeSrv:   noticeSrv.NewNoticeService(noticeRepo),
		emailSender: emailSender,
	}

	s.setupEventsSubscription()

	return s
}

func (s *BilliardServer) Login(ctx context.Context, username string, pwd string) (*dto.User, error) {
	auth, err := s.AuthSrv.GetUserAuthByIdentifier(ctx, auth.AuthTypePassword, username)
	if err != nil && errors.Is(err, exception.ErrUserAuthNotFound) {
		return nil, exception.ErrUserNotFound
	} else if err != nil {
		plog.Errorc(ctx, "get user auth by identifier error: %v", err)
		return nil, exception.ErrLoginFailed
	}

	if !password.CheckPasswordHash(pwd, auth.Credential) {
		plog.Errorc(ctx, "check password error: %v", err)
		return nil, exception.ErrPasswordNotCorrect
	}

	user, err := s.UserSrv.GetUserById(ctx, auth.UserId)
	if err != nil {
		plog.Errorc(ctx, "get user by id error: %v", err)
		return nil, exception.ErrLoginFailed
	}

	return dto.UserEntityToDto(user), nil
}

func (s *BilliardServer) WechatLogin(ctx context.Context, code string) (*dto.User, *wechat.WechatSessionKeyResponse, error) {
	wxSessionKey, err := wechat.GetSessionKey(ctx, code)
	if err != nil {
		plog.Errorc(ctx, "get session key error: %v", err)
		return nil, nil, exception.ErrLoginFailed
	}

	ua, err := s.AuthSrv.GetUserAuthByIdentifier(ctx, auth.AuthTypeWechat, wxSessionKey.OpenID)
	if err != nil && !errors.Is(err, exception.ErrUserAuthNotFound) {
		plog.Errorc(ctx, "get user auth by open id error: %v", err)
		return nil, nil, exception.ErrLoginFailed
	}

	var loginUser *user.User
	// first login, create a new user and userAuth
	if ua == nil {
		newUser := &user.User{
			Name:   fmt.Sprintf("微信用户%s", wxSessionKey.OpenID),
			Status: user.StatusActive,
		}
		loginUser, err = s.UserSrv.CreateUser(ctx, newUser)
		if err != nil {
			plog.Errorc(ctx, "create user error: %v", err)
			return nil, nil, exception.ErrLoginFailed
		}

		ua := &auth.Auth{
			AuthType:   auth.AuthTypeWechat,
			Identifier: wxSessionKey.OpenID,
			Credential: wxSessionKey.SessionKey,
		}

		if err = s.AuthSrv.CreateUserAuth(ctx, loginUser.UserId, ua); err != nil {
			plog.Errorc(ctx, "create user auth error: %v", err)
			return nil, nil, exception.ErrLoginFailed
		}
	} else {
		loginUser, err = s.UserSrv.GetUserById(ctx, ua.UserId)
		if err != nil {
			plog.Errorc(ctx, "get user by id error: %v", err)
			return nil, nil, exception.ErrLoginFailed
		}
	}

	return dto.UserEntityToDto(loginUser), wxSessionKey, nil
}

func (s *BilliardServer) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.User, error) {
	// check whether a same userAuth is exists
	a, err := s.AuthSrv.GetUserAuthByIdentifier(ctx, auth.AuthTypePassword, req.Username)
	if err != nil && !errors.Is(err, exception.ErrUserAuthNotFound) {
		plog.Errorc(ctx, "get user auth by identifier error: %v", err)
		return nil, err
	}

	if a != nil {
		return nil, exception.ErrUserAlreadyExists
	}

	hashPwd, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "hashPassword")
	}

	newUser := &user.User{
		Name:   req.Username,
		Status: user.StatusActive,
	}

	newAuth := &auth.Auth{
		AuthType:   auth.AuthTypePassword,
		Identifier: req.Username,
		Credential: hashPwd,
	}

	u, err := s.UserSrv.CreateUser(ctx, newUser)
	if err != nil {
		plog.Errorc(ctx, "create user error: %v", err)
		return nil, err
	}

	err = s.AuthSrv.CreateUserAuth(ctx, u.UserId, newAuth)
	if err != nil {
		plog.Errorc(ctx, "create user auth error: %v", err)
		return nil, err
	}

	return dto.UserEntityToDto(u), nil
}

func (s *BilliardServer) UpdateUserName(ctx context.Context, userId int, userName string) error {
	u := &user.User{
		UserId: userId,
		Name:   userName,
	}

	u, err := s.UserSrv.GetUserByName(ctx, userName)
	if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		plog.Errorc(ctx, "get user by name error: %v", err)
		return err
	}

	if u != nil && u.UserId == userId {
		return exception.ErrUpdateNameSameAsOld
	} else if u != nil {
		return exception.ErrUserNameAlreadyExists
	}

	updateUser := &user.User{
		UserId: userId,
		Name:   userName,
	}
	u, err = s.UserSrv.UpdateUser(ctx, updateUser)
	if errors.Is(err, exception.ErrUserNotFound) {
		return exception.ErrUserNotFound
	} else if err != nil {
		plog.Errorc(ctx, "update user name error: %v", err)
		return err
	}

	auth, err := s.AuthSrv.GetUserAuthByType(ctx, userId, auth.AuthTypePassword)
	if err != nil && !errors.Is(err, exception.ErrUserAuthNotFound) {
		plog.Errorc(ctx, "get user auth error: %v", err)
		return err
	}
	if auth == nil {
		return nil
	}

	auth.Identifier = userName
	if err := s.AuthSrv.UpdateUserAuth(ctx, auth); err != nil {
		plog.Errorc(ctx, "update user auth identifier error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) UpdateUserGender(ctx context.Context, userId int, gender int) error {
	u := &user.User{
		UserId: userId,
		Gender: user.Gender(gender),
	}

	u, err := s.UserSrv.UpdateUser(ctx, u)
	if errors.Is(err, exception.ErrUserNotFound) {
		return exception.ErrUserNotFound
	} else if err != nil {
		plog.Errorc(ctx, "update user gender error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) UploadAvatar(ctx context.Context, userId int, file *multipart.FileHeader) (string, error) {
	avatarUrl, err := s.UserSrv.UploadAvatar(ctx, userId, file)
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

func (s *BilliardServer) GetSystemNotice(ctx context.Context) ([]*dto.Notice, error) {
	notices, err := s.NoticeSrv.GetNoticeByType(ctx, notice.System)
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

func (s *BilliardServer) AddNotices(ctx context.Context, req *dto.AddNoticeRequest) error {
	notices := make([]*notice.Notice, 0, len(req.Contents))
	for _, content := range req.Contents {
		notices = append(notices, &notice.Notice{
			NoticeType: req.NoticeType,
			Message:    content,
		})
	}

	return s.NoticeSrv.AddNotices(ctx, notices)
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
	rs, err := s.RoomSrv.GetUserGameRooms(ctx, userId, false)
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
		GameType: shared.BilliardGameType(req.GameType),
		Icon:     req.IconUrl,
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

func (s *BilliardServer) UpdateGame(ctx context.Context, req *dto.UpdateGameRequest) error {
	g := &game.Game{
		GameId: req.GameId,
		Icon:   req.IconUrl,
		GameConfig: &game.Config{
			MaxPlayers: req.MaxPlayers,
			Desc:       req.Desc,
		},
	}
	err := s.GameSrv.UpdateGame(ctx, g)
	if err != nil {
		plog.Errorc(ctx, "update game error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) DeleteGame(ctx context.Context, gameId int) error {
	err := s.GameSrv.DeleteGame(ctx, gameId)
	if err != nil {
		plog.Errorc(ctx, "delete game error: %v", err)
		return err
	}

	return err
}

func (s *BilliardServer) UploadGameIcon(ctx context.Context, fh *multipart.FileHeader) (string, error) {
	iconUrl, err := s.GameSrv.UploadGameIcon(ctx, fh)
	if err != nil {
		plog.Errorc(ctx, "upload gameIcon error: %v", err)
		return "", exception.ErrUploadGameIcon
	}

	return iconUrl, nil
}

func (s *BilliardServer) CreateRoom(ctx context.Context, userId, gameId int) (*dto.GameRoom, error) {
	exists, err := s.UserSrv.UserExists(ctx, userId)
	if err != nil || !exists {
		plog.Errorc(ctx, "get user by id error: %v", err)
		return nil, exception.ErrUserNotFound
	}

	gr, err := s.RoomSrv.CreateGameRoom(ctx, userId, gameId)
	if err != nil {
		plog.Errorc(ctx, "create game room error: %v", err)
		return nil, err
	}

	err = s.RoomSrv.EnterGameRoom(ctx, "", userId, gr.RoomId)
	if err != nil {
		plog.Errorc(ctx, "enter game room error: %v", err)
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

func (s *BilliardServer) EnterGameRoom(ctx context.Context, virtualName string, userId int, roomId int) error {
	err := s.RoomSrv.EnterGameRoom(ctx, virtualName, userId, roomId)
	if err != nil {
		plog.Errorc(ctx, "enter game room error: %v", err)
		return err
	}

	// TODO: publish user enter room events
	// s.EventBus.Publish(room.NewEnterRoomEvent(roomId, enterUser))

	return nil
}

func (s *BilliardServer) LeaveGameRoom(ctx context.Context, virtualName string, userId int, roomId int) error {
	err := s.RoomSrv.QuitGameRoom(ctx, virtualName, userId, roomId)
	if err != nil {
		plog.Errorc(ctx, "leave game room error: %v", err)
		return err
	}

	// publish user leave room events
	s.EventBus.Publish(room.NewLeaveRoomEvent(virtualName, userId, roomId))

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

	auth := &auth.Auth{
		UserId:     userId,
		AuthType:   auth.AuthTypePhone,
		Identifier: phone,
	}
	if err := s.AuthSrv.CreateUserAuth(ctx, userId, auth); err != nil {
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

	auth := &auth.Auth{
		UserId:     userId,
		AuthType:   auth.AuthTypeEmail,
		Identifier: email,
	}
	if err := s.AuthSrv.CreateUserAuth(ctx, userId, auth); err != nil {
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
	key := fmt.Sprintf("email_code:%s", email)
	var cacheCode string
	if err := s.redisClient.Get(key, &cacheCode); err != nil {
		plog.Errorc(ctx, "get email code error: %v", err)
		return exception.ErrVerifyCode
	}

	if code != cacheCode {
		return exception.ErrVerifyCode
	}

	return nil
}

func (s *BilliardServer) generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (s *BilliardServer) checkAuthExists(ctx context.Context, authType auth.AuthType, identifier string) (bool, error) {
	_, err := s.AuthSrv.GetUserAuthByIdentifier(ctx, authType, identifier)
	if err != nil {
		plog.Errorc(ctx, "get user auth by identifier(%v) error: %v", authType, err)
		return false, err
	}

	return true, err
}

func (s *BilliardServer) CheckPhoneBind(ctx context.Context, phone string) (bool, error) {
	return s.checkAuthExists(ctx, auth.AuthTypePhone, phone)
}

func (s *BilliardServer) CheckEmailBind(ctx context.Context, email string) (bool, error) {
	return s.checkAuthExists(ctx, auth.AuthTypeEmail, email)
}
