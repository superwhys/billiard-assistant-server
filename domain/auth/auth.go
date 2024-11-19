// File:		auth.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package auth

type Auth struct {
	Id       int
	UserId   int
	AuthType AuthType
	// username or phone or wechatId or email
	Identifier string
	// password or something else (can be empty)
	Credential string
}

type AuthType int

const (
	AuthTypeUnknown AuthType = iota
	AuthTypeWechat
	AuthTypeEmail
	AuthTypePhone
	AuthTypePassword
)
