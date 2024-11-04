// File:		gen.go
// Created by:	Hoven
// Created on:	2024-10-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package main

import (
	"github.com/superwhys/snooker-assistant-server/pkg/dal/model"
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
	
	g.ApplyBasic(
		model.UserPo{},
		model.GamePo{},
		model.RoomPo{},
		model.NoticePo{},
	)
	
	g.Execute()
}
