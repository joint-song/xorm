// Copyright 2016 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-xorm/core"
)

// Ping test if database is ok
func (session *Session) Ping(ctx context.Context) error {
	if session.isAutoClose {
		defer session.Close()
	}

	session.engine.logger.Infof("PING DATABASE %v", session.engine.DriverName())
	return session.DB().PingContext(ctx)
}

// CreateTable create a table according a bean
func (session *Session) CreateTable(ctx context.Context, bean interface{}) error {
	if session.isAutoClose {
		defer session.Close()
	}

	return session.createTable(ctx, bean)
}

func (session *Session) createTable(ctx context.Context, bean interface{}) error {
	v := rValue(bean)
	if err := session.statement.setRefValue(v); err != nil {
		return err
	}

	sqlStr := session.statement.genCreateTableSQL()
	_, err := session.exec(ctx, sqlStr)
	return err
}

// CreateIndexes create indexes
func (session *Session) CreateIndexes(ctx context.Context, bean interface{}) error {
	if session.isAutoClose {
		defer session.Close()
	}

	return session.createIndexes(ctx, bean)
}

func (session *Session) createIndexes(ctx context.Context, bean interface{}) error {
	v := rValue(bean)
	if err := session.statement.setRefValue(v); err != nil {
		return err
	}

	sqls := session.statement.genIndexSQL()
	for _, sqlStr := range sqls {
		_, err := session.exec(ctx, sqlStr)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateUniques create uniques
func (session *Session) CreateUniques(ctx context.Context, bean interface{}) error {
	if session.isAutoClose {
		defer session.Close()
	}
	return session.createUniques(ctx, bean)
}

func (session *Session) createUniques(ctx context.Context, bean interface{}) error {
	v := rValue(bean)
	if err := session.statement.setRefValue(v); err != nil {
		return err
	}

	sqls := session.statement.genUniqueSQL()
	for _, sqlStr := range sqls {
		_, err := session.exec(ctx, sqlStr)
		if err != nil {
			return err
		}
	}
	return nil
}

// DropIndexes drop indexes
func (session *Session) DropIndexes(ctx context.Context, bean interface{}) error {
	if session.isAutoClose {
		defer session.Close()
	}

	return session.dropIndexes(ctx, bean)
}

func (session *Session) dropIndexes(ctx context.Context, bean interface{}) error {
	v := rValue(bean)
	if err := session.statement.setRefValue(v); err != nil {
		return err
	}

	sqls := session.statement.genDelIndexSQL()
	for _, sqlStr := range sqls {
		_, err := session.exec(ctx, sqlStr)
		if err != nil {
			return err
		}
	}
	return nil
}

// DropTable drop table will drop table if exist, if drop failed, it will return error
func (session *Session) DropTable(ctx context.Context, beanOrTableName interface{}) error {
	if session.isAutoClose {
		defer session.Close()
	}

	return session.dropTable(ctx, beanOrTableName)
}

func (session *Session) dropTable(ctx context.Context, beanOrTableName interface{}) error {
	tableName, err := session.engine.tableName(beanOrTableName)
	if err != nil {
		return err
	}

	var needDrop = true
	if !session.engine.dialect.SupportDropIfExists() {
		sqlStr, args := session.engine.dialect.TableCheckSql(tableName)
		results, err := session.queryBytes(ctx, sqlStr, args...)
		if err != nil {
			return err
		}
		needDrop = len(results) > 0
	}

	if needDrop {
		sqlStr := session.engine.Dialect().DropTableSql(tableName)
		_, err = session.exec(ctx, sqlStr)
		return err
	}
	return nil
}

// IsTableExist if a table is exist
func (session *Session) IsTableExist(ctx context.Context, beanOrTableName interface{}) (bool, error) {
	if session.isAutoClose {
		defer session.Close()
	}

	tableName, err := session.engine.tableName(beanOrTableName)
	if err != nil {
		return false, err
	}

	return session.isTableExist(ctx, tableName)
}

func (session *Session) isTableExist(ctx context.Context, tableName string) (bool, error) {
	sqlStr, args := session.engine.dialect.TableCheckSql(tableName)
	results, err := session.queryBytes(ctx, sqlStr, args...)
	return len(results) > 0, err
}

// IsTableEmpty if table have any records
func (session *Session) IsTableEmpty(ctx context.Context, bean interface{}) (bool, error) {
	v := rValue(bean)
	t := v.Type()

	if t.Kind() == reflect.String {
		if session.isAutoClose {
			defer session.Close()
		}
		return session.isTableEmpty(ctx, bean.(string))
	} else if t.Kind() == reflect.Struct {
		rows, err := session.Count(ctx, bean)
		return rows == 0, err
	}
	return false, errors.New("bean should be a struct or struct's point")
}

func (session *Session) isTableEmpty(ctx context.Context, tableName string) (bool, error) {
	var total int64
	sqlStr := fmt.Sprintf("select count(*) from %s", session.engine.Quote(tableName))
	err := session.queryRow(ctx, sqlStr).Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return true, err
	}

	return total == 0, nil
}

// find if index is exist according cols
func (session *Session) isIndexExist2(ctx context.Context, tableName string, cols []string, unique bool) (bool, error) {
	indexes, err := session.engine.dialect.GetIndexes(ctx, tableName)
	if err != nil {
		return false, err
	}

	for _, index := range indexes {
		if sliceEq(index.Cols, cols) {
			if unique {
				return index.Type == core.UniqueType, nil
			}
			return index.Type == core.IndexType, nil
		}
	}
	return false, nil
}

func (session *Session) addColumn(ctx context.Context, colName string) error {
	col := session.statement.RefTable.GetColumn(colName)
	sql, args := session.statement.genAddColumnStr(col)
	_, err := session.exec(ctx, sql, args...)
	return err
}

func (session *Session) addIndex(ctx context.Context, tableName, idxName string) error {
	index := session.statement.RefTable.Indexes[idxName]
	sqlStr := session.engine.dialect.CreateIndexSql(tableName, index)
	_, err := session.exec(ctx, sqlStr)
	return err
}

func (session *Session) addUnique(ctx context.Context, tableName, uqeName string) error {
	index := session.statement.RefTable.Indexes[uqeName]
	sqlStr := session.engine.dialect.CreateIndexSql(tableName, index)
	_, err := session.exec(ctx, sqlStr)
	return err
}

// Sync2 synchronize structs to database tables
func (session *Session) Sync2(ctx context.Context, beans ...interface{}) error {
	engine := session.engine

	if session.isAutoClose {
		session.isAutoClose = false
		defer session.Close()
	}

	tables, err := engine.DBMetas(ctx)
	if err != nil {
		return err
	}

	var structTables []*core.Table

	for _, bean := range beans {
		v := rValue(bean)
		table, err := engine.mapType(v)
		if err != nil {
			return err
		}
		structTables = append(structTables, table)
		var tbName = session.tbNameNoSchema(table)

		var oriTable *core.Table
		for _, tb := range tables {
			if strings.EqualFold(tb.Name, tbName) {
				oriTable = tb
				break
			}
		}

		if oriTable == nil {
			err = session.StoreEngine(session.statement.StoreEngine).createTable(ctx, bean)
			if err != nil {
				return err
			}

			err = session.createUniques(ctx, bean)
			if err != nil {
				return err
			}

			err = session.createIndexes(ctx, bean)
			if err != nil {
				return err
			}
		} else {
			for _, col := range table.Columns() {
				var oriCol *core.Column
				for _, col2 := range oriTable.Columns() {
					if strings.EqualFold(col.Name, col2.Name) {
						oriCol = col2
						break
					}
				}

				if oriCol != nil {
					expectedType := engine.dialect.SqlType(col)
					curType := engine.dialect.SqlType(oriCol)
					if expectedType != curType {
						if expectedType == core.Text &&
							strings.HasPrefix(curType, core.Varchar) {
							// currently only support mysql & postgres
							if engine.dialect.DBType() == core.MYSQL ||
								engine.dialect.DBType() == core.POSTGRES {
								engine.logger.Infof("Table %s column %s change type from %s to %s\n",
									tbName, col.Name, curType, expectedType)
								_, err = session.exec(ctx, engine.dialect.ModifyColumnSql(table.Name, col))
							} else {
								engine.logger.Warnf("Table %s column %s db type is %s, struct type is %s\n",
									tbName, col.Name, curType, expectedType)
							}
						} else if strings.HasPrefix(curType, core.Varchar) && strings.HasPrefix(expectedType, core.Varchar) {
							if engine.dialect.DBType() == core.MYSQL {
								if oriCol.Length < col.Length {
									engine.logger.Infof("Table %s column %s change type from varchar(%d) to varchar(%d)\n",
										tbName, col.Name, oriCol.Length, col.Length)
									_, err = session.exec(ctx, engine.dialect.ModifyColumnSql(table.Name, col))
								}
							}
						} else {
							if !(strings.HasPrefix(curType, expectedType) && curType[len(expectedType)] == '(') {
								engine.logger.Warnf("Table %s column %s db type is %s, struct type is %s",
									tbName, col.Name, curType, expectedType)
							}
						}
					} else if expectedType == core.Varchar {
						if engine.dialect.DBType() == core.MYSQL {
							if oriCol.Length < col.Length {
								engine.logger.Infof("Table %s column %s change type from varchar(%d) to varchar(%d)\n",
									tbName, col.Name, oriCol.Length, col.Length)
								_, err = session.exec(ctx, engine.dialect.ModifyColumnSql(table.Name, col))
							}
						}
					}
					if col.Default != oriCol.Default {
						engine.logger.Warnf("Table %s Column %s db default is %s, struct default is %s",
							tbName, col.Name, oriCol.Default, col.Default)
					}
					if col.Nullable != oriCol.Nullable {
						engine.logger.Warnf("Table %s Column %s db nullable is %v, struct nullable is %v",
							tbName, col.Name, oriCol.Nullable, col.Nullable)
					}
				} else {
					session.statement.RefTable = table
					session.statement.tableName = tbName
					err = session.addColumn(ctx, col.Name)
				}
				if err != nil {
					return err
				}
			}

			var foundIndexNames = make(map[string]bool)
			var addedNames = make(map[string]*core.Index)

			for name, index := range table.Indexes {
				var oriIndex *core.Index
				for name2, index2 := range oriTable.Indexes {
					if index.Equal(index2) {
						oriIndex = index2
						foundIndexNames[name2] = true
						break
					}
				}

				if oriIndex != nil {
					if oriIndex.Type != index.Type {
						sql := engine.dialect.DropIndexSql(tbName, oriIndex)
						_, err = session.exec(ctx, sql)
						if err != nil {
							return err
						}
						oriIndex = nil
					}
				}

				if oriIndex == nil {
					addedNames[name] = index
				}
			}

			for name2, index2 := range oriTable.Indexes {
				if _, ok := foundIndexNames[name2]; !ok {
					sql := engine.dialect.DropIndexSql(tbName, index2)
					_, err = session.exec(ctx, sql)
					if err != nil {
						return err
					}
				}
			}

			for name, index := range addedNames {
				if index.Type == core.UniqueType {
					session.statement.RefTable = table
					session.statement.tableName = tbName
					err = session.addUnique(ctx, tbName, name)
				} else if index.Type == core.IndexType {
					session.statement.RefTable = table
					session.statement.tableName = tbName
					err = session.addIndex(ctx, tbName, name)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	for _, table := range tables {
		var oriTable *core.Table
		for _, structTable := range structTables {
			if strings.EqualFold(table.Name, session.tbNameNoSchema(structTable)) {
				oriTable = structTable
				break
			}
		}

		if oriTable == nil {
			//engine.LogWarnf("Table %s has no struct to mapping it", table.Name)
			continue
		}

		for _, colName := range table.ColumnsSeq() {
			if oriTable.GetColumn(colName) == nil {
				engine.logger.Warnf("Table %s has column %s but struct has not related field", table.Name, colName)
			}
		}
	}
	return nil
}
