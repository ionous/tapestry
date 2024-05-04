package dot

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// supports requesting fields from the runtime
// where the top level is the id of the object or meta.Variables
type ObjValue struct {
	run  rt.Runtime
	name string
}

func NewTarget(run rt.Runtime, name string) Cursor {
	return ObjValue{run, name}
}

func (c ObjValue) CurrentValue() rt.Value {
	return rt.StringOf(c.name)
}

func (c ObjValue) SetAtIndex(int, rt.Value) (err error) {
	return fmt.Errorf("%s can't be indexed", c.name)
}

func (c ObjValue) GetAtIndex(int) (_ Cursor, err error) {
	err = fmt.Errorf("%s can't be indexed", c.name)
	return
}

func (c ObjValue) GetAtField(field string) (ret Cursor, err error) {
	if v, e := c.run.GetField(c.name, field); e != nil {
		err = e
	} else {
		ret = SubValue{c.run, v}
	}
	return
}

func (c ObjValue) SetAtField(field string, val rt.Value) error {
	return c.run.SetField(c.name, field, val)
}

type SubValue struct {
	run rt.Runtime
	val rt.Value
}

func (c SubValue) CurrentValue() rt.Value {
	return c.val
}

func (c SubValue) SetAtIndex(i int, newValue rt.Value) (err error) {
	if aff := c.val.Affinity(); !affine.IsList(aff) {
		err = fmt.Errorf("%s isn't a list", aff)
	} else if at, e := safe.Range(i, 0, c.val.Len()); e != nil {
		err = e
	} else if e := c.val.SetIndex(at, newValue); e != nil {
		err = e
	}
	return
}

func (c SubValue) GetAtIndex(i int) (ret Cursor, err error) {
	if aff := c.val.Affinity(); !affine.IsList(aff) {
		err = fmt.Errorf("%s isn't a list", aff)
	} else if i, e := safe.Range(i, 0, c.val.Len()); e != nil {
		err = e
	} else {
		next := c.val.Index(i)
		ret = SubValue{c.run, next}
	}
	return
}

func (c SubValue) SetAtField(field string, newValue rt.Value) (err error) {
	if aff := c.val.Affinity(); aff != affine.Record {
		err = fmt.Errorf("%s doesn't have fields", aff)
	} else if rec, ok := c.val.Record(); !ok {
		err = fmt.Errorf("cant set field %q into nil record", field)
	} else {
		err = rec.SetNamedField(field, newValue)
	}
	return
}

func (c SubValue) GetAtField(field string) (ret Cursor, err error) {
	if aff := c.val.Affinity(); aff != affine.Record {
		err = fmt.Errorf("%s doesn't have fields", aff)
	} else if rec, ok := c.val.Record(); !ok {
		err = fmt.Errorf("cant get field %q from nil record", field)
	} else if next, e := rec.GetNamedField(field); e != nil {
		err = e
	} else {
		ret = SubValue{c.run, next}
	}
	return
}
