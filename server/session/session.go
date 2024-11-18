// File:		session.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package sessionSrv

import (
	"context"
	"slices"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/superwhys/billiard-assistant-server/domain/session"
)

var _ session.ISessionService = (*sessionService)(nil)

type sessionService struct {
	sessMap     map[string]*session.Session
	roomSession map[int][]string
}

func (s *sessionService) CreateSession(ctx context.Context, playerID int, roomID int, conn *websocket.Conn) (*session.Session, error) {
	sess := session.NewSession(ctx, roomID, playerID, conn)
	s.sessMap[sess.ID] = sess
	
	if rs := s.roomSession[roomID]; rs == nil {
		s.roomSession[roomID] = make([]string, 0)
	}
	
	s.roomSession[roomID] = append(s.roomSession[roomID], sess.ID)
	
	return sess, nil
}

func (s *sessionService) readMessageLoop(sess *session.Session) chan *session.Message {
	msgChan := make(chan *session.Message)
	
	go func() {
		defer close(msgChan)
		
		for {
			var msg session.Message
			err := sess.Conn.ReadJSON(&msg)
			if err != nil {
				plog.Infof("session(%s) read message failed: %v", sess, err)
				return
			}
			
			msgChan <- &msg
		}
	}()
	
	return msgChan
}

func (s *sessionService) StartSession(sess *session.Session) {
	for {
		select {
		case <-sess.Ctx.Done():
			break
		case msg := <-s.readMessageLoop(sess):
			plog.Infof("read message %v", msg)
			
			switch msg.EventType {
			// TODO:
			}
		}
	}
}

func (s *sessionService) RemoveSession(sessionID string) error {
	sess, exists := s.sessMap[sessionID]
	if !exists {
		return errors.New("session not found")
	}
	
	delete(s.sessMap, sessionID)
	
	roomSess := s.roomSession[sess.RoomId]
	roomSess = slices.DeleteFunc(roomSess, func(id string) bool {
		return id == sessionID
	})
	
	return nil
}

func (s *sessionService) GetSessionByID(sessionID string) (*session.Session, error) {
	sess, exists := s.sessMap[sessionID]
	if !exists {
		return nil, errors.New("session not found")
	}
	
	return sess, nil
}

func (s *sessionService) BroadcastMessage(roomID int, message *session.Message) error {
	roomSess, exists := s.roomSession[roomID]
	if !exists {
		return errors.New("room not found")
	}
	
	for _, sessID := range roomSess {
		sess, err := s.GetSessionByID(sessID)
		if err != nil {
			plog.Errorf("room(%v) session(%v) not found", roomID, sessID)
			continue
		}
		
		err = sess.Conn.WriteJSON(message)
		if err != nil {
			plog.Errorf("session(%s) broadcast message failed: %v", sess, err)
			continue
		}
	}
	
	return nil
}
