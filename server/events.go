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

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/email"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/events"
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
	// e := event.Payload.(*room.EnterRoomEvent)
	// return s.SessionSrv.BroadcastMessage(e.RoomId, &session.Message{
	// 	EventType: event.EventType,
	// })
	// TODO: Implement player enget logic
	// broadcast new player joined event, etc.
	panic("not implemented")

}

func (s *BilliardServer) HandlePlayerLeaveEvent(event *events.EventMessage) error {
	// TODO: Implement player leave logic
	// broadcast player leave event, etc.
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
