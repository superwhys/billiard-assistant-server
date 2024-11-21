// File:		game.go
// Created by:	Hoven
// Created on:	2024-10-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import (
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/game"
)

type Game struct {
	GameId     int
	MaxPlayers int
	GameType   string
	IsActivate bool
	Desc       string
}

func GameEntityToDto(g *game.Game) *Game {
	game := &Game{
		GameId:     g.GameId,
		GameType:   g.GameType.String(),
		IsActivate: g.IsActivate,
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

type UpdateGameRequest struct {
	GameId     int    `uri:"gameId" binding:"gt=0"`
	MaxPlayers int    `json:"max_players"`
	IsActivate bool   `json:"is_activate"`
	IconUrl    string `json:"icon_url"`
	Desc       string `json:"desc"`
}

type CreateGameRequest struct {
	MaxPlayers int    `json:"max_players" binding:"required"`
	GameType   int    `json:"game_type" binding:"required"`
	IsActivate bool   `json:"is_activate"`
	IconUrl    string `json:"icon_url"`
	Desc       string `json:"desc"`
}

type DeleteGameRequest struct {
	GameId int `uri:"gameId"`
}
