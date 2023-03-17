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
	// fix: this duplicates qnaCall....
	// maybe it could all be moved into package rt/pattern?
	// ( some tweak to handle "startedPattern/endedPattern" would be needed.
	//   some ideas: maybe be writing to a global, using SetField with a start/end key,
	//               checking replace scope for a patter results type )
	//  see also maybe "Send" -- which uses "BuildPath"
	if res, e := pattern.NewResults(run, rec, aff); e != nil {
		err = e
	} else {
		oldScope := run.Stack.ReplaceScope(res)
		// ignores the initialization of locals during testing...
		var allFlags rt.Flags
		if rules, e := run.GetRules(rec.Kind().Name(), "", &allFlags); e != nil {
			err = e
		} else if e := res.ApplyRules(run, rules, allFlags); e != nil {
			err = e
		} else {
			ret, err = res.GetResult()
			if !res.ComputedResult() {
				err = rt.NoResult{}
			}
		}
		run.ReplaceScope(oldScope)
	}
	return
}
