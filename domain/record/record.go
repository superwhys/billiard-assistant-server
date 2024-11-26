// File:		record.go
// Created by:	Hoven
// Created on:	2024-11-24
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package record

import "time"

type RecordItem interface {
	UnmarshalFrom(jsonStr string) error
	GetRecordTime() time.Time
	GetRecordUser() (userId int)
}

type Record struct {
	ID            int
	RoomId        int
	Histories     []RecordItem
	CurrentRecord RecordItem
}

type IRecordService interface {
	HandleHistoryRecord(history RecordItem) error
}
