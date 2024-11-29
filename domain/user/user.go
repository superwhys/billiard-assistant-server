package user

type Room interface {
	GetRoomId() int
}

type BaseInfo struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
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
	UserId   int       `json:"user_id"`
	Name     string    `json:"name"`
	Gender   Gender    `json:"gender"`
	UserInfo *BaseInfo `json:"user_info"`

	Status Status `json:"status"`
	Role   Role   `json:"role"`

	Rooms []Room `json:"rooms"`
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

	if other.Gender != 0 && other.Gender != u.Gender {
		return true
	}

	if u.UserInfo == nil && other.UserInfo != nil {
		return true
	}

	if u.UserInfo != nil && other.UserInfo != nil && u.UserInfo.HasUpdate(other.UserInfo) {
		return true
	}

	return false
}

type Gender int

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
)

func (g Gender) String() string {
	switch g {
	case GenderUnknown:
		return "未知"
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "未知"
	}
}

func (g Gender) Parse(str string) Gender {
	switch str {
	case "未知":
		return GenderUnknown
	case "男":
		return GenderMale
	case "女":
		return GenderFemale
	default:
		return GenderUnknown
	}
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
	RolePro
)

func (r Role) IsPro() bool {
	return r == RolePro
}
