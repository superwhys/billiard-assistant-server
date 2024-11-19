package model

import (
	"time"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gorm.io/gorm"
)

type UserPo struct {
	ID     int    `gorm:"primaryKey"`
	Name   string `gorm:"type:varchar(50);not null"`
	Email  string `gorm:"type:varchar(100);not null"`
	Phone  string `gorm:"type:varchar(11)"`
	Avatar string `gorm:"type:varchar(255)"`
	Status user.Status
	Role   user.Role

	Rooms       []*RoomPo     `gorm:"many2many:room_users;"`
	UserAuthPos []*UserAuthPo `gorm:"constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *UserPo) TableName() string {
	return "users"
}

func (u *UserPo) FromEntity(ue *user.User) *UserPo {
	u.ID = ue.UserId
	u.Name = ue.Name
	u.Status = ue.Status
	u.Role = ue.Role
	if ue.UserInfo != nil {
		u.Email = ue.UserInfo.Email
		u.Phone = ue.UserInfo.Phone
		u.Avatar = ue.UserInfo.Avatar
	}

	return u
}

func (u *UserPo) ToEntity() *user.User {
	var rooms []user.Room
	for _, room := range u.Rooms {
		rooms = append(rooms, room.ToEntity())
	}

	return &user.User{
		UserId: u.ID,
		Name:   u.Name,
		Status: u.Status,
		Role:   u.Role,
		Rooms:  rooms,
		UserInfo: &user.BaseInfo{
			Email:  u.Email,
			Phone:  u.Phone,
			Avatar: u.Avatar,
		},
	}
}
