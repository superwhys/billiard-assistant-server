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

func (b *BaseInfo) HasUpdate(other *BaseInfo) bool {
	if other.Email != "" && other.Email != b.Email {
		return true
	}

	if other.Phone != "" && other.Phone != b.Phone {
		return true
	}

	if other.Avatar != "" && other.Avatar != b.Avatar {
		return true
	}

	return false
}

type User struct {
	UserId   int
	Name     string
	UserInfo *BaseInfo

	Status Status
	Role   Role

	Rooms []Room
}

// HasUpdate is used to determine whether certain options have been updated
// only Name, Status, Role, UserInfo support
func (u *User) HasUpdate(other *User) bool {
	if other.Name != "" && other.Name != u.Name {
		return true
	}

	if other.Status != 0 && other.Status != u.Status {
		return true
	}

	if other.Role != 0 && other.Role != u.Role {
		return true
	}

	if u.UserInfo.HasUpdate(other.UserInfo) {
		return true
	}

	return false
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
