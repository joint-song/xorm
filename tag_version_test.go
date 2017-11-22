// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type VersionS struct {
	Id      int64
	Name    string
	Ver     int       `xorm:"version"`
	Created time.Time `xorm:"created"`
}

func TestVersion1(t *testing.T) {
	assert.NoError(t, prepareEngine())

	err := testEngine.DropTables(context.Background(), new(VersionS))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = testEngine.CreateTables(context.Background(), new(VersionS))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	ver := &VersionS{Name: "sfsfdsfds"}
	_, err = testEngine.Insert(context.Background(), ver)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(ver)
	if ver.Ver != 1 {
		err = errors.New("insert error")
		t.Error(err)
		panic(err)
	}

	newVer := new(VersionS)
	has, err := testEngine.ID(ver.Id).Get(context.Background(), newVer)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	if !has {
		t.Error(errors.New(fmt.Sprintf("no version id is %v", ver.Id)))
		panic(err)
	}
	fmt.Println(newVer)
	if newVer.Ver != 1 {
		err = errors.New("insert error")
		t.Error(err)
		panic(err)
	}

	newVer.Name = "-------"
	_, err = testEngine.ID(ver.Id).Update(context.Background(), newVer)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if newVer.Ver != 2 {
		err = errors.New("update should set version back to struct")
		t.Error(err)
	}

	newVer = new(VersionS)
	has, err = testEngine.ID(ver.Id).Get(context.Background(), newVer)
	if err != nil {
		t.Error(err)
		panic(err)
	}
	fmt.Println(newVer)
	if newVer.Ver != 2 {
		err = errors.New("insert error")
		t.Error(err)
		panic(err)
	}
}

func TestVersion2(t *testing.T) {
	assert.NoError(t, prepareEngine())

	err := testEngine.DropTables(context.Background(), new(VersionS))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = testEngine.CreateTables(context.Background(), new(VersionS))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	var vers = []VersionS{
		{Name: "sfsfdsfds"},
		{Name: "xxxxx"},
	}
	_, err = testEngine.Insert(context.Background(), vers)
	if err != nil {
		t.Error(err)
		panic(err)
	}

	fmt.Println(vers)

	for _, v := range vers {
		if v.Ver != 1 {
			err := errors.New("version should be 1")
			t.Error(err)
			panic(err)
		}
	}
}
