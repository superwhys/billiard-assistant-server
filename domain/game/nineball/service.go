// File:		service.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package nineball

import "gitlab.hoven.com/billiard/billiard-assistant-server/domain/game"

type INineballService struct {
	game.IGameStrategy
}
