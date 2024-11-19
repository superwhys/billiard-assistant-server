// File:		auth.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package authSrv

import (
	"context"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/auth"
)

var _ auth.IAuthService = (*AuthService)(nil)

type AuthService struct {
	authRepo auth.IAuthRepo
}

func NewAuthService(authRepo auth.IAuthRepo) auth.IAuthService {
	return &AuthService{authRepo: authRepo}
}

func (as *AuthService) CreateUserAuth(ctx context.Context, userId int, auth *auth.Auth) error {
	return as.authRepo.CreateUserAuth(ctx, userId, auth)
}

func (as *AuthService) UpdateUserAuth(ctx context.Context, auth *auth.Auth) error {
	return as.authRepo.UpdateUserAuth(ctx, auth)
}

func (as *AuthService) DeleteUserAuth(ctx context.Context, authId int) error {
	return as.authRepo.DeleteUserAuth(ctx, authId)
}

func (as *AuthService) GetUserAuthByType(ctx context.Context, userId int, authType auth.AuthType) (*auth.Auth, error) {
	return as.authRepo.GetUserAuthByType(ctx, userId, authType)
}

func (as *AuthService) GetUserAuths(ctx context.Context, userId int) ([]*auth.Auth, error) {
	return as.authRepo.GetUserAuths(ctx, userId)
}

func (as *AuthService) GetUserAuthByIdentifier(ctx context.Context, authType auth.AuthType, identifier string) (*auth.Auth, error) {
	return as.authRepo.GetUserAuthByIdentifier(ctx, authType, identifier)
}
