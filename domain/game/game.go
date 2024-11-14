package game

import (
	"fmt"

	"github.com/superwhys/snooker-assistant-server/domain/shared"
)

type Game struct {
	GameId     int
	GameType   shared.SaGameType
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

func (c *Game) GetGameType() shared.SaGameType {
	return c.GameType
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

var strategyFactory = make(map[shared.SaGameType]IGameStrategy)

func RegisterStrategy(gt shared.SaGameType, strategy IGameStrategy) {
	strategyFactory[gt] = strategy
}

func NewGameStrategy(gameType shared.SaGameType) (IGameStrategy, error) {
	strategy, ok := strategyFactory[gameType]
	if !ok {
		return nil, fmt.Errorf("No strategy registered for the given game type: %v", gameType)
	}
	return strategy, nil
}
