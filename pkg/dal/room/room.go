package roomDal

import (
	"context"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
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
		GameStatus:    room.Preparing,
		WinLoseStatus: room.WinLoseUnknown,
	}

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		code := putils.RandString(7)
		if exists, _ := r.CheckRoomCodeExists(ctx, code); !exists {
			ro.RoomCode = code
			break
		}
	}

	if ro.RoomCode == "" {
		return nil, errors.New("failed to generate unique room code")
	}

	err := roomDb.WithContext(ctx).Create(ro)
	if err != nil {
		return nil, err
	}

	return ro.ToEntity(), nil
}

func (r *RoomRepoImpl) UpdateRoom(ctx context.Context, room *room.Room) error {
	roomDb := r.db.RoomPo

	roomPo := &model.RoomPo{
		ID:            room.RoomId,
		RoomCode:      room.RoomCode,
		GameStatus:    room.GameStatus,
		WinLoseStatus: room.WinLoseStatus,
	}
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
		IsVirtualPlayer: isVirtual,
		HeartbeatAt:     time.Now(),
	}
	if isVirtual {
		roomUser.VirtualName = userName
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

func (r *RoomRepoImpl) GetRoomGameType(ctx context.Context, roomId int) (shared.BilliardGameType, error) {
	roomDb := r.db.RoomPo
	gameDb := r.db.GamePo

	resp := &struct {
		GameID   int
		GameType shared.BilliardGameType
	}{}

	err := roomDb.WithContext(ctx).
		Select(roomDb.ID, gameDb.GameType).
		Join(gameDb, gameDb.ID.EqCol(roomDb.GameID)).
		Where(roomDb.ID.Eq(roomId)).
		Scan(resp)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shared.GameTypeUnkonwon, exception.ErrGameRoomNotFound
	} else if err != nil {
		return shared.GameTypeUnkonwon, err
	}

	return resp.GameType, nil
}

func (r *RoomRepoImpl) getRoom(ctx context.Context, condition ...gen.Condition) (*room.Room, error) {
	roomDb := r.db.RoomPo
	roomUserDb := r.db.RoomUserPo

	ro, err := roomDb.WithContext(ctx).
		Preload(roomDb.RoomUsers.Order(roomUserDb.CreatedAt.Asc())).
		Preload(roomDb.RoomUsers.User).
		Preload(roomDb.Owner).
		Preload(roomDb.Game).
		Preload(roomDb.RoomUsers).
		Where(condition...).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrGameRoomNotFound
	} else if err != nil {
		return nil, err
	}

	return ro.ToEntity(), nil
}

func (r *RoomRepoImpl) GetRoomById(ctx context.Context, roomId int) (*room.Room, error) {
	roomDb := r.db.RoomPo
	return r.getRoom(ctx, roomDb.ID.Eq(roomId))
}

func (r *RoomRepoImpl) GetRoomByRoomCode(ctx context.Context, roomCode string) (*room.Room, error) {
	roomDb := r.db.RoomPo
	return r.getRoom(ctx, roomDb.RoomCode.Eq(roomCode))
}

func (r *RoomRepoImpl) CheckRoomCodeExists(ctx context.Context, roomCode string) (bool, error) {
	roomDb := r.db.RoomPo

	count, err := roomDb.WithContext(ctx).
		Where(roomDb.RoomCode.Eq(roomCode)).
		Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
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

func (r *RoomRepoImpl) GetUserGameRooms(ctx context.Context, userId int) ([]*room.Room, error) {
	userDb := r.db.UserPo
	roomUserDb := r.db.RoomUserPo

	user, err := userDb.WithContext(ctx).Preload(userDb.RoomUsers).Where(userDb.ID.Eq(userId)).First()
	if err != nil {
		return nil, err
	}

	roomIds := putils.Map(user.RoomUsers, func(ru *model.RoomUserPo) int {
		return ru.RoomID
	})

	roomIds = putils.Dedup(roomIds)

	roomUsers, err := roomUserDb.WithContext(ctx).
		Preload(roomUserDb.Room).
		Preload(roomUserDb.Room.Game).
		Preload(roomUserDb.Room.Owner).
		Where(roomUserDb.RoomID.In(roomIds...)).
		Find()
	if err != nil {
		return nil, err
	}

	roomGroups := putils.GroupBy(roomUsers, func(ru *model.RoomUserPo) int {
		return ru.RoomID
	})

	var rooms []*room.Room
	for roomId, roomGroup := range roomGroups {
		room := roomGroup[0].Room.ToEntity()
		room.RoomId = roomId
		for _, ru := range roomGroup {
			room.Players = append(room.Players, ru.ToEntity())
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomRepoImpl) UpdateRoomUserHeartbeart(ctx context.Context, roomId, userId int) error {
	roomUserDb := r.db.RoomUserPo

	roomUser, err := roomUserDb.WithContext(ctx).
		Where(roomUserDb.RoomID.Eq(roomId)).
		Where(roomUserDb.UserID.Eq(userId)).
		First()
	if err != nil {
		return err
	}

	roomUser.HeartbeatAt = time.Now()
	return roomUserDb.WithContext(ctx).Save(roomUser)
}
