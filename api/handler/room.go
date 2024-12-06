package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/api/middlewares"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server/dto"
	"gorm.io/datatypes"
)

type RoomApp interface {
	GetUserGameRooms(ctx context.Context, userId int) ([]*dto.GameRoom, error)
	GetGameRoom(ctx context.Context, roomId int) (*dto.GameRoom, error)
	GetGameRoomByCode(ctx context.Context, roomCode string) (*dto.GameRoom, error)
	CreateRoom(ctx context.Context, userId, gameId int) (*dto.GameRoom, error)
	UpdateGameRoomStatus(ctx context.Context, userId int, gameRoom *dto.UpdateGameRoomRequest) error
	UpdateGameRoomExtra(ctx context.Context, userId int, gameRoom *dto.UpdateGameRoomExtraRequest) error
	DeleteRoom(ctx context.Context, userId, roomId int) error
	EnterGameRoom(ctx context.Context, roomId, currentUid int, virtualUser string) error
	LeaveGameRoom(ctx context.Context, roomId, currentUid int, virtualUser string) error
	CreateRoomSession(ctx context.Context, userId, roomId int, w http.ResponseWriter, r *http.Request) error
	StartGame(ctx context.Context, userId, roomId int, extra datatypes.JSONMap) error
	EndGame(ctx context.Context, userId, roomId int) error
}

type RoomHandler struct {
	roomApp    RoomApp
	middleware *middlewares.BilliardMiddleware
}

func NewRoomHandler(server *server.BilliardServer, middleware *middlewares.BilliardMiddleware) *RoomHandler {
	return &RoomHandler{
		roomApp:    server,
		middleware: middleware,
	}
}

func (r *RoomHandler) Init(router gin.IRouter) {
	roomAuth := router.Group("room", r.middleware.UserLoginRequired())
	roomAuth.POST("create", pgin.RequestResponseHandler(r.createGameRoom))
	roomAuth.GET("list", pgin.ResponseHandler(r.getUserGameRoom))
	roomAuth.GET("code/:roomCode", pgin.RequestResponseHandler(r.getRoomInfoByCode))

	roomIdAuth := roomAuth.Group(":roomId")
	roomIdAuth.GET("", pgin.RequestResponseHandler(r.getRoomInfo))
	roomIdAuth.PUT("update", pgin.RequestWithErrorHandler(r.updateGameRoomStatus))
	roomIdAuth.POST("update/extra", pgin.RequestWithErrorHandler(r.updateGameRoomExtra))
	roomIdAuth.POST("enter", pgin.RequestWithErrorHandler(r.enterGameRoom))
	roomIdAuth.POST("leave", pgin.RequestWithErrorHandler(r.leaveGameRoom))
	roomIdAuth.POST("start", pgin.RequestWithErrorHandler(r.startGame))
	roomIdAuth.POST("end", pgin.RequestWithErrorHandler(r.endGame))
	roomIdAuth.GET("ws", r.websocketHandler)
	roomIdAuth.DELETE("delete", pgin.RequestWithErrorHandler(r.deleteGameRoom))
}

func (r *RoomHandler) websocketHandler(ctx *gin.Context) {
	req := new(dto.UriRoomId)
	err := ctx.ShouldBindUri(req)
	if err != nil {
		plog.Errorc(ctx, "parse room id error: %v", err)
		return
	}

	userId, err := r.middleware.CurrentUserId(ctx)
	if err != nil {
		plog.Errorc(ctx, "get current user id error: %v", err)
		return
	}

	err = r.roomApp.CreateRoomSession(ctx, userId, req.RoomId, ctx.Writer, ctx.Request)
	if err != nil {
		plog.Errorc(ctx, "register room session error: %v", err)
		pgin.ErrorRet(400, err)
		return
	}
}

func (r *RoomHandler) getCurrentUser(ctx *gin.Context) (*user.User, error) {
	user, err := r.middleware.CurrentUser(ctx)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		plog.Errorc(ctx, "getCurrentUser error: %v", err)
		return nil, exception.ErrUserVerify
	}

	return user, nil
}

func (r *RoomHandler) getCurrentUserId(ctx *gin.Context) (int, error) {
	userId, err := r.middleware.CurrentUserId(ctx)
	if exception.CheckException(err) {
		return 0, errors.Cause(err)
	} else if err != nil {
		plog.Errorc(ctx, "getCurrentUserId error: %v", err)
		return 0, exception.ErrUserVerify
	}

	return userId, nil
}

func (r *RoomHandler) createGameRoom(ctx *gin.Context, req *dto.CreateGameRoomRequest) (*dto.GameRoom, error) {
	userId, err := r.getCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	gr, err := r.roomApp.CreateRoom(ctx, userId, req.GameId)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrCreateGameRoom
	}

	return gr, nil
}

func (r *RoomHandler) deleteGameRoom(ctx *gin.Context, req *dto.DeleteGameRoomRequest) error {
	userId, err := r.getCurrentUserId(ctx)
	if err != nil {
		return err
	}

	err = r.roomApp.DeleteRoom(ctx, userId, req.RoomId)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrDeleteGame
	}

	return nil
}

func (r *RoomHandler) getRoomInfo(ctx *gin.Context, req *dto.GetRoomRequest) (*dto.GameRoom, error) {
	gr, err := r.roomApp.GetGameRoom(ctx, req.RoomId)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetGameRoom
	}

	return gr, nil
}

func (r *RoomHandler) getRoomInfoByCode(ctx *gin.Context, req *dto.GetRoomByCodeRequest) (*dto.GameRoom, error) {
	gr, err := r.roomApp.GetGameRoomByCode(ctx, req.RoomCode)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetGameRoom
	}

	return gr, nil
}

func (r *RoomHandler) updateGameRoomStatus(ctx *gin.Context, req *dto.UpdateGameRoomRequest) error {
	userId, err := r.getCurrentUserId(ctx)
	if err != nil {
		return err
	}

	err = r.roomApp.UpdateGameRoomStatus(ctx, userId, req)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrUpdateGameRoom
	}

	return nil
}

func (r *RoomHandler) updateGameRoomExtra(ctx *gin.Context, req *dto.UpdateGameRoomRequest) error {
	userId, err := r.getCurrentUserId(ctx)
	if err != nil {
		return err
	}

	err = r.roomApp.UpdateGameRoomStatus(ctx, userId, req)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrUpdateGameRoom
	}

	return nil
}

// enterGameRoom support both virtual user and real user
// both virtual user and real user need to provide userName
// but virtual user need to set isVirtual to true
func (r *RoomHandler) enterGameRoom(ctx *gin.Context, req *dto.EnterGameRoomRequest) error {
	currentUid, err := r.middleware.CurrentUserId(ctx)
	if err != nil {
		return err
	}

	roomId, virtualUser := req.RoomId, req.VirtualUser

	err = r.roomApp.EnterGameRoom(ctx, roomId, currentUid, virtualUser)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrEnterGameRoom
	}

	return nil
}

func (r *RoomHandler) leaveGameRoom(ctx *gin.Context, req *dto.LeaveGameRoomRequest) error {
	currentUid, err := r.middleware.CurrentUserId(ctx)
	if err != nil {
		return err
	}

	roomId, virtualUser := req.RoomId, req.VirtualUser

	err = r.roomApp.LeaveGameRoom(ctx, roomId, currentUid, virtualUser)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrLeaveGameRoom
	}

	return nil
}

func (r *RoomHandler) getUserGameRoom(ctx *gin.Context) (*dto.GetUserGameRoomsResp, error) {
	userId, err := r.middleware.CurrentUserId(ctx)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		plog.Errorc(ctx, "get user game room error: %v", err)
		return nil, exception.ErrGetGameRoomList
	}

	rooms, err := r.roomApp.GetUserGameRooms(ctx, userId)
	if exception.CheckException(err) {
		return nil, errors.Cause(err)
	} else if err != nil {
		return nil, exception.ErrGetGameRoomList
	}

	return &dto.GetUserGameRoomsResp{Rooms: rooms}, nil
}

func (r *RoomHandler) startGame(ctx *gin.Context, req *dto.StartGameRequest) error {
	userId, err := r.getCurrentUserId(ctx)
	if err != nil {
		return err
	}

	err = r.roomApp.StartGame(ctx, userId, req.RoomId, req.Extra)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrStartGame
	}

	return nil
}

func (r *RoomHandler) endGame(ctx *gin.Context, req *dto.EndGameRequest) error {
	userId, err := r.getCurrentUserId(ctx)
	if err != nil {
		return err
	}

	err = r.roomApp.EndGame(ctx, userId, req.RoomId)
	if exception.CheckException(err) {
		return errors.Cause(err)
	} else if err != nil {
		return exception.ErrStartGame
	}

	return nil
}
