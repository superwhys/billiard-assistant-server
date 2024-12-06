// File:		storage.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package token

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultTTL = 3600 * time.Second
)

type Storage interface {
	SetValue(ctx context.Context, key string, value any, ttl time.Duration) error
	GetValue(ctx context.Context, key string, out any) error
	Del(ctx context.Context, key ...string) *redis.IntCmd
}

type Token interface {
	GetKey() string
}
