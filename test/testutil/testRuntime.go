package testutil

import (
	"io"
	"os"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/scope"
)

type panicTime struct {
	PanicRuntime
}

// Runtime - a simple runtime for testing
type Runtime struct {
	panicTime
	ObjectMap map[string]*g.Record
	scope.Chain
	*Kinds
}

func (x *Runtime) Writer() io.Writer {
	return os.Stdout
}

func (x *Runtime) SetField(target, field string, value g.Value) (err error) {
	switch target {
	case meta.Variables:
		err = x.Chain.SetFieldByName(field, g.CopyValue(value))
	case meta.ValueChanged:
		// unpack the real target and field
		switch target, field := field, value.String(); target {
		case meta.Variables:
			err = x.Chain.SetFieldDirty(field)
		default:
			// verify that the field exists...
			if _, ok := x.ObjectMap[field]; !ok {
				err = g.UnknownObject(field)
			}
		}
	default:
		if a, ok := x.ObjectMap[target]; !ok {
			err = g.UnknownField(target, field)
		} else {
			err = a.SetNamedField(field, value)
		}
	}
	return
}

func (x *Runtime) GetField(target, field string) (ret g.Value, err error) {
	switch target {
	case meta.ObjectId:
		if _, ok := x.ObjectMap[field]; !ok {
			err = g.UnknownObject(field)
		} else {
			// in the test runtime the name of the object is generally the same as id
			ret = g.StringOf(field)
		}

	// return type of an object
	case meta.ObjectKind:
		if a, ok := x.ObjectMap[field]; !ok {
			err = g.UnknownObject(field)
		} else {
			ret = g.StringOf(a.Type())
		}

		// hierarchy of an object's types ( a path )
	case meta.ObjectKinds:
		if a, ok := x.ObjectMap[field]; !ok {
			err = g.UnknownObject(field)
		} else {
			ret = g.StringsOf(g.Path(a.Kind()))
		}

	case meta.Variables:
		ret, err = x.Chain.FieldByName(field)

	default:
		if a, ok := x.ObjectMap[target]; !ok {
			err = g.UnknownField(target, field)
		} else {
			ret, err = a.GetNamedField(field)
		}
	}
	return
}
