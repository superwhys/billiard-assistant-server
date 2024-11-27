// File:		room_test.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package roomDal

import (
	"context"
	"testing"

	"github.com/go-puzzles/puzzles/pgorm"
	"github.com/stretchr/testify/assert"
)

func TestRoomRepo(t *testing.T) {
	/*
		instance: localhost:3306
		database: snooker
		username: root
		password: yang4869
	*/
	mysqlConf := &pgorm.MysqlConfig{
		Instance: "localhost:3306",
		Database: "snooker",
		Username: "root",
		Password: "yang4869",
	}

	db, err := mysqlConf.DialGorm()
	if err != nil {
		t.Fatal(err)
	}

	gameRepo := NewRoomRepo(db)

	t.Run("testGetRoomType", func(t *testing.T) {
		gt, err := gameRepo.GetRoomGameType(context.Background(), 8)
		if !assert.Nil(t, err) {
			return
		}

		t.Logf("game type: %v", gt)
	})
}
