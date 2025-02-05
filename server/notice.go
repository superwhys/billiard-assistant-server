// File:		notice.go
// Created by:	Hoven
// Created on:	2025-01-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"context"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/notice"
	"gitea.hoven.com/billiard/billiard-assistant-server/server/dto"
	"github.com/go-puzzles/puzzles/plog"
)

func (s *BilliardServer) GetNoticeList(ctx context.Context) ([]*dto.Notice, error) {
	notices, err := s.NoticeSrv.GetNoticeList(ctx)
	if err != nil {
		plog.Errorc(ctx, "get notice list error: %v", err)
		return nil, err
	}

	ret := make([]*dto.Notice, 0, len(notices))
	for _, n := range notices {
		ret = append(ret, dto.NoticeEntityToDto(n))
	}

	return ret, nil
}

func (s *BilliardServer) GetSystemNotice(ctx context.Context) ([]*dto.Notice, error) {
	notices, err := s.NoticeSrv.GetNoticeByType(ctx, notice.System)
	if err != nil {
		plog.Errorc(ctx, "get notice list error: %v", err)
		return nil, err
	}

	ret := make([]*dto.Notice, 0, len(notices))
	for _, n := range notices {
		ret = append(ret, dto.NoticeEntityToDto(n))
	}

	return ret, nil
}

func (s *BilliardServer) AddNotices(ctx context.Context, req *dto.AddNoticeRequest) error {
	notices := make([]*notice.Notice, 0, len(req.Contents))
	for _, content := range req.Contents {
		notices = append(notices, &notice.Notice{
			NoticeType: req.NoticeType,
			Message:    content,
		})
	}

	return s.NoticeSrv.AddNotices(ctx, notices)
}
