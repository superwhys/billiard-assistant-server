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
	GetOwnerRoomCount(ctx context.Context, userId int) (int64, error)
	GetUserGameRooms(ctx context.Context, userId int, justOwner bool) ([]*Room, error)
	AddUserToRoom(ctx context.Context, roomId, userId int, userName string, isVirtual bool) error
	RemoveUserFromRoom(ctx context.Context, roomId, userId int, userName string, isVirtual bool) error
}
