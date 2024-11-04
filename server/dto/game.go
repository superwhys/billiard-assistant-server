// File:		game.go
// Created by:	Hoven
// Created on:	2024-10-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import (
	"github.com/superwhys/snooker-assistant-server/domain/game"
)

type Game struct {
	GameId     int
	MaxPlayers int
	GameType   string
	Desc       string
}

func GameEntityToDto(g *game.Game) *Game {
	game := &Game{
		GameId:   g.GameId,
		GameType: g.GameType.String(),
	}
	
	if g.GameConfig != nil {
		game.MaxPlayers = g.GameConfig.MaxPlayers
		game.Desc = g.GameConfig.Desc
	}
	
	return game
}

type GetGameListResp struct {
	Games []*Game `json:"games"`
}

type CreateGameRequest struct {
	MaxPlayers int             `json:"max_players" binding:"required"`
	GameType   game.SaGameType `json:"game_type" binding:"required"`
	Desc       string          `json:"desc"`
}

type DeleteGameRequest struct {
	GameId int `uri:"gameId"`
}
