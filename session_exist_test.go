// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistStruct(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type RecordExist struct {
		Id   int64
		Name string
	}

	assertSync(t, new(RecordExist))

	has, err := testEngine.Exist(context.Background(), new(RecordExist))
	assert.NoError(t, err)
	assert.False(t, has)

	cnt, err := testEngine.Insert(context.Background(), &RecordExist{
		Name: "test1",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	has, err = testEngine.Exist(context.Background(), new(RecordExist))
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = testEngine.Exist(context.Background(), &RecordExist{
		Name: "test1",
	})
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = testEngine.Exist(context.Background(), &RecordExist{
		Name: "test2",
	})
	assert.NoError(t, err)
	assert.False(t, has)

	has, err = testEngine.Where("name = ?", "test1").Exist(context.Background(), &RecordExist{})
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = testEngine.Where("name = ?", "test2").Exist(context.Background(), &RecordExist{})
	assert.NoError(t, err)
	assert.False(t, has)

	has, err = testEngine.SQL("select * from record_exist where name = ?", "test1").Exist(context.Background())
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = testEngine.SQL("select * from record_exist where name = ?", "test2").Exist(context.Background())
	assert.NoError(t, err)
	assert.False(t, has)

	has, err = testEngine.Table("record_exist").Exist(context.Background())
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = testEngine.Table("record_exist").Where("name = ?", "test1").Exist(context.Background())
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = testEngine.Table("record_exist").Where("name = ?", "test2").Exist(context.Background())
	assert.NoError(t, err)
	assert.False(t, has)
}
