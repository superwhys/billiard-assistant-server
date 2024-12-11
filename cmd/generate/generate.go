// File:		gen.go
// Created by:	Hoven
// Created on:	2024-10-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gorm.io/gen"
)

//go:generate go run generate.go
func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../pkg/dal/base",
		WithUnitTest:  false,
		FieldNullable: true,
		Mode:          gen.WithQueryInterface,
	})

	// 直接使用模型
	g.ApplyBasic(
		&model.UserPo{},
		&model.UserAuthPo{},
		&model.RoomPo{},
		&model.GamePo{},
		&model.NoticePo{},
		&model.RoomUserPo{},
		&model.RecordPo{},
	)

	g.Execute()
}
