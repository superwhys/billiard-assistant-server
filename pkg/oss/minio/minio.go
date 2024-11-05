// File:		minio.go
// Created by:	Hoven
// Created on:	2024-11-05
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package minio

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/superwhys/snooker-assistant-server/pkg/oss"
)

var _ oss.IOSS = (*MinioOss)(nil)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string

	Bucket string
}

type MinioOss struct {
	*Config
	client *minio.Client
}

func NewMinioOss(conf *Config) *MinioOss {
	client, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: false,
	})
	plog.PanicError(err)

	return &MinioOss{
		Config: conf,
		client: client,
	}
}

func (m *MinioOss) checkUrl(u string) (string, error) {
	if !strings.Contains(u, "://") {
		u = "http://" + u
	}
	url, err := url.Parse(u)
	if err != nil {
		return "", errors.Wrap(err, "parseMinioURL")
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	return url.String(), nil
}

func (m *MinioOss) getFileExt(file string) string {
	return filepath.Ext(file)
}

func (m *MinioOss) parseObjInfo(obj string) (originName, objName string) {
	fileName := fmt.Sprintf("%d-%s", time.Now().UnixMilli(), uuid.New().String())

	objDir := filepath.Dir(obj)
	ext := m.getFileExt(obj)

	originName = filepath.Base(obj)
	objName = fmt.Sprintf("%s/%s%s", objDir, fileName, ext)

	return
}

func (m *MinioOss) UploadFile(ctx context.Context, size int64, objName string, obj io.Reader) (uri string, err error) {
	originName, objName := m.parseObjInfo(objName)

	putOpt := minio.PutObjectOptions{
		UserTags: map[string]string{
			"type": "avatar",
			"name": originName,
		},
	}
	_, err = m.client.PutObject(ctx, m.Bucket, objName, obj, size, putOpt)
	if err != nil {
		return "", errors.Wrap(err, "uploadMinio")
	}

	fileUrl := fmt.Sprintf("%s/%s/%s", m.Endpoint, m.Bucket, objName)

	return m.checkUrl(fileUrl)
}

func (m *MinioOss) DownloadFile(ctx context.Context, objName string, dest string) (filepath string, err error) {
	panic("not implemented") // TODO: Implement
}
