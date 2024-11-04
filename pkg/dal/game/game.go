package gameDal

import (
	"context"
	
	"github.com/pkg/errors"
	"github.com/superwhys/snooker-assistant-server/domain/game"
	"github.com/superwhys/snooker-assistant-server/pkg/dal/base"
	"github.com/superwhys/snooker-assistant-server/pkg/dal/model"
	"github.com/superwhys/snooker-assistant-server/pkg/exception"
	"gorm.io/gorm"
)

var _ game.IGameRepo = (*GameRepoImpl)(nil)

type GameRepoImpl struct {
	db *base.GameDB
}

func NewGameRepo(db *gorm.DB) *GameRepoImpl {
	return &GameRepoImpl{
		db: base.NewGameDB(db),
	}
}

func (g *GameRepoImpl) CreateGame(ctx context.Context, ge *game.Game) error {
	gamePo := new(model.GamePo)
	gamePo.FromEntity(ge)
	
	return g.db.WithContext(ctx).Create(gamePo)
}

func (g *GameRepoImpl) DeleteGame(ctx context.Context, gameId int) error {
	ret, err := g.db.WithContext(ctx).Where(g.db.ID.Eq(gameId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrGameNotFound
	} else if err != nil {
		return err
	}
	
	_, err = g.db.WithContext(ctx).Delete(ret)
	return err
}

func (g *GameRepoImpl) UpdateGame(ctx context.Context, game *game.Game) error {
	ret, err := g.db.WithContext(ctx).Where(g.db.ID.Eq(game.GameId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.ErrGameNotFound
	} else if err != nil {
		return err
	}
	
	ret.MaxPlayers = game.GameConfig.MaxPlayers
	
	return g.db.WithContext(ctx).Save(ret)
}

func (g *GameRepoImpl) GetGameById(ctx context.Context, gameId int) (*game.Game, error) {
	ret, err := g.db.WithContext(ctx).Where(g.db.ID.Eq(gameId)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.ErrGameNotFound
	} else if err != nil {
		return nil, err
	}
	
	return ret.ToEntity(), nil
}

func (g *GameRepoImpl) GetGameList(ctx context.Context) ([]*game.Game, error) {
	gameList, err := g.db.WithContext(ctx).Find()
	if err != nil {
		return nil, err
	}
	
	var ret []*game.Game
	for _, gamePo := range gameList {
		ret = append(ret, gamePo.ToEntity())
	}
	
	return ret, nil
}
