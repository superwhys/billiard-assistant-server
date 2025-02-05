package model

import (
	"time"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gorm.io/gorm"
)

type UserPo struct {
	ID int `gorm:"primaryKey"`

	RoomUsers []*RoomUserPo `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *UserPo) TableName() string {
	return "users"
}

func (u *UserPo) FromEntity(ue *user.User) *UserPo {
	u.ID = ue.UserId

	return u
}

func (u *UserPo) ToEntity() *user.User {
	if u == nil {
		return nil
	}

	var rooms []user.Room
	for _, room := range u.RoomUsers {
		rooms = append(rooms, room.ToEntity())
	}

	return &user.User{
		UserId: u.ID,
		Rooms:  rooms,
	}
}
