// File:		email.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package email

import (
	"context"

	"github.com/go-puzzles/puzzles/pqueue"
)

type EmailConf struct {
	Sender   string
	Password string
}

type AsyncSendTask interface {
	pqueue.Item
	Target() string
	Mesasge() string
}

type EmailSender interface {
	SendMsg(ctx context.Context, target string, msg []byte) error
	AsyncSendMsg(ctx context.Context, task AsyncSendTask) error
}
