package list_test

import (
	"strconv"
	"strings"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/rt/test"
)

type panicTime struct {
	rt.Panic
}
type listTime struct {
	panicTime
	src, res g.Value
	objs     map[string]*g.Record
	scope.ScopeStack
	sort  *pattern.BoolPattern
	remap *pattern.ActivityPattern
	kinds *test.Kinds
}

func B(i bool) rt.BoolEval       { return &core.Bool{i} }
func I(i int) rt.NumberEval      { return &core.Number{float64(i)} }
func T(i string) rt.TextEval     { return &core.Text{i} }
func V(i string) core.ObjectEval { return &core.GetVar{Name: T(i)} }

func FromTs(vs []string) (ret core.Assignment) {
	if len(vs) == 1 {
		ret = &core.FromText{&core.Text{vs[0]}}
	} else {
		ret = &core.FromTextList{&core.Texts{vs}}
	}
	return
}

// cmd to collect some text into a list of strings.
type Write struct {
	out  *[]string
	Text rt.TextEval
}

// Execute writes text to the runtime's current writer.
func (op *Write) Execute(run rt.Runtime) (err error) {
	if t, e := op.Text.GetText(run); e != nil {
		err = e
	} else {
		(*op.out) = append((*op.out), t)
	}
	return
}

func getNum(run rt.Runtime, op rt.NumberEval) (ret string) {
	if v, e := op.GetNumber(run); e != nil {
		ret = e.Error()
	} else {
		ret = strconv.FormatFloat(v, 'g', -1, 64)
	}
	return
}

func joinText(run rt.Runtime, op rt.TextListEval) (ret string) {
	if vs, e := op.GetTextList(run); e != nil {
		ret = e.Error()
	} else {
		ret = joinStrings(vs)
	}
	return
}

func joinStrings(vs []string) (ret string) {
	if len(vs) > 0 {
		ret = strings.Join(vs, ", ")
	} else {
		ret = "-"
	}
	return
}

func (lt *listTime) GetField(target, field string) (ret g.Value, err error) {
	if target == object.Variables && field == "src" {
		ret = lt.src
	} else if target == object.Variables && field == "res" {
		ret = lt.res
	} else if obj, ok := lt.objs[field]; target == object.Value && ok {
		ret, err = g.ValueOf(obj)
	} else {
		ret, err = lt.ScopeStack.GetField(target, field)
	}
	return
}

func (lt *listTime) SetField(target, field string, value g.Value) (err error) {
	if target != object.Variables {
		err = lt.ScopeStack.SetField(target, field, value)
	} else {
		switch field {
		case "src":
			lt.src = value
		case "res":
			lt.res = value
		default:
			err = lt.ScopeStack.SetField(target, field, value)
		}
	}
	return
}

func (lt *listTime) GetEvalByName(name string, pv interface{}) (err error) {
	if name == "sort" {
		ptr := pv.(*pattern.BoolPattern)
		(*ptr) = *lt.sort
	} else if name == "remap" {
		ptr := pv.(*pattern.ActivityPattern)
		(*ptr) = *lt.remap
	}
	return
}

func (lt *listTime) GetKindByName(name string) (*g.Kind, error) {
	return lt.kinds.GetKindByName(name)
}
