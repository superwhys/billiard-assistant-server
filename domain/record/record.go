// File:		record.go
// Created by:	Hoven
// Created on:	2024-11-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package record

import (
	"time"
)

type Action interface {
	GetActionRoomId() (roomId int)
	GetActionUser() (userId int)
	GetActionTime() time.Time
}

type RecordItem interface {
	GetRecordRoomId() (roomId int)
}

type Record struct {
	ID            int
	RoomId        int
	Histories     []Action
	CurrentRecord RecordItem
}
