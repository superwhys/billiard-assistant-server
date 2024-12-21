// File:		snooker.go
// Created by:	Hoven
// Created on:	2024-12-21
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package snooker

import "time"

type BallType string

const (
	RED    BallType = "red"
	YELLOW BallType = "yellow"
	GREEN  BallType = "green"
	BROWN  BallType = "brown"
	BLUE   BallType = "blue"
	PINK   BallType = "pink"
	BLACK  BallType = "black"
)

type ActionItem struct {
	PlayerIndex int      `json:"playerIndex"`
	PlayerName  string   `json:"playerName"`
	ScoreDelta  int      `json:"scoreDelta"`
	BallType    BallType `json:"ballType"`
}

type SnookerAction struct {
	UserId       int         `json:"userId"`
	RoomId       int         `json:"roomId"`
	UserName     string      `json:"userName"`
	History      *ActionItem `json:"history"`
	IsRevocation bool        `json:"isRevocation"`
	ActionTime   time.Time   `json:"actionTime"`
}

func (sa *SnookerAction) GetActionRoomId() (roomId int) {
	return sa.RoomId
}

func (sa *SnookerAction) GetActionUser() (userId int) {
	return sa.UserId
}

func (sa *SnookerAction) GetActionTime() time.Time {
	return sa.ActionTime
}

type SnookerPlayerRecord struct {
	RoomId      int    `json:"roomId"`
	Name        string `json:"name"`
	MainScore   int    `json:"mainScore"`
	FrameScore  int    `json:"frameScore"`
	TotalScore  int    `json:"totalScore"`
	IsRoomOwner bool   `json:"isRoomOwner"`
	IsNew       bool   `json:"isNew"`
	IsOnline    bool   `json:"isOnline"`
	IsVirtual   bool   `json:"isVirtual"`
}

func (sr *SnookerPlayerRecord) GetRecordRoomId() (roomId int) {
	return sr.RoomId
}
