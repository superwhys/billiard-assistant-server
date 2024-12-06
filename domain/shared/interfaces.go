package shared

import (
	"time"
)

type BaseGameConfig interface {
	GetMaxPlayers() int
}

type BaseGame interface {
	GetGameId() int
	GetMaxPlayers() int
	GetGameType() int
	GetIcon() string
	GetIsActivate() bool
	GetDesc() string
	GetGameConfig() BaseGameConfig
}

type BaseUser interface {
	GetUserId() int
	GetName() string
	GetEmail() string
	GetPhone() string
	GetAvatar() string
	GetGender() string
	GetStatus() int
	GetRole() int
	IsAdmin() bool
}

type BaseRoom interface {
	GetRoomId() int
	GetRoomCode() string
	GetOwner() BaseUser
	GetRecord() BaseRecord
	GetRoomPlayers() []RoomPlayer
	GetGameStatus() int
	GetExtra() map[string]any
	GetWinLoseStatus() string
	GetCreateAt() time.Time
	GetGame() BaseGame
}

type RoomPlayer interface {
	GetRoomId() int
	GetUserId() int
	GetUserName() string
	GetIsVirtual() bool
	GetHeartbeatAt() time.Time
}

type Action interface {
	GetActionRoomId() (roomId int)
	GetActionUser() (userId int)
	GetActionTime() time.Time
}

type RecordItem interface {
	GetRecordRoomId() (roomId int)
}

type BaseRecord interface {
	GetRecordId() int
	GetRoomId() int
	GetCurrentRecord() []RecordItem
	GetActions() []Action
}
