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

type User struct {
	UserId   int       `json:"user_id"`
	Name     string    `json:"name"`
	Gender   Gender    `json:"gender"`
	UserInfo *BaseInfo `json:"user_info"`

	Status Status `json:"status"`
	Role   Role   `json:"role"`

	Rooms []Room `json:"rooms"`
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
