// File:		config.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package models

import (
	"errors"

	"github.com/go-puzzles/pgorm"
	"github.com/go-puzzles/predis"
	"github.com/go-puzzles/puzzles/plog"
)

type SaConfig struct {
	WxAppId     string
	WxAppSecret string
	AvatarDir   string
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string

	Bucket string
}

func (c *MinioConfig) Validate() error {
	if c.Endpoint == "" || c.AccessKey == "" || c.SecretKey == "" || c.Bucket == "" {
		return errors.New("invalid minio config")
	}

	return nil
}

type parser func(out any) error

func ParseConfig(saConfParser, redisConfParser, mysqlConfParser, minioConfParaser parser) (*SaConfig, *predis.RedisConf, *pgorm.MysqlConfig, *MinioConfig) {
	saConfig := new(SaConfig)
	plog.PanicError(saConfParser(saConfig))

	redisConf := new(predis.RedisConf)
	plog.PanicError(redisConfParser(redisConf))

	mysqlConf := new(pgorm.MysqlConfig)
	plog.PanicError(mysqlConfParser(mysqlConf))

	minioConf := new(MinioConfig)
	plog.PanicError(minioConfParaser(minioConf))

	return saConfig, redisConf, mysqlConf, minioConf
}
