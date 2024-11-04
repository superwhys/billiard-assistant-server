package game

import (
	"context"
)

type IGameService interface {
	CreateGame(ctx context.Context, g *Game) error
	DeleteGame(ctx context.Context, gameId int) error
	GetGameList(ctx context.Context) ([]*Game, error)
}
