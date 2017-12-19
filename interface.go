// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/lingochamp/core"
)

// Interface defines the interface which Engine, EngineGroup and Session will implementate.
type Interface interface {
	AllCols() *Session
	Alias(alias string) *Session
	Asc(colNames ...string) *Session
	BufferSize(size int) *Session
	Cols(columns ...string) *Session
	Count(context.Context, ...interface{}) (int64, error)
	CreateIndexes(ctx context.Context, bean interface{}) error
	CreateUniques(ctx context.Context, bean interface{}) error
	Decr(column string, arg ...interface{}) *Session
	Desc(...string) *Session
	Delete(context.Context, interface{}) (int64, error)
	Distinct(columns ...string) *Session
	DropIndexes(ctx context.Context, bean interface{}) error
	Exec(context.Context, string, ...interface{}) (sql.Result, error)
	Exist(ctx context.Context, bean ...interface{}) (bool, error)
	Find(context.Context, interface{}, ...interface{}) error
	Get(context.Context, interface{}) (bool, error)
	GroupBy(keys string) *Session
	ID(interface{}) *Session
	In(string, ...interface{}) *Session
	Incr(column string, arg ...interface{}) *Session
	Insert(context.Context, ...interface{}) (int64, error)
	InsertOne(context.Context, interface{}) (int64, error)
	IsTableEmpty(ctx context.Context, bean interface{}) (bool, error)
	IsTableExist(ctx context.Context, beanOrTableName interface{}) (bool, error)
	Iterate(context.Context, interface{}, IterFunc) error
	Limit(int, ...int) *Session
	NoAutoCondition(...bool) *Session
	NotIn(string, ...interface{}) *Session
	Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *Session
	Omit(columns ...string) *Session
	OrderBy(order string) *Session
	Ping(ctx context.Context) error
	Query(ctx context.Context, sqlOrAgrs ...interface{}) (resultsSlice []map[string][]byte, err error)
	QueryInterface(ctx context.Context, sqlorArgs ...interface{}) ([]map[string]interface{}, error)
	QueryString(ctx context.Context, sqlorArgs ...interface{}) ([]map[string]string, error)
	Rows(ctx context.Context, bean interface{}) (*Rows, error)
	SetExpr(string, string) *Session
	SQL(interface{}, ...interface{}) *Session
	Sum(ctx context.Context, bean interface{}, colName string) (float64, error)
	SumInt(ctx context.Context, bean interface{}, colName string) (int64, error)
	Sums(ctx context.Context, bean interface{}, colNames ...string) ([]float64, error)
	SumsInt(ctx context.Context, bean interface{}, colNames ...string) ([]int64, error)
	Table(tableNameOrBean interface{}) *Session
	Unscoped() *Session
	Update(ctx context.Context, bean interface{}, condiBeans ...interface{}) (int64, error)
	UseBool(...string) *Session
	Where(interface{}, ...interface{}) *Session
}

// EngineInterface defines the interface which Engine, EngineGroup will implementate.
type EngineInterface interface {
	Interface

	Before(func(interface{})) *Session
	Charset(charset string) *Session
	CreateTables(context.Context, ...interface{}) error
	DBMetas(context.Context) ([]*core.Table, error)
	Dialect() core.Dialect
	DropTables(context.Context, ...interface{}) error
	DumpAllToFile(ctx context.Context, fp string, tp ...core.DbType) error
	GetColumnMapper() core.IMapper
	GetDefaultCacher() core.Cacher
	GetTableMapper() core.IMapper
	GetTZDatabase() *time.Location
	GetTZLocation() *time.Location
	NewSession() *Session
	NoAutoTime() *Session
	Quote(string) string
	SetDefaultCacher(core.Cacher)
	SetLogLevel(core.LogLevel)
	SetMapper(core.IMapper)
	SetTZDatabase(tz *time.Location)
	SetTZLocation(tz *time.Location)
	ShowSQL(show ...bool)
	Sync(context.Context, ...interface{}) error
	Sync2(context.Context, ...interface{}) error
	StoreEngine(storeEngine string) *Session
	TableInfo(bean interface{}) *Table
	UnMapType(reflect.Type)
}

var (
	_ Interface       = &Session{}
	_ EngineInterface = &Engine{}
	_ EngineInterface = &EngineGroup{}
)
