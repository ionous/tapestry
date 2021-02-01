package test

import (
	"git.sr.ht/~ionous/iffy/object"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/scope"
	"git.sr.ht/~ionous/iffy/rt/writer"
	"git.sr.ht/~ionous/iffy/test/testutil"
)

type panicTime struct {
	testutil.PanicRuntime
}

type testTime struct {
	panicTime
	objs map[string]*g.Record
	scope.ScopeStack
	testutil.PatternMap
	*testutil.Kinds
}

func (lt *testTime) Writer() writer.Output {
	return writer.NewStdout()
}

func (lt *testTime) GetField(target, field string) (ret g.Value, err error) {
	if obj, ok := lt.objs[field]; target == object.Value && ok {
		ret = g.RecordOf(obj)
	} else {
		ret, err = lt.ScopeStack.GetField(target, field)
	}
	return
}
