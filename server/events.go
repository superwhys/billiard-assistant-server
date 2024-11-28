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

	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/record"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/room"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/session"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/email"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/events"
	"gitlab.hoven.com/billiard/billiard-assistant-server/server/dto"
)

func (s *BilliardServer) setupEventsSubscription() {
	s.EventBus.Subscribe(events.ConnectHeartbeat, s.HandleHeartbeatEvent)
	s.EventBus.Subscribe(events.PlayerJoined, s.HandlePlayerEnterEvent)
	s.EventBus.Subscribe(events.PlayerLeft, s.HandlePlayerLeaveEvent)
	s.EventBus.Subscribe(events.GameStart, s.HandleGameStartEvent)
	s.EventBus.Subscribe(events.SendPhoneCode, s.HandleSendPhoneSMS)
	s.EventBus.Subscribe(events.SendEmailCode, s.HandleSendEmailCode)
	s.EventBus.Subscribe(events.RecordAction, s.HandleRecordAction)
}

func (s *BilliardServer) HandleHeartbeatEvent(events *events.EventMessage) error {
	roomId := int(events.Payload.(float64))
	userId := events.MessageOwner

	err := s.RoomSrv.UpdateRoomUserHeartbeart(context.TODO(), roomId, userId)
	if err != nil {
		return errors.Wrap(err, "updateRoomUserHeartbeart")
	}

	room, err := s.RoomSrv.GetRoomById(context.TODO(), roomId)
	if err != nil {
		return errors.Wrap(err, "getRoomById")
	}

	sess, err := s.SessionSrv.GetSessionByUserRoom(roomId, userId)
	if err != nil {
		return errors.Wrap(err, "getSessionByUserRoom")
	}

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

func (s *BilliardServer) HandleSendPhoneSMS(event *events.EventMessage) error {
	// TODO: 实现具体的短信发送逻辑
	// plog.Infoc(ctx, "sending SMS code %s to phone %s", code, phone)
	return nil
}

func (s *BilliardServer) HandleSendEmailCode(event *events.EventMessage) error {
	e := event.Payload.(*user.SendCodeEvent)

	msg, err := user.GenerateSendCodeEventMessage(e.Code)
	if err != nil {
		return err
	}

	return s.emailSender.AsyncSendMsg(context.Background(), &email.AsyncEmailTask{
		TargetEmail: e.Target,
		Msg:         string(msg),
	})
}
