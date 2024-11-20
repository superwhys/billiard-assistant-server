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

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/notice"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/user"
)

type User struct {
	UserId    int    `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Status    int    `json:"status,omitempty"`
	Role      int    `json:"role,omitempty"`
	AuthTypes []int  `json:"auth_types,omitempty"`
	IsAdmin   bool   `json:"is_admin,omitempty"`
}

func UserEntityToDto(u *user.User) *User {
	var authTypes []int

	user := &User{
		UserId:    u.UserId,
		Name:      u.Name,
		Status:    int(u.Status),
		Role:      int(u.Role),
		AuthTypes: authTypes,
		IsAdmin:   u.IsAdmin(),
	}

	if u.UserInfo != nil {
		user.Email = u.UserInfo.Email
		user.Phone = u.UserInfo.Phone
		user.Avatar = u.UserInfo.Avatar
	}

	return user
}

func UserDtoToEntity(u *User) *user.User {
	userEntity := &user.User{
		UserId: u.UserId,
		Name:   u.Name,
		UserInfo: &user.BaseInfo{
			Email:  u.Email,
			Phone:  u.Phone,
			Avatar: u.Avatar,
		},
		Status: user.Status(u.Status),
		Role:   user.Role(u.Role),
	}

	return userEntity
}

type WechatLoginRequest struct {
	Code string `json:"code"`
}

type WechatLoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
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

type UpdateUserNameRequest struct {
	UserName string `json:"username"`
}

type UploadAvatarResponse struct {
	AvatarUrl string `json:"avatar_url"`
}

type GetUserAvatarRequest struct {
	AvatarName string `uri:"avatar_name"`
}

type BindPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type BindEmailRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type SendPhoneCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type SendEmailCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

type CheckPhoneBindRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type CheckPhoneBindResponse struct {
	IsBound bool `json:"is_bound"`
}

type AddNoticeRequest struct {
	NoticeType notice.NoticeType `json:"notice_type" binding:"required"`
	Contents   []string          `json:"contents" binding:"required"`
}

type UploadGameIconResponse struct {
	IconUrl string `json:"icon_url"`
}
