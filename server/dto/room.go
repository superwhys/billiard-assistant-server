// File:		room.go
// Created by:	Hoven
// Created on:	2024-11-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import (
	"time"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

type GameRoom struct {
	RoomId   int    `json:"room_id,omitempty"`
	RoomCode string `json:"room_code,omitempty"`

	GameId   int    `json:"game_id,omitempty"`
	GameIcon string `json:"game_icon,omitempty"`

	MaxPlayer  int                     `json:"max_player,omitempty"`
	GameType   shared.BilliardGameType `json:"game_type,omitempty"`
	GameStatus room.Status             `json:"game_status,omitempty"`

	OwnerId       int                `json:"owner_id,omitempty"`
	Players       []*room.RoomPlayer `json:"players"`
	WinLoseStatus string             `json:"win_lose_status,omitempty"`
	CreateAt      time.Time          `json:"create_at,omitempty"`
}

func GameRoomEntityToDto(gr *room.Room) *GameRoom {
	gameRoom := &GameRoom{
		RoomId:        gr.RoomId,
		RoomCode:      gr.RoomCode,
		GameId:        gr.GameId,
		OwnerId:       gr.OwnerId,
		Players:       gr.Players,
		GameStatus:    gr.GameStatus,
		WinLoseStatus: gr.WinLoseStatus.String(),
		CreateAt:      gr.CreateAt,
	}

	if gr.Game != nil {
		gameRoom.MaxPlayer = gr.Game.GetMaxPlayers()
		gameRoom.GameType = gr.Game.GetGameType()
		gameRoom.GameIcon = gr.Game.GetIcon()
	}

	return gameRoom
}

type CreateGameRoomRequest struct {
	GameId int `json:"game_id"`
}

type GetRoomRequest struct {
	RoomId int `uri:"roomId"`
}

type GetRoomByCodeRequest struct {
	RoomCode string `uri:"roomCode"`
}

type UpdateGameRoomRequest struct {
	RoomId        int                `json:"room_id"`
	GameStatus    room.Status        `json:"game_status"`
	WinLoseStatus room.WinLoseStatus `json:"win_lose_status"`
}

type DeleteGameRoomRequest struct {
	RoomId int `json:"room_id"`
}

type EnterGameRoomRequest struct {
	RoomId    int    `json:"room_id"`
	UserName  string `json:"user_name"`
	IsVirtual bool   `json:"is_virtual"`
}

type LeaveGameRoomRequest struct {
	RoomId    int    `json:"room_id"`
	UserName  string `json:"user_name"`
	IsVirtual bool   `json:"is_virtual"`
}

type GetUserGameRoomsResp struct {
	Rooms []*GameRoom `json:"rooms"`
}

type PrepareGameRequest struct {
	RoomId int `json:"room_id"`
}

type StartGameRequest struct {
	RoomId int `json:"room_id"`
}
