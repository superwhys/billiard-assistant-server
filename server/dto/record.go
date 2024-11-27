// File:		record.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import (
	"encoding/json"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/record"
)

type Record struct {
	RoomId        int               `json:"room_id"`
	Histories     []record.Action   `json:"histories"`
	CurrentRecord record.RecordItem `json:"current_record"`
}

type Action struct {
	RoomId  int             `json:"room_id"`
	Actions []record.Action `json:"actions"`
}

type RoomUriRequest struct {
	RoomId int `uri:"roomId" binding:"required"`
}

type RoomActionRequest struct {
	RoomUriRequest
	Action json.RawMessage `json:"action"`
}

type RoomRecordRequest struct {
	RoomUriRequest
	Record json.RawMessage `json:"record"`
}

func RecordEntityToDto(r *record.Record) *Record {
	return &Record{
		RoomId:        r.RoomId,
		Histories:     r.Histories,
		CurrentRecord: r.CurrentRecord,
	}
}
