// File:		auth.go
// Created by:	Hoven
// Created on:	2025-01-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package server

import (
	"context"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/auth"
	"gitea.hoven.com/core/auth-core/domain/authentication"
	"github.com/go-puzzles/puzzles/plog"
)

func (s *BilliardServer) Logout(ctx context.Context, token string) error {
	err := s.AuthSrv.Logout(ctx, token)
	if err != nil {
		plog.Errorc(ctx, "logout error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) WechatLogin(ctx context.Context, device, code string) (*auth.Token, error) {
	resp, err := s.AuthSrv.WechatLogin(ctx, device, code)
	if err != nil {
		plog.Errorc(ctx, "wechat login error: %v", err)
		return nil, err
	}

	return resp, nil
}

func (s *BilliardServer) BindEmail(ctx context.Context, token string, email, code string) error {
	authPair := &auth.AuthPair{
		AuthType:       authentication.AuthTypeEmail,
		CredentialType: authentication.CredentialEmailCode,
	}

	identifierPair := &auth.IdentifierPair{
		Identifier: email,
		Credential: code,
	}
	err := s.AuthSrv.BindAuth(ctx, token, authPair, identifierPair)
	if err != nil {
		plog.Errorc(ctx, "bind email error: %v", err)
		return err
	}

	return nil
}

func (s *BilliardServer) SendEmailCode(ctx context.Context, email string) error {
	err := s.AuthSrv.SendEmailCode(ctx, email)
	if err != nil {
		plog.Errorc(ctx, "send email code error: %v", err)
		return err
	}

	return nil
}
