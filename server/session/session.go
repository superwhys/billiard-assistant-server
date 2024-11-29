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
	"net/http"
	"slices"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/session"
)

var _ session.ISessionService = (*sessionService)(nil)

type sessionService struct {
	sessMap           map[string]*session.Session
	roomSession       map[int][]string
	websocketUpgrader *websocket.Upgrader
}

func NewSessionService() *sessionService {
	return &sessionService{
		sessMap:     make(map[string]*session.Session),
		roomSession: make(map[int][]string),
		websocketUpgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *sessionService) CreateSession(ctx context.Context, playerID int, roomID int, w http.ResponseWriter, r *http.Request) (*session.Session, error) {
	ws, err := s.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		plog.Errorc(ctx, "upgrade websocket error: %v", err)
		return nil, err
	}

	sess := session.NewSession(ctx, roomID, playerID, ws)
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

		for msg := range sess.IterReadMessage() {
			msg.SetMessageOwner(sess.UserId)
			msgChan <- msg
		}
	}()

	return msgChan
}

func (s *sessionService) StartSession(sess *session.Session, sessionHandler session.SessionEventHandler) {
	defer s.RemoveSession(sess.ID)

	messageLoop := s.readMessageLoop(sess)
	for {
		select {
		case <-sess.Ctx.Done():
			plog.Debugc(sess.Ctx, "session done")
			break
		case msg, ok := <-messageLoop:
			if !ok {
				plog.Debugc(sess.Ctx, "connection closed")
				break
			}

			err := sessionHandler(sess.Ctx, msg)
			if err != nil {
				plog.Errorc(sess.Ctx, "handle message failed: %v", sess, err)
				continue
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

func (s *sessionService) GetSessionByUserRoom(roomId int, userId int) (*session.Session, error) {
	roomSess, exists := s.roomSession[roomId]
	if !exists {
		return nil, errors.New("room sessions not found")
	}

	for _, sessId := range roomSess {
		sess, err := s.GetSessionByID(sessId)
		if err != nil {
			plog.Errorf("room(%v) session(%v) not found", roomId, sessId)
			continue
		}
		if sess.UserId == userId {
			return sess, nil
		}
	}

	return nil, errors.New("user not found in room")
}

func (s *sessionService) BroadcastMessage(roomID, publishUserId int, message *session.Message) error {
	roomSess, exists := s.roomSession[roomID]
	if !exists {
		return errors.New("room sessions not found")
	}

	plog.Debugf("broadcasting room(%v), players cnt: %v", roomID, len(roomSess))

	for _, sessID := range roomSess {
		sess, err := s.GetSessionByID(sessID)
		if err != nil {
			plog.Errorf("room(%v) session(%v) not found", roomID, sessID)
			continue
		}

		if sess.UserId == publishUserId {
			continue
		}

		err = sess.SendMessage(message)
		if err != nil {
			plog.Errorf("session(%s) broadcast message failed: %v", sess, err)
			continue
		}
		plog.Debugf("room(%v) session(%v) broadcast message: %v", roomID, sess, message)
	}

	return nil
}
