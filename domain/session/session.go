// File:		session.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package session

import (
	"context"
	"fmt"
	
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/events"
)

/*
Session action:
1. EnterRoom -> Create a new Session -> BroadcastMessage to users who have previously joined the room
2. QuitRoom -> BroadcastMessage to user who in this room
*/

type Session struct {
	Ctx    context.Context
	ID     string
	RoomId int
	UserId int
	Conn   *websocket.Conn
	done   chan error
}

func NewSession(ctx context.Context, roomId int, userId int, conn *websocket.Conn) *Session {
	return &Session{
		Ctx:    ctx,
		ID:     uuid.New().String(),
		RoomId: roomId,
		UserId: userId,
		Conn:   conn,
		done:   make(chan error),
	}
}

func (s *Session) String() string {
	return fmt.Sprintf("room(%d) session(%s)", s.RoomId, s.ID)
}

func (s *Session) Wait() error {
	return <-s.done
}

func (s *Session) Close() error {
	defer close(s.done)
	return s.Conn.Close()
}

type Message struct {
	EventType events.EventType `json:"type"`
	Data      any              `json:"data"`
}
