package notice

import "context"

type INoticeRepo interface {
	GetNotices(ctx context.Context) ([]*Notice, error)
}
