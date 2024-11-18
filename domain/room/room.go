// File:		room.go
// Created by:	Hoven
// Created on:	2024-11-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package room

import (
	"github.com/superwhys/billiard-assistant-server/domain/shared"
)

type User interface {
	shared.BaseUser
}

type Game interface {
	shared.BaseGame
	GetMaxPlayers() int
	GetGameType() shared.BilliardGameType
}

type Player struct {
	User
	Prepared bool
}

type Room struct {
	RoomId        int
	GameId        int
	OwnerId       int
	Players       []Player
	Game          Game
	GameStatus    Status
	WinLoseStatus WinLoseStatus
}

func (r *Room) GetRoomId() int {
	return r.RoomId
}

func (r *Room) PlayerIds() []int {
	ids := make([]int, len(r.Players))
	for i, p := range r.Players {
		ids[i] = p.GetUserId()
	}
	return ids
}

func (r *Room) CanStart() bool {
	for _, p := range r.Players {
		if !p.Prepared {
			return false
		}
	}
	
	return true
}

func (r *Room) CanEnter() bool {
	return len(r.Players) < r.Game.GetMaxPlayers()
}

func (r *Room) IsOwner(userId int) bool {
	return r.OwnerId == userId
}

func (r *Room) IsInRoom(userId int) bool {
	for _, p := range r.Players {
		if p.GetUserId() == userId {
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
	Unknown = iota
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
	Preparing Status = iota
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
