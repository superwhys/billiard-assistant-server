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

type UserAuthPo struct {
	ID int `gorm:"primaryKey"`
	
	UserPoID   int `gorm:"index"`
	AuthType   user.AuthType
	Identifier string `gorm:"size:255;uniqueIndex:idx_auth_type_identifier"`
	Credential string `gorm:"size:255"`
	
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ua *UserAuthPo) TableName() string {
	return "user_auths"
}

func (u *UserPo) FromEntity(ue *user.User) *UserPo {
	var uas []*UserAuthPo
	for _, ua := range ue.UserAuths {
		uap := new(UserAuthPo)
		uap.FromEntity(ua)
		
		uas = append(uas, uap)
	}
	
	u.ID = ue.UserId
	u.Name = ue.Name
	u.UserAuthPos = uas
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
	
	var uas []*user.UserAuth
	for _, ua := range u.UserAuthPos {
		uas = append(uas, ua.ToEntity())
	}
	
	return &user.User{
		UserId:    u.ID,
		Name:      u.Name,
		Status:    u.Status,
		Role:      u.Role,
		Rooms:     rooms,
		UserAuths: uas,
		UserInfo: &user.BaseInfo{
			Email:  u.Email,
			Phone:  u.Phone,
			Avatar: u.Avatar,
		},
	}
}

func (ua *UserAuthPo) FromEntity(entity *user.UserAuth) {
	ua.ID = entity.Id
	ua.UserPoID = entity.UserId
	ua.AuthType = entity.AuthType
	ua.Identifier = entity.Identifier
	ua.Credential = entity.Credential
}

func (ua *UserAuthPo) ToEntity() *user.UserAuth {
	return &user.UserAuth{
		Id:         ua.ID,
		UserId:     ua.UserPoID,
		AuthType:   ua.AuthType,
		Identifier: ua.Identifier,
		Credential: ua.Credential,
	}
}
