// File:		report.go
// Created by:	Hoven
// Created on:	2024-11-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package recordDal

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/record"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/base"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ record.IRecordRepo = (*RecordRepoImpl)(nil)

type RecordRepoImpl struct {
	db *base.Query
}

func NewRecordRepo(db *gorm.DB) *RecordRepoImpl {
	ri := &RecordRepoImpl{
		db: base.Use(db),
	}

	return ri
}

func (r *RecordRepoImpl) GetRecordByRoomId(ctx context.Context, roomId int, recordTmpl reflect.Type) (*record.Record, error) {
	recordPo := r.db.RecordPo
	resp, err := recordPo.WithContext(ctx).Where(recordPo.RoomID.Eq(roomId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrRoomRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return resp.ToEntity(recordTmpl), nil
}

func (r *RecordRepoImpl) UpdateRoomRecord(ctx context.Context, roomId int, gameType shared.BilliardGameType, record []record.RecordItem) error {
	recordPo := r.db.RecordPo

	recordB, err := json.Marshal(record)
	if err != nil {
		return errors.Wrap(err, "marshalRecord")
	}

	return recordPo.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"data"}),
		},
	).Create(&model.RecordPo{
		RoomID:   roomId,
		Data:     datatypes.JSON(recordB),
		GameType: gameType,
	})
}
