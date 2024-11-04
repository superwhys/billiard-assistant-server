package roomSrv

import (
	"context"

	"github.com/go-puzzles/predis"
	"github.com/pkg/errors"
	"github.com/superwhys/snooker-assistant-server/domain/room"
	"github.com/superwhys/snooker-assistant-server/pkg/exception"
	"github.com/superwhys/snooker-assistant-server/pkg/locker"
)

var _ room.IRoomService = (*RoomService)(nil)

type RoomService struct {
	roomRepo room.IRoomRepo
	locker   *locker.Locker
}

func NewRoomService(remoRepo room.IRoomRepo, redisClient *predis.RedisClient) *RoomService {
	return &RoomService{
		roomRepo: remoRepo,
		locker:   locker.NewLocker(redisClient, locker.WithPrefix("sa:room")),
	}
}

func (r *RoomService) CreateGameRoom(ctx context.Context, userId, gameId int) (*room.Room, error) {
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

	if ro.OwnerId != userId {
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

func (r *RoomService) EnterGameRoom(ctx context.Context, userId, roomId int) error {
	if err := r.locker.Lock(roomId); err != nil {
		return errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(roomId)

	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if room.Game.GetMaxPlayers() >= len(room.Players) {
		return exception.ErrGameRoomFull
	}

	return r.roomRepo.AddUserToRoom(ctx, userId, roomId)
}

func (r *RoomService) QuitGameRoom(ctx context.Context, userId, roomId int) error {
	return r.roomRepo.RemoveUserFromRoom(ctx, userId, roomId)
}

func (r *RoomService) GetUserGameRooms(ctx context.Context, userId int) ([]*room.Room, error) {
	return r.roomRepo.GetUserGameRooms(ctx, userId)
}

func (r *RoomService) GetRoomById(ctx context.Context, roomId int) (*room.Room, error) {
	return r.roomRepo.GetRoomById(ctx, roomId)
}
