// File:		room.go
// Created by:	Hoven
// Created on:	2024-11-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package room

type Game interface {
	GetMaxPlayers() int
}

type Player interface {
	GetUserId() int
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
	Playing Status = iota
	Finish
)

func (s Status) String() string {
	switch s {
	case Playing:
		return "进行中"
	case Finish:
		return "已结束"
	default:
		return "未知"
	}
}
