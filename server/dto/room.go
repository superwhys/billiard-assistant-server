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

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gorm.io/datatypes"
)

type RoomPlayer struct {
	RoomId          int       `json:"room_id,omitempty"`
	UserId          int       `json:"user_id,omitempty"`
	UserName        string    `json:"user_name,omitempty"`
	IsVirtualPlayer bool      `json:"is_virtual_player,omitempty"`
	HeartbeatAt     time.Time `json:"heartbeat_at,omitempty"`
}

type GameRoom struct {
	RoomId        int            `json:"room_id,omitempty"`
	RoomCode      string         `json:"room_code,omitempty"`
	Game          *Game          `json:"game,omitempty"`
	Owner         *User          `json:"owner,omitempty"`
	Players       []*RoomPlayer  `json:"players"`
	Record        *Record        `json:"record,omitempty"`
	Extra         map[string]any `json:"extra,omitempty"`
	GameStatus    int            `json:"game_status,omitempty"`
	WinLoseStatus string         `json:"win_lose_status,omitempty"`
	CreateAt      time.Time      `json:"create_at,omitempty"`
}

func GameRoomEntityToDto(gr shared.BaseRoom) *GameRoom {
	if gr == nil {
		return nil
	}

	gameRoom := &GameRoom{
		RoomId:        gr.GetRoomId(),
		RoomCode:      gr.GetRoomCode(),
		GameStatus:    gr.GetGameStatus(),
		Extra:         gr.GetExtra(),
		WinLoseStatus: gr.GetWinLoseStatus(),
		CreateAt:      gr.GetCreateAt(),
	}

	if gr.GetGame() != nil {
		gameRoom.Game = GameEntityToDto(gr.GetGame())
	}

	if gr.GetOwner() != nil {
		gameRoom.Owner = UserEntityToDto(gr.GetOwner())
	}

	if gr.GetRecord() != nil {
		gameRoom.Record = RecordEntityToDto(gr.GetRecord())
	}

	var players []*RoomPlayer
	for _, p := range gr.GetRoomPlayers() {
		players = append(players, &RoomPlayer{
			RoomId:          p.GetRoomId(),
			UserId:          p.GetUserId(),
			UserName:        p.GetUserName(),
			IsVirtualPlayer: p.GetIsVirtual(),
			HeartbeatAt:     p.GetHeartbeatAt(),
		})
	}
	gameRoom.Players = players

	return gameRoom
}

type CreateGameRoomRequest struct {
	GameId int `json:"game_id"`
}

type UriRoomId struct {
	RoomId int `uri:"roomId"`
}

type GetRoomRequest struct {
	UriRoomId
}

type GetRoomByCodeRequest struct {
	RoomCode string `uri:"roomCode"`
}

type UpdateGameRoomRequest struct {
	UriRoomId
	GameStatus    room.Status        `json:"game_status"`
	WinLoseStatus room.WinLoseStatus `json:"win_lose_status"`
}

type UpdateGameRoomExtraRequest struct {
	UriRoomId
	Extra map[string]any `json:"extra"`
}

type DeleteGameRoomRequest struct {
	UriRoomId
}

type EnterGameRoomRequest struct {
	UriRoomId
	VirtualUser string `json:"virtual_user"`
}

type LeaveGameRoomRequest struct {
	UriRoomId
	VirtualUser string `json:"virtual_user"`
}

type GetUserGameRoomsResp struct {
	Rooms []*GameRoom `json:"rooms"`
}

type StartGameRequest struct {
	UriRoomId
	Extra datatypes.JSONMap `json:"extra"`
}

type EndGameRequest struct {
	UriRoomId
}
