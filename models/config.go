// File:		config.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package models

import (
	"time"

	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/oss/minio"
	"github.com/go-puzzles/puzzles/goredis"
	"github.com/go-puzzles/puzzles/pgorm"
	"github.com/go-puzzles/puzzles/plog"
)

type RoomConfig struct {
	UserMaxRoomCreateNumber int64
}

func (r *RoomConfig) SetDefault() {
	if r.UserMaxRoomCreateNumber == 0 {
		r.UserMaxRoomCreateNumber = 20
	}
}

type Config struct {
	BaseApi     string
	TokenPrefix string
	TokenTtl    time.Duration
	RoomConfig  *RoomConfig
}

func (c *Config) SetDefault() {
	if c.TokenPrefix == "" {
		c.TokenPrefix = "billiard"
	}

	if c.TokenTtl == 0 {
		c.TokenTtl = time.Hour
	}

	if c.RoomConfig == nil {
		c.RoomConfig = new(RoomConfig)
		c.RoomConfig.SetDefault()
	}
}

type parser func(out any) error

type Configs struct {
	SrvConf   *Config
	RedisConf *goredis.RedisConf
	MysqlConf *pgorm.MysqlConfig
	MinioConf *minio.MinioConfig
}

func ParseConfig(srvConfParser, redisConfParser, mysqlConfParser, minioConfParser parser) *Configs {
	srvConfig := new(Config)
	plog.PanicError(srvConfParser(srvConfig))

	redisConf := new(goredis.RedisConf)
	plog.PanicError(redisConfParser(redisConf))

	mysqlConf := new(pgorm.MysqlConfig)
	plog.PanicError(mysqlConfParser(mysqlConf))

	minioConf := new(minio.MinioConfig)
	plog.PanicError(minioConfParser(minioConf))

	return &Configs{
		SrvConf:   srvConfig,
		RedisConf: redisConf,
		MysqlConf: mysqlConf,
		MinioConf: minioConf,
	}
}
