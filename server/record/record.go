// File:		record.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package recordSrv

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-puzzles/puzzles/putils"
	"github.com/pkg/errors"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/record"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

var _ record.IRecordService = (*RecordService)(nil)

type RecordService struct {
	roomRepo            room.IRoomRepo
	recordRepo          record.IRecordRepo
	gameStrategyFactory map[shared.BilliardGameType]game.IGameStrategy
	gameRecordTmp       map[shared.BilliardGameType]reflect.Type
}

type RecordServiceOption func(rs *RecordService)

func WithGameStrategy(gameStrategies ...game.IGameStrategy) RecordServiceOption {
	return func(rs *RecordService) {
		for _, gameStrategy := range gameStrategies {
			rs.gameStrategyFactory[gameStrategy.GetGameType()] = gameStrategy
		}
	}
}

func WithGameRecordTmp(gameType shared.BilliardGameType, ri record.RecordItem) RecordServiceOption {
	return func(rs *RecordService) {
		t := reflect.TypeOf(ri)
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		rs.gameRecordTmp[gameType] = t
	}
}

func NewRecordService(
	recordRepo record.IRecordRepo,
	roomRepo room.IRoomRepo,
	opts ...RecordServiceOption,
) *RecordService {
	rs := &RecordService{
		roomRepo:            roomRepo,
		recordRepo:          recordRepo,
		gameStrategyFactory: make(map[shared.BilliardGameType]game.IGameStrategy),
		gameRecordTmp:       make(map[shared.BilliardGameType]reflect.Type),
	}

	for _, opt := range opts {
		opt(rs)
	}

	return rs
}

func (rs *RecordService) getSpecifyTmpl(gameType shared.BilliardGameType) (reflect.Type, bool) {
	if t, ok := rs.gameRecordTmp[gameType]; ok {
		return t, true
	}

	return nil, false
}

func (rs *RecordService) GetCurrentRecord(ctx context.Context, roomId int, gameType shared.BilliardGameType) (*record.Record, error) {
	recordTmpl, exists := rs.getSpecifyTmpl(gameType)
	if !exists {
		return nil, fmt.Errorf("Game record template not found: %v", gameType)
	}

	currentRecord, err := rs.recordRepo.GetRecordByRoomId(ctx, roomId, recordTmpl)
	if err != nil {
		return nil, err
	}

	return currentRecord, nil
}

func (rs *RecordService) getGameStrategy(gameType shared.BilliardGameType) (game.IGameStrategy, error) {
	strategy, ok := rs.gameStrategyFactory[gameType]
	if !ok {
		return nil, fmt.Errorf("Game strategy not found: %v", gameType)
	}
	return strategy, nil
}

func (rs *RecordService) getGameStrategyByRoomId(ctx context.Context, roomId int) (game.IGameStrategy, error) {
	gt, err := rs.roomRepo.GetRoomGameType(ctx, roomId)
	if err != nil {
		return nil, err
	}

	return rs.getGameStrategy(gt)
}

func (rs *RecordService) GetRoomActions(ctx context.Context, roomId int) ([]record.Action, error) {
	gameStrategy, err := rs.getGameStrategyByRoomId(ctx, roomId)
	if err != nil {
		return nil, err
	}

	gameActions, err := gameStrategy.GetRoomActions(ctx, roomId)
	if err != nil {
		return nil, err
	}

	actions := putils.Convert(gameActions, func(g game.Action) record.Action {
		return g
	})

	return actions, nil
}

func (rs *RecordService) HandleRecord(ctx context.Context, roomId int, gameType shared.BilliardGameType, rawRecord json.RawMessage) ([]record.RecordItem, error) {
	gameStrategy, err := rs.getGameStrategy(gameType)
	if err != nil {
		return nil, errors.Wrap(err, "getGameStrategy")
	}

	records, err := gameStrategy.UnmarshalRecord(rawRecord)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalRecord")
	}

	convertRecords := putils.Convert(records, func(r game.Record) record.RecordItem {
		return r
	})

	err = rs.recordRepo.UpdateRoomRecord(ctx, roomId, gameType, convertRecords)
	if err != nil {
		return nil, errors.Wrap(err, "updateRoomRecord")
	}

	return convertRecords, nil
}

func (rs *RecordService) HandleAction(ctx context.Context, gameType shared.BilliardGameType, rawAction json.RawMessage) (record.Action, error) {
	gameStrategy, err := rs.getGameStrategy(gameType)
	if err != nil {
		return nil, errors.Wrap(err, "getGameStrategy")
	}

	action, err := gameStrategy.UnmarshalAction(rawAction)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalAction")
	}

	err = gameStrategy.HandleAction(ctx, action)
	if err != nil {
		return nil, errors.Wrap(err, "handleAction")
	}

	return action, nil
}
