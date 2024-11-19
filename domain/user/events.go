package user

import (
	"bytes"
	"html/template"
	"time"

	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/events"
)

const (
	emailCodeSendTemplate = `
	尊敬的Billiard用户，您好！
	您的验证码是：{{.code }}。五分钟内有效，请尽快验证。请勿泄露您的验证码。
	`
)

func GenerateSendCodeEventMessage(code string) ([]byte, error) {
	tmpl, err := template.New("emailCode").Parse(emailCodeSendTemplate)
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"code": code,
	}
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return nil, errors.Wrap(err, "tmplExecute")
	}

	return buf.Bytes(), nil
}

type SendCodeEvent struct {
	Code     string
	Target   string // phone number or email address
	Message  string
	ExpireAt time.Time
}

func NewSendPhoneCodeEvent(phone, code string, expireAt time.Time) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.SendPhoneCode,
		Payload: &SendCodeEvent{
			Code:     code,
			Target:   phone,
			ExpireAt: expireAt,
		},
	}
}

func NewSendEmailCodeEvent(email, code string, expireAt time.Time) *events.EventMessage {
	return &events.EventMessage{
		EventType: events.SendEmailCode,
		Payload: &SendCodeEvent{
			Code:     code,
			Target:   email,
			ExpireAt: expireAt,
		},
	}
}
