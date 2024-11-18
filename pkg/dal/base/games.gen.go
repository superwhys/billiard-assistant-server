// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package base

import (
	"context"

	"github.com/superwhys/snooker-assistant-server/pkg/dal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func newGamePo(db *gorm.DB, opts ...gen.DOOption) gamePo {
	_gamePo := gamePo{}

	_gamePo.gamePoDo.UseDB(db, opts...)
	_gamePo.gamePoDo.UseModel(&model.GamePo{})

	tableName := _gamePo.gamePoDo.TableName()
	_gamePo.ALL = field.NewAsterisk(tableName)
	_gamePo.ID = field.NewInt(tableName, "id")
	_gamePo.MaxPlayers = field.NewInt(tableName, "max_players")
	_gamePo.GameType = field.NewInt(tableName, "game_type")
	_gamePo.Description = field.NewString(tableName, "description")
	_gamePo.CreatedAt = field.NewTime(tableName, "created_at")
	_gamePo.UpdatedAt = field.NewTime(tableName, "updated_at")
	_gamePo.DeletedAt = field.NewField(tableName, "deleted_at")

	_gamePo.fillFieldMap()

	return _gamePo
}

type gamePo struct {
	gamePoDo gamePoDo

	ALL         field.Asterisk
	ID          field.Int
	MaxPlayers  field.Int
	GameType    field.Int
	Description field.String
	CreatedAt   field.Time
	UpdatedAt   field.Time
	DeletedAt   field.Field

	fieldMap map[string]field.Expr
}

func (g gamePo) Table(newTableName string) *gamePo {
	g.gamePoDo.UseTable(newTableName)
	return g.updateTableName(newTableName)
}

func (g gamePo) As(alias string) *gamePo {
	g.gamePoDo.DO = *(g.gamePoDo.As(alias).(*gen.DO))
	return g.updateTableName(alias)
}

func (g *gamePo) updateTableName(table string) *gamePo {
	g.ALL = field.NewAsterisk(table)
	g.ID = field.NewInt(table, "id")
	g.MaxPlayers = field.NewInt(table, "max_players")
	g.GameType = field.NewInt(table, "game_type")
	g.Description = field.NewString(table, "description")
	g.CreatedAt = field.NewTime(table, "created_at")
	g.UpdatedAt = field.NewTime(table, "updated_at")
	g.DeletedAt = field.NewField(table, "deleted_at")

	g.fillFieldMap()

	return g
}

func (g *gamePo) WithContext(ctx context.Context) IGamePoDo { return g.gamePoDo.WithContext(ctx) }

func (g gamePo) TableName() string { return g.gamePoDo.TableName() }

func (g gamePo) Alias() string { return g.gamePoDo.Alias() }

func (g gamePo) Columns(cols ...field.Expr) gen.Columns { return g.gamePoDo.Columns(cols...) }

func (g *gamePo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := g.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (g *gamePo) fillFieldMap() {
	g.fieldMap = make(map[string]field.Expr, 7)
	g.fieldMap["id"] = g.ID
	g.fieldMap["max_players"] = g.MaxPlayers
	g.fieldMap["game_type"] = g.GameType
	g.fieldMap["description"] = g.Description
	g.fieldMap["created_at"] = g.CreatedAt
	g.fieldMap["updated_at"] = g.UpdatedAt
	g.fieldMap["deleted_at"] = g.DeletedAt
}

func (g gamePo) clone(db *gorm.DB) gamePo {
	g.gamePoDo.ReplaceConnPool(db.Statement.ConnPool)
	return g
}

func (g gamePo) replaceDB(db *gorm.DB) gamePo {
	g.gamePoDo.ReplaceDB(db)
	return g
}

type gamePoDo struct{ gen.DO }

type IGamePoDo interface {
	gen.SubQuery
	Debug() IGamePoDo
	WithContext(ctx context.Context) IGamePoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IGamePoDo
	WriteDB() IGamePoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IGamePoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IGamePoDo
	Not(conds ...gen.Condition) IGamePoDo
	Or(conds ...gen.Condition) IGamePoDo
	Select(conds ...field.Expr) IGamePoDo
	Where(conds ...gen.Condition) IGamePoDo
	Order(conds ...field.Expr) IGamePoDo
	Distinct(cols ...field.Expr) IGamePoDo
	Omit(cols ...field.Expr) IGamePoDo
	Join(table schema.Tabler, on ...field.Expr) IGamePoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IGamePoDo
	RightJoin(table schema.Tabler, on ...field.Expr) IGamePoDo
	Group(cols ...field.Expr) IGamePoDo
	Having(conds ...gen.Condition) IGamePoDo
	Limit(limit int) IGamePoDo
	Offset(offset int) IGamePoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IGamePoDo
	Unscoped() IGamePoDo
	Create(values ...*model.GamePo) error
	CreateInBatches(values []*model.GamePo, batchSize int) error
	Save(values ...*model.GamePo) error
	First() (*model.GamePo, error)
	Take() (*model.GamePo, error)
	Last() (*model.GamePo, error)
	Find() ([]*model.GamePo, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.GamePo, err error)
	FindInBatches(result *[]*model.GamePo, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.GamePo) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IGamePoDo
	Assign(attrs ...field.AssignExpr) IGamePoDo
	Joins(fields ...field.RelationField) IGamePoDo
	Preload(fields ...field.RelationField) IGamePoDo
	FirstOrInit() (*model.GamePo, error)
	FirstOrCreate() (*model.GamePo, error)
	FindByPage(offset int, limit int) (result []*model.GamePo, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IGamePoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (g gamePoDo) Debug() IGamePoDo {
	return g.withDO(g.DO.Debug())
}

func (g gamePoDo) WithContext(ctx context.Context) IGamePoDo {
	return g.withDO(g.DO.WithContext(ctx))
}

func (g gamePoDo) ReadDB() IGamePoDo {
	return g.Clauses(dbresolver.Read)
}

func (g gamePoDo) WriteDB() IGamePoDo {
	return g.Clauses(dbresolver.Write)
}

func (g gamePoDo) Session(config *gorm.Session) IGamePoDo {
	return g.withDO(g.DO.Session(config))
}

func (g gamePoDo) Clauses(conds ...clause.Expression) IGamePoDo {
	return g.withDO(g.DO.Clauses(conds...))
}

func (g gamePoDo) Returning(value interface{}, columns ...string) IGamePoDo {
	return g.withDO(g.DO.Returning(value, columns...))
}

func (g gamePoDo) Not(conds ...gen.Condition) IGamePoDo {
	return g.withDO(g.DO.Not(conds...))
}

func (g gamePoDo) Or(conds ...gen.Condition) IGamePoDo {
	return g.withDO(g.DO.Or(conds...))
}

func (g gamePoDo) Select(conds ...field.Expr) IGamePoDo {
	return g.withDO(g.DO.Select(conds...))
}

func (g gamePoDo) Where(conds ...gen.Condition) IGamePoDo {
	return g.withDO(g.DO.Where(conds...))
}

func (g gamePoDo) Order(conds ...field.Expr) IGamePoDo {
	return g.withDO(g.DO.Order(conds...))
}

func (g gamePoDo) Distinct(cols ...field.Expr) IGamePoDo {
	return g.withDO(g.DO.Distinct(cols...))
}

func (g gamePoDo) Omit(cols ...field.Expr) IGamePoDo {
	return g.withDO(g.DO.Omit(cols...))
}

func (g gamePoDo) Join(table schema.Tabler, on ...field.Expr) IGamePoDo {
	return g.withDO(g.DO.Join(table, on...))
}

func (g gamePoDo) LeftJoin(table schema.Tabler, on ...field.Expr) IGamePoDo {
	return g.withDO(g.DO.LeftJoin(table, on...))
}

func (g gamePoDo) RightJoin(table schema.Tabler, on ...field.Expr) IGamePoDo {
	return g.withDO(g.DO.RightJoin(table, on...))
}

func (g gamePoDo) Group(cols ...field.Expr) IGamePoDo {
	return g.withDO(g.DO.Group(cols...))
}

func (g gamePoDo) Having(conds ...gen.Condition) IGamePoDo {
	return g.withDO(g.DO.Having(conds...))
}

func (g gamePoDo) Limit(limit int) IGamePoDo {
	return g.withDO(g.DO.Limit(limit))
}

func (g gamePoDo) Offset(offset int) IGamePoDo {
	return g.withDO(g.DO.Offset(offset))
}

func (g gamePoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IGamePoDo {
	return g.withDO(g.DO.Scopes(funcs...))
}

func (g gamePoDo) Unscoped() IGamePoDo {
	return g.withDO(g.DO.Unscoped())
}

func (g gamePoDo) Create(values ...*model.GamePo) error {
	if len(values) == 0 {
		return nil
	}
	return g.DO.Create(values)
}

func (g gamePoDo) CreateInBatches(values []*model.GamePo, batchSize int) error {
	return g.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (g gamePoDo) Save(values ...*model.GamePo) error {
	if len(values) == 0 {
		return nil
	}
	return g.DO.Save(values)
}

func (g gamePoDo) First() (*model.GamePo, error) {
	if result, err := g.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.GamePo), nil
	}
}

func (g gamePoDo) Take() (*model.GamePo, error) {
	if result, err := g.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.GamePo), nil
	}
}

func (g gamePoDo) Last() (*model.GamePo, error) {
	if result, err := g.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.GamePo), nil
	}
}

func (g gamePoDo) Find() ([]*model.GamePo, error) {
	result, err := g.DO.Find()
	return result.([]*model.GamePo), err
}

func (g gamePoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.GamePo, err error) {
	buf := make([]*model.GamePo, 0, batchSize)
	err = g.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (g gamePoDo) FindInBatches(result *[]*model.GamePo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return g.DO.FindInBatches(result, batchSize, fc)
}

func (g gamePoDo) Attrs(attrs ...field.AssignExpr) IGamePoDo {
	return g.withDO(g.DO.Attrs(attrs...))
}

func (g gamePoDo) Assign(attrs ...field.AssignExpr) IGamePoDo {
	return g.withDO(g.DO.Assign(attrs...))
}

func (g gamePoDo) Joins(fields ...field.RelationField) IGamePoDo {
	for _, _f := range fields {
		g = *g.withDO(g.DO.Joins(_f))
	}
	return &g
}

func (g gamePoDo) Preload(fields ...field.RelationField) IGamePoDo {
	for _, _f := range fields {
		g = *g.withDO(g.DO.Preload(_f))
	}
	return &g
}

func (g gamePoDo) FirstOrInit() (*model.GamePo, error) {
	if result, err := g.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.GamePo), nil
	}
}

func (g gamePoDo) FirstOrCreate() (*model.GamePo, error) {
	if result, err := g.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.GamePo), nil
	}
}

func (g gamePoDo) FindByPage(offset int, limit int) (result []*model.GamePo, count int64, err error) {
	result, err = g.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = g.Offset(-1).Limit(-1).Count()
	return
}

func (g gamePoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = g.Count()
	if err != nil {
		return
	}

	err = g.Offset(offset).Limit(limit).Scan(result)
	return
}

func (g gamePoDo) Scan(result interface{}) (err error) {
	return g.DO.Scan(result)
}

func (g gamePoDo) Delete(models ...*model.GamePo) (result gen.ResultInfo, err error) {
	return g.DO.Delete(models)
}

func (g *gamePoDo) withDO(do gen.Dao) *gamePoDo {
	g.DO = *do.(*gen.DO)
	return g
}
