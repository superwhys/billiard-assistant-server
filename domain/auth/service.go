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
	WechatLogin(ctx context.Context, device string, code string) (*Token, error)
	Logout(ctx context.Context, token string) error
	BindAuth(ctx context.Context, token string, authPair *AuthPair, identifierPair *IdentifierPair) error
	UnbindAuth(ctx context.Context, token string, authPair *AuthPair, identifierPair *IdentifierPair) error
	SendEmailCode(ctx context.Context, email string) error
}
