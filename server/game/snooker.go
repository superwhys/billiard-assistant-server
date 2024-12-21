// File:		snooker.go
// Created by:	Hoven
// Created on:	2024-12-21
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package gameSrv

import (
	"context"
	"encoding/json"
	"fmt"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/game/snooker"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/locker"
	"github.com/go-puzzles/puzzles/goredis"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/pkg/errors"
)

var _ snooker.ISnookerService = (*SnookerService)(nil)

type SnookerService struct {
	redisClient *goredis.PuzzleRedisClient
	locker      *locker.Locker
}

func NewSnookerService(redisClient *goredis.PuzzleRedisClient) *SnookerService {
	return &SnookerService{
		redisClient: redisClient,
		locker:      locker.NewLocker(redisClient, locker.WithPrefix("nineball:record")),
	}
}

func (ss *SnookerService) GetGameType() shared.BilliardGameType {
	return shared.Snooker
}

func (ss *SnookerService) SetupGame(g shared.BaseGameConfig) []any {
	return nil
}

func (ss *SnookerService) UnmarshalAction(action json.RawMessage) (game.Action, error) {
	sa := new(snooker.SnookerAction)
	err := json.Unmarshal(action, sa)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal raw action to SnookerAction")
	}

	return sa, nil
}

func (ss *SnookerService) UnmarshalRecord(record json.RawMessage) ([]game.Record, error) {
	var pr []*snooker.SnookerPlayerRecord
	err := json.Unmarshal(record, &pr)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal raw record to SnookerPlayerRecord")
	}

	prs := putils.Convert(pr, func(r *snooker.SnookerPlayerRecord) game.Record {
		return r
	})

	return prs, nil
}

func (ss *SnookerService) HandleAction(ctx context.Context, action game.Action) error {
	ss.locker.Lock(ctx, action.GetActionRoomId())
	defer ss.locker.Unlock(ctx, action.GetActionRoomId())

	nineballAction, ok := action.(*snooker.SnookerAction)
	if !ok {
		return errors.New("action must be a NineballAction object")
	}

	err := ss.redisClient.LPushValue(ctx, ss.getActionKey(nineballAction.RoomId), nineballAction)
	return err
}

func (ss *SnookerService) getActionKey(roomId int) string {
	return fmt.Sprintf("nineball:action:%d", roomId)
}

func (ss *SnookerService) GetRoomActions(ctx context.Context, roomId int) ([]game.Action, error) {
	actions := make([]*snooker.SnookerAction, 0)
	actionKey := ss.getActionKey(roomId)

	err := ss.redisClient.RangeValue(ctx, actionKey, 0, -1, &actions)
	if err != nil {
		return nil, errors.Wrapf(err, "range room(%d) action", roomId)
	}

	converter := func(na *snooker.SnookerAction) game.Action {
		return na
	}
	gameActions := putils.Convert(actions, converter)

	return gameActions, nil
}
