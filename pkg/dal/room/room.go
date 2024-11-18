package roomDal

import (
	"context"
	
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/base"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gorm.io/gorm"
)

var _ room.IRoomRepo = (*RoomRepoImpl)(nil)

type RoomRepoImpl struct {
	db *base.Query
}

func NewRoomRepo(db *gorm.DB) *RoomRepoImpl {
	return &RoomRepoImpl{
		db: base.Use(db),
	}
}

func (r *RoomRepoImpl) CreateRoom(ctx context.Context, gameId, userId int) (*room.Room, error) {
	userDb := r.db.UserPo
	roomDb := r.db.RoomPo
	user, err := userDb.WithContext(ctx).Where(userDb.ID.Eq(userId)).First()
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
	
	err = roomDb.WithContext(ctx).Create(ro)
	if err != nil {
		return nil, err
	}
	
	if err := roomDb.Users.Model(ro).Append(user); err != nil {
		return nil, errors.Wrap(err, "appendRelation")
	}
	
	return ro.ToEntity(), nil
}

func (r *RoomRepoImpl) UpdateRoom(ctx context.Context, room *room.Room) error {
	roomDb := r.db.RoomPo
	
	roomPo := new(model.RoomPo).FromEntity(room)
	_, err := roomDb.WithContext(ctx).Where(roomDb.ID.Eq(room.RoomId)).Updates(roomPo)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrGameRoomNotFound
	}
	
	return nil
	
}

func (r *RoomRepoImpl) DeleteRoom(ctx context.Context, roomId int) error {
	roomDb := r.db.RoomPo
	_, err := roomDb.WithContext(ctx).Where(roomDb.ID.Eq(roomId)).Delete()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrGameRoomNotFound
	}
	
	return nil
}

func (r *RoomRepoImpl) updateUserRoom(ctx context.Context, userId, roomId int, operation func(*model.RoomPo, *model.UserPo) (room.User, error)) (room.User, error) {
	userDb := r.db.UserPo
	roomDb := r.db.RoomPo
	
	user, err := userDb.WithContext(ctx).Where(userDb.ID.Eq(userId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	
	room, err := roomDb.WithContext(ctx).Where(roomDb.ID.Eq(roomId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrGameRoomNotFound
	} else if err != nil {
		return nil, err
	}
	
	return operation(room, user)
}

func (r *RoomRepoImpl) enterRoom(room *model.RoomPo, user *model.UserPo) (room.User, error) {
	roomDb := r.db.RoomPo
	return user.ToEntity(), roomDb.Users.Model(room).Append(user)
}

func (r *RoomRepoImpl) leaveRoom(room *model.RoomPo, user *model.UserPo) (room.User, error) {
	roomDb := r.db.RoomPo
	return user.ToEntity(), roomDb.Users.Model(room).Delete(user)
}

func (r *RoomRepoImpl) AddUserToRoom(ctx context.Context, userId, roomId int) (room.User, error) {
	return r.updateUserRoom(ctx, userId, roomId, r.enterRoom)
}

func (r *RoomRepoImpl) RemoveUserFromRoom(ctx context.Context, userId, roomId int) (room.User, error) {
	return r.updateUserRoom(ctx, userId, roomId, r.leaveRoom)
}

func (r *RoomRepoImpl) GetRoomById(ctx context.Context, roomId int) (*room.Room, error) {
	roomDb := r.db.RoomPo
	roomUserDb := r.db.RoomUserPo
	
	ro, err := roomDb.WithContext(ctx).
		Preload(roomDb.Users).
		Preload(roomDb.Game).
		Where(roomDb.ID.Eq(roomId)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrGameRoomNotFound
	} else if err != nil {
		return nil, err
	}
	
	userRooms, err := roomUserDb.WithContext(ctx).Where(roomUserDb.RoomID.Eq(roomId)).Find()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrGameRoomNotFound
	} else if err != nil {
		return nil, err
	}
	
	userPrepared := make(map[int]bool)
	for _, ur := range userRooms {
		userPrepared[ur.UserID] = ur.Prepared
	}
	
	room := ro.ToEntity()
	for i, player := range room.Players {
		room.Players[i].Prepared = userPrepared[player.GetUserId()]
	}
	
	return room, nil
}

func (r *RoomRepoImpl) GetUserGameRooms(ctx context.Context, userId int) ([]*room.Room, error) {
	userDb := r.db.UserPo
	user, err := userDb.WithContext(ctx).
		Preload(userDb.Rooms).
		Preload(userDb.Rooms.Users).
		Where(userDb.ID.Eq(userId)).
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

func (r *RoomRepoImpl) UpdatePlayerPrepared(ctx context.Context, userId, roomId int, prepared bool) error {
	roomUserDb := r.db.RoomUserPo
	result, err := roomUserDb.WithContext(ctx).Where(
		roomUserDb.RoomID.Eq(roomId),
		roomUserDb.UserID.Eq(userId),
	).Update(roomUserDb.Prepared, prepared)
	if err != nil {
		return err
	}
	
	if result.RowsAffected == 0 {
		return exception.ErrUserNotInRoom
	}
	return nil
}
