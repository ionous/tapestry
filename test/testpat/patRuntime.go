package testpat

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
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

func (run *Runtime) Call(rec *g.Record, aff affine.Affinity) (ret g.Value, err error) {
	if res, e := pattern.NewResults(run, rec, aff); e != nil {
		err = e
	} else if oldScope, e := run.ReplaceScope(res, true); e != nil {
		err = e
	} else {
		var allFlags rt.Flags
		if rules, e := run.GetRules(rec.Kind().Name(), "", &allFlags); e != nil {
			err = e
		} else if e := res.ApplyRules(run, rules, allFlags); e != nil {
			err = e
		} else {
			ret, err = res.GetResult()
		}
		// only init can return an error
		run.ReplaceScope(oldScope, false)
	}
	return
}
