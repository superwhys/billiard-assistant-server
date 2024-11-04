// File:		notice.go
// Created by:	Hoven
// Created on:	2024-11-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package noticeSrv

import (
	"context"

	"github.com/pkg/errors"
	"github.com/superwhys/snooker-assistant-server/domain/notice"
)

var _ notice.INoticeService = (*NoticeService)(nil)

type NoticeService struct {
	noticeRepo notice.INoticeRepo
}

func NewNoticeService(noticeRepo notice.INoticeRepo) *NoticeService {
	return &NoticeService{noticeRepo: noticeRepo}
}

func (n *NoticeService) GetNoticeList(ctx context.Context) ([]*notice.Notice, error) {
	notices, err := n.noticeRepo.GetNotices(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getNoticeList")
	}
	return notices, nil
}
