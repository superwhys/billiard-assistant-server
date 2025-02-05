// File:		auth.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package auth

import (
	"gitea.hoven.com/core/auth-core/domain/authentication"
	"gitea.hoven.com/core/auth-core/server/dto"
)

type AuthPair struct {
	AuthType       authentication.AuthType
	CredentialType authentication.CredentialType
}

func (ap *AuthPair) ToAuthCoreAuthPair() *dto.AuthPair {
	return &dto.AuthPair{
		AuthType:       ap.AuthType,
		CredentialType: ap.CredentialType,
	}
}

type IdentifierPair struct {
	// username or phone or wechatId or email
	Identifier string
	// password or something else (can be empty)
	Credential string
}

func (ip *IdentifierPair) ToAuthCoreIdentifierPair() *dto.IdentifierPair {
	return &dto.IdentifierPair{
		Identifier: ip.Identifier,
		Credential: ip.Credential,
	}
}

type Token struct {
	UserId       int
	AccessToken  string
	RefreshToken string
}
