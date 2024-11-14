package shared

type BaseGameConfig interface {
	GetMaxPlayers() int
}

type BaseGame interface {
	GetGameId() int
	GetGameConfig() BaseGameConfig
}

type BaseUser interface {
	GetUserId() int
	GetName() string
}

type BaseRoom interface {
	GetRoomId() int
	PlayerIds() []int
}
