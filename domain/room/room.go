// File:		room.go
// Created by:	Hoven
// Created on:	2024-11-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package room

import (
	"time"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gorm.io/datatypes"
)

type Room struct {
	RoomId        int
	RoomCode      string
	GameId        int
	OwnerId       int
	Players       []shared.RoomPlayer
	Owner         shared.BaseUser
	Game          shared.BaseGame
	Record        shared.BaseRecord
	Extra         map[string]any
	GameStatus    Status
	WinLoseStatus WinLoseStatus
	CreateAt      time.Time
}

type RoomPlayer struct {
	RoomId          int
	UserId          int
	UserName        string
	IsVirtualPlayer bool
	HeartbeatAt     time.Time
}

func (r *Room) CanStart() bool {
	return r.GameStatus == Preparing
}

func (r *Room) CanEnter() bool {
	return len(r.Players) < r.Game.GetMaxPlayers()
}

func (r *Room) IsOwner(userId int) bool {
	return r.OwnerId == userId
}

func (r *Room) IsEnd() bool {
	return r.GameStatus == Finish
}

func (r *Room) IsInRoom(currentUserId int, virtualUser string) bool {
	if currentUserId == 0 && virtualUser == "" {
		return true
	}

	for _, p := range r.Players {
		if virtualUser != "" && p.GetUserName() == virtualUser {
			return true
		}

		if virtualUser == "" && p.GetUserId() == currentUserId {
			return true
		}
	}

	return false
}

func (r *Room) StartGame() {
	r.GameStatus = Playing
}

func (r *Room) SetExtra(extra datatypes.JSONMap) {
	r.Extra = extra
}

func (r *Room) EndGame() {
	r.GameStatus = Finish
}

type WinLoseStatus int

const (
	WinLoseUnknown = iota
	Win
	Lose
	Tie
)

func (gt WinLoseStatus) String() string {
	switch gt {
	case Win:
		return "胜利"
	case Lose:
		return "失败"
	case Tie:
		return "平局"
	default:
		return "未知"
	}
}

type Status int

const (
	StatusUnknown Status = iota
	Preparing
	Playing
	Finish
)

func (s Status) String() string {
	switch s {
	case Preparing:
		return "准备中"
	case Playing:
		return "进行中"
	case Finish:
		return "已结束"
	default:
		return "未知"
	}
}
