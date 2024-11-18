// File:		wechat.go
// Created by:	Hoven
// Created on:	2024-11-17
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package wechat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-puzzles/puzzles/pflags"
	"github.com/pkg/errors"
)

type WechatSessionKeyResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

var (
	wxAppId       = pflags.StringRequired("wxAppId", "Wx app id.")
	wxAppSecretId = pflags.StringRequired("wxAppSecret", "Wx app secret id.")

	wechatApi = "https://api.weixin.qq.com/sns/jscode2session"
)

func GetSessionKey(ctx context.Context, code string) (*WechatSessionKeyResponse, error) {
	u, err := url.Parse(wechatApi)
	if err != nil {
		return nil, errors.Wrap(err, "parse wechat api url failed")
	}

	q := u.Query()
	q.Add("appid", wxAppId())
	q.Add("secret", wxAppSecretId())
	q.Add("grant_type", "authorization_code")
	q.Add("js_code", code)

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrap(err, "get wechat api response failed")
	}

	sessionResp := &WechatSessionKeyResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(sessionResp); err != nil {
		return nil, errors.Wrap(err, "decode session")
	}

	if sessionResp.ErrCode != 0 {
		return nil, errors.New(sessionResp.ErrMsg)
	}
	return sessionResp, nil
}
