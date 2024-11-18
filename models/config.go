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
	
	"github.com/go-puzzles/puzzles/pgorm"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/predis"
)

type Config struct {
	AvatarDir string
	UserApi   string
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

func ParseConfig(srvConfParser, redisConfParser, mysqlConfParser, minioConfParser parser) (*Config, *predis.RedisConf, *pgorm.MysqlConfig, *MinioConfig) {
	srvConfig := new(Config)
	plog.PanicError(srvConfParser(srvConfig))
	
	redisConf := new(predis.RedisConf)
	plog.PanicError(redisConfParser(redisConf))
	
	mysqlConf := new(pgorm.MysqlConfig)
	plog.PanicError(mysqlConfParser(mysqlConf))
	
	minioConf := new(MinioConfig)
	plog.PanicError(minioConfParser(minioConf))
	
	return srvConfig, redisConf, mysqlConf, minioConf
}
