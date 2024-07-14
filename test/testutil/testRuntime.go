package testutil

import (
	"io"
	"os"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

type panicTime struct {
	PanicRuntime
}

// Runtime - a simple runtime for testing
type Runtime struct {
	panicTime
	Objects map[string]*rt.Record
	scope.Chain
	*Kinds
}

func (x *Runtime) Writer() io.Writer {
	return os.Stdout
}

func (x *Runtime) SetField(target, field string, value rt.Value) (err error) {
	switch target {
	case meta.Variables:
		err = x.Chain.SetFieldByName(field, rt.CopyValue(value))
	default:
		if a, ok := x.Objects[target]; !ok {
			err = rt.UnknownField(target, field)
		} else {
			err = a.SetNamedField(field, value)
		}
	}
	return
}

func (x *Runtime) GetField(target, field string) (ret rt.Value, err error) {
	switch target {
	case meta.ObjectId:
		if _, ok := x.Objects[field]; !ok {
			err = rt.UnknownObject(field)
		} else {
			// in the test runtime the name of the object is generally the same as id
			ret = rt.StringOf(field)
		}

	// return type of an object
	case meta.ObjectKind:
		if a, ok := x.Objects[field]; !ok {
			err = rt.UnknownObject(field)
		} else {
			ret = rt.StringOf(a.Name())
		}

		// hierarchy of an object's types ( a path )
	case meta.ObjectKinds:
		if a, ok := x.Objects[field]; !ok {
			err = rt.UnknownObject(field)
		} else {
			ret = rt.StringsOf(a.Ancestors())
		}

	case meta.Variables:
		ret, err = x.Chain.FieldByName(field)

	default:
		if a, ok := x.Objects[target]; !ok {
			err = rt.UnknownField(target, field)
		} else {
			ret, err = a.GetNamedField(field)
		}
	}
	return
}
