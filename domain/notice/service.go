package notice

import "context"

type INoticeService interface {
	GetNoticeList(ctx context.Context) ([]*Notice, error)
}
