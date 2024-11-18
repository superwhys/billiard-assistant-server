// File:		events.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"github.com/superwhys/billiard-assistant-server/domain/room"
	"github.com/superwhys/billiard-assistant-server/domain/session"
	"github.com/superwhys/billiard-assistant-server/domain/user"
	"github.com/superwhys/billiard-assistant-server/pkg/events"
	"github.com/superwhys/billiard-assistant-server/server/dto"
)

func (s *BilliardServer) setupEventsSubscription() {
	s.EventBus.Subscribe(events.PlayerJoined, s.HandlePlayerEnterEvent)
	s.EventBus.Subscribe(events.PlayerLeft, s.HandlePlayerLeaveEvent)
	s.EventBus.Subscribe(events.GameStart, s.HandleGameStartEvent)
	s.EventBus.Subscribe(events.SendPhoneCode, s.HandleSendPhoneSMS)
	s.EventBus.Subscribe(events.SendEmailCode, s.HandleSendEmailCode)
}

func (s *BilliardServer) HandlePlayerEnterEvent(event *events.EventMessage) error {
	// broadcast new player joined event
	e := event.Payload.(*room.EnterRoomEvent)
	u := e.User.(*user.User)
	return s.SessionSrv.BroadcastMessage(e.RoomId, &session.Message{
		EventType: event.EventType,
		Data:      dto.UserEntityToDto(u),
	})
}

func (s *BilliardServer) HandlePlayerLeaveEvent(event *events.EventMessage) error {
	// TODO: Implement player leave logic
	// broadcast new player joined event, etc.
	panic("not implemented")
}

func (s *BilliardServer) HandleGameStartEvent(event *events.EventMessage) error {
	// TODO: not implemented
	// Broadcast all user gameStart message
	// Initialize game state
	panic("not implemented")
}

func (s *BilliardServer) HandleSendPhoneSMS(event *events.EventMessage) error {
	// TODO: 实现具体的短信发送逻辑
	// plog.Infoc(ctx, "sending SMS code %s to phone %s", code, phone)
	return nil
}

func (s *BilliardServer) HandleSendEmailCode(event *events.EventMessage) error {
	// TODO: 实现具体的邮件发送逻辑
	// plog.Infoc(ctx, "sending email code %s to email %s", code, email)
	return nil
}
