package gameDal

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.hoven.com/billiard/billiard-assistant-server/domain/game"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/base"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gorm.io/gorm"
)

var _ game.IGameRepo = (*GameRepoImpl)(nil)

type GameRepoImpl struct {
	db *base.Query
}

func NewGameRepo(db *gorm.DB) *GameRepoImpl {
	return &GameRepoImpl{
		db: base.Use(db),
	}
}

func (g *GameRepoImpl) CreateGame(ctx context.Context, ge *game.Game) error {
	gamePo := new(model.GamePo)
	gamePo.FromEntity(ge)

	gameDb := g.db.GamePo
	return gameDb.WithContext(ctx).Create(gamePo)
}

func (g *GameRepoImpl) DeleteGame(ctx context.Context, gameId int) error {
	gameDb := g.db.GamePo
	ret, err := gameDb.WithContext(ctx).Where(gameDb.ID.Eq(gameId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrGameNotFound
	} else if err != nil {
		return err
	}

	_, err = gameDb.WithContext(ctx).Delete(ret)
	return err
}

func (g *GameRepoImpl) UpdateGame(ctx context.Context, game *game.Game) error {
	gameDb := g.db.GamePo
	ret, err := gameDb.WithContext(ctx).Where(gameDb.ID.Eq(game.GameId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrGameNotFound
	} else if err != nil {
		return err
	}

	ret.MaxPlayers = game.GameConfig.MaxPlayers

	return gameDb.WithContext(ctx).Save(ret)
}

func (g *GameRepoImpl) GetGameById(ctx context.Context, gameId int) (*game.Game, error) {
	gameDb := g.db.GamePo
	ret, err := gameDb.WithContext(ctx).Where(gameDb.ID.Eq(gameId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrGameNotFound
	} else if err != nil {
		return nil, err
	}

	return ret.ToEntity(), nil
}

func (g *GameRepoImpl) GetGameList(ctx context.Context) ([]*game.Game, error) {
	gameDb := g.db.GamePo
	gameList, err := gameDb.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}

	var ret []*game.Game
	for _, gamePo := range gameList {
		ret = append(ret, gamePo.ToEntity())
	}

	return ret, nil
}
