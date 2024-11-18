package model

import (
	"time"
	
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gorm.io/gorm"
)

type RoomUserPo struct {
	UserID    int  `gorm:"primaryKey"`
	RoomID    int  `gorm:"primaryKey"`
	Prepared  bool `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ur *RoomUserPo) TableName() string {
	return "room_users"
}

type RoomPo struct {
	ID     int `gorm:"primarykey"`
	GameID int
	Game   *GamePo
	
	OwnerID int
	Owner   *UserPo `gorm:"foreignKey:OwnerID"`
	
	Users []*UserPo `gorm:"many2many:room_users;"`
	
	GameStatus    room.Status
	WinLoseStatus room.WinLoseStatus
	
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (r *RoomPo) TableName() string {
	return "rooms"
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
	players := make([]room.Player, 0, len(r.Users))
	
	for _, u := range r.Users {
		players = append(players, room.Player{
			User:     u.ToEntity(),
			Prepared: false,
		})
	}
	
	return &room.Room{
		RoomId:        r.ID,
		GameId:        r.GameID,
		OwnerId:       r.OwnerID,
		Players:       players,
		Game:          r.Game.ToEntity(),
		GameStatus:    r.GameStatus,
		WinLoseStatus: r.WinLoseStatus,
	}
}
