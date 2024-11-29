// File:		getter.go
// Created by:	Hoven
// Created on:	2024-11-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package record

import (
	"github.com/go-puzzles/puzzles/putils"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

var _ shared.BaseRecord = (*Record)(nil)

func (r *Record) GetRecordId() int {
	return r.ID
}

func (r *Record) GetRoomId() int {
	return r.RoomId
}

func (r *Record) GetCurrentRecord() []shared.RecordItem {
	if r.CurrentRecord == nil {
		return nil
	}

	return putils.Convert(r.CurrentRecord, func(r RecordItem) shared.RecordItem {
		return r
	})
}

func (r *Record) GetActions() []shared.Action {
	return putils.Convert(r.Histories, func(a Action) shared.Action {
		return a
	})
}
