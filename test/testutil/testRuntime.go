package testutil

import (
	"io"
	"os"

	"git.sr.ht/~ionous/tapestry/rt"
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
	scope.Stack
	*Kinds
}

func (x *Runtime) Writer() io.Writer {
	return os.Stdout
}

// we ignore the request to initialize the scope during testing.
// that's only honored by the real thing.
func (x *Runtime) ReplaceScope(s rt.Scope, init bool) (ret rt.Scope, err error) {
	ret = x.Stack.ReplaceScope(s)
	return
}

func (x *Runtime) SetField(target, field string, value g.Value) (err error) {
	switch target {
	case meta.Variables:
		err = x.Stack.SetFieldByName(field, g.CopyValue(value))
	default:
		err = g.UnknownField(target, field)
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
			ret = g.StringsOf(a.Kind().Path())
		}

	case meta.Variables:
		ret, err = x.Stack.FieldByName(field)

	default:
		if a, ok := x.ObjectMap[target]; !ok {
			err = g.UnknownField(target, field)
		} else {
			ret, err = a.GetNamedField(field)
		}
	}
	return
}
