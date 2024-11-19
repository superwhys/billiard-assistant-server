// File:		repo.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package auth

import "context"

type IAuthRepo interface {
	CreateUserAuth(ctx context.Context, userId int, auth *Auth) error
	UpdateUserAuth(ctx context.Context, auth *Auth) error
	DeleteUserAuth(ctx context.Context, authId int) error
	GetUserAuths(ctx context.Context, userId int) ([]*Auth, error)
	GetUserAuthByType(ctx context.Context, userId int, authType AuthType) (*Auth, error)
	GetUserAuthByIdentifier(ctx context.Context, authType AuthType, identifier string) (*Auth, error)
}
