// File:		nineball.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package nineball

import (
	"encoding/json"
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
	PlayerIndex int
	PlayerName  string
	ScoreDelta  int
	ScoreType   ScoreType
	StateDelta  int
}

type RecordItem struct {
	Changes               []*ScoreChange
	PreviousCurrentPlayer int
}

type NineballRecord struct {
	UserId     int
	UserName   string
	History    *RecordItem
	ReportTime time.Time
}

func (nr *NineballRecord) UnmarshalFrom(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), nr)
}

func (nr *NineballRecord) GetRecordTime() time.Time {
	return nr.ReportTime
}

func (nr *NineballRecord) GetRecordUser() (userId int) {
	return nr.UserId
}
