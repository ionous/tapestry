package testutil

import (
	"git.sr.ht/~ionous/iffy/object"
	g "git.sr.ht/~ionous/iffy/rt/generic"
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
	PatternMap
	*Kinds
}

func (x *Runtime) Writer() writer.Output {
	return writer.NewStdout()
}

func (x *Runtime) SetField(target, field string, value g.Value) (err error) {
	switch target {
	case object.Variables:
		err = x.Stack.SetFieldByName(field, value)
	default:
		err = g.UnknownField(target, field)
	}
	return
}

func (x *Runtime) GetField(target, field string) (ret g.Value, err error) {
	switch target {
	case object.Variables:
		ret, err = x.Stack.FieldByName(field)
	case object.Value:
		if obj, ok := x.ObjectMap[field]; ok {
			ret = g.RecordOf(obj)
		} else {
			err = g.UnknownObject(field)
		}
	default:
		err = g.UnknownField(target, field)
	}
	return
}
