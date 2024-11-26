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
	"net/http"
)

type SessionEventHandler func(ctx context.Context, msg *Message) error

type ISessionService interface {
	CreateSession(ctx context.Context, playerID, roomID int, w http.ResponseWriter, r *http.Request) (*Session, error)
	StartSession(*Session, SessionEventHandler)
	RemoveSession(sessionID string) error
	GetSessionByID(sessionID string) (*Session, error)
	BroadcastMessage(roomID int, message *Message) error
}
