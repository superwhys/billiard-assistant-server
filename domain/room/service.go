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

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
)

type IRoomService interface {
	CreateGameRoom(ctx context.Context, u *user.User, gameId int) (*Room, error)
	DeleteGameRoom(ctx context.Context, roomId, userId int) error
	UpdateGameRoomStatus(ctx context.Context, room *Room) error
	GetUserGameRooms(ctx context.Context, userId int, justOwner bool) ([]*Room, error)
	GetRoomById(ctx context.Context, roomId int) (*Room, error)
	EnterGameRoom(ctx context.Context, virtualName string, userId, roomId int) error
	QuitGameRoom(ctx context.Context, virtualName string, userId, roomId int) error
	StartGame(ctx context.Context, userId, roomId int) (Game, error)
}
