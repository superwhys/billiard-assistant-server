// File:		service.go
// Created by:	Hoven
// Created on:	2024-11-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package room

import (
	"context"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
)

type IRoomService interface {
	CreateGameRoom(ctx context.Context, u *user.User, gameId int) (*Room, error)
	DeleteGameRoom(ctx context.Context, roomId, userId int) error
	UpdateGameRoomStatus(ctx context.Context, room *Room) error
	UpdateRoomUserHeartbeart(ctx context.Context, roomId, userId int) error
	GetUserGameRooms(ctx context.Context, userId int) ([]*Room, error)
	GetRoomGameType(ctx context.Context, roomId int) (shared.BilliardGameType, error)
	GetRoomById(ctx context.Context, roomId int) (*Room, error)
	GetRoomByCode(ctx context.Context, roomCode string) (*Room, error)
	EnterGameRoom(ctx context.Context, roomId int, enterUser shared.BaseUser, isVirtual bool) error
	QuitGameRoom(ctx context.Context, roomId int, quitUser shared.BaseUser, isVirtual bool) error
	StartGame(ctx context.Context, userId, roomId int) (shared.BaseGame, error)
	EndGame(ctx context.Context, userId, roomId int) error
}
