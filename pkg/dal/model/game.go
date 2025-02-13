package model

import (
	"time"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitea.hoven.com/billiard/billiard-assistant-server/domain/shared"
	"gorm.io/gorm"
)

type GamePo struct {
	ID int `gorm:"primarykey"`

	MaxPlayers  int
	GameType    shared.BilliardGameType
	Icon        string `gorm:"type:varchar(255)"`
	IsActivated bool
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
	g.Icon = ge.Icon
	g.IsActivated = ge.IsActivate
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
		GameId:     g.ID,
		GameType:   g.GameType,
		Icon:       g.Icon,
		IsActivate: g.IsActivated,
		GameConfig: &game.Config{
			MaxPlayers: g.MaxPlayers,
			Desc:       g.Description,
		},
	}
}
