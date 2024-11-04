// File:		exception.go
// Created by:	Hoven
// Created on:	2024-10-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package exception

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound       = NewSaException(404, "用户不存在")
	ErrPasswordNotCorrect = NewSaException(401, "密码不正确")
	ErrUserAlreadyExists  = NewSaException(400, "用户已存在")
	ErrRegisterUser       = NewSaException(400, "注册用户失败")
	ErrGetUserInfo        = NewSaException(400, "获取用户信息失败")
	ErrUpdateUserInfo     = NewSaException(400, "更新用户信息失败")
	ErrUploadAvatar       = NewSaException(400, "上传头像失败")
	ErrLoginFailed        = NewSaException(400, "登录失败")

	ErrCreateGame   = NewSaException(400, "创建游戏失败")
	ErrDeleteGame   = NewSaException(400, "删除游戏失败")
	ErrGameNotFound = NewSaException(404, "游戏不存在")

	ErrCreateGameRoom    = NewSaException(400, "创建房间失败")
	ErrUpdateGameRoom    = NewSaException(400, "更新房间状态失败")
	ErrGameRoomNotFound  = NewSaException(404, "游戏房间不存在")
	ErrEnterGameRoom     = NewSaException(400, "进入游戏房间失败")
	ErrGameRoomFull      = NewSaException(400, "房间已满人")
	ErrLeaveGameRoom     = NewSaException(400, "离开游戏房间失败")
	ErrRoomOwnerNotMatch = NewSaException(403, "不是访问拥有者")

	ErrGetGameList     = NewSaException(400, "获取游戏列表失败")
	ErrGetGameRoomList = NewSaException(400, "获取房间列表失败")
	ErrGetGameRoom     = NewSaException(400, "获取游戏房间失败")

	ErrGetNoticeList = NewSaException(400, "获取通知列表失败")
)

type SaException struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func NewSaException(code int, err string) *SaException {
	return &SaException{
		Code: code,
		Err:  err,
	}
}

func (se *SaException) Error() string {
	return fmt.Sprintf("错误: %s.", se.Err)
}

func CheckSaException(err error) bool {
	se := new(SaException)
	return errors.As(err, &se)
}
