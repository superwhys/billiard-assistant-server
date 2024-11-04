package model

import (
	"time"

	"github.com/superwhys/snooker-assistant-server/domain/user"
	"gorm.io/gorm"
)

type UserPo struct {
	ID int `gorm:"primarykey"`

	Name     string `gorm:"type:varchar(20);not null"`
	Password string `gorm:"type:varchar(200);not null"`

	// WechatLogin
	WechatId string `gorm:"type:varchar(100)"`

	Email  string `gorm:"type:varchar(100)"`
	Phone  string `gorm:"type:varchar(100)"`
	Avatar string `gorm:"type:varchar(255)"`
	Status user.Status
	Role   user.Role
	Rooms  []*RoomPo `gorm:"many2many:user_rooms;"`

	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLoginAt time.Time      `gorm:"default:null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (u *UserPo) TableName() string {
	return "users"
}

func (u *UserPo) FromEntity(ue *user.User) *UserPo {
	u.ID = ue.UserId
	u.Name = ue.Name
	u.Password = ue.Password
	u.WechatId = ue.WechatId
	u.Status = ue.Status
	u.Role = ue.Role
	u.LastLoginAt = ue.LastLoginAt
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
		UserId:      u.ID,
		Name:        u.Name,
		Password:    u.Password,
		WechatId:    u.WechatId,
		Status:      u.Status,
		Role:        u.Role,
		Rooms:       rooms,
		LastLoginAt: u.LastLoginAt,
		UserInfo: &user.BaseInfo{
			Email:  u.Email,
			Phone:  u.Phone,
			Avatar: u.Avatar,
		},
	}
}
