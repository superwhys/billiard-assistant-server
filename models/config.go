// File:		config.go
// Created by:	Hoven
// Created on:	2024-10-28
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package models

import (
	"github.com/go-puzzles/pgorm"
	"github.com/go-puzzles/predis"
	"github.com/go-puzzles/puzzles/plog"
)

type SaConfig struct {
	WxAppId     string
	WxAppSecret string
	AvatarDir   string
}

type parser func(out any) error

func ParseConfig(saConfParser, redisConfParser, mysqlConfParser parser) (*SaConfig, *predis.RedisConf, *pgorm.MysqlConfig) {
	saConfig := new(SaConfig)
	plog.PanicError(saConfParser(saConfig))

	redisConf := new(predis.RedisConf)
	plog.PanicError(redisConfParser(redisConf))

	mysqlConf := new(pgorm.MysqlConfig)
	plog.PanicError(mysqlConfParser(mysqlConf))

	return saConfig, redisConf, mysqlConf
}
