// File:		locker.go
// Created by:	Hoven
// Created on:	2024-11-04
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package locker

import (
	"fmt"
	"time"

	"github.com/go-puzzles/predis"
)

const (
	defaultPrefix = "sa:locker"
	defaultRetry  = 3
	defaultTTL    = time.Second * 5
)

type Locker struct {
	prefix string
	retry  int
	ttl    time.Duration
	client *predis.RedisClient
}

type LockerOption func(*Locker)

func WithPrefix(prefix string) LockerOption {
	return func(l *Locker) {
		l.prefix = prefix
	}
}

func WithRetry(retry int) LockerOption {
	return func(l *Locker) {
		l.retry = retry
	}
}

func WithTTL(dura time.Duration) LockerOption {
	return func(l *Locker) {
		l.ttl = dura
	}
}

func NewLocker(client *predis.RedisClient, opts ...LockerOption) *Locker {
	l := &Locker{
		client: client,
		prefix: defaultPrefix,
		ttl:    defaultTTL,
		retry:  defaultRetry,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Locker) Lock(key any) (err error) {
	lockKey := fmt.Sprintf("%s:%v", l.prefix, key)

	return l.client.LockWithBlock(lockKey, l.retry, l.ttl)
}

func (l *Locker) Unlock(key any) (err error) {
	lockKey := fmt.Sprintf("%s:%v", l.prefix, key)

	return l.client.UnLock(lockKey)
}
