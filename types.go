package xorm

import (
	"reflect"

	"github.com/lingochamp/core"
)

var (
	ptrPkType = reflect.TypeOf(&core.PK{})
	pkType    = reflect.TypeOf(core.PK{})
)
