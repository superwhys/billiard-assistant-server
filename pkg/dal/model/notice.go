package model

import (
	"time"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/notice"
	"gorm.io/gorm"
)

type NoticePo struct {
	ID int `gorm:"primarykey"`

	NoticeType notice.NoticeType
	Message    string

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
		NoticeType: n.NoticeType,
		Message:    n.Message,
	}
}

func (n *NoticePo) FromEntity(nt *notice.Notice) {
	if nt == nil {
		return
	}

	n.NoticeType = nt.NoticeType
	n.Message = nt.Message
}
