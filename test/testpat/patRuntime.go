package testpat

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
)

type Runtime struct {
	Map // pattern definitions
	testutil.Runtime
}

func (run *Runtime) GetField(object, field string) (ret g.Value, err error) {
	if object != meta.PatternLabels {
		ret, err = run.Runtime.GetField(object, field)
	} else if p, ok := run.Map[field]; !ok {
		err = errutil.New("unknown pattern", field)
	} else {
		ret = g.StringsOf(p.GetLabels())
	}
	return
}

func (run *Runtime) Call(name string, aff affine.Affinity, keys []string, vals []g.Value) (ret g.Value, err error) {
	// create a container to hold results of args, locals, and the pending return value
	if kind, e := run.GetKindByName(name); e != nil {
		err = e
	} else if rec, e := safe.FillRecord(run, kind.NewRecord(), keys, vals); e != nil {
		err = e
	} else if labels, e := run.GetField(meta.PatternLabels, name); e != nil {
		err = e
	} else {
		var allFlags rt.Flags
		res := pattern.NewResults(rec, labels.Strings(), aff)
		oldScope := run.Stack.ReplaceScope(res)
		// ignores the initialization of locals during testing...
		if rules, e := run.GetRules(rec.Kind().Name(), "", &allFlags); e != nil {
			err = e
		} else if e := res.ApplyRules(run, rules, allFlags); e != nil {
			err = e
		} else {
			ret, err = res.GetResult()
			if !res.ComputedResult() {
				err = errutil.Fmt("%w calling %s test pattern %q", rt.NoResult, aff, rec.Kind().Name())
			}
		}
		run.ReplaceScope(oldScope)
	}
	return
}
