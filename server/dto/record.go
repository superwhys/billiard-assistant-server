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
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

type Record struct {
	RoomId        int                 `json:"room_id"`
	Histories     []shared.Action     `json:"histories"`
	CurrentRecord []shared.RecordItem `json:"current_records"`
}

func RecordEntityToDto(r shared.BaseRecord) *Record {
	if r == nil {
		return nil
	}

	return &Record{
		RoomId:        r.GetRoomId(),
		Histories:     r.GetActions(),
		CurrentRecord: r.GetCurrentRecord(),
	}
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
	Records json.RawMessage `json:"records"`
}
