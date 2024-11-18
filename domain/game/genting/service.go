// File:		service.go
// Created by:	Hoven
// Created on:	2024-11-14
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package genting

import "gitlab.hoven.com/billiard/billiard-assistant-server/domain/game"

type IGentingService interface {
	game.IGameStrategy
}
