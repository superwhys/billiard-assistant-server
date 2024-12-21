// File:		service.go
// Created by:	Hoven
// Created on:	2024-12-21
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package snooker

import "gitea.hoven.com/billiard/billiard-assistant-server/domain/game"

type ISnookerService interface {
	game.IGameStrategy
}
