// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-xorm/core"
	"github.com/stretchr/testify/assert"
)

func TestJoinLimit(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type Salary struct {
		Id  int64
		Lid int64
	}

	type CheckList struct {
		Id  int64
		Eid int64
	}

	type Empsetting struct {
		Id   int64
		Name string
	}

	assert.NoError(t, testEngine.Sync2(context.Background(), new(Salary), new(CheckList), new(Empsetting)))

	var emp Empsetting
	cnt, err := testEngine.Insert(context.Background(), &emp)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var checklist = CheckList{
		Eid: emp.Id,
	}
	cnt, err = testEngine.Insert(context.Background(), &checklist)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var salary = Salary{
		Lid: checklist.Id,
	}
	cnt, err = testEngine.Insert(context.Background(), &salary)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var salaries []Salary
	err = testEngine.Table("salary").
		Join("INNER", "check_list", "check_list.id = salary.lid").
		Join("LEFT", "empsetting", "empsetting.id = check_list.eid").
		Limit(10, 0).
		Find(context.Background(), &salaries)
	assert.NoError(t, err)
}

func assertSync(t *testing.T, beans ...interface{}) {
	for _, bean := range beans {
		assert.NoError(t, testEngine.DropTables(context.Background(), bean))
		assert.NoError(t, testEngine.Sync2(context.Background(), bean))
	}
}

func TestWhere(t *testing.T) {
	assert.NoError(t, prepareEngine())

	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.Where("(id) > ?", 2).Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	err = testEngine.Where("(id) > ?", 2).And("(id) < ?", 10).Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)
}

func TestFind(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)

	err := testEngine.Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}

	users2 := make([]Userinfo, 0)
	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	err = testEngine.SQL("select * from "+testEngine.Quote(userinfo)).Find(context.Background(), &users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func TestFind2(t *testing.T) {
	assert.NoError(t, prepareEngine())
	users := make([]*Userinfo, 0)

	assertSync(t, new(Userinfo))

	err := testEngine.Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

type Team struct {
	Id int64
}

type TeamUser struct {
	OrgId  int64
	Uid    int64
	TeamId int64
}

func TestFind3(t *testing.T) {
	assert.NoError(t, prepareEngine())
	err := testEngine.Sync2(context.Background(), new(Team), new(TeamUser))
	if err != nil {
		t.Error(err)
		panic(err.Error())
	}

	var teams []Team
	err = testEngine.Cols("`team`.id").
		Where("`team_user`.org_id=?", 1).
		And("`team_user`.uid=?", 2).
		Join("INNER", "`team_user`", "`team_user`.team_id=`team`.id").
		Find(context.Background(), &teams)
	if err != nil {
		t.Error(err)
		panic(err.Error())
	}
}

func TestFindMap(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	users := make(map[int64]Userinfo)
	err := testEngine.Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

func TestFindMap2(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	users := make(map[int64]*Userinfo)
	err := testEngine.Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	for id, user := range users {
		fmt.Println(id, user)
	}
}

func TestDistinct(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	_, err := testEngine.Insert(context.Background(), &Userinfo{
		Username: "lunny",
	})
	assert.NoError(t, err)

	users := make([]Userinfo, 0)
	departname := testEngine.GetTableMapper().Obj2Table("Departname")
	err = testEngine.Distinct(departname).Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users) != 1 {
		t.Error(err)
		panic(errors.New("should be one record"))
	}

	fmt.Println(users)

	type Depart struct {
		Departname string
	}

	users2 := make([]Depart, 0)
	err = testEngine.Distinct(departname).Table(new(Userinfo)).Find(context.Background(), &users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if len(users2) != 1 {
		t.Error(err)
		panic(errors.New("should be one record"))
	}
	fmt.Println(users2)
}

func TestOrder(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.OrderBy("id desc").Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	users2 := make([]Userinfo, 0)
	err = testEngine.Asc("id", "username").Desc("height").Find(context.Background(), &users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users2)
}

func TestHaving(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.GroupBy("username").Having("username='xlw'").Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	/*users = make([]Userinfo, 0)
	err = testEngine.Cols("id, username").GroupBy("username").Having("username='xlw'").Find(&users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)*/
}

func TestOrderSameMapper(t *testing.T) {
	assert.NoError(t, prepareEngine())
	testEngine.UnMapType(rValue(new(Userinfo)).Type())

	mapper := testEngine.GetTableMapper()
	testEngine.SetMapper(core.SameMapper{})

	defer func() {
		testEngine.UnMapType(rValue(new(Userinfo)).Type())
		testEngine.SetMapper(mapper)
	}()

	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.OrderBy("(id) desc").Find(context.Background(), &users)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users)

	users2 := make([]Userinfo, 0)
	err = testEngine.Asc("(id)", "Username").Desc("Height").Find(context.Background(), &users2)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(users2)
}

func TestHavingSameMapper(t *testing.T) {
	assert.NoError(t, prepareEngine())
	testEngine.UnMapType(rValue(new(Userinfo)).Type())

	mapper := testEngine.GetTableMapper()
	testEngine.SetMapper(core.SameMapper{})
	defer func() {
		testEngine.UnMapType(rValue(new(Userinfo)).Type())
		testEngine.SetMapper(mapper)
	}()
	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.GroupBy("`Username`").Having("`Username`='xlw'").Find(context.Background(), &users)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(users)
}

func TestFindInts(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var idsInt64 []int64
	err := testEngine.Table(userinfo).Cols("id").Desc("id").Find(context.Background(), &idsInt64)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsInt64)

	var idsInt32 []int32
	err = testEngine.Table(userinfo).Cols("id").Desc("id").Find(context.Background(), &idsInt32)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsInt32)

	var idsInt []int
	err = testEngine.Table(userinfo).Cols("id").Desc("id").Find(context.Background(), &idsInt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsInt)

	var idsUint []uint
	err = testEngine.Table(userinfo).Cols("id").Desc("id").Find(context.Background(), &idsUint)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsUint)

	type MyInt int
	var idsMyInt []MyInt
	err = testEngine.Table(userinfo).Cols("id").Desc("id").Find(context.Background(), &idsMyInt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsMyInt)
}

func TestFindStrings(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))
	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	username := testEngine.GetColumnMapper().Obj2Table("Username")
	var idsString []string
	err := testEngine.Table(userinfo).Cols(username).Desc("id").Find(context.Background(), &idsString)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsString)
}

func TestFindMyString(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))
	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	username := testEngine.GetColumnMapper().Obj2Table("Username")

	var idsMyString []MyString
	err := testEngine.Table(userinfo).Cols(username).Desc("id").Find(context.Background(), &idsMyString)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsMyString)
}

func TestFindInterface(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	username := testEngine.GetColumnMapper().Obj2Table("Username")
	var idsInterface []interface{}
	err := testEngine.Table(userinfo).Cols(username).Desc("id").Find(context.Background(), &idsInterface)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(idsInterface)
}

func TestFindSliceBytes(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var ids [][][]byte
	err := testEngine.Table(userinfo).Desc("id").Find(context.Background(), &ids)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range ids {
		fmt.Println(record)
	}
}

func TestFindSlicePtrString(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var ids [][]*string
	err := testEngine.Table(userinfo).Desc("id").Find(context.Background(), &ids)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range ids {
		fmt.Println(record)
	}
}

func TestFindMapBytes(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var ids []map[string][]byte
	err := testEngine.Table(userinfo).Desc("id").Find(context.Background(), &ids)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range ids {
		fmt.Println(record)
	}
}

func TestFindMapPtrString(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var ids []map[string]*string
	err := testEngine.Table(userinfo).Desc("id").Find(context.Background(), &ids)
	assert.NoError(t, err)
	for _, record := range ids {
		fmt.Println(record)
	}
}

func TestFindBit(t *testing.T) {
	type FindBitStruct struct {
		Id  int64
		Msg bool `xorm:"bit"`
	}

	assert.NoError(t, prepareEngine())
	assertSync(t, new(FindBitStruct))

	cnt, err := testEngine.Insert(context.Background(), []FindBitStruct{
		{
			Msg: false,
		},
		{
			Msg: true,
		},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)

	var results = make([]FindBitStruct, 0, 2)
	err = testEngine.Find(context.Background(), &results)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(results))
}
