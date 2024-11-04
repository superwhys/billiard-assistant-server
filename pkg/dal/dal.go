package dal

import (
	"github.com/go-puzzles/pgorm"
	"github.com/superwhys/snooker-assistant-server/pkg/dal/model"
)

func AllTables() []pgorm.SqlModel {
	return []pgorm.SqlModel{
		&model.UserPo{},
		&model.RoomPo{},
		&model.GamePo{},
		&model.NoticePo{},
	}
}
