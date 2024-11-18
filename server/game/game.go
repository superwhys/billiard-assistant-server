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
	
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
)

var _ game.IGameService = (*GameService)(nil)

type GameService struct {
	gameRepo game.IGameRepo
	userRepo user.IUserRepo
}

func NewGameService(gameRepo game.IGameRepo, userRepo user.IUserRepo) *GameService {
	return &GameService{
		gameRepo: gameRepo,
		userRepo: userRepo,
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
