// File:		record.go
// Created by:	Hoven
// Created on:	2024-11-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package record

import (
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

type Action interface {
	shared.Action
}

type RecordItem interface {
	shared.RecordItem
}

type Record struct {
	ID            int
	RoomId        int
	Histories     []Action
	CurrentRecord RecordItem
}
