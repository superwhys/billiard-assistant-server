package game

import (
	"context"
)

type IGameRepo interface {
	CreateGame(ctx context.Context, g *Game) error
	DeleteGame(ctx context.Context, gameId int) error
	UpdateGame(ctx context.Context, g *Game) error
	GetGameById(ctx context.Context, gameId int) (*Game, error)
	GetGameList(ctx context.Context) ([]*Game, error)
}
