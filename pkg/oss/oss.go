// File:		oss.go
// Created by:	Hoven
// Created on:	2024-11-05
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package oss

import (
	"context"
	"io"

	"github.com/go-puzzles/puzzles/pgin"
)

type OssSourceType string

const (
	SourceAvatar   OssSourceType = "avatar"
	SourceGameIcon OssSourceType = "game_icon"
)

func (st OssSourceType) String() string {
	return string(st)
}

type IOSS interface {
	pgin.Router
	UploadFile(ctx context.Context, size int64, sourceType OssSourceType, objName string, obj io.Reader) (uri string, err error)
	DownloadFile(ctx context.Context, objName string, dest string) (filepath string, err error)
	GetFile(ctx context.Context, objName string, w io.Writer) error
}
