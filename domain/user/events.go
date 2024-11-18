package user

import (
	"time"

	"github.com/superwhys/snooker-assistant-server/pkg/events"
)

type SendCodeEvent struct {
	Code     string
	Target   string // phone number or email address
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
