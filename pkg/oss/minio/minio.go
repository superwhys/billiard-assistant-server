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
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/oss"
)

var _ oss.IOSS = (*MinioOss)(nil)

type MinioOss struct {
	*MinioConfig
	client  *minio.Client
	userApi string
}

func NewMinioOss(userApi string, conf *MinioConfig) *MinioOss {
	m := &MinioOss{
		MinioConfig: conf,
	}

	var err error
	m.userApi, err = m.checkUrl(userApi)
	plog.PanicError(err)

	discoverAddr := discover.GetAddress(conf.Endpoint)
	conf.Endpoint = discoverAddr

	m.client, err = minio.New(discoverAddr, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: false,
	})
	plog.PanicError(err)

	exists, err := m.client.BucketExists(context.TODO(), conf.Bucket)
	plog.PanicError(err)
	if !exists {
		plog.Fatalf("bucket %s not exists", conf.Bucket)
	}

	return m
}

func (m *MinioOss) checkUrl(u string) (string, error) {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "https://" + u
	}
	url, err := url.Parse(u)
	if err != nil {
		return "", errors.Wrap(err, "parseMinioURL")
	}

	if url.Scheme == "" {
		url.Scheme = "https"
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

	// https://billiard.superwhys.top/api/v1/user/avatar/1731850656800-d887240b-0177-44c7-853d-69f14b7cf874.jpeg
	return fmt.Sprintf("%s/%s", m.userApi, objName), nil
}

func (m *MinioOss) GetFile(ctx context.Context, objName string, w io.Writer) error {
	object, err := m.client.GetObject(ctx, m.Bucket, objName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()

	_, err = io.Copy(w, object)
	if err != nil {
		return errors.Wrap(err, "getMinioObject")
	}

	return nil
}

func (m *MinioOss) DownloadFile(ctx context.Context, objName string, dest string) (filepath string, err error) {
	panic("not implemented") // TODO: Implement
}
