// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lingochamp/core"
	"github.com/stretchr/testify/assert"
)

func TestGetVar(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type GetVar struct {
		Id      int64  `xorm:"autoincr pk"`
		Msg     string `xorm:"varchar(255)"`
		Age     int
		Money   float32
		Created time.Time `xorm:"created"`
	}

	assert.NoError(t, testEngine.Sync2(context.Background(), new(GetVar)))

	var data = GetVar{
		Msg:   "hi",
		Age:   28,
		Money: 1.5,
	}
	_, err := testEngine.InsertOne(context.Background(), data)
	assert.NoError(t, err)

	var msg string
	has, err := testEngine.Table("get_var").Cols("msg").Get(context.Background(), &msg)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, "hi", msg)

	var age int
	has, err = testEngine.Table("get_var").Cols("age").Get(context.Background(), &age)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, 28, age)

	var age2 int64
	has, err = testEngine.Table("get_var").Cols("age").
		Where("age > ?", 20).
		And("age < ?", 30).
		Get(context.Background(), &age2)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.EqualValues(t, 28, age2)

	var money float64
	has, err = testEngine.Table("get_var").Cols("money").Get(context.Background(), &money)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, "1.5", fmt.Sprintf("%.1f", money))

	var valuesString = make(map[string]string)
	has, err = testEngine.Table("get_var").Get(context.Background(), &valuesString)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, 5, len(valuesString))
	assert.Equal(t, "1", valuesString["id"])
	assert.Equal(t, "hi", valuesString["msg"])
	assert.Equal(t, "28", valuesString["age"])
	assert.Equal(t, "1.5", valuesString["money"])

	// for mymysql driver, interface{} will be []byte, so ignore it currently
	if testEngine.Dialect().DriverName() != "mymysql" {
		var valuesInter = make(map[string]interface{})
		has, err = testEngine.Table("get_var").Where("id = ?", 1).Select("*").Get(context.Background(), &valuesInter)
		assert.NoError(t, err)
		assert.Equal(t, true, has)
		assert.Equal(t, 5, len(valuesInter))
		assert.EqualValues(t, 1, valuesInter["id"])
		assert.Equal(t, "hi", fmt.Sprintf("%s", valuesInter["msg"]))
		assert.EqualValues(t, 28, valuesInter["age"])
		assert.Equal(t, "1.5", fmt.Sprintf("%v", valuesInter["money"]))
	}

	var valuesSliceString = make([]string, 5)
	has, err = testEngine.Table("get_var").Get(context.Background(), &valuesSliceString)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, "1", valuesSliceString[0])
	assert.Equal(t, "hi", valuesSliceString[1])
	assert.Equal(t, "28", valuesSliceString[2])
	assert.Equal(t, "1.5", valuesSliceString[3])

	var valuesSliceInter = make([]interface{}, 5)
	has, err = testEngine.Table("get_var").Get(context.Background(), &valuesSliceInter)
	assert.NoError(t, err)
	assert.Equal(t, true, has)

	v1, err := convertInt(valuesSliceInter[0])
	assert.NoError(t, err)
	assert.EqualValues(t, 1, v1)

	assert.Equal(t, "hi", fmt.Sprintf("%s", valuesSliceInter[1]))

	v3, err := convertInt(valuesSliceInter[2])
	assert.NoError(t, err)
	assert.EqualValues(t, 28, v3)

	v4, err := convertFloat(valuesSliceInter[3])
	assert.NoError(t, err)
	assert.Equal(t, "1.5", fmt.Sprintf("%v", v4))
}

func TestGetStruct(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type UserinfoGet struct {
		Uid   int `xorm:"pk autoincr"`
		IsMan bool
	}

	assert.NoError(t, testEngine.Sync2(context.Background(), new(UserinfoGet)))

	var err error
	if testEngine.Dialect().DBType() == core.MSSQL {
		_, err = testEngine.Exec(context.Background(), "SET IDENTITY_INSERT userinfo_get ON")
		assert.NoError(t, err)
	}
	cnt, err := testEngine.Insert(context.Background(), &UserinfoGet{Uid: 2})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	user := UserinfoGet{Uid: 2}
	has, err := testEngine.Get(context.Background(), &user)
	assert.NoError(t, err)
	assert.True(t, has)

	type NoIdUser struct {
		User   string `xorm:"unique"`
		Remain int64
		Total  int64
	}

	assert.NoError(t, testEngine.Sync2(context.Background(), &NoIdUser{}))

	userCol := testEngine.GetColumnMapper().Obj2Table("User")
	_, err = testEngine.Where("`"+userCol+"` = ?", "xlw").Delete(context.Background(), &NoIdUser{})
	assert.NoError(t, err)

	cnt, err = testEngine.Insert(context.Background(), &NoIdUser{"xlw", 20, 100})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	noIdUser := new(NoIdUser)
	has, err = testEngine.Where("`"+userCol+"` = ?", "xlw").Get(context.Background(), noIdUser)
	assert.NoError(t, err)
	assert.True(t, has)
}

func TestGetSlice(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type UserinfoSlice struct {
		Uid   int `xorm:"pk autoincr"`
		IsMan bool
	}

	assertSync(t, new(UserinfoSlice))

	var users []UserinfoSlice
	has, err := testEngine.Get(context.Background(), &users)
	assert.False(t, has)
	assert.Error(t, err)
}

func TestGetError(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type GetError struct {
		Uid   int `xorm:"pk autoincr"`
		IsMan bool
	}

	assertSync(t, new(GetError))

	var info = new(GetError)
	has, err := testEngine.Get(context.Background(), &info)
	assert.False(t, has)
	assert.Error(t, err)

	has, err = testEngine.Get(context.Background(), info)
	assert.False(t, has)
	assert.NoError(t, err)
}

func TestJSONString(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type JsonString struct {
		Id      int64
		Content string `xorm:"json"`
	}
	type JsonJson struct {
		Id      int64
		Content []string `xorm:"json"`
	}

	assertSync(t, new(JsonJson))

	_, err := testEngine.Insert(context.Background(), &JsonJson{
		Content: []string{"1", "2"},
	})
	assert.NoError(t, err)

	var js JsonString
	has, err := testEngine.Table("json_json").Get(context.Background(), &js)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 1, js.Id)
	assert.EqualValues(t, `["1","2"]`, js.Content)

	var jss []JsonString
	err = testEngine.Table("json_json").Find(context.Background(), &jss)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(jss))
	assert.EqualValues(t, `["1","2"]`, jss[0].Content)
}
