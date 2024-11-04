// File:		storage.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package token

import "time"

const (
	defaultTTL = 3600 * time.Second
)

type Storage interface {
	SetWithTTL(key string, value any, ttl time.Duration) error
	Get(key string, out any) error
	Delete(key string) error
}

type Token interface {
	GetKey() string
}
