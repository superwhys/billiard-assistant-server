// File:		locker.go
// Created by:	Hoven
// Created on:	2024-11-04
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package locker

import (
	"context"
	"fmt"
	"time"

	"github.com/go-puzzles/puzzles/goredis"
)

const (
	defaultPrefix  = "billiard:locker"
	defaultTTL     = time.Second * 5
	defaultTimeout = time.Second * 10
)

type Locker struct {
	prefix  string
	ttl     time.Duration
	timeout time.Duration
	client  *goredis.PuzzleRedisClient
}

type LockerOption func(*Locker)

func WithPrefix(prefix string) LockerOption {
	return func(l *Locker) {
		l.prefix = prefix
	}
}

func WithTimeout(to time.Duration) LockerOption {
	return func(l *Locker) {
		l.timeout = to
	}
}

func WithTTL(dura time.Duration) LockerOption {
	return func(l *Locker) {
		l.ttl = dura
	}
}

func NewLocker(client *goredis.PuzzleRedisClient, opts ...LockerOption) *Locker {
	l := &Locker{
		client:  client,
		prefix:  defaultPrefix,
		ttl:     defaultTTL,
		timeout: defaultTimeout,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Locker) Lock(ctx context.Context, key any) (err error) {
	lockKey := fmt.Sprintf("%s:%v", l.prefix, key)
	return l.client.TryLockWithTimeout(ctx, lockKey, l.ttl, l.timeout)
}

func (l *Locker) Unlock(ctx context.Context, key any) (err error) {
	lockKey := fmt.Sprintf("%s:%v", l.prefix, key)
	return l.client.Unlock(ctx, lockKey)
}
