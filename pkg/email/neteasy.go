// File:		neteasy.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"time"

	"github.com/go-puzzles/puzzles/goredis"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/pqueue"
	"github.com/pkg/errors"
)

const (
	from       = "Billiard Assistant"
	subject    = "Billiard 邮箱验证码"
	smtpServer = "smtp.163.com"
	smtpPort   = "465"
	asyncQueue = "billiard-email-send"
)

type NetEasyEmailSender struct {
	conf    *EmailConf
	auth    smtp.Auth
	tlsConf *tls.Config
	queue   *pqueue.RedisQueue[*AsyncEmailTask]
}

func NewNetEasySender(conf *EmailConf, redisClient *goredis.PuzzleRedisClient) *NetEasyEmailSender {
	auth := smtp.PlainAuth("", conf.Sender, conf.Password, smtpServer)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer,
	}

	queue := pqueue.NewRedisQueueWithClient[*AsyncEmailTask](redisClient, asyncQueue)
	return &NetEasyEmailSender{
		conf:    conf,
		auth:    auth,
		tlsConf: tlsConfig,
		queue:   queue,
	}
}

func (e *NetEasyEmailSender) LoopAsyncTask(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		task, err := e.queue.Dequeue()
		if errors.Is(err, pqueue.QueueEmptyError) {
			time.Sleep(time.Millisecond * 500)
			continue
		} else if err != nil {
			plog.Errorc(ctx, "async task dequeue error: %v", err)
			return errors.Wrap(err, "asyncTaskDequeue")
		}

		if err := e.SendMsg(ctx, task.Target(), []byte(task.Message())); err != nil {
			plog.Errorc(ctx, "send email to %s error: %v", task.Target(), err)
			time.Sleep(time.Millisecond * 500)
			continue
		}
	}
}

func (e *NetEasyEmailSender) getClient(target string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", smtpServer+":"+smtpPort, e.tlsConf)
	if err != nil {
		return nil, errors.Wrap(err, "dialSmtp")
	}

	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		return nil, errors.Wrap(err, "newSmtpClient")
	}

	if err = client.Auth(e.auth); err != nil {
		return nil, errors.Wrap(err, "smtpAuth")
	}

	if err = client.Mail(e.conf.Sender); err != nil {
		return nil, errors.Wrap(err, "smtpMail")
	}
	if err = client.Rcpt(target); err != nil {
		return nil, errors.Wrap(err, "smtpRcpt")
	}

	return client, nil
}

func (e *NetEasyEmailSender) wrapMsg(subject, target string, msg []byte) []byte {
	subject = fmt.Sprintf("Subject: %s!\r\n,", subject)
	fromHeader := fmt.Sprintf("From: %s\r\n", from)
	toHeader := fmt.Sprintf("To: %s\r\n", target)
	contentTypeHeader := "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
	mimeVersionHeader := "MIME-Version: 1.0\r\n"

	data := []byte(fromHeader + toHeader + subject + mimeVersionHeader + contentTypeHeader + "\r\n")
	return append(data, msg...)
}

func (e *NetEasyEmailSender) SendMsg(ctx context.Context, target string, msg []byte) error {
	client, err := e.getClient(target)
	if err != nil {
		return errors.Wrapf(err, "getSmtpClient: %v", target)
	}
	defer client.Quit()

	w, err := client.Data()
	if err != nil {
		return errors.Wrapf(err, "getWriter: %v", target)
	}
	defer w.Close()

	_, err = w.Write(e.wrapMsg(subject, target, msg))
	if err != nil {
		return errors.Wrapf(err, "writeData: %v", target)
	}

	return nil
}

func (e *NetEasyEmailSender) AsyncSendMsg(ctx context.Context, task *AsyncEmailTask) error {
	return e.queue.Enqueue(task)
}
