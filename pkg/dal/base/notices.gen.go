// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package base

import (
	"context"

	"gitlab.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func newNoticePo(db *gorm.DB, opts ...gen.DOOption) noticePo {
	_noticePo := noticePo{}

	_noticePo.noticePoDo.UseDB(db, opts...)
	_noticePo.noticePoDo.UseModel(&model.NoticePo{})

	tableName := _noticePo.noticePoDo.TableName()
	_noticePo.ALL = field.NewAsterisk(tableName)
	_noticePo.ID = field.NewInt(tableName, "id")
	_noticePo.NoticeType = field.NewString(tableName, "notice_type")
	_noticePo.Message = field.NewString(tableName, "message")
	_noticePo.CreatedAt = field.NewTime(tableName, "created_at")
	_noticePo.UpdatedAt = field.NewTime(tableName, "updated_at")
	_noticePo.DeletedAt = field.NewField(tableName, "deleted_at")

	_noticePo.fillFieldMap()

	return _noticePo
}

type noticePo struct {
	noticePoDo noticePoDo

	ALL        field.Asterisk
	ID         field.Int
	NoticeType field.String
	Message    field.String
	CreatedAt  field.Time
	UpdatedAt  field.Time
	DeletedAt  field.Field

	fieldMap map[string]field.Expr
}

func (n noticePo) Table(newTableName string) *noticePo {
	n.noticePoDo.UseTable(newTableName)
	return n.updateTableName(newTableName)
}

func (n noticePo) As(alias string) *noticePo {
	n.noticePoDo.DO = *(n.noticePoDo.As(alias).(*gen.DO))
	return n.updateTableName(alias)
}

func (n *noticePo) updateTableName(table string) *noticePo {
	n.ALL = field.NewAsterisk(table)
	n.ID = field.NewInt(table, "id")
	n.NoticeType = field.NewString(table, "notice_type")
	n.Message = field.NewString(table, "message")
	n.CreatedAt = field.NewTime(table, "created_at")
	n.UpdatedAt = field.NewTime(table, "updated_at")
	n.DeletedAt = field.NewField(table, "deleted_at")

	n.fillFieldMap()

	return n
}

func (n *noticePo) WithContext(ctx context.Context) INoticePoDo { return n.noticePoDo.WithContext(ctx) }

func (n noticePo) TableName() string { return n.noticePoDo.TableName() }

func (n noticePo) Alias() string { return n.noticePoDo.Alias() }

func (n noticePo) Columns(cols ...field.Expr) gen.Columns { return n.noticePoDo.Columns(cols...) }

func (n *noticePo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := n.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (n *noticePo) fillFieldMap() {
	n.fieldMap = make(map[string]field.Expr, 6)
	n.fieldMap["id"] = n.ID
	n.fieldMap["notice_type"] = n.NoticeType
	n.fieldMap["message"] = n.Message
	n.fieldMap["created_at"] = n.CreatedAt
	n.fieldMap["updated_at"] = n.UpdatedAt
	n.fieldMap["deleted_at"] = n.DeletedAt
}

func (n noticePo) clone(db *gorm.DB) noticePo {
	n.noticePoDo.ReplaceConnPool(db.Statement.ConnPool)
	return n
}

func (n noticePo) replaceDB(db *gorm.DB) noticePo {
	n.noticePoDo.ReplaceDB(db)
	return n
}

type noticePoDo struct{ gen.DO }

type INoticePoDo interface {
	gen.SubQuery
	Debug() INoticePoDo
	WithContext(ctx context.Context) INoticePoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() INoticePoDo
	WriteDB() INoticePoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) INoticePoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) INoticePoDo
	Not(conds ...gen.Condition) INoticePoDo
	Or(conds ...gen.Condition) INoticePoDo
	Select(conds ...field.Expr) INoticePoDo
	Where(conds ...gen.Condition) INoticePoDo
	Order(conds ...field.Expr) INoticePoDo
	Distinct(cols ...field.Expr) INoticePoDo
	Omit(cols ...field.Expr) INoticePoDo
	Join(table schema.Tabler, on ...field.Expr) INoticePoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) INoticePoDo
	RightJoin(table schema.Tabler, on ...field.Expr) INoticePoDo
	Group(cols ...field.Expr) INoticePoDo
	Having(conds ...gen.Condition) INoticePoDo
	Limit(limit int) INoticePoDo
	Offset(offset int) INoticePoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) INoticePoDo
	Unscoped() INoticePoDo
	Create(values ...*model.NoticePo) error
	CreateInBatches(values []*model.NoticePo, batchSize int) error
	Save(values ...*model.NoticePo) error
	First() (*model.NoticePo, error)
	Take() (*model.NoticePo, error)
	Last() (*model.NoticePo, error)
	Find() ([]*model.NoticePo, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.NoticePo, err error)
	FindInBatches(result *[]*model.NoticePo, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.NoticePo) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) INoticePoDo
	Assign(attrs ...field.AssignExpr) INoticePoDo
	Joins(fields ...field.RelationField) INoticePoDo
	Preload(fields ...field.RelationField) INoticePoDo
	FirstOrInit() (*model.NoticePo, error)
	FirstOrCreate() (*model.NoticePo, error)
	FindByPage(offset int, limit int) (result []*model.NoticePo, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) INoticePoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (n noticePoDo) Debug() INoticePoDo {
	return n.withDO(n.DO.Debug())
}

func (n noticePoDo) WithContext(ctx context.Context) INoticePoDo {
	return n.withDO(n.DO.WithContext(ctx))
}

func (n noticePoDo) ReadDB() INoticePoDo {
	return n.Clauses(dbresolver.Read)
}

func (n noticePoDo) WriteDB() INoticePoDo {
	return n.Clauses(dbresolver.Write)
}

func (n noticePoDo) Session(config *gorm.Session) INoticePoDo {
	return n.withDO(n.DO.Session(config))
}

func (n noticePoDo) Clauses(conds ...clause.Expression) INoticePoDo {
	return n.withDO(n.DO.Clauses(conds...))
}

func (n noticePoDo) Returning(value interface{}, columns ...string) INoticePoDo {
	return n.withDO(n.DO.Returning(value, columns...))
}

func (n noticePoDo) Not(conds ...gen.Condition) INoticePoDo {
	return n.withDO(n.DO.Not(conds...))
}

func (n noticePoDo) Or(conds ...gen.Condition) INoticePoDo {
	return n.withDO(n.DO.Or(conds...))
}

func (n noticePoDo) Select(conds ...field.Expr) INoticePoDo {
	return n.withDO(n.DO.Select(conds...))
}

func (n noticePoDo) Where(conds ...gen.Condition) INoticePoDo {
	return n.withDO(n.DO.Where(conds...))
}

func (n noticePoDo) Order(conds ...field.Expr) INoticePoDo {
	return n.withDO(n.DO.Order(conds...))
}

func (n noticePoDo) Distinct(cols ...field.Expr) INoticePoDo {
	return n.withDO(n.DO.Distinct(cols...))
}

func (n noticePoDo) Omit(cols ...field.Expr) INoticePoDo {
	return n.withDO(n.DO.Omit(cols...))
}

func (n noticePoDo) Join(table schema.Tabler, on ...field.Expr) INoticePoDo {
	return n.withDO(n.DO.Join(table, on...))
}

func (n noticePoDo) LeftJoin(table schema.Tabler, on ...field.Expr) INoticePoDo {
	return n.withDO(n.DO.LeftJoin(table, on...))
}

func (n noticePoDo) RightJoin(table schema.Tabler, on ...field.Expr) INoticePoDo {
	return n.withDO(n.DO.RightJoin(table, on...))
}

func (n noticePoDo) Group(cols ...field.Expr) INoticePoDo {
	return n.withDO(n.DO.Group(cols...))
}

func (n noticePoDo) Having(conds ...gen.Condition) INoticePoDo {
	return n.withDO(n.DO.Having(conds...))
}

func (n noticePoDo) Limit(limit int) INoticePoDo {
	return n.withDO(n.DO.Limit(limit))
}

func (n noticePoDo) Offset(offset int) INoticePoDo {
	return n.withDO(n.DO.Offset(offset))
}

func (n noticePoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) INoticePoDo {
	return n.withDO(n.DO.Scopes(funcs...))
}

func (n noticePoDo) Unscoped() INoticePoDo {
	return n.withDO(n.DO.Unscoped())
}

func (n noticePoDo) Create(values ...*model.NoticePo) error {
	if len(values) == 0 {
		return nil
	}
	return n.DO.Create(values)
}

func (n noticePoDo) CreateInBatches(values []*model.NoticePo, batchSize int) error {
	return n.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (n noticePoDo) Save(values ...*model.NoticePo) error {
	if len(values) == 0 {
		return nil
	}
	return n.DO.Save(values)
}

func (n noticePoDo) First() (*model.NoticePo, error) {
	if result, err := n.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.NoticePo), nil
	}
}

func (n noticePoDo) Take() (*model.NoticePo, error) {
	if result, err := n.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.NoticePo), nil
	}
}

func (n noticePoDo) Last() (*model.NoticePo, error) {
	if result, err := n.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.NoticePo), nil
	}
}

func (n noticePoDo) Find() ([]*model.NoticePo, error) {
	result, err := n.DO.Find()
	return result.([]*model.NoticePo), err
}

func (n noticePoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.NoticePo, err error) {
	buf := make([]*model.NoticePo, 0, batchSize)
	err = n.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (n noticePoDo) FindInBatches(result *[]*model.NoticePo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return n.DO.FindInBatches(result, batchSize, fc)
}

func (n noticePoDo) Attrs(attrs ...field.AssignExpr) INoticePoDo {
	return n.withDO(n.DO.Attrs(attrs...))
}

func (n noticePoDo) Assign(attrs ...field.AssignExpr) INoticePoDo {
	return n.withDO(n.DO.Assign(attrs...))
}

func (n noticePoDo) Joins(fields ...field.RelationField) INoticePoDo {
	for _, _f := range fields {
		n = *n.withDO(n.DO.Joins(_f))
	}
	return &n
}

func (n noticePoDo) Preload(fields ...field.RelationField) INoticePoDo {
	for _, _f := range fields {
		n = *n.withDO(n.DO.Preload(_f))
	}
	return &n
}

func (n noticePoDo) FirstOrInit() (*model.NoticePo, error) {
	if result, err := n.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.NoticePo), nil
	}
}

func (n noticePoDo) FirstOrCreate() (*model.NoticePo, error) {
	if result, err := n.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.NoticePo), nil
	}
}

func (n noticePoDo) FindByPage(offset int, limit int) (result []*model.NoticePo, count int64, err error) {
	result, err = n.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = n.Offset(-1).Limit(-1).Count()
	return
}

func (n noticePoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = n.Count()
	if err != nil {
		return
	}

	err = n.Offset(offset).Limit(limit).Scan(result)
	return
}

func (n noticePoDo) Scan(result interface{}) (err error) {
	return n.DO.Scan(result)
}

func (n noticePoDo) Delete(models ...*model.NoticePo) (result gen.ResultInfo, err error) {
	return n.DO.Delete(models)
}

func (n *noticePoDo) withDO(do gen.Dao) *noticePoDo {
	n.DO = *do.(*gen.DO)
	return n
}
