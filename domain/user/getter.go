// File:		getter.go
// Created by:	Hoven
// Created on:	2024-11-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package user

import "gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"

var _ shared.BaseUser = (*User)(nil)

func (u *User) GetUserId() int {
	if u == nil {
		return 0
	}
	return u.UserId
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetEmail() string {
	if u.UserInfo == nil {
		return ""
	}
	return u.UserInfo.Email
}

func (u *User) GetPhone() string {
	if u.UserInfo == nil {
		return ""
	}
	return u.UserInfo.Phone
}

func (u *User) GetAvatar() string {
	if u.UserInfo == nil {
		return ""
	}
	return u.UserInfo.Avatar
}

func (u *User) GetGender() string {
	return u.Gender.String()
}

func (u *User) GetStatus() int {
	return int(u.Status)
}

func (u *User) GetRole() int {
	return int(u.Role)
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
