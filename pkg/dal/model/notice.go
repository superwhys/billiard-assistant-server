package model

import (
	"time"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/notice"
	"gorm.io/gorm"
)

type NoticePo struct {
	ID int `gorm:"primarykey"`

	Message string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (n *NoticePo) TableName() string {
	return "notices"
}

func (n *NoticePo) ToEntity() *notice.Notice {
	if n == nil {
		return nil
	}

	return &notice.Notice{
		Message: n.Message,
	}
}
