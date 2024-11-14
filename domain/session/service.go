// File:		service.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package session

import (
	"context"

	"github.com/gorilla/websocket"
)

type ISessionService interface {
	CreateSession(ctx context.Context, playerID, roomID int, conn *websocket.Conn) (*Session, error)
	StartSession(*Session)
	RemoveSession(sessionID string) error
	GetSessionByID(sessionID string) (*Session, error)
	BroadcastMessage(roomID int, message *Message) error
}
