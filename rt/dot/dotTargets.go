package dot

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

type rootCursor struct {
	run  rt.Runtime
	name string
}

// supports requesting fields from the runtime
// where the name is either the id of an object or meta.Variables
func MakeCursor(run rt.Runtime, name string) Cursor {
	return rootCursor{run, name}
}

func (c rootCursor) CurrentValue() rt.Value {
	return rt.StringOf(c.name)
}

func (c rootCursor) SetAtIndex(int, rt.Value) error {
	return fmt.Errorf("%s can't be indexed", c.name)
}

func (c rootCursor) GetAtIndex(int) (_ Cursor, err error) {
	err = fmt.Errorf("%s can't be indexed", c.name)
	return
}

func (c rootCursor) GetAtField(field string) (ret Cursor, err error) {
	if v, e := c.run.GetField(c.name, field); e != nil {
		err = e
	} else {
		ret = valCursor{c.run, v}
	}
	return
}

func (c rootCursor) SetAtField(field string, val rt.Value) error {
	return c.run.SetField(c.name, field, val)
}

type valCursor struct {
	ks  rt.Kinds // targeted
	val rt.Value // possibly a record or a list
}

// supports requesting fields from a record
// uses the runtime to create sub-records when needed.
func MakeValueCursor(ks rt.Kinds, val rt.Value) Cursor {
	return valCursor{ks, val}
}

func (c valCursor) CurrentValue() rt.Value {
	return c.val
}

// errors if the current value isn't a list.
func (c valCursor) SetAtIndex(i int, newValue rt.Value) (err error) {
	if aff := c.val.Affinity(); !affine.IsList(aff) {
		err = fmt.Errorf("%s isn't a list", aff)
	} else if at, e := safe.Range(i, 0, c.val.Len()); e != nil {
		err = e
	} else if e := c.val.SetIndex(at, newValue); e != nil {
		err = e
	}
	return
}

// errors if the current value isn't a list.
func (c valCursor) GetAtIndex(i int) (ret Cursor, err error) {
	if aff := c.val.Affinity(); !affine.IsList(aff) {
		err = fmt.Errorf("%s isn't a list", aff)
	} else if i, e := safe.Range(i, 0, c.val.Len()); e != nil {
		err = e
	} else {
		next := c.val.Index(i)
		ret = valCursor{c.ks, next}
	}
	return
}

// copies the incoming value ( because varient SetFieldByName does )
// errors if the current value isn't a record.
func (c valCursor) SetAtField(field string, newValue rt.Value) (err error) {
	if aff := c.val.Affinity(); aff != affine.Record {
		err = fmt.Errorf("%s doesn't have fields", aff)
	} else {
		err = c.val.SetFieldByName(field, newValue)
	}
	return
}

// errors if the current value isn't a record.
func (c valCursor) GetAtField(field string) (ret Cursor, err error) {
	if aff := c.val.Affinity(); aff != affine.Record {
		err = fmt.Errorf("%s doesn't have fields", aff)
	} else if v, e := safe.EnsureField(c.ks, c.val.Record(), field); e != nil {
		err = e
	} else {
		ret = valCursor{c.ks, v}
	}
	return
}
