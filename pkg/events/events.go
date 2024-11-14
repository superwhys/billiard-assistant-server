// File:		events.go
// Created by:	Hoven
// Created on:	2024-11-06
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package events

import (
	"sync"
	
	"github.com/go-puzzles/puzzles/plog"
)

type EventType int

const (
	PlayerJoined EventType = iota
	PlayerLeft
	PlayerPrepare
	GameStart
	GameEnd
)

type EventMessage struct {
	EventType EventType
	Payload   any
}

type eventFunc func(event *EventMessage) error

type EventFunc struct {
	name string
	fn   eventFunc
}

type EventBus struct {
	subscribers map[EventType][]*EventFunc
	lock        sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]*EventFunc),
		lock:        sync.RWMutex{},
	}
}

func (eb *EventBus) handleFn(ef *EventFunc, event *EventMessage) {
	defer func() {
		if err := recover(); err != nil {
			plog.Errorf("recovered from panic in handler event func(%s) for %s: %v", ef.name, event.EventType, err)
		}
	}()
	
	if err := ef.fn(event); err != nil {
		plog.Errorf("handler event func(%s) for %s error: %v", ef.name, event.EventType, err)
	}
	plog.Debugf("handler event func(%s) for %s completed", ef.name, event.EventType)
	
}

func (eb *EventBus) Publish(event *EventMessage) {
	eb.lock.RLock()
	defer eb.lock.RUnlock()
	
	subscribeFns, exists := eb.subscribers[event.EventType]
	if !exists {
		plog.Errorf("eventType %v no subscribers registered", event.EventType)
		return
	}
	
	for _, ef := range subscribeFns {
		go eb.handleFn(ef, event)
	}
}

func (eb *EventBus) Subscribe(eventType EventType, fn eventFunc) {
	eb.lock.Lock()
	defer eb.lock.Unlock()
	
	ef := &EventFunc{
		name: plog.GetFuncName(fn),
		fn:   fn,
	}
	
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ef)
}
