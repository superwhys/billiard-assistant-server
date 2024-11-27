package game

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

type Game struct {
	GameId     int
	GameType   shared.BilliardGameType
	Icon       string
	IsActivate bool
	GameConfig *Config
}

func (c *Game) GetGameId() int {
	return c.GameId
}

func (c *Game) GetGameConfig() shared.BaseGameConfig {
	return c.GameConfig
}

func (c *Game) GetMaxPlayers() int {
	return c.GameConfig.MaxPlayers
}

func (c *Game) GetGameType() shared.BilliardGameType {
	return c.GameType
}

func (c *Game) GetIcon() string {
	return c.Icon
}

type Config struct {
	MaxPlayers int
	Desc       string
}

func (c *Config) GetMaxPlayers() int {
	return c.MaxPlayers
}

type Action interface {
	GetActionRoomId() (roomId int)
	GetActionUser() (userId int)
	GetActionTime() time.Time
}

type Record interface {
	GetRecordRoomId() (roomId int)
}

type IGameStrategy interface {
	GetGameType() shared.BilliardGameType
	SetupGame(g shared.BaseGameConfig) []any
	UnmarshalAction(action json.RawMessage) (Action, error)
	UnmarshalRecord(record json.RawMessage) (Record, error)
	HandleAction(ctx context.Context, action Action) error
	GetRoomActions(ctx context.Context, roomId int) ([]Action, error)
}
