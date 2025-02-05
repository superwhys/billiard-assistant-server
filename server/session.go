// File:		session.go
// Created by:	Hoven
// Created on:	2025-01-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"context"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/session"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/events"
	"github.com/pkg/errors"
)

func (s *BilliardServer) handleSessionMessage(ctx context.Context, msg *session.Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}

	s.EventBus.Publish(&events.EventMessage{
		EventType:    msg.EventType,
		MessageOwner: msg.GetMessageOwner(),
		Payload:      msg,
	})
	return nil
}
