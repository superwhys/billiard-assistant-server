// File:		room.go
// Created by:	Hoven
// Created on:	2025-01-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"context"
	"encoding/json"
	"net/http"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/record"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/events"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitea.hoven.com/billiard/billiard-assistant-server/server/dto"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
)

func (s *BilliardServer) GetUserGameRooms(ctx context.Context, userId int) ([]*dto.GameRoom, error) {
	rs, err := s.RoomSrv.GetUserGameRooms(ctx, userId)
	if err != nil {
		plog.Errorc(ctx, "get user game rooms error: %v", err)
		return nil, err
	}

	ret := make([]*dto.GameRoom, 0, len(rs))
	for _, r := range rs {
		ret = append(ret, dto.GameRoomEntityToDto(r))
	}

	return ret, nil
}

func (s *BilliardServer) CreateRoom(ctx context.Context, userId, gameId int) (*dto.GameRoom, error) {
	gr, err := s.RoomSrv.CreateGameRoom(ctx, userId, gameId)
	if err != nil {
		plog.Errorc(ctx, "create game room error: %v", err)
		return nil, err
	}

	err = s.RoomSrv.EnterGameRoom(ctx, gr.RoomId, userId, "")
	if err != nil {
		plog.Errorc(ctx, "enter game room error: %v", err)
		return nil, err
	}

	return dto.GameRoomEntityToDto(gr), nil
}

func (s *BilliardServer) UpdateGameRoomStatus(ctx context.Context, userId int, req *dto.UpdateGameRoomRequest) error {
	gr := &room.Room{
		RoomId:        req.RoomId,
		OwnerId:       userId,
		GameStatus:    req.GameStatus,
		WinLoseStatus: req.WinLoseStatus,
	}

	err := s.RoomSrv.UpdateGameRoomStatus(ctx, gr)
	if err != nil {
		plog.Errorc(ctx, "update game room status error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) UpdateGameRoomExtra(ctx context.Context, userId int, gameRoom *dto.UpdateGameRoomExtraRequest) error {
	gr := &room.Room{
		RoomId:  gameRoom.RoomId,
		OwnerId: userId,
		Extra:   gameRoom.Extra,
	}

	err := s.RoomSrv.UpdateGameRoomStatus(ctx, gr)
	if err != nil {
		plog.Errorc(ctx, "update game room extra error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) DeleteRoom(ctx context.Context, userId int, roomId int) error {
	err := s.RoomSrv.DeleteGameRoom(ctx, userId, roomId)
	if err != nil {
		plog.Errorc(ctx, "delete game room error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) EnterGameRoom(ctx context.Context, roomId, currentUid int, virtualUser string) (err error) {
	err = s.RoomSrv.EnterGameRoom(ctx, currentUid, roomId, virtualUser)
	if err != nil {
		plog.Errorc(ctx, "enter game room error: %v", err)
		return err
	}

	isVirtual := virtualUser != ""
	s.EventBus.Publish(room.NewEnterRoomEvent(roomId, currentUid, virtualUser, isVirtual))

	return nil
}

func (s *BilliardServer) LeaveGameRoom(ctx context.Context, roomId, currentUid int, virtualUser string) (err error) {
	err = s.RoomSrv.QuitGameRoom(ctx, roomId, currentUid, virtualUser)
	if err != nil {
		plog.Errorc(ctx, "leave game room error: %v", err)
		return err
	}

	isVirtual := virtualUser != ""

	// publish user leave room events
	s.EventBus.Publish(room.NewLeaveRoomEvent(roomId, currentUid, virtualUser, isVirtual))

	return nil
}

func (s *BilliardServer) GetGameRoom(ctx context.Context, roomId int) (*dto.GameRoom, error) {
	r, err := s.RoomSrv.GetRoomById(ctx, roomId)
	if err != nil {
		plog.Errorc(ctx, "get game room error: %v", err)
		return nil, err
	}

	record, err := s.RecordSrv.GetCurrentRecord(ctx, roomId, shared.BilliardGameType(r.Game.GetGameType()))
	if err != nil && !errors.Is(err, exception.ErrRoomRecordNotFound) {
		plog.Errorc(ctx, "get current record error: %v", err)
		return nil, err
	}

	if record != nil {
		r.Record = record
	}

	return dto.GameRoomEntityToDto(r), nil
}

func (s *BilliardServer) GetGameRoomByCode(ctx context.Context, roomCode string) (*dto.GameRoom, error) {
	r, err := s.RoomSrv.GetRoomByCode(ctx, roomCode)
	if err != nil {
		plog.Errorc(ctx, "get game room error: %v", err)
		return nil, err
	}

	reocrd, err := s.RecordSrv.GetCurrentRecord(ctx, r.GetRoomId(), shared.BilliardGameType(r.Game.GetGameType()))
	if err != nil && !errors.Is(err, exception.ErrRoomRecordNotFound) {
		plog.Errorc(ctx, "get current record error: %v", err)
		return nil, err
	}
	if reocrd != nil {
		r.Record = reocrd
	}

	return dto.GameRoomEntityToDto(r), nil
}

func (s *BilliardServer) CreateRoomSession(ctx context.Context, userId, roomId int, w http.ResponseWriter, r *http.Request) error {
	_, err := s.RoomSrv.GetRoomById(ctx, roomId)
	if errors.Is(err, exception.ErrGameRoomNotFound) {
		return err
	} else if err != nil {
		plog.Errorc(ctx, "get room error: %v", err)
		return err
	}

	sess, err := s.SessionSrv.CreateSession(ctx, userId, roomId, w, r)
	if err != nil {
		plog.Errorc(ctx, "register room session error: %v", err)
		return err
	}
	s.EventBus.Publish(room.NewPlayerOnlineOfflineEvent(events.PlayerOnline, roomId, userId))

	defer func() {
		s.SessionSrv.RemoveSession(sess.ID)
		sess.Close()
		s.EventBus.Publish(room.NewPlayerOnlineOfflineEvent(events.PlayerOffline, roomId, userId))
	}()

	go s.SessionSrv.StartSession(sess, s.handleSessionMessage)

	return sess.Wait()
}

func (s *BilliardServer) HandleRoomAction(ctx context.Context, roomId, userId int, rawAction json.RawMessage) error {
	gameType, err := s.RoomSrv.GetRoomGameType(ctx, roomId)
	if err != nil {
		plog.Errorc(ctx, "get room game type error: %v", err)
		return err
	}

	action, err := s.RecordSrv.HandleAction(ctx, gameType, rawAction)
	if err != nil {
		plog.Errorc(ctx, "handle room action error: %v", err)
		return err
	}

	s.EventBus.Publish(record.NewActionEvent(roomId, userId, action))

	return nil
}

func (s *BilliardServer) HandleRoomRecord(ctx context.Context, roomId int, rawRecord json.RawMessage) error {
	gameType, err := s.RoomSrv.GetRoomGameType(ctx, roomId)
	if err != nil {
		plog.Errorc(ctx, "get room game type error: %v", err)
		return err
	}

	_, err = s.RecordSrv.HandleRecord(ctx, roomId, gameType, rawRecord)
	if err != nil {
		plog.Errorc(ctx, "handle room record error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) GetRoomActions(ctx context.Context, roomId int) (*dto.Action, error) {
	action, err := s.RecordSrv.GetRoomActions(ctx, roomId)
	if err != nil {
		plog.Errorc(ctx, "get room actions error: %v", err)
		return nil, err
	}
	return &dto.Action{
		Actions: action,
		RoomId:  roomId,
	}, nil
}

func (s *BilliardServer) GetRoomRecoed(ctx context.Context, roomId int) (*dto.Record, error) {
	gameType, err := s.RoomSrv.GetRoomGameType(ctx, roomId)
	if err != nil {
		plog.Errorc(ctx, "get room game type error: %v", err)
		return nil, err
	}

	record, err := s.RecordSrv.GetCurrentRecord(ctx, roomId, gameType)
	if err != nil {
		plog.Errorc(ctx, "get room record error: %v", err)
		return nil, err
	}
	return dto.RecordEntityToDto(record), nil
}
