// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package base

import (
	"context"

	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/dal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func newUserPo(db *gorm.DB, opts ...gen.DOOption) userPo {
	_userPo := userPo{}

	_userPo.userPoDo.UseDB(db, opts...)
	_userPo.userPoDo.UseModel(&model.UserPo{})

	tableName := _userPo.userPoDo.TableName()
	_userPo.ALL = field.NewAsterisk(tableName)
	_userPo.ID = field.NewInt(tableName, "id")
	_userPo.Name = field.NewString(tableName, "name")
	_userPo.Email = field.NewString(tableName, "email")
	_userPo.Phone = field.NewString(tableName, "phone")
	_userPo.Avatar = field.NewString(tableName, "avatar")
	_userPo.Gender = field.NewInt(tableName, "gender")
	_userPo.Status = field.NewInt(tableName, "status")
	_userPo.Role = field.NewInt(tableName, "role")
	_userPo.CreatedAt = field.NewTime(tableName, "created_at")
	_userPo.UpdatedAt = field.NewTime(tableName, "updated_at")
	_userPo.DeletedAt = field.NewField(tableName, "deleted_at")
	_userPo.RoomUsers = userPoHasManyRoomUsers{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("RoomUsers", "model.RoomUserPo"),
		Room: struct {
			field.RelationField
			Game struct {
				field.RelationField
			}
			Owner struct {
				field.RelationField
				RoomUsers struct {
					field.RelationField
				}
				UserAuthPos struct {
					field.RelationField
				}
			}
			RoomUsers struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("RoomUsers.Room", "model.RoomPo"),
			Game: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("RoomUsers.Room.Game", "model.GamePo"),
			},
			Owner: struct {
				field.RelationField
				RoomUsers struct {
					field.RelationField
				}
				UserAuthPos struct {
					field.RelationField
				}
			}{
				RelationField: field.NewRelation("RoomUsers.Room.Owner", "model.UserPo"),
				RoomUsers: struct {
					field.RelationField
				}{
					RelationField: field.NewRelation("RoomUsers.Room.Owner.RoomUsers", "model.RoomUserPo"),
				},
				UserAuthPos: struct {
					field.RelationField
				}{
					RelationField: field.NewRelation("RoomUsers.Room.Owner.UserAuthPos", "model.UserAuthPo"),
				},
			},
			RoomUsers: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("RoomUsers.Room.RoomUsers", "model.RoomUserPo"),
			},
		},
		User: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("RoomUsers.User", "model.UserPo"),
		},
	}

	_userPo.UserAuthPos = userPoHasManyUserAuthPos{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("UserAuthPos", "model.UserAuthPo"),
	}

	_userPo.fillFieldMap()

	return _userPo
}

type userPo struct {
	userPoDo userPoDo

	ALL       field.Asterisk
	ID        field.Int
	Name      field.String
	Email     field.String
	Phone     field.String
	Avatar    field.String
	Gender    field.Int
	Status    field.Int
	Role      field.Int
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	RoomUsers userPoHasManyRoomUsers

	UserAuthPos userPoHasManyUserAuthPos

	fieldMap map[string]field.Expr
}

func (u userPo) Table(newTableName string) *userPo {
	u.userPoDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userPo) As(alias string) *userPo {
	u.userPoDo.DO = *(u.userPoDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userPo) updateTableName(table string) *userPo {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt(table, "id")
	u.Name = field.NewString(table, "name")
	u.Email = field.NewString(table, "email")
	u.Phone = field.NewString(table, "phone")
	u.Avatar = field.NewString(table, "avatar")
	u.Gender = field.NewInt(table, "gender")
	u.Status = field.NewInt(table, "status")
	u.Role = field.NewInt(table, "role")
	u.CreatedAt = field.NewTime(table, "created_at")
	u.UpdatedAt = field.NewTime(table, "updated_at")
	u.DeletedAt = field.NewField(table, "deleted_at")

	u.fillFieldMap()

	return u
}

func (u *userPo) WithContext(ctx context.Context) IUserPoDo { return u.userPoDo.WithContext(ctx) }

func (u userPo) TableName() string { return u.userPoDo.TableName() }

func (u userPo) Alias() string { return u.userPoDo.Alias() }

func (u userPo) Columns(cols ...field.Expr) gen.Columns { return u.userPoDo.Columns(cols...) }

func (u *userPo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userPo) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 13)
	u.fieldMap["id"] = u.ID
	u.fieldMap["name"] = u.Name
	u.fieldMap["email"] = u.Email
	u.fieldMap["phone"] = u.Phone
	u.fieldMap["avatar"] = u.Avatar
	u.fieldMap["gender"] = u.Gender
	u.fieldMap["status"] = u.Status
	u.fieldMap["role"] = u.Role
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["updated_at"] = u.UpdatedAt
	u.fieldMap["deleted_at"] = u.DeletedAt

}

func (u userPo) clone(db *gorm.DB) userPo {
	u.userPoDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userPo) replaceDB(db *gorm.DB) userPo {
	u.userPoDo.ReplaceDB(db)
	return u
}

type userPoHasManyRoomUsers struct {
	db *gorm.DB

	field.RelationField

	Room struct {
		field.RelationField
		Game struct {
			field.RelationField
		}
		Owner struct {
			field.RelationField
			RoomUsers struct {
				field.RelationField
			}
			UserAuthPos struct {
				field.RelationField
			}
		}
		RoomUsers struct {
			field.RelationField
		}
	}
	User struct {
		field.RelationField
	}
}

func (a userPoHasManyRoomUsers) Where(conds ...field.Expr) *userPoHasManyRoomUsers {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a userPoHasManyRoomUsers) WithContext(ctx context.Context) *userPoHasManyRoomUsers {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a userPoHasManyRoomUsers) Session(session *gorm.Session) *userPoHasManyRoomUsers {
	a.db = a.db.Session(session)
	return &a
}

func (a userPoHasManyRoomUsers) Model(m *model.UserPo) *userPoHasManyRoomUsersTx {
	return &userPoHasManyRoomUsersTx{a.db.Model(m).Association(a.Name())}
}

type userPoHasManyRoomUsersTx struct{ tx *gorm.Association }

func (a userPoHasManyRoomUsersTx) Find() (result []*model.RoomUserPo, err error) {
	return result, a.tx.Find(&result)
}

func (a userPoHasManyRoomUsersTx) Append(values ...*model.RoomUserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a userPoHasManyRoomUsersTx) Replace(values ...*model.RoomUserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a userPoHasManyRoomUsersTx) Delete(values ...*model.RoomUserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a userPoHasManyRoomUsersTx) Clear() error {
	return a.tx.Clear()
}

func (a userPoHasManyRoomUsersTx) Count() int64 {
	return a.tx.Count()
}

type userPoHasManyUserAuthPos struct {
	db *gorm.DB

	field.RelationField
}

func (a userPoHasManyUserAuthPos) Where(conds ...field.Expr) *userPoHasManyUserAuthPos {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a userPoHasManyUserAuthPos) WithContext(ctx context.Context) *userPoHasManyUserAuthPos {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a userPoHasManyUserAuthPos) Session(session *gorm.Session) *userPoHasManyUserAuthPos {
	a.db = a.db.Session(session)
	return &a
}

func (a userPoHasManyUserAuthPos) Model(m *model.UserPo) *userPoHasManyUserAuthPosTx {
	return &userPoHasManyUserAuthPosTx{a.db.Model(m).Association(a.Name())}
}

type userPoHasManyUserAuthPosTx struct{ tx *gorm.Association }

func (a userPoHasManyUserAuthPosTx) Find() (result []*model.UserAuthPo, err error) {
	return result, a.tx.Find(&result)
}

func (a userPoHasManyUserAuthPosTx) Append(values ...*model.UserAuthPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a userPoHasManyUserAuthPosTx) Replace(values ...*model.UserAuthPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a userPoHasManyUserAuthPosTx) Delete(values ...*model.UserAuthPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a userPoHasManyUserAuthPosTx) Clear() error {
	return a.tx.Clear()
}

func (a userPoHasManyUserAuthPosTx) Count() int64 {
	return a.tx.Count()
}

type userPoDo struct{ gen.DO }

type IUserPoDo interface {
	gen.SubQuery
	Debug() IUserPoDo
	WithContext(ctx context.Context) IUserPoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserPoDo
	WriteDB() IUserPoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserPoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserPoDo
	Not(conds ...gen.Condition) IUserPoDo
	Or(conds ...gen.Condition) IUserPoDo
	Select(conds ...field.Expr) IUserPoDo
	Where(conds ...gen.Condition) IUserPoDo
	Order(conds ...field.Expr) IUserPoDo
	Distinct(cols ...field.Expr) IUserPoDo
	Omit(cols ...field.Expr) IUserPoDo
	Join(table schema.Tabler, on ...field.Expr) IUserPoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserPoDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserPoDo
	Group(cols ...field.Expr) IUserPoDo
	Having(conds ...gen.Condition) IUserPoDo
	Limit(limit int) IUserPoDo
	Offset(offset int) IUserPoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserPoDo
	Unscoped() IUserPoDo
	Create(values ...*model.UserPo) error
	CreateInBatches(values []*model.UserPo, batchSize int) error
	Save(values ...*model.UserPo) error
	First() (*model.UserPo, error)
	Take() (*model.UserPo, error)
	Last() (*model.UserPo, error)
	Find() ([]*model.UserPo, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserPo, err error)
	FindInBatches(result *[]*model.UserPo, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.UserPo) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserPoDo
	Assign(attrs ...field.AssignExpr) IUserPoDo
	Joins(fields ...field.RelationField) IUserPoDo
	Preload(fields ...field.RelationField) IUserPoDo
	FirstOrInit() (*model.UserPo, error)
	FirstOrCreate() (*model.UserPo, error)
	FindByPage(offset int, limit int) (result []*model.UserPo, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserPoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (u userPoDo) Debug() IUserPoDo {
	return u.withDO(u.DO.Debug())
}

func (u userPoDo) WithContext(ctx context.Context) IUserPoDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userPoDo) ReadDB() IUserPoDo {
	return u.Clauses(dbresolver.Read)
}

func (u userPoDo) WriteDB() IUserPoDo {
	return u.Clauses(dbresolver.Write)
}

func (u userPoDo) Session(config *gorm.Session) IUserPoDo {
	return u.withDO(u.DO.Session(config))
}

func (u userPoDo) Clauses(conds ...clause.Expression) IUserPoDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userPoDo) Returning(value interface{}, columns ...string) IUserPoDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userPoDo) Not(conds ...gen.Condition) IUserPoDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userPoDo) Or(conds ...gen.Condition) IUserPoDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userPoDo) Select(conds ...field.Expr) IUserPoDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userPoDo) Where(conds ...gen.Condition) IUserPoDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userPoDo) Order(conds ...field.Expr) IUserPoDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userPoDo) Distinct(cols ...field.Expr) IUserPoDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userPoDo) Omit(cols ...field.Expr) IUserPoDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userPoDo) Join(table schema.Tabler, on ...field.Expr) IUserPoDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userPoDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserPoDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userPoDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserPoDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userPoDo) Group(cols ...field.Expr) IUserPoDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userPoDo) Having(conds ...gen.Condition) IUserPoDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userPoDo) Limit(limit int) IUserPoDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userPoDo) Offset(offset int) IUserPoDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userPoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserPoDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userPoDo) Unscoped() IUserPoDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userPoDo) Create(values ...*model.UserPo) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userPoDo) CreateInBatches(values []*model.UserPo, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userPoDo) Save(values ...*model.UserPo) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userPoDo) First() (*model.UserPo, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPo), nil
	}
}

func (u userPoDo) Take() (*model.UserPo, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPo), nil
	}
}

func (u userPoDo) Last() (*model.UserPo, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPo), nil
	}
}

func (u userPoDo) Find() ([]*model.UserPo, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserPo), err
}

func (u userPoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserPo, err error) {
	buf := make([]*model.UserPo, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userPoDo) FindInBatches(result *[]*model.UserPo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userPoDo) Attrs(attrs ...field.AssignExpr) IUserPoDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userPoDo) Assign(attrs ...field.AssignExpr) IUserPoDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userPoDo) Joins(fields ...field.RelationField) IUserPoDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userPoDo) Preload(fields ...field.RelationField) IUserPoDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userPoDo) FirstOrInit() (*model.UserPo, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPo), nil
	}
}

func (u userPoDo) FirstOrCreate() (*model.UserPo, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserPo), nil
	}
}

func (u userPoDo) FindByPage(offset int, limit int) (result []*model.UserPo, count int64, err error) {
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

func (u userPoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userPoDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userPoDo) Delete(models ...*model.UserPo) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userPoDo) withDO(do gen.Dao) *userPoDo {
	u.DO = *do.(*gen.DO)
	return u
}
