package model

import (
	"time"
	
	"github.com/superwhys/billiard-assistant-server/domain/notice"
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
	return &notice.Notice{
		Message: n.Message,
	}
}
