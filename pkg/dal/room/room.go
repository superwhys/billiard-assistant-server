package roomDal

import (
	"context"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/putils"
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
	return r.db.Transaction(func(tx *base.Query) error {
		roomDb := tx.RoomPo
		roomUserPo := tx.RoomUserPo

		_, err := roomDb.WithContext(ctx).Where(roomDb.ID.Eq(roomId)).Delete()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.ErrGameRoomNotFound
		}

		_, err = roomUserPo.WithContext(ctx).Where(roomUserPo.RoomID.Eq(roomId)).Delete()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		return nil
	})
}

type enterLeaveOperaion func(ctx context.Context, roomId, userId int, userName string, isVirtual bool) error

func (r *RoomRepoImpl) updateUserRoom(ctx context.Context, roomId, userId int, userName string, isVirtual bool, operation enterLeaveOperaion) error {
	userDb := r.db.UserPo
	roomDb := r.db.RoomPo

	if userId == 0 && (isVirtual && userName == "") {
		return errors.New("oneof userId or virtualUser must be provided")
	}

	roomCount, err := roomDb.WithContext(ctx).Where(roomDb.ID.Eq(roomId)).Count()
	if err != nil {
		return err
	}

	if roomCount == 0 {
		return exception.ErrGameRoomNotFound
	}

	if userId != 0 {
		userCnt, err := userDb.WithContext(ctx).Where(userDb.ID.Eq(userId)).Count()
		if err != nil || userCnt == 0 {
			if err != nil {
				plog.Errorc(ctx, "count userId count error: %v", err)
			}
			return exception.ErrUserNotFound
		}
	}

	return operation(ctx, roomId, userId, userName, isVirtual)
}

func (r *RoomRepoImpl) enterRoom(ctx context.Context, roomId, userId int, userName string, isVirtual bool) error {
	roomUserPo := r.db.RoomUserPo

	roomUser := &model.RoomUserPo{
		RoomID:          roomId,
		UserID:          userId,
		VirtualName:     userName,
		IsVirtualPlayer: isVirtual,
	}
	return roomUserPo.WithContext(ctx).Create(roomUser)
}

func (r *RoomRepoImpl) leaveRoom(ctx context.Context, roomId, userId int, userName string, isVirtual bool) error {
	roomUserPo := r.db.RoomUserPo

	condition := []gen.Condition{}
	if !isVirtual {
		condition = append(condition, roomUserPo.UserID.Eq(userId))
	} else {
		condition = append(condition, roomUserPo.VirtualName.Eq(userName))
	}

	_, err := roomUserPo.WithContext(ctx).Where(condition...).Delete()
	return err
}

func (r *RoomRepoImpl) AddUserToRoom(ctx context.Context, roomId, userId int, userName string, isVirtual bool) error {
	return r.updateUserRoom(ctx, roomId, userId, userName, isVirtual, r.enterRoom)
}

func (r *RoomRepoImpl) RemoveUserFromRoom(ctx context.Context, roomId, userId int, userName string, isVirtual bool) error {
	return r.updateUserRoom(ctx, roomId, userId, userName, isVirtual, r.leaveRoom)
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

	userRooms, err := roomUserDb.WithContext(ctx).
		Preload(roomUserDb.User).
		Where(roomUserDb.RoomID.Eq(roomId)).Find()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrGameRoomNotFound
	} else if err != nil {
		return nil, err
	}

	roomEntity := ro.ToEntity()
	for _, ur := range userRooms {
		roomEntity.Players = append(roomEntity.Players, ur.ToEntity())
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
	roomUserDb := r.db.RoomUserPo
	roomDb := r.db.RoomPo
	userDb := r.db.UserPo

	user, err := userDb.WithContext(ctx).
		Preload(userDb.Rooms.Order(roomDb.CreatedAt.Desc())).
		Preload(userDb.Rooms.Players).
		Preload(userDb.Rooms.Game).
		Where(userDb.ID.Eq(userId)).
		First()
	if err != nil {
		return nil, err
	}

	roomIds := putils.Map(user.Rooms, func(r *model.RoomPo) int {
		return r.ID
	})

	roomUsers, err := roomUserDb.WithContext(ctx).
		Preload(roomUserDb.User).
		Where(roomUserDb.RoomID.In(roomIds...)).Find()
	if err != nil {
		return nil, err
	}

	roomUserEntities := putils.Map(roomUsers, func(ru *model.RoomUserPo) *room.RoomPlayer {
		return ru.ToEntity()
	})

	roomUsersGroup := putils.GroupBy(roomUserEntities, func(rp *room.RoomPlayer) int {
		return rp.RoomId
	})

	var rooms []*room.Room
	for _, room := range user.Rooms {
		roomEntity := room.ToEntity()
		roomEntity.Players = roomUsersGroup[room.ID]
		rooms = append(rooms, roomEntity)
	}

	return rooms, nil
}
