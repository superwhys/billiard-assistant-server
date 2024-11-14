// File:		events.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package room

import (
	"github.com/superwhys/snooker-assistant-server/pkg/events"
)

type EnterRoomEvent struct {
	RoomId int
	User   User
}

func NewEnterRoomEvent(roomId int, user User) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.PlayerJoined,
		Payload:   &EnterRoomEvent{RoomId: roomId, User: user},
	}
}

type LeaveRoomEvent struct {
	UserId int
	RoomId int
}

func NewLeaveRoomEvent(userId, roomId int) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.PlayerLeft,
		Payload:   &LeaveRoomEvent{UserId: userId, RoomId: roomId},
	}
}

type PrepareEvent struct {
	UserId int
	RoomId int
}

func NewPrepareEvent(userId, roomId int) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.PlayerPrepare,
		Payload:   &PrepareEvent{UserId: userId, RoomId: roomId},
	}
}

type GameStartEvent struct {
	roomId  int
	payload any
}

func NewGameStartEvent(roomId int, payload any) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.GameStart,
		Payload:   &GameStartEvent{roomId, payload},
	}
}
