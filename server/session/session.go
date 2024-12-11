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
	"sync"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/session"
)

var _ session.ISessionService = (*sessionService)(nil)

type sessionService struct {
	sessMap           *sync.Map
	roomSession       *sync.Map
	websocketUpgrader *websocket.Upgrader
}

func NewSessionService() *sessionService {
	return &sessionService{
		sessMap:     &sync.Map{},
		roomSession: &sync.Map{},
		websocketUpgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *sessionService) getSession(sessId string) (*session.Session, error) {
	v, ok := s.sessMap.Load(sessId)
	if !ok {
		return nil, errors.New("session not found")
	}

	return v.(*session.Session), nil
}

func (s *sessionService) getRoomSessions(roomId int) ([]string, error) {
	v, ok := s.roomSession.Load(roomId)
	if !ok {
		return nil, errors.New("room session not found")
	}

	return v.([]string), nil
}

func (s *sessionService) CreateSession(ctx context.Context, playerID int, roomID int, w http.ResponseWriter, r *http.Request) (*session.Session, error) {
	ws, err := s.websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		plog.Errorc(ctx, "upgrade websocket error: %v", err)
		return nil, err
	}

	sess := session.NewSession(ctx, roomID, playerID, ws)

	s.sessMap.Store(sess.ID, sess)
	roomSessions, exists := s.roomSession.LoadOrStore(roomID, []string{sess.ID})
	if exists {
		roomSessions := append(roomSessions.([]string), sess.ID)
		s.roomSession.Store(roomID, roomSessions)
	}

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
			return
		case msg, ok := <-messageLoop:
			if !ok {
				plog.Debugc(sess.Ctx, "connection closed")
				return
			}

			msg.Sess = sess
			err := sessionHandler(sess.Ctx, msg)
			if err != nil {
				plog.Errorc(sess.Ctx, "handle message failed: %v", sess, err)
				continue
			}
		}
	}
}

func (s *sessionService) RemoveSession(sessionID string) error {
	sess, err := s.getSession(sessionID)
	if err != nil {
		return errors.Wrap(err, "getSession")
	}

	s.sessMap.Delete(sessionID)
	s.roomSession.Range(func(key, value any) bool {
		if key != sess.RoomId {
			return true
		}

		roomSessions := value.([]string)
		roomSessions = slices.DeleteFunc(roomSessions, func(id string) bool {
			return id == sessionID
		})
		s.roomSession.Store(key, roomSessions)
		return true
	})

	return nil
}

func (s *sessionService) GetSessionByID(sessionID string) (*session.Session, error) {
	sess, err := s.getSession(sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "getSession")
	}

	return sess, nil
}

func (s *sessionService) GetSessionByUserRoom(roomId int, userId int) (*session.Session, error) {
	rs, err := s.getRoomSessions(roomId)
	if err != nil {
		return nil, errors.Wrap(err, "getRoomSessions")
	}

	for _, sid := range rs {
		sess, err := s.getSession(sid)
		if err != nil {
			plog.Errorf("get room(%v) session(%v) err: %v", roomId, sid, err)
			continue
		}

		if sess.UserId != userId {
			continue
		}

		return sess, nil
	}

	return nil, errors.New("user session not found in room")
}

func (s *sessionService) BroadcastMessage(roomID, publishUserId int, message *session.Message) error {
	roomSess, err := s.getRoomSessions(roomID)
	if err != nil {
		return errors.Wrap(err, "getRoomSessions")
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
