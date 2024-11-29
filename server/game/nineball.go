// File:		nineball.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package gameSrv

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-puzzles/puzzles/predis"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/game/nineball"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/locker"
)

var _ nineball.INineballService = (*NineballService)(nil)

type NineballService struct {
	redisClient *predis.RedisClient
	locker      *locker.Locker
}

func NewNineballService(redisClient *predis.RedisClient) *NineballService {
	return &NineballService{
		redisClient: redisClient,
		locker:      locker.NewLocker(redisClient, locker.WithPrefix("nineball:record")),
	}
}

func (ns *NineballService) GetGameType() shared.BilliardGameType {
	return shared.NineBall
}

func (ns *NineballService) SetupGame(g shared.BaseGameConfig) []any {
	return nil
}

func (ns *NineballService) UnmarshalAction(rawAction json.RawMessage) (game.Action, error) {
	na := new(nineball.NineballAction)
	err := json.Unmarshal(rawAction, na)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal raw record to NineballAction")
	}

	return na, nil
}

func (ns *NineballService) UnmarshalRecord(rawRecord json.RawMessage) ([]game.Record, error) {
	var pr []*nineball.PlayerRecord
	err := json.Unmarshal(rawRecord, &pr)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal raw record to PlayerRecord")
	}

	prs := putils.Convert(pr, func(r *nineball.PlayerRecord) game.Record {
		return r
	})

	return prs, nil
}

func (ns *NineballService) HandleAction(ctx context.Context, action game.Action) error {
	ns.locker.Lock(action.GetActionRoomId())
	defer ns.locker.Unlock(action.GetActionRoomId())

	nineballAction, ok := action.(*nineball.NineballAction)
	if !ok {
		return errors.New("action must be a NineballAction object")
	}

	_, err := ns.redisClient.LPush(ns.getActionKey(nineballAction.RoomId), nineballAction)
	return err
}

func (ns *NineballService) GetRoomActions(ctx context.Context, roomId int) ([]game.Action, error) {
	actions := make([]*nineball.NineballAction, 0)
	actionKey := ns.getActionKey(roomId)

	err := ns.redisClient.LRange(actionKey, 0, -1, &actions)
	if err != nil {
		return nil, errors.Wrapf(err, "range room(%d) action", roomId)
	}

	converter := func(na *nineball.NineballAction) game.Action {
		return na
	}
	gameActions := putils.Convert(actions, converter)

	return gameActions, nil
}

func (ns *NineballService) getActionKey(roomId int) string {
	return fmt.Sprintf("nineball:action:%d", roomId)
}
