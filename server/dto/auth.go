// File:		auth.go
// Created by:	Hoven
// Created on:	2025-01-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import "errors"

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (req *RegisterRequest) Validate() error {
	if req.Username == "" {
		return errors.New("missing account")
	}

	if req.Password == "" {
		return errors.New("missing password")
	}

	return nil
}

type WechatLoginRequest struct {
	DeviceId string `headers:"deviceId" binding:"required"`
	Code     string `json:"code"`
}

type LoginRequest struct {
	DeviceId string `headers:"deviceId" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}
