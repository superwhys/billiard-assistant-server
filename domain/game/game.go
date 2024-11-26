package game

import (
	"fmt"

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

type IGameStrategy interface {
	SetupGame(g shared.BaseGameConfig) []any
	HandleAction()
}

// TODO: may be delete
var strategyFactory = make(map[shared.BilliardGameType]IGameStrategy)

func RegisterStrategy(gt shared.BilliardGameType, strategy IGameStrategy) {
	strategyFactory[gt] = strategy
}

func NewGameStrategy(gameType shared.BilliardGameType) (IGameStrategy, error) {
	strategy, ok := strategyFactory[gameType]
	if !ok {
		return nil, fmt.Errorf("No strategy registered for the given game type: %v", gameType)
	}
	return strategy, nil
}
