// File:		repo.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package record

import (
	"context"
	"reflect"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

type IRecordRepo interface {
	GetRecordByRoomId(ctx context.Context, roomId int, recordTmpl reflect.Type) (*Record, error)
	UpdateRoomRecord(ctx context.Context, gameType shared.BilliardGameType, record RecordItem) error
}
