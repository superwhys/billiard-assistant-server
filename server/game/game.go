// File:		game.go
// Created by:	Hoven
// Created on:	2024-10-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package gameSrv

import (
	"context"
	"mime/multipart"

	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/oss"
)

var _ game.IGameService = (*GameService)(nil)

type GameService struct {
	gameRepo game.IGameRepo
	oss      oss.IOSS
}

func NewGameService(gameRepo game.IGameRepo, oss oss.IOSS) *GameService {
	return &GameService{
		gameRepo: gameRepo,
		oss:      oss,
	}
}

func (gs *GameService) CreateGame(ctx context.Context, g *game.Game) error {
	return gs.gameRepo.CreateGame(ctx, g)
}

func (gs *GameService) DeleteGame(ctx context.Context, gameId int) error {
	return gs.gameRepo.DeleteGame(ctx, gameId)
}

func (gs *GameService) GetGameList(ctx context.Context) ([]*game.Game, error) {
	return gs.gameRepo.GetGameList(ctx)
}

func (gs *GameService) UploadGameIcon(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	avatarUrl, err := gs.oss.UploadFile(ctx, file.Size, oss.SourceGameIcon, file.Filename, src)
	if err != nil {
		return "", errors.Wrap(err, "uploadAvatar")
	}

	return avatarUrl, nil
}
