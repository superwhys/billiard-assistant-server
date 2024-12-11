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

func newRoomPo(db *gorm.DB, opts ...gen.DOOption) roomPo {
	_roomPo := roomPo{}

	_roomPo.roomPoDo.UseDB(db, opts...)
	_roomPo.roomPoDo.UseModel(&model.RoomPo{})

	tableName := _roomPo.roomPoDo.TableName()
	_roomPo.ALL = field.NewAsterisk(tableName)
	_roomPo.ID = field.NewInt(tableName, "id")
	_roomPo.RoomCode = field.NewString(tableName, "room_code")
	_roomPo.GameID = field.NewInt(tableName, "game_id")
	_roomPo.OwnerID = field.NewInt(tableName, "owner_id")
	_roomPo.Extra = field.NewField(tableName, "extra")
	_roomPo.GameStatus = field.NewInt(tableName, "game_status")
	_roomPo.WinLoseStatus = field.NewInt(tableName, "win_lose_status")
	_roomPo.CreatedAt = field.NewTime(tableName, "created_at")
	_roomPo.UpdatedAt = field.NewTime(tableName, "updated_at")
	_roomPo.DeletedAt = field.NewField(tableName, "deleted_at")
	_roomPo.RoomUsers = roomPoHasManyRoomUsers{
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

	_roomPo.Game = roomPoBelongsToGame{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Game", "model.GamePo"),
	}

	_roomPo.Owner = roomPoBelongsToOwner{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Owner", "model.UserPo"),
	}

	_roomPo.fillFieldMap()

	return _roomPo
}

type roomPo struct {
	roomPoDo roomPoDo

	ALL           field.Asterisk
	ID            field.Int
	RoomCode      field.String
	GameID        field.Int
	OwnerID       field.Int
	Extra         field.Field
	GameStatus    field.Int
	WinLoseStatus field.Int
	CreatedAt     field.Time
	UpdatedAt     field.Time
	DeletedAt     field.Field
	RoomUsers     roomPoHasManyRoomUsers

	Game roomPoBelongsToGame

	Owner roomPoBelongsToOwner

	fieldMap map[string]field.Expr
}

func (r roomPo) Table(newTableName string) *roomPo {
	r.roomPoDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r roomPo) As(alias string) *roomPo {
	r.roomPoDo.DO = *(r.roomPoDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *roomPo) updateTableName(table string) *roomPo {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewInt(table, "id")
	r.RoomCode = field.NewString(table, "room_code")
	r.GameID = field.NewInt(table, "game_id")
	r.OwnerID = field.NewInt(table, "owner_id")
	r.Extra = field.NewField(table, "extra")
	r.GameStatus = field.NewInt(table, "game_status")
	r.WinLoseStatus = field.NewInt(table, "win_lose_status")
	r.CreatedAt = field.NewTime(table, "created_at")
	r.UpdatedAt = field.NewTime(table, "updated_at")
	r.DeletedAt = field.NewField(table, "deleted_at")

	r.fillFieldMap()

	return r
}

func (r *roomPo) WithContext(ctx context.Context) IRoomPoDo { return r.roomPoDo.WithContext(ctx) }

func (r roomPo) TableName() string { return r.roomPoDo.TableName() }

func (r roomPo) Alias() string { return r.roomPoDo.Alias() }

func (r roomPo) Columns(cols ...field.Expr) gen.Columns { return r.roomPoDo.Columns(cols...) }

func (r *roomPo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *roomPo) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 13)
	r.fieldMap["id"] = r.ID
	r.fieldMap["room_code"] = r.RoomCode
	r.fieldMap["game_id"] = r.GameID
	r.fieldMap["owner_id"] = r.OwnerID
	r.fieldMap["extra"] = r.Extra
	r.fieldMap["game_status"] = r.GameStatus
	r.fieldMap["win_lose_status"] = r.WinLoseStatus
	r.fieldMap["created_at"] = r.CreatedAt
	r.fieldMap["updated_at"] = r.UpdatedAt
	r.fieldMap["deleted_at"] = r.DeletedAt

}

func (r roomPo) clone(db *gorm.DB) roomPo {
	r.roomPoDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r roomPo) replaceDB(db *gorm.DB) roomPo {
	r.roomPoDo.ReplaceDB(db)
	return r
}

type roomPoHasManyRoomUsers struct {
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

func (a roomPoHasManyRoomUsers) Where(conds ...field.Expr) *roomPoHasManyRoomUsers {
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

func (a roomPoHasManyRoomUsers) WithContext(ctx context.Context) *roomPoHasManyRoomUsers {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a roomPoHasManyRoomUsers) Session(session *gorm.Session) *roomPoHasManyRoomUsers {
	a.db = a.db.Session(session)
	return &a
}

func (a roomPoHasManyRoomUsers) Model(m *model.RoomPo) *roomPoHasManyRoomUsersTx {
	return &roomPoHasManyRoomUsersTx{a.db.Model(m).Association(a.Name())}
}

type roomPoHasManyRoomUsersTx struct{ tx *gorm.Association }

func (a roomPoHasManyRoomUsersTx) Find() (result []*model.RoomUserPo, err error) {
	return result, a.tx.Find(&result)
}

func (a roomPoHasManyRoomUsersTx) Append(values ...*model.RoomUserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a roomPoHasManyRoomUsersTx) Replace(values ...*model.RoomUserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a roomPoHasManyRoomUsersTx) Delete(values ...*model.RoomUserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a roomPoHasManyRoomUsersTx) Clear() error {
	return a.tx.Clear()
}

func (a roomPoHasManyRoomUsersTx) Count() int64 {
	return a.tx.Count()
}

type roomPoBelongsToGame struct {
	db *gorm.DB

	field.RelationField
}

func (a roomPoBelongsToGame) Where(conds ...field.Expr) *roomPoBelongsToGame {
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

func (a roomPoBelongsToGame) WithContext(ctx context.Context) *roomPoBelongsToGame {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a roomPoBelongsToGame) Session(session *gorm.Session) *roomPoBelongsToGame {
	a.db = a.db.Session(session)
	return &a
}

func (a roomPoBelongsToGame) Model(m *model.RoomPo) *roomPoBelongsToGameTx {
	return &roomPoBelongsToGameTx{a.db.Model(m).Association(a.Name())}
}

type roomPoBelongsToGameTx struct{ tx *gorm.Association }

func (a roomPoBelongsToGameTx) Find() (result *model.GamePo, err error) {
	return result, a.tx.Find(&result)
}

func (a roomPoBelongsToGameTx) Append(values ...*model.GamePo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a roomPoBelongsToGameTx) Replace(values ...*model.GamePo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a roomPoBelongsToGameTx) Delete(values ...*model.GamePo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a roomPoBelongsToGameTx) Clear() error {
	return a.tx.Clear()
}

func (a roomPoBelongsToGameTx) Count() int64 {
	return a.tx.Count()
}

type roomPoBelongsToOwner struct {
	db *gorm.DB

	field.RelationField
}

func (a roomPoBelongsToOwner) Where(conds ...field.Expr) *roomPoBelongsToOwner {
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

func (a roomPoBelongsToOwner) WithContext(ctx context.Context) *roomPoBelongsToOwner {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a roomPoBelongsToOwner) Session(session *gorm.Session) *roomPoBelongsToOwner {
	a.db = a.db.Session(session)
	return &a
}

func (a roomPoBelongsToOwner) Model(m *model.RoomPo) *roomPoBelongsToOwnerTx {
	return &roomPoBelongsToOwnerTx{a.db.Model(m).Association(a.Name())}
}

type roomPoBelongsToOwnerTx struct{ tx *gorm.Association }

func (a roomPoBelongsToOwnerTx) Find() (result *model.UserPo, err error) {
	return result, a.tx.Find(&result)
}

func (a roomPoBelongsToOwnerTx) Append(values ...*model.UserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a roomPoBelongsToOwnerTx) Replace(values ...*model.UserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a roomPoBelongsToOwnerTx) Delete(values ...*model.UserPo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a roomPoBelongsToOwnerTx) Clear() error {
	return a.tx.Clear()
}

func (a roomPoBelongsToOwnerTx) Count() int64 {
	return a.tx.Count()
}

type roomPoDo struct{ gen.DO }

type IRoomPoDo interface {
	gen.SubQuery
	Debug() IRoomPoDo
	WithContext(ctx context.Context) IRoomPoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IRoomPoDo
	WriteDB() IRoomPoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IRoomPoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRoomPoDo
	Not(conds ...gen.Condition) IRoomPoDo
	Or(conds ...gen.Condition) IRoomPoDo
	Select(conds ...field.Expr) IRoomPoDo
	Where(conds ...gen.Condition) IRoomPoDo
	Order(conds ...field.Expr) IRoomPoDo
	Distinct(cols ...field.Expr) IRoomPoDo
	Omit(cols ...field.Expr) IRoomPoDo
	Join(table schema.Tabler, on ...field.Expr) IRoomPoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRoomPoDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRoomPoDo
	Group(cols ...field.Expr) IRoomPoDo
	Having(conds ...gen.Condition) IRoomPoDo
	Limit(limit int) IRoomPoDo
	Offset(offset int) IRoomPoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRoomPoDo
	Unscoped() IRoomPoDo
	Create(values ...*model.RoomPo) error
	CreateInBatches(values []*model.RoomPo, batchSize int) error
	Save(values ...*model.RoomPo) error
	First() (*model.RoomPo, error)
	Take() (*model.RoomPo, error)
	Last() (*model.RoomPo, error)
	Find() ([]*model.RoomPo, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RoomPo, err error)
	FindInBatches(result *[]*model.RoomPo, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.RoomPo) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRoomPoDo
	Assign(attrs ...field.AssignExpr) IRoomPoDo
	Joins(fields ...field.RelationField) IRoomPoDo
	Preload(fields ...field.RelationField) IRoomPoDo
	FirstOrInit() (*model.RoomPo, error)
	FirstOrCreate() (*model.RoomPo, error)
	FindByPage(offset int, limit int) (result []*model.RoomPo, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRoomPoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (r roomPoDo) Debug() IRoomPoDo {
	return r.withDO(r.DO.Debug())
}

func (r roomPoDo) WithContext(ctx context.Context) IRoomPoDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r roomPoDo) ReadDB() IRoomPoDo {
	return r.Clauses(dbresolver.Read)
}

func (r roomPoDo) WriteDB() IRoomPoDo {
	return r.Clauses(dbresolver.Write)
}

func (r roomPoDo) Session(config *gorm.Session) IRoomPoDo {
	return r.withDO(r.DO.Session(config))
}

func (r roomPoDo) Clauses(conds ...clause.Expression) IRoomPoDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r roomPoDo) Returning(value interface{}, columns ...string) IRoomPoDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r roomPoDo) Not(conds ...gen.Condition) IRoomPoDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r roomPoDo) Or(conds ...gen.Condition) IRoomPoDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r roomPoDo) Select(conds ...field.Expr) IRoomPoDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r roomPoDo) Where(conds ...gen.Condition) IRoomPoDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r roomPoDo) Order(conds ...field.Expr) IRoomPoDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r roomPoDo) Distinct(cols ...field.Expr) IRoomPoDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r roomPoDo) Omit(cols ...field.Expr) IRoomPoDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r roomPoDo) Join(table schema.Tabler, on ...field.Expr) IRoomPoDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r roomPoDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRoomPoDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r roomPoDo) RightJoin(table schema.Tabler, on ...field.Expr) IRoomPoDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r roomPoDo) Group(cols ...field.Expr) IRoomPoDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r roomPoDo) Having(conds ...gen.Condition) IRoomPoDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r roomPoDo) Limit(limit int) IRoomPoDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r roomPoDo) Offset(offset int) IRoomPoDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r roomPoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRoomPoDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r roomPoDo) Unscoped() IRoomPoDo {
	return r.withDO(r.DO.Unscoped())
}

func (r roomPoDo) Create(values ...*model.RoomPo) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r roomPoDo) CreateInBatches(values []*model.RoomPo, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r roomPoDo) Save(values ...*model.RoomPo) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r roomPoDo) First() (*model.RoomPo, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomPo), nil
	}
}

func (r roomPoDo) Take() (*model.RoomPo, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomPo), nil
	}
}

func (r roomPoDo) Last() (*model.RoomPo, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomPo), nil
	}
}

func (r roomPoDo) Find() ([]*model.RoomPo, error) {
	result, err := r.DO.Find()
	return result.([]*model.RoomPo), err
}

func (r roomPoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.RoomPo, err error) {
	buf := make([]*model.RoomPo, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r roomPoDo) FindInBatches(result *[]*model.RoomPo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r roomPoDo) Attrs(attrs ...field.AssignExpr) IRoomPoDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r roomPoDo) Assign(attrs ...field.AssignExpr) IRoomPoDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r roomPoDo) Joins(fields ...field.RelationField) IRoomPoDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r roomPoDo) Preload(fields ...field.RelationField) IRoomPoDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r roomPoDo) FirstOrInit() (*model.RoomPo, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomPo), nil
	}
}

func (r roomPoDo) FirstOrCreate() (*model.RoomPo, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.RoomPo), nil
	}
}

func (r roomPoDo) FindByPage(offset int, limit int) (result []*model.RoomPo, count int64, err error) {
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

func (r roomPoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r roomPoDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r roomPoDo) Delete(models ...*model.RoomPo) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *roomPoDo) withDO(do gen.Dao) *roomPoDo {
	r.DO = *do.(*gen.DO)
	return r
}
