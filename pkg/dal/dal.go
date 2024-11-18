package dal

import (
	"github.com/go-puzzles/puzzles/pgorm"
	"github.com/superwhys/billiard-assistant-server/pkg/dal/model"
)

func AllTables() []pgorm.SqlModel {
	return []pgorm.SqlModel{
		&model.UserPo{},
		&model.UserAuthPo{},
		&model.RoomPo{},
		&model.GamePo{},
		&model.NoticePo{},
		&model.RoomUserPo{},
	}
}
