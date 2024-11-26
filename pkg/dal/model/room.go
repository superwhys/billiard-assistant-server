package model

import (
	"time"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gorm.io/gorm"
)

type RoomUserPo struct {
	ID int `gorm:"primaryKey"`

	RoomID int
	Room   *RoomPo `gorm:"foreignKey:RoomID"`
	UserID int
	User   *UserPo `gorm:"foreignKey:UserID"`

	VirtualName     string
	IsVirtualPlayer bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ur *RoomUserPo) TableName() string {
	return "room_users"
}

type RoomPo struct {
	ID int `gorm:"primarykey"`

	GameID int
	Game   *GamePo `gorm:"foreignKey:GameID"`

	OwnerID int
	Owner   *UserPo `gorm:"foreignKey:OwnerID"`

	Players []*UserPo `gorm:"many2many:room_users;foreignKey:ID;joinForeignKey:RoomID;References:ID;joinReferences:user_id"`

	GameStatus    room.Status
	WinLoseStatus room.WinLoseStatus

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (r *RoomPo) TableName() string {
	return "rooms"
}

func (r *RoomUserPo) ToEntity() *room.RoomPlayer {
	var userName string
	if r.IsVirtualPlayer {
		userName = r.VirtualName
	} else if r.User != nil {
		userName = r.User.Name
	}

	return &room.RoomPlayer{
		RoomId:          r.RoomID,
		UserId:          r.UserID,
		UserName:        userName,
		IsVirtualPlayer: r.IsVirtualPlayer,
	}
}

func (r *RoomPo) FromEntity(gr *room.Room) *RoomPo {
	r.ID = gr.RoomId
	r.GameID = gr.GameId
	r.OwnerID = gr.OwnerId

	r.GameStatus = gr.GameStatus
	r.WinLoseStatus = gr.WinLoseStatus
	return r
}

func (r *RoomPo) ToEntity() *room.Room {
	if r == nil {
		return nil
	}

	var game room.Game
	if r.Game != nil {
		game = r.Game.ToEntity()
	}

	return &room.Room{
		RoomId:        r.ID,
		GameId:        r.GameID,
		OwnerId:       r.OwnerID,
		Game:          game,
		GameStatus:    r.GameStatus,
		WinLoseStatus: r.WinLoseStatus,
		CreateAt:      r.CreatedAt,
	}
}
