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

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/auth"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitea.hoven.com/core/auth-core/api/grpc/interceptor"
	"gitea.hoven.com/core/auth-core/api/grpc/proto/authenticationpb"
	"gitea.hoven.com/core/auth-core/api/grpc/proto/verifycodepb"
	"gitea.hoven.com/core/auth-core/server/dto"
	"google.golang.org/grpc/metadata"
)

var _ auth.IAuthService = (*AuthService)(nil)

type AuthService struct {
	authClient       authenticationpb.AuthCoreAuthenticationHandlerClient
	verifycodeClient verifycodepb.AuthCoreVerifyCodeHandlerClient
}

func NewAuthService(
	authClient authenticationpb.AuthCoreAuthenticationHandlerClient,
	verifycodeClient verifycodepb.AuthCoreVerifyCodeHandlerClient,
) auth.IAuthService {
	return &AuthService{
		authClient:       authClient,
		verifycodeClient: verifycodeClient,
	}
}

func (as *AuthService) injectToken(ctx context.Context, token string) context.Context {
	md := metadata.Pairs(interceptor.TokenMetadataKey, token)
	return metadata.NewOutgoingContext(ctx, md)
}

func (as *AuthService) AccountLogin(ctx context.Context, device, username, password string) (*auth.Token, error) {
	resp, err := as.authClient.AccountLogin(ctx, &dto.AccountLoginRequest{
		DeviceId: device,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, exception.ParseGrpcError(err)
	}

	return &auth.Token{
		UserId:       int(resp.GetUserId()),
		AccessToken:  resp.GetToken().AccessToken,
		RefreshToken: resp.GetToken().RefreshToken,
	}, nil
}

func (as *AuthService) AccountRegister(ctx context.Context, username, password string) error {
	_, err := as.authClient.AccountRegister(ctx, &dto.AccountRegisterRequest{
		IdentifierPair: &dto.IdentifierPair{
			Identifier: username,
			Credential: password,
		},
	})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	return nil
}

func (as *AuthService) WechatLogin(ctx context.Context, device string, code string) (*auth.Token, error) {
	resp, err := as.authClient.WechatLogin(ctx, &dto.WechatLoginRequest{
		DeviceId: device,
		Code:     code,
	})
	if err != nil {
		return nil, exception.ParseGrpcError(err)
	}

	return &auth.Token{
		UserId:       int(resp.GetUserId()),
		AccessToken:  resp.GetToken().AccessToken,
		RefreshToken: resp.GetToken().RefreshToken,
	}, nil
}

func (as *AuthService) Logout(ctx context.Context, token string) error {
	ctx = as.injectToken(ctx, token)
	_, err := as.authClient.Logout(ctx, &dto.LogoutRequest{})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	return nil
}

func (as *AuthService) BindAuth(ctx context.Context, token string, authPair *auth.AuthPair, identifierPair *auth.IdentifierPair) error {
	ctx = as.injectToken(ctx, token)

	_, err := as.authClient.BindAuth(ctx, &dto.BindAuthRequest{
		AuthPair:       authPair.ToAuthCoreAuthPair(),
		IdentifierPair: identifierPair.ToAuthCoreIdentifierPair(),
	})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	return nil
}

func (as *AuthService) UnbindAuth(ctx context.Context, token string, authPair *auth.AuthPair, identifierPair *auth.IdentifierPair) error {
	ctx = as.injectToken(ctx, token)

	_, err := as.authClient.UnbindAuth(ctx, &dto.UnbindAuthRequest{
		AuthPair:       authPair.ToAuthCoreAuthPair(),
		IdentifierPair: identifierPair.ToAuthCoreIdentifierPair(),
	})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	return nil
}

func (as *AuthService) SendEmailCode(ctx context.Context, email string) error {
	_, err := as.verifycodeClient.SendEmailVerifyCode(ctx, &dto.SendEmailVerifyCodeRequest{Email: email})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	return nil
}
