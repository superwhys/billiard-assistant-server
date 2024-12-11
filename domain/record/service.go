// File:		service.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package record

import (
	"context"
	"encoding/json"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

type IRecordService interface {
	HandleRecord(ctx context.Context, roomId int, gameType shared.BilliardGameType, rawRecord json.RawMessage) ([]RecordItem, error)
	HandleAction(ctx context.Context, gameType shared.BilliardGameType, rawAction json.RawMessage) (Action, error)
	GetCurrentRecord(ctx context.Context, roomId int, gameType shared.BilliardGameType) (*Record, error)
	GetRoomActions(ctx context.Context, roomId int) ([]Action, error)
}
