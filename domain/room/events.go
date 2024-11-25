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
	UserId      int
	VirtualName string
	RoomId      int
}

func NewEnterRoomEvent(virtualName string, userId, roomId int) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.PlayerJoined,
		Payload:   &EnterRoomEvent{VirtualName: virtualName, UserId: userId, RoomId: roomId},
	}
}

type LeaveRoomEvent struct {
	UserId      int
	VirtualName string
	RoomId      int
}

func NewLeaveRoomEvent(virtualName string, userId, roomId int) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.PlayerLeft,
		Payload:   &LeaveRoomEvent{VirtualName: virtualName, UserId: userId, RoomId: roomId},
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
