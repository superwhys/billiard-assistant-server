package dal

import (
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"github.com/go-puzzles/puzzles/pgorm"
)

func AllTables() []pgorm.SqlModel {
	return []pgorm.SqlModel{
		&model.UserPo{},
		&model.RoomPo{},
		&model.GamePo{},
		&model.NoticePo{},
		&model.RoomUserPo{},
		&model.RecordPo{},
	}
}
