package game

import (
	"context"
	"mime/multipart"
)

type IGameService interface {
	CreateGame(ctx context.Context, g *Game) error
	DeleteGame(ctx context.Context, gameId int) error
	UpdateGame(ctx context.Context, g *Game) error
	GetGameList(ctx context.Context) ([]*Game, error)
	UploadGameIcon(ctx context.Context, fh *multipart.FileHeader) (string, error)
}
