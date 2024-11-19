package model

import (
	"time"

	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gorm.io/gorm"
)

type GamePo struct {
	ID int `gorm:"primarykey"`

	MaxPlayers  int
	GameType    shared.BilliardGameType
	Description string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (g *GamePo) TableName() string {
	return "games"
}

func (g *GamePo) FromEntity(ge *game.Game) {
	g.ID = ge.GameId
	g.GameType = ge.GameType
	if ge.GameConfig != nil {
		g.MaxPlayers = ge.GameConfig.MaxPlayers
		g.Description = ge.GameConfig.Desc
	}
}

func (g *GamePo) ToEntity() *game.Game {
	if g == nil {
		return nil
	}

	return &game.Game{
		GameId:   g.ID,
		GameType: g.GameType,
		GameConfig: &game.Config{
			MaxPlayers: g.MaxPlayers,
			Desc:       g.Description,
		},
	}
}
