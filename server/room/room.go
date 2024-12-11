package roomSrv

import (
	"context"

	"github.com/go-puzzles/puzzles/goredis"
	"github.com/pkg/errors"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitea.hoven.com/billiard/billiard-assistant-server/models"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/locker"
	"gorm.io/datatypes"
)

var _ room.IRoomService = (*RoomService)(nil)

type RoomService struct {
	roomRepo   room.IRoomRepo
	locker     *locker.Locker
	roomConfig *models.RoomConfig
}

func NewRoomService(remoRepo room.IRoomRepo, redisClient *goredis.PuzzleRedisClient, roomConfig *models.RoomConfig) *RoomService {
	return &RoomService{
		roomRepo:   remoRepo,
		roomConfig: roomConfig,
		locker:     locker.NewLocker(redisClient, locker.WithPrefix("billiard:room")),
	}
}

func (r *RoomService) CreateGameRoom(ctx context.Context, u *user.User, gameId int) (*room.Room, error) {
	roomCnt, err := r.roomRepo.GetOwnerRoomCount(ctx, u.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "getOwnerRoomCount")
	}

	if !u.Role.IsPro() && roomCnt >= r.roomConfig.UserMaxRoomCreateNumber {
		return nil, exception.ErrRoomUserMaxCreateNumber
	}

	room, err := r.roomRepo.CreateRoom(ctx, gameId, u.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "createRoom")
	}

	return room, nil
}

func (r *RoomService) CheckRoomCodeExists(ctx context.Context, roomCode string) (bool, error) {
	return r.roomRepo.CheckRoomCodeExists(ctx, roomCode)
}

func (r *RoomService) DeleteGameRoom(ctx context.Context, roomId, userId int) error {
	ro, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getGameRoom: %d", roomId)
	}

	if !ro.IsOwner(userId) {
		return exception.ErrRoomOwnerNotMatch
	}

	return r.roomRepo.DeleteRoom(ctx, roomId)
}

func (r *RoomService) UpdateGameRoomStatus(ctx context.Context, gameRoom *room.Room) error {
	if err := r.roomRepo.UpdateRoom(ctx, gameRoom); err != nil {
		return errors.Wrapf(err, "updateGameRoom: %d", gameRoom.RoomId)
	}

	return nil
}

func (r *RoomService) UpdateRoomUserHeartbeart(ctx context.Context, roomId, userId int) error {
	if err := r.roomRepo.UpdateRoomUserHeartbeart(ctx, roomId, userId); err != nil {
		return errors.Wrapf(err, "updateUserHeartbeart: %d, %d", roomId, userId)
	}

	return nil
}

func (r *RoomService) EnterGameRoom(ctx context.Context, roomId int, currentUser shared.BaseUser, virtualUser string) error {
	if err := r.locker.Lock(ctx, roomId); err != nil {
		return errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(ctx, roomId)

	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if room.IsEnd() {
		return exception.ErrGameRoomEnd
	}

	// not allow to add a virtual player while not roomOwner
	if virtualUser != "" && !room.IsOwner(currentUser.GetUserId()) {
		return exception.ErrNotRoomOwner
	}

	if room.IsInRoom(currentUser.GetUserId(), virtualUser) {
		return exception.ErrAlreadyInRoom
	}

	if !room.CanEnter() {
		return exception.ErrGameRoomFull
	}

	addUser := currentUser.GetName()
	if virtualUser != "" {
		addUser = virtualUser
	}

	return r.roomRepo.AddUserToRoom(ctx, roomId, currentUser.GetUserId(), addUser, virtualUser != "")
}

func (r *RoomService) QuitGameRoom(ctx context.Context, roomId int, currentUser shared.BaseUser, virtualUser string) error {
	if err := r.locker.Lock(ctx, roomId); err != nil {
		return errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(ctx, roomId)

	isVirtual := virtualUser != ""

	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if room.IsEnd() {
		return exception.ErrGameRoomEnd
	}

	// not allow to quit virtual player while not roomOwner
	if isVirtual && !room.IsOwner(currentUser.GetUserId()) {
		return exception.ErrNotRoomOwner
	}

	if !room.IsInRoom(currentUser.GetUserId(), virtualUser) {
		return exception.ErrNotInRoom
	}

	leaveUser := currentUser.GetName()
	if isVirtual {
		leaveUser = virtualUser
	}

	return r.roomRepo.RemoveUserFromRoom(ctx, roomId, currentUser.GetUserId(), leaveUser, isVirtual)
}

func (r *RoomService) GetUserGameRooms(ctx context.Context, userId int) ([]*room.Room, error) {
	rooms, err := r.roomRepo.GetUserGameRooms(ctx, userId)
	if err != nil {
		return nil, errors.Wrapf(err, "getUserGameRooms: %d", userId)
	}

	return rooms, nil
}

func (r *RoomService) GetRoomGameType(ctx context.Context, roomId int) (shared.BilliardGameType, error) {
	return r.roomRepo.GetRoomGameType(ctx, roomId)
}

func (r *RoomService) GetRoomById(ctx context.Context, roomId int) (*room.Room, error) {
	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return nil, errors.Wrapf(err, "getRoom: %d", roomId)
	}

	return room, nil
}

func (r *RoomService) GetRoomByCode(ctx context.Context, roomCode string) (*room.Room, error) {
	room, err := r.roomRepo.GetRoomByRoomCode(ctx, roomCode)
	if err != nil {
		return nil, errors.Wrapf(err, "getRoom: %s", roomCode)
	}

	return room, nil
}

func (r *RoomService) StartGame(ctx context.Context, userId, roomId int, extra datatypes.JSONMap) (shared.BaseGame, error) {
	if err := r.locker.Lock(ctx, roomId); err != nil {
		return nil, errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(ctx, roomId)

	ro, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return nil, errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if !ro.IsOwner(userId) {
		return nil, exception.ErrNotRoomOwner
	}

	if !ro.CanStart() {
		return nil, exception.ErrStartGame
	}

	ro.StartGame()
	ro.SetExtra(extra)
	if err := r.roomRepo.UpdateRoom(ctx, ro); err != nil {
		return nil, errors.Wrap(err, "update room status")
	}

	return ro.Game, nil
}

func (r *RoomService) EndGame(ctx context.Context, userId, roomId int) error {
	if err := r.locker.Lock(ctx, roomId); err != nil {
		return errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(ctx, roomId)

	ro, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if !ro.IsOwner(userId) {
		return exception.ErrNotRoomOwner
	}

	ro.EndGame()
	if err := r.roomRepo.UpdateRoom(ctx, ro); err != nil {
		return errors.Wrap(err, "update room status")
	}

	return nil
}
