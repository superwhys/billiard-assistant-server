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
	"time"

	"github.com/go-puzzles/puzzles/pgorm"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/predis"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/email"
)

type Config struct {
	AvatarDir   string
	UserApi     string
	TokenPrefix string
	TokenTtl    time.Duration
}

func (c *Config) SetDefault() {
	if c.TokenPrefix == "" {
		c.TokenPrefix = "billiard"
	}

	if c.TokenTtl == 0 {
		c.TokenTtl = time.Hour
	}
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

type Configs struct {
	SrvConf   *Config
	RedisConf *predis.RedisConf
	MysqlConf *pgorm.MysqlConfig
	MinioConf *MinioConfig
	EmailConf *email.EmailConf
}

func ParseConfig(srvConfParser, redisConfParser, mysqlConfParser, minioConfParser, emailConfParser parser) *Configs {
	srvConfig := new(Config)
	plog.PanicError(srvConfParser(srvConfig))

	redisConf := new(predis.RedisConf)
	plog.PanicError(redisConfParser(redisConf))

	mysqlConf := new(pgorm.MysqlConfig)
	plog.PanicError(mysqlConfParser(mysqlConf))

	minioConf := new(MinioConfig)
	plog.PanicError(minioConfParser(minioConf))

	emailConf := new(email.EmailConf)
	plog.PanicError(emailConfParser(emailConf))

	return &Configs{
		SrvConf:   srvConfig,
		RedisConf: redisConf,
		MysqlConf: mysqlConf,
		MinioConf: minioConf,
		EmailConf: emailConf,
	}
}
