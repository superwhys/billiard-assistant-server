package roomDal

import (
	"context"
	
	"github.com/pkg/errors"
	"github.com/superwhys/snooker-assistant-server/domain/room"
	"github.com/superwhys/snooker-assistant-server/pkg/dal/base"
	"github.com/superwhys/snooker-assistant-server/pkg/dal/model"
	"github.com/superwhys/snooker-assistant-server/pkg/exception"
	"gorm.io/gorm"
)

var _ room.IRoomRepo = (*RoomRepoImpl)(nil)

type RoomRepoImpl struct {
	roomDb *base.RoomDB
	userDb *base.UserDB
}

func NewRoomRepo(db *gorm.DB) *RoomRepoImpl {
	return &RoomRepoImpl{
		roomDb: base.NewRoomDB(db),
		userDb: base.NewUserDB(db),
	}
}

func (r *RoomRepoImpl) CreateRoom(ctx context.Context, gameId, userId int) (*room.Room, error) {
	user, err := r.userDb.WithContext(ctx).Where(r.userDb.ID.Eq(userId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	
	ro := &model.RoomPo{
		GameID:        gameId,
		OwnerID:       userId,
		GameStatus:    room.Playing,
		WinLoseStatus: room.Unknown,
	}
	
	err = r.roomDb.WithContext(ctx).Create(ro)
	if err != nil {
		return nil, err
	}
	
	if err := r.roomDb.Users.Model(ro).Append(user); err != nil {
		return nil, errors.Wrap(err, "appendRelation")
	}
	
	return ro.ToEntity(), nil
}

func (r *RoomRepoImpl) UpdateRoom(ctx context.Context, room *room.Room) error {
	roomPo := new(model.RoomPo).FromEntity(room)
	_, err := r.roomDb.WithContext(ctx).Where(r.roomDb.ID.Eq(room.RoomId)).Updates(roomPo)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrGameRoomNotFound
	}
	
	return nil
	
}

func (r *RoomRepoImpl) DeleteRoom(ctx context.Context, roomId int) error {
	_, err := r.roomDb.WithContext(ctx).Where(r.roomDb.ID.Eq(roomId)).Delete()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrGameRoomNotFound
	}
	
	return nil
}

func (r *RoomRepoImpl) updateUserRoom(ctx context.Context, userId, roomId int, operation func(*model.RoomPo, *model.UserPo) error) error {
	user, err := r.userDb.WithContext(ctx).Where(r.userDb.ID.Eq(userId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrUserNotFound
	} else if err != nil {
		return err
	}
	
	room, err := r.roomDb.WithContext(ctx).Where(r.roomDb.ID.Eq(roomId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrGameRoomNotFound
	} else if err != nil {
		return err
	}
	
	return operation(room, user)
}

func (r *RoomRepoImpl) enterRoom(room *model.RoomPo, user *model.UserPo) error {
	return r.roomDb.Users.Model(room).Append(user)
}

func (r *RoomRepoImpl) leaveRoom(room *model.RoomPo, user *model.UserPo) error {
	return r.roomDb.Users.Model(room).Delete(user)
}

func (r *RoomRepoImpl) AddUserToRoom(ctx context.Context, userId, roomId int) error {
	return r.updateUserRoom(ctx, userId, roomId, r.enterRoom)
}

func (r *RoomRepoImpl) RemoveUserFromRoom(ctx context.Context, userId, roomId int) error {
	return r.updateUserRoom(ctx, userId, roomId, r.leaveRoom)
}

func (r *RoomRepoImpl) GetRoomById(ctx context.Context, roomId int) (*room.Room, error) {
	ro, err := r.roomDb.WithContext(ctx).
		Preload(r.roomDb.Users).
		Preload(r.roomDb.Game).
		Where(r.roomDb.ID.Eq(roomId)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrGameRoomNotFound
	} else if err != nil {
		return nil, err
	}
	
	return ro.ToEntity(), nil
}

func (r *RoomRepoImpl) GetUserGameRooms(ctx context.Context, userId int) ([]*room.Room, error) {
	user, err := r.userDb.WithContext(ctx).
		Preload(r.userDb.Rooms).
		Preload(r.userDb.Rooms.Users).
		Where(r.userDb.ID.Eq(userId)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	
	ret := make([]*room.Room, 0, len(user.Rooms))
	for _, room := range user.Rooms {
		ret = append(ret, room.ToEntity())
	}
	
	return ret, nil
}
