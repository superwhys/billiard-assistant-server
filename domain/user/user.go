package user

import "time"

type Room interface {
	GetRoomId() int
}

type User struct {
	UserId      int
	Name        string
	Password    string
	WechatId    string
	Status      Status
	UserInfo    *BaseInfo
	Role        Role
	Rooms       []Room
	LastLoginAt time.Time
}

type BaseInfo struct {
	Email  string
	Phone  string
	Avatar string
}

func (u *User) GetUserId() int {
	return u.UserId
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

type Status int

const (
	StatusUnknown Status = iota
	StatusActive
	StatusInactive
)

type Role int

const (
	RoleUnknown Role = iota
	RoleUser
	RoleAdmin
)
