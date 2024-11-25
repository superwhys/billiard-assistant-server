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

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
)

type User interface {
	shared.BaseUser
}

type Game interface {
	shared.BaseGame
	GetMaxPlayers() int
	GetGameType() shared.BilliardGameType
	GetIcon() string
}

type Player struct {
	UserId          int
	UserName        string
	IsVirtualPlayer bool
}

type Room struct {
	RoomId        int
	GameId        int
	OwnerId       int
	Players       []*Player
	Game          Game
	GameStatus    Status
	WinLoseStatus WinLoseStatus
	CreateAt      time.Time
}

func (r *Room) GetRoomId() int {
	return r.RoomId
}

func (r *Room) CanStart() bool {
	return true
}

func (r *Room) CanEnter() bool {
	return len(r.Players) < r.Game.GetMaxPlayers()
}

func (r *Room) IsOwner(userId int) bool {
	return r.OwnerId == userId
}

func (r *Room) IsInRoom(virtualName string, userId int) bool {
	if virtualName == "" && userId == 0 {
		return false
	}

	for _, p := range r.Players {
		if p.IsVirtualPlayer && p.UserName == virtualName {
			return true
		}

		if p.UserId == userId {
			return true
		}
	}

	return false
}

func (r *Room) StartGame() {
	r.GameStatus = Playing
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
