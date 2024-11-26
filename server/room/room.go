package roomSrv

import (
	"context"
	"time"

	"github.com/go-puzzles/puzzles/predis"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitlab.hoven.com/billiard/billiard-assistant-server/models"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/locker"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/roomcode"
)

var _ room.IRoomService = (*RoomService)(nil)

type RoomService struct {
	roomRepo          room.IRoomRepo
	locker            *locker.Locker
	roomConfig        *models.RoomConfig
	roomCodeGenerator *roomcode.RoomCodeGenerator
}

func NewRoomService(remoRepo room.IRoomRepo, redisClient *predis.RedisClient, roomConfig *models.RoomConfig) *RoomService {
	return &RoomService{
		roomRepo:          remoRepo,
		roomConfig:        roomConfig,
		locker:            locker.NewLocker(redisClient, locker.WithPrefix("billiard:room")),
		roomCodeGenerator: roomcode.NewRoomCodeGenerator(redisClient, time.Hour*24*3),
	}
}

func (r *RoomService) CreateGameRoom(ctx context.Context, u *user.User, gameId int) (*room.Room, error) {
	roomCnt, err := r.roomRepo.GetOwnerRoomCount(ctx, u.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "getOwnerRoomCount")
	}

	if roomCnt >= r.roomConfig.UserMaxRoomCreateNumber {
		return nil, exception.ErrRoomUserMaxCreateNumber
	}

	room, err := r.roomRepo.CreateRoom(ctx, gameId, u.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "createRoom")
	}

	roomCode, err := r.roomCodeGenerator.GenerateCode(ctx, room.RoomId)
	if err != nil {
		return nil, errors.Wrap(err, "generateRoomCode")
	}
	room.RoomCode = roomCode

	return room, nil
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

func (r *RoomService) EnterGameRoom(ctx context.Context, roomId, userId int, userName string, isVirtual bool) error {
	if err := r.locker.Lock(roomId); err != nil {
		return errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(roomId)

	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if room.IsEnd() {
		return exception.ErrGameRoomEnd
	}

	if room.IsInRoom(userName, userId) {
		return exception.ErrAlreadyInRoom
	}

	if !room.CanEnter() {
		return exception.ErrGameRoomFull
	}

	return r.roomRepo.AddUserToRoom(ctx, roomId, userId, userName, isVirtual)
}

func (r *RoomService) QuitGameRoom(ctx context.Context, roomId, userId int, userName string, isVirtual bool) error {
	if err := r.locker.Lock(roomId); err != nil {
		return errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(roomId)

	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getRoom: %d", roomId)
	}

	if room.IsEnd() {
		return exception.ErrGameRoomEnd
	}

	if !room.IsInRoom(userName, userId) {
		return exception.ErrNotInRoom
	}

	return r.roomRepo.RemoveUserFromRoom(ctx, roomId, userId, userName, isVirtual)
}

func (r *RoomService) GetUserGameRooms(ctx context.Context, userId int, justOwner bool) ([]*room.Room, error) {
	return r.roomRepo.GetUserGameRooms(ctx, userId, justOwner)
}

func (r *RoomService) GetRoomById(ctx context.Context, roomId int) (*room.Room, error) {
	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return nil, errors.Wrapf(err, "getRoom: %d", roomId)
	}

	roomCode, err := r.roomCodeGenerator.GetRoomCode(ctx, roomId)
	if err != nil && !errors.Is(err, redis.ErrNil) {
		return nil, errors.Wrap(err, "getRoomCode")
	}

	if roomCode != "" {
		room.RoomCode = roomCode
		return room, nil
	}

	roomCode, err = r.roomCodeGenerator.GenerateCode(ctx, roomId)
	if err != nil {
		return nil, errors.Wrap(err, "generateRoomCode")
	}

	room.RoomCode = roomCode
	return room, nil
}

func (r *RoomService) GetRoomByCode(ctx context.Context, roomCode string) (*room.Room, error) {
	roomId, err := r.roomCodeGenerator.GetRoomId(ctx, roomCode)
	if err != nil {
		return nil, errors.Wrapf(err, "getRoomId: %s", roomCode)
	}

	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return nil, errors.Wrapf(err, "getRoom: %d", roomId)
	}

	room.RoomCode = roomCode
	return room, nil
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

func (r *RoomService) EndGame(ctx context.Context, userId, roomId int) error {
	if err := r.locker.Lock(roomId); err != nil {
		return errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(roomId)

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

	err = r.roomCodeGenerator.DeleteCode(ctx, roomId)
	if err != nil {
		return errors.Wrap(err, "delete roomCode")
	}

	return nil
}
