// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type UserinfoDelete struct {
		Uid   int64 `xorm:"id pk not null autoincr"`
		IsMan bool
	}

	assert.NoError(t, testEngine.Sync2(context.Background(), new(UserinfoDelete)))

	user := UserinfoDelete{Uid: 1}
	cnt, err := testEngine.Insert(context.Background(), &user)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Delete(context.Background(), &UserinfoDelete{Uid: user.Uid})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	user.Uid = 0
	user.IsMan = true
	has, err := testEngine.ID(1).Get(context.Background(), &user)
	assert.NoError(t, err)
	assert.False(t, has)

	cnt, err = testEngine.Insert(context.Background(), &user)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Where("id=?", user.Uid).Delete(context.Background(), &UserinfoDelete{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	user.Uid = 0
	user.IsMan = true
	has, err = testEngine.ID(2).Get(context.Background(), &user)
	assert.NoError(t, err)
	assert.False(t, has)
}

func TestDeleted(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type Deleted struct {
		Id        int64 `xorm:"pk"`
		Name      string
		DeletedAt time.Time `xorm:"deleted"`
	}

	err := testEngine.DropTables(context.Background(), &Deleted{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(context.Background(), &Deleted{})
	assert.NoError(t, err)

	_, err = testEngine.InsertOne(context.Background(), &Deleted{Id: 1, Name: "11111"})
	assert.NoError(t, err)

	_, err = testEngine.InsertOne(context.Background(), &Deleted{Id: 2, Name: "22222"})
	assert.NoError(t, err)

	_, err = testEngine.InsertOne(context.Background(), &Deleted{Id: 3, Name: "33333"})
	assert.NoError(t, err)

	// Test normal Find()
	var records1 []Deleted
	err = testEngine.Where("`"+testEngine.GetColumnMapper().Obj2Table("Id")+"` > 0").Find(context.Background(), &records1, &Deleted{})
	assert.EqualValues(t, 3, len(records1))

	// Test normal Get()
	record1 := &Deleted{}
	has, err := testEngine.ID(1).Get(context.Background(), record1)
	assert.NoError(t, err)
	assert.True(t, has)

	// Test Delete() with deleted
	affected, err := testEngine.ID(1).Delete(context.Background(), &Deleted{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affected)

	has, err = testEngine.ID(1).Get(context.Background(), &Deleted{})
	assert.NoError(t, err)
	assert.False(t, has)

	var records2 []Deleted
	err = testEngine.Where("`"+testEngine.GetColumnMapper().Obj2Table("Id")+"` > 0").Find(context.Background(), &records2)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(records2))

	// Test no rows affected after Delete() again.
	affected, err = testEngine.ID(1).Delete(context.Background(), &Deleted{})
	assert.NoError(t, err)
	assert.EqualValues(t, 0, affected)

	// Deleted.DeletedAt must not be updated.
	affected, err = testEngine.ID(2).Update(context.Background(), &Deleted{Name: "2", DeletedAt: time.Now()})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affected)

	record2 := &Deleted{}
	has, err = testEngine.ID(2).Get(context.Background(), record2)
	assert.NoError(t, err)
	assert.True(t, record2.DeletedAt.IsZero())

	// Test find all records whatever `deleted`.
	var unscopedRecords1 []Deleted
	err = testEngine.Unscoped().Where("`"+testEngine.GetColumnMapper().Obj2Table("Id")+"` > 0").Find(context.Background(), &unscopedRecords1, &Deleted{})
	assert.NoError(t, err)
	assert.EqualValues(t, 3, len(unscopedRecords1))

	// Delete() must really delete a record with Unscoped()
	affected, err = testEngine.Unscoped().ID(1).Delete(context.Background(), &Deleted{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, affected)

	var unscopedRecords2 []Deleted
	err = testEngine.Unscoped().Where("`"+testEngine.GetColumnMapper().Obj2Table("Id")+"` > 0").Find(context.Background(), &unscopedRecords2, &Deleted{})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(unscopedRecords2))

	var records3 []Deleted
	err = testEngine.Where("`"+testEngine.GetColumnMapper().Obj2Table("Id")+"` > 0").And("`"+testEngine.GetColumnMapper().Obj2Table("Id")+"`> 1").
		Or("`"+testEngine.GetColumnMapper().Obj2Table("Id")+"` = ?", 3).Find(context.Background(), &records3)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(records3))
}

func TestCacheDelete(t *testing.T) {
	assert.NoError(t, prepareEngine())

	oldCacher := testEngine.GetDefaultCacher()
	cacher := NewLRUCacher(NewMemoryStore(), 1000)
	testEngine.SetDefaultCacher(cacher)

	type CacheDeleteStruct struct {
		Id int64
	}

	err := testEngine.CreateTables(context.Background(), &CacheDeleteStruct{})
	assert.NoError(t, err)

	_, err = testEngine.Insert(context.Background(), &CacheDeleteStruct{})
	assert.NoError(t, err)

	aff, err := testEngine.Delete(context.Background(), &CacheDeleteStruct{
		Id: 1,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, aff, 1)

	aff, err = testEngine.Unscoped().Delete(context.Background(), &CacheDeleteStruct{
		Id: 1,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, aff, 0)

	testEngine.SetDefaultCacher(oldCacher)
}

func TestUnscopeDelete(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type UnscopeDeleteStruct struct {
		Id        int64
		Name      string
		DeletedAt time.Time `xorm:"deleted"`
	}

	assertSync(t, new(UnscopeDeleteStruct))

	cnt, err := testEngine.Insert(context.Background(), &UnscopeDeleteStruct{
		Name: "test",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var nowUnix = time.Now().Unix()
	var s UnscopeDeleteStruct
	cnt, err = testEngine.ID(1).Delete(context.Background(), &s)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	assert.EqualValues(t, nowUnix, s.DeletedAt.Unix())

	var s1 UnscopeDeleteStruct
	has, err := testEngine.ID(1).Get(context.Background(), &s1)
	assert.NoError(t, err)
	assert.False(t, has)

	var s2 UnscopeDeleteStruct
	has, err = testEngine.ID(1).Unscoped().Get(context.Background(), &s2)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "test", s2.Name)
	assert.EqualValues(t, nowUnix, s2.DeletedAt.Unix())

	cnt, err = testEngine.ID(1).Unscoped().Delete(context.Background(), new(UnscopeDeleteStruct))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var s3 UnscopeDeleteStruct
	has, err = testEngine.ID(1).Get(context.Background(), &s3)
	assert.NoError(t, err)
	assert.False(t, has)

	var s4 UnscopeDeleteStruct
	has, err = testEngine.ID(1).Unscoped().Get(context.Background(), &s4)
	assert.NoError(t, err)
	assert.False(t, has)
}
