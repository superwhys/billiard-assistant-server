// File:		events.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"github.com/superwhys/snooker-assistant-server/domain/room"
	"github.com/superwhys/snooker-assistant-server/domain/session"
	"github.com/superwhys/snooker-assistant-server/domain/user"
	"github.com/superwhys/snooker-assistant-server/pkg/events"
	"github.com/superwhys/snooker-assistant-server/server/dto"
)

func (s *SaServer) setupEventsSubscription() {
	s.EventBus.Subscribe(events.PlayerJoined, s.HandlePlayerEnterEvent)
	s.EventBus.Subscribe(events.PlayerLeft, s.HandlePlayerLeaveEvent)
	s.EventBus.Subscribe(events.GameStart, s.HandleGameStartEvent)
	s.EventBus.Subscribe(events.SendPhoneCode, s.HandleSendPhoneSMS)
	s.EventBus.Subscribe(events.SendEmailCode, s.HandleSendEmailCode)
}

func (s *SaServer) HandlePlayerEnterEvent(event *events.EventMessage) error {
	// broadcast new player joined event
	e := event.Payload.(*room.EnterRoomEvent)
	u := e.User.(*user.User)
	return s.SessionSrv.BroadcastMessage(e.RoomId, &session.Message{
		EventType: event.EventType,
		Data:      dto.UserEntityToDto(u),
	})
}

func (s *SaServer) HandlePlayerLeaveEvent(event *events.EventMessage) error {
	// TODO: Implement player leave logic
	// broadcast new player joined event, etc.
	panic("not implemented")
}

func (s *SaServer) HandleGameStartEvent(event *events.EventMessage) error {
	// TODO: not implemented
	// Broadcast all user gameStart message
	// Initialize game state
	panic("not implemented")
}

func (s *SaServer) HandleSendPhoneSMS(event *events.EventMessage) error {
	// TODO: 实现具体的短信发送逻辑
	// plog.Infoc(ctx, "sending SMS code %s to phone %s", code, phone)
	return nil
}

func (s *SaServer) HandleSendEmailCode(event *events.EventMessage) error {
	// TODO: 实现具体的邮件发送逻辑
	// plog.Infoc(ctx, "sending email code %s to email %s", code, email)
	return nil
}
