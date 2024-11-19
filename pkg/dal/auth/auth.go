// File:		auth.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package authDal

import (
	"context"
	"errors"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/auth"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/base"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gorm.io/gorm"
)

var _ auth.IAuthRepo = (*AuthRepoImpl)(nil)

type AuthRepoImpl struct {
	db *base.Query
}

func NewAuthRepo(db *gorm.DB) *AuthRepoImpl {
	return &AuthRepoImpl{db: base.Use(db)}
}

func (ar *AuthRepoImpl) CreateUserAuth(ctx context.Context, userId int, auth *auth.Auth) error {
	userDb := ar.db.UserPo
	userAuthDb := ar.db.UserAuthPo

	userAuthPo := new(model.UserAuthPo)
	userAuthPo.FromEntity(auth)
	userAuthPo.UserPoID = userId
	if err := userAuthDb.WithContext(ctx).Create(userAuthPo); err != nil {
		return err
	}

	return userDb.UserAuthPos.Model(&model.UserPo{ID: userId}).Append(userAuthPo)
}

func (ar *AuthRepoImpl) UpdateUserAuth(ctx context.Context, auth *auth.Auth) error {
	userAuthPo := new(model.UserAuthPo)
	userAuthPo.FromEntity(auth)

	userAuthDb := ar.db.UserAuthPo
	_, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.ID.Eq(auth.Id)).
		Updates(userAuthPo)
	return err
}

func (ar *AuthRepoImpl) DeleteUserAuth(ctx context.Context, authId int) error {
	userAuthDb := ar.db.UserAuthPo
	_, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.ID.Eq(authId)).
		Delete()
	return err
}

func (ar *AuthRepoImpl) GetUserAuths(ctx context.Context, userId int) ([]*auth.Auth, error) {
	userAuthDb := ar.db.UserAuthPo
	auths, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.UserPoID.Eq(userId)).
		Find()
	if err != nil {
		return nil, err
	}

	result := make([]*auth.Auth, 0, len(auths))
	for _, auth := range auths {
		result = append(result, auth.ToEntity())
	}
	return result, nil
}

func (ar *AuthRepoImpl) GetUserAuthByType(ctx context.Context, userId int, authType auth.AuthType) (*auth.Auth, error) {
	userAuthDb := ar.db.UserAuthPo
	auth, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.UserPoID.Eq(userId)).
		Where(userAuthDb.AuthType.Eq(int(authType))).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserAuthNotFound
	} else if err != nil {
		return nil, err
	}

	return auth.ToEntity(), nil
}

func (ar *AuthRepoImpl) GetUserAuthByIdentifier(ctx context.Context, authType auth.AuthType, identifier string) (*auth.Auth, error) {
	userAuthDb := ar.db.UserAuthPo
	auth, err := userAuthDb.WithContext(ctx).
		Where(userAuthDb.AuthType.Eq(int(authType))).
		Where(userAuthDb.Identifier.Eq(identifier)).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrUserAuthNotFound
	} else if err != nil {
		return nil, err
	}

	return auth.ToEntity(), nil
}
