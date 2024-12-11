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
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/cores/discover"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/oss"
)

var _ oss.IOSS = (*MinioOss)(nil)

type MinioOss struct {
	*MinioConfig
	client  *minio.Client
	baseApi string
}

func NewMinioOss(baseApi string, conf *MinioConfig) *MinioOss {
	m := &MinioOss{
		MinioConfig: conf,
		baseApi:     baseApi,
	}

	var err error
	m.baseApi, err = m.checkUrl(baseApi)
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

func (m *MinioOss) Init(router gin.IRouter) {
	minioGroup := router.Group("minio")
	minioGroup.GET(":sourceType/:sourceName", pgin.RequestHandler(m.getMinioSourceHandler))
}

type GetMinioSourceRequest struct {
	SourceType string `uri:"sourceType" binding:"required"`
	SourceName string `uri:"sourceName" binding:"required"`
}

func (m *MinioOss) getMinioSourceHandler(ctx *gin.Context, req *GetMinioSourceRequest) {
	objName := fmt.Sprintf("%s/%s", req.SourceType, req.SourceName)
	if err := m.GetFile(ctx, objName, ctx.Writer); err != nil {
		plog.Errorc(ctx, "get minio source(%s) error: %v", objName, err)
		ctx.Status(http.StatusBadRequest)
		return
	}
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

func (m *MinioOss) generateObjName(obj string) string {
	fileName := fmt.Sprintf("%d-%s", time.Now().UnixMilli(), uuid.New().String())
	ext := m.getFileExt(obj)
	return fileName + ext
}

func (m *MinioOss) UploadFile(ctx context.Context, size int64, sourceType oss.OssSourceType, objName string, obj io.Reader) (uri string, err error) {
	newObjName := m.generateObjName(objName)

	putOpt := minio.PutObjectOptions{
		UserTags: map[string]string{
			"type": sourceType.String(),
			"name": objName,
		},
	}

	newObjName = fmt.Sprintf("%s/%s", sourceType.String(), newObjName)
	_, err = m.client.PutObject(ctx, m.Bucket, newObjName, obj, size, putOpt)
	if err != nil {
		return "", errors.Wrap(err, "uploadMinio")
	}

	// https://billiard.superwhys.top/api/v1/minio/avatar/1731850656800-d887240b-0177-44c7-853d-69f14b7cf874.jpeg
	return fmt.Sprintf("%s/minio/%s", m.baseApi, newObjName), nil
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
