// File:		notice.go
// Created by:	Hoven
// Created on:	2024-10-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import "gitlab.hoven.com/billiard/billiard-assistant-server/domain/notice"

type Notice struct {
	Message string
}

type GetNoticeListResp struct {
	Notices []*Notice `json:"notices"`
}

func NoticeEntityToDto(n *notice.Notice) *Notice {
	return &Notice{
		Message: n.Message,
	}
}
