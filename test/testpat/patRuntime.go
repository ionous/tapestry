package testpat

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/rt/scope"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
)

type Runtime struct {
	Map // pattern definitions
	testutil.Runtime
}

func (run *Runtime) GetField(object, field string) (ret rt.Value, err error) {
	if object != meta.PatternLabels {
		ret, err = run.Runtime.GetField(object, field)
	} else if p, ok := run.Map[field]; !ok {
		err = errutil.New("unknown pattern", field)
	} else {
		ret = rt.StringsOf(p.GetLabels())
	}
	return
}

// fix? follows from qna, but isnt an exact copy:
// improving the way inits work would probably help...
func (run *Runtime) Call(name string, aff affine.Affinity, keys []string, vals []rt.Value) (ret rt.Value, err error) {
	if kind, e := run.GetKindByName(name); e != nil {
		err = e
	} else if rec, e := pattern.InitRecord(run, kind, keys, vals); e != nil {
		err = e
	} else if field, e := pattern.GetResultField(run, kind); e != nil {
		err = e
	} else {
		newScope := scope.FromRecord(rec)
		oldScope := run.Chain.ReplaceScope(newScope)
		if rules, e := run.GetRules(name); e != nil {
			err = e
		} else if res, e := rules.Calls(run, newScope, field); e != nil {
			err = e
		} else {
			ret, err = res.GetResult(run, aff)
		}
		run.Chain.RestoreScope(oldScope)
	}
	return
}
