// File:		events.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package room

import (
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/events"
)

type EnterRoomEvent struct {
	UserId    int
	UserName  string
	IsVirtual bool
	RoomId    int
}

func NewEnterRoomEvent(roomId, userId int, userName string, isVirtual bool) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.PlayerJoined,
		Payload: &EnterRoomEvent{
			UserId:    userId,
			UserName:  userName,
			IsVirtual: isVirtual,
			RoomId:    roomId,
		},
	}
}

type LeaveRoomEvent struct {
	UserId    int
	UserName  string
	IsVirtual bool
	RoomId    int
}

func NewLeaveRoomEvent(roomId, userId int, userName string, isVirtual bool) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.PlayerLeft,
		Payload: &LeaveRoomEvent{
			UserId:    userId,
			UserName:  userName,
			IsVirtual: isVirtual,
			RoomId:    roomId,
		},
	}
}

type GameStartEvent struct {
	RoomId int
	Game   Game
}

func NewGameStartEvent(roomId int, game Game) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.GameStart,
		Payload: &GameStartEvent{
			RoomId: roomId,
			Game:   game,
		},
	}
}

type GameEndEvent struct {
	RoomId int
}

func NewGameEndEvent(roomId int) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.GameEnd,
		Payload: &GameEndEvent{
			RoomId: roomId,
		},
	}
}
