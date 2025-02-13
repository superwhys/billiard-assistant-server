// File:		events.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"context"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/record"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/session"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/events"
	"gitea.hoven.com/billiard/billiard-assistant-server/server/dto"
	"github.com/pkg/errors"
)

func (s *BilliardServer) setupEventsSubscription() {
	s.EventBus.Subscribe(events.ConnectHeartbeat, s.HandleHeartbeatEvent)
	s.EventBus.Subscribe(events.PlayerJoined, s.HandlePlayerEnterEvent)
	s.EventBus.Subscribe(events.PlayerLeft, s.HandlePlayerLeaveEvent)
	s.EventBus.Subscribe(events.GameStart, s.HandleGameStartEvent)
	s.EventBus.Subscribe(events.GameEnd, s.HandleGameEndEvent)
	s.EventBus.Subscribe(events.RecordAction, s.HandleRecordAction)
	s.EventBus.Subscribe(events.PlayerOnline, s.HandlePlayerOnlineOfflineEvent)
	s.EventBus.Subscribe(events.PlayerOffline, s.HandlePlayerOnlineOfflineEvent)
}

func (s *BilliardServer) HandlePlayerOnlineOfflineEvent(events *events.EventMessage) error {
	msg := events.Payload.(*room.PlayerOnlineOfflineEvent)

	r, err := s.RoomSrv.GetRoomById(context.TODO(), msg.RoomId)
	if err != nil {
		return errors.Wrap(err, "getRoomById")
	}

	if r.GameStatus != room.Playing {
		return nil
	}

	return s.SessionSrv.BroadcastMessage(msg.RoomId, msg.UserId, &session.Message{
		EventType: events.EventType,
		Data:      msg.UserId,
	})
}

func (s *BilliardServer) HandleHeartbeatEvent(events *events.EventMessage) error {
	msg := events.Payload.(*session.Message)
	roomId := int(msg.Data.(float64))
	userId := events.MessageOwner

	err := s.RoomSrv.UpdateRoomUserHeartbeart(context.TODO(), roomId, userId)
	if err != nil {
		return errors.Wrap(err, "updateRoomUserHeartbeart")
	}

	room, err := s.RoomSrv.GetRoomById(context.TODO(), roomId)
	if err != nil {
		return errors.Wrap(err, "getRoomById")
	}

	sess := msg.Sess
	return sess.SendMessage(&session.Message{
		EventType: events.EventType,
		Data:      dto.GameRoomEntityToDto(room),
	})
}

func (s *BilliardServer) HandleRecordAction(event *events.EventMessage) error {
	e, ok := event.Payload.(*record.ActionEvent)
	if !ok {
		return errors.New("invalid payload type for record action event")
	}

	return s.SessionSrv.BroadcastMessage(e.RoomId, e.UserId, &session.Message{
		EventType: event.EventType,
		Data:      e.Action,
	})
}

func (s *BilliardServer) HandlePlayerEnterEvent(event *events.EventMessage) error {
	// broadcast new player joined event
	e, ok := event.Payload.(*room.EnterRoomEvent)
	if !ok {
		return errors.New("invalid payload type for player joined event")
	}
	return s.SessionSrv.BroadcastMessage(e.RoomId, e.UserId, &session.Message{
		EventType: event.EventType,
		Data:      e,
	})
}

func (s *BilliardServer) HandlePlayerLeaveEvent(event *events.EventMessage) error {
	e, ok := event.Payload.(*room.LeaveRoomEvent)
	if !ok {
		return errors.New("invalid payload type for player leave event")
	}
	// broadcast user leave event
	return s.SessionSrv.BroadcastMessage(e.RoomId, e.UserId, &session.Message{
		EventType: event.EventType,
		Data:      e,
	})
}

func (s *BilliardServer) HandleGameStartEvent(event *events.EventMessage) error {
	e, ok := event.Payload.(*room.GameStartEvent)
	if !ok {
		return errors.New("invalid payload type for game start event")
	}

	return s.SessionSrv.BroadcastMessage(e.RoomId, e.UserId, &session.Message{
		EventType: event.EventType,
		Data:      e,
	})
}

func (s *BilliardServer) HandleGameEndEvent(event *events.EventMessage) error {
	e, ok := event.Payload.(*room.GameEndEvent)
	if !ok {
		return errors.New("invalid payload type for game end event")
	}

	return s.SessionSrv.BroadcastMessage(e.RoomId, e.UserId, &session.Message{
		EventType: event.EventType,
		Data:      e,
	})
}
