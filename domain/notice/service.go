package notice

import "context"

type INoticeService interface {
	GetNoticeList(ctx context.Context) ([]*Notice, error)
	GetNoticeByType(ctx context.Context, bt NoticeType) ([]*Notice, error)
	AddNotices(ctx context.Context, notices []*Notice) error
}
