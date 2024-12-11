// File:		game.go
// Created by:	Hoven
// Created on:	2024-10-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import "gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"

type Game struct {
	GameId     int    `json:"game_id"`
	MaxPlayers int    `json:"max_players"`
	GameType   int    `json:"game_type"`
	Icon       string `json:"icon"`
	IsActivate bool   `json:"is_activate"`
	Desc       string `json:"desc"`
}

func GameEntityToDto(g shared.BaseGame) *Game {
	if g == nil {
		return nil
	}

	game := &Game{
		GameId:     g.GetGameId(),
		GameType:   g.GetGameType(),
		Icon:       g.GetIcon(),
		IsActivate: g.GetIsActivate(),
		MaxPlayers: g.GetMaxPlayers(),
		Desc:       g.GetDesc(),
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
