// File:		getter.go
// Created by:	Hoven
// Created on:	2024-11-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package room

import (
	"time"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

var _ shared.BaseRoom = (*Room)(nil)

func (r *Room) GetRoomId() int {
	return r.RoomId
}

func (r *Room) GetRoomCode() string {
	return r.RoomCode
}

func (r *Room) GetOwner() shared.BaseUser {
	if r.Owner == nil {
		return nil
	}

	return r.Owner
}

func (r *Room) GetRecord() shared.BaseRecord {
	if r.Record == nil {
		return nil
	}
	return r.Record
}

func (r *Room) GetRoomPlayers() []shared.RoomPlayer {
	return r.Players
}

func (r *Room) GetGameStatus() int {
	return int(r.GameStatus)
}

func (r *Room) GetWinLoseStatus() string {
	return r.WinLoseStatus.String()
}

func (r *Room) GetCreateAt() time.Time {
	return r.CreateAt
}

func (r *Room) GetGame() shared.BaseGame {
	if r.Game == nil {
		return nil
	}
	return r.Game
}

var _ shared.RoomPlayer = (*RoomPlayer)(nil)

func (r *RoomPlayer) GetRoomId() int {
	return r.RoomId
}

func (r *RoomPlayer) GetUserId() int {
	return r.UserId
}

func (r *RoomPlayer) GetUserName() string {
	return r.UserName
}

func (r *RoomPlayer) GetIsVirtual() bool {
	return r.IsVirtualPlayer
}

func (r *RoomPlayer) GetHeartbeatAt() time.Time {
	return r.HeartbeatAt
}
