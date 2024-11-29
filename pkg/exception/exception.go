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
	ErrUserVerify            = NewBilliardException(400, "用户校验失败")
	ErrUserNotFound          = NewBilliardException(404, "用户不存在")
	ErrPasswordNotCorrect    = NewBilliardException(401, "密码不正确")
	ErrUserAlreadyExists     = NewBilliardException(400, "用户已存在")
	ErrRegisterUser          = NewBilliardException(400, "注册用户失败")
	ErrGetUserInfo           = NewBilliardException(400, "获取用户信息失败")
	ErrUpdateNameSameAsOld   = NewBilliardException(400, "更新名称相同")
	ErrUserNameAlreadyExists = NewBilliardException(400, "更新名称已存在")
	ErrUpdateUserName        = NewBilliardException(400, "更新用户昵称失败")
	ErrUpdateUserGender      = NewBilliardException(400, "更新性别失败")
	ErrUploadAvatar          = NewBilliardException(400, "上传头像失败")
	ErrGetAvatar             = NewBilliardException(400, "获取头像失败")
	ErrLoginFailed           = NewBilliardException(400, "登录失败")
	ErrLogoutFailed          = NewBilliardException(400, "退出登录失败")
	ErrUserAuthNotFound      = NewBilliardException(404, "用户认证信息不存在")

	ErrCreateGame     = NewBilliardException(400, "创建游戏失败")
	ErrDeleteGame     = NewBilliardException(400, "删除游戏失败")
	ErrUpdateGame     = NewBilliardException(400, "更新游戏信息失败")
	ErrGameNotFound   = NewBilliardException(404, "游戏不存在")
	ErrUploadGameIcon = NewBilliardException(400, "上传游戏图标失败")

	ErrCreateGameRoom          = NewBilliardException(400, "创建房间失败")
	ErrRoomUserMaxCreateNumber = NewBilliardException(400, "创建房间数到达上限")
	ErrUpdateGameRoom          = NewBilliardException(400, "更新房间状态失败")
	ErrGameRoomNotFound        = NewBilliardException(404, "游戏房间不存在")
	ErrGameRoomEnd             = NewBilliardException(404, "游戏已结束")
	ErrEnterGameRoom           = NewBilliardException(400, "进入游戏房间失败")
	ErrAlreadyInRoom           = NewBilliardException(400, "已经在房间中")
	ErrNotInRoom               = NewBilliardException(400, "不在该房间")
	ErrGameRoomFull            = NewBilliardException(400, "房间已满人")
	ErrLeaveGameRoom           = NewBilliardException(400, "离开游戏房间失败")
	ErrRoomOwnerNotMatch       = NewBilliardException(403, "不是房间拥有者")

	ErrGetGameList     = NewBilliardException(400, "获取游戏列表失败")
	ErrGetGameRoomList = NewBilliardException(400, "获取房间列表失败")
	ErrGetGameRoom     = NewBilliardException(400, "获取游戏房间失败")

	ErrGetNoticeList   = NewBilliardException(400, "获取通知列表失败")
	ErrGetSystemNotice = NewBilliardException(400, "获取系统通知失败")
	ErrAddNotice       = NewBilliardException(400, "添加通知失败")

	ErrUserNotInRoom = NewBilliardException(400, "用户不在房间内")
	ErrNotRoomOwner  = NewBilliardException(400, "不是房主")
	ErrPrepareGame   = NewBilliardException(400, "准备游戏失败")
	ErrStartGame     = NewBilliardException(400, "开始游戏失败")

	ErrBindPhone = NewBilliardException(400, "绑定手机号失败")
	ErrBindEmail = NewBilliardException(400, "绑定邮箱失败")

	ErrSendPhoneCode = NewBilliardException(400, "发送手机验证码失败")
	ErrSendEmailCode = NewBilliardException(400, "发送邮箱验证码失败")
	ErrVerifyCode    = NewBilliardException(400, "验证码错误或已过期")

	ErrHandleRecord       = NewBilliardException(400, "处理更新分数统计失败")
	ErrHandleAction       = NewBilliardException(400, "处理更新历史记录失败")
	ErrGetRoomRecord      = NewBilliardException(400, "获取房间分数统计失败")
	ErrGetRoomAction      = NewBilliardException(400, "获取房间历史操作记录失败")
	ErrRoomRecordNotFound = NewBilliardException(400, "房间记录不存在")
)

type BilliardException struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func NewBilliardException(code int, err string) *BilliardException {
	return &BilliardException{
		Code: code,
		Err:  err,
	}
}

func (se *BilliardException) Error() string {
	return fmt.Sprintf("%s.", se.Err)
}

func CheckException(err error) bool {
	se := new(BilliardException)
	return errors.As(err, &se)
}
