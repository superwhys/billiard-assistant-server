package notice

import "context"

type INoticeRepo interface {
	GetNotices(ctx context.Context) ([]*Notice, error)
	GetNoticesByType(ctx context.Context, bt NoticeType) ([]*Notice, error)
	AddNotice(ctx context.Context, n *Notice) error
	AddNoticeBatch(ctx context.Context, ns []*Notice) error
}
