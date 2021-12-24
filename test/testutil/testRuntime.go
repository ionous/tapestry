package testutil

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/meta"
	"git.sr.ht/~ionous/iffy/rt/scope"
	"git.sr.ht/~ionous/iffy/rt/writer"
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

func (x *Runtime) Writer() writer.Output {
	return writer.NewStdout()
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
		err = x.Stack.SetFieldByName(field, value)
	default:
		err = g.UnknownField(target, field)
	}
	return
}

func (x *Runtime) GetField(target, field string) (ret g.Value, err error) {
	switch target {
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

	case meta.ObjectValue:
		if obj, ok := x.ObjectMap[field]; ok {
			ret = g.RecordOf(obj)
		} else {
			err = g.UnknownObject(field)
		}

	case meta.Variables:
		ret, err = x.Stack.FieldByName(field)

	default:
		err = g.UnknownField(target, field)
	}
	return
}
