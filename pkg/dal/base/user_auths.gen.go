// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package base

import (
	"context"

	"github.com/superwhys/billiard-assistant-server/pkg/dal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func newUserAuthPo(db *gorm.DB, opts ...gen.DOOption) userAuthPo {
	_userAuthPo := userAuthPo{}

	_userAuthPo.userAuthPoDo.UseDB(db, opts...)
	_userAuthPo.userAuthPoDo.UseModel(&model.UserAuthPo{})

	tableName := _userAuthPo.userAuthPoDo.TableName()
	_userAuthPo.ALL = field.NewAsterisk(tableName)
	_userAuthPo.ID = field.NewInt(tableName, "id")
	_userAuthPo.UserPoID = field.NewInt(tableName, "user_po_id")
	_userAuthPo.AuthType = field.NewInt(tableName, "auth_type")
	_userAuthPo.Identifier = field.NewString(tableName, "identifier")
	_userAuthPo.Credential = field.NewString(tableName, "credential")
	_userAuthPo.CreatedAt = field.NewTime(tableName, "created_at")
	_userAuthPo.UpdatedAt = field.NewTime(tableName, "updated_at")
	_userAuthPo.DeletedAt = field.NewField(tableName, "deleted_at")

	_userAuthPo.fillFieldMap()

	return _userAuthPo
}

type userAuthPo struct {
	userAuthPoDo userAuthPoDo

	ALL        field.Asterisk
	ID         field.Int
	UserPoID   field.Int
	AuthType   field.Int
	Identifier field.String
	Credential field.String
	CreatedAt  field.Time
	UpdatedAt  field.Time
	DeletedAt  field.Field

	fieldMap map[string]field.Expr
}

func (u userAuthPo) Table(newTableName string) *userAuthPo {
	u.userAuthPoDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userAuthPo) As(alias string) *userAuthPo {
	u.userAuthPoDo.DO = *(u.userAuthPoDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userAuthPo) updateTableName(table string) *userAuthPo {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt(table, "id")
	u.UserPoID = field.NewInt(table, "user_po_id")
	u.AuthType = field.NewInt(table, "auth_type")
	u.Identifier = field.NewString(table, "identifier")
	u.Credential = field.NewString(table, "credential")
	u.CreatedAt = field.NewTime(table, "created_at")
	u.UpdatedAt = field.NewTime(table, "updated_at")
	u.DeletedAt = field.NewField(table, "deleted_at")

	u.fillFieldMap()

	return u
}

func (u *userAuthPo) WithContext(ctx context.Context) IUserAuthPoDo {
	return u.userAuthPoDo.WithContext(ctx)
}

func (u userAuthPo) TableName() string { return u.userAuthPoDo.TableName() }

func (u userAuthPo) Alias() string { return u.userAuthPoDo.Alias() }

func (u userAuthPo) Columns(cols ...field.Expr) gen.Columns { return u.userAuthPoDo.Columns(cols...) }

func (u *userAuthPo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userAuthPo) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 8)
	u.fieldMap["id"] = u.ID
	u.fieldMap["user_po_id"] = u.UserPoID
	u.fieldMap["auth_type"] = u.AuthType
	u.fieldMap["identifier"] = u.Identifier
	u.fieldMap["credential"] = u.Credential
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["updated_at"] = u.UpdatedAt
	u.fieldMap["deleted_at"] = u.DeletedAt
}

func (u userAuthPo) clone(db *gorm.DB) userAuthPo {
	u.userAuthPoDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userAuthPo) replaceDB(db *gorm.DB) userAuthPo {
	u.userAuthPoDo.ReplaceDB(db)
	return u
}

type userAuthPoDo struct{ gen.DO }

type IUserAuthPoDo interface {
	gen.SubQuery
	Debug() IUserAuthPoDo
	WithContext(ctx context.Context) IUserAuthPoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserAuthPoDo
	WriteDB() IUserAuthPoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserAuthPoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserAuthPoDo
	Not(conds ...gen.Condition) IUserAuthPoDo
	Or(conds ...gen.Condition) IUserAuthPoDo
	Select(conds ...field.Expr) IUserAuthPoDo
	Where(conds ...gen.Condition) IUserAuthPoDo
	Order(conds ...field.Expr) IUserAuthPoDo
	Distinct(cols ...field.Expr) IUserAuthPoDo
	Omit(cols ...field.Expr) IUserAuthPoDo
	Join(table schema.Tabler, on ...field.Expr) IUserAuthPoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserAuthPoDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserAuthPoDo
	Group(cols ...field.Expr) IUserAuthPoDo
	Having(conds ...gen.Condition) IUserAuthPoDo
	Limit(limit int) IUserAuthPoDo
	Offset(offset int) IUserAuthPoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserAuthPoDo
	Unscoped() IUserAuthPoDo
	Create(values ...*model.UserAuthPo) error
	CreateInBatches(values []*model.UserAuthPo, batchSize int) error
	Save(values ...*model.UserAuthPo) error
	First() (*model.UserAuthPo, error)
	Take() (*model.UserAuthPo, error)
	Last() (*model.UserAuthPo, error)
	Find() ([]*model.UserAuthPo, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserAuthPo, err error)
	FindInBatches(result *[]*model.UserAuthPo, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.UserAuthPo) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserAuthPoDo
	Assign(attrs ...field.AssignExpr) IUserAuthPoDo
	Joins(fields ...field.RelationField) IUserAuthPoDo
	Preload(fields ...field.RelationField) IUserAuthPoDo
	FirstOrInit() (*model.UserAuthPo, error)
	FirstOrCreate() (*model.UserAuthPo, error)
	FindByPage(offset int, limit int) (result []*model.UserAuthPo, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserAuthPoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (u userAuthPoDo) Debug() IUserAuthPoDo {
	return u.withDO(u.DO.Debug())
}

func (u userAuthPoDo) WithContext(ctx context.Context) IUserAuthPoDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userAuthPoDo) ReadDB() IUserAuthPoDo {
	return u.Clauses(dbresolver.Read)
}

func (u userAuthPoDo) WriteDB() IUserAuthPoDo {
	return u.Clauses(dbresolver.Write)
}

func (u userAuthPoDo) Session(config *gorm.Session) IUserAuthPoDo {
	return u.withDO(u.DO.Session(config))
}

func (u userAuthPoDo) Clauses(conds ...clause.Expression) IUserAuthPoDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userAuthPoDo) Returning(value interface{}, columns ...string) IUserAuthPoDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userAuthPoDo) Not(conds ...gen.Condition) IUserAuthPoDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userAuthPoDo) Or(conds ...gen.Condition) IUserAuthPoDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userAuthPoDo) Select(conds ...field.Expr) IUserAuthPoDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userAuthPoDo) Where(conds ...gen.Condition) IUserAuthPoDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userAuthPoDo) Order(conds ...field.Expr) IUserAuthPoDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userAuthPoDo) Distinct(cols ...field.Expr) IUserAuthPoDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userAuthPoDo) Omit(cols ...field.Expr) IUserAuthPoDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userAuthPoDo) Join(table schema.Tabler, on ...field.Expr) IUserAuthPoDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userAuthPoDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserAuthPoDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userAuthPoDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserAuthPoDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userAuthPoDo) Group(cols ...field.Expr) IUserAuthPoDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userAuthPoDo) Having(conds ...gen.Condition) IUserAuthPoDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userAuthPoDo) Limit(limit int) IUserAuthPoDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userAuthPoDo) Offset(offset int) IUserAuthPoDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userAuthPoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserAuthPoDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userAuthPoDo) Unscoped() IUserAuthPoDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userAuthPoDo) Create(values ...*model.UserAuthPo) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userAuthPoDo) CreateInBatches(values []*model.UserAuthPo, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userAuthPoDo) Save(values ...*model.UserAuthPo) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userAuthPoDo) First() (*model.UserAuthPo, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAuthPo), nil
	}
}

func (u userAuthPoDo) Take() (*model.UserAuthPo, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAuthPo), nil
	}
}

func (u userAuthPoDo) Last() (*model.UserAuthPo, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAuthPo), nil
	}
}

func (u userAuthPoDo) Find() ([]*model.UserAuthPo, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserAuthPo), err
}

func (u userAuthPoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserAuthPo, err error) {
	buf := make([]*model.UserAuthPo, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userAuthPoDo) FindInBatches(result *[]*model.UserAuthPo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userAuthPoDo) Attrs(attrs ...field.AssignExpr) IUserAuthPoDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userAuthPoDo) Assign(attrs ...field.AssignExpr) IUserAuthPoDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userAuthPoDo) Joins(fields ...field.RelationField) IUserAuthPoDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userAuthPoDo) Preload(fields ...field.RelationField) IUserAuthPoDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userAuthPoDo) FirstOrInit() (*model.UserAuthPo, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAuthPo), nil
	}
}

func (u userAuthPoDo) FirstOrCreate() (*model.UserAuthPo, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserAuthPo), nil
	}
}

func (u userAuthPoDo) FindByPage(offset int, limit int) (result []*model.UserAuthPo, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userAuthPoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userAuthPoDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userAuthPoDo) Delete(models ...*model.UserAuthPo) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userAuthPoDo) withDO(do gen.Dao) *userAuthPoDo {
	u.DO = *do.(*gen.DO)
	return u
}
