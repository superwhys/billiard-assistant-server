// File:		auth.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package model

import (
	"time"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/auth"
	"gorm.io/gorm"
)

type UserAuthPo struct {
	ID int `gorm:"primaryKey"`

	UserPoID   int           `gorm:"index"`
	AuthType   auth.AuthType `gorm:"uniqueIndex:idx_auth_type_identifier"`
	Identifier string        `gorm:"size:255;uniqueIndex:idx_auth_type_identifier"`
	Credential string        `gorm:"size:255"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ua *UserAuthPo) TableName() string {
	return "user_auths"
}

func (ua *UserAuthPo) FromEntity(entity *auth.Auth) {
	ua.ID = entity.Id
	ua.UserPoID = entity.UserId
	ua.AuthType = entity.AuthType
	ua.Identifier = entity.Identifier
	ua.Credential = entity.Credential
}

func (ua *UserAuthPo) ToEntity() *auth.Auth {
	if ua == nil {
		return nil
	}

	return &auth.Auth{
		Id:         ua.ID,
		UserId:     ua.UserPoID,
		AuthType:   ua.AuthType,
		Identifier: ua.Identifier,
		Credential: ua.Credential,
	}
}
