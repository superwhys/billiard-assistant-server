// File:		nineball.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package nineball

import (
	"time"
)

type ScoreType string

const (
	FOUL        ScoreType = "foul"
	NORMAL      ScoreType = "normal"
	BIG_GOLD    ScoreType = "bigGold"
	SMALL_GOLD  ScoreType = "smallGold"
	GOLDEN_NINE ScoreType = "goldenNine"
)

type ScoreChange struct {
	PlayerIndex int       `json:"player_index"`
	PlayerName  string    `json:"player_name"`
	ScoreDelta  int       `json:"score_delta"`
	ScoreType   ScoreType `json:"score_type"`
	StateDelta  int       `json:"state_delta"`
}

type ActionItem struct {
	Changes               []*ScoreChange `json:"changes"`
	PreviousCurrentPlayer int            `json:"previous_current_player"`
}

type NineballAction struct {
	UserId     int         `json:"user_id"`
	RoomId     int         `json:"room_id"`
	UserName   string      `json:"user_name"`
	History    *ActionItem `json:"history"`
	ActionTime time.Time   `json:"action_time"`
}

func (nr *NineballAction) GetActionRoomId() (roomId int) {
	return nr.RoomId
}

func (nr *NineballAction) GetActionUser() (userId int) {
	return nr.UserId
}

func (nr *NineballAction) GetActionTime() time.Time {
	return nr.ActionTime
}

type ScoreStats struct {
	Fouls      int `json:"foulCount"`
	Normal     int `json:"normalCount"`
	BigGold    int `json:"bigGoldCount"`
	SmallGold  int `json:"smallGoldCount"`
	GoldenNine int `json:"goldenNineCount"`
}

type PlayerRecord struct {
	RoomId      int         `json:"room_id"`
	Name        string      `json:"name"`
	Score       int         `json:"score"`
	Stats       *ScoreStats `json:"stats"`
	IsRoomOwner bool        `json:"isRoomOwner"`
	IsNew       bool        `json:"isNew"`
	IsOnline    bool        `json:"isOnline"`
	IsVirtual   bool        `json:"isVirtual"`
}

func (pr *PlayerRecord) GetRecordRoomId() (roomId int) {
	return pr.RoomId
}
