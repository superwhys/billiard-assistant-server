package roomSrv

import (
	"context"

	"github.com/go-puzzles/puzzles/predis"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
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

func (r *RoomService) EnterGameRoom(ctx context.Context, roomId int, enterUser shared.BaseUser, isVirtual bool) error {
	if err := r.locker.Lock(roomId); err != nil {
		return errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(roomId)

	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getRoom: %d", roomId)
	}

	// not allow to add a virtual player while not roomOwner
	if !room.IsOwner(enterUser.GetUserId()) && isVirtual {
		return exception.ErrNotRoomOwner
	}

	if room.IsEnd() {
		return exception.ErrGameRoomEnd
	}

	if room.IsInRoom(isVirtual, enterUser.GetName(), enterUser.GetUserId()) {
		return exception.ErrAlreadyInRoom
	}

	if !room.CanEnter() {
		return exception.ErrGameRoomFull
	}

	return r.roomRepo.AddUserToRoom(ctx, roomId, enterUser.GetUserId(), enterUser.GetName(), isVirtual)
}

func (r *RoomService) QuitGameRoom(ctx context.Context, roomId int, quitUser shared.BaseUser, isVirtual bool) error {
	if err := r.locker.Lock(roomId); err != nil {
		return errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(roomId)

	room, err := r.roomRepo.GetRoomById(ctx, roomId)
	if err != nil {
		return errors.Wrapf(err, "getRoom: %d", roomId)
	}

	// not allow to quit a virtual player while not roomOwner
	if !room.IsOwner(quitUser.GetUserId()) && isVirtual {
		return exception.ErrNotRoomOwner
	}

	if room.IsEnd() {
		return exception.ErrGameRoomEnd
	}

	if !room.IsInRoom(isVirtual, quitUser.GetName(), quitUser.GetUserId()) {
		return exception.ErrNotInRoom
	}

	return r.roomRepo.RemoveUserFromRoom(ctx, roomId, quitUser.GetUserId(), quitUser.GetName(), isVirtual)
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

func (r *RoomService) StartGame(ctx context.Context, userId, roomId int) (shared.BaseGame, error) {
	if err := r.locker.Lock(roomId); err != nil {
		return nil, errors.Wrap(err, "lock room")
	}
	defer r.locker.Unlock(roomId)

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

	return nil
}
