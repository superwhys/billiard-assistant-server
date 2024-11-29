package model

import (
	"time"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
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

	CreatedAt   time.Time
	UpdatedAt   time.Time
	HeartbeatAt time.Time
}

func (ur *RoomUserPo) TableName() string {
	return "room_users"
}

type RoomPo struct {
	ID int `gorm:"primarykey"`

	RoomCode string `gorm:"unique"`

	GameID int
	Game   *GamePo `gorm:"foreignKey:GameID"`

	OwnerID int
	Owner   *UserPo `gorm:"foreignKey:OwnerID"`

	RoomUsers []*RoomUserPo `gorm:"foreignKey:RoomID"`

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
		HeartbeatAt:     r.HeartbeatAt,
	}
}

func (r *RoomPo) ToEntity() *room.Room {
	if r == nil {
		return nil
	}

	var game shared.BaseGame
	if r.Game != nil {
		game = r.Game.ToEntity()
	}

	var owner shared.BaseUser
	if r.Owner != nil {
		owner = r.Owner.ToEntity()
	}

	var players []shared.RoomPlayer
	for _, ur := range r.RoomUsers {
		players = append(players, ur.ToEntity())
	}

	return &room.Room{
		RoomId:        r.ID,
		RoomCode:      r.RoomCode,
		GameId:        r.GameID,
		OwnerId:       r.OwnerID,
		Game:          game,
		Owner:         owner,
		Players:       players,
		GameStatus:    r.GameStatus,
		WinLoseStatus: r.WinLoseStatus,
		CreateAt:      r.CreatedAt,
	}
}
