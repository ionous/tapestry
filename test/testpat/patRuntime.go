package testpat

import (
	"git.sr.ht/~ionous/tapestry/affine"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/scope"
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

// fix? follows from qna, but isnt an exact copy:
// improving the way inits work would probably help...
func (run *Runtime) Call(name string, aff affine.Affinity, keys []string, vals []g.Value) (ret g.Value, err error) {
	if kind, e := run.GetKindByName(name); e != nil {
		err = e
	} else if rec, e := safe.FillRecord(run, kind.NewRecord(), keys, vals); e != nil {
		err = e
	} else if field, e := pattern.GetResultField(run, kind); e != nil {
		err = e
	} else {
		oldScope := run.Chain.ReplaceScope(scope.FromRecord(run, rec))
		if rules, e := run.GetRules(name); e != nil {
			err = e
		} else if res, e := rules.Call(run, rec, field); e != nil {
			err = e
		} else {
			ret, err = res.GetResult(run, aff)
		}
		run.Chain.RestoreScope(oldScope)
	}
	return
}
