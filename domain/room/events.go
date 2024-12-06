// File:		events.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package room

import (
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
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
	UserId int
	Game   shared.BaseGame
	Extra  any
}

func NewGameStartEvent(roomId, userId int, game shared.BaseGame, extra any) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.GameStart,
		Payload: &GameStartEvent{
			RoomId: roomId,
			UserId: userId,
			Game:   game,
			Extra:  extra,
		},
	}
}

type GameEndEvent struct {
	RoomId int
	UserId int
}

func NewGameEndEvent(roomId, userId int) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.GameEnd,
		Payload: &GameEndEvent{
			RoomId: roomId,
			UserId: userId,
		},
	}
}
