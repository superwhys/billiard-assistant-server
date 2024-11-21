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
)

type GameRoom struct {
	RoomId int `json:"room_id"`

	GameId     int    `json:"game_id"`
	GameIcon   string `json:"game_icon"`
	GameType   string `json:"game_type"`
	GameStatus string `json:"game_status"`

	OwnerId       int       `json:"owner_id"`
	Players       []int     `json:"players"`
	WinLoseStatus string    `json:"win_lose_status"`
	CreateAt      time.Time `json:"create_at"`
}

func GameRoomEntityToDto(gr *room.Room) *GameRoom {
	gameRoom := &GameRoom{
		RoomId:        gr.RoomId,
		GameId:        gr.GameId,
		OwnerId:       gr.OwnerId,
		Players:       gr.PlayerIds(),
		GameStatus:    gr.GameStatus.String(),
		WinLoseStatus: gr.WinLoseStatus.String(),
		CreateAt:      gr.CreateAt,
	}

	if gr.Game != nil {
		gameRoom.GameType = gr.Game.GetGameType().String()
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

type UpdateGameRoomRequest struct {
	RoomId        int                `json:"room_id"`
	GameStatus    room.Status        `json:"game_status"`
	WinLoseStatus room.WinLoseStatus `json:"win_lose_status"`
}

type DeleteGameRoomRequest struct {
	RoomId int `json:"game_id"`
}

type EnterGameRoomRequest struct {
	RoomId int `json:"room_id"`
}

type LeaveGameRoomRequest struct {
	RoomId int `json:"room_id"`
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
