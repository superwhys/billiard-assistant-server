package noticeDal

import (
	"context"
	
	"github.com/superwhys/billiard-assistant-server/domain/notice"
	"github.com/superwhys/billiard-assistant-server/pkg/dal/base"
	"gorm.io/gorm"
)

var _ notice.INoticeRepo = (*NoticeRepoImpl)(nil)

type NoticeRepoImpl struct {
	db *base.Query
}

func NewNoticeRepo(db *gorm.DB) *NoticeRepoImpl {
	return &NoticeRepoImpl{base.Use(db)}
}

func (n NoticeRepoImpl) GetNotices(ctx context.Context) ([]*notice.Notice, error) {
	noticeList, err := n.db.NoticePo.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	
	var ret []*notice.Notice
	for _, notice := range noticeList {
		ret = append(ret, notice.ToEntity())
	}
	
	return ret, nil
}
