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
	"encoding/json"
	"fmt"
	"iter"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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
	conn   *websocket.Conn
	done   chan error
}

func NewSession(ctx context.Context, roomId int, userId int, conn *websocket.Conn) *Session {
	sess := &Session{
		Ctx:    ctx,
		ID:     uuid.New().String(),
		RoomId: roomId,
		UserId: userId,
		conn:   conn,
		done:   make(chan error),
	}

	sess.Ctx = plog.With(sess.Ctx, sess.String())
	return sess
}

func (s *Session) String() string {
	return fmt.Sprintf("room(%d) user(%d) session(%s)", s.RoomId, s.UserId, s.ID)
}

func (s *Session) IterReadMessage() iter.Seq[*Message] {
	return func(yield func(*Message) bool) {
		for {
			msg := new(Message)
			err := s.conn.ReadJSON(&msg)
			if err != nil {
				plog.Errorc(s.Ctx, "read message failed: %v", err)
				s.done <- err
				return
			}

			if !yield(msg) {
				break
			}
		}
	}
}

func (s *Session) SendMessage(message *Message) error {
	b, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "marshalMessage")
	}

	return s.conn.WriteMessage(websocket.TextMessage, b)
}

func (s *Session) Wait() error {
	return <-s.done
}

func (s *Session) Close() error {
	return s.conn.Close()
}

type Message struct {
	ownerId   int
	Sess      *Session         `json:"-"`
	EventType events.EventType `json:"type"`
	Data      any              `json:"data"`
}

func (m *Message) GetMessageOwner() int {
	return m.ownerId
}

func (m *Message) SetMessageOwner(ownerId int) {
	m.ownerId = ownerId
}
