package scope

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

// a scope which provides a single named field.
type SingleValue struct {
	name     string
	val      g.Value
	readOnly bool
}

//
func NewSingleValue(name string, val g.Value) *SingleValue {
	return &SingleValue{name: name, val: val, readOnly: false}
}

// creates a scope that provides a single named value;
// errors if the author tries to *replace* the value,
// noting that altering the *contents* of a value does not cause an error.
// ( ex. altering the fields of a record value or the elements of an array. )
func NewReadOnlyValue(name string, val g.Value) rt.Scope {
	return &SingleValue{name: name, val: val, readOnly: true}
}

// backdoor for updating the value of held by this scope.
// panics if the value was created read-only
func (k *SingleValue) SetValue(val g.Value) {
	if k.readOnly {
		panic("setting a readOnly value")
	}
	k.val = val
}

func (k *SingleValue) FieldByName(field string) (ret g.Value, err error) {
	if found := len(field) > 0 && field == k.name; !found {
		err = g.UnknownVariable(field)
	} else {
		ret = k.val
	}
	return
}

func (k *SingleValue) SetFieldByName(field string, val g.Value) (err error) {
	if found := len(field) > 0 && field == k.name; !found {
		err = g.UnknownVariable(field)
	} else if k.readOnly {
		err = errutil.New("the", k.name, "is read-only")
	} else {
		k.val = val
	}
	return
}

// placeholder method: for now, determines whether the field exists
func (k *SingleValue) SetFieldDirty(field string) (err error) {
	if found := len(field) > 0 && field == k.name; !found {
		err = g.UnknownVariable(field)
	} else if k.readOnly {
		err = errutil.New("the", k.name, "is read-only")
	}
	return
}
