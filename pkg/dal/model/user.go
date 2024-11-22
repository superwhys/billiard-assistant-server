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
	Gender user.Gender
	Status user.Status
	Role   user.Role

	Rooms       []*RoomPo     `gorm:"many2many:room_users;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:room_id"`
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
	u.Gender = ue.Gender
	u.Role = ue.Role
	if ue.UserInfo != nil {
		u.Email = ue.UserInfo.Email
		u.Phone = ue.UserInfo.Phone
		u.Avatar = ue.UserInfo.Avatar
	}

	return u
}

func (u *UserPo) ToEntity() *user.User {
	if u == nil {
		return nil
	}

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
		Gender: u.Gender,
		UserInfo: &user.BaseInfo{
			Email:  u.Email,
			Phone:  u.Phone,
			Avatar: u.Avatar,
		},
	}
}
