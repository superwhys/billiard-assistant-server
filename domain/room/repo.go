// File:		repo.go
// Created by:	Hoven
// Created on:	2024-11-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package room

import (
	"context"
)

type IRoomRepo interface {
	CreateRoom(ctx context.Context, userId, gameId int) (*Room, error)
	UpdateRoom(ctx context.Context, room *Room) error
	DeleteRoom(ctx context.Context, roomId int) error
	GetRoomById(ctx context.Context, roomId int) (*Room, error)
	GetUserGameRooms(ctx context.Context, userId int) ([]*Room, error)
	AddUserToRoom(ctx context.Context, userId, roomId int) (User, error)
	RemoveUserFromRoom(ctx context.Context, userId, roomId int) (User, error)
	UpdatePlayerPrepared(ctx context.Context, userId, roomId int, prepared bool) error
}
