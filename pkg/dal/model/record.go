// File:		record.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package model

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/go-puzzles/puzzles/plog"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/record"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RecordPo struct {
	ID int `gorm:"primaryKey"`

	RoomID   int     `gorm:"unique"`
	Room     *RoomPo `gorm:"foreignKey:RoomID"`
	GameType shared.BilliardGameType

	Data datatypes.JSON

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (r *RecordPo) TableName() string {
	return "records"
}

func (r *RecordPo) ToEntity(t reflect.Type) *record.Record {
	if r == nil {
		return nil
	}
	ptrType := reflect.PointerTo(t)
	sliceType := reflect.SliceOf(ptrType)
	sliceValue := reflect.New(sliceType).Interface()

	if err := json.Unmarshal(r.Data, sliceValue); err != nil {
		plog.Errorf("json.Unmarshal recordItem slice(%v) error: %v", t.Name(), err)
		return nil
	}

	slice := reflect.ValueOf(sliceValue).Elem()

	records := &record.Record{
		ID:            r.ID,
		RoomId:        r.RoomID,
		CurrentRecord: make([]record.RecordItem, slice.Len()),
	}

	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i).Interface().(record.RecordItem)
		records.CurrentRecord[i] = item
	}

	return records
}
