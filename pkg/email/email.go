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
)

type EmailConf struct {
	Sender   string
	Password string
}

type AsyncEmailTask struct {
	TargetEmail string
	Msg         string
}

func (t *AsyncEmailTask) Target() string {
	return t.TargetEmail
}

func (t *AsyncEmailTask) Message() string {
	return t.Msg
}

func (t *AsyncEmailTask) Key() string {
	return "AsyncEmailTask" + t.TargetEmail
}

type EmailSender interface {
	SendMsg(ctx context.Context, target string, msg []byte) error
	AsyncSendMsg(ctx context.Context, task *AsyncEmailTask) error
}
