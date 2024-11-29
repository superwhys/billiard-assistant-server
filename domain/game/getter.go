// File:		getter.go
// Created by:	Hoven
// Created on:	2024-11-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package game

import (
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

var _ shared.BaseGame = (*Game)(nil)

func (c *Game) GetGameId() int {
	return c.GameId
}

func (c *Game) GetMaxPlayers() int {
	return c.GameConfig.MaxPlayers
}

func (c *Game) GetGameType() int {
	return int(c.GameType)
}

func (c *Game) GetIcon() string {
	return c.Icon
}

func (c *Game) GetIsActivate() bool {
	return c.IsActivate
}

func (c *Game) GetDesc() string {
	return c.GameConfig.Desc
}

func (c *Game) GetGameConfig() shared.BaseGameConfig {
	return c.GameConfig
}
