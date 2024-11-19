// File:		service.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package auth

import "context"

type IAuthService interface {
	// Auth management
	CreateUserAuth(ctx context.Context, userId int, auth *Auth) error
	UpdateUserAuth(ctx context.Context, auth *Auth) error
	DeleteUserAuth(ctx context.Context, authId int) error
	GetUserAuths(ctx context.Context, userId int) ([]*Auth, error)
	GetUserAuthByIdentifier(ctx context.Context, authType AuthType, identifier string) (*Auth, error)
}
