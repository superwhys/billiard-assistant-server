package roomSrv

import (
	"context"

	"github.com/go-puzzles/puzzles/predis"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/models"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/locker"
)

var _ room.IRoomService = (*RoomService)(nil)

type RoomService struct {
	roomRepo   room.IRoomRepo
	locker     *locker.Locker
	roomConfig *models.RoomConfig
}

func NewRoomService(remoRepo room.IRoomRepo, redisClient *predis.RedisClient, roomConfig *models.RoomConfig) *RoomService {
	return &RoomService{
		roomRepo: remoRepo,
		locker:   locker.NewLocker(redisClient, locker.WithPrefix("billiard:room")),
	}
}

func (r *RoomService) CreateGameRoom(ctx context.Context, userId, gameId int) (*room.Room, error) {
	roomCnt, err := r.roomRepo.GetOwnerRoomCount(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "getOwnerRoomCount")
	}

	if roomCnt+1 > r.roomConfig.UserMaxRoomCreateNumber {
		return nil, exception.ErrRoomUserMaxCreateNumber
	}

	ro, err := r.roomRepo.CreateRoom(ctx, gameId, userId)
	if err != nil {
		return nil, errors.Wrap(err, "createGameRoom")
	}

	return ro, nil
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

func (r *RoomService) EnterGameRoom(ctx context.Context, userId, roomId int) (room.User, error) {
	if err := r.locker.Lock(roomId); err != nil {
		return nil, errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(roomId)

	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return nil, errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if !room.CanEnter() {
		return nil, exception.ErrGameRoomFull
	}

	return r.roomRepo.AddUserToRoom(ctx, userId, roomId)
}

func (r *RoomService) QuitGameRoom(ctx context.Context, userId, roomId int) (room.User, error) {
	return r.roomRepo.RemoveUserFromRoom(ctx, userId, roomId)
}

func (r *RoomService) GetUserGameRooms(ctx context.Context, userId int, justOwner bool) ([]*room.Room, error) {
	return r.roomRepo.GetUserGameRooms(ctx, userId, justOwner)
}

func (r *RoomService) GetRoomById(ctx context.Context, roomId int) (*room.Room, error) {
	return r.roomRepo.GetRoomById(ctx, roomId)
}

func (r *RoomService) PrepareGame(ctx context.Context, userId, roomId int) error {
	ro, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if !ro.IsInRoom(userId) {
		return exception.ErrUserNotInRoom
	}

	return r.roomRepo.UpdatePlayerPrepared(ctx, userId, roomId, true)
}

func (r *RoomService) StartGame(ctx context.Context, userId, roomId int) (room.Game, error) {
	if err := r.locker.Lock(roomId); err != nil {
		return nil, errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(roomId)

	ro, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return nil, errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if ro.IsOwner(userId) {
		return nil, exception.ErrNotRoomOwner
	}

	if !ro.CanStart() {
		return nil, exception.ErrPlayerNotReady
	}

	ro.StartGame()
	if err := r.roomRepo.UpdateRoom(ctx, ro); err != nil {
		return nil, errors.Wrap(err, "update room status")
	}

	return ro.Game, nil
}
