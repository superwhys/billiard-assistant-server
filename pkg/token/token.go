// File:		token.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package token

import (
	"fmt"
	"time"

	"github.com/go-puzzles/puzzles/plog"
)

type Manager struct {
	storage     Storage
	cachePrefix string
	cacheTTL    time.Duration
}

type ManagerOption func(*Manager)

func WithCachePrefix(prefix string) ManagerOption {
	return func(tm *Manager) {
		tm.cachePrefix = prefix
	}
}

func WithCacheTTL(ttl time.Duration) ManagerOption {
	return func(tm *Manager) {
		tm.cacheTTL = ttl
	}
}

func NewManager(storage Storage, opts ...ManagerOption) *Manager {
	tm := &Manager{
		storage:  storage,
		cacheTTL: defaultTTL,
	}

	for _, opt := range opts {
		opt(tm)
	}

	return tm
}

func (tm *Manager) getKey(t Token) string {
	return tm.getKeyWithKeyId(t.GetKey(), t)
}

func (tm *Manager) getKeyWithKeyId(id string, t Token) string {
	key := plog.GetStructName(t)
	if tm.cachePrefix != "" {
		key = fmt.Sprintf("%s:%s", tm.cachePrefix, key)
	}

	return fmt.Sprintf("%v:%v", key, id)
}

func (tm *Manager) Save(t Token) error {
	return tm.storage.SetWithTTL(tm.getKey(t), t, tm.cacheTTL)
}

func (tm *Manager) Read(tokenId string, t Token) error {
	key := tm.getKeyWithKeyId(tokenId, t)
	return tm.storage.Get(key, t)
}

func (tm *Manager) Remove(t Token) error {
	return tm.storage.Delete(tm.getKey(t))
}
