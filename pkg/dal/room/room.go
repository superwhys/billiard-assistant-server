package roomDal

import (
	"context"
	"slices"

	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/base"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gorm.io/gen"
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
	roomDb := r.db.RoomPo

	ro := &model.RoomPo{
		GameID:        gameId,
		OwnerID:       userId,
		GameStatus:    room.Playing,
		WinLoseStatus: room.WinLoseUnknown,
	}

	err := roomDb.WithContext(ctx).Create(ro)
	if err != nil {
		return nil, err
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

type enterLeaveOperaion func(ctx context.Context, virtualUser string, userId, roomId int) error

func (r *RoomRepoImpl) updateUserRoom(ctx context.Context, virtualUser string, userId int, roomId int, operation enterLeaveOperaion) error {
	userDb := r.db.UserPo
	roomDb := r.db.RoomPo

	if userId == 0 && virtualUser == "" {
		return errors.New("oneof userId or virtualUser must be provided")
	}

	userName := virtualUser
	if userId != 0 {
		user, err := userDb.WithContext(ctx).Where(userDb.ID.Eq(userId)).First()
		if errors.Is(err, exception.ErrUserNotFound) {
			return err
		} else if err != nil {
			return err
		}
		userName = user.Name
	}

	roomCount, err := roomDb.WithContext(ctx).Where(roomDb.ID.Eq(roomId)).Count()
	if err != nil {
		return err
	}

	if roomCount == 0 {
		return exception.ErrGameRoomNotFound
	}

	return operation(ctx, userName, userId, roomId)
}

func (r *RoomRepoImpl) enterRoom(ctx context.Context, virtualUser string, userId, roomId int) error {
	roomUserPo := r.db.RoomUserPo

	roomUser := &model.RoomUserPo{
		RoomID:          roomId,
		UserID:          userId,
		UserName:        virtualUser,
		IsVirtualPlayer: userId == 0,
	}
	return roomUserPo.WithContext(ctx).Create(roomUser)
}

func (r *RoomRepoImpl) leaveRoom(ctx context.Context, virtualUser string, userId, roomId int) error {
	roomUserPo := r.db.RoomUserPo

	condition := []gen.Condition{}
	if userId != 0 {
		condition = append(condition, roomUserPo.UserID.Eq(userId))
	} else {
		condition = append(condition, roomUserPo.UserName.Eq(virtualUser))
	}

	_, err := roomUserPo.WithContext(ctx).Where(condition...).Delete()
	return err
}

func (r *RoomRepoImpl) AddUserToRoom(ctx context.Context, virtualUser string, userId, roomId int) error {
	return r.updateUserRoom(ctx, virtualUser, userId, roomId, r.enterRoom)
}

func (r *RoomRepoImpl) RemoveUserFromRoom(ctx context.Context, virtualUser string, userId, roomId int) error {
	return r.updateUserRoom(ctx, virtualUser, userId, roomId, r.leaveRoom)
}

func (r *RoomRepoImpl) GetRoomById(ctx context.Context, roomId int) (*room.Room, error) {
	roomDb := r.db.RoomPo
	roomUserDb := r.db.RoomUserPo

	ro, err := roomDb.WithContext(ctx).
		Preload(roomDb.Owner, roomDb.Game).
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

	roomEntity := ro.ToEntity()
	for _, user := range userRooms {
		roomEntity.Players = append(roomEntity.Players, &room.RoomPlayer{
			UserId:          user.UserID,
			UserName:        user.UserName,
			IsVirtualPlayer: user.IsVirtualPlayer,
		})
	}
	roomEntity.Game = ro.Game.ToEntity()

	return roomEntity, nil
}

func (r *RoomRepoImpl) GetOwnerRoomCount(ctx context.Context, userId int) (int64, error) {
	roomDb := r.db.RoomPo

	count, err := roomDb.WithContext(ctx).
		Where(roomDb.OwnerID.Eq(userId)).
		Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RoomRepoImpl) GetUserGameRooms(ctx context.Context, userId int, justOwner bool) ([]*room.Room, error) {
	roomUserPo := r.db.RoomUserPo

	userRooms, err := roomUserPo.WithContext(ctx).
		Where(roomUserPo.UserID.Eq(userId)).
		Find()
	if err != nil {
		return nil, err
	}

	if len(userRooms) == 0 {
		return []*room.Room{}, nil
	}

	roomIds := make([]int, 0, len(userRooms))
	for _, ur := range userRooms {
		roomIds = append(roomIds, ur.RoomID)
	}

	allRoomUsers, err := roomUserPo.WithContext(ctx).
		Preload(roomUserPo.Room).
		Preload(roomUserPo.Room.Game).
		Where(roomUserPo.RoomID.In(roomIds...)).
		Find()
	if err != nil {
		return nil, err
	}

	roomMap := make(map[int]*room.Room)
	for _, ru := range allRoomUsers {
		roomId := ru.RoomID
		if _, exists := roomMap[roomId]; !exists {
			roomMap[roomId] = ru.Room.ToEntity()
			roomMap[roomId].Game = ru.Room.Game.ToEntity()
		}

		roomMap[roomId].Players = append(roomMap[roomId].Players, &room.RoomPlayer{
			UserId:          ru.UserID,
			UserName:        ru.UserName,
			IsVirtualPlayer: ru.IsVirtualPlayer,
		})
	}

	ret := make([]*room.Room, 0, len(roomMap))
	for _, r := range roomMap {
		ret = append(ret, r)
	}

	slices.SortStableFunc[[]*room.Room, *room.Room](ret, func(a, b *room.Room) int {
		if a.CreateAt.Before(b.CreateAt) {
			return 1
		} else if a.CreateAt.After(b.CreateAt) {
			return -1
		}
		return 0
	})

	return ret, nil
}
