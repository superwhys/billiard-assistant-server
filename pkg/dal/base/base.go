// File:		base.go
// Created by:	Hoven
// Created on:	2024-10-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package base

import (
	"gorm.io/gorm"
)

type GameDB struct {
	gamePo
}

func NewGameDB(db *gorm.DB) *GameDB {
	return &GameDB{
		gamePo: newGamePo(db),
	}
}

type RoomDB struct {
	roomPo
}

func NewRoomDB(db *gorm.DB) *RoomDB {
	return &RoomDB{
		roomPo: newRoomPo(db),
	}
}

type UserDB struct {
	userPo
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{
		userPo: newUserPo(db),
	}
}

type NoticeDB struct {
	noticePo
}

func NewNoticeDB(db *gorm.DB) *NoticeDB {
	return &NoticeDB{
		noticePo: newNoticePo(db),
	}
}
