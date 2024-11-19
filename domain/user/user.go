package user

type Room interface {
	GetRoomId() int
}

type BaseInfo struct {
	Email    string
	Phone    string
	Avatar   string
	Password string
}

type User struct {
	UserId   int
	Name     string
	UserInfo *BaseInfo

	Status Status
	Role   Role

	Rooms []Room
}

func (u *User) GetUserId() int {
	return u.UserId
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetAvatar() string {
	if u.UserInfo == nil {
		return ""
	}

	return u.UserInfo.Avatar
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
