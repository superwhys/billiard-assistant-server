// File:		game.go
// Created by:	Hoven
// Created on:	2025-01-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"context"
	"mime/multipart"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitea.hoven.com/billiard/billiard-assistant-server/server/dto"
	"github.com/go-puzzles/puzzles/plog"
	"gorm.io/datatypes"
)

func (s *BilliardServer) GetGameList(ctx context.Context) ([]*dto.Game, error) {
	gameList, err := s.GameSrv.GetGameList(ctx)
	if err != nil {
		plog.Errorc(ctx, "get game list error: %v", err)
		return nil, err
	}

	ret := make([]*dto.Game, 0, len(gameList))
	for _, g := range gameList {
		ret = append(ret, dto.GameEntityToDto(g))
	}

	return ret, nil
}

func (s *BilliardServer) CreateGame(ctx context.Context, req *dto.CreateGameRequest) (*dto.Game, error) {

	g := &game.Game{
		GameType: shared.BilliardGameType(req.GameType),
		Icon:     req.IconUrl,
		GameConfig: &game.Config{
			MaxPlayers: req.MaxPlayers,
			Desc:       req.Desc,
		},
	}
	err := s.GameSrv.CreateGame(ctx, g)
	if err != nil {
		plog.Errorc(ctx, "create game error: %v", err)
		return nil, err
	}

	return dto.GameEntityToDto(g), nil
}

func (s *BilliardServer) UpdateGame(ctx context.Context, req *dto.UpdateGameRequest) error {
	g := &game.Game{
		GameId: req.GameId,
		Icon:   req.IconUrl,
		GameConfig: &game.Config{
			MaxPlayers: req.MaxPlayers,
			Desc:       req.Desc,
		},
	}
	err := s.GameSrv.UpdateGame(ctx, g)
	if err != nil {
		plog.Errorc(ctx, "update game error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) DeleteGame(ctx context.Context, gameId int) error {
	err := s.GameSrv.DeleteGame(ctx, gameId)
	if err != nil {
		plog.Errorc(ctx, "delete game error: %v", err)
		return err
	}

	return err
}

func (s *BilliardServer) UploadGameIcon(ctx context.Context, fh *multipart.FileHeader) (string, error) {
	iconUrl, err := s.GameSrv.UploadGameIcon(ctx, fh)
	if err != nil {
		plog.Errorc(ctx, "upload gameIcon error: %v", err)
		return "", exception.ErrUploadGameIcon
	}

	return iconUrl, nil
}

func (s *BilliardServer) StartGame(ctx context.Context, userId, roomId int, extra datatypes.JSONMap) error {
	currentGame, err := s.RoomSrv.StartGame(ctx, userId, roomId, extra)
	if err != nil {
		plog.Errorc(ctx, "start game error: %v", err)
		return err
	}

	s.EventBus.Publish(room.NewGameStartEvent(roomId, userId, currentGame, extra))
	return nil
}

func (s *BilliardServer) EndGame(ctx context.Context, userId, roomId int) error {
	err := s.RoomSrv.EndGame(ctx, userId, roomId)
	if err != nil {
		plog.Errorc(ctx, "end game error: %v", err)
		return err
	}

	s.EventBus.Publish(room.NewGameEndEvent(roomId, userId))
	return nil
}
