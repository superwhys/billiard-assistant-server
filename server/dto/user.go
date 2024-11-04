// File:		game.go
// Created by:	Hoven
// Created on:	2024-10-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dto

import (
	"errors"
	"time"
	
	"github.com/superwhys/snooker-assistant-server/domain/user"
)

type User struct {
	UserId      int       `json:"user_id"`
	Name        string    `json:"name"`
	WechatId    string    `json:"wechat_id"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Avatar      string    `json:"avatar"`
	Status      int       `json:"status"`
	LastLoginAt time.Time `json:"last_login_at"`
}

func UserEntityToDto(u *user.User) *User {
	user := &User{
		UserId:   u.UserId,
		Name:     u.Name,
		WechatId: u.WechatId,
		Status:   int(u.Status),
	}
	
	if u.UserInfo != nil {
		user.Email = u.UserInfo.Email
		user.Phone = u.UserInfo.Phone
		user.Avatar = u.UserInfo.Avatar
	}
	
	return user
}

func UserDtoToEntity(u *User) *user.User {
	return &user.User{
		UserId:   u.UserId,
		Name:     u.Name,
		WechatId: u.WechatId,
		UserInfo: &user.BaseInfo{
			Email:  u.Email,
			Phone:  u.Phone,
			Avatar: u.Avatar,
		},
		Status: user.Status(u.Status),
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	WechatId string `json:"wechat_id"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
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

type RegisterResponse struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
}

type GetUserInfoRequest struct {
	UserId int `json:"userId" uri:"userId"`
}

type GetUserInfoResponse struct {
}

type UpdateUserRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type UploadAvatarResponse struct {
	AvatarUrl string `json:"avatar_url"`
}
