// File:		events.go
// Created by:	Hoven
// Created on:	2024-11-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package record

import "gitea.hoven.com/billiard/billiard-assistant-server/pkg/events"

type ActionEvent struct {
	RoomId int
	UserId int
	Action Action
}

func NewActionEvent(roomId, userId int, action Action) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.RecordAction,
		Payload: &ActionEvent{
			RoomId: roomId,
			UserId: userId,
			Action: action,
		},
	}
}
