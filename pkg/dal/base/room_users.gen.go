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

func newRoomUserPo(db *gorm.DB, opts ...gen.DOOption) roomUserPo {
	_roomUserPo := roomUserPo{}

	_roomUserPo.roomUserPoDo.UseDB(db, opts...)
	_roomUserPo.roomUserPoDo.UseModel(&model.RoomUserPo{})

	tableName := _roomUserPo.roomUserPoDo.TableName()
	_roomUserPo.ALL = field.NewAsterisk(tableName)
	_roomUserPo.ID = field.NewInt(tableName, "id")
	_roomUserPo.RoomID = field.NewInt(tableName, "room_id")
	_roomUserPo.UserID = field.NewInt(tableName, "user_id")
	_roomUserPo.UserName = field.NewString(tableName, "user_name")
	_roomUserPo.IsVirtualPlayer = field.NewBool(tableName, "is_virtual_player")
	_roomUserPo.CreatedAt = field.NewTime(tableName, "created_at")
	_roomUserPo.UpdatedAt = field.NewTime(tableName, "updated_at")
	_roomUserPo.HeartbeatAt = field.NewTime(tableName, "heartbeat_at")
	_roomUserPo.Room = roomUserPoBelongsToRoom{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Room", "model.RoomPo"),
		Game: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Room.Game", "model.GamePo"),
		},
		Owner: struct {
			field.RelationField
			RoomUsers struct {
				field.RelationField
				Room struct {
					field.RelationField
				}
				User struct {
					field.RelationField
				}
			}
		}{
			RelationField: field.NewRelation("Room.Owner", "model.UserPo"),
			RoomUsers: struct {
				field.RelationField
				Room struct {
					field.RelationField
				}
				User struct {
					field.RelationField
				}
			}{
				RelationField: field.NewRelation("Room.Owner.RoomUsers", "model.RoomUserPo"),
				Room: struct {
					field.RelationField
				}{
					RelationField: field.NewRelation("Room.Owner.RoomUsers.Room", "model.RoomPo"),
				},
				User: struct {
					field.RelationField
				}{
					RelationField: field.NewRelation("Room.Owner.RoomUsers.User", "model.UserPo"),
				},
			},
		},
		RoomUsers: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Room.RoomUsers", "model.RoomUserPo"),
		},
	}

	_roomUserPo.User = roomUserPoBelongsToUser{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("User", "model.UserPo"),
	}

	_roomUserPo.fillFieldMap()

	return _roomUserPo
}

type roomUserPo struct {
	roomUserPoDo roomUserPoDo

	ALL             field.Asterisk
	ID              field.Int
	RoomID          field.Int
	UserID          field.Int
	UserName        field.String
	IsVirtualPlayer field.Bool
	CreatedAt       field.Time
	UpdatedAt       field.Time
	HeartbeatAt     field.Time
	Room            roomUserPoBelongsToRoom

	User roomUserPoBelongsToUser

	fieldMap map[string]field.Expr
}

func (r roomUserPo) Table(newTableName string) *roomUserPo {
	r.roomUserPoDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r roomUserPo) As(alias string) *roomUserPo {
	r.roomUserPoDo.DO = *(r.roomUserPoDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *roomUserPo) updateTableName(table string) *roomUserPo {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewInt(table, "id")
	r.RoomID = field.NewInt(table, "room_id")
	r.UserID = field.NewInt(table, "user_id")
	r.UserName = field.NewString(table, "user_name")
	r.IsVirtualPlayer = field.NewBool(table, "is_virtual_player")
	r.CreatedAt = field.NewTime(table, "created_at")
	r.UpdatedAt = field.NewTime(table, "updated_at")
	r.HeartbeatAt = field.NewTime(table, "heartbeat_at")

	r.fillFieldMap()

	return r
}

func (r *roomUserPo) WithContext(ctx context.Context) IRoomUserPoDo {
	return r.roomUserPoDo.WithContext(ctx)
}

func (r roomUserPo) TableName() string { return r.roomUserPoDo.TableName() }

func (r roomUserPo) Alias() string { return r.roomUserPoDo.Alias() }

func (r roomUserPo) Columns(cols ...field.Expr) gen.Columns { return r.roomUserPoDo.Columns(cols...) }

func (r *roomUserPo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *roomUserPo) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 10)
	r.fieldMap["id"] = r.ID
	r.fieldMap["room_id"] = r.RoomID
	r.fieldMap["user_id"] = r.UserID
	r.fieldMap["user_name"] = r.UserName
	r.fieldMap["is_virtual_player"] = r.IsVirtualPlayer
	r.fieldMap["created_at"] = r.CreatedAt
	r.fieldMap["updated_at"] = r.UpdatedAt
	r.fieldMap["heartbeat_at"] = r.HeartbeatAt

}

func (r roomUserPo) clone(db *gorm.DB) roomUserPo {
	r.roomUserPoDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r roomUserPo) replaceDB(db *gorm.DB) roomUserPo {
	r.roomUserPoDo.ReplaceDB(db)
	return r
}

type roomUserPoBelongsToRoom struct {
	db *gorm.DB

	field.RelationField

	Game struct {
		field.RelationField
	}
	Owner struct {
		field.RelationField
		RoomUsers struct {
			field.RelationField
			Room struct {
				field.RelationField
			}
			User struct {
				field.RelationField
			}
		}
	}
	RoomUsers struct {
		field.RelationField
	}
}

func (a roomUserPoBelongsToRoom) Where(conds ...field.Expr) *roomUserPoBelongsToRoom {
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

func (a roomUserPoBelongsToRoom) WithContext(ctx context.Context) *roomUserPoBelongsToRoom {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a roomUserPoBelongsToRoom) Session(session *gorm.Session) *roomUserPoBelongsToRoom {
	a.db = a.db.Session(session)
	return &a
}

func (a roomUserPoBelongsToRoom) Model(m *model.RoomUserPo) *roomUserPoBelongsToRoomTx {
	return &roomUserPoBelongsToRoomTx{a.db.Model(m).Association(a.Name())}
}

type roomUserPoBelongsToRoomTx struct{ tx *gorm.Association }

func (a roomUserPoBelongsToRoomTx) Find() (result *model.RoomPo, err error) {
	return result, a.tx.Find(&result)
}

func (a roomUserPoBelongsToRoomTx) Append(values ...*model.RoomPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a roomUserPoBelongsToRoomTx) Replace(values ...*model.RoomPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a roomUserPoBelongsToRoomTx) Delete(values ...*model.RoomPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a roomUserPoBelongsToRoomTx) Clear() error {
	return a.tx.Clear()
}

func (a roomUserPoBelongsToRoomTx) Count() int64 {
	return a.tx.Count()
}

type roomUserPoBelongsToUser struct {
	db *gorm.DB

	field.RelationField
}

func (a roomUserPoBelongsToUser) Where(conds ...field.Expr) *roomUserPoBelongsToUser {
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

func (a roomUserPoBelongsToUser) WithContext(ctx context.Context) *roomUserPoBelongsToUser {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a roomUserPoBelongsToUser) Session(session *gorm.Session) *roomUserPoBelongsToUser {
	a.db = a.db.Session(session)
	return &a
}

func (a roomUserPoBelongsToUser) Model(m *model.RoomUserPo) *roomUserPoBelongsToUserTx {
	return &roomUserPoBelongsToUserTx{a.db.Model(m).Association(a.Name())}
}

type roomUserPoBelongsToUserTx struct{ tx *gorm.Association }

func (a roomUserPoBelongsToUserTx) Find() (result *model.UserPo, err error) {
	return result, a.tx.Find(&result)
}

func (a roomUserPoBelongsToUserTx) Append(values ...*model.UserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a roomUserPoBelongsToUserTx) Replace(values ...*model.UserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a roomUserPoBelongsToUserTx) Delete(values ...*model.UserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a roomUserPoBelongsToUserTx) Clear() error {
	return a.tx.Clear()
}

func (a roomUserPoBelongsToUserTx) Count() int64 {
	return a.tx.Count()
}

type roomUserPoDo struct{ gen.DO }

type IRoomUserPoDo interface {
	gen.SubQuery
	Debug() IRoomUserPoDo
	WithContext(ctx context.Context) IRoomUserPoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IRoomUserPoDo
	WriteDB() IRoomUserPoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IRoomUserPoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRoomUserPoDo
	Not(conds ...gen.Condition) IRoomUserPoDo
	Or(conds ...gen.Condition) IRoomUserPoDo
	Select(conds ...field.Expr) IRoomUserPoDo
	Where(conds ...gen.Condition) IRoomUserPoDo
	Order(conds ...field.Expr) IRoomUserPoDo
	Distinct(cols ...field.Expr) IRoomUserPoDo
	Omit(cols ...field.Expr) IRoomUserPoDo
	Join(table schema.Tabler, on ...field.Expr) IRoomUserPoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRoomUserPoDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRoomUserPoDo
	Group(cols ...field.Expr) IRoomUserPoDo
	Having(conds ...gen.Condition) IRoomUserPoDo
	Limit(limit int) IRoomUserPoDo
	Offset(offset int) IRoomUserPoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRoomUserPoDo
	Unscoped() IRoomUserPoDo
	Create(values ...*model.RoomUserPo) error
	CreateInBatches(values []*model.RoomUserPo, batchSize int) error
	Save(values ...*model.RoomUserPo) error
	First() (*model.RoomUserPo, error)
	Take() (*model.RoomUserPo, error)
	Last() (*model.RoomUserPo, error)
	Find() ([]*model.RoomUserPo, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RoomUserPo, err error)
	FindInBatches(result *[]*model.RoomUserPo, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.RoomUserPo) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRoomUserPoDo
	Assign(attrs ...field.AssignExpr) IRoomUserPoDo
	Joins(fields ...field.RelationField) IRoomUserPoDo
	Preload(fields ...field.RelationField) IRoomUserPoDo
	FirstOrInit() (*model.RoomUserPo, error)
	FirstOrCreate() (*model.RoomUserPo, error)
	FindByPage(offset int, limit int) (result []*model.RoomUserPo, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRoomUserPoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (r roomUserPoDo) Debug() IRoomUserPoDo {
	return r.withDO(r.DO.Debug())
}

func (r roomUserPoDo) WithContext(ctx context.Context) IRoomUserPoDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r roomUserPoDo) ReadDB() IRoomUserPoDo {
	return r.Clauses(dbresolver.Read)
}

func (r roomUserPoDo) WriteDB() IRoomUserPoDo {
	return r.Clauses(dbresolver.Write)
}

func (r roomUserPoDo) Session(config *gorm.Session) IRoomUserPoDo {
	return r.withDO(r.DO.Session(config))
}

func (r roomUserPoDo) Clauses(conds ...clause.Expression) IRoomUserPoDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r roomUserPoDo) Returning(value interface{}, columns ...string) IRoomUserPoDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r roomUserPoDo) Not(conds ...gen.Condition) IRoomUserPoDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r roomUserPoDo) Or(conds ...gen.Condition) IRoomUserPoDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r roomUserPoDo) Select(conds ...field.Expr) IRoomUserPoDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r roomUserPoDo) Where(conds ...gen.Condition) IRoomUserPoDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r roomUserPoDo) Order(conds ...field.Expr) IRoomUserPoDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r roomUserPoDo) Distinct(cols ...field.Expr) IRoomUserPoDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r roomUserPoDo) Omit(cols ...field.Expr) IRoomUserPoDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r roomUserPoDo) Join(table schema.Tabler, on ...field.Expr) IRoomUserPoDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r roomUserPoDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRoomUserPoDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r roomUserPoDo) RightJoin(table schema.Tabler, on ...field.Expr) IRoomUserPoDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r roomUserPoDo) Group(cols ...field.Expr) IRoomUserPoDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r roomUserPoDo) Having(conds ...gen.Condition) IRoomUserPoDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r roomUserPoDo) Limit(limit int) IRoomUserPoDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r roomUserPoDo) Offset(offset int) IRoomUserPoDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r roomUserPoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRoomUserPoDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r roomUserPoDo) Unscoped() IRoomUserPoDo {
	return r.withDO(r.DO.Unscoped())
}

func (r roomUserPoDo) Create(values ...*model.RoomUserPo) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r roomUserPoDo) CreateInBatches(values []*model.RoomUserPo, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r roomUserPoDo) Save(values ...*model.RoomUserPo) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r roomUserPoDo) First() (*model.RoomUserPo, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomUserPo), nil
	}
}

func (r roomUserPoDo) Take() (*model.RoomUserPo, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomUserPo), nil
	}
}

func (r roomUserPoDo) Last() (*model.RoomUserPo, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomUserPo), nil
	}
}

func (r roomUserPoDo) Find() ([]*model.RoomUserPo, error) {
	result, err := r.DO.Find()
	return result.([]*model.RoomUserPo), err
}

func (r roomUserPoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RoomUserPo, err error) {
	buf := make([]*model.RoomUserPo, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r roomUserPoDo) FindInBatches(result *[]*model.RoomUserPo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r roomUserPoDo) Attrs(attrs ...field.AssignExpr) IRoomUserPoDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r roomUserPoDo) Assign(attrs ...field.AssignExpr) IRoomUserPoDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r roomUserPoDo) Joins(fields ...field.RelationField) IRoomUserPoDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r roomUserPoDo) Preload(fields ...field.RelationField) IRoomUserPoDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r roomUserPoDo) FirstOrInit() (*model.RoomUserPo, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomUserPo), nil
	}
}

func (r roomUserPoDo) FirstOrCreate() (*model.RoomUserPo, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomUserPo), nil
	}
}

func (r roomUserPoDo) FindByPage(offset int, limit int) (result []*model.RoomUserPo, count int64, err error) {
	result, err = r.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = r.Offset(-1).Limit(-1).Count()
	return
}

func (r roomUserPoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r roomUserPoDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r roomUserPoDo) Delete(models ...*model.RoomUserPo) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *roomUserPoDo) withDO(do gen.Dao) *roomUserPoDo {
	r.DO = *do.(*gen.DO)
	return r
}
