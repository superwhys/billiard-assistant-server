package noticeDal

import (
	"context"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/notice"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/base"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gorm.io/gorm"
)

var _ notice.INoticeRepo = (*NoticeRepoImpl)(nil)

type NoticeRepoImpl struct {
	db *base.Query
}

func NewNoticeRepo(db *gorm.DB) *NoticeRepoImpl {
	return &NoticeRepoImpl{base.Use(db)}
}

func (n *NoticeRepoImpl) GetNotices(ctx context.Context) ([]*notice.Notice, error) {
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

func (n *NoticeRepoImpl) GetNoticesByType(ctx context.Context, bt notice.NoticeType) ([]*notice.Notice, error) {
	noticeDb := n.db.NoticePo

	noticeList, err := noticeDb.WithContext(ctx).Where(noticeDb.NoticeType.Eq(bt.String())).Find()
	if err != nil {
		return nil, err
	}

	var ret []*notice.Notice
	for _, notice := range noticeList {
		ret = append(ret, notice.ToEntity())
	}

	return ret, nil
}

func (n *NoticeRepoImpl) AddNotice(ctx context.Context, nc *notice.Notice) error {
	noticeDb := n.db.NoticePo

	ncPo := new(model.NoticePo)
	ncPo.FromEntity(nc)
	return noticeDb.WithContext(ctx).Create(ncPo)
}

func (n *NoticeRepoImpl) AddNoticeBatch(ctx context.Context, ns []*notice.Notice) error {
	noticePos := make([]*model.NoticePo, 0, len(ns))
	for _, nc := range ns {
		ncPo := new(model.NoticePo)
		ncPo.FromEntity(nc)
		noticePos = append(noticePos, ncPo)
	}

	return n.db.NoticePo.WithContext(ctx).CreateInBatches(noticePos, 100)
}
